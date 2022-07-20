package main

import (
	"testing"
	"os"
)

func TestGet_metadata(t *testing.T) {
	path := "/"
	want := 200
	_, status := get_metadata(path)
	if status != want {
		t.Errorf("get_metadata(%s) returned %d, want %d", path, status, want)
	}
}

func Test_account_info(t *testing.T) {
	want := 200
	_, status := account_info()
	if status != want {
		t.Errorf("account_info() returned %d, want %d", status, want)
	}
}

func Test_get_file(t *testing.T) {
	path := "IMG_20220711_083355.jpg"
	want := 200
	content, status := get_file(path)
	os.WriteFile(path, content, 0644)
	if status != want {
		t.Errorf("get_file(%s) returned %d, want %d", path, status, want)
	}
}


func Test_send_file(t *testing.T) {
	filepath := "main.go"
	want := 200
	status := send_file(filepath)
	if status != want {
		t.Errorf("send_file(%s) returned %d, want %d", filepath, status, want)
	}

}

func Test_delete_file(t *testing.T) {
	filepath := "main.go"
	want := 200
	status := delete_file(filepath)
	if status != want {
		t.Errorf("delete_file(%s) returned %d, want %d", filepath, status, want)
	}
}


func Test_create_dir(t *testing.T) {
	folderpath := "go_created_me"
	want := 403
	status := create_dir(folderpath)
	if status != want {
		t.Errorf("create_dir(%s) returned %d, want %d", folderpath, status, want)
	}
}
