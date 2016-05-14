
lightproxy is a lightweight HTTP reverse proxy written in Golang aims
to accelerate client requesting. It caches web resources in main memory and 
implements a non-blocking design, so providing lightning fast speed for you.

lightproxy takes a rather conservative strategy in deciding whether 
to cache a certain resources, rules includes:

1. requests of types other than GET are never cached
2. URLs containing query strings are not cached
3. if Set-Cookie was found in response's header, it is not cached

Apache Bench testing
-------------

On Windows 7 Professional,

######1. start a web server first:
> godoc -http :9009

######2. run our proxy program:
> lightproxy.exe

######3. run benchmark against backend server:
> ab -c 100 -n 10000 http://127.0.0.1:9009/pkg/

outputs:
<pre><code>This is ApacheBench, Version 2.3 <$Revision: 655654 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking 127.0.0.1 (be patient)
Completed 1000 requests
Completed 2000 requests
Completed 3000 requests
Completed 4000 requests
Completed 5000 requests
Completed 6000 requests
Completed 7000 requests
Completed 8000 requests
Completed 9000 requests
Completed 10000 requests
Finished 10000 requests

Server Software:
Server Hostname:        127.0.0.1
Server Port:            9009

Document Path:          /pkg/
Document Length:        45904 bytes

Concurrency Level:      100
Time taken for tests:   34.586 seconds
Complete requests:      10000
Failed requests:        0
Write errors:           0
Total transferred:      460000000 bytes
HTML transferred:       459040000 bytes
Requests per second:    289.13 [#/sec] (mean)
Time per request:       345.860 [ms] (mean)
Time per request:       3.459 [ms] (mean, across all concurrent requests)
Transfer rate:          12988.46 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    0   0.5      0       5
Processing:    13  344  83.2    316     949
Waiting:        7  338  81.1    311     948
Total:         13  344  83.2    316     949

Percentage of the requests served within a certain time (ms)
  50%    316
  66%    337
  75%    361
  80%    381
  90%    443
  95%    493
  98%    609
  99%    724
 100%    949 (longest request)</code></pre>

######4. run benchmark against the proxy program:
> ab -c 100 -n 10000 http://127.0.0.1:9090/pkg

outputs:
<pre><code>This is ApacheBench, Version 2.3 <$Revision: 655654 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking 127.0.0.1 (be patient)
Completed 1000 requests
Completed 2000 requests
Completed 3000 requests
Completed 4000 requests
Completed 5000 requests
Completed 6000 requests
Completed 7000 requests
Completed 8000 requests
Completed 9000 requests
Completed 10000 requests
Finished 10000 requests


Server Software:
Server Hostname:        127.0.0.1
Server Port:            9090

Document Path:          /pkg/
Document Length:        45904 bytes

Concurrency Level:      100
Time taken for tests:   10.090 seconds
Complete requests:      10000
Failed requests:        0
Write errors:           0
Total transferred:      460000000 bytes
HTML transferred:       459040000 bytes
Requests per second:    991.12 [#/sec] (mean)
Time per request:       100.896 [ms] (mean)
Time per request:       1.009 [ms] (mean, across all concurrent requests)
Transfer rate:          44523.05 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    0   0.6      0      14
Processing:    20  100  34.1     98     300
Waiting:        4   38  28.0     28     236
Total:         20  100  34.2     98     301

Percentage of the requests served within a certain time (ms)
  50%     98
  66%    113
  75%    120
  80%    125
  90%    142
  95%    163
  98%    187
  99%    193
 100%    301 (longest request)
</code></pre>
