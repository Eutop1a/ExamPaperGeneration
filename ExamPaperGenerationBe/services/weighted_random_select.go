package services

import (
	"graduation/entity"
	"graduation/mapper"
	"math"
	"math/rand"
	"time"
)

const (
	TARGET_TOTAL_SCORE = 100.0
	MAX_SIMILARITY     = 0.3 // 30% 相似度阈值
)

// KnowledgePointWeight 知识点权重
type KnowledgePointWeight struct {
	Label1 string  `json:"label1"`
	Weight float64 `json:"weight"`
}

// QuestionTypeRequirement 题型要求
type QuestionTypeRequirement struct {
	MinCount    int     // 最少题目数量
	TargetScore float64 // 目标分数
}

// WeightedRandomSelect 加权随机选题
func WeightedRandomSelect(
	questions []entity.QuestionBank,
	averageDifficulty float64,
	knowledgeWeights []KnowledgePointWeight,
	questionTypeRequirements map[string]QuestionTypeRequirement,
	selectedQuestionIds []int, // 手动选择的题目ID
	excludedQuestionIds []int, // 排除已选过的题目ID
) []entity.QuestionBank {
	// 初始化随机数生成器
	rand.Seed(time.Now().UnixNano())

	// 创建已排除题目ID的映射，用于快速查找
	excludedIdsMap := make(map[int]bool)
	for _, id := range excludedQuestionIds {
		excludedIdsMap[id] = true
	}

	// 创建已选择题目ID的映射，用于快速查找
	selectedIdsMap := make(map[int]bool)
	for _, id := range selectedQuestionIds {
		selectedIdsMap[id] = true
	}

	// 过滤掉已排除的题目
	var filteredQuestions []entity.QuestionBank
	for _, q := range questions {
		if !excludedIdsMap[q.ID] {
			filteredQuestions = append(filteredQuestions, q)
		}
	}

	// 按题型和知识点分组
	questionsByTypeAndKnowledge := make(map[string]map[string][]entity.QuestionBank)
	for _, q := range filteredQuestions {
		if _, ok := questionsByTypeAndKnowledge[q.TopicType]; !ok {
			questionsByTypeAndKnowledge[q.TopicType] = make(map[string][]entity.QuestionBank)
		}
		questionsByTypeAndKnowledge[q.TopicType][q.Label1] = append(questionsByTypeAndKnowledge[q.TopicType][q.Label1], q)
	}

	var selectedQuestions []entity.QuestionBank
	totalScore := 0.0

	// 第一步：处理手动选择的题目
	manualQuestions := make([]entity.QuestionBank, 0)
	for _, q := range questions {
		if selectedIdsMap[q.ID] { // 使用传入的已选择题目ID
			manualQuestions = append(manualQuestions, q)
			totalScore += q.Score
		}
	}

	// 将手动选择的题目添加到结果中
	selectedQuestions = append(selectedQuestions, manualQuestions...)

	// 第二步：优先选择简答题，确保至少5道
	shortAnswerQuestions := questionsByTypeAndKnowledge["简答题"]
	if len(shortAnswerQuestions) > 0 {
		// 计算每个知识点的权重
		knowledgeWeightsMap := make(map[string]float64)
		totalWeight := 0.0
		for _, kw := range knowledgeWeights {
			knowledgeWeightsMap[kw.Label1] = kw.Weight
			totalWeight += kw.Weight
		}

		// 选择至少5道简答题
		selectedCount := 0
		for selectedCount < 5 {
			// 根据权重随机选择知识点
			r := rand.Float64() * totalWeight
			var selectedKnowledge string
			var currentWeight float64
			for knowledge, weight := range knowledgeWeightsMap {
				currentWeight += weight
				if r <= currentWeight {
					selectedKnowledge = knowledge
					break
				}
			}

			// 从选中的知识点中随机选择一道题
			if questions, ok := shortAnswerQuestions[selectedKnowledge]; ok && len(questions) > 0 {
				// 选择难度最接近平均难度的题目
				bestIndex := 0
				bestDiff := math.MaxFloat64
				for i, q := range questions {
					diff := math.Abs(float64(q.Difficulty) - averageDifficulty)
					if diff < bestDiff {
						bestDiff = diff
						bestIndex = i
					}
				}

				// 添加选中的题目
				selectedQuestions = append(selectedQuestions, questions[bestIndex])
				selectedCount++
				totalScore += questions[bestIndex].Score

				// 从可选题目中移除已选题目
				shortAnswerQuestions[selectedKnowledge] = append(questions[:bestIndex], questions[bestIndex+1:]...)
			}
		}
	}

	// 第三步：满足其他题型的最少数量要求
	for questionType, requirement := range questionTypeRequirements {
		// 跳过简答题，因为已经在第二步处理过了
		if questionType == "简答题" {
			continue
		}

		typeQuestions := questionsByTypeAndKnowledge[questionType]
		if len(typeQuestions) == 0 {
			continue
		}

		// 计算每个知识点的权重
		knowledgeWeightsMap := make(map[string]float64)
		totalWeight := 0.0
		for _, kw := range knowledgeWeights {
			knowledgeWeightsMap[kw.Label1] = kw.Weight
			totalWeight += kw.Weight
		}

		// 选择题目直到满足最少数量要求
		selectedCount := 0
		for selectedCount < requirement.MinCount {
			// 根据权重随机选择知识点
			r := rand.Float64() * totalWeight
			var selectedKnowledge string
			var currentWeight float64
			for knowledge, weight := range knowledgeWeightsMap {
				currentWeight += weight
				if r <= currentWeight {
					selectedKnowledge = knowledge
					break
				}
			}

			// 从选中的知识点中随机选择一道题
			if questions, ok := typeQuestions[selectedKnowledge]; ok && len(questions) > 0 {
				// 选择难度最接近平均难度的题目
				bestIndex := 0
				bestDiff := math.MaxFloat64
				for i, q := range questions {
					diff := math.Abs(float64(q.Difficulty) - averageDifficulty)
					if diff < bestDiff {
						bestDiff = diff
						bestIndex = i
					}
				}

				// 添加选中的题目
				selectedQuestions = append(selectedQuestions, questions[bestIndex])
				selectedCount++
				totalScore += questions[bestIndex].Score

				// 从可选题目中移除已选题目
				typeQuestions[selectedKnowledge] = append(questions[:bestIndex], questions[bestIndex+1:]...)
			}
		}
	}

	// 第四步：循环选择题目直到总分达到100分
	for totalScore < TARGET_TOTAL_SCORE {
		// 计算还需要多少分
		remainingScore := TARGET_TOTAL_SCORE - totalScore

		// 计算每个知识点的目标分数
		knowledgeTargetScores := make(map[string]float64)
		totalWeight := 0.0
		for _, kw := range knowledgeWeights {
			totalWeight += kw.Weight
		}
		for _, kw := range knowledgeWeights {
			knowledgeTargetScores[kw.Label1] = (kw.Weight / totalWeight) * remainingScore
		}

		// 调整分数为整数，优先分配给简答题
		adjustedScores := adjustScoresToIntegers(knowledgeTargetScores, remainingScore)

		// 为每个知识点选择题目直到达到目标分数
		for knowledge, targetScore := range adjustedScores {
			currentScore := 0.0

			// 优先选择简答题
			questionTypes := []string{"简答题", "选择题", "判断题", "填空题"}
			for _, questionType := range questionTypes {
				questions := questionsByTypeAndKnowledge[questionType]
				if len(questions[knowledge]) == 0 {
					continue
				}

				for currentScore < targetScore {
					// 选择难度最接近平均难度的题目
					bestIndex := 0
					bestDiff := math.MaxFloat64
					for i, q := range questions[knowledge] {
						diff := math.Abs(float64(q.Difficulty) - averageDifficulty)
						if diff < bestDiff {
							bestDiff = diff
							bestIndex = i
						}
					}

					if bestDiff == math.MaxFloat64 {
						break
					}

					// 添加选中的题目
					selectedQuestions = append(selectedQuestions, questions[knowledge][bestIndex])
					currentScore += questions[knowledge][bestIndex].Score
					totalScore += questions[knowledge][bestIndex].Score

					// 从可选题目中移除已选题目
					questions[knowledge] = append(questions[knowledge][:bestIndex], questions[knowledge][bestIndex+1:]...)

					// 如果总分超过100分，移除最后一道题
					if totalScore > TARGET_TOTAL_SCORE {
						lastQuestion := selectedQuestions[len(selectedQuestions)-1]
						selectedQuestions = selectedQuestions[:len(selectedQuestions)-1]
						totalScore -= lastQuestion.Score
						break
					}
				}
			}
		}

		// 如果总分仍然不足100分，继续循环
		if totalScore < TARGET_TOTAL_SCORE {
			continue
		}
	}

	return selectedQuestions
}

