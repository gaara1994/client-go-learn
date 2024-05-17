package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// 加载kubeconfig文件以建立连接
	kubeconfigPath := filepath.Join(homeDir(), ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		panic(err.Error())
	}

	// 创建一个新的Kubernetes API客户端实例
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// 获取所有节点的信息
	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	// 遍历节点并打印注解
	for _, node := range nodes.Items {
		if node.Name == "testnode1" {
			annotations := node.Annotations
			fmt.Printf("Node Name: %s\n", node.Name)
			fmt.Println("Annotations:")
			for key, value := range annotations {
				fmt.Printf("\t%s: %s\n", key, value)
			}
			fmt.Println("--------------------")
		}

	}
}

// homeDir helps find the user's home directory, no matter what platform we're on.
func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}