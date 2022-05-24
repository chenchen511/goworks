package controller

import (
	"goworks/common"
	"goworks/model"
	"goworks/response"

	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ICategoryController interface {
	RestController
}
type CategoryController struct {
	DB *gorm.DB
}

func NeCategoryController() ICategoryController {
	db := common.GetDB()
	db.AutoMigrate(model.Category{})
	return CategoryController{DB: db}
}

func (c CategoryController) Create(ctx *gin.Context) {
	var requesteCategory model.Category
	ctx.Bind(&requesteCategory)
	if requesteCategory.Name == "" {
		response.Fail(ctx, nil, "数据验证错误，分类名必填")
	}
	c.DB.Create(&requesteCategory)
	response.Success(ctx, gin.H{"category": requesteCategory}, "")
}
func (c CategoryController) Update(ctx *gin.Context) {

	var requesteCategory model.Category
	ctx.Bind(&requesteCategory)
	if requesteCategory.Name == "" {
		response.Fail(ctx, nil, "数据验证错误，分类名必填")
	}
	categoryId, _ := strconv.Atoi(ctx.Params.ByName("id"))
	var updateCategory model.Category
	if err := c.DB.First(&updateCategory, categoryId).Error; err != nil {
		response.Fail(ctx, nil, "分类不存在")
	}
	c.DB.Model(&updateCategory).Update("name", requesteCategory.Name)

	response.Success(ctx, gin.H{"category": updateCategory}, "修改成功")

}
func (c CategoryController) Show(ctx *gin.Context) {
	categoryId, _ := strconv.Atoi(ctx.Params.ByName("id"))
	var category model.Category
	if err := c.DB.First(&category, categoryId).Error; err != nil {
		response.Fail(ctx, nil, "分类不存在")
	}
	response.Success(ctx, gin.H{"category": category}, "")
}

func (c CategoryController) Delete(ctx *gin.Context) {
	categoryId, _ := strconv.Atoi(ctx.Params.ByName("id"))
	// var category model.Category
	// if err := c.DB.First(&category, categoryId).Error; err != nil {
	// 	response.Fail(ctx, nil, "分类不存在")
	// }
	if err := c.DB.Delete(model.Category{}, categoryId); err != nil {
		response.Fail(ctx, nil, "删除失败")
	}
	response.Success(ctx, nil, "")
}
