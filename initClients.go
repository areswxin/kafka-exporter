package main

import (
	"fmt"
	"github.com/krallistic/kazoo-go"
	"gopkg.in/Shopify/sarama.v1"
)

var zookeeperClient *kazoo.Kazoo
var brokerClient sarama.Client

//初始化连接
func initClients() {
	fmt.Println("Init zookeeper client with connection string: ", *zookeeperConnect)
	var err error
	zookeeperClient, err = kazoo.NewKazooFromConnectionString(*zookeeperConnect, nil)
	if err != nil {
		fmt.Println("Error Init zookeeper client with connection string:", *zookeeperConnect)
		panic(err)
	}
	// 获取kafka地址
	brokers, err := zookeeperClient.BrokerList()
	if err != nil {
		fmt.Println("Error reading brokers from zk")
		panic(err)
	}

	//初始化kafka
	fmt.Println("Init Kafka Client with Brokers:", brokers)
	config := sarama.NewConfig()
	brokerClient, err = sarama.NewClient(brokers, config)

	if err != nil {
		fmt.Println("Error Init Kafka Client")
		panic(err)
	}
	fmt.Println("Done Init Clients")
}