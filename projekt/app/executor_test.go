package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"code-executor/executorpb"
)

// A testdata mappa szkriptjeit futtatja le, és összeveti az eredményt a várttal.
// Új eset: egy új .py fájl a testdata mappában és egy új sor a táblázatban.
func TestExecuteScripts(t *testing.T) {
	old := execTimeout
	execTimeout = time.Second // az infinite_loop.py miatt ne kelljen 5 mp-et várni
	defer func() { execTimeout = old }()

	tests := []struct {
		file           string
		wantExit       int32
		wantStdout     string // pontos egyezés
		stderrContains string // üres: nem ellenőrizzük
	}{
		{
			file:       "fibonacci.py",
			wantExit:   0,
			wantStdout: "0: 0\n1: 1\n2: 1\n3: 2\n4: 3\n5: 5\n6: 8\n7: 13\n8: 21\n9: 34\n",
		},
		{
			file:       "stdlib.py",
			wantExit:   0,
			wantStdout: `{"most_common": "alma", "sqrt_2": 1.4142, "unique": 3}` + "\n",
		},
		{
			file:           "exception.py",
			wantExit:       1,
			wantStdout:     "a hiba előtti kimenet még megjelenik\n",
			stderrContains: "ZeroDivisionError",
		},
		{
			file:           "infinite_loop.py",
			wantExit:       -1,
			wantStdout:     "",
			stderrContains: "execution timed out",
		},
	}

	for _, tc := range tests {
		t.Run(tc.file, func(t *testing.T) {
			code, err := os.ReadFile(filepath.Join("testdata", tc.file))
			if err != nil {
				t.Fatalf("nem sikerült beolvasni: %v", err)
			}

			res, err := (&server{}).Execute(context.Background(), &executorpb.ExecuteRequest{Code: string(code)})
			if err != nil {
				t.Fatalf("váratlan hiba: %v", err)
			}
			if res.ExitCode != tc.wantExit {
				t.Errorf("kilépési kód: kaptam %d, vártam %d (stderr: %q)", res.ExitCode, tc.wantExit, res.Stderr)
			}
			if res.Stdout != tc.wantStdout {
				t.Errorf("stdout:\nkaptam: %q\nvártam: %q", res.Stdout, tc.wantStdout)
			}
			if tc.stderrContains != "" && !strings.Contains(res.Stderr, tc.stderrContains) {
				t.Errorf("a stderr nem tartalmazza ezt: %q\nstderr: %q", tc.stderrContains, res.Stderr)
			}
			if res.Version != version {
				t.Errorf("verzió: kaptam %q, vártam %q", res.Version, version)
			}
		})
	}
}

func TestExecuteEmptyCode(t *testing.T) {
	_, err := (&server{}).Execute(context.Background(), &executorpb.ExecuteRequest{Code: ""})
	if status.Code(err) != codes.InvalidArgument {
		t.Fatalf("InvalidArgument hibát vártam, kaptam: %v", err)
	}
}
