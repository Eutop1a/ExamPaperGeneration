package services

import (
	"graduation/entity"
	"math"
	"math/rand"
	"sort"
	"time"
)

// RandomSelectTopic 保持原有函数签名不变
func RandomSelectTopic(dataSource []entity.QuestionBank, targetDiff float64, selectCount int) ([]entity.QuestionBank, []int) {
	// 初始化结果列表
	result := make([]entity.QuestionBank, 0, selectCount)
	questionIds := make([]int, 0, selectCount)

	// 参数校验
	if len(dataSource) == 0 || targetDiff < 1.0 || targetDiff > 5.0 || selectCount <= 0 {
		return result, questionIds
	}

	// 预排序，便于后续处理
	sort.SliceStable(dataSource, func(i, j int) bool {
		return dataSource[i].Difficulty < dataSource[j].Difficulty
	})

	// 创建工作集合副本
	workingSet := make([]entity.QuestionBank, len(dataSource))
	copy(workingSet, dataSource)

	// 计算历史题目权重
	historyWeights := calculateHistoryWeights(workingSet)

	// 清空当前试卷的已选题目
	updateSelectedQuestions(nil)

	// 选题过程
	for remaining := selectCount; remaining > 0 && len(workingSet) > 0; remaining-- {
		selected, ok := selectQuestionV2(workingSet, targetDiff, historyWeights)
		if !ok {
			break
		}
		result = append(result, selected)
		questionIds = append(questionIds, selected.ID)
		workingSet = removeQuestion(workingSet, selected.ID)

		// 更新当前试卷的已选题目
		updateSelectedQuestions(result)
	}

	return result, questionIds
}

// 计算历史题目权重
func calculateHistoryWeights(questions []entity.QuestionBank) map[int]float64 {
	weights := make(map[int]float64)
	currentTime := time.Now()

	for _, q := range questions {
		if lastUsed, exists := getHistoryQuestionTime(q.ID); exists {
			// 计算时间衰减因子（两年内线性衰减）
			yearsDiff := currentTime.Sub(lastUsed).Hours() / (24 * 365)
			if yearsDiff < 2 {
				weights[q.ID] = 1.0 - yearsDiff/2 // 衰减系数从1.0到0
			}
		}
	}
	return weights
}

// 计算两个题目的相似度
func calculateQuestionSimilarity(q1, q2 entity.QuestionBank) float64 {
	// 1. 章节相似度 (30%)
	chapterSim := 0.0
	if q1.Chapter1 == q2.Chapter1 {
		chapterSim += 0.15
	}
	if q1.Chapter2 == q2.Chapter2 {
		chapterSim += 0.15
	}

	// 2. 标签相似度 (30%)
	labelSim := 0.0
	if q1.Label1 == q2.Label1 {
		labelSim += 0.15
	}
	if q1.Label2 == q2.Label2 {
		labelSim += 0.15
	}

	// 3. 难度相似度 (20%)
	diffSim := 1.0 - math.Abs(float64(q1.Difficulty-q2.Difficulty))/5.0
	if diffSim < 0 {
		diffSim = 0
	}

	// 4. 题目类型相似度 (20%)
	typeSim := 0.0
	if q1.TopicType == q2.TopicType {
		typeSim = 0.2
	}

	return chapterSim + labelSim + diffSim*0.2 + typeSim
}

func selectQuestionV2(candidates []entity.QuestionBank, targetDiff float64,
	historyWeights map[int]float64) (entity.QuestionBank, bool) {

	if len(candidates) == 0 {
		return entity.QuestionBank{}, false
	}

	// 使用带权重的随机选择
	type candidate struct {
		q     entity.QuestionBank
		score float64
	}

	// 构建候选列表并计算得分
	var candidatesPool []candidate
	var minScore = math.MaxFloat64

	for _, q := range candidates {
		// 1. 难度差异分数 (40%)
		difficultyDelta := math.Abs(float64(q.Difficulty) - targetDiff)
		difficultyScore := difficultyDelta * 0.4

		// 2. 历史题目惩罚 (30%)
		historyScore := historyWeights[q.ID] * 0.3

		// 3. 与已选题目的相似度惩罚 (30%)
		similarityScore := 0.0
		for _, selected := range getSelectedQuestions() {
			sim := calculateQuestionSimilarity(q, selected)
			if sim > similarityScore {
				similarityScore = sim
			}
		}
		similarityScore *= 0.3

		// 总分数
		totalScore := difficultyScore + historyScore + similarityScore

		if totalScore < minScore {
			minScore = totalScore
		}
		candidatesPool = append(candidatesPool, candidate{q: q, score: totalScore})
	}

	// 构建权重分布（使用更陡峰的指数衰减）
	var weights []float64
	for _, c := range candidatesPool {
		// 使用更大的衰减系数，使得分数差异更明显
		weight := math.Exp(-(c.score - minScore) * 4)
		weights = append(weights, weight)
	}

	// 执行加权随机选择
	sumWeights := 0.0
	for _, w := range weights {
		sumWeights += w
	}

	r := rand.Float64() * sumWeights
	for i, w := range weights {
		if r < w {
			return candidatesPool[i].q, true
		}
		r -= w
	}

	// 保底返回最后一个
	return candidatesPool[len(candidatesPool)-1].q, true
}

