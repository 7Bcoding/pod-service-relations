package client

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"pod-service-relations/logging"
)

func NewKubeClientSet(ctx context.Context, masterURL string, clusterFilePath string) *kubernetes.Clientset {
	// 在 kubeconfig 中使用当前上下文
	// path-to-kubeconfig -- 例如 /root/kubernetes_config/config
	//kubeApiServerConfig := config.NewKubeApiServerConfig()
	kubeConfig, err := clientcmd.BuildConfigFromFlags(masterURL, clusterFilePath)
	if err != nil {
		logging.GetLogger().Errorln(fmt.Printf("get kubeConfig got error: %s", err))
	}
	// 创建 clientSet
	clientset, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		logging.GetLogger().Errorln(fmt.Printf("get clientset got error: %s", err))
	}

	return clientset
}

func GetPods(ctx context.Context, masterURL string, clusterFilePath string) {
	clientset := NewKubeClientSet(ctx, masterURL, clusterFilePath)
	pods, _ := clientset.CoreV1().Pods("").List(context.TODO(), v1.ListOptions{})
	for _, pod := range pods.Items {
		fmt.Printf("namespace:%v\n name:%v\n labels:%v\n volumes:%v\n containers:%v\n",
			pod.Namespace, pod.Name, pod.Labels, pod.Spec.Volumes, pod.Spec.Containers)
	}

}

func GetServiceList(ctx context.Context, masterURL string, clusterFilePath string) []string {
	var serviceList []string
	clientset := NewKubeClientSet(ctx, masterURL, clusterFilePath)
	services, err := clientset.CoreV1().Services("").List(context.TODO(), v1.ListOptions{})
	if err != nil {
		logging.GetLogger().Errorln(fmt.Printf("GetServiceList get service list from k8s cluster failed: %s", err))
	}
	for _, serviceItem := range services.Items {
		//fmt.Printf("namespace:%v\n service_name:%v\n Selector:%v \n", serviceItem.Namespace, serviceItem.Name, serviceItem.Spec.Selector)
		serviceList = append(serviceList, serviceItem.Name)
	}
	return serviceList
}

func GetConfigMaps(ctx context.Context, masterURL string, clusterFilePath string) {
	clientset := NewKubeClientSet(ctx, masterURL, clusterFilePath)
	configMaps, err := clientset.CoreV1().ConfigMaps("").List(context.TODO(), v1.ListOptions{})
	if err != nil {
		logging.GetLogger().Errorln(fmt.Printf("get configMaps from k8s cluster failed: %s", err))
	}
	for _, ConfigMap := range configMaps.Items {
		fmt.Printf("namespace:%v\n configmap_name:%v\n data:%v \n", ConfigMap.Namespace, ConfigMap.Name, ConfigMap.Data)
		for key, value := range ConfigMap.Data {
			fmt.Printf("key: %s, value: %s \n", key, value)
		}
	}
}

// GetPodVolumeConfigMapMapping 获取Pod和configMap的映射关系Map
func GetPodVolumeConfigMapMapping(ctx context.Context, masterURL string, clusterFilePath string) map[string]string {
	var podVolumeConfigMapData map[string]string
	clientset := NewKubeClientSet(ctx, masterURL, clusterFilePath)
	podVolumeConfigMapData = make(map[string]string)
	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), v1.ListOptions{})
	if err != nil {
		logging.GetLogger().Errorln(fmt.Printf("GetPodVolumeConfigMapMapping get pods data from k8s cluster failed: %s", err))
	}
	for _, pod := range pods.Items {
		for _, volume := range pod.Spec.Volumes {
			//fmt.Printf("volume_info: %v\n", volume)
			if volume.ConfigMap != nil {
				podVolumeConfigMapData[pod.Name] = volume.ConfigMap.Name
			}
		}
	}
	return podVolumeConfigMapData
}

// GetPodServiceMap 获取Pod和Service的映射关系Map
func GetPodServiceMap(ctx context.Context, masterURL string, clusterFilePath string) map[string]string {
	var (
		podServiceMap           map[string]string
		serviceLabelSelectorMap map[string]string
		podServiceNameCount     map[string]int
	)
	clientset := NewKubeClientSet(ctx, masterURL, clusterFilePath)
	serviceLabelSelectorMap = make(map[string]string)
	podServiceMap = make(map[string]string)
	podServiceNameCount = make(map[string]int)
	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), v1.ListOptions{})
	if err != nil {
		logging.GetLogger().Errorln(fmt.Printf("GetPodServiceMap get pods data from k8s cluster failed: %s", err))
	}
	serviceLabelSelectorMap = GetServiceLabelSelectorMapping(ctx, masterURL, clusterFilePath)
	for _, pod := range pods.Items {
		for key, value := range pod.Labels {
			podLabel := key + ":" + value
			if serviceLabelSelectorMap[podLabel] != "" {
				serviceName := serviceLabelSelectorMap[podLabel]
				if podServiceNameCount[pod.Name+":"+serviceName] == 0 {
					podServiceNameCount[pod.Name+":"+serviceName]++
					podServiceMap[pod.Name] = podServiceMap[pod.Name] + serviceName + ";"
				}
			}
		}
	}
	return podServiceMap
}

