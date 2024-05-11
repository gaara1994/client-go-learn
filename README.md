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



# 1.ClientSet

`ClientSet`是一个结构体，它封装了一系列的客户端对象，每个客户端对象对应Kubernetes API的一个特定资源或一组资源。比如，你可以通过`ClientSet.CoreV1().Pods(namespace).List(...)`来列出指定命名空间中的所有Pod，或者使用`ClientSet.AppsV1().Deployments(namespace).Create(...)`来创建一个新的Deployment。

## 1. 功能与使用

- **资源操作**：提供了创建、读取、更新和删除（CRUD）Kubernetes资源的方法，涵盖了Pods、Services、Deployments、ConfigMaps等各种资源类型。
- **命名空间操作**：支持跨命名空间的操作，允许你针对特定命名空间执行API调用。
- **发现与列举**：可以用来发现API服务器上的资源类型，并列举资源实例。
- **配置与认证**：`ClientSet`的创建通常需要一个配置对象（`rest.Config`），该配置包含了访问API服务器所需的认证信息、服务器地址、TLS设置等。

## 2. 创建`ClientSet`

通常，你会使用`kubernetes.NewForConfig()`函数来从一个`rest.Config`实例创建一个`ClientSet`。下面是一个简单的示例：

### 1.k8s集群外

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



### 2.k8s集群内

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



# 2.不同资源的管理

`ClientSet`包含了针对不同组的客户端。在Kubernetes的Go客户端(`client-go`)中，`ClientSet`是一个综合性的结构，它整合了多个API组的客户端接口。API组是用来对Kubernetes API中的资源进行分类和版本管理的一种方式，支持扩展性和版本控制。

创建一个`ClientSet`时，它会为多个核心及扩展组预初始化客户端，让你能够与诸如以下资源进行互动：

- 核心组 (`core/v1`)：Pod、服务、配置映射、密钥、节点等。
- 应用组 (`apps/v1`)：部署、有状态集、守护进程集、副本集、任务、定时任务等。
- 批处理组 (`batch/v1`, `batch/v1beta1`)：任务、定时任务等。
- 网络组 (`networking.k8s.io/v1`)：入口、网络策略等。
- 以及其他更多，依据Kubernetes版本和可用的扩展而定。

每个组的客户端都提供了特定于该组资源的方法，让你能够执行创建、读取、更新、删除等操作以及其他与那些资源相关的活动。

简化的来说，`ClientSet`内部结构是这样的组织，展示它如何聚合不同组的客户端：

```go
type ClientSet struct {
    coreV1                        *corev1.CoreV1Client
    appsV1                        *appsv1.AppsV1Client
    batchV1                       *batchv1.BatchV1Client
    // ... 其他组的接口
}
```

你可以通过如`ClientSet.CoreV1()`、`ClientSet.AppsV1()`等方法访问这些接口，进而与相应的Kubernetes资源进行互动。

例如，若要创建一个部署，你会使用`AppsV1Interface`：

```go
deploymentClient := clientset.AppsV1().Deployments(namespace)
deployment, err := deploymentClient.Create(context.TODO(), deploymentSpec, metav1.CreateOptions{})
```

这种`ClientSet`的模块化设计确保了你的代码能够轻松地与广泛的Kubernetes资源互动，同时保持清晰和结构化的API调用方式。



k8s有不同类型的资源，不同类型的资源要调用client-go不同子包的方法，下面是子包的列表。

```shell
➜  typed git:(master) ✗ pwd
/home/yantao/go/src/github.com/gaara1994/client-go/kubernetes/typed
➜  typed git:(master) ✗ ll 
总用量 88K
drwxrwxr-x 5 yantao yantao 4.0K 5月  10 11:13 admissionregistration
drwxrwxr-x 3 yantao yantao 4.0K 5月  10 11:13 apiserverinternal
drwxrwxr-x 5 yantao yantao 4.0K 5月  10 11:13 apps
drwxrwxr-x 5 yantao yantao 4.0K 5月  10 11:13 authentication
drwxrwxr-x 4 yantao yantao 4.0K 5月  10 11:13 authorization
drwxrwxr-x 6 yantao yantao 4.0K 5月  10 11:13 autoscaling
drwxrwxr-x 4 yantao yantao 4.0K 5月  10 11:13 batch
drwxrwxr-x 5 yantao yantao 4.0K 5月  10 11:13 certificates
drwxrwxr-x 4 yantao yantao 4.0K 5月  10 11:13 coordination
drwxrwxr-x 3 yantao yantao 4.0K 5月  10 11:13 core
drwxrwxr-x 4 yantao yantao 4.0K 5月  10 11:13 discovery
drwxrwxr-x 4 yantao yantao 4.0K 5月  10 11:13 events
drwxrwxr-x 3 yantao yantao 4.0K 5月  10 11:13 extensions
drwxrwxr-x 6 yantao yantao 4.0K 5月  10 11:13 flowcontrol
drwxrwxr-x 5 yantao yantao 4.0K 5月  10 11:13 networking
drwxrwxr-x 5 yantao yantao 4.0K 5月  10 11:13 node
drwxrwxr-x 4 yantao yantao 4.0K 5月  10 11:13 policy
drwxrwxr-x 5 yantao yantao 4.0K 5月  10 11:13 rbac
drwxrwxr-x 3 yantao yantao 4.0K 5月  10 11:13 resource
drwxrwxr-x 5 yantao yantao 4.0K 5月  10 11:13 scheduling
drwxrwxr-x 5 yantao yantao 4.0K 5月  10 11:13 storage
drwxrwxr-x 3 yantao yantao 4.0K 5月  10 11:13 storagemigration

```

