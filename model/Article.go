package model

import (
	"Bluebell/utils/errmsg"
	"github.com/jinzhu/gorm"
)

type Article struct {
	Category Category `gorm:"foreignkey:Cid"`
	gorm.Model
	Title   string `gorm:"type:varchar(100);not null" json:"title"`
	Cid     int    `gorm:"type:int" json:"cid"`
	Desc    string `gorm:"type:varchar(200)" json:"desc"`
	Content string `gorm:"type:longtext" json:"content"`
	Img     string `gorm:"type:varchar(100)" json:"img"`
}

// 新增文章
func CreateArt(data *Article) int {
	err := db.Create(&data).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}

// 查询分类下的所有文章
func GetCateArt(id int, pageSize int, pageNum int) ([]Article, int, int) {
	var caetArtList []Article
	var total int
	err = db.Preload("Category").Limit(pageSize).Offset((pageNum-1)*pageSize).Where("cid=?", id).Find(&caetArtList).Count(&total).Error
	if err != nil {
		return nil, errmsg.ERROR_CATE_NOT_EXIST, 0
	}
	return caetArtList, errmsg.SUCCESS, total
}

// 查询单个文章
func GetArtInfo(id int) (Article, int) {
	var art Article
	err = db.Preload("Category").Where("id=?", id).First(&art).Error
	if err != nil {
		return art, errmsg.ERROR_ART_NOT_EXIST
	}
	return art, errmsg.SUCCESS
}

// 查询文章列表
func GetArt(title string, pageSize int, pageNum int) ([]Article, int, int) {
	var articleList []Article
	var total int
	var err error
	if title == "" {
		err = db.Order("Updated_At DESC").Preload("Category").Find(&articleList).Limit(pageSize).Offset((pageNum - 1) * pageSize).Error
		db.Model(&articleList).Count(&total)
		//err = db.Preload("Category").Find(&articleList).Count(&total).Limit(pageSize).Offset((pageNum - 1) * pageSize).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, errmsg.ERROR, 0
		}
		return articleList, errmsg.SUCCESS, total
	}
	err = db.Order("Update_At DESC").Preload("Category").Where("title LIKE ?", title+"%").
		Find(&articleList).Limit(pageSize).Offset((pageNum - 1) * pageSize).Error
	//err = dbCount(&total).Error
	db.Model(&articleList).Count(&total)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, errmsg.ERROR, 0
	}
	return articleList, errmsg.SUCCESS, total

}

// 编辑文章
func EditArt(id int, data *Article) int {
	var art Article
	var maps = make(map[string]interface{})
	maps["title"] = data.Title
	maps["cid"] = data.Cid
	maps["desc"] = data.Desc
	maps["content"] = data.Content
	maps["img"] = data.Img
	err = db.Model(&art).Where("id=?", id).Update(maps).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}

// 删除分类
func DeleteArt(id int) int {
	var art Article
	err = db.Where("id=?", id).Delete(&art).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCESS
}
