package heisiwu

import (
	"fmt"
	"testing"
)

func Test_getSetu(t *testing.T) {
	title, imageUrls := getSetu("heisi", 11)
	fmt.Printf("title: %s", title)
	println()
	fmt.Printf("imageUrls: %v", imageUrls)
}
