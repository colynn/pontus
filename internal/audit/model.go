package audit

import (
	app "github.com/colynn/pontus/internal"
	"github.com/colynn/pontus/internal/db"

	log "unknwon.dev/clog/v2"
)

// audit type
const (
	Created = iota + 1
	Updated
	Deleted
)

// SysAudit ..
type SysAudit struct {
	ID         int    `json:"ID" gorm:"primary_key;AUTO_INCREMENT"` // 编码
	UserID     int    `json:"UserID" gorm:"type:int"`               //
	UserName   string `json:"UserName" gorm:"type:varchar(128)"`    //
	InstanceID string `json:"InstanceID" gorm:"type:varchar(64)"`   //
	Content    string `gorm:"type:varchar(1048);"`                  // 更新字段
	Type       int    `json:"type" gorm:"type:int"`                 // audit type
	db.BaseModel
}

// TableName ..
func (SysAudit) TableName() string {
	return "sys_audit"
}

// GetPage ..
func (r *SysAudit) GetPage(req *app.PaginationReq, instanceID string) ([]*Rsp, int, error) {
	var doc []*SysAudit

	table := db.Eloquent.Select("*").Table(r.TableName()).Where("deleted_at is NULL").Where("instance_id = ?", instanceID)

	var count int

	if err := table.Order("created_at desc").Offset((req.PageIndex - 1) * req.PageSize).Limit(req.PageSize).Find(&doc).Error; err != nil {
		return nil, 0, err
	}
	table.Count(&count)

	rsp := make([]*Rsp, 0, len(doc))
	for _, item := range doc {
		content := Content{}
		content.Struct(item.Content)

		rspItem := &Rsp{
			ID:        item.ID,
			UserID:    item.UserID,
			Type:      item.Type,
			Username:  item.UserName,
			CreatedAt: item.CreatedAt,
			Content:   content,
		}
		rsp = append(rsp, rspItem)
	}
	return rsp, count, nil
}

// Insert ..
func (r *SysAudit) Insert(content interface{}) (id int, err error) {
	r.Content = getContentStr(content)
	result := db.Eloquent.Table(r.TableName()).Create(&r)
	if result.Error != nil {
		err = result.Error
		return
	}
	id = r.ID
	return
}

// getContentStr
func getContentStr(content interface{}) (contentstr string) {
	var err error
	switch content.(type) {
	case Content:
		ContentStruct, _ := content.(Content)
		contentstr, err = ContentStruct.String()
		if err != nil {
			log.Error("content is []UpdateItme type, but can't use String func ")
			contentstr = err.Error()
		}
	case string:
		contentstr, _ = content.(string)
	default:
		log.Warn("content can not parse to string")
	}
	return
}
