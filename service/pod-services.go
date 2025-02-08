package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"os"
	"pod-service-relations/client"
	"pod-service-relations/config"
	"pod-service-relations/dao"
	"pod-service-relations/database"
	"pod-service-relations/logging"
	"pod-service-relations/model"
	"strings"
	"time"
)

func GetPodToServiceRelations(ctx context.Context) {
	var (
		clusterFileMap         map[string]string
		clusterMasterURLMap    map[string]string
		serviceList            []string
		configMapDataMap       map[string]string
		podVolumeConfigMapData map[string]string
		podToServiceRelations  map[string][]string
		podNamespaceMap        map[string]string
		podServiceFilter       map[string]int
		podContainersEnv       map[string]string
		podServiceModel        *model.PodService
		podServiceModelList    []*model.PodService
		podServiceMap          map[string]string
		podStatusPhaseMap      map[string]string
		podServiceFieldsMap    map[string]interface{}
		abnormalPodFieldsMap   map[string]interface{}
		abnormalPodServiceInfo []*model.PodService
		affectPodCount         map[string]int
		affectServiceCount     map[string]int
		affectPods             []string
		affectServices         []string
		cluster                string
	)
	ignoreServiceList := "http;redis;mysql;ss;db;service"
	kubeApiServerConfig := config.NewKubeApiServerConfig()
	clusterFileMap = make(map[string]string)
	affectPodCount = make(map[string]int)
	affectServiceCount = make(map[string]int)
	clusterMasterURLMap = make(map[string]string)
	podServiceModelList = make([]*model.PodService, 1)
	affectPods = make([]string, 1)
	affectServices = make([]string, 1)
	fileDirPath := kubeApiServerConfig.KubeConfigPath
	fileAbsPaths, err := client.GetAllFiles(fileDirPath)
	dbConn := database.GetDB()
	ps := dao.NewPodServiceDao(dbConn)
	ap := dao.NewAbnormalPodDao(dbConn)
	if err != nil {
		logging.GetLogger().Errorln(fmt.Sprintf("get kube config dir path failed: %s", err))
	}
	for _, fileAbsPath := range fileAbsPaths {
		filePathArr := strings.Split(fileAbsPath, fileDirPath)
		if len(filePathArr) > 0 {
			if filePathArr[1] != "" {
				fileName := filePathArr[1]
				fileNameArr := strings.Split(fileName, "_")
				clusterName := ""
				if len(fileNameArr) > 0 {
					clusterName = fileNameArr[0]
					clusterFileMap[clusterName] = fileName
				}
				file, err := os.Open(fileAbsPath)
				if err != nil {
					logging.GetLogger().Errorln(err)
				}
				defer file.Close()
				// 设置配置文件的名字
				viper.SetConfigName(fileName)
				// 设置配置文件的类型
				viper.SetConfigType("yaml")
				// 添加配置文件的路径，指定 config 目录下寻找
				viper.AddConfigPath(fileDirPath)
				// 寻找配置文件并读取
				err = viper.ReadInConfig()
				if err != nil {
					logging.GetLogger().Errorln(fmt.Errorf("fatal error while viper read config file: %v", err))
				}
				clustersConfig := viper.Get("clusters")
				if clusters, ok := clustersConfig.([]interface{}); ok {
					clusterConfig := clusters[0]
					if clusterData, ok := clusterConfig.(map[string]interface{}); ok {
						if cluster, ok := clusterData["cluster"].(map[string]interface{}); ok {
							if masterURL, ok := cluster["server"].(string); ok {
								clusterMasterURLMap[clusterName] = masterURL
							}
						}
					}
				}
			}
		}
	}

	for clusterName, clusterfileName := range clusterFileMap {
		configMapDataMap = make(map[string]string)
		podVolumeConfigMapData = make(map[string]string)
		podToServiceRelations = make(map[string][]string)
		podNamespaceMap = make(map[string]string)
		podServiceFilter = make(map[string]int)
		podContainersEnv = make(map[string]string)
		podServiceMap = make(map[string]string)
		podStatusPhaseMap = make(map[string]string)
		podServiceFieldsMap = make(map[string]interface{})
		masterURL := clusterMasterURLMap[clusterName]
		kubeConfigPath := kubeApiServerConfig.KubeConfigPath + clusterfileName
		// 获取serviceName列表
		serviceList = client.GetServiceList(ctx, masterURL, kubeConfigPath)
		// 获取configmap的数据
		configMapDataMap = client.GetConfigMapsDataMapping(ctx, masterURL, kubeConfigPath)
		// 获取pod-configmap映射关系
		podVolumeConfigMapData = client.GetPodVolumeConfigMapMapping(ctx, masterURL, kubeConfigPath)
		// 获取pod-namespace映射关系
		podNamespaceMap = client.GetPodNamespaceMap(ctx, masterURL, kubeConfigPath)
		// 获取pod中container env的数据
		podContainersEnv = client.GetPodsContainersEnv(ctx, masterURL, kubeConfigPath)
		// 获取pod-service的映射关系
		podServiceMap = client.GetPodServiceMap(ctx, masterURL, kubeConfigPath)
		// 获取pod当前状态
		podStatusPhaseMap = client.GetPodStatusPhase(ctx, masterURL, kubeConfigPath)

		// service匹配configMap
		for _, serviceName := range serviceList {
			for podName, configMapName := range podVolumeConfigMapData {
				configMapData := configMapDataMap[configMapName]
				if strings.Contains(ignoreServiceList, serviceName) {
					continue
				}
				// 如果podName的configmap里包含podName1的serviceName，说明该serviceName被podName的configmap调用
				if strings.Contains(configMapData, serviceName+".") || strings.Contains(configMapData, serviceName+":") {
					if podServiceFilter[podName+":"+serviceName] < 1 {
						podServiceFilter[podName+":"+serviceName]++
						podToServiceRelations[podName] = append(podToServiceRelations[podName], serviceName)
					}
				}
			}
		}
		// service匹配container env
		for _, serviceName := range serviceList {
			for podName, envValue := range podContainersEnv {
				if strings.Contains(ignoreServiceList, serviceName) {
					continue
				}
				if strings.Contains(envValue, serviceName) {
					if podServiceFilter[podName+":"+serviceName] < 1 {
						podServiceFilter[podName+":"+serviceName]++
						podToServiceRelations[podName] = append(podToServiceRelations[podName], serviceName)
					}
				}
			}
		}

		for podName, podToServices := range podToServiceRelations {
			namespace := podNamespaceMap[podName]
			podServiceModel = &model.PodService{}
			// 保存clusterName时把斜杠去掉
			clusterNameArr := strings.Split(clusterName, "/")
			if len(clusterNameArr) > 0 {
				podServiceModel.Cluster = clusterNameArr[1]
			}
			podServiceModel.PodName = podName
			podServiceModel.PodStatus = podStatusPhaseMap[podName]
			podServiceModel.Namespace = namespace
			podServiceModel.ServiceName = podServiceMap[podName]
			podServiceModel.PodToService = strings.Replace(strings.Trim(fmt.Sprint(podToServices), "[]"), " ", ";", -1)
			podServiceModelList = append(podServiceModelList, podServiceModel)
			pod, err := ps.Get(podName)
			if err != nil {
				// 查不到，就插入数据
				if err == gorm.ErrRecordNotFound {
					if err := ps.Create(podServiceModel); err != nil {
						logging.GetLogger().Errorln(fmt.Sprintf("create pod service record failed: %s", err))
						fmt.Printf("create pod service record failed: %s", err)
					}
				} else {
					logging.GetLogger().Errorln(fmt.Printf("get pod service record failed: %s", err))
					fmt.Printf("get pod service record failed: %s", err)
				}
				//	查得到，就更新数据
			} else {
				podServiceModel.ID = pod.ID
				now := time.Now().In(time.Local) // 获取本地时间
				beijingTime, _ := time.ParseInLocation("2006-01-02 15:04:05", now.Format("2006-01-02 15:04:05"), time.Local)
				podServiceModel.UpdatedAt = beijingTime
				podServiceBytes, err := json.Marshal(&podServiceModel)
				if err != nil {
					logging.GetLogger().Errorln(fmt.Printf("json marshal failed: %s", err))
					fmt.Printf("json marshal failed: %s", err)
				}
				err = json.Unmarshal(podServiceBytes, &podServiceFieldsMap)
				if err != nil {
					logging.GetLogger().Errorln(fmt.Printf("json unmarshal failed: %s", err))
					fmt.Printf("json unmarshal failed: %s", err)
				}
				delete(podServiceFieldsMap, "created_at")
				err = ps.Update(pod, podServiceFieldsMap)
				if err != nil {
					logging.GetLogger().Errorln(fmt.Printf("update podService to DB failed: %s", err))
					fmt.Printf("update podService to DB failed: %s", err)
				}
			}
		}
		clusterNameArr := strings.Split(clusterName, "/")
		if len(clusterNameArr) > 0 {
			cluster = clusterNameArr[1]
		} else {
			cluster = ""
		}
		abnormalPodServiceInfo, err = ps.GetAbnormalPodInfoByStatus(cluster)
		for _, abnormalPodServiceData := range abnormalPodServiceInfo {
			podConditionStatus := client.GetPodConditionStatus(ctx, masterURL, kubeConfigPath, abnormalPodServiceData.Namespace, abnormalPodServiceData.PodName)
			podNameSplit := strings.Split(abnormalPodServiceData.PodName, "-")
			podNameArr := podNameSplit[:len(podNameSplit)-2]
			originPodName := strings.Join(podNameArr, "-")
			podNum, err := ps.GetRunningPodNumByName(originPodName)
			if err != nil {
				logging.GetLogger().Errorln(fmt.Printf("get pod num by pod name failed: %s", err))
			}
			abnormalPodModel := &model.AbnormalPod{}
			abnormalPodFieldsMap = make(map[string]interface{})

			// running状态的pod数量小于1，则存入db
			abnormalPod, err := ap.Get(originPodName)
			if podNum <= 0 || podConditionStatus == false {
				for _, podServiceData := range podServiceModelList {
					if podServiceData != nil && abnormalPodServiceData.ServiceName != "" &&
						podServiceData.PodToService != "" && abnormalPodServiceData.ServiceName != ";" {
						abnormalServiceList := strings.Split(abnormalPodServiceData.ServiceName, ";")
						for _, abnormalServiceName := range abnormalServiceList {
							if strings.Contains(podServiceData.PodToService, abnormalServiceName) && abnormalServiceName != "" {
								if strings.Contains(ignoreServiceList, podServiceData.PodToService) ||
									strings.Contains(ignoreServiceList, abnormalServiceName) {
									continue
								}
								if podServiceData.ServiceName != "" && podServiceData.PodName != "" {
									if affectServiceCount[podServiceData.ServiceName] < 1 {
										affectServices = append(affectServices, podServiceData.ServiceName)
										affectServiceCount[podServiceData.ServiceName]++
									}
									affectPodNameSplit := strings.Split(podServiceData.PodName, "-")
									affectPodNameArr := affectPodNameSplit[:len(affectPodNameSplit)-2]
									affectPodName := strings.Join(affectPodNameArr, "-")
									if affectPodCount[affectPodName] < 1 {
										affectPods = append(affectPods, affectPodName)
										affectPodCount[affectPodName]++
									}
								}
							}
						}
					}
				}
				abnormalPodModel.AffectServices = strings.Join(affectServices, ";")
				abnormalPodModel.AffectPods = strings.Join(affectPods, ";")
				affectServices = nil
				affectPods = nil
				if err != nil {
					if err == gorm.ErrRecordNotFound {
						if originPodName != "" {
							abnormalPodModel.PodName = originPodName
							abnormalPodModel.ServiceName = abnormalPodServiceData.ServiceName
							abnormalPodModel.Cluster = abnormalPodServiceData.Cluster
							abnormalPodModel.Namespace = abnormalPodServiceData.Namespace
							if err := ap.Create(abnormalPodModel); err != nil {
								logging.GetLogger().Errorln(fmt.Printf("create abnormal pod record failed: %s", err))
								fmt.Printf("create abnormal pod record failed: %s", err)
							}
						} else {
							logging.GetLogger().Errorln(fmt.Printf("get pod service record failed: %s", err))
						}
					}
				} else {

					abnormalPod.PodName = originPodName
					abnormalPod.ServiceName = abnormalPodServiceData.ServiceName

					now := time.Now().In(time.Local) // 获取本地时间
					beijingTime, _ := time.ParseInLocation("2006-01-02 15:04:05", now.Format("2006-01-02 15:04:05"), time.Local)
					abnormalPod.UpdatedAt = beijingTime
					abnormalPodBytes, err := json.Marshal(&abnormalPod)
					if err != nil {
						logging.GetLogger().Errorln(fmt.Printf("json marshal failed: %s", err))
						fmt.Printf("json marshal failed: %s", err)
					}
					err = json.Unmarshal(abnormalPodBytes, &abnormalPodFieldsMap)
					if err != nil {
						logging.GetLogger().Errorln(fmt.Printf("json unmarshal failed: %s", err))
						fmt.Printf("json unmarshal failed: %s", err)
					}
					delete(abnormalPodFieldsMap, "created_at")
					err = ap.Update(abnormalPod, abnormalPodFieldsMap)
					if err != nil {
						logging.GetLogger().Errorln(fmt.Printf("update podService to DB failed: %s", err))
						fmt.Printf("update podService to DB failed: %s", err)
					}
				}
			} else {
				// pod副本数大于0，该pod已脱离异常情况，从db中删除
				if err != nil {
					// 没查到，略过
					//logging.GetLogger().Errorln(fmt.Printf("get pod service record failed: %s", err))
					fmt.Printf("get abnormal pod record failed: %s", err)
				} else {
					// 查到了，删除
					if podNum > 0 && podConditionStatus == true {
						err := ap.Delete(abnormalPod.ID, abnormalPod.PodName)
						if err != nil {
							logging.GetLogger().Errorln(fmt.Printf("delete abnormalPod from DB failed: %s", err))
							fmt.Printf("delete abnormalPod from DB failed: %s", err)
						}
					}
				}
			}
		}
	}
}

