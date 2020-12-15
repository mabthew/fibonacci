# Fibonacci API
A fault-tolerant, high-throughput API that steps through the Fibonacci sequence.

## Usage
This program is written in Go, so with Go installed, you can clone, build, and run with these commands.

```
git clone https://github.com/mabthew/fibonacci.git && 
    cd fibonacci/ &&
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
The groupcache/lru package was used to implement an LRU cache storing a set of Fibonacci values surrounding the current index. The initial implementation was a map with caching in mind from the start, but changed in favor of the groupcache/lru package to simplify evictions and underlying logic. While simple addition is enough to find the previous or next number in the sequence, LRU caching ensures numbers close to the current index are easily retreivable with no big number arithmetic.

### Recovery
To give the API the ability to recover from a crash, a goroutine that writes the current index to disk every 2 seconds so that the user won't lose more than a couple seconds of calculations from an unexpected crash. Upon starting or restarting, the server attempts to read the backup file and if it finds an index, it builds the series up to that index before becoming available to new requests.

### Server
 HttpRouter is used for the server implementation due to its high performance and small memory footprint. 

### big.Int
The math/big package is used because the big.Int type allows for calculations outside of the capabilities of Go's built-in integer types. 

## Performance
### Throughput
The API was tested using `wrk` with an empty cache and a starting index of 0. It was consistently able to handle ~2,600 requests/second over a 2 minute duration with no issues. Due to the more time-intensive calculations at higher numbers, this rate drops over longer durations. However, the API was still capable of processing ~1,100 requests/seconds at a 5 minute duration.

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

The `/previous` endpoint was tested with a starting index of 300,000 so there would be a lot of room to move towards 0. Running the same `wrk` command from above against the `/previous` endpoint showed that the throughput was higher than the `/next` endpoint. This was most likely because the cache was populated with the indices being retrieved, resulting in a lot of cache hits and saving time on big number arithmetic.

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

**Random load testing**<br/>
Several different load testing tools were implemented to see how a randomized program would perform- Apache JMeter, Apache Bench, Vegeta, and http_load. There were issues with each, but the http_load results seemed the most reliable. http_load was showing ~1000 requests/second for a 5s duration but would drastically drop to ~500 requests/second at durations over a minute, and from all the testing performed it didn't seem to be an issue with the API. I don't trust this metric because the caching along with an even distribution of previous, current, and next calls should have resulted in mostly cache hits. I would expect a test to perform with a similar or higher throughput than the previous tests.

### Performance on a small machine
To test this on a small machine, I spun up a docker container allowing 1 CPU and 512MB RAM. The docker container can be built and run using the following commands.

```
docker build .  
docker run -p 8080:8080 -m 512m --cpus 1 <image_id>
```

The performance was slightly lower than on my machine, but still over 1,000 requests/second.
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

## Testing 
Functional testing is available in the `fibonacci_test.go` file, and run using the command
```
go test -v
```
