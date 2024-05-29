# Macbook Pro: M2 Pro, 32GB

Wrk 16 connections:

```bash
➜  go-fastcounters git:(main) wrk http://127.0.0.1:8080/ --latency -t8 -c64 -d60s
Running 1m test @ http://127.0.0.1:8080/
  8 threads and 64 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   390.94us    1.76ms  87.51ms   99.39%
    Req/Sec    25.35k     3.88k   57.20k    84.54%
  Latency Distribution
     50%  281.00us
     75%  357.00us
     90%  454.00us
     99%    1.37ms
  12114742 requests in 1.00m, 1.44GB read
Requests/sec: 201726.14
Transfer/sec:     24.62MB
```

```bash
➜  go-fastcounters git:(main) wrk http://127.0.0.1:8080/ --latency -t8 -c16 -d60s
Running 1m test @ http://127.0.0.1:8080/
  8 threads and 16 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     6.37ms   49.80ms 627.86ms   98.09%
    Req/Sec    14.82k     1.65k   16.44k    92.11%
  Latency Distribution
     50%  122.00us
     75%  147.00us
     90%  217.00us
     99%  303.44ms
  6948928 requests in 1.00m, 848.26MB read
Requests/sec: 115738.14
Transfer/sec:     14.13MB
```

Skip http:

```bash
➜  go-fastcounters git:(main) ✗ go test . -run Bench -v
=== RUN   TestBench
             Read WAL in:           5 seconds
             Read WAL at:   1,955,352 msgs/sec
 Commands read from file:  10,011,117
    Next sequence number:  10,011,118

Elapsed time:  17.570290125s
Commands processed:  10000000
Commands per second:  569142.0716129699
--- PASS: TestBench (22.69s)
PASS
ok  	github.com/mj1618/go-fastcounters	23.149s
```
