# Missing Semester: Sharded Map

## Description

This project aims to develop an **HTTP API** that mimics the functionality of a **HashMap** while overcoming the constraint of data distribution across multiple machines. The primary challenge lies in designing and implementing a **sharding scheme** to distribute and manage data efficiently across a cluster of machines.

## Features
***HashMap Functionality*** : _The ```API``` provides basic ```HashMap``` operations including inserting, retrieving, and deleting key-value pairs._

***Sharding*** : _Data is distributed across multiple maps that mimic the behavior of machines using a sharding scheme, ensuring load balancing (Consistent hashing) and fault tolerance._

***HTTP API*** : _Utilizes ```HTTP``` endpoints for communication, making it accessible for use with different apps._

***Concurrency*** : _Handles ``concurrent`` read and write operations safely to maintain data consistency._

## Implementation

***I chose Golang because it's a good opportunity to learn this language, and also I am intersted in the way how this programming language support concurrency.***

### Golang Sync package : 

The implementation of the project include the sync package as the most used package by using ``sync.atomic`` that provides methods of achieving atomic operations. Also ``sync.Mutex`` to prevent race conditions when accessing data concurrently from different goroutines.

The ``sync`` package provides a special type of maps which is ``sync.Map``. In Golang docs ``sync.Map`` is considered as special map that can handle concurrent access to data within the map itself by using ``atomic`` operations.

I had to choose between two different implementations : 

1. The first implementation is based on the use of a normal map with ``sync.RWMutex``.

2. The second implementation use ``sync.map`` which by default includes features like ``sync.atomic`` to handle concurrency.

To choose the implementation, I refered to the documentation of ``sync.Map``:

***The Map type is optimized for two common use cases:***

>_(1) when the entry for a given key is only ever written once but read many times, as in caches that only grow, in caches systems_.

>_(2) when multiple goroutines read, write, and overwrite entries for disjoint sets of keys_.

I chose the _2nd_ proposition since it's compatible with our use case.

Also, using ``atomic operations`` are better in terms of performance compared to
``Mutexes``.

Here is a simple benchmark of how atomic operations perform better than ``Mutexes``:

```
pkg: missing_semester/sharded_maps
cpu: AMD Ryzen 5 PRO 5650U with Radeon Graphics     
BenchmarkMutexAdd-12     	124022059	         9.673 ns/op
BenchmarkAtomicAdd-12    	289258807	         4.434 ns/op
PASS
ok  	missing_semester/sharded_maps	4.184s
```

> You can find the _source code_ of this benchmark test in ``benchmark_test.go`` file.

### Sharded Map structure:

1. A sharded map is composed of multiple ```sync.map```.
2. Each shard within the sharded map is value of specific key.
3. To distribute the traffic across the sharded map, I used ***consistent hashing*** to identify
the next shard we will direct the traffic to.
4. I used a hash algorithm to hash the key of the key-value to insert, then find which of our shards will handle the operation.

## How to run the project:
You can run the project by using this command:
```
go run .
```

**Note**: _the application server runs on port 8080, which can be modified within the main file_.

**There is no external dependency needed for the project since the api is built using the golang http package from the standard library.**

In order to build the project to a single executable file:
```
go build
```

## Testing:

*For each endpoint operation (post, get, delete), I'm only able to have ~210 concurrent requests locally because of hardware limitations.*.

***Even, using fastHttp that is already implemented in go Fiber framework, the program was limited to the same number of concurrent request which is ~210 requests.***

***Each endpoint of the api has its own test file***

***For example,*** ```api_post_test.go```  ***to test*** ```POST /api <key>: <value>```.

```
Server listening on :8080
POST request handled successfully for key 'key95'
POST request handled successfully for key 'key143'
POST request handled successfully for key 'key200'
POST request handled successfully for key 'key142'
POST request handled successfully for key 'key4'  
POST request handled successfully for key 'key90'
...
```