1. **核心资源(Core Resources)**: 调用 `clientset.CoreV1()`
   - Pods
   - Services
   - ConfigMaps
   - Secrets
   - PersistentVolumes 和 PersistentVolumeClaims
   - Nodes
   - Namespaces
2. **Apps 资源(Apps Resources)**: 调用 `clientset.AppsV1()`
   - Deployments
   - ReplicaSets
   - StatefulSets
   - DaemonSets
   - Jobs
   - CronJobs
3. **网络资源(Networking Resources)**: 调用 `clientset.NetworkingV1()`
   - Ingresses
   - NetworkPolicies
4. **存储资源(Storage Resources)**: 调用 `clientset.StorageV1()`
   - StorageClasses
5. **扩展资源(Extensions/CustomResources)**: 调用 `clientset.ExtensionsV1beta1()`
   - CustomResourceDefinitions (CRDs) —— 虽然直接通过标准`ClientSet`创建CRDs比较复杂，通常需要使用`apiextensions.k8s.io/v1` API组的客户端，但可以通过扩展`ClientSet`或使用`DynamicClient`来操作。
   - 各种自定义资源对象（由CRDs定义）
6. **安全与认证资源(Security and Authentication)**: 调用 `clientset.AuthorizationV1()`
   - Roles 和 ClusterRoles
   - RoleBindings 和 ClusterRoleBindings
   - ServiceAccounts
7. **API 资源(API Resources)**: 调用 `clientset.OpenAPIV3()`
   - APIServices
   - CustomResourceDefinitions
8. **其他资源**: 
   - Events	调用 `clientset.EventsV1()`
   - Endpoints   
   - HorizontalPodAutoscalers
   - PodDisruptionBudgets
   - LimitRanges
   - ResourceQuotas  ResourceV1alpha2

下面是整理的列表

| 资源类型                             | Clientset类型                 | 具体的资源                                                   |
| ------------------------------------ | :---------------------------- | ------------------------------------------------------------ |
| 核心资源(Core Resources)             | clientset.CoreV1()            | Pods  Services  ConfigMaps  Secrets  PersistentVolumes PersistentVolumeClaims  Nodes  Namespaces |
| Apps 资源(Apps Resources)            | clientset.AppsV1()            | Deployments  ReplicaSets  StatefulSets  DaemonSets  Jobs  CronJobs |
| 网络资源(Networking Resources)       | clientset.NetworkingV1()      | Ingresses  NetworkPolicies                                   |
| 存储资源(Storage Resources)          | clientset.StorageV1()         | StorageClasses                                               |
| 扩展资源(Extensions/CustomResources) | clientset.ExtensionsV1beta1() | CustomResourceDefinitions (CRDs)                             |
| RBAC资源                             | clientset.RbacV1()            | Roles                                                        |
| 认证与授权策略                       | clientset.AuthenticationV1()  |                                                              |
| API 资源(API Resources)              | clientset.OpenAPIV3()         | APIServices CustomResourceDefinitions                        |



### 1.核心资源(Core Resources)

#### 1.Pods

创建 监听 删除 demo3/main.go

```
Pod资源在Kubernetes中被认为是基本上不可变的，这意味着一旦Pod被创建，你不应该直接修改其定义，包括更换镜像。Pod的设计原则是围绕着它的 immutability（不变性）和 declarative configuration（声明式配置）理念构建的。当需要改变Pod的属性，比如镜像版本，推荐的做法是通过操作更高层次的抽象资源来间接实现，比如：
Deployments: 用于无状态应用，支持滚动更新、回滚等特性。
StatefulSets: 针对有状态应用，同样支持更新策略，同时保持Pod的唯一标识和稳定的存储。
```



#### 

