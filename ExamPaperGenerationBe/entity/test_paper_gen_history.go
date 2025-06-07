package entity

import (
	"time"

	"gorm.io/gorm"
)

// TestPaperGenHistory 表示试卷生成历史实体
type TestPaperGenHistory struct {
	ID                  int       `gorm:"primaryKey;column:id" json:"id"`
	TestPaperUID        string    `gorm:"column:test_paper_uid" json:"test_paper_uid"`
	TestPaperName       string    `gorm:"column:test_paper_name" json:"test_paper_name"`
	QuestionCount       int       `gorm:"column:question_count" json:"question_count"`
	AverageDifficulty   float64   `gorm:"column:average_difficulty" json:"average_difficulty"`
	UpdateTime          time.Time `gorm:"column:update_time" json:"update_time"`
	Username            string    `gorm:"column:username" json:"username"`
	SimilarityThreshold float64   `gorm:"column:similarity_threshold" json:"similarity_threshold"`
}

// BeforeCreate 在创建记录前设置更新时间
func (t *TestPaperGenHistory) BeforeCreate(tx *gorm.DB) error {
	t.UpdateTime = time.Now()
	return nil
}

func (t *TestPaperGenHistory) TableName() string {
	return "testpapergenhistory" // 明确指定表名
}
