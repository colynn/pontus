package audit

import (
	"fmt"
	"testing"
)

func TestGetContentStr(t *testing.T) {
	contentStr := "hellworld"
	strs := getContentStr(contentStr)
	if strs == "hellworld" {
		fmt.Printf("str type success\n")
	}
}

func TestGetContentStruct(t *testing.T) {
	contentStruct := Content{
		{
			Field:       "Test",
			OriginValue: "1",
			NewValue:    "3",
		},
	}
	str02 := getContentStr(contentStruct)
	if len(str02) > 0 {
		fmt.Printf("stuct type success, origin struct: %+v\nstrings: %s\n", contentStruct, str02)
	} else {
		fmt.Printf("stuct type failed, strings: %s\n", str02)
	}
}
