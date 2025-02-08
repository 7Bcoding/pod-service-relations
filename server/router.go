package server

import (
	"github.com/gin-gonic/gin"

	"pod-service-relations/config"
	"pod-service-relations/controller"
)

func newRoute(conf config.ServerConfig) *gin.Engine {
	// production. For debug set gin.SetMode(gin.DebugMode)
	//gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	// register default middlewares
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	// register custom or 3rd party middlewares
	//store, _ := redis.NewStore(conf.Session.Redis.MaxIdle, "tcp", conf.Session.Redis.URI, "", []byte(conf.Session.Secret))
	//router.Use(sessions.Sessions(conf.Session.Name, store))
	router.Use(logMW())

	pc := controller.ProductController{}
	sc := controller.ServiceController{}
	ps := controller.PodServiceController{}

	routerRead := router.Group("api/v1", controller.PaginationMiddleware)
	{
		// service & service version
		routerRead.GET("service/:id", sc.GetService)
		routerRead.GET("services", sc.GetServices)
		// product
		routerRead.GET("product/:id", pc.GetProduct)
		routerRead.GET("products", pc.GetProducts)
		routerRead.GET("pod_services", ps.GetPodServiceList)
	}

	//routerWrite := router.Group("api/v1", loginMW())
	//{
	//	// service & service version
	//	routerWrite.POST("services", sc.AddService)
	//	routerWrite.PATCH("service/:id", sc.UpdateService)
	//	routerWrite.DELETE("services/:id", sc.DeleteService)
	//	routerWrite.POST("services/version", sc.AddServiceVersion)
	//
	//	// product
	//	routerWrite.POST("products", pc.AddProduct)
	//	routerWrite.PATCH("products/:id", pc.UpdateProduct)
	//	routerWrite.DELETE("products/:id", pc.DeleteProduct)
	//}

	return router
}
