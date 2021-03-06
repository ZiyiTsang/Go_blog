package article

import (
	"Go_blog/pkg/logTool"
	"Go_blog/pkg/model"
	"Go_blog/pkg/typesTool"
)

func Get(idStr string) (Article, error) {
	var article Article
	id := typesTool.StringToint64(idStr)
	if err := model.DB.First(&article, id).Error; err != nil {
		return article, err
	}
	return article, nil
}

func GetAll() ([]Article, error) {
	var articles []Article
	if err := model.DB.Find(&articles).Error; err != nil {
		return articles, err
	}
	return articles, nil
}

func (a Article) CreateWithTitleBody() (int64, error) {
	tObj := model.DB.Exec("insert into articles(title,body,time) VALUES (?,?,now())", a.Title, a.Body)
	row := tObj.RowsAffected
	if err := tObj.Error; err != nil {
		return 0, err
	}
	return row, nil
}
func (a Article) Create() (int64, error) {
	tObj := model.DB.Create(&a)
	row := tObj.RowsAffected
	if err := tObj.Error; err != nil {
		return 0, err
	}
	return row, nil
}

func (article *Article) Update() (rowsAffected int64, err error) {
	result := model.DB.Save(&article)
	if err = result.Error; err != nil {
		logTool.CheckError(err)
		return 0, err
	}
	return result.RowsAffected, nil
}
func (article Article) Delete() (rowsAffected int64, err error) {
	result := model.DB.Delete(&article)
	if err = result.Error; err != nil {
		logTool.CheckError(err)
		return 0, err
	}

	return result.RowsAffected, nil
}
