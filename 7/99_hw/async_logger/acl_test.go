package main

import (
	"reflect"
	"testing"
)

var ACLDataTest string = `{
	"logger1":          ["/main.Admin/Logging"],
	"logger2":          ["/main.Admin/Logging"],
	"stat1":            ["/main.Admin/Statistics"],
	"stat2":            ["/main.Admin/Statistics"],
	"biz_user":         ["/main.Biz/Check", "/main.Biz/Add"],
	"biz_admin":        ["/main.Biz/*"],
	"after_disconnect": ["/main.Biz/Add"]
}`

func TestNewACL(t *testing.T) {
	ACLData := make(map[string][]string)
	ACLData["logger1"] = []string{"/main.Admin/Logging"}
	ACLData["logger2"] = []string{"/main.Admin/Logging"}
	ACLData["stat1"] = []string{"/main.Admin/Statistics"}
	ACLData["stat2"] = []string{"/main.Admin/Statistics"}
	ACLData["biz_user"] = []string{"/main.Biz/Check", "/main.Biz/Add"}
	ACLData["biz_admin"] = []string{"/main.Biz/*"}
	ACLData["after_disconnect"] = []string{"/main.Biz/Add"}

	type args struct {
		ACLData string
	}
	tests := []struct {
		name    string
		args    args
		want    *ACL
		wantErr bool
	}{
		{
			name: "TestNewACL #1",
			args: args{
				ACLData: ACLDataTest,
			},
			want: &ACL{
				Directory: ACLData,
			},
			wantErr: false,
		},
		{
			name: "TestNewACL #2",
			args: args{
				ACLData: `{"logger1":"/main.Admin/Logging"`,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewACL(tt.args.ACLData)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewACL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewACL() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestACL_CheckAccess(t *testing.T) {
	ACL, _ := NewACL(ACLDataTest)

	type args struct {
		consumer        string
		RequestedMethod string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "#1: ACL.CheckAccess: allowed",
			args: args{
				consumer:        "logger1",
				RequestedMethod: "/main.Admin/Logging",
			},
			want: true,
		},
		{
			name: "#2: ACL.CheckAccess: Not allowed",
			args: args{
				consumer:        "stat1",
				RequestedMethod: "/main.Admin/Logging",
			},
			want: false,
		},
		{
			name: "#3: ACL.CheckAccess: Invalid consumer",
			args: args{
				consumer:        "unknown",
				RequestedMethod: "/main.Admin/Logging",
			},
			want: false,
		},
		{
			name: "#4: ACL.CheckAccess: Invalid method",
			args: args{
				consumer:        "stat1",
				RequestedMethod: "/main.Admin/UnknownMethod",
			},
			want: false,
		},
		// Админ имеет доступ ко всем методам
		// ест на верный парсинг ["/main.Biz/*"]
		{
			name: "#5: ACL.CheckAccess: All methods allowed",
			args: args{
				consumer:        "biz_admin",
				RequestedMethod: "/main.Biz/check",
			},
			want: true,
		},
		{
			name: "#6: ACL.CheckAccess: Empty consumer",
			args: args{
				consumer:        "",
				RequestedMethod: "/main.Biz/check",
			},
			want: false,
		},
		{
			name: "#7: ACL.CheckAccess: Empty methods",
			args: args{
				consumer:        "biz_admin",
				RequestedMethod: "123",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ACL.CheckAccess(tt.args.consumer, tt.args.RequestedMethod); got != tt.want {
				t.Errorf("CheckAccess() = %v, want %v", got, tt.want)
			}
		})
	}
}
