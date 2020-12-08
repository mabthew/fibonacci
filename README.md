# Fibonacci API
A fault-tolerant, high throughput API that steps through the Fibonacci sequence.

## Usage
This program is written in Go, so with Go installed, you can clone, build, and run with these commands.

```
git clone https://github.com/mabthew/fibonacci.git && 
    cd /fibonacci
    go get github.com/julienschmidt/httprouter && 
    go get github.com/golang/groupcache/lru && 
    go build . && 
    ./fibonacci
```


The server exposes the following endpoints:

* `GET /current`: returns the current number in the sequence
* `GET /next`: returns the next number in the sequence
* `GET /previous`: returns the previous number in the sequence


## Implementation Details
### Caching
I used the groupcache/lru package to implement an LRU cache storing a set of Fibonacci values surrounding the current index. My initial implementation was a map because I had caching in mind from the start, but once I found the groupcache/lru package I replaced the map and let it handle all the evictions and underlying logic. While simple addition is enough to find the previous or next number in the sequence, I chose LRU caching to ensure a numbers close to the current would be easily retreivable with no extra large number math.

### Recovery
To give the API the ability to recover from a crash, I kicked off a goroutine that writes the current index to disk every 2 seconds so that the user won't lose more than a couple seconds of calculations from an unexpected crash. Upon starting or restarting, the server attempts to read the backup file and if it finds an index, it builds the series up to that index before becoming available to new requests.

### Server
For the server implementation I chose to use HttpRouter. It's been almost a year since I used Go so instead of implementing my own I researched the top Go HTTP frameworks decided on HttpRouter for the simplicity of implementation, along with its high performance and small memory footprint. 

### big.Int
In order to reach later numbers in the Fibonacci series, I used big.Int imported by math/big that allow for calculations outside of the capabilities of Go's built-in integer types. 

## Performance

### Throughput
I tested throughput using `wrk` and found that starting with a clear cache and the current index at 0, the API was consistently able to handle well over 1,000 requests/second over a 2 minute duration with no issues. Due to the calculations at higher numbers, this rate drops over longer durations, but I was still seeing ~1,100 requests/seconds at a 5 minute duration.

**`/next`**

```
➜  fibonacci git:(main) ✗ wrk -t12 -c14 -d120s http://127.0.0.1:8080/next
Running 2m test @ http://127.0.0.1:8080/next
  12 threads and 14 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     4.85ms    4.22ms  86.22ms   76.13%
    Req/Sec   221.46    335.59     5.66k    95.29%
  317341 requests in 2.00m, 9.84GB read
Requests/sec:   2642.60
Transfer/sec:     83.92MB
```

**`/previous`**

To test the `/previous` endpoint, I set the backup file to 300,000 so there would be a lot of room to move towards 0. I ran the same command above against the `/previous` endpoint and the throughput was higher than the `/next` endpoint. This was most likely because the cache was populated with the indices being retrieved, resulting in a lot of cache hits and saving time on big number arithmetic.

```
➜  fibonacci git:(main) ✗ wrk -t12 -c14 -d120s http://127.0.0.1:8080/previous
Running 2m test @ http://127.0.0.1:8080/previous
  12 threads and 14 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     2.39ms    3.11ms  63.98ms   84.40%
    Req/Sec   664.63      1.65k   10.65k    92.94%
  950537 requests in 2.00m, 9.91GB read
Requests/sec:   7914.98
Transfer/sec:     84.51MB
```

**Random load testing**

The  order to test I tried several different load testing tools- Apache JMeter, Apache Bench, Vegeta, and http_load. I ran into issues with each, but the http_load results seemed the most reliable, so I'll address those. http_load was showing ~1000 requests/second for a 5s duration but would drastically drop to ~500 requests/second at durations over a minute, and from what I could see it didn't seem to be an issue with the API. I don't trust this metric because the caching along with an even distribution of previous, current, and next calls should have resulted in mostly cache hits. I would expect a test to perform with a similar or higher throughput than the previous tests.

### Performance on a small machine
To test this on a small machine, I spun up a docker container allowing 1 CPU and 512MB RAM. I've included the dockerfile and to run it can be run using these commands.

```
docker build .  
docker run -p 8080:8080 -m 512m --cpus <image_id>

```

The performance was slightly lower than on my machine but still over 1,000 requests/second.
```
➜  fibonacci git:(main) ✗ wrk -t12 -c14 -d60s http://127.0.0.1:8080/next
Running 1m test @ http://127.0.0.1:8080/next
  12 threads and 14 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    19.46ms   23.09ms 191.49ms   80.19%
    Req/Sec    86.41     47.19   350.00     78.91%
  62043 requests in 1.00m, 391.59MB read
Requests/sec:   1032.18
Transfer/sec:      6.51MB
```