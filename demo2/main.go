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

	//获取所有的Pods
	list, err := clientset.CoreV1().Pods("").List(context.Background(),metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	for _, pod := range list.Items {
		//打印pod的名字
		fmt.Println(pod.Name)
	}
}

// go build -o app