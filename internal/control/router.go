package control

import (
	domainService "im/internal/control/domain/service"
	errorService "im/internal/control/error/service"
	menuService "im/internal/control/menu/service"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "im/docs/control"
)

func NewRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	baseRouter := router.Group("/im_con")

	baseRouter.GET("/swagger/*any", RegistryDoc())
	errorRouter := baseRouter.Group("/error")
	{
		errorRouter.POST("/uploadErrlog", errorService.ErrorService.HandlerUpload)
		errorRouter.POST("/searchErrlog", errorService.ErrorService.HandlerSearch)
	}
	menuGroup := baseRouter.Group("/menu")
	{
		menuGroup.GET("", menuService.MenuService.MenuList)
		menuGroup.POST("", menuService.MenuService.MenuAdd)
		menuGroup.GET("/:id", menuService.MenuService.MenuGet)
		menuGroup.PUT("/:id", menuService.MenuService.MenuUpdate)
		menuGroup.DELETE("/:id", menuService.MenuService.MenuDelete)
	}
	domainGroup := baseRouter.Group("/domain")
	{
		domainGroup.POST("/domain_list", domainService.DomainSiteService.DomainList)
		domainGroup.POST("/domain_add", domainService.DomainSiteService.AddDomain)
		domainGroup.POST("/domain_remove", domainService.DomainSiteService.RemoveDomain)
		domainGroup.POST("/app_domain_list", domainService.DomainSiteService.AppDomainList)

		domainGroup.POST("/warning_list", domainService.DomainSiteService.WarningList)
		domainGroup.POST("/warning_add", domainService.DomainSiteService.AddWarning)
	}

	menuConfigGroup := baseRouter.Group("/menu_config")
	{
		menuConfigGroup.GET("", menuService.MenuService.HandlerGetMenuConfigTime)
	}
	return router
}

func RegistryDoc() gin.HandlerFunc {
	return ginSwagger.WrapHandler(swaggerFiles.Handler)
}
