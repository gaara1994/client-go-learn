package main

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
)
func main() {
	//创建集群内的配置
	config,err := rest.InClusterConfig()
	if err != nil {
		log.Fatal("rest.InClusterConfig err:",err)
	}
	clientset,err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal("kubernetes.NewForConfig err:",err)
	}

	// 获取节点列表
	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatal("Nodes().List err:",err)
	}

	// 打印每个节点的名称和内存信息
	for _, node := range nodes.Items {
		fmt.Printf("Node Name: %v\n", node.Name)
		for _, status := range node.Status.Conditions {
			if status.Type == "MemoryPressure" {
				fmt.Printf("Memory Pressure: %v\n", status.Status)
			}
		}
		// 注意：直接获取内存使用量可能需要解析Status.Capacity和Status.Allocatable字段，
		// 但这些字段提供的信息是总量而非实时使用量。实时内存使用量通常需要通过metrics-server获取。
	}
}

// go build -o app