package mapper

import (
	"graduation/entity"
	"time"

	"gorm.io/gorm"
)

// TestPaperGenHistoryGormMapper 实现测试试卷生成历史映射器接口
type TestPaperGenHistoryGormMapper struct {
	db *gorm.DB
}

// NewTestPaperGenHistoryGormMapper 创建测试试卷生成历史 GORM 映射器实例
func NewTestPaperGenHistoryGormMapper() *TestPaperGenHistoryGormMapper {
	return &TestPaperGenHistoryGormMapper{db: DB}
}

// InsertTestPaperGenHistory 插入测试试卷生成历史记录
func (m *TestPaperGenHistoryGormMapper) InsertTestPaperGenHistory(testPaperGenHistory entity.TestPaperGenHistory) (int64, error) {
	result := m.db.Create(&testPaperGenHistory)
	return result.RowsAffected, result.Error
}

// QueryAllTestPaperGenHistory 查询所有测试试卷生成历史记录
func (m *TestPaperGenHistoryGormMapper) QueryAllTestPaperGenHistory() ([]entity.TestPaperGenHistory, error) {
	var histories []entity.TestPaperGenHistory
	result := m.db.Order("update_time desc").Find(&histories)
	return histories, result.Error
}

// GetTestPaperNameByTestPaperUid 根据测试试卷 UID 获取测试试卷名称
func (m *TestPaperGenHistoryGormMapper) GetTestPaperNameByTestPaperUid(testPaperUID string) ([]string, error) {
	var names []string
	result := m.db.Model(&entity.TestPaperGenHistory{}).Where("test_paper_uid = ?", testPaperUID).Pluck("test_paper_name", &names)
	return names, result.Error
}

// UpdateTestPaperTime 更新测试试卷的更新时间
func (m *TestPaperGenHistoryGormMapper) UpdateTestPaperTime(testPaperUID string, date time.Time) (int64, error) {
	result := m.db.Model(&entity.TestPaperGenHistory{}).Where("test_paper_uid = ?", testPaperUID).Update("update_time", date)
	return result.RowsAffected, result.Error
}

// DeleteTestPaperGenHistoryByTestPaperUid 根据测试试卷 UID 删除测试试卷生成历史记录
func (m *TestPaperGenHistoryGormMapper) DeleteTestPaperGenHistoryByTestPaperUid(testPaperUID string) (int64, error) {
	result := m.db.Where("test_paper_uid = ?", testPaperUID).Delete(&entity.TestPaperGenHistory{})
	return result.RowsAffected, result.Error
}

// CreateTestPaperGenHistory creates a new test paper generation history record
func CreateTestPaperGenHistory(history *entity.TestPaperGenHistory) error {
	return DB.Create(history).Error
}
