package cmdb

import (
	"errors"
	"time"

	"github.com/colynn/pontus/internal/db"
	"github.com/colynn/pontus/pkg/intersect"

	"github.com/jinzhu/gorm"
	log "unknwon.dev/clog/v2"
)

/*
sync aliyun cloud resource data
*/

// CloudAsset model defined
type CloudAsset struct {
	db.BaseModel
	ID           int    `gorm:"primary_key;AUTO_INCREMENT"  json:"id"`           // 编码
	InstanceID   string `gorm:"primary_key;type:varchar(128)" json:"InstanceID"` // eg. i-xxx
	InstanceName string `gorm:"type:varchar(128)" json:"InstanceName"`           // eg. shb-manage-app-wps-node01
	Manufactory  string `gorm:"type:varchar(128)" json:"Manufactory"`            // eg. aliyun
	Type         string `gorm:"type:varchar(64)" json:"Type"`                    // eg. ecs,slb,
	ZoneID       string `gorm:"type:varchar(64)" json:"ZoneID"`                  // eg. cn-shanghai-b
	VpcID        string `gorm:"type:varchar(64)" json:"VpcID"`
	SerialNumber string `gorm:"type:varchar(64)" json:"SerialNumber"` // eg. 8e8a6016-15eb-4da7-a177-6f423dfa339a

	InstanceType            string     `gorm:"type:varchar(128)" json:"InstanceType"` // eg. ecs.sn2ne.2xlarge
	CPU                     int        `gorm:"type:int(11)" json:"CPU"`
	Memory                  int        `gorm:"type:int(11)" json:"Memory"`
	IoOptimized             bool       `gorm:"type:bool" json:"IoOptimized"`
	OSType                  string     `gorm:"type:varchar(128)" json:"OSType"`
	OSName                  string     `gorm:"type:varchar(128)" json:"OSName"`
	PublicIP                string     `gorm:"type:varchar(128)" json:"PublicIP"`
	EipAddress              string     `gorm:"type:varchar(128)" json:"EipAddress"`
	PrivateIP               string     `gorm:"type:varchar(128)" json:"PrivateIp"`
	InternetMaxBandwidthOut int        `gorm:"type:int(11)" json:"InternetMaxBandwidthOut"`
	InstanceChargeType      string     `gorm:"type:varchar(64)" json:"InstanceChargeType"` //实例的计费方式 PostPaid：按量付费 PrePaid：包年包月
	Status                  string     `gorm:"type:varchar(128)" json:"Status"`
	CreationTime            time.Time  `json:"CreationTime"`
	ExpiredTime             *time.Time `json:"ExpiredTime"`
	// Tags         []ServerTag `gorm:"many2many:server_tags;"`
	Tags []AssetTag `gorm:"association_autoupdate:false;association_autocreate:false;foreignkey:InstanceID;association_foreignkey:ID;jointable_foreignkey:asset_instance_id;many2many:cmdb_asset_tags;"`
}

// TableName ..
func (CloudAsset) TableName() string {
	return "cmdb_cloud_asset"
}

// AssetTag server tags model defined
type AssetTag struct {
	db.BaseModel
	ID          int    `gorm:"primary_key;AUTO_INCREMENT"  json:"id"`
	TagKey      string `gorm:"unique_index:idx_key_value;type:varchar(128)" json:"tag_key"`   // eg. business
	TagValue    string `gorm:"unique_index:idx_key_value;type:varchar(128)" json:"tag_value"` // eg. tms
	Description string `gorm:"type:varchar(128)" json:"description"`
	// Assets      []CloudAsset `gorm:"many2many:cmdb_asset_tags"`
}

// TableName ..
func (AssetTag) TableName() string {
	return "cmdb_asset_tag"
}

// AssetTags ..
type AssetTags struct {
	AssetInstanceID string `gorm:"asset_instance_id"`
	AssetTagID      int    `gorm:"asset_tag_id"`
}

// TableName ..
func (AssetTags) TableName() string {
	return "cmdb_asset_tags"
}

