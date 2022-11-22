package cmdb

import (
	"errors"
	"time"

	"github.com/colynn/pontus/internal/account/user"
	"github.com/colynn/pontus/internal/db"
	"github.com/colynn/pontus/tools"

	log "unknwon.dev/clog/v2"
)

/* ----- model defined start ------   */

// TangibleAsset status
const (
	StatusInUse           = iota + 1
	StatusInactive        // 闲置中
	StatusNeedToScrapped  // 需报废
	StatusAlreadyScrapped // 已报废
	StatusNeedToSurrender // 需退租
	StatusSurrendered     // 已退租
)

// TangibleAsset ..
type TangibleAsset struct {
	db.BaseModel
	ID int `gorm:"primary_key;AUTO_INCREMENT"  json:"id"` // 编码

	AssetNumber string `gorm:"type:varchar(128)" json:"AssetNumber"`            // 原资产编号
	InstanceID  string `gorm:"primary_key;type:varchar(128)" json:"InstanceID"` // eg. i-xxx
	// InstanceName string `gorm:"type:varchar(128)" json:"InstanceName"`
	Manufactory     string `gorm:"type:varchar(128)" json:"Manufactory"`    // 品牌 eg. aliyun、苹果
	Type            string `gorm:"type:varchar(64)" json:"Type"`            // eg. 笔记本电脑/台式电脑/显示器/服务器
	Configuration   string `gorm:"type:varchar(1024)" json:"Configuration"` // 配置
	ProcurementType string `gorm:"type:varchar(64)" json:"ProcurementType"` // 采购类型 自采/租赁

	// 类型为服务器,有效
	VmwareEnabled bool    `json:"VmwareEnabled" gorm:"type:char(1);"` // vmware 是否启用
	CPU           int     `gorm:"type:int(11)" json:"CPU"`
	Memory        float32 `gorm:"type:float(32)" json:"Memory"`
	Disk          int     `gorm:"type:int(11)" json:"Disk"`
	OSName        string  `gorm:"type:varchar(128)" json:"OSName"`
	PrivateIP     string  `gorm:"type:varchar(128)" json:"PrivateIP"`
	PublicIP      string  `gorm:"type:varchar(128)" json:"PublicIP"`
	Description   string  `gorm:"type:varchar(1024)" json:"Description"`

	InvoiceTime       *time.Time `json:"InvoiceTime"`                             // 发票日期，入账日期
	DeliveryDate      *time.Time `json:"DeliveryDate"`                            // 到货日期
	InvoiceNumber     string     `json:"InvoiceNumber"`                           // 发票号码
	PretaxAmount      float32    `gorm:"type:float(32)" json:"PretaxAmount"`      // 税前金额（税额）
	PretaxGrossAmount float32    `gorm:"type:float(32)" json:"PretaxGrossAmount"` // 总金额
	InvoiceDiscount   float32    `gorm:"type:float(32)" json:"InvoiceAmount"`     // 发票金额
	AssetLife         int        `gorm:"type:int" json:"AssetLife"`               // 使用寿命
	WarrantyPeriod    int        `gorm:"type:int" json:"WarrantyPeriod"`          // 保修期限
	ResidualRatio     float32    `gorm:"type:float(32)" json:"ResidualRatio"`

	InventoryTime *time.Time `json:"InventoryTime"`                         // 盘点时间
	SerialNumber  string     `gorm:"type:varchar(128)" json:"SerialNumber"` // eg. 8e8a6016-15eb-4da7-a177-6f423dfa339a
	Region        string     `gorm:"type:varchar(64)" json:"Region"`        // eg. 地区
	RecipientID   int        `gorm:"type:int" json:"RecipientID"`           // 领取人ID

	Status          int        `gorm:"type:int" json:"Status"`              // 资产状态
	ScrappedTime    *time.Time `json:"scrappedTime"`                        // 报废时间
	ScrappedPrice   float32    `gorm:"type:float(32)" json:"ScrappedPrice"` // 报废价格
	SurrenderedTime *time.Time `json:"SurrenderedTime"`                     // 退租时间

	Recipient  string `gorm:"-"`
	Department string `gorm:"-"`
	SecondDept string `gorm:"-"`
}

// TableName ..
func (TangibleAsset) TableName() string {
	return "cmdb_tangible_asset"
}

// AssetUser ..
type AssetUser struct {
	db.BaseModel
	ID              int        `gorm:"primary_key;AUTO_INCREMENT"  json:"id"` // 编码
	StartedAt       *time.Time `json:"StartedAt"`                             // 开始时间
	EndedAt         *time.Time `json:"EndedAt"`                               // 结束时间
	UserID          int        `gorm:"type:int" json:"userID"`
	AssetInstanceID string     `gorm:"type:varchar(128)" json:"AssetInstanceID"`
}

// TableName ..
func (AssetUser) TableName() string {
	return "cmdb_asset_user"
}