//func GetServicesListByYamlFile(ctx context.Context) []string {
//	var serviceNameList []string
//
//	fmt.Println(serviceNameList)
//	fileDirPath := ""
//	fileAbsPaths, err := client.GetAllFiles(fileDirPath)
//	if err != nil {
//		logging.GetLogger().Errorln(fmt.Sprintf("get service yaml files and dirs path failed: %s", err))
//	}
//	fmt.Println(fileAbsPaths)
//	for _, fileAbsPath := range fileAbsPaths {
//		filePathArr := strings.Split(fileAbsPath, fileDirPath+"\\")
//		if len(filePathArr) > 0 {
//			if filePathArr[1] != "" {
//				fileNameComplete := filePathArr[1]
//				// 筛选去除部分文件
//				if strings.HasPrefix(fileNameComplete, ".") || strings.HasSuffix(fileNameComplete, ".md") ||
//					strings.HasSuffix(fileNameComplete, ".sh") {
//				} else {
//					if strings.HasSuffix(fileNameComplete, ".yml") {
//						// 得到.yml格式的文件
//						//fileData, err := os.ReadFile(fileAbsPath)
//						//if err != nil {
//						//	logging.GetLogger().Errorln(fmt.Sprintf("get service name list data failed: %s", err))
//						//}
//						//fileText := string(fileData)
//						//fmt.Println(fileText)
//						file, err := os.Open(fileAbsPath)
//						if err != nil {
//							logging.GetLogger().Errorln(err)
//						}
//						defer file.Close()
//						// 设置配置文件的名字
//						fileNames := strings.Split(fileNameComplete, ".")
//						fileName := fileNames[0]
//						viper.SetConfigName(fileName)
//						// 设置配置文件的类型
//						viper.SetConfigType("yaml")
//						// 添加配置文件的路径，指定 config 目录下寻找
//						viper.AddConfigPath(fileDirPath)
//						// 寻找配置文件并读取
//						err = viper.ReadInConfig()
//						if err != nil {
//							logging.GetLogger().Errorln((fmt.Errorf("fatal error config file: %w", err)))
//						}
//						if serviceName, ok := viper.Get("metadata.name").(string); ok {
//							serviceNameList = append(serviceNameList, serviceName)
//						}
//					}
//				}
//			}
//		}
//	}
//	fmt.Println(serviceNameList)
//	return serviceNameList
//}
//
//func MatchConfigMapByYamlFile(ctx context.Context, serviceNameList []string) {
//	fileDirPath := ""
//	fileAbsPaths, err := client.GetAllFiles(fileDirPath)
//	if err != nil {
//		logging.GetLogger().Errorln(fmt.Sprintf("get configmap yaml files and dirs path failed: %s", err))
//	}
//	fmt.Println(fileAbsPaths)
//	for _, fileAbsPath := range fileAbsPaths {
//		filePathArr := strings.Split(fileAbsPath, fileDirPath+"\\")
//		if len(filePathArr) > 0 {
//			if filePathArr[1] != "" {
//				configMapfileName := filePathArr[1]
//				// 筛选去除部分文件
//				if strings.HasPrefix(configMapfileName, ".") || strings.HasSuffix(configMapfileName, ".md") ||
//					strings.HasSuffix(configMapfileName, ".sh") {
//				} else {
//					if strings.HasSuffix(configMapfileName, ".yml") {
//						fileData, err := os.ReadFile(fileAbsPath)
//						if err != nil {
//							logging.GetLogger().Errorln(fmt.Sprintf("get service name list data failed: %s", err))
//						}
//						fileText := string(fileData)
//						//fmt.Println(fileText)
//						for _, serviceName := range serviceNameList {
//							if strings.Contains(fileText, serviceName) {
//								fmt.Println(fmt.Sprintf("configMapFileName: %s, serviceName: %s",
//									configMapfileName, serviceName))
//							}
//						}
//					}
//				}
//			}
//		}
//	}
//}