// InsertOrUpdate ..
func (s *CloudAsset) InsertOrUpdate(item *CloudAsset) (id int, nonstanardInstance NonStandardInstance, err error) {
	existItem := CloudAsset{}
	db.Eloquent.Table(s.TableName()).Where("instance_id = ?", item.InstanceID).First(&existItem)
	if existItem.ID != 0 {
		existItem.InstanceID = item.InstanceID
		existItem.InstanceName = item.InstanceName
		existItem.Manufactory = item.Manufactory
		existItem.Type = item.Type
		existItem.ZoneID = item.ZoneID
		existItem.VpcID = item.VpcID
		existItem.SerialNumber = item.SerialNumber
		existItem.InstanceType = item.InstanceType
		existItem.CPU = item.CPU
		existItem.Memory = item.Memory
		existItem.IoOptimized = item.IoOptimized
		existItem.OSType = item.OSType
		existItem.OSName = item.OSName
		existItem.PublicIP = item.PublicIP
		existItem.EipAddress = item.EipAddress
		existItem.PrivateIP = item.PrivateIP
		existItem.InternetMaxBandwidthOut = item.InternetMaxBandwidthOut
		existItem.InstanceChargeType = item.InstanceChargeType
		existItem.Status = item.Status
		existItem.CreationTime = item.CreationTime
		existItem.ExpiredTime = item.ExpiredTime
		if err = db.Eloquent.Table(s.TableName()).Save(&existItem).Error; err != nil {
			return
		}
	} else {
		// 添加数据
		if err = db.Eloquent.Table(s.TableName()).Create(&item).Error; err != nil {
			log.Error("create asset item, error: %v", err.Error())
			return
		}
	}

	// origin
	origin := CloudAsset{}
	if err = db.Eloquent.Table(s.TableName()).Preload("Tags").Where("instance_id = ?", item.InstanceID).First(&origin).Error; err != nil {
		return
	}
	// when tags already delete, need delete item
	needDeteteTags := getNeedDeleteTags(origin.Tags, item.Tags, item.InstanceID)
	_ = s.deleteOriginTags(needDeteteTags)

	for _, _item := range item.Tags {
		t := &AssetTag{
			TagKey:   _item.TagKey,
			TagValue: _item.TagValue,
		}
		exsited := &AssetTag{}
		tagID := 0
		db.Eloquent.Table(t.TableName()).Where("tag_key = ? and tag_value = ?", t.TagKey, t.TagValue).First(exsited)
		if exsited.ID == 0 {
			if err := db.Eloquent.Table(t.TableName()).Create(t).Error; err != nil {
				log.Error("tags insert key: %v, value: %v error: %v", t.TagKey, t.TagValue, err.Error())
				continue
			}
			tagID = t.ID
		} else {
			tagID = exsited.ID
		}

		st := &AssetTags{
			AssetInstanceID: item.InstanceID,
			AssetTagID:      tagID,
		}
		var tagMappingCount int
		if err := db.Eloquent.Table(st.TableName()).Where("asset_instance_id = ? and asset_tag_id = ?", s.InstanceID, tagID).Count(&tagMappingCount).Error; err != nil {
			log.Error("query tag mapping error")
		}
		if tagMappingCount == 0 {
			if err := db.Eloquent.Table(st.TableName()).Create(st).Error; err != nil {
				log.Error("asset tags insert error")
				continue
			}
		}
	}
	id = s.ID
	// TODO: aliyun app service comment
	// if item.Type != "ecs" {
	// 	return
	// }
	// verifyTime, err := istools.ParseStrToDate("2020-08-26", "YYYY-MM-DD", "Asia/Shanghai")
	// if istools.TimeComparison(item.CreationTime, verifyTime) == 1 {
	// 	// new ecs parse hostname
	// 	serviceItem := &CloudServiceAsset{}
	// 	serviceItem, err := generateAliServiceItem(item, serviceItem, allDepts, allProducts, allServices)
	// 	if err != nil {
	// 		log.Error("generate ali service item error: %s", err.Error())
	// 		nonstanardInstance = NonStandardInstance{
	// 			IP:           item.PrivateIP,
	// 			InstanceID:   item.InstanceID,
	// 			InstanceName: item.InstanceName,
	// 		}
	// 	} else {
	// 		id, err := serviceItem.InsertOrUpdate(serviceItem)
	// 		if err != nil {
	// 			log.Trace("upate ali service item error: %s", err.Error())
	// 		} else {
	// 			log.Trace("cmdb_cloud_service_asset id: %d, update success", id)
	// 		}
	// 	}

	// }
	return
}

// GetTagIDByKeyAndValue ..
func (s *CloudAsset) GetTagIDByKeyAndValue(key, value string) (tag AssetTag, err error) {
	exsited := &AssetTag{}
	if err = db.Eloquent.Table(exsited.TableName()).Where("tag_key = ? and tag_value = ?", key, value).First(&tag).Error; err != nil {
		return
	}
	return
}

