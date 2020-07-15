package common

import (
	"fmt"
	httpext "github.com/go-playground/pkg/v5/net/http"
	"github.com/go-playground/pure/v5"
	"net/http"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// QueryParams parses a URL query and gets the items from the database based on it.
type QueryParams struct {
	Limit  int    `form:"limit" validate:"min=0,max=50"`
	Offset int    `form:"offset"`
	Sort   string `form:"sort" validate:"oneof=asc desc"`
	SortBy string `form:"sortby"`
	Expand bool   `form:"expand"`
}

// Query runs a query against the database.
func (q *QueryParams) Query(table string, model interface{}, db *gorm.DB) error {
	v := validator.New()
	err := v.Struct(q)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}
		for _, err := range err.(validator.ValidationErrors) {
			return fmt.Errorf("query: invalid value %v for field %v", err.Value(), err.Field())
		}
	}

	tx := db.Table(table).Offset(q.Offset).Limit(q.Limit).Order(fmt.Sprintf("%s %s", q.SortBy, q.Sort))

	if !q.Expand {
		tx.Select("id")
	}

	tx.Scan(model)
	if tx.Error != nil {
		return err
	}

	return nil
}

// Setup parses a url query.
func (q *QueryParams) Setup(r *http.Request) error {
	q.Limit = 20
	q.Offset = 0
	q.Sort = "asc"
	q.SortBy = "id"
	q.Expand = false

	return pure.DecodeQueryParams(r, httpext.QueryParams, &q)
}
