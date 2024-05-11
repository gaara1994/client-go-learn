package main

import (
	"context"
	"flag"
	"fmt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
)
func main() {
	// 创建字符串指针变量存储kubeconfig的文件地址
	var kubeconfigPath *string
	// 获取家目录
	home := homedir.HomeDir()
	if home != "" {
		// 有家目录时，设置默认的kubeconfig路径，并允许用户通过命令行参数 --kubeconfig=*** ip覆盖
		kubeconfigPath = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		// 没有检测到家目录，仍然定义kubeconfig参数，但提示用户必须提供（此时应避免暗示使用命令行传入家目录）
		kubeconfigPath = flag.String("kubeconfig", "", "Please provide the absolute path to the kubeconfig file")
	}

	// 确保在使用kubeconfigPath之前调用flag.Parse()解析命令行参数
	flag.Parse()
	fmt.Println("kubeconfig为：",*kubeconfigPath)

	config,err := clientcmd.BuildConfigFromFlags("",*kubeconfigPath)
	if err != nil {
		// 使用Fatal，打印日志后程序会退出
		log.Fatal(err)
	}

	//创建 clientset
	clientset,err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
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
