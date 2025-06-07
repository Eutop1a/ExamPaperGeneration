package controller

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"graduation/mapper"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetSimilarityThreshold 设置用户的相似度阈值
// @Summary 设置用户的相似度阈值
// @Description 设置当前用户的相似度阈值，用于控制题目重复度
// @Tags 用户设置
// @Accept json
// @Produce json
// @Param threshold body float64 true "相似度阈值(0-1)"
// @Success 200 {object} gin.H
// @Router /user/similarity [put]
func SetSimilarityThreshold(ctx *gin.Context) {
	session := sessions.Default(ctx)
	username := session.Get("username")
	if username == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var threshold struct {
		Threshold float64 `json:"threshold" binding:"required,min=0,max=1"`
	}
	if err := ctx.ShouldBindJSON(&threshold); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新用户的相似度阈值
	if err := updateUserSimilarityThreshold(username.(string), threshold.Threshold); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set similarity threshold"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Similarity threshold updated successfully"})
}

// updateUserSimilarityThreshold 更新用户的相似度阈值
func updateUserSimilarityThreshold(username string, threshold float64) error {
	// 验证阈值范围
	if threshold < 0 || threshold > 1 {
		return fmt.Errorf("threshold must be between 0 and 1")
	}

	// 更新用户相似度阈值
	if err := mapper.UpdateUserSimilarityThreshold(username, threshold); err != nil {
		return fmt.Errorf("failed to update user similarity threshold: %v", err)
	}

	return nil
}
