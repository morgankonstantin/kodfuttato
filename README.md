# Kódfuttató

Go nyelvű gRPC mikroszolgáltatás, amely a kapott Python kódot futtatja
(5 mp-es időkorláttal), és visszaadja a kimenetet. A kód a folyamat standard bemenetén érkezik,
így teljes, többsoros szkriptek is küldhetők, méretkorlát nélkül.

## Verziók

| Git tag  | Ág                   | Viselkedés                                                        |
|----------|----------------------|-------------------------------------------------------------------|
| `v1.0.0` | `main`               | Stabil verzió.                                                    |
| `v1.1.0` | `fault/memory-spike` | Kérésenként 50 MB memóriát foglal és 30 mp-ig megtartja.          |
| `v1.2.0` | `fault/latency`      | Minden kérés előtt 300–1500 ms véletlenszerű késleltetés.         |

A két hibás verzió a `v1.0.0`-ból ágazik le, és csak az `app/executor.go`-ban tér el tőle.

## API

`executor.v1.CodeExecutor/Execute`, alapértelmezett port: `50051` (a `PORT` környezeti változóval módosítható).
A szerződés: `app/proto/executor.proto`. A szerveren be van kapcsolva a gRPC reflection.

## Előfeltételek

- Go 1.22+ és `python3` (a lokális futtatáshoz és a tesztekhez)
- Docker
- `protoc`, `protoc-gen-go v1.35.1`, `protoc-gen-go-grpc 1.5.1`

## Lokális fejlesztés

```bash
cd app
go test ./...
go run .
```

A unit tesztek az `app/testdata/` mappában lévő `.py` szkripteket futtatják le, és az eredményt
az `executor_test.go` táblázatában megadott értékekkel vetik össze. Új teszteset: egy új `.py` fájl
a `testdata` mappában és egy új sor a táblázatban.

Kézi hívás (ha a `grpcurl` telepítve van):

```bash
grpcurl -plaintext -d '{"code":"print(6*7)"}' localhost:50051 executor.v1.CodeExecutor/Execute
```

## Proto generálása

```bash
cd app
protoc --go_out=. --go_opt=module=code-executor \
       --go-grpc_out=. --go-grpc_opt=module=code-executor \
       proto/executor.proto
```

## Konténerizálás

Futtatás:

```bash
docker run --rm -p 50051:50051 code-executor:v1.0.0
```
