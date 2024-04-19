package router

import (
	"embed"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"html/template"
	"log"
	"net/http"
	"vland.live/app/global"
	"vland.live/app/initialize/validate"
	"vland.live/app/internal/api"
	"vland.live/app/middleware"
)

//go:embed html/*
var htmlFS embed.FS

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

	router.SetHTMLTemplate(template.Must(template.New("").ParseFS(htmlFS, "html/*")))

	aiManageRouter := router.Group("/ai-manage")
	aiManageRouter.GET("/chat", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", nil)
	})

	// ai
	{
		collector := api.NewAICollector()
		//aiManageRouter.POST("/v1/chat/completions", middleware.JwtAuth(), collector.V1ChatCompletions) // 验证版
		aiManageRouter.POST("/v1/chat/completions", collector.V1ChatCompletions) // 无验证版
	}

	// project
	{
		collector := api.NewProjectCollector()
		//projectsRouter := aiManageRouter.Group("/projects", middleware.JwtAuth()) // 验证版
		projectsRouter := aiManageRouter.Group("/projects") // 无验证版
		projectsRouter.POST("", collector.CreateAdmin)
		projectsRouter.PUT("/:id", collector.UpdateAdmin)
		projectsRouter.DELETE("/:id", collector.DeleteAdmin)
		projectsRouter.GET("", collector.SearchAdmin)
	}

	// prompt
	{
		collector := api.NewPromptCollector()
		// 管理后台用
		//promptsRouter := aiManageRouter.Group("/prompts", middleware.JwtAuth()) // 验证版
		promptsRouter := aiManageRouter.Group("/prompts") // 无验证版
		promptsRouter.POST("", collector.CreateAdmin)
		promptsRouter.PUT("/:id", collector.UpdateAdmin)
		promptsRouter.PUT("/:id/locked", collector.UpdateAdminLocked)
		promptsRouter.DELETE("/:id", collector.DeleteAdmin)
		promptsRouter.GET("", collector.SearchAdmin)
		promptsRouter.GET("/:id", collector.SearchAdminDetail)
		promptsRouter.POST("/:id/versions/save", collector.SaveAdminVersion)
		promptsRouter.POST("/:id/versions", collector.CreateAdminPromptVersion)
		promptsRouter.PATCH("/:id/versions/:version_id/is_production", collector.UpdateAdminPromptVersionIsProduction)
		promptsRouter.PATCH("/:id/versions/:version_id/name", collector.UpdateAdminPromptVersionName)
		promptsRouter.DELETE("/:id/versions/:version_id", collector.DeleteAdminPromptVersion)

		// 内部调用
		promptsInternalRouter := router.Group("/internal/prompts")
		promptsInternalRouter.POST("/:id", collector.Chat) // 内部调用，根据prompt请求
		// 外部调用
		promptsExternalRouter := aiManageRouter.Group("/prompts")
		promptsExternalRouter.POST("/:id/chat", collector.ChatExternal) // 外部调用聊天
	}

	log.Println("router init success")
	return router
}
