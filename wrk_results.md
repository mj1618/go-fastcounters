# Macbook Pro: M2 Pro, 32GB

Wrk 16 connections:

```bash
➜  go-fastcounters git:(main) ✗ wrk http://127.0.0.1:8080 --latency -d10s -t4 -c64
Running 10s test @ http://127.0.0.1:8080
  4 threads and 64 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   147.62us  110.31us   4.17ms   93.17%
    Req/Sec    30.95k    22.60k   76.90k    53.31%
  Latency Distribution
     50%  121.00us
     75%  141.00us
     90%  219.00us
     99%  599.00us
  745164 requests in 10.10s, 90.96MB read
Requests/sec:  73777.36
Transfer/sec:      9.01MB
```

```bash
➜  go-fastcounters git:(main) ✗ wrk http://127.0.0.1:8080 --latency -d10s -t4 -c64
Running 10s test @ http://127.0.0.1:8080
  4 threads and 64 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   175.38us  129.82us   7.82ms   91.73%
    Req/Sec    31.21k    10.05k   71.49k    79.77%
  Latency Distribution
     50%  134.00us
     75%  213.00us
     90%  289.00us
     99%  609.00us
  1089533 requests in 10.10s, 133.00MB read
Requests/sec: 107870.18
Transfer/sec:     13.17MB
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
