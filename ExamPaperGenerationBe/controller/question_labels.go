package controller

import (
	"graduation/entity"
	"graduation/mapper"
	"graduation/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// GetAllQuestionLabels 获取所有题目标签
func GetAllQuestionLabels(c *gin.Context) {
	var questionLabels []entity.QuestionLabels
	mapper.DB.Find(&questionLabels)
	resp := utils.Make200Resp("successfully get all question labels", questionLabels)
	c.String(http.StatusOK, resp)
}

// GetDistinctChapter1 获取不同的 chapter_1
func GetDistinctChapter1(c *gin.Context) {
	var distinctChapter1 []entity.QuestionLabels
	mapper.DB.Distinct("chapter_1").Find(&distinctChapter1)
	var chapter1List []string
	for _, label := range distinctChapter1 {
		chapter1List = append(chapter1List, label.Chapter1)
	}
	resp := utils.Make200Resp("successfully get chapter1", chapter1List)
	c.String(http.StatusOK, resp)
}

// GetDistinctChapter2 获取不同的 chapter_2
func GetDistinctChapter2(c *gin.Context) {
	var distinctChapter2 []entity.QuestionLabels
	mapper.DB.Distinct("chapter_2").Find(&distinctChapter2)
	var chapter2List []string
	for _, label := range distinctChapter2 {
		chapter2List = append(chapter2List, label.Chapter2)
	}
	resp := utils.Make200Resp("successfully get chapter2", chapter2List)
	c.String(http.StatusOK, resp)
}

// GetChapter2ByChapter1 根据 chapter_1 获取对应的 chapter_2
func GetChapter2ByChapter1(c *gin.Context) {
	chapter1 := c.Query("chapter1")
	var chapter2ByChapter1 []entity.QuestionLabels
	mapper.DB.Where("chapter_1 = ?", chapter1).Find(&chapter2ByChapter1)
	var chapter2List []string
	for _, label := range chapter2ByChapter1 {
		chapter2List = append(chapter2List, label.Chapter2)
	}
	resp := utils.Make200Resp("successfully get chapter2 by chapter1", chapter2List)
	c.String(http.StatusOK, resp)
}

// GetDistinctLabel1 获取不同的 label_1
func GetDistinctLabel1(c *gin.Context) {
	var distinctLabel1 []entity.QuestionLabels
	mapper.DB.Distinct("label_1").Find(&distinctLabel1)
	var label1List []string
	for _, label := range distinctLabel1 {
		label1List = append(label1List, label.Label1)
	}
	resp := utils.Make200Resp("successfully get label1", label1List)
	c.String(http.StatusOK, resp)
}

// GetDistinctLabel2 获取不同的 label_2
func GetDistinctLabel2(c *gin.Context) {
	var distinctLabel2 []entity.QuestionLabels
	mapper.DB.Distinct("label_2").Find(&distinctLabel2)
	var label2List []string
	for _, label := range distinctLabel2 {
		label2List = append(label2List, label.Label2)
	}
	resp := utils.Make200Resp("successfully get label2", label2List)
	c.String(http.StatusOK, resp)
}

// CreateLabel 创建标签
// @Summary 创建知识点标签
// @Description 创建新的知识点标签
// @Tags 知识点标签
// @Accept json
// @Produce json
// @Param label body entity.QuestionLabels true "标签信息"
// @Success 200 {object} entity.QuestionLabels
// @Router /labels [post]
func CreateLabel(ctx *gin.Context) {
	var label entity.QuestionLabels
	if err := ctx.ShouldBindJSON(&label); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	now := time.Now()
	label.CreatedAt = now
	label.UpdatedAt = now
	if err := mapper.DB.Model(&entity.QuestionLabels{}).Create(&label).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create the label. Please try again later."})
		return
	}

	ctx.JSON(http.StatusOK, label)
}

// UpdateLabel 更新标签
// @Summary 更新知识点标签
// @Description 更新指定的知识点标签
// @Tags 知识点标签
// @Accept json
// @Produce json
// @Param id path int true "标签ID"
// @Param label body entity.QuestionLabels true "标签信息"
// @Success 200 {object} entity.QuestionLabels
// @Router /labels/{id} [put]
func UpdateLabel(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var label entity.QuestionLabels
	if err := ctx.ShouldBindJSON(&label); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	label.ID = id
	// 先检查记录是否存在
	var existingLabel entity.QuestionLabels
	result := mapper.DB.Model(&entity.QuestionLabels{}).Where("id = ?", id).First(&existingLabel)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Label not found"})
		return
	}

	// 开始事务
	tx := mapper.DB.Begin()
	if tx.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	// 更新标签记录
	if err := tx.Model(&entity.QuestionLabels{}).Where("id = ?", id).Updates(&label).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update label"})
		return
	}

	// 更新题库中对应的标签
	if err := tx.Debug().Model(&entity.QuestionBank{}).
		Where("chapter_1 = ? AND chapter_2 = ? AND label_1 = ? AND label_2 = ?",
			existingLabel.Chapter1, existingLabel.Chapter2, existingLabel.Label1, existingLabel.Label2).
		Updates(map[string]interface{}{
			"chapter_1": label.Chapter1,
			"chapter_2": label.Chapter2,
			"label_1":   label.Label1,
			"label_2":   label.Label2,
		}).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update question bank labels"})
		return
	}

	// 更新历史记录中对应的标签
	if err := tx.Debug().Model(&entity.QuestionGenHistory{}).
		Where("chapter_1 = ? AND chapter_2 = ? AND label_1 = ? AND label_2 = ?",
			existingLabel.Chapter1, existingLabel.Chapter2, existingLabel.Label1, existingLabel.Label2).
		Updates(map[string]interface{}{
			"chapter_1": label.Chapter1,
			"chapter_2": label.Chapter2,
			"label_1":   label.Label1,
			"label_2":   label.Label2,
		}).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update question history labels"})
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	ctx.JSON(http.StatusOK, existingLabel)
}

// DeleteLabel 删除标签
// @Summary 删除知识点标签
// @Description 删除指定的知识点标签
// @Tags 知识点标签
// @Produce json
// @Param id path int true "标签ID"
// @Success 200 {object} gin.H
// @Router /labels/{id} [delete]
func DeleteLabel(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var label entity.QuestionLabels
	if err := mapper.DB.Model(&label).Where("id = ?", id).Delete(&label).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Label deleted successfully"})
}
