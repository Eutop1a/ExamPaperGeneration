package services

import (
	"fmt"
	"graduation/entity"
	"graduation/mapper"
	"sync"
	"time"
)

var (
	historyQuestionIDs = make(map[int]struct{})
	historyMutex       sync.RWMutex
	lastRefreshTime    time.Time
	cacheExpiration    = 30 * time.Minute // 缓存刷新间隔
	twoYears           = -2 * 365 * 24 * time.Hour
)

// 初始化时自动刷新缓存
func init() {
	refreshHistoryCache()
	go autoRefreshCache()
}

func autoRefreshCache() {
	ticker := time.NewTicker(cacheExpiration)
	defer ticker.Stop()

	for range ticker.C {
		refreshHistoryCache()
	}
}

func refreshHistoryCache() {
	fmt.Println("开始刷新历史缓存...")
	defer fmt.Println("缓存刷新完成，最后更新时间:", lastRefreshTime)
	// 查询两年内的试卷UID
	var paperUIDs []string
	mapper.DB.Model(&entity.TestPaperGenHistory{}).
		Where("update_time >= ?", time.Now().Add(twoYears)).
		Pluck("test_paper_uid", &paperUIDs)

	// 查询对应的题目ID
	var questionIDs []int
	if len(paperUIDs) > 0 {
		mapper.DB.Model(&entity.QuestionGenHistory{}).
			Where("test_paper_uid IN ?", paperUIDs).
			Pluck("question_bank_id", &questionIDs)
	}

	// 更新缓存
	historyMutex.Lock()
	defer historyMutex.Unlock()

	historyQuestionIDs = make(map[int]struct{})
	for _, id := range questionIDs {
		historyQuestionIDs[id] = struct{}{}
	}
	lastRefreshTime = time.Now()
	// ...原有逻辑...
	fmt.Printf("加载到 %d 个历史题目ID\n", len(questionIDs))
}

// 判断题目是否在历史缓存中
func isHistoricalQuestion(id int) bool {
	historyMutex.RLock()
	defer historyMutex.RUnlock()

	_, exists := historyQuestionIDs[id]
	return exists
}
