package router

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"log"
	"vland.live/app/global"
	"vland.live/app/initialize/validate"
	"vland.live/app/internal/api"
	"vland.live/app/middleware"
)

////go:embed html/*
//var htmlFS embed.FS

func Init() *gin.Engine {

	router := gin.New()

	// 注册验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := v.RegisterValidation("objectId", validate.ObjectId); err != nil {
			log.Panicln("注册【objectId】验证器失败", err)
		}
		if err := v.RegisterValidation("sliceRepeat", validate.SliceRepeat); err != nil {
			log.Panicln("注册【sliceRepeat】验证器失败", err)
		}
	}

	router.Use(middleware.Cors(), middleware.GinLogger(global.Logger), middleware.GinRecovery(true, global.Logger)) // 处理跨域

	//router.SetHTMLTemplate(template.Must(template.New("").ParseFS(htmlFS, "html/*")))

	hollywoodRouter := router.Group("/hollywood")

	// Screenplay
	{
		collector := api.NewScreenplayCollector()
		_router := hollywoodRouter.Group("/screenplays")
		_router.POST("", collector.CreateAdmin)
		_router.PUT("/:id", collector.UpdateAdmin)
		_router.DELETE("/:id", collector.DeleteAdmin)

		_router.GET("", collector.Search)
	}

	// Member
	{
		collector := api.NewMemberCollector()
		_router := hollywoodRouter.Group("/members")
		_router.POST("/enter", collector.Enter)
		_router.PATCH("/screenplay", collector.UpdateScreenplayReq)
		_router.PATCH("/refresh/attributes", collector.RefreshAttributes)
		_router.POST("/start-or-nextround", collector.StartOrNextRound)
		_router.POST("/outcome", collector.ClickButtonOutcome)
		_router.POST("/outcome-news", collector.ClickButtonNewsByOutcome)
		_router.GET("", collector.Search)
	}
	log.Println("router init success")
	return router
}
