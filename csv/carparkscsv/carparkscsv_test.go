package carparkscsv

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name    string
		args    args
		want    []Row
		wantErr bool
	}{
		{
			name: "missing file",
			args: args{
				filePath: "nonexistingfile.csv",
			},
			want:    []Row{},
			wantErr: true,
		},
		{
			name: "empty file",
			args: args{
				filePath: "csv/emptyfile.csv",
			},
			want:    []Row{},
			wantErr: false,
		},
		{
			name: "header but no rows",
			args: args{
				filePath: "csv/onlyheader.csv",
			},
			want:    []Row{},
			wantErr: false,
		},
		{
			name: "happy path - header + 1 row",
			args: args{
				filePath: "csv/headerplusone.csv",
			},
			want: []Row{
				{
					Number:  "ABC",
					Address: "Home",
					XCoord:  1.0,
					YCoord:  1.0,
				},
			},
			wantErr: false,
		},
		{
			name: "happy path - header + multiple rows",
			args: args{
				filePath: "csv/many.csv",
			},
			want: []Row{
				{
					Number:  "XXX",
					Address: "Home 1",
					XCoord:  1.1,
					YCoord:  22.1,
				},
				{
					Number:  "YYY",
					Address: "Home 2",
					XCoord:  333.2,
					YCoord:  4444.2,
				},
				{
					Number:  "ZZZ",
					Address: "Home 3",
					XCoord:  55555.3,
					YCoord:  666666.3,
				},
				{
					Number:  "ZZZ",
					Address: "Home 3",
					XCoord:  55555.3,
					YCoord:  666666.3,
				},
			},
			wantErr: false,
		},
		{
			name: "invalid number",
			args: args{
				filePath: "csv/invalid1.csv",
			},
			want:    []Row{},
			wantErr: true,
		},
		{
			name: "invalid address",
			args: args{
				filePath: "csv/invalid2.csv",
			},
			want:    []Row{},
			wantErr: true,
		},
		{
			name: "invalid x coord",
			args: args{
				filePath: "csv/invalid3.csv",
			},
			want:    []Row{},
			wantErr: true,
		},
		{
			name: "invalid y coord",
			args: args{
				filePath: "csv/invalid4.csv",
			},
			want:    []Row{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseCSV() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseCSV() got = %v, want %v", got, tt.want)
			}
		})
	}
}
