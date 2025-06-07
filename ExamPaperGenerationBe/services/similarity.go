package services

const (
	similarityThreshold = 0.7 // 相似度阈值
	historyWeight       = 0.3 // 历史影响权重
)

// 计算两套试卷的Jaccard相似度
func jaccardSimilarity(a, b []int) float64 {
	set := make(map[int]struct{})
	intersection := 0

	for _, id := range a {
		set[id] = struct{}{}
	}

	for _, id := range b {
		if _, exists := set[id]; exists {
			intersection++
		}
	}

	union := len(a) + len(b) - intersection
	if union == 0 {
		return 0
	}
	return float64(intersection) / float64(union)
}

// 获取综合相似度评分
func GetSimilarityScore(current []int, history [][]int) float64 {
	total := 0.0
	for _, past := range history {
		total += jaccardSimilarity(current, past)
	}
	return total / float64(len(history)+1)
}
