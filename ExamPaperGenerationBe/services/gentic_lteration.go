package services

import (
	"graduation/entity"
	"math"
	"math/rand"
	"sort"
)

// 参数调整：

// // 在NewGeneticIteration中可调整：
// gi.similarityWeight = 0.5 // 加大相似度权重

// GeneticIteration 遗传算法迭代
type GeneticIteration struct {
	questions                []entity.QuestionBank
	averageDifficulty        float64
	knowledgeWeights         []KnowledgePointWeight
	questionTypeRequirements map[string]QuestionTypeRequirement
	selectedQuestionIds      []int // 手动选择的题目ID
	excludedQuestionIds      []int // 排除已选过的题目ID
	populationSize           int
	generationCount          int
	mutationRate             float64
	eliteCount               int
	similarityThreshold      float64
	Variance                 []float64 // 方差列表
}

// NewGeneticIteration 创建遗传算法迭代器
func NewGeneticIteration(
	questions []entity.QuestionBank,
	averageDifficulty float64,
	knowledgeWeights []KnowledgePointWeight,
	questionTypeRequirements map[string]QuestionTypeRequirement,
	selectedQuestionIds []int,
	excludedQuestionIds []int,
	iterationsNum int,
) *GeneticIteration {
	// 根据迭代次数动态调整种群大小
	populationSize := 40 // 默认种群大小
	if iterationsNum > 5000 {
		populationSize = 30 // 迭代次数多时，减小种群大小
	} else if iterationsNum < 1000 {
		populationSize = 50 // 迭代次数少时，增大种群大小
	}

	return &GeneticIteration{
		questions:                questions,
		averageDifficulty:        averageDifficulty,
		knowledgeWeights:         knowledgeWeights,
		questionTypeRequirements: questionTypeRequirements,
		selectedQuestionIds:      selectedQuestionIds,
		excludedQuestionIds:      excludedQuestionIds,
		populationSize:           populationSize,
		generationCount:          iterationsNum,
		mutationRate:             0.05,               // 降低变异率，提高收敛性
		eliteCount:               populationSize / 3, // 增加精英保留比例
		similarityThreshold:      0.3,
		Variance:                 []float64{},
	}
}

// Run 运行遗传算法
func (gi *GeneticIteration) Run() []entity.QuestionBank {
	// 初始化种群
	population := gi.initializePopulation()
	bestFitness := -1.0
	bestSolution := population[0]
	noImprovementCount := 0
	const maxNoImprovement = 100 // 增加无改进的容忍度

	// 迭代进化
	for generation := 0; generation < gi.generationCount; generation++ {
		// 评估适应度
		fitnesses := gi.evaluateFitness(population)

		// 计算方差
		variance := gi.calculateVariance(fitnesses)
		gi.Variance = append(gi.Variance, variance)

		// 更新最优解
		currentBestFitness := -1.0
		currentBestIndex := 0
		for i, fitness := range fitnesses {
			if fitness > currentBestFitness {
				currentBestFitness = fitness
				currentBestIndex = i
			}
		}

		if currentBestFitness > bestFitness {
			bestFitness = currentBestFitness
			bestSolution = population[currentBestIndex]
			noImprovementCount = 0
		} else {
			noImprovementCount++
		}

		// 如果连续多代没有改进，提前结束
		if noImprovementCount >= maxNoImprovement && generation > 200 { // 至少迭代200代
			break
		}

		// 选择精英
		elites := gi.selectElites(population, fitnesses)

		// 生成新一代
		newPopulation := make([][]entity.QuestionBank, gi.populationSize)
		copy(newPopulation[:gi.eliteCount], elites)

		// 使用锦标赛选择进行交叉和变异
		for i := gi.eliteCount; i < gi.populationSize; i++ {
			parent1 := gi.selectParent(population, fitnesses)
			parent2 := gi.selectParent(population, fitnesses)
			child := gi.crossover(parent1, parent2)
			child = gi.mutate(child)
			newPopulation[i] = child
		}

		population = newPopulation
	}

	return bestSolution
}

// initializePopulation 初始化种群
func (gi *GeneticIteration) initializePopulation() [][]entity.QuestionBank {
	population := make([][]entity.QuestionBank, gi.populationSize)
	for i := 0; i < gi.populationSize; i++ {
		solution := gi.generateRandomSolution()
		population[i] = solution
	}
	return population
}

