package common

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestGetJSONRPCRequestBody(t *testing.T) {
	type args struct {
		jsonrpc string
		method  string
		params  []string
		id      int64
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "正向测试",
			args: args{
				jsonrpc: "2.0",
				method:  "add",
				params:  []string{"1", "2", "3"},
				id:      1,
			},
			want:    []byte(`{"jsonrpc":"2.0","method":"add","params":["1","2","3"],"id":1}`),
			wantErr: false,
		},
		{
			name: "空参数测试",
			args: args{
				jsonrpc: "2.0",
				method:  "subtract",
				params:  []string{},
				id:      2,
			},
			want:    []byte(`{"jsonrpc":"2.0","method":"subtract","params":[],"id":2}`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getJSONRPCRequestBody(tt.args.jsonrpc, tt.args.method, tt.args.params, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("getJSONRPCRequestBody() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getJSONRPCRequestBody() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSendHTTPRequest(t *testing.T) {
	type args struct {
		url     string
		headers map[string]string
		body    io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "正向测试",
			args: args{
				url:     "http://localhost:8545",
				headers: map[string]string{"Content-Type": "application/json"},
				body:    bytes.NewBuffer([]byte(`{"name":"test","phone":13682937384,"password":"123456"}`)),
			},
			want:    []byte(`{"code":200,"message":"success}`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建mock Http服务
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"code":200,"message":"success}`))
			}))
			defer mockServer.Close()

			got, err := sendHTTPRequest(mockServer.URL, tt.args.headers, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("sendHTTPRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sendHTTPRequest() got = %v, want %v", got, tt.want)
			}
		})
	}
}
