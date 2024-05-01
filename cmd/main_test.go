package main

import "testing"

func Test_next(t *testing.T) {
	type args struct {
		mode string
		str  string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "next major",
			args: args{
				mode: "major",
				str:  "0.0.0",
			},
			want:    "1.0.0",
			wantErr: false,
		},
		{
			name: "next minor",
			args: args{
				mode: "minor",
				str:  "0.0.0",
			},
			want:    "0.1.0",
			wantErr: false,
		},
		{
			name: "next patch",
			args: args{
				mode: "patch",
				str:  "0.0.0",
			},
			want:    "0.0.1",
			wantErr: false,
		},
		{
			name: "next invalid",
			args: args{
				mode: "major",
				str:  "-0.0.0",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := next(tt.args.mode, tt.args.str)
			if (err != nil) != tt.wantErr {
				t.Errorf("next() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("next() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_strip(t *testing.T) {
	type args struct {
		mode string
		str  string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "strip all",
			args: args{
				mode: "all",
				str:  "0.0.0-x+y",
			},
			want:    "0.0.0",
			wantErr: false,
		},
		{
			name: "strip pre",
			args: args{
				mode: "pre",
				str:  "0.0.0-x+y",
			},
			want:    "0.0.0+y",
			wantErr: false,
		},
		{
			name: "strip build",
			args: args{
				mode: "build",
				str:  "0.0.0-x+y",
			},
			want:    "0.0.0-x",
			wantErr: false,
		},
		{
			name: "strip invalid",
			args: args{
				mode: "all",
				str:  "-0.0.0",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := strip(tt.args.mode, tt.args.str)
			if (err != nil) != tt.wantErr {
				t.Errorf("strip() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("strip() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_valid(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "valid",
			args:    args{str: "0.0.0-x+y"},
			wantErr: false,
		},
		{
			name:    "invalid",
			args:    args{str: "-0.0.0"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := valid(tt.args.str); (err != nil) != tt.wantErr {
				t.Errorf("valid() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
