# client-go 学习
`client-go`是Kubernetes项目的官方Go语言客户端库，它为开发者提供了一组强大的工具和接口来与Kubernetes集群进行交互。这个库使得开发者能够编写Go应用程序来管理Kubernetes上的各种资源，比如Deployments、Services、Ingresses、Pods、Namespaces、Nodes等，执行创建、读取、更新和删除（CRUD）操作。

## 客户端类型

`client-go`主要提供了以下几种客户端对象：

1. **RESTClient**：这是最基础的客户端类型，它封装了HTTP请求，支持JSON和Protocol Buffers格式的数据交换。你可以直接使用它来发送RESTful风格的请求到Kubernetes API服务器。

2. **DiscoveryClient**：这个客户端用于发现Kubernetes API中的资源和API版本。当你需要动态地获取集群支持的API资源时，会用到它。

3. **ClientSet**：这是最常用的客户端类型，它为Kubernetes API中的每个资源类型都提供了对应的客户端对象，方便进行高度类型化的操作。通过ClientSet，你可以直接与特定资源进行交互，如`clientset.CoreV1().Pods("namespace").Get(name, metav1.GetOptions{})`来获取Pod。

4. **DynamicClient**：当需要处理未知或动态类型的资源时，DynamicClient就非常有用。它允许你以一种不关心具体资源类型的方式与API交互，适用于那些运行时才知道具体资源类型的场景。

## 如何使用

使用`client-go`通常遵循以下步骤：

1. **设置Kubernetes配置**：首先，你需要设置一个`rest.Config`对象，它通常来源于kubeconfig文件（如`~/.kube/config`），用于配置如何连接到Kubernetes集群。

    ```go
    config, err := rest.InClusterConfig()
    if err != nil {
        // 如果是在集群内运行（如作为Pod的一部分），使用InClusterConfig
        // 否则，可能需要使用rest.LoadFromConfig来从kubeconfig文件加载配置
    }
    ```

2. **创建ClientSet**：使用配置好的`rest.Config`实例化一个`clientset.Clientset`。

    ```go
    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        panic(err)
    }
    ```

3. **执行操作**：现在你可以使用`clientset`来进行各种资源操作了，比如列出所有Pods：

    ```go
    pods, err := clientset.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{})
    if err != nil {
        panic(err)
    }
    for _, pod := range pods.Items {
        fmt.Printf("Pod Name: %s, Status: %s\n", pod.Name, pod.Status.Phase)
    }
    ```

## 学习资源

- **官方文档**：查阅Kubernetes官方网站上的client-go文档，是最权威的学习资源。
- **GitHub示例**：client-go的GitHub仓库中有许多示例代码，位于`examples/`目录下，可以作为学习和参考的起点。
- **在线教程和视频**：如B站和知乎上的教程视频，这些通常提供更直观的操作演示和解释。
- **源码分析文章**：深入阅读一些技术社区的文章，特别是关于源码分析的内容，能帮助你更深入理解其内部机制。

记住，随着Kubernetes和client-go版本的更新，具体的API和用法可能会有所变化，因此始终参考最新的官方文档是明智的选择。



# 1.创建客户端

## 1.k8s集群外

参考 `client-go/examples/out-of-cluster-client-configuration/main.go`

```go
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

```



## 2.k8s集群内

参考 `client-go/examples/in-cluster-client-configuration/main.go`

```go
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
```



