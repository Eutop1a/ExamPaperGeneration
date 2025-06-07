package controller

import (
	"fmt"
	"graduation/entity"
	"graduation/mapper"
	"graduation/services"
	"graduation/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// 处理 /getQuestionGenHistoriesByTestPaperUid 请求
func GetQuestionGenHistoriesByTestPaperUid(c *gin.Context) {
	testPaperUid := c.Query("test_paper_uid")
	var questionGenHistories []entity.QuestionGenHistory
	mapper.DB.Where("test_paper_uid = ?", testPaperUid).Find(&questionGenHistories)
	resp := utils.Make200Resp("Success", questionGenHistories)
	c.String(http.StatusOK, resp)
}

// 处理 /deleteQuestionGenHistoryByTestPaperUid 请求
func DeleteQuestionGenHistoryByTestPaperUid(c *gin.Context) {
	testPaperUid := c.Query("test_paper_uid")
	var delQuestionCount int64
	mapper.DB.Where("test_paper_uid = ?", testPaperUid).Delete(&entity.QuestionGenHistory{}).Count(&delQuestionCount)
	var delTestPaperCount int64
	mapper.DB.Where("test_paper_uid = ?", testPaperUid).Delete(&entity.TestPaperGenHistory{}).Count(&delTestPaperCount)
	response := map[string]interface{}{
		"delQuestionCount":  delQuestionCount,
		"delTestPaperCount": delTestPaperCount,
	}
	resp := utils.Make200Resp("Success", response)
	c.String(http.StatusOK, resp)
}

// 处理 /updateQuestionGenHistory 请求
func UpdateQuestionGenHistory(c *gin.Context) {
	testPaperUid := c.Query("test_paper_uid")
	questionBankIdsStr := c.QueryArray("question_bank_id")
	var questionBankIds []int
	for _, idStr := range questionBankIdsStr {
		var id int
		fmt.Sscanf(idStr, "%d", &id)
		questionBankIds = append(questionBankIds, id)
	}

	var names []string
	mapper.DB.Table("test_paper_gen_histories").Where("test_paper_uid = ?", testPaperUid).Pluck("test_paper_name", &names)
	var testPaperName string
	if len(names) > 0 {
		testPaperName = names[0]
	}

	var questions []entity.QuestionBank
	for _, id := range questionBankIds {
		var question entity.QuestionBank
		if err := mapper.DB.Where("id = ?", id).First(&question).Error; err == nil {
			questions = append(questions, question)
		}
	}

	date := time.Now()
	var questionGenHistories []entity.QuestionGenHistory
	for _, q := range questions {
		questionGenHistory := entity.QuestionGenHistory{
			TestPaperUID:   testPaperUid,
			TestPaperName:  testPaperName,
			QuestionBankID: q.ID,
			Topic:          q.Topic,
			Answer:         q.Answer,
			TopicType:      q.TopicType,
			Score:          q.Score,
			Difficulty:     q.Difficulty,
			Chapter1:       q.Chapter1,
			Chapter2:       q.Chapter2,
			Label1:         q.Label1,
			Label2:         q.Label2,
			UpdateTime:     date,
		}
		questionGenHistories = append(questionGenHistories, questionGenHistory)
	}

	// 更新时间
	updateRes := mapper.DB.Model(&entity.TestPaperGenHistory{}).
		Where("test_paper_uid = ?", testPaperUid).
		Update("update_time", date)
	// 删除旧的
	deleteRes := mapper.DB.Where("test_paper_uid = ?", testPaperUid).Delete(&entity.QuestionGenHistory{})
	// 插入新的
	insertRes := mapper.DB.Create(&questionGenHistories)

	resp := utils.Make200Resp("Success", updateRes.RowsAffected+deleteRes.RowsAffected+insertRes.RowsAffected)
	c.String(http.StatusOK, resp)
}

// 处理 /reExportTestPaper 请求
func ReExportTestPaper(c *gin.Context) {
	testPaperUid := c.Query("test_paper_uid")
	var questionGenHistories []entity.QuestionGenHistory
	mapper.DB.Where("test_paper_uid = ?", testPaperUid).Find(&questionGenHistories)

	var questionBanks []entity.QuestionBank
	for _, item := range questionGenHistories {
		var question entity.QuestionBank
		if err := mapper.DB.Where("id = ?", item.QuestionBankID).First(&question).Error; err == nil {
			questionBanks = append(questionBanks, question)
		}
	}

	// 分类题目
	var tktQuestions, xztQuestions, pdtQuestions, jdtQuestions []entity.QuestionBank
	for _, q := range questionBanks {
		switch q.TopicType {
		case "填空题":
			tktQuestions = append(tktQuestions, q)
		case "选择题":
			xztQuestions = append(xztQuestions, q)
		case "判断题":
			pdtQuestions = append(pdtQuestions, q)
		case "简答题":
			jdtQuestions = append(jdtQuestions, q)
		}
	}

	var totalScore float64
	totalCount := 0
	contents := ""
	questionNumber := 1 // 统一题号，从1开始

	// 添加选择题
	if len(xztQuestions) > 0 {
		// 计算每小题分数
		scorePerQuestion := xztQuestions[0].Score
		contents += fmt.Sprintf("[SECTION_TITLE]一、选择题（本大题共%d小题，每小题%.1f分，共%.1f分）[/SECTION_TITLE]\r\r", len(xztQuestions), scorePerQuestion, float64(len(xztQuestions))*scorePerQuestion)
		for _, q := range xztQuestions {
			totalScore += q.Score
			scoreStr := formatScore(q.Score)
			contents = fmt.Sprintf("%s%d、（本题%s分）%s", contents, questionNumber, scoreStr, q.Topic)

			// 如果有图片，添加图片标记
			if q.TopicImagePath != "" {
				contents += fmt.Sprintf(" [IMAGE:%s]", q.TopicImagePath)
			}

			contents += "\r\r[QUESTION_END]\r\r" // 使用特殊标记分隔题目
			questionNumber++
			totalCount++
		}
	}

	// 添加填空题
	if len(tktQuestions) > 0 {
		// 计算每小题分数
		scorePerQuestion := tktQuestions[0].Score
		contents += fmt.Sprintf("[SECTION_TITLE]二、填空题（本大题共%d小题，每小题%.1f分，共%.1f分）[/SECTION_TITLE]\r\r", len(tktQuestions), scorePerQuestion, float64(len(tktQuestions))*scorePerQuestion)
		for _, q := range tktQuestions {
			totalScore += q.Score
			scoreStr := formatScore(q.Score)
			contents = fmt.Sprintf("%s%d、（本题%s分）%s", contents, questionNumber, scoreStr, q.Topic)

			// 如果有图片，添加图片标记
			if q.TopicImagePath != "" {
				contents += fmt.Sprintf(" [IMAGE:%s]", q.TopicImagePath)
			}

			contents += "\r\r[QUESTION_END]\r\r" // 使用特殊标记分隔题目
			questionNumber++
			totalCount++
		}
	}

	// 添加判断题
	if len(pdtQuestions) > 0 {
		// 计算每小题分数
		scorePerQuestion := pdtQuestions[0].Score
		contents += fmt.Sprintf("[SECTION_TITLE]三、判断题（本大题共%d小题，每小题%.1f分，共%.1f分）[/SECTION_TITLE]\r\r", len(pdtQuestions), scorePerQuestion, float64(len(pdtQuestions))*scorePerQuestion)
		for _, q := range pdtQuestions {
			totalScore += q.Score
			scoreStr := formatScore(q.Score)
			contents = fmt.Sprintf("%s%d、（本题%s分）%s", contents, questionNumber, scoreStr, q.Topic)

			// 如果有图片，添加图片标记
			if q.TopicImagePath != "" {
				contents += fmt.Sprintf(" [IMAGE:%s]", q.TopicImagePath)
			}

			contents += "\r\r[QUESTION_END]\r\r" // 使用特殊标记分隔题目
			questionNumber++
			totalCount++
		}
	}

	// 添加简答题
	if len(jdtQuestions) > 0 {
		// 计算每小题分数
		scorePerQuestion := jdtQuestions[0].Score
		contents += fmt.Sprintf("[SECTION_TITLE]四、简答题（本大题共%d小题，每小题%.1f分，共%.1f分）[/SECTION_TITLE]\r\r", len(jdtQuestions), scorePerQuestion, float64(len(jdtQuestions))*scorePerQuestion)
		for _, q := range jdtQuestions {
			totalScore += q.Score
			scoreStr := formatScore(q.Score)
			contents = fmt.Sprintf("%s%d、（本题%s分）%s", contents, questionNumber, scoreStr, q.Topic)

			// 如果有图片，添加图片标记
			if q.TopicImagePath != "" {
				contents += fmt.Sprintf(" [IMAGE:%s]", q.TopicImagePath)
			}

			contents += "\r\r[QUESTION_END]\r\r" // 使用特殊标记分隔题目
			questionNumber++
			totalCount++
		}
	}

	mapData := map[string]string{
		"total_score": fmt.Sprintf("%s", formatScore(totalScore)),
		"total_count": fmt.Sprintf("%d", totalCount),
		"contents":    contents,
	}

	// 使用新的 WordExporterGooxml 导出 Word 文档
	wE := services.NewWordExporterGooxml(mapData)
	file, err := wE.ExportTestPaper(1)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	downloadFile(c, file)
}

// 处理 /exportAnswer 请求
func ExportAnswer(c *gin.Context) {
	testPaperUid := c.Query("test_paper_uid")
	var questionGenHistories []entity.QuestionGenHistory
	mapper.DB.Where("test_paper_uid = ?", testPaperUid).Find(&questionGenHistories)

	var questionBanks []entity.QuestionBank
	for _, item := range questionGenHistories {
		var question entity.QuestionBank
		if err := mapper.DB.Where("id = ?", item.QuestionBankID).First(&question).Error; err == nil {
			questionBanks = append(questionBanks, question)
		}
	}

	// 分类题目
	var tktQuestions, xztQuestions, pdtQuestions, jdtQuestions []entity.QuestionBank
	for _, q := range questionBanks {
		switch q.TopicType {
		case "填空题":
			tktQuestions = append(tktQuestions, q)
		case "选择题":
			xztQuestions = append(xztQuestions, q)
		case "判断题":
			pdtQuestions = append(pdtQuestions, q)
		case "简答题":
			jdtQuestions = append(jdtQuestions, q)
		}
	}

	var totalScore float64
	totalCount := 0
	contents := ""
	questionNumber := 1 // 统一题号，从1开始

	// 添加选择题答案
	if len(xztQuestions) > 0 {
		// 计算每小题分数
		scorePerQuestion := xztQuestions[0].Score
		contents += fmt.Sprintf("[SECTION_TITLE]一、选择题答案（本大题共%d小题，每小题%.1f分，共%.1f分）[/SECTION_TITLE]\r\r", len(xztQuestions), scorePerQuestion, float64(len(xztQuestions))*scorePerQuestion)
		for _, q := range xztQuestions {
			totalScore += q.Score
			scoreStr := formatScore(q.Score)
			contents = fmt.Sprintf("%s%d、（本题%s分）%s", contents, questionNumber, scoreStr, q.Answer)
			contents += "\r\r[QUESTION_END]\r\r" // 使用特殊标记分隔题目
			questionNumber++
			totalCount++
		}
	}

	// 添加填空题答案
	if len(tktQuestions) > 0 {
		// 计算每小题分数
		scorePerQuestion := tktQuestions[0].Score
		contents += fmt.Sprintf("[SECTION_TITLE]二、填空题答案（本大题共%d小题，每小题%.1f分，共%.1f分）[/SECTION_TITLE]\r\r", len(tktQuestions), scorePerQuestion, float64(len(tktQuestions))*scorePerQuestion)
		for _, q := range tktQuestions {
			totalScore += q.Score
			scoreStr := formatScore(q.Score)
			contents = fmt.Sprintf("%s%d、（本题%s分）%s", contents, questionNumber, scoreStr, q.Answer)
			contents += "\r\r[QUESTION_END]\r\r" // 使用特殊标记分隔题目
			questionNumber++
			totalCount++
		}
	}

	// 添加判断题答案
	if len(pdtQuestions) > 0 {
		// 计算每小题分数
		scorePerQuestion := pdtQuestions[0].Score
		contents += fmt.Sprintf("[SECTION_TITLE]三、判断题答案（本大题共%d小题，每小题%.1f分，共%.1f分）[/SECTION_TITLE]\r\r", len(pdtQuestions), scorePerQuestion, float64(len(pdtQuestions))*scorePerQuestion)
		for _, q := range pdtQuestions {
			totalScore += q.Score
			scoreStr := formatScore(q.Score)
			contents = fmt.Sprintf("%s%d、（本题%s分）%s", contents, questionNumber, scoreStr, q.Answer)
			contents += "\r\r[QUESTION_END]\r\r" // 使用特殊标记分隔题目
			questionNumber++
			totalCount++
		}
	}

	// 添加简答题答案
	if len(jdtQuestions) > 0 {
		// 计算每小题分数
		scorePerQuestion := jdtQuestions[0].Score
		contents += fmt.Sprintf("[SECTION_TITLE]四、简答题答案（本大题共%d小题，每小题%.1f分，共%.1f分）[/SECTION_TITLE]\r\r", len(jdtQuestions), scorePerQuestion, float64(len(jdtQuestions))*scorePerQuestion)
		for _, q := range jdtQuestions {
			totalScore += q.Score
			scoreStr := formatScore(q.Score)
			contents = fmt.Sprintf("%s%d、（本题%s分）%s", contents, questionNumber, scoreStr, q.Answer)
			contents += "\r\r[QUESTION_END]\r\r" // 使用特殊标记分隔题目
			questionNumber++
			totalCount++
		}
	}

	mapData := map[string]string{
		"total_score": fmt.Sprintf("%s", formatScore(totalScore)),
		"total_count": fmt.Sprintf("%d", totalCount),
		"contents":    contents,
	}

	// 使用新的 WordExporterGooxml 导出 Word 文档
	wE := services.NewWordExporterGooxml(mapData)
	file, err := wE.ExportTestPaper(2)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	downloadFile(c, file)
}

// 格式化分数输出
func formatScore(score float64) string {
	if score == float64(int(score)) {
		return fmt.Sprintf("%d", int(score))
	}
	return fmt.Sprintf("%.1f", score)
}
