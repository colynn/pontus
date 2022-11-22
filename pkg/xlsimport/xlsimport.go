package xlsimport

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/colynn/pontus/tools"

	execl "github.com/360EntSecGroup-Skylar/excelize"
)

// XlsImport ..
type XlsImport struct{}

// DeptItem ..
type DeptItem struct {
	PrimaryDept string `json:"primary_dept,omitempty"`
	SecondDept  string `json:"second_dept,omitempty"`
	Business    string `json:"business,omitempty"`
	ServcieName string
}

// ServiceData aliyun service item data
type ServiceData struct {
	PrimaryDept string
	SecondDept  string
	Env         string
	ProductLine string
	ServiceName string
	IPAddress   string
	Owner       string
}

// TerminalDevice ..
type TerminalDevice struct {
	AssetNumber       string
	Configuration     string
	InvoiceDiscount   float32
	PretaxGrossAmount float32
	PretaxAmount      float32
	WarrantyPeriod    int //unit month
	AssetLife         int //
	DeliveryDate      *time.Time
	InvoiceNumber     string
	Status            int
	InventoryTime     *time.Time
	ProcurementType   string // 采购类型

	// common
	Manufactory  string
	Type         string // 物理机 physical
	Department   string
	Region       string
	SerialNumber string
	Recipients   RecipientItem
	Description  string

	// physical
	IP            string
	VmwareEnabled bool
	OSName        string
	CPU           int
	Memory        float32
	Disk          int
	SecondDept    string
}

// RecipientItem ..
type RecipientItem struct {
	User    string
	StartAt *time.Time
	EndAt   *time.Time
}

// NewXlsImport ..
func NewXlsImport() *XlsImport {
	return &XlsImport{}
}

// ImportPhysicalDevice ..
func (xls *XlsImport) ImportPhysicalDevice(file string) ([]*TerminalDevice, error) {
	if file == "" {
		return nil, fmt.Errorf("file path is required")
	}
	if !fileExists(file) {
		return nil, fmt.Errorf("file:%s is not exist", file)
	}
	f, err := execl.OpenFile(file)
	if err != nil {
		return nil, err
	}
	rows := f.GetRows("物理机")
	items := make([]*TerminalDevice, 0, len(rows))
	for _, row := range rows {
		if len(row) <= 1 {
			// ingore empty line and sub-title line
			continue
		}
		if row[0] == "品牌" {
			// ingore line hader
			continue
		}
		item := &TerminalDevice{
			Type: "physical",
			// 品牌
			Manufactory: row[0],
			// 序列号
			SerialNumber: row[1],
			IP:           row[2],
			VmwareEnabled: func(item string) bool {
				if item == "是" {
					return true
				}
				return false
			}(row[4]),

			// 操作系统
			OSName: row[5],

			// CPU
			CPU: func(item string) int {
				var err error
				return tools.StrToInt(err, item)
			}(row[6]),

			// 内存
			Memory: func(item string) float32 {
				value, _ := tools.StringToFloat32(item)
				return value
			}(row[7]),

			// Disk
			Disk: func(item string) int {
				var err error
				return tools.StrToInt(err, item)
			}(row[8]),

			// 归属部门
			Department: row[9],
			SecondDept: row[10],

			// 领用人/领用时间
			Recipients: func(user string) RecipientItem {
				if user == "" {
					// 此状态下无领用人
					return RecipientItem{}
				}
				// 领用时间默认为当前时间
				startAt := time.Now()
				return RecipientItem{
					User:    user,
					StartAt: &startAt,
				}
			}(row[11]),
			// 备注
			Description: row[12],
			Region:      row[13],
		}
		items = append(items, item)
	}
	return items, nil
}

