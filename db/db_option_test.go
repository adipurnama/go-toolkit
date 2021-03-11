package db_test

import (
	"reflect"
	"testing"

	"github.com/adipurnama/go-toolkit/db"
)

func TestNewDatabaseOption(t *testing.T) {
	type args struct {
		host     string
		port     int
		username string
		password string
		dbName   string
	}

	tests := []struct {
		name    string
		args    args
		want    *db.Option
		wantErr bool
	}{
		{
			"host empty, return error",
			args{
				port: 1212,
			},
			nil,
			true,
		},
		{
			"port empty, return error",
			args{
				host: "localhost",
			},
			nil,
			true,
		},
		{
			"port & host exists, no error",
			args{
				host: "localhost",
				port: 1212,
			},
			&db.Option{
				Host:             "localhost",
				Port:             1212,
				ConnectionOption: db.DefaultConnectionOption(),
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := db.NewDatabaseOption(tt.args.host, tt.args.port, tt.args.username, tt.args.password, tt.args.dbName, nil)

			if (err != nil) != tt.wantErr {
				t.Errorf("NewDatabaseOption() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got == nil {
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDatabaseOption() = %v, want %v", got, tt.want)
			}
		})
	}
}