// ServerList ..
func (s *CloudAsset) ServerList(req *AssetReq) ([]*CloudAsset, int, error) {
	var doc []*CloudAsset
	query := s.CloudAssetQuery(req)
	var count int
	if err := query.Preload("Tags").Offset((req.PageIndex - 1) * req.PageSize).Limit(req.PageSize).Find(&doc).Error; err != nil {
		return nil, 0, err
	}
	query.Count(&count)
	return doc, count, nil
}

// GetExportData ..
func (s *CloudAsset) GetExportData(req *AssetReq) ([]*CloudAsset, error) {
	var items []*CloudAsset
	query := s.CloudAssetQuery(req)
	if err := query.Preload("Tags").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// AssetStatistics ..
func (s *CloudAsset) AssetStatistics() (*AssetStatisticsResp, error) {
	assetItems := []*CategoryItem{}
	query := db.Eloquent.Table(s.TableName())
	query.Select("type as name, count(instance_id) as value").Group("type").Scan(&assetItems)

	// environment
	envStatisItems, err := statisticsBasedTag("environment")
	if err != nil {
		// TOOD: logger
		return nil, err
	}

	// fl-department
	flDepartMentStatisItems, err := statisticsBasedTag("fl-department")
	if err != nil {
		// TOOD: logger
		return nil, err
	}
	ecsStatisItems := []*CategoryItem{}
	ecsRegionItems := []*CategoryItem{}

	slbItems := []*CategoryItem{}
	ossItems := []*CategoryItem{}
	nasItems := []*CategoryItem{}
	query.Select("status as name, count(instance_id) as value").Where("type = ?", "ecs").Group("status").Scan(&ecsStatisItems)
	query.Select("zone_id as name, count(instance_id) as value").Where("type = ?", "ecs").Group("zone_id").Scan(&ecsRegionItems)

	// other
	query.Select("status as name, count(instance_id) as value").Where("type = ?", "slb").Group("status").Scan(&slbItems)
	query.Select("status as name, count(instance_id) as value").Where("type = ?", "oss").Group("status").Scan(&ossItems)
	query.Select("status as name, count(instance_id) as value").Where("type = ?", "nas").Group("status").Scan(&nasItems)

	for _, item := range ecsRegionItems {
		zoneRunning := &CategoryItem{}
		query.Select("count(instance_id) as running").Where("status = 'Running' and zone_id = ?", item.Name).Scan(&zoneRunning)
		item.Running = zoneRunning.Running
		item.Stopped = item.Value - zoneRunning.Running
	}
	rsp := &AssetStatisticsResp{
		AssetCategories:      assetItems,
		EnvCategories:        envStatisItems,
		DepartmentCategories: flDepartMentStatisItems,
		ServerCategories:     ecsStatisItems,
		ECSRegionCategories:  ecsRegionItems,
		NASCategories:        nasItems,
		OSSCategories:        ossItems,
		SLBCategories:        slbItems,
	}

	return rsp, nil
}

// Insert ..
func (t *AssetTag) Insert(req *CreateTagReq) (id int, err error) {
	t.TagKey = req.TagKey
	t.TagValue = req.TagValue
	t.Description = req.Description

	var count int
	db.Eloquent.Table(t.TableName()).Where("tag_key = ? and tag_value = ?", req.TagKey, req.TagValue).Count(&count)
	if count > 0 {
		err = errors.New("tag already existed")
		return 0, err
	}

	//添加数据
	if err = db.Eloquent.Table(t.TableName()).Create(&t).Error; err != nil {
		return 0, err
	}
	return t.ID, nil
}

func specificFilterQuery(query *gorm.DB, key, value string, slices []interface{}) []interface{} {
	servers := []*CloudAsset{}
	signalQuery := query.Where("cmdb_asset_tag.tag_key = ? and cmdb_asset_tag.tag_value= ?", key, value)
	signalQuery.Find(&servers)
	for _, server := range servers {
		slices = append(slices, server.InstanceID)
	}
	return slices
}

func statisticsBasedTag(key string) ([]*CategoryItem, error) {
	// TODO: use TableName()
	tagQuery := db.Eloquent.Table("cmdb_asset_tag")
	IDQuery := db.Eloquent.Table("cmdb_asset_tags")
	tags := []*AssetTag{}
	if err := tagQuery.Where("tag_key = ?", key).Find(&tags).Error; err != nil {
		log.Error("tag_key %s query errror: %s", key, err.Error())
		return nil, err
	}
	items := []*CategoryItem{}
	var count int
	for _, tag := range tags {
		IDQuery.Where("asset_tag_id = ?", tag.ID).Count(&count)
		if count == 0 {
			continue
		}
		items = append(items, &CategoryItem{
			Name:  tag.TagValue,
			Value: count,
		})
	}
	return items, nil
}

func (s *CloudAsset) deleteOriginTags(tags []AssetTag) error {
	for _, deleteItem := range tags {
		tag, err := s.GetTagIDByKeyAndValue(deleteItem.TagKey, deleteItem.TagValue)
		if err != nil {
			log.Error("when delete existed tag, get tag by key error: %s, skipp delete this tag relevance", err.Error())
			continue
		}
		// TODO: change to use TableName
		if err := db.Eloquent.Table("cmdb_asset_tags").Where("asset_instance_id = ? and asset_tag_id = ?", s.InstanceID, tag.ID).Delete(&AssetTags{}).Error; err != nil {
			log.Error("when delete tags error: %s", err.Error())
			continue
		}
	}
	return nil
}

// CloudAssetQuery ..
func (s *CloudAsset) CloudAssetQuery(req *AssetReq) *gorm.DB {
	query := db.Eloquent.Table(s.TableName())

	joinQuery := query.Joins("inner join cmdb_asset_tags on cmdb_asset_tags.asset_instance_id=cmdb_cloud_asset.instance_id").
		Joins("inner join cmdb_asset_tag on cmdb_asset_tags.asset_tag_id = cmdb_asset_tag.id")
	instanceIDsWithPrimary := []interface{}{}
	instanceIDsWithSecond := []interface{}{}
	instanceIDsWithBusiness := []interface{}{}
	instanceIDsWithEnv := []interface{}{}
	tagFilterEnable := false
	tagFilterTime := 0
	tagFilterIds := []interface{}{}

	if req.FlDepartment != "" {
		tagFilterEnable = true
		tagFilterTime++
		instanceIDsWithPrimary = specificFilterQuery(joinQuery, "fl-department", req.FlDepartment, instanceIDsWithPrimary)
		if tagFilterTime == 1 {
			tagFilterIds = instanceIDsWithPrimary
		} else {
			tagFilterIds = intersect.Hash(instanceIDsWithPrimary, tagFilterIds)
		}
	}

	if req.TlDepartment != "" {
		tagFilterEnable = true
		tagFilterTime++
		instanceIDsWithSecond = specificFilterQuery(joinQuery, "tl-department", req.TlDepartment, instanceIDsWithSecond)
		if tagFilterTime == 1 {
			tagFilterIds = instanceIDsWithSecond
		} else {
			tagFilterIds = intersect.Hash(instanceIDsWithSecond, tagFilterIds)
		}
	}

	if req.Business != "" {
		tagFilterEnable = true
		tagFilterTime++
		instanceIDsWithBusiness = specificFilterQuery(joinQuery, "business", req.Business, instanceIDsWithBusiness)
		if tagFilterTime == 1 {
			tagFilterIds = instanceIDsWithBusiness
		} else {
			tagFilterIds = intersect.Hash(instanceIDsWithBusiness, tagFilterIds)
		}
	}
	if req.Environment != "" {
		tagFilterEnable = true
		tagFilterTime++
		instanceIDsWithEnv = specificFilterQuery(joinQuery, "environment", req.Environment, instanceIDsWithEnv)
		if tagFilterTime == 1 {
			tagFilterIds = instanceIDsWithEnv
		} else {
			tagFilterIds = intersect.Hash(instanceIDsWithEnv, tagFilterIds)
		}
	}
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	if req.Type != "" {
		query = query.Where("type = ?", req.Type)
	}

	if req.OSType != "" {
		query = query.Where("os_type = ?", req.OSType)
	}

	if req.AssetIP != "" {
		query = query.Where("private_ip like ?", "%"+req.AssetIP+"%").Or("public_ip like ?", "%"+req.AssetIP+"%").Or("eip_address like ?", "%"+req.AssetIP+"%")
	}

	if req.AssetName != "" {
		query = query.Where("instance_name like ?", "%"+req.AssetName+"%")
	}

	if req.StartAt != "" && req.EndAt != "" {
		startAt, _ := time.Parse("2006-01-02", req.StartAt)
		endAt, _ := time.Parse("2006-01-02", req.EndAt)
		query = query.Where("creation_time >= ? AND creation_time <= ?", startAt, endAt)
	}

	if tagFilterEnable {
		query = query.Where("instance_id in (?)", tagFilterIds)
	}
	return query
}