// generateRandomSolution 生成随机解
func (gi *GeneticIteration) generateRandomSolution() []entity.QuestionBank {
	var solution []entity.QuestionBank
	totalScore := 0.0
	questionTypes := []string{"简答题", "选择题", "判断题", "填空题"}

	// 添加手动选择的题目
	for _, q := range gi.questions {
		if contains(gi.selectedQuestionIds, q.ID) {
			solution = append(solution, q)
			totalScore += q.Score
		}
	}

	// 按题型和知识点分组
	questionsByTypeAndKnowledge := make(map[string]map[string][]entity.QuestionBank)
	for _, q := range gi.questions {
		if !contains(gi.excludedQuestionIds, q.ID) && !contains(gi.selectedQuestionIds, q.ID) {
			if _, ok := questionsByTypeAndKnowledge[q.TopicType]; !ok {
				questionsByTypeAndKnowledge[q.TopicType] = make(map[string][]entity.QuestionBank)
			}
			questionsByTypeAndKnowledge[q.TopicType][q.Label1] = append(questionsByTypeAndKnowledge[q.TopicType][q.Label1], q)
		}
	}

	// 优先选择简答题，确保至少5道
	shortAnswerQuestions := questionsByTypeAndKnowledge["简答题"]
	if len(shortAnswerQuestions) > 0 {
		selectedCount := 0
		for selectedCount < 5 && totalScore < TARGET_TOTAL_SCORE {
			for knowledge, questions := range shortAnswerQuestions {
				if len(questions) > 0 {
					index := rand.Intn(len(questions))
					solution = append(solution, questions[index])
					totalScore += questions[index].Score
					shortAnswerQuestions[knowledge] = append(questions[:index], questions[index+1:]...)
					selectedCount++
					if totalScore >= TARGET_TOTAL_SCORE {
						break
					}
				}
			}
		}
	}

	// 满足其他题型的最少数量要求
	for questionType, requirement := range gi.questionTypeRequirements {
		if questionType == "简答题" {
			continue
		}

		typeQuestions := questionsByTypeAndKnowledge[questionType]
		if len(typeQuestions) == 0 {
			continue
		}

		selectedCount := 0
		for selectedCount < requirement.MinCount && totalScore < TARGET_TOTAL_SCORE {
			for knowledge, questions := range typeQuestions {
				if len(questions) > 0 {
					index := rand.Intn(len(questions))
					solution = append(solution, questions[index])
					totalScore += questions[index].Score
					typeQuestions[knowledge] = append(questions[:index], questions[index+1:]...)
					selectedCount++
					if totalScore >= TARGET_TOTAL_SCORE {
						break
					}
				}
			}
		}
	}

	// 循环选择题目直到总分达到100分
	maxAttempts := 100 // 添加最大尝试次数
	attempts := 0
	for totalScore < TARGET_TOTAL_SCORE && attempts < maxAttempts {
		attempts++
		// 计算每个知识点的目标分数
		knowledgeTargetScores := make(map[string]float64)
		totalWeight := 0.0
		for _, kw := range gi.knowledgeWeights {
			totalWeight += kw.Weight
		}
		for _, kw := range gi.knowledgeWeights {
			knowledgeTargetScores[kw.Label1] = (kw.Weight / totalWeight) * (TARGET_TOTAL_SCORE - totalScore)
		}

		// 调整分数为整数，优先分配给简答题
		adjustedScores := adjustScoresToIntegers(knowledgeTargetScores, TARGET_TOTAL_SCORE-totalScore)

		// 为每个知识点选择题目直到达到目标分数
		for knowledge, targetScore := range adjustedScores {
			currentScore := 0.0

			// 优先选择简答题
			for _, questionType := range questionTypes {
				questions := questionsByTypeAndKnowledge[questionType]
				if len(questions[knowledge]) == 0 {
					continue
				}

				for currentScore < targetScore && totalScore < TARGET_TOTAL_SCORE {
					index := rand.Intn(len(questions[knowledge]))
					solution = append(solution, questions[knowledge][index])
					currentScore += questions[knowledge][index].Score
					totalScore += questions[knowledge][index].Score
					questions[knowledge] = append(questions[knowledge][:index], questions[knowledge][index+1:]...)

					if totalScore > TARGET_TOTAL_SCORE {
						lastQuestion := solution[len(solution)-1]
						solution = solution[:len(solution)-1]
						totalScore -= lastQuestion.Score
						break
					}
				}
			}
		}

		// 如果总分仍然小于100分，尝试添加一个合适分数的题目
		if totalScore < TARGET_TOTAL_SCORE {
			remainingScore := TARGET_TOTAL_SCORE - totalScore
			// 遍历所有题型和知识点，寻找合适分数的题目
			for _, questionType := range questionTypes {
				for knowledge, questions := range questionsByTypeAndKnowledge[questionType] {
					for i, q := range questions {
						if math.Abs(q.Score-remainingScore) < 0.1 { // 允许0.1分的误差
							solution = append(solution, q)
							totalScore += q.Score
							questionsByTypeAndKnowledge[questionType][knowledge] = append(questions[:i], questions[i+1:]...)
							break
						}
					}
					if totalScore >= TARGET_TOTAL_SCORE {
						break
					}
				}
				if totalScore >= TARGET_TOTAL_SCORE {
					break
				}
			}
		}

		if totalScore >= TARGET_TOTAL_SCORE {
			break
		}
	}

	// 如果尝试次数达到上限仍未达到100分，重新生成解
	if totalScore != TARGET_TOTAL_SCORE {
		return gi.generateRandomSolution()
	}

	return solution
}

