package main

import (
	"context"
	"flag"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"os"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/metrics/pkg/client/clientset/versioned"
)

func main() {
	kubeconfig := flag.String("kubeconfig", filepath.Join(homeDir(), ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	flag.Parse()

	// 使用kubeconfig文件创建config
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 创建metrics clientset
	metricsClient, err := versioned.NewForConfig(config)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 创建kubernetes clientset以获取节点信息
	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 获取所有节点的度量数据
	nodesMetrics, err := metricsClient.MetricsV1beta1().NodeMetricses().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Println(err)
		return
	}

	// 获取所有节点的信息
	nodes, err := kubeClient.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Println(err)
		return
	}

	// 遍历节点，匹配度量数据和节点信息以输出内存使用和总内存大小
	nodeMap := make(map[string]*corev1.Node)
	for _, node := range nodes.Items {
		nodeMap[node.Name] = &node
	}

	for _, nodeMetric := range nodesMetrics.Items {
		node, exists := nodeMap[nodeMetric.Name]
		if !exists {
			fmt.Printf("Node %s not found in node list.\n", nodeMetric.Name)
			continue
		}

		totalMemory := node.Status.Capacity["memory"]
		fmt.Printf("Node Name: %s\n", nodeMetric.Name)
		fmt.Printf("CPU Usage: %s\n", nodeMetric.Usage.Cpu().String())
		fmt.Printf("Memory Usage: %s\n", nodeMetric.Usage.Memory().String())
		fmt.Printf("Total Memory: %s\n", totalMemory.String())
	}
}

// homeDir helps find the user's home directory, no matter what platform we're on.
func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}