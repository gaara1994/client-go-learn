package main

import (
	"context"
	"flag"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"log"
	"path/filepath"
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
	serviceName := "my-service"
	nameSpace := "client-test"
	//定义一个 Kubernetes service 对象的实例
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:    serviceName,
			Namespace: nameSpace,
			Labels: map[string]string{
				"app": "my-app",
			},
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
			Ports: []corev1.ServicePort{
				{
					Port:       80,
					TargetPort: intstr.FromInt(9376),
					Protocol:   corev1.ProtocolTCP,
				},
			},
			Selector: map[string]string{
				"app": "my-app",
			},
		},
	}
	//1.创建service
	serv, err := clientset.CoreV1().Services(nameSpace).Create(context.TODO(), service, metav1.CreateOptions{})
	if err != nil {
		log.Fatal("Create err ",err)
	}
	fmt.Println("Services 添加成功")
	fmt.Println(serv.Namespace)
	fmt.Println(serv.Name)

	//2.获取 Services
	old,err := clientset.CoreV1().Services(nameSpace).Get(context.TODO(),serviceName,metav1.GetOptions{})
	if err != nil {
		log.Fatal("Get err ",err)
	}
	fmt.Println("old",old.Name)


	//2.更新 Services
	old.Name = "new-service"
	updateServ, err := clientset.CoreV1().Services(nameSpace).Update(context.TODO(), serv, metav1.UpdateOptions{})
	if err != nil {
		log.Fatal("Update err ",err)
	}
	fmt.Println("新的updateServ：",updateServ.Name)

}
