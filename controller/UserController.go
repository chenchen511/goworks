package controller

import (
	_ "go/token"
	"goworks/common"
	"goworks/dto"
	"goworks/model"
	"goworks/response"
	"goworks/until"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Register(ctx *gin.Context) {
	DB := common.GetDB()
	//获取参数
	//使用map
	// var requesMap =make(map[string]string)
	// json.NewDecoder(ctx.Request.Body).Decode(&requesMap)
	//结构体
	var requesUser = model.User{}
	// json.NewDecoder(ctx.Request.Body).Decode(&requesUser)
	ctx.Bind(&requesUser)
	name := requesUser.Name
	telephone := requesUser.Telephone
	password := requesUser.Password
	if len(telephone) != 11 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "手机号必须为11位")
		return
	}
	if len(password) < 6 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "密码不能少于6位")
		return
	}
	if len(name) == 0 {
		name = until.RandomString(10)
	}
	if isTelephoneExist(DB, telephone) {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 422, "msg": "用户已存在",
		})
		return
	}
	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 422, "msg": "加密错误",
		})
		return
	}
	newUser := model.User{
		Name:      name,
		Telephone: telephone,
		Password:  string(hasedPassword),
	}
	DB.Create(&newUser)

	token, err := common.ReleaseToken(newUser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 500, "msg": "系统异常",
		})
		log.Panicf("token generate error:%v", err)
		return
	}
	response.Success(ctx, gin.H{"token": token}, "注册成功")

}
func Login(ctx *gin.Context) {
	DB := common.GetDB()
	telephone := ctx.PostForm("telephone")
	password := ctx.PostForm("password")
	if len(telephone) != 11 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 422, "msg": "手机号必须为11位",
		})
		return
	}
	if len(password) < 6 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 422, "msg": "密码不能少于6位",
		})
		return
	}
	var user model.User
	DB.Where("telephone = ?", telephone).First(&user)
	if user.ID == 0 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 422, "msg": "用户不存在",
		})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 422, "msg": "密码错误！",
		})
		return
	}
	token, err := common.ReleaseToken(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 500, "msg": "系统异常",
		})
		log.Panicf("token generate error:%v", err)
		return
	}
	response.Success(ctx, gin.H{"token": token}, "登陆成功")
}
func Info(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200, "data": gin.H{"user": dto.ToUserDto(user.(model.User))},
	})
}
func isTelephoneExist(db *gorm.DB, telephone string) bool {
	var user model.User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}
