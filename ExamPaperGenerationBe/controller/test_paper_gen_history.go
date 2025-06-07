package controller

import (
	"github.com/gin-gonic/gin"
	"graduation/entity"
	"graduation/mapper"
	"graduation/utils"
	"log"
	"net/http"
	"time"
)

// TestPaperGenHistoryMapper 定义测试试卷生成历史映射器接口
type TestPaperGenHistoryMapper interface {
	InsertTestPaperGenHistory(testPaperGenHistory entity.TestPaperGenHistory) (int64, error)
	QueryAllTestPaperGenHistory() ([]entity.TestPaperGenHistory, error)
	GetTestPaperNameByTestPaperUid(testPaperUID string) ([]string, error)
	UpdateTestPaperTime(testPaperUID string, date time.Time) (int64, error)
	DeleteTestPaperGenHistoryByTestPaperUid(testPaperUID string) (int64, error)
}

var testMapper = mapper.NewTestPaperGenHistoryGormMapper()

func GetAllTestPaperGenHistory(c *gin.Context) {
	info, err := testMapper.QueryAllTestPaperGenHistory()
	if err != nil {
		log.Println(err)
		utils.Make500Resp(err.Error())
	}
	resp := utils.Make200Resp("successfully get all test paper gen history", info)
	c.String(http.StatusOK, resp)
}
