package main

import (
	"fmt"
	"gopkg.in/Shopify/sarama.v1"
	"regexp"
	"strconv"
	"time"
)

func updateOffsets() {
	startTime := time.Now()
	fmt.Println("Updating offset stats, Time: ", time.Now())
	oldConsumerGroups, err := zookeeperClient.Consumergroups()

	if err != nil {
		fmt.Println("Error reading consumergroup offsets: ", err)
		initClients()
		return
	}
	// 返回topiclist
	topics, err := zookeeperClient.Topics()

	if err != nil {
		initClients()
		return
	}
	// 获取topic
	for _, topic := range topics {
		if *topicsFilter != "" {
			// 排除topic
			match, err := regexp.MatchString(*topicsFilter, topic.Name)

			if err != nil {
				fmt.Println("Invalid Regex: ", err)
				panic("Exiting..")
			}

			if !match {
				fmt.Println("Filtering out " + topic.Name)
				continue
			}
		}
		// 获取partitions
		partitions, _ := topic.Partitions()
		for _, partition := range partitions {

			currentOffset, err := brokerClient.GetOffset(topic.Name, partition.ID, sarama.OffsetNewest)

			if err != nil  {
				fmt.Println("Error reading offsets from broker for topic, partition: ", topic.Name, partition, err)
				initClients()
				return
			}

			//fmt.Println(oldConsumerGroups)
			for _, group := range oldConsumerGroups {
				offset, _ := group.FetchOffset(topic.Name, partition.ID)
				if offset > 0 {
					consumerGroupLabels := map[string]string{"consumergroup": group.Name, "topic": topic.Name, "partition": strconv.Itoa(int(partition.ID))}
					//kafka_consumergroup_current_offset{cluster="kafka-cluster",consumergroup="test",partition="0",topic="test"} 1645
					consumerGroupGougeVec.With(consumerGroupLabels).Set(float64(offset))

					//kafka_consumergroup_logsize{cluster="kafka-cluster",consumergroup="test",partition="2",topic="test"} 103390
					consumerGroupLogSizeGougeVec.With(consumerGroupLabels).Set(float64(currentOffset))

					consumerGroupLag := currentOffset - offset
					//kafka_consumergroup_lag{cluster="kafka-cluster",consumergroup="test",partition="0",topic="test"} 6.64290361e+08
					consumerGroupLagGougeVec.With(consumerGroupLabels).Set(float64(consumerGroupLag))
				}

			}
		}
	}
	fmt.Println("Done updating offset stats in: ", time.Since(startTime))
}