package cmdb

import (
	"fmt"

	"github.com/colynn/pontus/internal/account/user"
	"github.com/colynn/pontus/pkg/xlsimport"
	"github.com/colynn/pontus/tools/uuid"

	log "unknwon.dev/clog/v2"
)

// CreateTagReq ..
type CreateTagReq struct {
	TagKey      string `json:"tag_key,omitempty"`
	TagValue    string `json:"tag_value,omitempty"`
	Description string `json:"description,omitempty"`
}

// PaginationReq ..
type PaginationReq struct {
	PageIndex int `json:"page_index"`
	PageSize  int `json:"page_size"`
}

// AssetReq ..
type AssetReq struct {
	PaginationReq
	FlDepartment string
	TlDepartment string
	Business     string
	Environment  string
	Type         string
	AssetName    string
	AssetIP      string
	Status       string
	StartAt      string
	EndAt        string
	OSType       string
}

// TangibleAssetReq ..
type TangibleAssetReq struct {
	PaginationReq

	Business     string
	InstanceID   string
	SerialNumber string
	AssetName    string // 资产名称

	// Physical
	PrivateIP string
	OSName    string

	// is used
	PCList        bool
	Type          string
	PrimaryDeptID int
	SecondDeptID  int
	UserName      string // 领用人
	Status        int
	StartAt       string
	EndAt         string
}

// TangibleAssetUpdate ..
type TangibleAssetUpdate struct {
	ID          int    `json:"id"`          // 编码
	AssetNumber string `json:"AssetNumber"` // 原资产编号
	InstanceID  string `json:"InstanceID"`  // eg. i-xxx
	// InstanceName string `gorm:"type:varchar(128)" json:"InstanceName"`
	Manufactory     string `json:"Manufactory"`     // 品牌 eg. aliyun、苹果
	Type            string `json:"Type"`            // eg. 笔记本电脑/台式电脑/显示器/服务器
	Configuration   string `json:"Configuration"`   // 配置
	ProcurementType string `json:"ProcurementType"` // 采购类型 自采/租赁

	// 类型为服务器,有效
	VmwareEnabled bool    `json:"VmwareEnabled" gorm:"type:char(1);"` // vmware 是否启用
	CPU           int     `json:"CPU"`
	Memory        float32 `json:"Memory"`
	Disk          int     `json:"Disk"`
	PrivateIP     string  `json:"PrivateIp"`
	PublicIP      string  `json:"PublicIP"`

	InvoiceTime       string  `json:"InvoiceTime"`       // 发票日期，入账日期
	DeliveryDate      string  `json:"DeliveryDate"`      // 到货日期
	InvoiceNumber     string  `json:"InvoiceNumber"`     // 发票号码
	PretaxAmount      float32 `json:"PretaxAmount"`      // 税前金额（税额）
	PretaxGrossAmount float32 `json:"PretaxGrossAmount"` // 总金额
	InvoiceDiscount   float32 `json:"InvoiceAmount"`     // 发票金额
	AssetLife         int     `json:"AssetLife"`         // 使用寿命
	WarrantyPeriod    int     `json:"WarrantyPeriod"`    // 保修期限
	ResidualRatio     float32 `json:"ResidualRatio"`

	InventoryTime string `json:"InventoryTime"` // 盘点时间
	SerialNumber  string `json:"SerialNumber"`  // eg. 8e8a6016-15eb-4da7-a177-6f423dfa339a
	Region        string `json:"Region"`        // eg. 地区
	RecipientID   int    `json:"RecipientID"`   // 领取人ID

	Status          int     `json:"Status"`        // 资产状态
	ScrappedTime    string  `json:"scrappedTime"`  // 报废时间
	ScrappedPrice   float32 `json:"ScrappedPrice"` // 报废价格
	SurrenderedTime string  // 退租时间
	Description     string
	Recipient       string
}

// CategoryItem ..
type CategoryItem struct {
	Name    string `json:"name"`
	Value   int    `json:"value"`
	Running int
	Stopped int
}

// AssetStatisticsResp ..
type AssetStatisticsResp struct {
	AssetCategories      []*CategoryItem `json:"asset_categories"`
	EnvCategories        []*CategoryItem `json:"env_categories"`
	DepartmentCategories []*CategoryItem `json:"department_categories"`
	ServerCategories     []*CategoryItem `json:"server_categories"`
	ECSRegionCategories  []*CategoryItem `json:"ecsregion_categories"`
	SLBCategories        []*CategoryItem `json:"slb_categories"`
	OSSCategories        []*CategoryItem `json:"oss_categories"`
	NASCategories        []*CategoryItem `json:"nas_categories"`
}

// InstanceDisplayName ..
type InstanceDisplayName struct {
	InstanceID  string
	DisplayName string
	ZoneID      string
}

// NonStandardInstance ..
type NonStandardInstance struct {
	IP           string
	InstanceID   string
	InstanceName string
}

// Service ..
type Service struct {
	CloudAsset    *CloudAsset
	TangibleAsset *TangibleAsset
}

// NewService ..
func NewService() *Service {
	return &Service{
		CloudAsset:    &CloudAsset{},
		TangibleAsset: &TangibleAsset{},
	}
}

