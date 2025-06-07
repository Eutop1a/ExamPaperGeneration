package entity

import (
	"gorm.io/gorm"
	"time"
)

// QuestionBank 表示问题库实体
type QuestionBank struct {
	ID              int       `gorm:"primaryKey;column:id" json:"id"`
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
	TopicImagePath  string    `gorm:"column:topic_image_path" json:"topic_image_path"`
	UpdateTime      time.Time `gorm:"column:update_time" json:"update_time"`
}

// BeforeCreate 在创建记录前设置更新时间
func (q *QuestionBank) BeforeCreate(tx *gorm.DB) error {
	q.UpdateTime = time.Now()
	return nil
}

func (q *QuestionBank) TableName() string {
	return "questionbank" // 明确指定表名
}
