package controller

import (
	"fmt"
	"graduation/entity"
	"graduation/mapper"
	"graduation/services"
	"graduation/utils"

	"github.com/gin-contrib/sessions"

	"net/http"
	"os"
	"time"

	"encoding/base64"
	"io/ioutil"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RandomSelectRequest 随机抽题请求
type RandomSelectRequest struct {
	SelectedTopicIds         []int                                       `json:"selectedTopicIds"`
	AverageDifficulty        float64                                     `json:"averageDifficulty"`
	GenerateRange            []string                                    `json:"generateRange"`
	KnowledgeWeights         []services.KnowledgePointWeight             `json:"knowledgeWeights"`
	QuestionTypeRequirements map[string]services.QuestionTypeRequirement `json:"questionTypeRequirements"`
	IterationsNum            int                                         `json:"iterationsNum"` // 添加迭代次数参数
}

// QuestionTypeRequirement 题型要求
type QuestionTypeRequirement struct {
	MinCount    int     `json:"minCount"`
	TargetScore float64 `json:"targetScore"`
}

// RandomSelect 随机选题
func RandomSelect(c *gin.Context) {
	var request RandomSelectRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid request parameters", "error": err.Error()})
		return
	}

	// 获取所有题目
	questions, err := getAllQuestions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to get questions", "error": err.Error()})
		return
	}

	// 过滤题目
	filteredQuestions := filterQuestionsByTopicAndRange(questions, request.SelectedTopicIds, request.GenerateRange)
	if len(filteredQuestions) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "No questions found for the selected topics and range"})
		return
	}
	session := sessions.Default(c)
	username := session.Get("username")
	var userInfo entity.User
	if err := mapper.DB.Model(&entity.User{}).Where("username = ?", username).First(&userInfo).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Username not exist", "error": err.Error()})
		return
	}
	// 获取需要排除的题目ID
	excludedQuestionIds := services.GetExcludedQuestionIds(userInfo.SimilarityThreshold)
	//excludedQuestionIds = nil
	// 转换题型要求
	questionTypeRequirements := make(map[string]services.QuestionTypeRequirement)
	for qType, req := range request.QuestionTypeRequirements {
		questionTypeRequirements[qType] = services.QuestionTypeRequirement{
			MinCount:    req.MinCount,
			TargetScore: req.TargetScore,
		}
	}

	// 随机选题
	selectedQuestions := services.WeightedRandomSelect(
		filteredQuestions,
		request.AverageDifficulty,
		request.KnowledgeWeights,
		questionTypeRequirements,
		request.SelectedTopicIds,
		excludedQuestionIds,
	)

	// 按题型分类
	var tktList, xztList, pdtList, jdtList []entity.QuestionBank
	var totalScore float64

	for _, q := range selectedQuestions {
		totalScore += q.Score
		switch q.TopicType {
		case "填空题":
			tktList = append(tktList, q)
		case "选择题":
			xztList = append(xztList, q)
		case "判断题":
			pdtList = append(pdtList, q)
		case "简答题":
			jdtList = append(jdtList, q)
		}
	}

	// 转换图片为base64
	convertImagesToBase64(tktList)
	convertImagesToBase64(xztList)
	convertImagesToBase64(pdtList)
	convertImagesToBase64(jdtList)

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Success",
		"data": gin.H{
			"TKTList":    tktList,
			"XZTList":    xztList,
			"PDTList":    pdtList,
			"JDTList":    jdtList,
			"TotalScore": totalScore,
		},
	})
}

// getAllQuestions 获取所有题目
func getAllQuestions() ([]entity.QuestionBank, error) {
	var questions []entity.QuestionBank
	result := mapper.DB.Find(&questions)
	if result.Error != nil {
		return nil, result.Error
	}
	return questions, nil
}

// filterQuestionsByTopicAndRange 根据主题ID和范围过滤题目
func filterQuestionsByTopicAndRange(questions []entity.QuestionBank, selectedTopicIds []int, generateRange []string) []entity.QuestionBank {
	var filteredQuestions []entity.QuestionBank

	// 创建已选主题ID的映射，用于快速查找
	selectedTopicIdsMap := make(map[int]bool)
	for _, id := range selectedTopicIds {
		selectedTopicIdsMap[id] = true
	}

	// 创建生成范围的映射，用于快速查找
	generateRangeMap := make(map[string]bool)
	for _, range_ := range generateRange {
		generateRangeMap[range_] = true
	}

	// 过滤题目
	for _, q := range questions {
		// 如果生成范围不为空，且题目不在生成范围内，跳过
		if len(generateRange) > 0 && !generateRangeMap[q.Label1] {
			continue
		}

		filteredQuestions = append(filteredQuestions, q)
	}

	return filteredQuestions
}

