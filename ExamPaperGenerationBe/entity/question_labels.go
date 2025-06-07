package entity

import "time"

// QuestionLabels 表示题目标签实体
type QuestionLabels struct {
	ID        int       `gorm:"primaryKey;column:id" json:"id"`
	Chapter1  string    `gorm:"column:chapter_1" json:"chapter_1"`
	Chapter2  string    `gorm:"column:chapter_2" json:"chapter_2"`
	Label1    string    `gorm:"column:label_1" json:"label_1"`
	Label2    string    `gorm:"column:label_2" json:"label_2"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (q *QuestionLabels) TableName() string {
	return "questionlabels" // 明确指定表名
}
