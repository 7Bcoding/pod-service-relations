package main

import (
	"fmt"
	"pod-service-relations/config"
	"pod-service-relations/database"
	"pod-service-relations/model"
)

func main() {
	config.InitConfigs()
	database.Init()
	db := database.GetDB()
	err := db.AutoMigrate(&model.PodService{})
	if err != nil {
		fmt.Println(err)
	}
	err = db.AutoMigrate(&model.AbnormalPod{})
	if err != nil {
		fmt.Println(err)
	}
}
