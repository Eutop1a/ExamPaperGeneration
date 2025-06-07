package mapper

import (
	"gorm.io/gorm"
	"graduation/entity"
	"os"
)

// QuestionBankMapper 接口定义
type QuestionBankMapper struct {
	db *gorm.DB
}

// NewQuestionBankMapper 创建一个新的 QuestionBankMapper 实例
func NewQuestionBankMapper() *QuestionBankMapper {
	return &QuestionBankMapper{
		db: DB,
	}
}

// GetAllQuestionBank 获取所有题库记录，按更新时间降序排列
func (m *QuestionBankMapper) GetAllQuestionBank() ([]entity.QuestionBank, error) {
	var questionBanks []entity.QuestionBank
	result := m.db.Order("update_time desc").Find(&questionBanks)
	return questionBanks, result.Error
}

// GetDistinctTopicType 获取所有不同的题目类型
func (m *QuestionBankMapper) GetDistinctTopicType() ([]string, error) {
	var topicTypes []string
	result := m.db.Model(&entity.QuestionBank{}).Distinct("topic_type").Pluck("topic_type", &topicTypes)
	return topicTypes, result.Error
}

// SearchQuestionByTopic 根据题目类型和关键字搜索题目
func (m *QuestionBankMapper) SearchQuestionByTopic(topicType, keyword string) ([]entity.QuestionBank, error) {
	var questionBanks []entity.QuestionBank
	query := m.db
	if topicType != "" {
		query = query.Where("topic_type =?", topicType)
	}
	if keyword != "" {
		query = query.Where("topic LIKE ?", "%"+keyword+"%")
	}
	result := query.Find(&questionBanks)
	return questionBanks, result.Error
}

// InsertSingleQuestionBank 插入单条题库记录
func (m *QuestionBankMapper) InsertSingleQuestionBank(questionBank *entity.QuestionBank) (int64, error) {
	result := m.db.Create(questionBank)
	return result.RowsAffected, result.Error
}

// DeleteSingleQuestionBank 根据 ID 删除单条题库记录
func (m *QuestionBankMapper) DeleteSingleQuestionBank(id int) (int64, error) {
	var tmp entity.QuestionBank
	if err := m.db.First(&tmp, id).Error; err != nil {
		return 0, err
	}
	// 删除数据库记录
	result := m.db.Delete(&entity.QuestionBank{}, id)
	if result.Error != nil {
		return result.RowsAffected, result.Error
	}
	// 如果记录存在且图片路径不为空，则删除图片文件
	if result.RowsAffected > 0 && tmp.TopicImagePath != "" {
		err := os.Remove(tmp.TopicImagePath)
		if err != nil {
			return result.RowsAffected, err
		}
	}
	return result.RowsAffected, nil
}

// GetQuestionBankById 根据 ID 获取题库记录
func (m *QuestionBankMapper) GetQuestionBankById(id int) ([]entity.QuestionBank, error) {
	var questionBanks []entity.QuestionBank
	result := m.db.Where("id =?", id).Find(&questionBanks)
	return questionBanks, result.Error
}

// UpdateSingleQuestionBank 更新单条题库记录
func (m *QuestionBankMapper) UpdateSingleQuestionBank(questionBank *entity.QuestionBank) (int64, error) {
	//result := m.db.Model(questionBank).Where("id =?", questionBank.ID).Updates(questionBank)
	result := m.db.Model(questionBank).Updates(questionBank)
	return result.RowsAffected, result.Error
}

// GetAvgDifficultyByIds 根据题目的 ID 列表查询题目的平均难度
func (m *QuestionBankMapper) GetAvgDifficultyByIds(ids []int) (float64, error) {
	var totalDifficulty float64
	var count int64
	result := m.db.Model(&entity.QuestionBank{}).Where("id IN ?", ids).Select("SUM(difficulty)").Scan(&totalDifficulty)
	if result.Error != nil {
		return 0, result.Error
	}
	result = m.db.Model(&entity.QuestionBank{}).Where("id IN ?", ids).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	if count == 0 {
		return 0, nil
	}
	return totalDifficulty / float64(count), nil
}

// GetQuestionBankByIds 根据 ID 列表和生成范围查询不在 ID 列表里的题目
func (m *QuestionBankMapper) GetQuestionBankByIds(ids []int, generateRange []string) ([]entity.QuestionBank, error) {
	var questionBanks []entity.QuestionBank
	query := m.db
	if len(ids) > 0 {
		query = query.Where("id NOT IN ?", ids)
	}
	if len(generateRange) > 0 {
		query = query.Where("topic_type IN ?", generateRange)
	}
	result := query.Find(&questionBanks)
	return questionBanks, result.Error
}

// GetAll 获取所有题库记录
func (m *QuestionBankMapper) GetAll() ([]entity.QuestionBank, error) {
	var questionBanks []entity.QuestionBank
	result := m.db.Find(&questionBanks)
	return questionBanks, result.Error
}

// GetDistinctLabel1FromQuestionBank 获取题库表的不同 label1
func (m *QuestionBankMapper) GetDistinctLabel1FromQuestionBank() ([]string, error) {
	var label1s []string
	result := m.db.Model(&entity.QuestionBank{}).Distinct("label_1").Pluck("label_1", &label1s)
	return label1s, result.Error
}

// GetQuestionBankCountByLabel1 根据 label1 查询题目数量
func (m *QuestionBankMapper) GetQuestionBankCountByLabel1(label1 string) (int, error) {
	var count int64
	result := m.db.Model(&entity.QuestionBank{}).Where("label_1 =?", label1).Count(&count)
	return int(count), result.Error
}

// GetDistinctScoreFromQuestionBank 获取题库表的不同分数
func (m *QuestionBankMapper) GetDistinctScoreFromQuestionBank() ([]float64, error) {
	var scores []float64
	result := m.db.Model(&entity.QuestionBank{}).Distinct("score").Pluck("score", &scores)
	return scores, result.Error
}

// GetQuestionBankCountByScore 根据分数查询题目数量
func (m *QuestionBankMapper) GetQuestionBankCountByScore(score float64) (int, error) {
	var count int64
	result := m.db.Model(&entity.QuestionBank{}).Where("score =?", score).Count(&count)
	return int(count), result.Error
}