/* ----- model defined end ------   */

func getUserIDByName(allUser []user.SysUser, userName string) int {
	if userName == "" || userName == "备用" || userName == "已回收" {
		return 0
	}
	for _, user := range allUser {
		if userName == user.RealName {
			return user.ID
		}
	}
	return 0
}

// GetPage ..
func (r *TangibleAsset) GetPage(req *TangibleAssetReq) ([]*TangibleAsset, int, error) {
	var doc []*TangibleAsset

	table := db.Eloquent.Select("*").Table(r.TableName()).Where("deleted_at is NULL")

	if req.InstanceID != "" {
		table = table.Where("instance_id = ?", req.InstanceID)
	}
	if req.Status != 0 {
		table = table.Where("status = ? ", req.Status)
	}
	if req.SerialNumber != "" {
		table = table.Where("serial_number = ?", req.SerialNumber)
	}

	if req.PrimaryDeptID != 0 {
		table = table.Where("primary_dept = ?", req.PrimaryDeptID)
	}

	if req.SecondDeptID != 0 {
		table = table.Where("second_dept_id = ?", req.SecondDeptID)
	}

	if req.PrivateIP != "" {
		table = table.Where("private_ip like ?", "%"+req.PrivateIP+"%")
	}
	if req.OSName != "" {
		table = table.Where("os_name like ?", "%"+req.OSName+"%")
	}
	if req.PCList {
		table = table.Where("type != ?", "physical")
	}

	if req.Type != "" {
		table = table.Where("type = ?", req.Type)
	}

	if req.UserName != "" {
		userIDs, err := getCurrentRecipientID(req.UserName)
		if err != nil {
			log.Error("when get recipient ids error: %s", err.Error())
		}
		log.Trace("recipient ids: %+v", userIDs)
		ids := make([]int, 0, len(userIDs))
		for _, item := range userIDs {
			ids = append(ids, item.UserID)
		}
		table = table.Where("recipient_id in (?)", ids)

	}

	if req.StartAt != "" && req.EndAt != "" {
		startAt, _ := time.Parse("2006-01-02", req.StartAt)
		endAt, _ := time.Parse("2006-01-02", req.EndAt)
		table = table.Where("delivery_date >= ? AND delivery_date <= ?", startAt, endAt)
	}

	var count int

	if err := table.Offset((req.PageIndex - 1) * req.PageSize).Limit(req.PageSize).Find(&doc).Error; err != nil {
		return nil, 0, err
	}
	table.Count(&count)
	// assign username
	for _, item := range doc {
		item.Recipient = GetUserRealName(item.RecipientID)
	}
	return doc, count, nil
}

// Get ..
func (r *TangibleAsset) Get() (tangibleAsset TangibleAsset, err error) {
	table := db.Eloquent.Table(r.TableName()).Where("deleted_at is NULL")
	if r.InstanceID != "" {
		table = table.Where("instance_id = ?", r.InstanceID)
	}
	if r.InvoiceNumber != "" {
		table = table.Where("invoice_number = ? ", r.InvoiceNumber)
	}
	if r.InvoiceNumber != "" {
		table = table.Where("invoice_number = ?", r.InvoiceNumber)
	}
	if r.AssetNumber != "" {
		table = table.Where("asset_number = ?", r.AssetNumber)
	}
	if r.Manufactory != "" {
		table = table.Where("manufactory = ?", r.Manufactory)
	}
	if r.SerialNumber != "" {
		table = table.Where("serial_number = ?", r.SerialNumber)
	}
	if r.Type != "" {
		table = table.Where("type = ?", r.Type)
	}
	if err = table.First(&tangibleAsset).Error; err != nil {
		return
	}

	tangibleAsset.Recipient = GetUserRealName(tangibleAsset.RecipientID)
	return
}

