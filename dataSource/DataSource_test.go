package dataSource

import "testing"

func Test_prepareMysqlDatabase(t *testing.T) {

}

func Test_extractDatabaseName(t *testing.T) {
	type args struct {
		mysqlDriverArgs string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
	}{
		{name: "case1", args: args{mysqlDriverArgs: "root:P@4word@(test.xxxxx.com:3308)/dbname?charset=utf8mb4"},
			want: "dbname", want1: "root:P@4word@(test.xxxxx.com:3308)/?charset=utf8mb4"},
		{name: "case2", args: args{mysqlDriverArgs: "root:P@4word@(test.xxxxx.com:3308)/?charset=utf8mb4"},
			want: "", want1: "root:P@4word@(test.xxxxx.com:3308)/?charset=utf8mb4"},
		{name: "case3", args: args{mysqlDriverArgs: "root:P@4word@(test.xxxxx.com:3308)?charset=utf8mb4"},
			want: "", want1: "root:P@4word@(test.xxxxx.com:3308)?charset=utf8mb4"},
		{name: "case4", args: args{mysqlDriverArgs: "root:P@4word@(test.xxxxx.com:3308)/dbname"},
			want: "dbname", want1: "root:P@4word@(test.xxxxx.com:3308)/"},
		{name: "case5", args: args{mysqlDriverArgs: "root:P@4word@(test.xxxxx.com:3308)/"},
			want: "", want1: "root:P@4word@(test.xxxxx.com:3308)/"},
		{name: "case6", args: args{mysqlDriverArgs: "root:P@4word@(test.xxxxx.com:3308)"},
			want: "", want1: "root:P@4word@(test.xxxxx.com:3308)"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := extractDatabaseName(tt.args.mysqlDriverArgs)
			if got != tt.want {
				t.Errorf("extractDatabaseName() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("extractDatabaseName() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
