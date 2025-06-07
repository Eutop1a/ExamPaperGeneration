package entity

import (
	"gorm.io/gorm"
	"time"
)

// QuestionGenHistory 表示问题生成历史实体
type QuestionGenHistory struct {
	ID              int       `gorm:"primaryKey;column:id" json:"id"`
	TestPaperUID    string    `gorm:"column:test_paper_uid" json:"test_paper_uid"`
	TestPaperName   string    `gorm:"column:test_paper_name" json:"test_paper_name"`
	QuestionBankID  int       `gorm:"column:question_bank_id" json:"question_bank_id"`
	Topic           string    `gorm:"column:topic" json:"topic"`
	TopicMaterialID int       `gorm:"column:topic_material_id" json:"topic_material_id"`
	Answer          string    `gorm:"column:answer" json:"answer"`
	TopicType       string    `gorm:"column:topic_type" json:"topic_type"`
	Score           float64   `gorm:"column:score" json:"score"`
	Difficulty      int       `gorm:"column:difficulty" json:"difficulty"`
	Chapter1        string    `gorm:"column:chapter_1" json:"chapter_1"`
	Chapter2        string    `gorm:"column:chapter_2" json:"chapter_2"`
	Label1          string    `gorm:"column:label_1" json:"label_1"`
	Label2          string    `gorm:"column:label_2" json:"label_2"`
	UpdateTime      time.Time `gorm:"column:update_time" json:"update_time"`
}

// BeforeCreate 在创建记录前设置更新时间
func (q *QuestionGenHistory) BeforeCreate(tx *gorm.DB) error {
	q.UpdateTime = time.Now()
	return nil
}

func (q *QuestionGenHistory) TableName() string {
	return "questiongenhistory" // 明确指定表名
}
