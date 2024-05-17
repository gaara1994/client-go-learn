package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"path/filepath"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func applyResourceFromFile(clientset kubernetes.Interface, dynamicClient dynamic.Interface, filename string) error {
	// 读取YAML文件
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	// 解码YAML到unstructured.Unstructured对象
	decode := yaml.NewYAMLOrJSONDecoder(strings.NewReader(string(data)), 1024)
	obj := &unstructured.Unstructured{}
	if err := decode.Decode(obj); err != nil {
		return err
	}

	// 根据GroupVersionResource应用资源
	gvk := obj.GroupVersionKind()
	gvr := dynamicClient.Resource(gvk.GroupVersion().WithResource(gvk.Kind))


	// 检查是否已存在（可选）
	existing, err := gvr.Namespace(obj.GetNamespace()).Get(context.TODO(), obj.GetName(), v1.GetOptions{})
	if err == nil {
		// 如果已存在，则更新（可选）
		obj.SetResourceVersion(existing.GetResourceVersion())
		_, err = gvr.Namespace(obj.GetNamespace()).Update(context.TODO(), obj, v1.UpdateOptions{})
		if err != nil {
			return err
		}
		// 注意：这里注释了更新的逻辑，如果你需要更新资源，可以取消注释
		fmt.Printf("Resource %s/%s already exists, skipping creation\n", obj.GetNamespace(), obj.GetName())
		return nil
	}

	// 创建资源
	_, err = gvr.Namespace(obj.GetNamespace()).Create(context.TODO(), obj, v1.CreateOptions{})
	if err != nil {
		return err
	}

	fmt.Printf("Resource %s/%s created\n", obj.GetNamespace(), obj.GetName())
	return nil
}

func main() {
	kubeconfig := flag.String("kubeconfig", filepath.Join(homedir.HomeDir(), ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	filename := flag.String("filename", "resource.yaml", "path to the Kubernetes resource YAML file")
	flag.Parse()

	// 使用kubeconfig文件创建config
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// 创建clientset和dynamicClient
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// 应用资源
	if err := applyResourceFromFile(clientset, dynamicClient, *filename); err != nil {
		panic(err.Error())
	}
}