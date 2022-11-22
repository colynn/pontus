package xlsclient_test

import (
	"fmt"
	"testing"

	"github.com/colynn/pontus/pkg/xlsclient"
)

func TestTernimalDeveice(t *testing.T) {
	client := xlsclient.NewXlsClient()
	items, err := client.RenderPrimayAndSecondDept("/Users/colynn/Documents/Code/pontus/pontus/assets/电脑及显示器采购明细表2019.xlsx", "/Users/colynn/Documents/Code/pontus/pontus/pkg/xlsclient/assets/花名册.xlsx")
	if err != nil {
		fmt.Printf("improt devices error: %s", err.Error())
		return
	}
	fmt.Printf("ternimal devices parse line: %v\n", len(items))
}
