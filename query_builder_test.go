package goje

import (
	"testing"
)

func TestSQLConditionBuilder(t *testing.T) {
	type args struct {
		Queries []QueryInterface
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   int
		wantErr bool
	}{
		{
			name:    "Empty",
			args:    args{},
			want:    " ",
			want1:   0,
			wantErr: false,
		},
		{
			name: "Simple one argument",
			args: args{
				Queries: []QueryInterface{
					Where("id=?", 1),
				},
			},
			want:    "  WHERE id=?",
			want1:   1,
			wantErr: false,
		},
		{
			name: "Complex arguments with group,order,limit,offset",
			args: args{
				Queries: []QueryInterface{
					Where("id=?", 1),
					Order("id DESC"),
					Limit(100),
					Offset(1),
					GroupBy("id"),
					GroupBy("name"),
					Having("id > 1"),
					Having("LENGTH(name) = 1"),
				},
			},
			want:    "  WHERE id=? GROUP BY `id`,`name` HAVING id > 1 AND LENGTH(name) = 1 ORDER BY id DESC LIMIT ? OFFSET ?",
			want1:   3,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := SQLConditionBuilder(tt.args.Queries)
			if (err != nil) != tt.wantErr {
				t.Errorf("SQLConditionBuilder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SQLConditionBuilder() got = \"%v\", want \"%v\"", got, tt.want)
			}
			if len(got1) != tt.want1 {
				t.Errorf("SQLConditionBuilder() len(got1) = %v, want len = %v", got1, tt.want1)
			}
		})
	}
}

func TestSelectQueryBuilder(t *testing.T) {
	type args struct {
		Tablename string
		Columns   []string
		Queries   []QueryInterface
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   int
		wantErr bool
	}{
		{
			name: "multiple args",
			args: args{
				Tablename: "users",
				Queries: []QueryInterface{
					Where("user_id=?", 1),
					Order("user_id DESC"),
					Limit(100),
					Offset(1),
					GroupBy("user_id"),
					GroupBy("name"),
					Having("user_id > 1"),
					Having("LENGTH(name) = 1"),
					InnerJoin("baskets", "baskets.user_id = users.id"),
				},
				Columns: []string{
					"user_id",
					"name",
					"GROUP_CONCAT(baskets.products) as pids",
				},
			},
			want:    "SELECT `user_id`,`name`,GROUP_CONCAT(baskets.products) as pids  FROM users  INNER JOIN baskets ON baskets.user_id = users.id  WHERE user_id=? GROUP BY `user_id`,`name` HAVING user_id > 1 AND LENGTH(name) = 1 ORDER BY user_id DESC LIMIT ? OFFSET ?",
			want1:   3,
			wantErr: false,
		},
		{
			name: "multiple args with joins",
			args: args{
				Tablename: "users",
				Queries: []QueryInterface{
					Where("user_id=?", 1),
					Order("user_id DESC"),
					Limit(100),
					Offset(1),
					GroupBy("user_id"),
					GroupBy("name"),
					Having("user_id > 1"),
					Having("LENGTH(name) = 1"),
					LeftJoin("products", "baskets.product_id = products.id"),
					InnerJoin("baskets", "baskets.user_id = users.id"),
				},
				Columns: []string{
					"user_id",
					"name",
					"GROUP_CONCAT(baskets.products) as pids",
				},
			},
			want:    "SELECT `user_id`,`name`,GROUP_CONCAT(baskets.products) as pids  FROM users  LEFT JOIN products ON baskets.product_id = products.id  INNER JOIN baskets ON baskets.user_id = users.id  WHERE user_id=? GROUP BY `user_id`,`name` HAVING user_id > 1 AND LENGTH(name) = 1 ORDER BY user_id DESC LIMIT ? OFFSET ?",
			want1:   3,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := SelectQueryBuilder(tt.args.Tablename, tt.args.Columns, tt.args.Queries)
			if (err != nil) != tt.wantErr {
				t.Errorf("SelectQueryBuilder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SelectQueryBuilder() got = %v, want %v", got, tt.want)
			}
			if len(got1) != tt.want1 {
				t.Errorf("SelectQueryBuilder() got1 = len(%v), want len = %v", got1, tt.want1)
			}
		})
	}
}