// ImportTerminalDevice ..
func (xls *XlsImport) ImportTerminalDevice(file string) ([]*TerminalDevice, error) {
	if file == "" {
		return nil, fmt.Errorf("file path is required")
	}
	if !fileExists(file) {
		return nil, fmt.Errorf("file:%s is not exist", file)
	}
	f, err := execl.OpenFile(file)
	if err != nil {
		return nil, err
	}
	style, _ := f.NewStyle(`{"number_format":14}`)
	f.SetCellStyle("终端设备PC", "L2", "L10000", style)
	rows := f.GetRows("终端设备PC")
	items := make([]*TerminalDevice, 0, len(rows))
	for _, row := range rows {
		if len(row) <= 1 {
			// ingore empty line and sub-title line
			continue
		}
		if row[0] == "编号" {
			// ingore line hader
			continue
		}
		item := &TerminalDevice{
			AssetNumber: row[0],
			Type: func(item string) string {
				switch item {
				case "笔记本":
					return "笔记本电脑"
				case "台式机":
					return "台式电脑"
				default:
					return item
				}
			}(row[1]),
			Manufactory:   row[2],
			Configuration: row[3],

			// 发票号码
			InvoiceNumber: row[4],

			// 发票金额
			InvoiceDiscount: func(item string) float32 {
				value, _ := tools.StringToFloat32(item)
				return value
			}(row[5]),

			// 总金额
			PretaxGrossAmount: func(item string) float32 {
				value, _ := tools.StringToFloat32(item)
				return value
			}(row[6]),

			// 税前金额（税额）
			PretaxAmount: func(item string) float32 {
				value, _ := tools.StringToFloat32(item)
				return value
			}(row[7]),

			// 保修期限
			WarrantyPeriod: func(item string) int {
				if strings.Contains(item, "1年") {
					return 12
				}
				return 0
			}(row[8]),

			// 使用寿命
			// AssetLife: row[9]

			// 领用人/领用时间
			Recipients: func(user, changeDate string) RecipientItem {
				if user == "备用" || user == "已回收" || user == "已退租" {
					// 此状态下无领用人
					return RecipientItem{}
				}
				var startAt time.Time
				if changeDate == "" {
					startAt = time.Now()
				} else {
					startAt, _ = time.Parse("01-02-06", changeDate)
				}

				return RecipientItem{
					User:    user,
					StartAt: &startAt,
				}
			}(row[10], row[11]),

			// 归属部门
			Department:   row[12],
			Region:       row[13],
			SerialNumber: row[14],

			// 盘点时间
			// InventoryTime: row[15],

			// 采购类型
			ProcurementType: row[16],

			// 状态
			// StatusInUse           1 // 使用中
			// StatusInactive      	 2 // 闲置中
			// StatusNeedToScrapped  3 // 需报废
			// StatusAlreadyScrapped 4 // 已报废
			// StatusNeedToSurrender 5 // 需退租
			// StatusSurrendered     6 // 已退租
			Status: func(user string) int {
				if user == "备用" || user == "已回收" {
					return 2
				}
				if user == "已退租" {
					return 6
				}
				return 1
			}(row[10]),

			// 到货日期
			DeliveryDate: func(item string) *time.Time {
				if strings.Contains(item, "年") && strings.Contains(item, "月") {
					splitItem := strings.Split(item, "年")
					year := splitItem[0]
					monthSplit := strings.Split(splitItem[1], "月")
					month := monthSplit[0]
					timeNow, err := time.Parse("2006-01", fmt.Sprintf("%s-%s", year, month))
					if err != nil {
						return nil
					}
					return &timeNow
				}
				return nil
			}(row[18]),
		}
		items = append(items, item)
	}
	return items, nil
}

// ImportDept ..
func (xls *XlsImport) ImportDept(file string) ([]*DeptItem, error) {
	if file == "" {
		return nil, fmt.Errorf("file path is required")
	}
	if !fileExists(file) {
		return nil, fmt.Errorf("file:%s is not exist", file)
	}
	f, err := execl.OpenFile(file)
	if err != nil {
		return nil, err
	}

	// Get all the rows in the Sheet1.
	rows := f.GetRows("Sheet1")

	items := make([]*DeptItem, 0, len(rows))
	for i, row := range rows {
		if i == 0 {
			// ingore line hader
			continue
		}
		if row[3] == "" || row[4] == "" {
			// 产品线/服务名为空时 ignore
			continue
		}
		item := &DeptItem{
			PrimaryDept: strings.TrimSpace(row[0]),
			SecondDept:  strings.TrimSpace(row[1]),
			Business:    strings.TrimSpace(row[3]),
			ServcieName: strings.TrimSpace(row[4]),
		}
		items = append(items, item)
	}
	return items, nil
}

// ParseAliServiceData ..
func (xls *XlsImport) ParseAliServiceData(file string) ([]ServiceData, error) {
	if file == "" {
		return nil, fmt.Errorf("file path is required")
	}
	if !fileExists(file) {
		return nil, fmt.Errorf("file:%s is not exist", file)
	}
	f, err := execl.OpenFile(file)
	if err != nil {
		return nil, err
	}

	// Get all the rows in the Sheet1.
	rows := f.GetRows("Sheet1")
	items := make([]ServiceData, 0, len(rows))
	for i, row := range rows {
		if i == 0 {
			continue
		}
		item := ServiceData{
			PrimaryDept: strings.TrimSpace(row[0]),
			SecondDept:  strings.TrimSpace(row[1]),
			Env:         strings.TrimSpace(row[2]),
			ProductLine: strings.TrimSpace(row[3]),
			ServiceName: strings.TrimSpace(row[4]),
			IPAddress:   strings.TrimSpace(row[5]),
			Owner:       strings.TrimSpace(row[6]),
		}
		items = append(items, item)
	}
	return items, nil
}
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
