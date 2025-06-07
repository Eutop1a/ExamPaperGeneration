package mapper

import (
	"gorm.io/gorm"
	"graduation/entity"
)

// QuestionGenHistoryMapper 接口定义
type QuestionGenHistoryMapper struct {
	db *gorm.DB
}

// NewQuestionGenHistoryMapper 创建一个新的 QuestionGenHistoryMapper 实例
func NewQuestionGenHistoryMapper() *QuestionGenHistoryMapper {
	return &QuestionGenHistoryMapper{
		db: DB,
	}
}

// InsertQuestionGenHistories 插入多条问题生成历史记录
func (m *QuestionGenHistoryMapper) InsertQuestionGenHistories(list []entity.QuestionGenHistory) (int64, error) {
	result := m.db.Create(&list)
	return result.RowsAffected, result.Error
}

// GetQuestionGenHistoriesByTestPaperUid 根据试卷 UID 获取问题生成历史记录
func (m *QuestionGenHistoryMapper) GetQuestionGenHistoriesByTestPaperUid(testPaperUid string) ([]entity.QuestionGenHistory, error) {
	var histories []entity.QuestionGenHistory
	result := m.db.Where("test_paper_uid =?", testPaperUid).Find(&histories)
	return histories, result.Error
}

// DeleteQuestionGenHistoryByTestPaperUid 根据试卷 UID 删除问题生成历史记录
func (m *QuestionGenHistoryMapper) DeleteQuestionGenHistoryByTestPaperUid(testPaperUid string) (int64, error) {
	result := m.db.Where("test_paper_uid =?", testPaperUid).Delete(&entity.QuestionGenHistory{})
	return result.RowsAffected, result.Error
}