// evaluateFitness 评估适应度
func (gi *GeneticIteration) evaluateFitness(population [][]entity.QuestionBank) []float64 {
	fitnesses := make([]float64, len(population))
	for i, solution := range population {
		fitnesses[i] = gi.calculateFitness(solution)
	}
	return fitnesses
}

// calculateFitness 计算适应度
func (gi *GeneticIteration) calculateFitness(solution []entity.QuestionBank) float64 {
	totalScore := 0.0
	for _, q := range solution {
		totalScore += q.Score
	}

	// 总分必须为100分
	if totalScore != TARGET_TOTAL_SCORE {
		return 0.0
	}

	// 计算难度偏差
	difficultyDeviation := 0.0
	for _, q := range solution {
		difficultyDeviation += math.Abs(float64(q.Difficulty) - gi.averageDifficulty)
	}
	difficultyDeviation /= float64(len(solution))
	difficultyDeviation = math.Min(difficultyDeviation, 1.0) // 归一化

	// 计算知识点覆盖偏差
	knowledgeCoverage := make(map[string]float64)
	for _, q := range solution {
		knowledgeCoverage[q.Label1] += q.Score
	}

	knowledgeCoverageDeviation := 0.0
	totalWeight := 0.0
	for _, kw := range gi.knowledgeWeights {
		totalWeight += kw.Weight
	}
	for _, kw := range gi.knowledgeWeights {
		expectedScore := (kw.Weight / totalWeight) * TARGET_TOTAL_SCORE
		actualScore := knowledgeCoverage[kw.Label1]
		knowledgeCoverageDeviation += math.Abs(actualScore - expectedScore)
	}
	knowledgeCoverageDeviation = math.Min(knowledgeCoverageDeviation/TARGET_TOTAL_SCORE, 1.0) // 归一化

	// 计算题型要求偏差
	typeRequirementDeviation := 0.0
	typeCounts := make(map[string]int)
	for _, q := range solution {
		typeCounts[q.TopicType]++
	}
	for questionType, requirement := range gi.questionTypeRequirements {
		actualCount := typeCounts[questionType]
		if actualCount < requirement.MinCount {
			typeRequirementDeviation += float64(requirement.MinCount - actualCount)
		}
	}
	typeRequirementDeviation = math.Min(typeRequirementDeviation/float64(len(gi.questionTypeRequirements)), 1.0) // 归一化

	// 计算相似度惩罚
	similarityPenalty := gi.calcSimilarityPenalty(solution)
	similarityPenalty = math.Min(similarityPenalty/float64(len(solution)), 1.0) // 归一化

	// 综合评分，调整权重
	fitness := 100.0 * (1.0 - (difficultyDeviation * 0.2) - (knowledgeCoverageDeviation * 0.3) - (typeRequirementDeviation * 0.3) - (similarityPenalty * 0.2))
	return math.Max(0.0, fitness)
}

// selectElites 选择精英
func (gi *GeneticIteration) selectElites(population [][]entity.QuestionBank, fitnesses []float64) [][]entity.QuestionBank {
	elites := make([][]entity.QuestionBank, gi.eliteCount)
	indices := make([]int, len(fitnesses))
	for i := range indices {
		indices[i] = i
	}
	sort.Slice(indices, func(i, j int) bool {
		return fitnesses[indices[i]] > fitnesses[indices[j]]
	})
	for i := 0; i < gi.eliteCount; i++ {
		elites[i] = population[indices[i]]
	}
	return elites
}

