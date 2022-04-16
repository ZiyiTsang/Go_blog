package article

import (
	"Go_blog/pkg/model"
	"Go_blog/pkg/typesTool"
)

func Get(id_str string) (Article, error) {
	var article Article
	id := typesTool.StringToint64(id_str)
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