// InsertOrUpdatePhysicalDevice ..
func (c *Service) InsertOrUpdatePhysicalDevice(filePath string) error {
	xls := xlsimport.NewXlsImport()
	items, err := xls.ImportPhysicalDevice(filePath)
	if err != nil {
		log.Error("when inser update physical device, parse xls error: %s", err.Error())
		return err
	}
	allUsers, err := getAllUsers()
	if err != nil {
		log.Error("get all dept/users error: %s", err.Error())
		return fmt.Errorf("get all dept/users error: %s", err.Error())
	}
	syncPool := make(chan bool, 50)
	for _, item := range items {
		go func(item *xlsimport.TerminalDevice) {
			syncPool <- true
			err := c.InsertOrUpdateTangibleDeviceOneItem(item, allUsers)
			if err != nil {
				log.Error("inster or update item error: %s", err.Error())
			}
			<-syncPool
		}(item)
	}
	return nil
}

// InsertOrUpdateTerminalDevice ..
func (c *Service) InsertOrUpdateTerminalDevice(filePath string) error {
	xls := xlsimport.NewXlsImport()
	items, err := xls.ImportTerminalDevice(filePath)
	if err != nil {
		log.Error("when inser update terminal device, parse xls error: %s", err.Error())
		return err
	}
	allUsers, err := getAllUsers()
	if err != nil {
		log.Error("get all dept/users error: %s", err.Error())
		return fmt.Errorf("get all dept/users error: %s", err.Error())
	}
	syncPool := make(chan bool, 50)
	for _, item := range items {
		go func(item *xlsimport.TerminalDevice) {
			syncPool <- true
			err := c.InsertOrUpdateTangibleDeviceOneItem(item, allUsers)
			if err != nil {
				log.Error("inster or update item error: %s", err.Error())
			}
			<-syncPool
		}(item)
	}
	return nil
}

// InsertOrUpdateTangibleDeviceOneItem ..
// pc, physical
func (c *Service) InsertOrUpdateTangibleDeviceOneItem(item *xlsimport.TerminalDevice, allUsers []user.SysUser) error {
	instanceID := uuid.GetUUID()
	tangibleInstance := &TangibleAsset{
		// Common
		Manufactory:  item.Manufactory,
		Type:         item.Type,
		SerialNumber: item.SerialNumber,
		Region:       item.Region,
		Description:  item.Description,

		// physical
		PrivateIP:     item.IP,
		VmwareEnabled: item.VmwareEnabled,
		OSName:        item.OSName,
		CPU:           item.CPU,
		Memory:        item.Memory,
		Disk:          item.Disk,

		// pc
		AssetNumber:     item.AssetNumber,
		Configuration:   item.Configuration,
		WarrantyPeriod:  item.WarrantyPeriod,
		DeliveryDate:    item.DeliveryDate,
		InvoiceNumber:   item.InvoiceNumber,
		InvoiceTime:     item.InventoryTime,
		InventoryTime:   item.InventoryTime,
		ProcurementType: item.ProcurementType,
		Status:          item.Status,
		// pc - 金额
		PretaxAmount:      item.PretaxAmount,
		PretaxGrossAmount: item.PretaxGrossAmount,
		InvoiceDiscount:   item.InvoiceDiscount,
	}
	id, err := tangibleInstance.InsertOrUpdate(instanceID)
	if err != nil {
		log.Error("item: %s insert or update error: %s", instanceID, err.Error())
	}
	tangibleModel, err := c.TangibleAsset.GetTangibleAssetByID(id)
	if err != nil {
		log.Error("get asset item error: %s", err.Error())
	}
	instanceID = tangibleModel.InstanceID
	if instanceID == "" {
		instanceID = uuid.GetUUID()
		tangibleModel.InstanceID = instanceID
		_, err := tangibleModel.Update(tangibleModel.ID)
		if err != nil {
			log.Error("update item error: %s", err.Error())
		}
	}
	// update recipient
	userID := getUserIDByName(allUsers, item.Recipients.User)
	if userID != 0 {
		err = generateRecipientItem(item.Recipients, instanceID, userID)
		if err != nil {
			log.Warn("get user id error: %s", err.Error())
		} else {
			tangibleModel.RecipientID = userID
		}
	}

	_, err = tangibleModel.Update(tangibleModel.ID)
	if err != nil {
		log.Trace("update item error: %s", err.Error())
	}

	return nil
}

// generateRecipientItem ..
func generateRecipientItem(userItem xlsimport.RecipientItem, instanceID string, userID int) error {
	// 用户有效时
	assetUserInstance := &AssetUser{
		UserID:          userID,
		StartedAt:       userItem.StartAt,
		EndedAt:         userItem.EndAt,
		AssetInstanceID: instanceID,
	}
	_, err := assetUserInstance.InsertOrUpdate()
	if err != nil {
		log.Error("insert asseet user item error: %s", err.Error())
		return err
	}
	return nil
}

func getAllUsers() (allUsers []user.SysUser, err error) {
	allUsers, err = user.NewService().UserDao.GetList()
	if err != nil {
		return
	}
	return
}
