package xlsclient

import (
	"fmt"
	"os"
	"time"

	execl "github.com/360EntSecGroup-Skylar/excelize"
)

// XlsClient ..
type XlsClient struct{}

// DeptItem ..
type DeptItem struct {
	PrimaryDept string `json:"primary_dept,omitempty"`
	SecondDept  string `json:"second_dept,omitempty"`
	Business    string `json:"business,omitempty"`
}

// TerminalDevice ..
type TerminalDevice struct {
	AssetNumber     string
	Type            string
	Manufactory     string
	Configuration   string
	InvoiceDiscount float32
	WarrantyPeriod  int //unit month
	DeliveryDate    *time.Time
	InvoiceNumber   string
	Region          string
	SerialNumber    string
	Status          int
	Recipients      []RecipientItem
}

// RecipientItem ..
type RecipientItem struct {
	User    string
	StartAt *time.Time
	EndAt   *time.Time
}

// SurrenderedTime   time.Time `json:"SurrenderedTime"`
// InvoiceTime       time.Time `json:"InvoiceTime"`                             // 发票日期，入账日期
// InvoiceNumber     string    `gorm:"type:vachar(128)" json:"InvoiceNumber"`   // 发票号码
// PretaxAmount      float32   `gorm:"type:float(32)" json:"PretaxAmount"`      // 税前金额（税额）
// PretaxGrossAmount float32   `gorm:"type:float(32)" json:"PretaxGrossAmount"` // 总金额
// InvoiceDiscount   float32   `gorm:"type:float(32)" json:"InvoiceAmount"`     // 发票金额
// AssetLife         int       `gorm:"type:int" json:"AssetLife"`               // 使用寿命
// WarrantyPeriod    int       `gorm:"type:int" json:"WarrantyPeriod"`          // 保修期限
// ResidualRatio     float32   `gorm:"type(float(32)" json:"ResidualRatio"`

// InstanceID  string `gorm:"primary_key;type:varchar(128)" json:"InstanceID"` // eg. i-xxx
// AssetNumber string `gorm:"type:vachar(128)" json:"AssetNumber"`             // 原资产编号
// // InstanceName string `gorm:"type:varchar(128)" json:"InstanceName"`
// Manufactory string `gorm:"type:varchar(128)" json:"Manufactory"` // 品牌 eg. aliyun、苹果

// Type          string `gorm:"type:varchar(64)" json:"Type"` // eg. 笔记本电脑/台式电脑/显示器
// Configuration string `gorm:"type:varchar(1024)" json:"Configuration"`
// SerialNumber  string `gorm:"type:varchar(128)" json:"SerialNumber"` // eg. 8egtg6016-15eb-4da7-a177-6f423dfa339a
// Region        string `gorm:"type:varchar(64)" json:"Region"`        // eg. 地区
// RecipientID   int    `gorm:"type:int" json:"RecipientID"`           // 领取人ID

// NewXlsClient ..
func NewXlsClient() *XlsClient {
	return &XlsClient{}
}

// RenderPrimayAndSecondDept ..
func (xls *XlsClient) RenderPrimayAndSecondDept(assetfile, deptfile string) ([]*TerminalDevice, error) {
	if assetfile == "" || deptfile == "" {
		return nil, fmt.Errorf("file path is required")
	}
	if !fileExists(assetfile) {
		return nil, fmt.Errorf("file:%s is not exist", assetfile)
	}

	if !fileExists(deptfile) {
		return nil, fmt.Errorf("file:%s is not exist", deptfile)
	}

	f1, err := execl.OpenFile(assetfile)
	if err != nil {
		return nil, err
	}

	deptf, err := execl.OpenFile(deptfile)
	if err != nil {
		return nil, err
	}

	roster := deptf.GetRows("Sheet1")

	rows := f1.GetRows("上海")

	// rows = append(rows, f.GetRows("北京")...)
	for index, row := range rows {
		if len(row) <= 1 {
			// ingore empty line and sub-title line
			continue
		}
		if row[1] == "资产编号" {
			// ingore line hader
			continue
		}
		if row[9] == "" || row[9] == "备用" || row[9] == "已回收" {
			// ingore
			continue
		}
		user := row[9]
		fmt.Printf("index: %d row-9: %s\n", index, row[9])
		for _, item := range roster {
			if item[3] == user {
				fmt.Printf("roster: %s, %s, %s, %s\n", item[0], item[1], item[2], item[3])
			}
		}
	}
	return nil, nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
