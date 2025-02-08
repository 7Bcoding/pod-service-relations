package controller

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"pod-service-relations/dao"
	"pod-service-relations/database"
)

type PodServiceController struct{}

func (ps PodServiceController) GetPodService(ctx *gin.Context) {
	podName := ctx.Param("pod_name")
	//if podID == 0 {
	//	respondBadRequest(ctx, fmt.Errorf("bad requeset product id"))
	//	return
	//}
	params := struct {
		Verbose bool `form:"verbose"`
	}{}
	_ = ctx.ShouldBindQuery(&params)

	podServiceDao := dao.NewPodServiceDao(database.GetDB())

	product, err := podServiceDao.Get(podName)
	if err == gorm.ErrRecordNotFound {
		respondNotFound(ctx, err)
		return
	} else if err != nil {
		respondInternalError(ctx, err)
		return
	}
	respondSuccess(ctx, product)

}

func (ps PodServiceController) GetPodServiceList(ctx *gin.Context) {
	podServiceDao := dao.NewPodServiceDao(database.GetDB())
	params := struct {
		Verbose bool `form:"verbose"`
	}{}
	_ = ctx.ShouldBindQuery(&params)

	podServices, err := podServiceDao.ALL(nil)
	switch err {
	case nil, gorm.ErrRecordNotFound:
		respondSuccess(ctx, podServices)
	default:
		respondInternalError(ctx, fmt.Errorf("database error: %w", err))
	}
}

func (ps PodServiceController) AddPodService(ctx *gin.Context) {
	respondSuccess(ctx, "todo")
}

func (ps PodServiceController) UpdateProduct(ctx *gin.Context) {
	respondSuccess(ctx, "todo")
}

func (ps PodServiceController) DeleteProduct(ctx *gin.Context) {
	respondSuccess(ctx, "todo")
}
