# Macbook Pro: M2 Pro, 32GB

Wrk 16 connections:

```bash
➜  go-fastcounters git:(main) ✗ wrk http://127.0.0.1:8080 --latency -d10s -t1  -c16
Running 10s test @ http://127.0.0.1:8080
  1 threads and 16 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   130.11us   70.56us   4.13ms   88.77%
    Req/Sec    77.06k     5.59k   85.78k    58.42%
  Latency Distribution
     50%  112.00us
     75%  147.00us
     90%  202.00us
     99%  300.00us
  773941 requests in 10.10s, 95.95MB read
Requests/sec:  76624.75
Transfer/sec:      9.50MB
```

Wrk 8 threads, 16 connections:

```bash
➜  go-fastcounters git:(main) ✗ wrk http://127.0.0.1:8080 --latency -d20s -t8 -c 16
Running 20s test @ http://127.0.0.1:8080
  8 threads and 16 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   138.11us   73.05us   5.64ms   95.06%
    Req/Sec     8.61k     2.64k   15.75k    83.93%
  Latency Distribution
     50%  130.00us
     75%  144.00us
     90%  165.00us
     99%  290.00us
  714587 requests in 20.10s, 88.59MB read
Requests/sec:  35551.57
Transfer/sec:      4.41MB
```

```bash
➜  go-fastcounters git:(main) ✗ wrk http://127.0.0.1:8080 --latency -d10s -t8 -c 16
Running 10s test @ http://127.0.0.1:8080
  8 threads and 16 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   139.71us   67.08us   4.16ms   92.60%
    Req/Sec     9.81k     3.42k   16.25k    68.94%
  Latency Distribution
     50%  129.00us
     75%  144.00us
     90%  177.00us
     99%  302.00us
  471798 requests in 10.10s, 58.49MB read
Requests/sec:  46698.54
Transfer/sec:      5.79MB
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
