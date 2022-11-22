package xlsimport_test

import (
	"fmt"
	"testing"

	"github.com/colynn/pontus/pkg/xlsimport"

	"github.com/stretchr/testify/assert"
)

func TestImportExecl(t *testing.T) {
	client := xlsimport.NewXlsImport()
	deptInfo, err := client.ImportDept("/Users/colynn/Documents/Code/pontus/backend/assets/aliServiceData.xlsx")
	assert.Equal(t, nil, err)
	if err != nil {
		fmt.Printf("improt dept error: %s", err.Error())
		return
	}
	fmt.Printf("dept parse line: %v\n", len(deptInfo))
}

func TestTernimalDeveice(t *testing.T) {
	client := xlsimport.NewXlsImport()
	items, err := client.ImportTerminalDevice("/Users/colynn/Documents/Code/pontus/backend/assets/pc.xlsx")
	if err != nil {
		fmt.Printf("improt devices error: %s", err.Error())
		return
	}
	fmt.Printf("ternimal devices parse line: %v\n", len(items))
	for _, item := range items {
		fmt.Printf("item.started_at: %v\n", item.Recipients.StartAt)
	}
}

func TestPhysicalData(t *testing.T) {
	client := xlsimport.NewXlsImport()
	items, err := client.ImportPhysicalDevice("/Users/colynn/Documents/Code/pontus/backend/assets/ph.xlsx")
	assert.Equal(t, nil, err)
	if err != nil {
		fmt.Printf("improt devices error: %s", err.Error())
		return
	}
	fmt.Printf("ternimal devices parse line: %v\n", len(items))
}
