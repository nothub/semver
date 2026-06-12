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
			name:    "next major",
			args:    args{mode: "major", str: "1.2.3"},
			want:    "2.0.0",
			wantErr: false,
		},
		{
			name:    "next minor",
			args:    args{mode: "minor", str: "1.2.3"},
			want:    "1.3.0",
			wantErr: false,
		},
		{
			name:    "next patch",
			args:    args{mode: "patch", str: "1.2.3"},
			want:    "1.2.4",
			wantErr: false,
		},
		{
			name:    "next strips pre-release and build",
			args:    args{mode: "patch", str: "1.2.3-alpha+build"},
			want:    "1.2.4",
			wantErr: false,
		},
		{
			name:    "next invalid version",
			args:    args{mode: "major", str: "-0.0.0"},
			wantErr: true,
		},
		{
			name:    "next invalid mode",
			args:    args{mode: "bogus", str: "1.2.3"},
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

func Test_tags(t *testing.T) {
	tests := []struct {
		name    string
		str     string
		want    string
		wantErr bool
	}{
		{
			name: "release",
			str:  "1.2.3",
			want: "1.2.3 1.2 1",
		},
		{
			name: "pre-release appended to all three tags",
			str:  "1.2.3-alpha.1",
			want: "1.2.3-alpha.1 1.2-alpha.1 1-alpha.1",
		},
		{
			name: "build metadata stripped",
			str:  "1.2.3+build.5",
			want: "1.2.3 1.2 1",
		},
		{
			name: "pre-release kept build stripped",
			str:  "1.2.3-rc.1+build.5",
			want: "1.2.3-rc.1 1.2-rc.1 1-rc.1",
		},
		{
			name:    "invalid version",
			str:     "-0.0.0",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tags(tt.str)
			if (err != nil) != tt.wantErr {
				t.Errorf("tags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("tags() got = %v, want %v", got, tt.want)
			}
		})
	}
}
