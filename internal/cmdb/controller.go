package cmdb

import (
	"fmt"
	"sync"

	"github.com/colynn/pontus/internal/pkg/customerror"
	"github.com/colynn/pontus/internal/pkg/response"
	"github.com/colynn/pontus/tools"

	excel "github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	log "unknwon.dev/clog/v2"
)

// CreateTag ..
func CreateTag(c *gin.Context) {
	tag := &AssetTag{}
	req := CreateTagReq{}
	err := c.MustBindWith(&req, binding.JSON)
	customerror.HasError(err, "非法数据格式", 500)
	id, err := tag.Insert(&req)
	customerror.HasError(err, "", 500)
	response.OK(c, id, "添加成功")
}

// CreateMutilpleServer ..
func CreateMutilpleServer(c *gin.Context) {
	servers := []*CloudAsset{}
	err := c.MustBindWith(&servers, binding.JSON)
	customerror.HasError(err, "非法数据格式", 500)

	nonStandard := make(chan NonStandardInstance)
	log.Info("sync aliyun asset items: %v", len(servers))
	syncPool := make(chan bool, 40)
	var wg sync.WaitGroup
	for _, item := range servers {
		syncPool <- true
		wg.Add(1)
		go func(item *CloudAsset) {
			defer wg.Done()
			_, nonStanardInstance, err := item.InsertOrUpdate(item)
			if nonStanardInstance.InstanceID != "" {
				nonStandard <- nonStanardInstance
			}
			if err != nil {
				log.Error("item: %s update error: %s", item.InstanceID, err.Error())
			}
			<-syncPool
		}(item)
	}
	go func() {
		wg.Wait()
		close(nonStandard)
	}()
	log.Trace("update total len: %d, non standard ecs instances len: %v", len(servers), len(nonStandard))
	for item := range nonStandard {
		log.Trace("--->item private ip: %s, instance name: %s", item.IP, item.InstanceName)
	}
	// TODO: send email when instance name is invalid or add product/service item.
	response.OK(c, "", "添加成功")
}

// GetServersByPagination  ..
func GetServersByPagination(c *gin.Context) {
	req := AssetReq{}
	s := &CloudAsset{}
	err := c.Bind(&req)
	customerror.HasError(err, "", 400)
	result, count, err := s.ServerList(&req)
	customerror.HasError(err, "", 500)
	response.PageOK(c, result, count, req.PageIndex, req.PageSize, "success")
}

// AssetStatistics ..
func AssetStatistics(c *gin.Context) {
	s := CloudAsset{}
	result, err := s.AssetStatistics()
	customerror.HasError(err, "", 500)
	response.OK(c, result, "")
}

// ExportCloudAsset ..
func ExportCloudAsset(c *gin.Context) {
	req := AssetReq{}
	s := &CloudAsset{}
	err := c.Bind(&req)
	customerror.HasError(err, "", 500)
	result, err := s.GetExportData(&req)
	customerror.HasError(err, "", 500)
	w := c.Writer
	filename := fmt.Sprintf("aliyun-asset-%s.xlsx", tools.GetUniqString())
	w.Header().Set("Content-Type", "multipart/form-data")
	w.Header().Set("Content-disposition", "attachment; filename="+filename)
	// 设置 excel 头行
	titles := []string{"实例ID", "实例名称", "可用区", "内网IP", "配置", "操作系统", "环境", "一级部门", "二级部门", "产品线", "状态"}
	xlsx := excel.NewFile()
	xlsx.SetSheetRow("Sheet1", "A1", &titles)
	var env, primary, second, business string
	for k, item := range result {
		axis := "A" + fmt.Sprintf("%d", k+2)
		env = getTagValueBaseONKey("environment", item)
		primary = getTagValueBaseONKey("primary", item)
		second = getTagValueBaseONKey("second", item)
		business = getTagValueBaseONKey("business", item)

		xlsx.SetSheetRow(
			"Sheet1", axis, &[]interface{}{
				item.InstanceID,
				item.InstanceName,
				item.ZoneID,
				item.PrivateIP,
				item.InstanceType,
				item.OSType,
				env,
				primary,
				second,
				business,
				item.Status,
			})
	}
	xlsx.Write(w)
}

func getTagValueBaseONKey(key string, item *CloudAsset) string {
	var tagValue string
	hasKey := false
	for _, tag := range item.Tags {
		if tag.TagKey == key {
			tagValue = tag.TagValue
			hasKey = true
			break
		}
	}
	if !hasKey {
		return ""
	}
	return tagValue
}
