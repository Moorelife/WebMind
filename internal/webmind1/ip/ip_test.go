package ip

import (
	"strings"
	"testing"
)

func TestGetLocalIP(t *testing.T) {
	tests := []struct {
		name    string
		want    []string
		wantErr bool
	}{
		{name: "GetLocalIPGoodFlow", want: []string{""}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetLocalIP()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLocalIP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) < 1 {
				t.Error("GetLocalIP() got no valid local IPs")
			}
		})
	}
}

func TestGetOutboundIP(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{name: "GetOutboundIPGoodFlow", want: "", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetOutboundIP()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetOutboundIP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			parts := strings.Split(got, ".")
			if len(parts) != 4 {
				t.Errorf("GetOutboundIP() address has no four parts: %v", got)
				return
			}
		})
	}
}

func TestGetPublicIP(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{name: "GetOutboundIPGoodFlow", want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetPublicIP()
			parts := strings.Split(got, ".")
			if len(parts) != 4 {
				t.Errorf("GetOutboundIP() address has no four parts: %v", got)
				return
			}
		})
	}
}
