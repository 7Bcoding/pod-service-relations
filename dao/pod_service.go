package dao

import (
	"fmt"
	"go/types"
	"pod-service-relations/logging"
	"time"

	"gorm.io/gorm"

	"pod-service-relations/model"
)

type podServiceDao struct {
	db *gorm.DB
}

type PodServiceDao interface {
	Create(PodService *model.PodService) error
	Get(PodName string) (*model.PodService, error)
	ALL(filters map[string]interface{}) ([]*model.PodService, error)
	Delete(ID uint32, PodName string) error
	Update(Pod *model.PodService, fields map[string]interface{}) error
	GetAbnormalPodInfoByStatus(clusterName string) ([]*model.PodService, error)
	GetRunningPodNumByName(podName string) (int64, error)
}

func NewPodServiceDao(db *gorm.DB) PodServiceDao {
	return &podServiceDao{db: db}
}

// Create
/*
* insert Pod
* PARAMS:
*	- Pod: model Pod
* RETURN:
* 	- error: if success, nil, else gorm error and other error
 */
func (ps podServiceDao) Create(PodService *model.PodService) error {
	return ps.db.Create(PodService).Error
}

// Get
/*
* select single Pod
* params:
*	- id: pk
*	- name: unique key
* RETURN:
*	- model.Pod
*	- error: if success, nil, else gorm error and other error
 */
func (ps podServiceDao) Get(PodName string) (*model.PodService, error) {
	if PodName == "" {
		return nil, fmt.Errorf(
			"Pod_name must set one: %s", PodName,
		)
	}
	Pod := new(model.PodService)
	db := ps.db
	if PodName != "" {
		db = db.Where("pod_name = ?", PodName)
	}
	err := db.First(Pod).Error
	if err != nil {
		return nil, err
	}
	return Pod, nil
}

// ALL
/*
* select batch Pods
* params:
*	- filters: map[string]interface{}
*		key: db column
*		value: interface{}, maybe uint32, string, slice
*		like {"id": 1, "name": ["test1", "test2"]}
* RETURN:
*	- model.Pod slice
*	- error: if success, nil, else gorm error and other error
 */
func (ps podServiceDao) ALL(filters map[string]interface{}) ([]*model.PodService, error) {
	var PodServices []*model.PodService
	db := ps.db
	for k, v := range filters {
		switch value := v.(type) {
		case string, uint32, time.Time:
			db = db.Where(k+"= ?", value)
		case types.Slice:
			db = db.Where(k+"IN ?", value)
		default:
			return nil, fmt.Errorf("unsupport value type: %s", v)
		}
	}
	err := db.Find(&PodServices).Error
	return PodServices, err
}

// Delete
/*
* delete single Pod
* PARAMS:
* 	- id: pk
*	- name: unique key
* RETURN:
*	- error: if success, nil, else gorm error and other error
 */
func (ps podServiceDao) Delete(id uint32, name string) error {
	db := ps.db
	// id or name must set, else will delete all tables
	if id > 0 || name != "" {
		if id > 0 {
			db = db.Where("id = ?", id)
		}
		if name != "" {
			db = db.Where("pod_name = ?", name)
		}
	} else {
		return fmt.Errorf("id or name must set one")
	}

	return db.Delete(model.PodService{}).Error
}

// Update
/*
* update single Pod
* PARAMS:
* 	- action: model.Pod
*	- filters: map[string]interface{}
*		key: db column
*		value: interface{}, maybe uint32, string
*		like {"bin_repo": "baidu@123"}
* RETURN:
* 	- error: if success nil, else gorm error or other error
 */
func (ps podServiceDao) Update(Pod *model.PodService, fields map[string]interface{}) error {
	return ps.db.Model(Pod).Updates(fields).Error
}

func (ps podServiceDao) GetAbnormalPodInfoByStatus(clusterName string) ([]*model.PodService, error) {
	var PodsInfo []*model.PodService
	sql := `SELECT * FROM ` + model.PodServiceTABLE + ` WHERE cluster='` + clusterName + `' AND pod_status NOT IN ('Succeeded', 'Completed', 'Terminating', 'ContainerCreating') AND pod_status IS NOT NULL`
	db := ps.db
	res := db.Table("pod_services").Raw(sql).Scan(&PodsInfo)
	if res.Error != nil {
		logging.GetLogger().Errorln(fmt.Printf("get pods info by status failed: %s", res.Error))
	}
	return PodsInfo, nil
}

func (ps podServiceDao) GetRunningPodNumByName(podName string) (int64, error) {
	var count int64
	db := ps.db
	res := db.Model(&model.PodService{}).Where("pod_name like ? and pod_status = 'Running' and pod_status is not null", podName+"%").Count(&count)
	if res.Error != nil {
		logging.GetLogger().Errorln(fmt.Printf("get pods info by pod_name failed: %s", res.Error))
	}
	return count, nil
}
