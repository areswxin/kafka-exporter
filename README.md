# kafka-exporter

kakfa0.8版本Prometheus exporter，通过zookeeper获取offset，输出Metrrics如下:
 - `kafka_consumergroup_current_offset` kafka Offset 值
 - `kafka_consumergroup_logsize`  kafka logsize 值
 - `kafka_consumergroup_lag` 近似的kafka lag值

参考[kafka-prometheus-exporter](https://github.com/krallistic/kafka-prometheus-exporter)修改

# Build 
Build go binary `go build -o kafka-exporter`
 
# Usage
`
./kafka-exporter -command
`

# Command line Options:

| Argument | Description | Default |
| --- | --- | --- |
| `listen-address` | The address on which to expose the web interface and generated Prometheus metrics. | `:8080`
| `telemetry-path` | Path under which to expose metrics. | `/metrics`
| `zookeeper-connect` | Zookeeper connection string | `localhost:2181`
| `cluster-name` | Name of the Kafka cluster used in static label |`kafka-cluster` 
| `refresh-interval` | Seconds to sleep in between refreshes | `15`

### Example Output:
```yaml
# HELP kafka_consumergroup_current_offset Current Offset of a ConsumerGroup at Topic/Partition
# TYPE kafka_consumergroup_current_offset gauge
kafka_consumergroup_current_offset{cluster="kafka-cluster",consumergroup="test",partition="0",topic="test"} 1645
kafka_consumergroup_current_offset{cluster="kafka-cluster",consumergroup="test",partition="1",topic="test"} 1578
....
# HELP kafka_consumergroup_lag Current approximate Lag of a ConsumerGroup at Topic/Partition
# TYPE kafka_consumergroup_lag gauge
kafka_consumergroup_lag{cluster="kafka-cluster",consumergroup="test",partition="0",topic="test"} 6.64292682e+08
kafka_consumergroup_lag{cluster="kafka-cluster",consumergroup="test",partition="1",topic="test"} 8.1979849e+08
...
# HELP kafka_consumergroup_logsize Current logsize of a ConsumerGroup at Topic/Partition
# TYPE kafka_consumergroup_logsize gauge
kafka_consumergroup_logsize{cluster="kafka-cluster",consumergroup="test",partition="0",topic="test"} 6.64294327e+08
kafka_consumergroup_logsize{cluster="kafka-cluster",consumergroup="test",partition="1",topic="test"} 8.19800068e+08
...
```


