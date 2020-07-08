package common

import (
	"log"
	"strconv"
)

type QueryParams struct {
	Limit  int    `json:"limit,omitempty"`
	Offset int    `json:"offset,omitempty"`
	Sort   string `json:"sort,omitempty"`
	Expand bool   `json:"expand,omitempty"`
}

type TaskQueryParams struct {
	QueryParams
	Nodes []string
	Jobs  []string
}

func (q *QueryParams) Query(table string) string {
	log.Println(q)

	switch q.Limit {
	case 0:
		if q.Offset != 0 {
			return "SELECT jsonb_agg(id) FROM " + table + "LIMIT 20 OFFSET" + strconv.Itoa(q.Offset*20) + ";"
		}
	default:
		if q.Offset != 0 {
			return "SELECT jsonb_agg(id) FROM " + table + "LIMIT " + strconv.Itoa(q.Limit) + " OFFSET" + strconv.Itoa(q.Offset*q.Limit) + ";"
		}
		return "SELECT jsonb_agg(id) FROM " + table + " LIMIT " + strconv.Itoa(q.Limit) + ";"
	}
	return "SELECT jsonb_agg(id) FROM " + table + " LIMIT 20;"
}
