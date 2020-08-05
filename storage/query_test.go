package storage

import (
	"github.com/chabad360/covey/models"
	"github.com/google/go-cmp/cmp"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestQueryParams_Query(t *testing.T) {
	DB.Delete(&models.Job{}, "id != ''")
	AddJob(j)
	AddJob(j2)
	type fields struct {
		Limit  int
		Offset int
		Sort   string
		SortBy string
		Expand bool
	}
	type args struct {
		table string
		model interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		want    interface{}
	}{
		{"success", fields{10, 0, "asc", "name", false}, args{"jobs", &[]string{}}, false, &[]string{j2.ID, j.ID}},
		{"limitAndOffset", fields{1, 1, "asc", "name", false}, args{"jobs", &[]string{}}, false, &[]string{j.ID}},
		{"fail", fields{0, 1, "asc", "name", false}, args{"jobs", &[]string{}}, true, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Query{
				Limit:  tt.fields.Limit,
				Offset: tt.fields.Offset,
				Sort:   tt.fields.Sort,
				SortBy: tt.fields.SortBy,
				Expand: tt.fields.Expand,
			}
			if err := q.Query(tt.args.table, tt.args.model); (err != nil) != tt.wantErr {
				t.Errorf("Query() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !cmp.Equal(tt.want, tt.args.model) && !tt.wantErr {
				t.Errorf("Query() = %v, want %v", tt.args.model, tt.want)
			}
		})
	}
}

func TestQueryParams_Setup(t *testing.T) {
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    Query
	}{
		// revive:disable:line-length-limit
		{"work", args{httptest.NewRequest("GET", "/test?limit=2&offset=2&expand=true&sort=desc&sortby=name", nil)}, false, Query{2, 2, "desc", "name", true}},
		// revive:enable:line-length-limit
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var q Query
			if err := q.Setup(tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("Setup() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !cmp.Equal(q, tt.want) {
				t.Errorf("Setup() = %v, want %v", q, tt.want)
			}
		})
	}
}
