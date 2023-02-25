package heisiwu

import (
	"reflect"
	"testing"
)

func TestGetTextLink(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"case1", args{"http://hs.heisiwu.com/heisi/page/17"}, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetTextLink(tt.args.url); !reflect.DeepEqual(len(got), tt.want) {
				t.Errorf("ReadHtmlAndGetHref() = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestGetImageLink(t *testing.T) {
	type args struct {
		url      string
		selector string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"case1", args{"http://hs.heisiwu.com/heisi/33230.html", "img[loading=\"lazy\"]"}, 75},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetImageLink(tt.args.url, tt.args.selector); !reflect.DeepEqual(len(got), tt.want) {
				t.Errorf("GetImageLink() = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestDownloadImage(t *testing.T) {
	type args struct {
		url        string
		referer    string
		folderPath string
	}
	tests := []struct {
		name string
		args args
	}{
		{"case1", args{"http://hs.heisiwu.com/wp-content/uploads/2023/01/5eba022e40fd294.jpg", "http://hs.heisiwu.com/baisi/62019.html", "heisiwu"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			DownloadImage(tt.args.url, tt.args.referer, tt.args.folderPath)
		})
	}
}
