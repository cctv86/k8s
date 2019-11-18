package  main

import (
	"flag"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"fmt"
	"k8s.io/client-go/kubernetes/typed/apps/v1beta1"
	"k8s.io/client-go/kubernetes/typed/core/v1"
	apiv1 "k8s.io/api/core/v1"
)


func main() {
	kubeconfig := flag.String("kubeconfig","/root/.kube/config","abs path to the kubeconfig file")

	config, err := clientcmd.BuildConfigFromFlags("",*kubeconfig)

	if err != nil {
		panic(err)
	}

	clientset, err := kubernetes.NewForConfig(config)

	if err != nil {
		panic(err)
	}
	
	//获取pod的日志
	podsName := ""  //pod的名称
	nsName := ""	// namespace
	i, err := clientset.CoreV1().Pods(nsName).GetLogs(
		podsName, &corev1.PodLogOptions{
			Follow:true,
			Container: "",
			//TailLines: &line,
		}).Stream()

	if err != nil {
		panic(err)
	}
	
	b ,err := ioutil.ReadAll(i)
	fmt.Println(err)
	fmt.Println(string(b))	
	

	
	deploymentclient := clientset.AppsV1beta1().Deployments(apiv1.NamespaceDefault)

	DeploymentListDemo(deploymentclient)

	podclient := clientset.CoreV1().Pods(apiv1.NamespaceDefault)
	
	PodListDemo(podclient)
	
	nodeclient := clientset.CoreV1().Nodes()
	
	NodeListDemo(nodeclient)

}

// deployment
func DeploymentListDemo(deploymentclient v1beta1.DeploymentInterface) {

	deployment, err := deploymentclient.List(metav1.ListOptions{})

	if err != nil {
		panic(err)
	}

	for _, i := range deployment.Items {
		fmt.Printf("%s have %d replices\n", i.Name, *i.Spec.Replicas)
	}

	// deployment名字
	fmt.Println(deployment.Items[0].Name)

	// deployment label
	fmt.Println(deployment.Items[1].Labels)

	// 所在namespace
	fmt.Println(deployment.Items[1].Namespace)

	fmt.Println(deployment.Items[1].Status.ReadyReplicas)
	fmt.Println(deployment.Items[1].Status.Conditions[0].Status)
}


// pod

func PodListDemo(podclient v1.PodInterface) {
	pods, err := podclient.List(metav1.ListOptions{})

	if err != nil {
		panic(err)
	}

	fmt.Println(pods.Items[1].Name)	//名称
	fmt.Println(pods.Items[1].CreationTimestamp)	// 创建时间
	fmt.Println(pods.Items[1].Labels)	// label
	fmt.Println(pods.Items[1].Namespace)	// 命名空间
	fmt.Println(pods.Items[1].Status.HostIP)	// 宿主机IP地址
	fmt.Println(pods.Items[1].Status.PodIP) // pod的IP地址
	fmt.Println(pods.Items[1].Status.StartTime) // 开始事件
	fmt.Println(pods.Items[1].Status.Phase) // 状态
	fmt.Println(pods.Items[1].Status.ContainerStatuses[0].RestartCount)   //重启次数
	fmt.Println(pods.Items[1].Status.ContainerStatuses[0].Image) //获取重启时间

}


// node
func NodeListDemo(nodeclient v1.NodeInterface) {
	nodes, err := nodeclient.List(metav1.ListOptions{})

	if err != nil {
		panic(err)
	}

	fmt.Println(nodes.Items[0].Name)	// 宿主机的主机名
	fmt.Println(nodes.Items[0].CreationTimestamp)    //加入集群时间
	fmt.Println(nodes.Items[0].Status.NodeInfo)	// 宿主机的操作系统信息，docker版本信息
	fmt.Println(nodes.Items[0].Status.Conditions[len(nodes.Items[0].Status.Conditions)-1].Type) // 状态,ready
	fmt.Println(nodes.Items[0].Status.Allocatable.Memory().String())  // 内存信息，单位字节
}

//
