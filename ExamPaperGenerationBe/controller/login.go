package controller

import (
	"fmt"
	"graduation/entity"
	"graduation/mapper"
	"graduation/utils"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// 处理 /permission_denied 请求
func PermissionDenied(c *gin.Context) {
	response := utils.Make403Resp("Permission denied")
	c.String(http.StatusOK, response)
}

// 处理 /getLoginStatus 请求
func GetLoginStatus(c *gin.Context) {
	session := sessions.Default(c)
	username := session.Get("username")
	userRole := session.Get("user_role")
	lastLogin := session.Get("last_login")

	if username != nil && userRole != nil && lastLogin != nil {
		retData := map[string]interface{}{
			"username":   username,
			"user_role":  userRole,
			"last_login": lastLogin,
		}
		response := utils.Make200Resp("Success", retData)
		c.String(http.StatusOK, response)
	} else {
		response := utils.Make403Resp("Permission denied")
		c.String(http.StatusOK, response)
	}
}

// 处理 /Login 请求
func Login(c *gin.Context) {
	var user entity.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user.LastLogin = time.Now()
	var userList []entity.User
	//mapper.DB.Where("username = ? AND password = ?", user.Username, user.Password).Find(&userList)
	mapper.DB.Where("username = ? and enable = 1", user.Username).Find(&userList)
	if len(userList) > 0 {
		if !utils.DecPassword(user.Password, userList[0].Password) { // 密码错误
			response := utils.Make403Resp("Permission denied")
			c.String(http.StatusOK, response)
			return
		}
		mapper.DB.Model(&user).Where("username = ?", user.Username).Update("last_login", user.LastLogin)
		session := sessions.Default(c)
		session.Set("username", userList[0].Username)
		session.Set("user_role", userList[0].UserRole)
		session.Set("last_login", userList[0].LastLogin)
		if err := session.Save(); err != nil {
			response := utils.Make500Resp("保存session失败")
			c.String(http.StatusInternalServerError, response)
			return
		}
		retData := map[string]interface{}{
			"username":   userList[0].Username,
			"user_role":  userList[0].UserRole,
			"last_login": userList[0].LastLogin.Format(time.RFC3339),
		}
		response := utils.Make200Resp("Success", retData)
		c.String(http.StatusOK, response)
	} else {
		response := utils.Make403Resp("Permission denied")
		c.String(http.StatusOK, response)
	}
}

// 处理 /logout 请求
func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	response := utils.Make200Resp("Success", "退出登陆成功")
	c.String(http.StatusOK, response)
}

// 处理 /registered 请求
func Registered(c *gin.Context) {
	var user entity.User
	var err error
	if err = c.ShouldBindJSON(&user); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		fmt.Println(err.Error())
		return
	}
	var userByUsername []entity.User
	mapper.DB.Where("username = ?", user.Username).Find(&userByUsername)
	if len(userByUsername) > 0 {
		response := utils.Make500Resp("用户名重复")
		c.String(http.StatusInternalServerError, response)
		fmt.Println(response)
		return
	}
	user.LastLogin = time.Now()
	user.Enable = 0
	user.Password, err = utils.EncPassword(user.Password)
	user.SimilarityThreshold = 0.5
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		fmt.Println(err.Error())
		return
	}
	result := mapper.DB.Create(&user)
	if result.Error != nil {
		response := utils.Make500Resp("注册失败")
		c.String(http.StatusInternalServerError, response)
		fmt.Println(response)
		return
	}
	response := utils.Make200Resp("Success", result.RowsAffected)
	c.String(http.StatusOK, response)
}

// 处理 /getApplyUser 请求
func GetApplyUser(c *gin.Context) {
	session := sessions.Default(c)
	userRole := session.Get("user_role")
	if userRole == nil || userRole.(string) != "admin" {
		response := utils.Make403Resp("Permission denied")
		c.String(http.StatusForbidden, response)
		return
	}
	var applyUser []entity.User
	mapper.DB.Where("enable = 0").Find(&applyUser)
	response := utils.Make200Resp("Success", applyUser)
	c.String(http.StatusOK, response)
}

// 处理 /getAllUser 请求
func GetAllUser(c *gin.Context) {
	var allUser []entity.User
	mapper.DB.Select("id, username, user_role, last_login, enable").
		Find(&allUser)
	response := utils.Make200Resp("Success", allUser)
	c.String(http.StatusOK, response)
}

// 处理 /deleteUser 请求
func DeleteUser(c *gin.Context) {
	username := c.Query("username")
	result := mapper.DB.Where("username = ?", username).Delete(&entity.User{})
	response := utils.Make200Resp("Success", result.RowsAffected)
	c.String(http.StatusOK, response)
}

// 处理 /passApply 请求
func PassApply(c *gin.Context) {
	username := c.Query("username")
	result := mapper.DB.Model(&entity.User{}).Where("username = ?", username).Update("enable", 1)
	response := utils.Make200Resp("Success", result.RowsAffected)
	c.String(http.StatusOK, response)
}

// 处理 /deleteApply 请求
func DeleteApply(c *gin.Context) {
	username := c.Query("username")
	result := mapper.DB.Where("username = ? AND enable = 0", username).Delete(&entity.User{})
	response := utils.Make200Resp("Success", result.RowsAffected)
	c.String(http.StatusOK, response)
}
