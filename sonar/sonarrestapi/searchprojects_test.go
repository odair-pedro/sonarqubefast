package sonarrestapi

import (
	"encoding/json"
	"errors"
	"sonarci/net"
	"sonarci/sonar"
	"testing"
)

func Test_restApi_SearchProjects(t *testing.T) {
	mockOk := &mockConnection{doGet: func(route string) (<-chan []byte, <-chan error) {
		wrapper := &searchProjectsWrapper{Components: []searchProject{{"id", "org", "key", "name", "visibility"}}}
		buff, _ := json.Marshal(wrapper)

		chOk := make(chan []byte, 1)
		chOk <- buff

		chEr := make(chan error, 1)
		chEr <- nil
		return chOk, chEr
	}}
	mockComponentsNil := &mockConnection{doGet: func(route string) (<-chan []byte, <-chan error) {
		wrapper := &searchProjectsWrapper{}
		buff, _ := json.Marshal(wrapper)

		chOk := make(chan []byte, 1)
		chOk <- buff

		chEr := make(chan error, 1)
		chEr <- nil
		return chOk, chEr
	}}
	mockComponentsEmpty := &mockConnection{doGet: func(route string) (<-chan []byte, <-chan error) {
		wrapper := &searchProjectsWrapper{Components: []searchProject{}}
		buff, _ := json.Marshal(wrapper)

		chOk := make(chan []byte, 1)
		chOk <- buff

		chEr := make(chan error, 1)
		chEr <- nil
		return chOk, chEr
	}}
	mockInvalidJson := &mockConnection{doGet: func(route string) (<-chan []byte, <-chan error) {
		chOk := make(chan []byte, 1)
		chOk <- []byte{}

		chEr := make(chan error, 1)
		chEr <- nil
		return chOk, chEr
	}}
	mockError := &mockConnection{doGet: func(route string) (<-chan []byte, <-chan error) {
		chError := make(chan error, 1)
		chError <- errors.New("failure")
		return nil, chError
	}}

	type fields struct {
		Connection net.Connection
		Server     string
	}
	type args struct {
		projects string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    sonar.Project
		wantErr bool
	}{
		{"ok", fields{Connection: mockOk, Server: "http://server"}, args{projects: "project"}, sonar.Project{Id: "id", Organization: "org", Key: "key", Name: "name", Visibility: "visibility"}, false},
		{"components-nil", fields{Connection: mockComponentsNil, Server: "http://server"}, args{projects: "project"}, sonar.Project{}, false},
		{"components-empty", fields{Connection: mockComponentsEmpty, Server: "http://server"}, args{projects: "project"}, sonar.Project{}, false},
		{"invalid-json", fields{Connection: mockInvalidJson, Server: "http://server"}, args{projects: "project"}, sonar.Project{}, true},
		{"error", fields{Connection: mockError, Server: "http://server"}, args{projects: "project"}, sonar.Project{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			restApi := &restApi{
				Connection: tt.fields.Connection,
				Server:     tt.fields.Server,
			}
			chGot, err := restApi.SearchProjects(tt.args.projects)
			if (err != nil) != tt.wantErr {
				t.Errorf("SearchProjects() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if chGot == nil {
				t.Errorf("SearchProjects() got nil, want %v", tt.wantErr)
				return
			}
			if !tt.wantErr {
				got := <-chGot
				if got != tt.want {
					t.Errorf("SearchProjects() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
