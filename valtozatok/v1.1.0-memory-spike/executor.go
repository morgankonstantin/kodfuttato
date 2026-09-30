package main

import (
	"bytes"
	"context"
	"errors"
	"log"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"code-executor/executorpb"
)

// v2 hibája: minden kérés lefoglal spikeSize memóriát, és spikeHold ideig megtartja.
// Terhelés alatt a lefoglalt memória összeadódik, így a konténer memóriahasználata felszalad.
const (
	spikeSize = 50 << 20 // 50 MB
	spikeHold = 30 * time.Second
)

// execTimeout a Python futás felső időkorlátja. Változó, hogy a teszt rövidebbre állíthassa.
var execTimeout = 5 * time.Second

type server struct {
	executorpb.UnimplementedCodeExecutorServer
}

func (s *server) Execute(ctx context.Context, req *executorpb.ExecuteRequest) (*executorpb.ExecuteResponse, error) {
	if req.GetCode() == "" {
		return nil, status.Error(codes.InvalidArgument, "a code mező nem lehet üres")
	}

	memorySpike()

	res, err := runPython(ctx, req.GetCode())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "a futtatás nem indult el: %v", err)
	}

	log.Printf("execute: exit=%d duration=%dms", res.ExitCode, res.DurationMs)
	res.Version = version
	return res, nil
}

// runPython lefuttatja a kódot egy python3 folyamatban. A kód a standard bemeneten érkezik
// ("python3 -"), így a hosszát nem korlátozza a parancssori argumentumok mérethatára.
// Hibát csak akkor ad vissza, ha a folyamat el sem indult; a Python oldali hibák
// a kilépési kódban és a stderr-ben jelennek meg.
func runPython(ctx context.Context, code string) (*executorpb.ExecuteResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, execTimeout)
	defer cancel()

	var stdout, stderr bytes.Buffer
	cmd := exec.CommandContext(ctx, "python3", "-")
	cmd.Stdin = strings.NewReader(code)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	start := time.Now()
	err := cmd.Run()
	res := &executorpb.ExecuteResponse{DurationMs: time.Since(start).Milliseconds()}

	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		res.Stdout = stdout.String()
		res.Stderr = stderr.String() + "execution timed out\n"
		res.ExitCode = -1
		return res, nil
	}

	var exitErr *exec.ExitError
	if err != nil && !errors.As(err, &exitErr) {
		return nil, err
	}

	res.Stdout = stdout.String()
	res.Stderr = stderr.String()
	res.ExitCode = int32(cmd.ProcessState.ExitCode())
	return res, nil
}

// memorySpike lefoglal spikeSize bájtot, és spikeHold ideig életben tartja.
func memorySpike() {
	buf := make([]byte, spikeSize)
	for i := 0; i < len(buf); i += 4096 {
		buf[i] = 1 // minden lapot érintünk, hogy a memória ténylegesen lefoglalódjon
	}
	go func() {
		time.Sleep(spikeHold)
		runtime.KeepAlive(buf)
	}()
}