// adjustScoresToIntegers 调整分数为整数，确保总和不变，并尽量减少奇数
func adjustScoresToIntegers(scores map[string]float64, totalScore float64) map[string]float64 {
	// 创建结果map
	result := make(map[string]float64)

	// 计算每个知识点的整数部分和小数部分
	type scoreInfo struct {
		knowledge string
		integer   int
		fraction  float64
	}

	scoreInfos := make([]scoreInfo, 0, len(scores))
	var totalInteger int

	for knowledge, score := range scores {
		integer := int(math.Floor(score))
		fraction := score - float64(integer)
		scoreInfos = append(scoreInfos, scoreInfo{knowledge, integer, fraction})
		totalInteger += integer
		result[knowledge] = float64(integer)
	}

	// 计算需要分配的小数点
	remainingPoints := int(totalScore) - totalInteger

	// 按小数部分从大到小排序
	for i := 0; i < len(scoreInfos)-1; i++ {
		for j := i + 1; j < len(scoreInfos); j++ {
			if scoreInfos[i].fraction < scoreInfos[j].fraction {
				scoreInfos[i], scoreInfos[j] = scoreInfos[j], scoreInfos[i]
			}
		}
	}

	// 分配剩余的小数点，优先分配给偶数
	for i := 0; i < remainingPoints && i < len(scoreInfos); i++ {
		// 如果当前分数是奇数，且还有下一个知识点，且下一个知识点的小数部分也很大
		if result[scoreInfos[i].knowledge] > 0 && int(result[scoreInfos[i].knowledge])%2 == 1 &&
			i+1 < len(scoreInfos) && scoreInfos[i+1].fraction > 0.3 {
			// 跳过当前知识点，分配给下一个
			result[scoreInfos[i+1].knowledge] += 1
			// 将当前知识点移到后面
			scoreInfos[i], scoreInfos[i+1] = scoreInfos[i+1], scoreInfos[i]
		} else {
			result[scoreInfos[i].knowledge] += 1
		}
	}

	return result
}

