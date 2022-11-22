package middleware

import (
	"io/ioutil"
	"reflect"
	"strings"
	"time"

	"github.com/colynn/pontus/internal/audit"
	"github.com/colynn/pontus/internal/cmdb"
	"github.com/colynn/pontus/internal/pkg/accountinfo"

	istools "github.com/isbrick/tools"

	"github.com/gin-gonic/gin"
	log "unknwon.dev/clog/v2"
)

// SysAudit 编辑审计
func SysAudit() gin.HandlerFunc {

	return func(c *gin.Context) {

		// 请求方式
		reqMethod := c.Request.Method
		updateType := getUpdateType(reqMethod)

		// 请求路由
		reqURI := c.Request.RequestURI
		instanceID := c.Param("instanceID")
		var content interface{}
		if reqMethod == "PUT" {
			// before 处理请求前
			if strings.HasPrefix(reqURI, "/api/v1/assets/tangibles/") {
				//
				if instanceID == "" {
					c.Next()
					return
				}
				originItem := getAssetItemByInstanceID(instanceID)

				// 处理请求
				c.Next()
				newItem := getAssetItemByInstanceID(instanceID)
				content = getUpdateContent(originItem, newItem)
				createAuditRecord(accountinfo.GetUserID(c), updateType, accountinfo.GetUserName(c), instanceID, content)

			} else {
				log.Trace("did not defined audit, skip")
			}
		} else {
			// 处理请求
			c.Next()
			if reqMethod == "POST" {
				content = "创建条目"
			}
			switch reqMethod {
			case "POST":
				content = "创建条目"
			case "DELETE":
				content = "删除条目"
			default:
				log.Warn("request method %s, is unexception", reqMethod)
			}
			createAuditRecord(accountinfo.GetUserID(c), updateType, accountinfo.GetUserName(c), instanceID, content)
		}
	}
}

// create audit record
func createAuditRecord(userID, updateType int, username, instanceID string, content interface{}) {
	if content == nil || content == "" {
		// ignore create audit record
		log.Trace("content is nil, skip audit record")
		return
	}
	auditInstance := &audit.SysAudit{
		Type:       updateType,
		UserID:     userID,
		UserName:   username,
		InstanceID: instanceID,
	}
	_, err := auditInstance.Insert(content)
	if err != nil {
		log.Error("create audit log error: %s", err.Error())
	}
}

func getAssetItemByInstanceID(instanceID string) (assetItem cmdb.TangibleAsset) {
	var err error
	assetItem.InstanceID = instanceID
	assetItem, err = assetItem.Get()
	if err != nil {
		log.Error("get asset item %s error: %s", assetItem.InstanceID, err.Error())
	}
	log.Trace("instance id: %s", assetItem.InstanceID)
	return
}

// getUpdateContent ..
func getUpdateContent(originItem, newItem cmdb.TangibleAsset) (content interface{}) {
	var updateItem []audit.UpdateItem

	originRef := reflect.ValueOf(originItem)
	typeOfOrigin := originRef.Type()

	new := reflect.ValueOf(newItem)
	typeOfNew := new.Type()
	for i := 0; i < originRef.NumField(); i++ {
		if typeOfOrigin.Field(i).Name == "BaseModel" {
			continue
		}
		for j := 0; j < new.NumField(); j++ {
			if typeOfOrigin.Field(j).Name == "BaseModel" {
				continue
			}
			if typeOfOrigin.Field(i).Name == typeOfNew.Field(j).Name {
				if istools.IsSliceContainsStr([]string{"DeliveryDate", "InventoryTime", "InvoiceTime"}, typeOfOrigin.Field(i).Name) {
					// TODO: change it to pkg.
					originTimeValue := originRef.Field(i).Elem()
					newTimeValue := new.Field(j).Elem()
					var originTime, newTime interface{}

					if originTimeValue.IsValid() {
						originTime = originTimeValue.Interface().(time.Time)
					} else {
						originTime = ""
					}

					if newTimeValue.IsValid() {
						newTime = newTimeValue.Interface().(time.Time)
					} else {
						newTime = ""
					}
					if originTime == newTime {
						continue
					}
					log.Trace("origin Time: %v, %v", originTime, newTime)
				} else {
					if originRef.Field(i).Interface() == new.Field(j).Interface() {
						continue
					}
				}
				log.Trace("orgin value: %v, new value: %v", originRef.Field(i).Interface(), new.Field(j).Interface())
				item := audit.UpdateItem{
					Field:       typeOfOrigin.Field(i).Name,
					OriginValue: originRef.Field(i).Interface(),
					NewValue:    new.Field(j).Interface(),
				}
				updateItem = append(updateItem, item)
			}
		}
	}
	if len(updateItem) == 0 {
		content = nil
		log.Trace("content reset to nil")
	} else {
		content = audit.Content(updateItem)
	}
	return
}

func getUpdateType(reqMethod string) (updateType int) {
	switch reqMethod {
	case "PUT":
		updateType = 2
	case "POST":
		updateType = 1
	case "DELETE":
		updateType = 3
	default:
		log.Trace("request method: %s unknown", reqMethod)
	}
	return updateType
}

func catchRequestBody(c *gin.Context) []byte {
	// Read the Body content
	var bodyBytes []byte
	if c.Request.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
	}
	return bodyBytes
}