// selectParent 选择父代
func (gi *GeneticIteration) selectParent(population [][]entity.QuestionBank, fitnesses []float64) []entity.QuestionBank {
	// 使用锦标赛选择，增加锦标赛规模
	tournamentSize := 5 // 增加锦标赛规模
	bestIndex := rand.Intn(len(population))
	bestFitness := fitnesses[bestIndex]

	for i := 1; i < tournamentSize; i++ {
		index := rand.Intn(len(population))
		if fitnesses[index] > bestFitness {
			bestIndex = index
			bestFitness = fitnesses[index]
		}
	}

	return population[bestIndex]
}

// crossover 交叉
func (gi *GeneticIteration) crossover(parent1, parent2 []entity.QuestionBank) []entity.QuestionBank {
	if len(parent1) == 0 || len(parent2) == 0 {
		return gi.generateRandomSolution()
	}

	// 确保两个父代长度相同
	minLen := min(len(parent1), len(parent2))
	if minLen < 2 {
		return gi.generateRandomSolution()
	}

	// 使用多点交叉，确保交叉点在有效范围内
	crossoverPoints := make([]int, 2)
	crossoverPoints[0] = rand.Intn(minLen-1) + 1                                   // 确保至少有一个元素在交叉点之前
	crossoverPoints[1] = rand.Intn(minLen-crossoverPoints[0]) + crossoverPoints[0] // 确保第二个交叉点在第一个之后

	child := make([]entity.QuestionBank, 0, minLen)
	child = append(child, parent1[:crossoverPoints[0]]...)
	child = append(child, parent2[crossoverPoints[0]:crossoverPoints[1]]...)
	child = append(child, parent1[crossoverPoints[1]:]...)

	// 确保总分为100分
	totalScore := 0.0
	for _, q := range child {
		totalScore += q.Score
	}

	if totalScore != TARGET_TOTAL_SCORE {
		return gi.generateRandomSolution()
	}

	return child
}

// mutate 变异
func (gi *GeneticIteration) mutate(solution []entity.QuestionBank) []entity.QuestionBank {
	if len(solution) == 0 {
		return gi.generateRandomSolution()
	}

	if rand.Float64() < gi.mutationRate {
		// 只变异一个基因
		mutationPoint := rand.Intn(len(solution))
		newSolution := make([]entity.QuestionBank, len(solution))
		copy(newSolution, solution)

		// 找到相同知识点和题型的其他题目
		originalQuestion := solution[mutationPoint]
		var candidates []entity.QuestionBank
		for _, q := range gi.questions {
			if q.Label1 == originalQuestion.Label1 &&
				q.TopicType == originalQuestion.TopicType &&
				!contains(gi.excludedQuestionIds, q.ID) &&
				!contains(gi.selectedQuestionIds, q.ID) {
				candidates = append(candidates, q)
			}
		}

		if len(candidates) > 0 {
			newQuestion := candidates[rand.Intn(len(candidates))]
			newSolution[mutationPoint] = newQuestion

			// 检查总分
			totalScore := 0.0
			for _, q := range newSolution {
				totalScore += q.Score
			}

			if totalScore == TARGET_TOTAL_SCORE {
				return newSolution
			}
		}
	}

	return solution
}

// findBestSolution 找到最优解
func (gi *GeneticIteration) findBestSolution(population [][]entity.QuestionBank) []entity.QuestionBank {
	bestFitness := -1.0
	bestSolution := population[0]

	for _, solution := range population {
		fitness := gi.calculateFitness(solution)
		if fitness > bestFitness {
			bestFitness = fitness
			bestSolution = solution
		}
	}

	return bestSolution
}

// calcSimilarityPenalty 计算相似度惩罚
func (gi *GeneticIteration) calcSimilarityPenalty(solution []entity.QuestionBank) float64 {
	penalty := 0.0
	for i, q1 := range solution {
		for j := i + 1; j < len(solution); j++ {
			q2 := solution[j]
			if q1.ID == q2.ID {
				penalty += 1.0
			}
		}
	}
	return penalty
}

// min 返回两个整数中的较小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// contains 检查切片中是否包含某个值
func contains(slice []int, value int) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

// calculateVariance 计算适应度的方差
func (gi *GeneticIteration) calculateVariance(fitnesses []float64) float64 {
	if len(fitnesses) == 0 {
		return 0.0
	}

	mean := 0.0
	for _, fitness := range fitnesses {
		mean += fitness
	}
	mean /= float64(len(fitnesses))

	variance := 0.0
	for _, fitness := range fitnesses {
		diff := fitness - mean
		variance += diff * diff
	}
	variance /= float64(len(fitnesses))

	return variance
}