// GetExcludedQuestionIds 获取需要排除的题目ID列表
func GetExcludedQuestionIds(similarityThreshold float64) []int {
	// 获取最近两年内的所有题目ID
	twoYearsAgo := time.Now().AddDate(-2, 0, 0)
	var questionIds []int

	// 一次性查询所有符合条件的题目ID
	if err := mapper.DB.Model(&entity.QuestionGenHistory{}).
		Select("DISTINCT questiongenhistory.question_bank_id").
		Joins("JOIN testpapergenhistory ON testpapergenhistory.test_paper_uid = questiongenhistory.test_paper_uid").
		Where("testpapergenhistory.update_time >= ?", twoYearsAgo).
		Pluck("question_bank_id", &questionIds).Error; err != nil {
		return nil
	}

	// 如果没有题目，直接返回空切片
	if len(questionIds) == 0 {
		return nil
	}

	// 根据相似度阈值计算需要排除的题目数量
	excludeRatio := 1 - similarityThreshold
	excludeCount := int(float64(len(questionIds)) * excludeRatio)
	if excludeCount <= 0 {
		return nil
	}

	// 随机打乱题目ID
	rand.Shuffle(len(questionIds), func(i, j int) {
		questionIds[i], questionIds[j] = questionIds[j], questionIds[i]
	})

	// 返回需要排除的题目ID
	return questionIds[:excludeCount]
}

// CalculateTestPaperSimilarity 计算两套试卷的相似度
func CalculateTestPaperSimilarity(paper1Ids []int, paper2Ids []int) float64 {
	if len(paper1Ids) == 0 || len(paper2Ids) == 0 {
		return 0
	}

	// 创建ID映射，用于快速查找
	paper1Map := make(map[int]bool)
	for _, id := range paper1Ids {
		paper1Map[id] = true
	}

	// 计算重复的题目数量
	duplicateCount := 0
	for _, id := range paper2Ids {
		if paper1Map[id] {
			duplicateCount++
		}
	}

	// 计算相似度：重复题目数量 / 较小试卷的题目数量
	minCount := len(paper1Ids)
	if len(paper2Ids) < minCount {
		minCount = len(paper2Ids)
	}

	return float64(duplicateCount) / float64(minCount)
}