// GetTangibleAssetByID  ..
func (r *TangibleAsset) GetTangibleAssetByID(id int) (*TangibleAsset, error) {
	item := TangibleAsset{}
	table := db.Eloquent.Table(r.TableName())
	if id == 0 {
		return nil, errors.New("get tangible id is required")
	}
	if err := table.Where("id = ?", id).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// Insert ..
func (r *TangibleAsset) Insert() (id int, err error) {
	existRole := TangibleAsset{}
	existRole.InstanceID = r.InstanceID
	_, existErr := existRole.Get()
	if existErr == nil {
		err = errors.New("资产唯一标识不允许重复，请修改后重试")
		return
	}
	result := db.Eloquent.Table("cmdb_tangible_asset").Create(&r)
	if result.Error != nil {
		err = result.Error
		return
	}
	id = r.ID
	return
}

// InsertOrUpdate ..
func (r *TangibleAsset) InsertOrUpdate(instanceID string) (id int, err error) {
	//
	existedModel, existErr := r.Get()
	if existErr != nil {
		r.InstanceID = instanceID
		if err = db.Eloquent.Table(r.TableName()).Create(&r).Error; err != nil {
			id = r.ID
			return
		}
	}
	if err = db.Eloquent.Table(r.TableName()).Model(&existedModel).Updates(&r).Error; err != nil {
		return
	}
	id = existedModel.ID
	return
}

// UpdateItem ..修改 based on instanceID
func (r *TangibleAsset) UpdateItem(instanceID string) (update TangibleAsset, err error) {
	if err = db.Eloquent.Table(r.TableName()).Where("instance_id = ?", instanceID).First(&update).Error; err != nil {
		return
	}

	// 闲置状态时，会清空联系人
	if r.Status == 2 {
		r.RecipientID = 0
	}

	if update.RecipientID != r.RecipientID {
		err := r.UpdateAndInsertAssetUser(update.RecipientID, r.RecipientID, update.InstanceID)
		if err != nil {
			log.Warn("update recipient error: %s, retry..", err.Error())
			err = r.UpdateAndInsertAssetUser(update.RecipientID, r.RecipientID, update.InstanceID)
			log.Error("update recipient error: %s, need repair data manually", err.Error())
		}
	}

	//参数1:是要修改的数据
	//参数2:是修改的数据
	if err = db.Eloquent.Table(r.TableName()).Model(&update).Updates(&r).Error; err != nil {
		return
	}

	if r.RecipientID == 0 {
		err = db.Eloquent.Table(r.TableName()).Save(&r).Error
	}
	update.Recipient = GetUserRealName(r.RecipientID)
	return
}

// UpdateAndInsertAssetUser ..
func (r *TangibleAsset) UpdateAndInsertAssetUser(originID, userID int, instanceID string) error {
	//
	assetUser := AssetUser{
		UserID:          originID,
		AssetInstanceID: instanceID,
	}
	if err := assetUser.UpdateEndAt(); err != nil {
		log.Error("when update end at: userid: %d error: %s", originID, err.Error())
	}
	// userID 为空时，不再创建 assetUser item
	if userID == 0 {
		return nil
	}
	userModel := user.SysUser{}
	userItem, err := userModel.Get(userID)
	if err != nil {
		return err
	}
	timeValue := tools.GetCurrentTime("")
	newAssetUserItem := AssetUser{
		UserID:    userItem.ID,
		StartedAt: &timeValue,
	}
	id, err := newAssetUserItem.InsertOrUpdate()
	log.Trace("add asset user(%d), item: %d success", id)
	return nil
}

// DeleteItem ..
func (r *TangibleAsset) DeleteItem() (origin TangibleAsset, err error) {
	if err = db.Eloquent.Table(r.TableName()).Where("instance_id = ?", r.InstanceID).First(&origin).Error; err != nil {
		return
	}

	deletedTime := time.Now()
	r.DeletedAt = &deletedTime
	if err = db.Eloquent.Table(r.TableName()).Model(&origin).Updates(&r).Error; err != nil {
		return
	}
	return

}

// Update ..修改
func (r *TangibleAsset) Update(id int) (update TangibleAsset, err error) {
	if err = db.Eloquent.Table(r.TableName()).First(&update, id).Error; err != nil {
		return
	}
	//参数1:是要修改的数据
	//参数2:是修改的数据
	if err = db.Eloquent.Table(r.TableName()).Model(&update).Updates(&r).Error; err != nil {
		return
	}
	return
}

// Get ..
func (u *AssetUser) Get() (assetUser AssetUser, err error) {
	table := db.Eloquent.Table(u.TableName()).Where("deleted_at is NULL")

	if u.AssetInstanceID != "" {
		table = table.Where("asset_instance_id = ? ", u.AssetInstanceID)
	}
	if u.UserID != 0 {
		table = table.Where("user_id = ?", u.UserID)
	}
	if u.StartedAt != nil {
		table = table.Where("started_at = ?", u.StartedAt)
	}

	if err = table.First(&assetUser).Error; err != nil {
		return
	}
	return
}

// InsertOrUpdate ..
func (u *AssetUser) InsertOrUpdate() (id int, err error) {
	existedModel, existErr := u.Get()
	if existErr != nil {
		if err = db.Eloquent.Table(u.TableName()).Create(&u).Error; err != nil {
			id = u.ID
			return
		}
	}
	if err = db.Eloquent.Table(u.TableName()).Model(&existedModel).Updates(&u).Error; err != nil {
		return
	}
	id = existedModel.ID
	return
}

// UpdateEndAt ..
func (u *AssetUser) UpdateEndAt() error {
	originItem, existErr := u.Get()
	if existErr != nil {
		return existErr
	}
	endAt := tools.GetCurrentTime("")
	originItem.EndedAt = &endAt
	err := db.Eloquent.Table(u.TableName()).Save(&originItem).Error
	return err
}
