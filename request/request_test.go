package request

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func Test_Request(t *testing.T) {
	const token = "123456790"
	type args struct {
		newPR PR
		token string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "Not Authorized test",
			args: args{
				newPR: PR{
					Title: "title",
					Base:  "base",
					Head:  "head",
				},
				token: "test",
			},
			want:    http.StatusUnauthorized,
			wantErr: false,
		},
		{
			name: "No Title test",
			args: args{
				newPR: PR{
					Title: "",
					Base:  "base",
					Head:  "head",
				},
				token: token,
			},
			want:    http.StatusInternalServerError,
			wantErr: false,
		},
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "token "+token {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		decoder := json.NewDecoder(r.Body)
		var pr PR
		err := decoder.Decode(&pr)
		if err != nil {
			panic(err)
		}
		if pr.Title == "" || pr.Head == "" || pr.Base == "" {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}))
	defer ts.Close()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Request(tt.args.newPR, ts.URL, tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("Request() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.StatusCode, tt.want) {
				t.Errorf("Request() = %v, want %v", got, tt.want)
			}
		})
	}
}
