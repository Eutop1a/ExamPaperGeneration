package mapper

import (
	"gorm.io/gorm"
	"graduation/entity"
)

// QuestionLabelsMapper 接口定义
type QuestionLabelsMapper struct {
	db *gorm.DB
}

// NewQuestionLabelsMapper 创建一个新的 QuestionLabelsMapper 实例
func NewQuestionLabelsMapper() *QuestionLabelsMapper {
	return &QuestionLabelsMapper{
		db: DB,
	}
}

// GetAllQuestionLabels 获取所有题目标签
func (m *QuestionLabelsMapper) GetAllQuestionLabels() ([]entity.QuestionLabels, error) {
	var labels []entity.QuestionLabels
	result := m.db.Find(&labels)
	return labels, result.Error
}

// GetDistinctChapter1 获取不同的 chapter_1
func (m *QuestionLabelsMapper) GetDistinctChapter1() ([]entity.QuestionLabels, error) {
	var labels []entity.QuestionLabels
	result := m.db.Distinct("chapter_1").Find(&labels)
	return labels, result.Error
}

// GetDistinctChapter2 获取不同的 chapter_2
func (m *QuestionLabelsMapper) GetDistinctChapter2() ([]entity.QuestionLabels, error) {
	var labels []entity.QuestionLabels
	result := m.db.Distinct("chapter_2").Find(&labels)
	return labels, result.Error
}

// GetChapter2ByChapter1 根据 chapter_1 获取对应的 chapter_2
func (m *QuestionLabelsMapper) GetChapter2ByChapter1(chapter1 string) ([]entity.QuestionLabels, error) {
	var labels []entity.QuestionLabels
	result := m.db.Distinct("chapter_2").Where("chapter_1 =?", chapter1).Find(&labels)
	return labels, result.Error
}

// GetDistinctLabel1 获取不同的 label_1
func (m *QuestionLabelsMapper) GetDistinctLabel1() ([]entity.QuestionLabels, error) {
	var labels []entity.QuestionLabels
	result := m.db.Distinct("label_1").Find(&labels)
	return labels, result.Error
}

// GetDistinctLabel2 获取不同的 label_2
func (m *QuestionLabelsMapper) GetDistinctLabel2() ([]entity.QuestionLabels, error) {
	var labels []entity.QuestionLabels
	result := m.db.Distinct("label_2").Find(&labels)
	return labels, result.Error
}