// removeQuestion 从切片中删除指定题目（保持顺序）
func removeQuestion(slice []entity.QuestionBank, id int) []entity.QuestionBank {
	for i, q := range slice {
		if q.ID == id {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

// 当前试卷已选题目缓存
var currentSelectedQuestions []entity.QuestionBank

// 获取历史题目的使用时间
func getHistoryQuestionTime(id int) (time.Time, bool) {
	// 使用全局题目缓存
	historyMutex.RLock()
	defer historyMutex.RUnlock()

	if _, exists := historyQuestionIDs[id]; exists {
		// 注意：这里我们使用当前时间作为上次使用时间
		// 在实际应用中，应该从数据库中获取真实的使用时间
		return time.Now().Add(-30 * 24 * time.Hour), true // 模拟30天前使用过
	}
	return time.Time{}, false
}

// 获取当前试卷已选题目
func getSelectedQuestions() []entity.QuestionBank {
	return currentSelectedQuestions
}

// 更新已选题目列表
func updateSelectedQuestions(questions []entity.QuestionBank) {
	currentSelectedQuestions = questions
}

// 初始化时加载历史数据
func init() {
	rand.Seed(time.Now().UnixNano())                          // 初始化随机种子
	currentSelectedQuestions = make([]entity.QuestionBank, 0) // 初始化已选题目列表
}

//
//// RandomSelectTopic 随机选题核心算法
//func RandomSelectTopic(dataSource []entity.QuestionBank, targetDiff float64, selectCount int) []entity.QuestionBank {
//	result := make([]entity.QuestionBank, 0, selectCount)
//
//	// 参数校验
//	if len(dataSource) == 0 ||
//		targetDiff < 1.0 ||
//		targetDiff > 5.0 ||
//		selectCount <= 0 {
//		return result
//	}
//	// 预排序
//	sort.SliceStable(dataSource, func(i, j int) bool {
//		return dataSource[i].Difficulty < dataSource[j].Difficulty
//	})
//
//	workingSet := make([]entity.QuestionBank, len(dataSource))
//	copy(workingSet, dataSource) // 避免修改原始数据
//	for remaining := selectCount; remaining > 0 && len(workingSet) > 0; remaining-- {
//		// 单次选择流程
//		if selected, ok := selectQuestion(workingSet, targetDiff); ok {
//			result = append(result, selected)
//			workingSet = removeQuestion(workingSet, selected.ID)
//		} else {
//			break // 无可用题目时提前终止
//		}
//	}
//
//	return result
//}

// selectQuestion 单次选题逻辑
func selectQuestion(candidates []entity.QuestionBank, targetDiff float64) (entity.QuestionBank, bool) {
	if len(candidates) == 0 {
		return entity.QuestionBank{}, false
	}

	// 第一阶段：寻找最小差值
	minDelta := math.MaxFloat64
	var minCandidates []entity.QuestionBank

	for _, q := range candidates {
		currentDelta := math.Abs(float64(q.Difficulty) - targetDiff)

		switch {
		case currentDelta < minDelta:
			minDelta = currentDelta
			minCandidates = []entity.QuestionBank{q}
		case currentDelta == minDelta:
			minCandidates = append(minCandidates, q)
		}
	}

	// 第二阶段：从候选中随机选择
	if len(minCandidates) == 0 {
		return entity.QuestionBank{}, false
	}
	return minCandidates[rand.Intn(len(minCandidates))], true
}

/* 使用示例
func main() {
    // 准备测试数据
    questions := []QuestionBank{
        {ID: 1, Difficulty: 3.2},
        {ID: 2, Difficulty: 4.5},
        {ID: 3, Difficulty: 2.8},
        {ID: 4, Difficulty: 4.5},
        {ID: 5, Difficulty: 3.9},
    }

    // 执行选题
    selected := RandomSelectTopic(questions, 4.0, 3)

    // 输出结果
    fmt.Println("Selected questions:")
    for _, q := range selected {
        fmt.Printf("ID: %d, Difficulty: %.1f\n", q.ID, q.Difficulty)
    }
}
*/