// GeneticSelect 遗传算法抽题
func GeneticSelect(c *gin.Context) {
	var payload RandomSelectRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取所有题目
	questions, err := getAllQuestions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 根据选中的知识点和生成范围过滤题目
	filteredQuestions := filterQuestionsByTopicAndRange(questions, payload.SelectedTopicIds, payload.GenerateRange)
	session := sessions.Default(c)
	username := session.Get("username")
	var userInfo entity.User
	if err := mapper.DB.Model(&entity.User{}).Where("username = ?", username).First(&userInfo).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Username not exist", "error": err.Error()})
		return
	}
	// 获取已选过的题目ID
	excludedQuestionIds := services.GetExcludedQuestionIds(userInfo.SimilarityThreshold)

	// 转换题型要求格式
	questionTypeRequirements := make(map[string]services.QuestionTypeRequirement)
	for qType, req := range payload.QuestionTypeRequirements {
		questionTypeRequirements[qType] = services.QuestionTypeRequirement{
			MinCount: req.MinCount,
		}
	}

	// 创建遗传算法实例并运行
	genIter := services.NewGeneticIteration(
		filteredQuestions,
		payload.AverageDifficulty,
		payload.KnowledgeWeights,
		questionTypeRequirements,
		payload.SelectedTopicIds, // 手动选择的题目ID列表
		excludedQuestionIds,
		payload.IterationsNum, // 传入迭代次数
	)
	selectedQuestions := genIter.Run()

	// 按题型分类
	var TKTList, XZTList, PDTList, JDTList []entity.QuestionBank
	for _, q := range selectedQuestions {
		switch q.TopicType {
		case "填空题":
			TKTList = append(TKTList, q)
		case "选择题":
			XZTList = append(XZTList, q)
		case "判断题":
			PDTList = append(PDTList, q)
		case "简答题":
			JDTList = append(JDTList, q)
		}
	}

	// 转换图片为base64
	convertImagesToBase64(TKTList)
	convertImagesToBase64(XZTList)
	convertImagesToBase64(PDTList)
	convertImagesToBase64(JDTList)

	response := map[string]interface{}{
		"TKTList":  TKTList,
		"XZTList":  XZTList,
		"PDTList":  PDTList,
		"JDTList":  JDTList,
		"variance": genIter.Variance,
	}
	// 返回结果
	c.String(http.StatusOK, utils.Make200Resp("success", response))
}

// QuestionGen 处理 /questionGen 请求
func QuestionGen(c *gin.Context) {
	session := sessions.Default(c)
	username, ok := session.Get("username").(string)
	fmt.Println("username: ", username)
	if !ok || username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username is missing"})
		return
	}

	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request body: %v", err)})
		return
	}

	questionIdList := getIntList(payload, "questionIdList")
	TKTIdList := getIntList(payload, "TKTIdList")
	XZTIdList := getIntList(payload, "XZTIdList")
	PDTIdList := getIntList(payload, "PDTIdList")
	JDTIdList := getIntList(payload, "JDTIdList")
	testPaperName := getString(payload, "testPaperName")

	questionBanks := getQuestionBanks(questionIdList, TKTIdList, XZTIdList, PDTIdList, JDTIdList)
	if len(questionBanks) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid questions found"})
		return
	}

	wE := services.NewWordGenerator()
	str, _ := wE.GenerateTestPaper(questionBanks, testPaperName, username)
	fmt.Println(str)

	response := map[string]interface{}{
		"status":  200,
		"message": "Success",
	}
	c.String(http.StatusOK, utils.Make200Resp("Success", response))
}

// QuestionGen2 处理 /questionGen2 请求
func QuestionGen2(c *gin.Context) {
	session := sessions.Default(c)
	username, ok := session.Get("username").(string)
	if !ok || username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username is missing"})
		return
	}

	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request body: %v", err)})
		return
	}

	questionIdList := getIntList(payload, "questionIdList")
	TKTIdList := getIntList(payload, "TKTIdList")
	XZTIdList := getIntList(payload, "XZTIdList")
	PDTIdList := getIntList(payload, "PDTIdList")
	JDTIdList := getIntList(payload, "JDTIdList")
	testPaperName := getString(payload, "testPaperName")

	questionBanks := getQuestionBanks(questionIdList, TKTIdList, XZTIdList, PDTIdList, JDTIdList)
	if len(questionBanks) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid questions found"})
		return
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to export Word document: %v", err)})
		return
	}

	logHistory(questionBanks, testPaperName, username, file)
	downloadFile(c, file)
}

