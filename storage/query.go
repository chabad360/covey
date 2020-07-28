package storage

import (
	"fmt"
	httpext "github.com/go-playground/pkg/v5/net/http"
	"github.com/go-playground/pure/v5"
	"net/http"

	"github.com/go-playground/validator/v10"
)

// Query parses a URL query and gets the items from the database based on it.
type Query struct {
	Limit  int    `form:"limit" validate:"min=1,max=50"`
	Offset int    `form:"offset"`
	Sort   string `form:"sort" validate:"oneof=asc desc"`
	SortBy string `form:"sortby"`
	Expand bool   `form:"expand"`
}

// Query runs a query against the database.
func (q *Query) Query(table string, model interface{}) error {
	v := validator.New()
	errs := v.Struct(q)
	if errs != nil {
		if _, ok := errs.(*validator.InvalidValidationError); ok {
			return errs
		}
		for _, err := range errs.(validator.ValidationErrors) {
			return fmt.Errorf("query: invalid value %v for field %v", err.Value(), err.Field())
		}
	}

	tx := DB.Table(table).Offset(q.Offset).Limit(q.Limit).Order(fmt.Sprintf("%s %s", q.SortBy, q.Sort))

	if !q.Expand {
		tx.Select("id")
	}

	tx.Scan(model)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

// Setup parses a url query.
func (q *Query) Setup(r *http.Request) error {
	q.Limit = 20
	q.Offset = 0
	q.Sort = "asc"
	q.SortBy = "id"
	q.Expand = false

	return pure.DecodeQueryParams(r, httpext.QueryParams, &q)
}