// GetServiceLabelSelectorMapping 获取service和label-selector的映射关系Map
func GetServiceLabelSelectorMapping(ctx context.Context, masterURL string, clusterFilePath string) map[string]string {
	clientset := NewKubeClientSet(ctx, masterURL, clusterFilePath)
	var serviceLabelSelectorMap map[string]string
	serviceLabelSelectorMap = make(map[string]string)
	services, err := clientset.CoreV1().Services("").List(context.TODO(), v1.ListOptions{})
	if err != nil {
		logging.GetLogger().Errorln(fmt.Printf("GetServiceLabelSelectorMapping get services data from k8s cluster failed: %s", err))
	}
	for _, serviceItem := range services.Items {
		for key, value := range serviceItem.Spec.Selector {
			selectorLabel := key + ":" + value
			serviceLabelSelectorMap[selectorLabel] = serviceItem.Name
		}
	}
	return serviceLabelSelectorMap
}

func GetConfigMapsDataMapping(ctx context.Context, masterURL string, clusterFilePath string) map[string]string {
	clientset := NewKubeClientSet(ctx, masterURL, clusterFilePath)
	var configMapDataMap map[string]string
	var dataValues string
	configMapDataMap = make(map[string]string)
	configMaps, err := clientset.CoreV1().ConfigMaps("").List(context.TODO(), v1.ListOptions{})
	if err != nil {
		logging.GetLogger().Errorln(fmt.Printf("GetConfigMapsDataMapping get configMaps data from k8s cluster failed: %s", err))
	}
	for _, ConfigMap := range configMaps.Items {
		for key, value := range ConfigMap.Data {
			dataValues = dataValues + "\n" + key + ":" + value
		}
		configMapDataMap[ConfigMap.Name] = dataValues
		dataValues = ""
	}
	return configMapDataMap
}

func GetPodNamespaceMap(ctx context.Context, masterURL string, clusterFilePath string) map[string]string {
	clientset := NewKubeClientSet(ctx, masterURL, clusterFilePath)
	var podNamespaceMap map[string]string
	podNamespaceMap = make(map[string]string)
	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), v1.ListOptions{})
	if err != nil {
		logging.GetLogger().Errorln(fmt.Printf("GetPodNamespaceMap get pods data from k8s cluster failed: %s", err))
	}
	for _, pod := range pods.Items {
		podNamespaceMap[pod.Name] = pod.Namespace
	}
	return podNamespaceMap
}

func GetPodsContainersEnv(ctx context.Context, masterURL string, clusterFilePath string) map[string]string {
	clientset := NewKubeClientSet(ctx, masterURL, clusterFilePath)
	var podContainersEnv map[string]string
	podContainersEnv = make(map[string]string)
	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), v1.ListOptions{})
	if err != nil {
		logging.GetLogger().Errorln(fmt.Printf("GetPodsContainersEnv get pods data from k8s cluster failed: %s", err))
	}
	for _, pod := range pods.Items {
		envValue := ""
		for _, container := range pod.Spec.Containers {
			for _, env := range container.Env {
				envValue += "#" + env.Value
			}
		}
		podContainersEnv[pod.Name] = envValue
		envValue = ""
	}
	return podContainersEnv
}

func GetPodStatusPhase(ctx context.Context, masterURL string, clusterFilePath string) map[string]string {
	clientset := NewKubeClientSet(ctx, masterURL, clusterFilePath)
	var podStatusPhaseMap map[string]string
	podStatusPhaseMap = make(map[string]string)

	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), v1.ListOptions{})
	if err != nil {
		logging.GetLogger().Errorln(fmt.Printf("GetPodStatusPhase get pods data from k8s cluster failed: %s", err))
	}
	for _, pod := range pods.Items {
		podStatusPhaseMap[pod.Name] = string(pod.Status.Phase)
	}
	return podStatusPhaseMap
}

func GetPodConditionStatus(ctx context.Context, masterURL string, clusterFilePath string, namespace string, podName string) bool {
	clientset := NewKubeClientSet(ctx, masterURL, clusterFilePath)
	pod, err := clientset.CoreV1().Pods(namespace).Get(ctx, podName, v1.GetOptions{})
	if err != nil {
		logging.GetLogger().Errorln(fmt.Printf("GetPodConditionStatus get pods data from k8s cluster failed: %s", err))
	}
	for _, podCondition := range pod.Status.Conditions {
		if podCondition.Status != "True" {
			return false
		}
	}
	return true
}