// 处理 /getFile 请求
func GetFile(c *gin.Context) {
	genWord := services.NewWordGenerator()
	file := genWord.GetFile()
	downloadFile(c, file)
}

// 从数据库中获取题目列表
func getQuestionBanks(ids ...[]int) []entity.QuestionBank {
	var allIds []int
	for _, idList := range ids {
		allIds = append(allIds, idList...)
	}

	var questionBanks []entity.QuestionBank
	mapper.DB.Where("id IN ?", allIds).Find(&questionBanks)
	return questionBanks
}

// 记录历史记录
func logHistory(questionBanks []entity.QuestionBank, testPaperName, username string, file *os.File) {
	date := time.Now()
	uid := fmt.Sprintf("%s_%s_%d", file.Name(), uuid.New().String(), date.Unix())

	testPaperGenHistory := entity.TestPaperGenHistory{
		TestPaperUID:      uid,
		TestPaperName:     testPaperName,
		QuestionCount:     len(questionBanks),
		AverageDifficulty: calculateAverageDifficulty(questionBanks),
		UpdateTime:        date,
		Username:          username,
	}
	mapper.DB.Create(&testPaperGenHistory)
	var questionGenHistoryList []entity.QuestionGenHistory
	for _, q := range questionBanks {
		questionGenHistory := entity.QuestionGenHistory{
			TestPaperUID:    uid,
			TestPaperName:   testPaperName,
			QuestionBankID:  q.ID,
			Topic:           q.Topic,
			TopicMaterialID: q.TopicMaterialID,
			Answer:          q.Answer,
			TopicType:       q.TopicType,
			Score:           q.Score,
			Difficulty:      q.Difficulty,
			Chapter1:        q.Chapter1,
			Chapter2:        q.Chapter2,
			Label1:          q.Label1,
			Label2:          q.Label2,
			UpdateTime:      date,
		}
		questionGenHistoryList = append(questionGenHistoryList, questionGenHistory)
	}
	mapper.DB.Create(&questionGenHistoryList)
}

// 计算平均难度
func calculateAverageDifficulty(questionBanks []entity.QuestionBank) float64 {
	if len(questionBanks) == 0 {
		return 0
	}
	var totalDifficulty int
	for _, q := range questionBanks {
		totalDifficulty += q.Difficulty
	}
	return float64(totalDifficulty) / float64(len(questionBanks))
}

// 下载文件
func downloadFile(c *gin.Context, file *os.File) {
	if file == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error getting file stats: %v", err)})
		return
	}

	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=TestPaperExport_%d.docx", time.Now().Unix()))
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")
	c.Header("Last-Modified", time.Now().String())
	c.Header("Content-Length", fmt.Sprintf("%d", stat.Size()))
	c.Header("Content-Type", "application/octet-stream")

	c.File(file.Name())
}

// 从 map 中获取 int 列表
func getIntList(m map[string]interface{}, key string) []int {
	var list []int
	if val, ok := m[key].([]interface{}); ok {
		for _, v := range val {
			if num, ok := v.(float64); ok {
				list = append(list, int(num))
			}
		}
	}
	return list
}

// 从 map 中获取 string 列表
func getStringList(m map[string]interface{}, key string) []string {
	var list []string
	if val, ok := m[key].([]interface{}); ok {
		for _, v := range val {
			if str, ok := v.(string); ok {
				list = append(list, str)
			}
		}
	}
	return list
}

// 从 map 中获取 float64 类型的值
func getFloat64(m map[string]interface{}, key string) float64 {
	if val, ok := m[key].(float64); ok {
		return val
	}
	return 0
}

// 从 map 中获取 string 类型的值
func getString(m map[string]interface{}, key string) string {
	if val, ok := m[key].(string); ok {
		return val
	}
	return ""
}

// convertImagesToBase64 将题目列表中的图片转换为base64格式
func convertImagesToBase64(questions []entity.QuestionBank) {
	for i := range questions {
		if questions[i].TopicImagePath != "" {
			// 读取图片文件
			imgData, err := ioutil.ReadFile(questions[i].TopicImagePath)
			if err != nil {
				log.Printf("Error reading image file %s: %v", questions[i].TopicImagePath, err)
				continue
			}

			// 转换为base64
			base64Str := base64.StdEncoding.EncodeToString(imgData)
			// 添加base64前缀
			questions[i].TopicImagePath = "data:image/jpeg;base64," + base64Str
		}
	}
}
