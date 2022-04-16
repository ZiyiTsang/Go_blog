package article

import (
	"Go_blog/pkg/route"
	"strconv"
)

type Article struct {
	ID    int64
	Title string
	Body  string
	Time  string
}

func (a Article) Link() string {
	return route.Name2URL("articles.show", "id", strconv.FormatInt(a.ID, 10))
}
