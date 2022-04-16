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

func (a Article) CreateWithTitleBody() (int64, error) {
	t_obj := model.DB.Exec("insert into articles(title,body,time) VALUES (?,?,now())", a.Title, a.Body)
	row := t_obj.RowsAffected
	if err := t_obj.Error; err != nil {
		return 0, err
	}
	return row, nil
}
