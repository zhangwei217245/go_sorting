## Running and Inspecting

### Running the Program

```bash
docker run -m 2g -d \
             --name go_sorting \
             -p 9092:9092 \
             -e KAFKA_ADVERTISED_HOST_NAME=localhost \
             -e KAFKA_CREATE_TOPICS="source:4:1" \
             zhangwei217245/go_sorting:latest
```
### Inspecting the logs

To see the go_sorting program was run by the shell script

```bash
cat /nohup.out
```

To see the logs of the go_sorting program:

```bash
cd /tmp/go_sorting
tail -f /tmp/go_sorting_*.log
```

We have file different types of logs, such as `info`, `warn`, `error`, `debug`, and `timing`.

The `timing` log shows the detailed timing. 

### Viewing the source code

Please see `/app` directory.

## Design Considerations

### ID Generator

To support generating IDs from multiple threads, I reserve 
the higher-order 5 bits for thread ID. With the rest 27 bits, each thread will have sufficient ID space for 100 million IDs. 
On a int32 basis, there should be 2^31 available IDs. 

I currently do not consider machine IDs, Unix timestamp etc. in our design, as for such a demo program running on a single machine, there is no need to pursue additional complexity. 

### Kconnect and Kafka.

Kconnect module basically contains a set of functions that implements producer and consumer functionalities. 

As we have 4 cores, so we set the number of partitions to be 4. Then, for each partition, there will be a goroutine working on it to retrieve the message. 

The program was written to be scalable, so when you have more partitions, there will be more goroutines.

### Memory Restriction

As there are only 2GB memory, we consider running Kafka with 4 topics (1 topic featuring 2 partitions, others with only 1 partition), and during the test, I set a limitation of 800MB to watch over the performance, it was fine. So I consider roughly 800MB for Kafka. 

For the data processor program, as it has to sort 100 million rows according to different keys using less than 1GB memory space, therefore I use external sorting. I consider dividing the whole source message stream into `n`
chunks so that each chunk will have `c/n` messages, given `c` as the total number of messages. By default, `n=10`, and buffering `10,000,000` messages will only take roughly 495MB. Whenever the buffer is full, I sort the buffer according to three specified sorting criteria respectively and send the sorted data into different files. 

Finally, I merge the sorted result with the help of minHeap. For each sorting criteria, I open `n` files (30 opened files in total in my default case), and I merge the sorted result with a dedicated minHeap for the particular sorting criteria. Whenever the minHeap is full of `n` elements, I pop-up the top element and send it to the corresponding Kafka topic, and I also read the next record from the file where the pushed elements came from.

When merging the sorted result, as I setup 3 goroutines, there will be 3 minHeaps with size `n`, therefore, it will only take O(n) space where n is the number of chunks for the entire message stream.

## Optimizations

For optimizations, I did the following.

### Channel buffer size of the consumer.

I setup the buffer size of the consumer channel as 256. This will allow 256 messages to be buffered in memory and hence will reduce the overhead that goroutines get blocked by the zero-sized channel.
This will also allow a lower memory space consumption, as we do not want a large piece of memory to be buffered. 

### Possible optimizations

There could be chances that I can increase the number of file descriptors that any one process may open. This will allow more sorted files to be created for external sorting and hence will boost up the parallelism. 
In fact, I did not do so because, we only have 4 cores and limited memory of 2GB, generating a larger number of sorted intermediate files can also be a challenge for such a limited hardware configuration.

## Difficulties and efforts

For me, it is the first time to write a program using Golang. Previously, I had no experience with Golang. But this is an awesome programming language, easy to use and easy to understand. The most beautiful feature is the concurrency model and "sharing memory by communication" is really something wonderful but also something a Java programmer has to adapt himself to.

If given more time, I'd rather study a little bit more about its concurrency model and different ways of using channels. It was a lot of fun.

Also, for this mini-project, I also did some study on Kafka. I had some experience with RabbitMQ, so Kafka is not a big problem for me. But to achieve the best I can, I still did some study on selecting the best Golang driver for kafka, and finally I choose Shopify/sarama. It is said to be good for both performance and developer-friendliness. As compared to another high-performance driver - confluence-kafka driver, it does not rely on the installation of dynamic link library.

I understand I already exceeded the 9-day time limitation and I actually take 2 more days. Given that I have no experience with Golang and Kafka and that I also have heavy workload preparing my paper publication and other job interviews, I think I did my best. 

But doing the project just made me learn and I appreciate the opportunity you guys provided for me to be able to enjoy such a wonderful process. I'm not expecting much, but I would like to get your feedback, especially on where I can improve. Thanks!

