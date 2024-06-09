# Macbook Pro: M2 Pro, 32GB

````bash
wrk http://127.0.0.1:8080/move --latency -t1 -c32 -d10s -s ./scripts/wrk.lua

State: map[Games:10000 Games Complete:0 Moves:0 asfd:0]
Running 10s test @ http://127.0.0.1:8080/move
  1 threads and 32 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   393.92us  132.63us   5.16ms   84.34%
    Req/Sec    80.21k    13.44k  128.35k    75.00%
  Latency Distribution
     50%  387.00us
     75%  436.00us
     90%  501.00us
     99%  708.00us
  796750 requests in 10.00s, 98.02MB read
Requests/sec:  79673.40
Transfer/sec:      9.80MB
State: map[Games:20000 Games Complete:10000 Moves:796782 asfd:13]
signal: terminated```
````
