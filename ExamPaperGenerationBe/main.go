package main

import (
	"encoding/gob"
	"graduation/component"
	"graduation/config"
	"graduation/controller"
	"log"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func setupMiddleware(r *gin.Engine) {
	gob.Register(time.Time{})
	store := cookie.NewStore([]byte("secret"))
	store.Options(sessions.Options{
		MaxAge: 3600, // 1 hour
	})
	r.Use(sessions.Sessions("mysession", store))
	r.Use(config.Cors())
	// 添加登录拦截器
	r.Use(component.LoginHandlerInterceptor())
}

func registerAuthRoutes(r *gin.Engine) {
	r.GET("/permission_denied", controller.PermissionDenied)
	r.GET("/getLoginStatus", controller.GetLoginStatus)
	r.POST("/login", controller.Login)
	r.POST("/logout", controller.Logout)
	r.POST("/registered", controller.Registered)
	r.PUT("/similarity", controller.SetSimilarityThreshold)
}

func registerUserManagementRoutes(r *gin.Engine) {
	r.GET("/getApplyUser", controller.GetApplyUser)
	r.GET("/getAllUser", controller.GetAllUser)
	r.GET("/deleteUser", controller.DeleteUser)
	r.GET("/passApply", controller.PassApply)
	r.GET("/deleteApply", controller.DeleteApply)
}

func registerQuestionBankRoutes(r *gin.Engine, qBan *controller.QuestionBankController) {
	// Question bank management
	r.GET("/getAllQuestionBank", qBan.GetAllQuestionBank)
	r.GET("/getQuestionBank", qBan.GetQuestionBank)
	r.GET("/getTopicType", qBan.GetTopicType)
	r.GET("/searchQuestionByTopic", qBan.SearchQuestionByTopic)
	r.POST("/insertSingleQuestionBank", qBan.InsertSingleQuestionBank)
	r.POST("/insertSingleQuestionBankWithImg", qBan.InsertSingleQuestionBankWithImg)
	r.GET("/deleteSingleQuestionBank", qBan.DeleteSingleQuestionBank)
	r.GET("/getQuestionBankById", qBan.GetQuestionBankById)
	r.POST("/updateQuestionBankById", qBan.UpdateQuestionBankById)
	r.POST("/upload", qBan.UploadFile)
	r.GET("/getEachChapterCount", qBan.GetEachChapterCount)
	r.GET("/getEachScoreCount", qBan.GetEachScoreCount)
}

func registerQuestionGenRoutes(r *gin.Engine) {
	// Question generation and history
	r.GET("/getQuestionGenHistoriesByTestPaperUid", controller.GetQuestionGenHistoriesByTestPaperUid)
	r.GET("/deleteQuestionGenHistoryByTestPaperUid", controller.DeleteQuestionGenHistoryByTestPaperUid)
	r.POST("/updateQuestionGenHistory", controller.UpdateQuestionGenHistory)
	r.GET("/reExportTestPaper", controller.ReExportTestPaper)
	r.GET("/exportAnswer", controller.ExportAnswer)
	r.GET("/getAllTestPaperGenHistory", controller.GetAllTestPaperGenHistory)

	// Question metadata
	r.GET("/getAllQuestionLabels", controller.GetAllQuestionLabels)
	r.GET("/getDistinctChapter1", controller.GetDistinctChapter1)
	r.GET("/getDistinctChapter2", controller.GetDistinctChapter2)
	r.GET("/getChapter2ByChapter1", controller.GetChapter2ByChapter1)
	r.GET("/getDistinctLabel1", controller.GetDistinctLabel1)
	r.GET("/getDistinctLabel2", controller.GetDistinctLabel2)

	// Question generation algorithms
	r.POST("/randomSelect", controller.RandomSelect)
	r.POST("/geneticSelect", controller.GeneticSelect)
	r.POST("/questionGen", controller.QuestionGen)
	r.POST("/questionGen2", controller.QuestionGen2)
	r.POST("/getFile", controller.GetFile)
}

func registerLabelRoutes(r *gin.Engine) {
	// 题目标签相关路由
	labelsGroup := r.Group("/labels")
	{
		labelsGroup.POST("", controller.CreateLabel)
		labelsGroup.PUT("/:id", controller.UpdateLabel)
		labelsGroup.DELETE("/:id", controller.DeleteLabel)
	}
}

func main() {
	r := gin.Default()

	// Setup middleware
	setupMiddleware(r)

	// Register routes
	registerAuthRoutes(r)
	registerUserManagementRoutes(r)
	qBan := controller.NewQuestionBankController()
	registerQuestionBankRoutes(r, qBan)
	registerQuestionGenRoutes(r)
	registerLabelRoutes(r)

	// Start server
	addr := ":8081"
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// 说明：最后组成的试卷的整体难度怎么算出来的

// 必须实现
// 组卷：对于不同章节要有权重和数量限制，总分控制为100
// 题目支持图片
// 添加对于 大小知识点的 crud
// 试卷相似度手动指定

// 根据实际情况
// 试卷生成模板可以手动导入
// 题目重复情况的处理，相似度 90% 为不同题目
