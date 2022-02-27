# Map Reduce

This demonstrates how map reduce can be run for a 
- Synchronous path
- Asynchronous with channels
-  Asynchronous with an external queue and worker instances on containers

Live stream recording can be found here: https://t.co/O7IMYcxRgO

The example task is word count, a popular one which showcases the efficacy of distributed computing in a simple sense.

## Methods
All methods are tagged with their respective functions
 - synchronous: Utilizes no concurrency
 - asynchronous: Utilizes native concurrency with channels
 - asynchronous_with_external_queue: Utilizes concurrent processing with individual instances and queues for scheduling tasks.
 
## Running 
Methods 1 and 2 can be done with:

```bash
go run *.go
```

Method 3 with queues and containers can be executed with:
```bash
docker build -t map_reduce .
docker-compose up
```

Wait until the master finishes then Ctrl-C to close.
