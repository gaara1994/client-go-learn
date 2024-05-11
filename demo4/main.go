package main

import (
	"context"
	"flag"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"log"
	"path/filepath"
	"time"
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
	fmt.Println("kubeconfig为：", *kubeconfigPath)

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfigPath)
	if err != nil {
		// 使用Fatal，打印日志后程序会退出
		log.Fatal(err)
	}

	//创建 clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}
	podName := "client-test-alpine-1"
	nameSpace := "client-test"
	//定义一个 Kubernetes Pod 对象的实例
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: "client-test",
			Labels:    nil,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    "alpine-sleep",
					Image:   "alpine",
					Command: []string{"/bin/sh", "-c", "echo start-sleep && sleep 100"},
				},
			},
		},
	}
	//1.创建Pods
	podInfo, err := clientset.CoreV1().Pods("client-test").Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		log.Fatal("Create err ",err)
	}
	fmt.Println("添加成功")
	fmt.Println(podInfo.Namespace)
	fmt.Println(podInfo.Name)


	//2.监听Pod
	watch, err := clientset.CoreV1().Pods("client-test").Watch(context.TODO(), metav1.ListOptions{Watch: true})
	if err != nil {
		log.Fatal("Watch err ",err)
	}
	deletedChan := make(chan bool,0)
	defer watch.Stop()
	go func() {
		for event := range watch.ResultChan() {
			fmt.Println("pod changed==========================================================================")
			fmt.Println("event.Type===", event.Type)
			switch event.Type {
			//case "MODIFIED":
			//	fmt.Println("有更新")
			case "DELETED":
				fmt.Println("删除成功")
				deletedChan<- true
			}
		}
	}()

	//3.获取当前Pod实例
	time.Sleep(60 * time.Second)
	oldPod, err := clientset.CoreV1().Pods(nameSpace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		return
	}
	fmt.Println("当前注解：", oldPod.ObjectMeta.Annotations)

	// 更新Pod的注解，例如添加一个新的注解或者修改已存在的注解
	if oldPod.ObjectMeta.Annotations == nil {
		oldPod.ObjectMeta.Annotations = make(map[string]string)
	}
	oldPod.ObjectMeta.Annotations["update-example"] = "This pod was updated at " + time.Now().Format(time.RFC3339)

	//3.更新pod
	/*
	Pod资源在Kubernetes中被认为是基本上不可变的，这意味着一旦Pod被创建，你不应该直接修改其定义，包括更换镜像。Pod的设计原则是围绕着它的 immutability（不变性）和 declarative configuration（声明式配置）理念构建的。当需要改变Pod的属性，比如镜像版本，推荐的做法是通过操作更高层次的抽象资源来间接实现，比如：
	Deployments: 用于无状态应用，支持滚动更新、回滚等特性。
	StatefulSets: 针对有状态应用，同样支持更新策略，同时保持Pod的唯一标识和稳定的存储。
	*/
	updatedPod, err := clientset.CoreV1().Pods(nameSpace).Update(context.TODO(),oldPod,metav1.UpdateOptions{})
	if err != nil {
		log.Fatalf("Failed to update Pod %q: %v", podName, err)
	}
	fmt.Println("新的注解：", updatedPod.ObjectMeta.Annotations)


	//4.删除pod
	time.Sleep(100 * time.Second)
	err = clientset.CoreV1().Pods(nameSpace).Delete(context.TODO(), podName, metav1.DeleteOptions{})
	if err != nil {
		log.Fatal("Delete err ",err)
	}
	fmt.Println("删除中")

	res := <-deletedChan
	fmt.Println("删除结果：",res)
}
