package controller

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"pod-service-relations/dao"
	"pod-service-relations/database"
	"pod-service-relations/model"
	"pod-service-relations/utils"
)

type ServiceController struct{}

func (sc ServiceController) GetService(ctx *gin.Context) {
	serviceName := ctx.Param("service_name")
	params := struct {
		BNS     string `form:"bns"`
		Verbose bool   `form:"verbose"`
		//BinVersion  string `form:"bin_version"`
		//Confversion string `form:"conf_version"`
	}{}
	_ = ctx.ShouldBindQuery(&params)

	serviceDao := dao.NewServiceDao(database.GetDB())

	service, err := serviceDao.Get(0, serviceName)
	if err == gorm.ErrRecordNotFound {
		respondNotFound(ctx, err)
		return
	} else if err != nil {
		respondInternalError(ctx, err)
		return
	}
	respondSuccess(ctx, service)
}

func (sc ServiceController) GetServices(ctx *gin.Context) {
	respondSuccess(ctx, "todo")
}

func (sc ServiceController) AddService(ctx *gin.Context) {
	params := struct {
		Name        string `json:"name" binding:"required"`
		Arch        string `json:"arch"`
		OS          string `json:"os"`
		SyncType    string `json:"sync_type"`
		BinRepo     string `json:"bin_repo"`
		BinType     string `json:"bin_type"`
		BinFtpAddr  string `json:"bin_ftp_addr"`
		ConfRepo    string `json:"conf_repo"`
		ConfType    string `json:"conf_type"`
		ConfFtpAddr string `json:"conf_ftp_addr"`
	}{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		respondBadRequest(ctx, err)
		return
	}
	var service = new(model.Service)
	utils.CopyFields(params, service)

	service, err := dao.NewServiceDao(database.GetDB()).Create(service)
	if err != nil {
		respondInternalError(ctx, fmt.Errorf("database error: %w", err))
		return
	}
	respondSuccess(ctx, service)
}

func (sc ServiceController) DeleteService(ctx *gin.Context) {
	serviceID := uint32(getIntParam(ctx, "service_id"))
	if serviceID == 0 {
		respondBadRequest(ctx, fmt.Errorf("invalid service_id"))
		return
	}

	tx := database.GetDB().Begin()
	if err := tx.Error; err != nil {
		respondInternalError(ctx, fmt.Errorf("database error: %w", err))
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	sd := dao.NewServiceDao(tx)
	_, err := sd.Get(serviceID, "")
	if err == gorm.ErrRecordNotFound {
		respondNotFound(ctx, err)
		return
	} else if err != nil {
		respondInternalError(ctx, fmt.Errorf("database error: %w", err))
		return
	}
	if err := sd.Delete(serviceID, ""); err != nil {
		respondInternalError(ctx, fmt.Errorf("database error: %w", err))
		return
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		respondInternalError(ctx, fmt.Errorf("database error: %w", err))
		return
	}
	respondNoContent(ctx, nil)
}

func (sc ServiceController) UpdateService(ctx *gin.Context) {
	serviceID := uint32(getIntParam(ctx, "service_id"))
	if serviceID == 0 {
		respondBadRequest(ctx, fmt.Errorf("invalid service_id"))
	}
	params := struct {
		SyncType    string `json:"sync_type" json-validator:"optional"`
		BinRepo     string `json:"bin_repo" json-validator:"optional"`
		BinType     string `json:"bin_type" json-validator:"optional"`
		BinFtpAddr  string `json:"bin_ftp_addr" json-validator:"optional"`
		ConfRepo    string `json:"conf_repo" json-validator:"optional"`
		ConfType    string `json:"conf_type" json-validator:"optional"`
		ConfFtpAddr string `json:"conf_ftp_addr" json-validator:"optional"`
	}{}
	fieldMap, err := utils.BindJSONWithContext(&params, ctx)
	if err != nil {
		respondBadRequest(ctx, err)
		return
	}

	tx := database.GetDB().Begin()
	if err := tx.Error; err != nil {
		respondInternalError(ctx, fmt.Errorf("database error: %w", err))
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	sd := dao.NewServiceDao(tx)
	service, err := sd.Get(serviceID, "")
	if err == gorm.ErrRecordNotFound {
		respondNotFound(ctx, err)
		return
	} else if err != nil {
		respondInternalError(ctx, fmt.Errorf("database error: %w", err))
		return
	}

	if err := sd.Update(service, fieldMap); err != nil {
		respondInternalError(ctx, fmt.Errorf("database error: %w", err))
		return
	}
	if err := tx.Commit().Error; err != nil {
		respondInternalError(ctx, fmt.Errorf("database error: %w", err))
		return
	}
	respondSuccess(ctx, service)
}

func (sc ServiceController) AddServiceVersion(ctx *gin.Context) {
	respondSuccess(ctx, "todo")
}
