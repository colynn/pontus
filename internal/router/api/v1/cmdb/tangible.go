package cmdb

import (
	"errors"
	"time"

	"github.com/colynn/pontus/internal/pkg/customerror"
	"github.com/colynn/pontus/internal/pkg/response"
	"github.com/colynn/pontus/tools"
	"github.com/colynn/pontus/tools/file"

	cmdbsvc "github.com/colynn/pontus/internal/cmdb"

	log "unknwon.dev/clog/v2"

	"github.com/gin-gonic/gin"
)

// ImportPhysicalData ..
func ImportPhysicalData(c *gin.Context) {
	f, fh, err := c.Request.FormFile("file")
	customerror.HasError(err, "", 500)

	if fh.Size > 1024*1024 {
		customerror.HasError(errors.New("导入文件不允许超过1M"), "", 500)
	}
	defer f.Close()
	cmdbSvc := cmdbsvc.NewService()
	filePath, err := file.StorageImportFile(f, fh)
	customerror.HasError(err, "", 500)
	err = cmdbSvc.InsertOrUpdatePhysicalDevice(filePath)
	customerror.HasError(err, "", 500)
	response.OK(c, filePath, "添加成功")
}

// GetPhysicallist ..
func GetPhysicallist(c *gin.Context) {
	req := cmdbsvc.TangibleAssetReq{}
	cmdbSvc := cmdbsvc.NewService()
	err := c.Bind(&req)
	if req.PageIndex == 0 {
		// reset to default value
		req.PageIndex = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}
	req.Type = "physical"
	log.Trace("params: %+v", req)
	result, count, err := cmdbSvc.TangibleAsset.GetPage(&req)
	customerror.HasError(err, "", 500)
	response.PageOK(c, result, count, req.PageIndex, req.PageSize, "success")
}

// ImportTerminalDeviceData ..
func ImportTerminalDeviceData(c *gin.Context) {
	f, fh, err := c.Request.FormFile("file")
	customerror.HasError(err, "", 500)

	if fh.Size > 1024*1024 {
		customerror.HasError(errors.New("导入文件不允许超过1M"), "", 500)
	}
	defer f.Close()
	cmdbSvc := cmdbsvc.NewService()
	filePath, err := file.StorageImportFile(f, fh)
	customerror.HasError(err, "", 500)
	err = cmdbSvc.InsertOrUpdateTerminalDevice(filePath)
	customerror.HasError(err, "", 500)
	response.OK(c, filePath, "添加成功")
}

// GetPClist ..
func GetPClist(c *gin.Context) {
	req := cmdbsvc.TangibleAssetReq{}
	cmdbSvc := cmdbsvc.NewService()
	err := c.Bind(&req)
	if req.PageIndex == 0 {
		// reset to default value
		req.PageIndex = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}
	req.PCList = true
	log.Trace("params: %+v", req)
	result, count, err := cmdbSvc.TangibleAsset.GetPage(&req)
	customerror.HasError(err, "", 500)
	response.PageOK(c, result, count, req.PageIndex, req.PageSize, "success")
}

// GetTangibleAssetInfo ..
func GetTangibleAssetInfo(c *gin.Context) {
	data := cmdbsvc.TangibleAsset{}
	data.InstanceID = c.Param("instanceID")
	result, err := data.Get()
	customerror.HasError(err, "抱歉未找到相关信息", 500)
	response.OK(c, result, "")
}

// UpdateTangibleAssetInfo ..
func UpdateTangibleAssetInfo(c *gin.Context) {
	instanceID := c.Param("instanceID")
	req := cmdbsvc.TangibleAssetUpdate{}
	err := c.ShouldBindJSON(&req)
	customerror.HasError(err, "更新参数错误", 500)
	log.Trace("update params: %+v", req)
	// req to model
	dataModel := cmdbsvc.TangibleAsset{
		InstanceID:      instanceID,
		ID:              req.ID,
		AssetNumber:     req.AssetNumber,
		Manufactory:     req.Manufactory,
		Type:            req.Type,
		Configuration:   req.Configuration,
		ProcurementType: req.ProcurementType,
		VmwareEnabled:   req.VmwareEnabled,
		CPU:             req.CPU,
		Memory:          req.Memory,
		Disk:            req.Disk,
		PrivateIP:       req.PrivateIP,
		PublicIP:        req.PublicIP,
		InvoiceTime: func() *time.Time {
			if req.InvoiceTime == "" {
				return nil
			}
			invoiceTime, err := tools.ParseStrToDateTime(req.InvoiceTime)
			if err != nil {
				log.Error("parse InvoiceTime error: %s", err.Error())
				return nil
			}
			return &invoiceTime
		}(),
		DeliveryDate: func() *time.Time {
			if req.DeliveryDate == "" {
				return nil
			}
			deliveryDate, err := tools.ParseStrToDate(req.DeliveryDate)
			if err != nil {
				log.Error("parse deliveryDate error: %s", err.Error())
				return nil
			}
			return &deliveryDate
		}(),
		InvoiceNumber:     req.InvoiceNumber,
		PretaxAmount:      req.PretaxAmount,
		PretaxGrossAmount: req.PretaxGrossAmount,
		InvoiceDiscount:   req.InvoiceDiscount,
		AssetLife:         req.AssetLife,
		WarrantyPeriod:    req.WarrantyPeriod,
		ResidualRatio:     req.ResidualRatio,
		InventoryTime: func() *time.Time {
			if req.InventoryTime == "" {
				return nil
			}
			inventoryTime, err := tools.ParseStrToDate(req.InventoryTime)
			if err != nil {
				log.Error("parse inventoryTime error: %s", err.Error())
				return nil
			}
			return &inventoryTime
		}(),
		SerialNumber: req.SerialNumber,
		Region:       req.Region,
		RecipientID:  req.RecipientID,
		Status:       req.Status,
		// ScrappedTime:      req.ScrappedTime,
		ScrappedPrice: req.ScrappedPrice,
		// SurrenderedTime:   req.SurrenderedTime,
		Recipient:   req.Recipient,
		Description: req.Description,
	}
	result, err := dataModel.UpdateItem(dataModel.InstanceID)
	customerror.HasError(err, "抱歉未找到相关信息", 500)
	response.OK(c, result, "")
}

// DeleteTangibleAssetInfo ..
func DeleteTangibleAssetInfo(c *gin.Context) {
	data := cmdbsvc.TangibleAsset{}
	data.InstanceID = c.Param("instanceID")
	result, err := data.DeleteItem()
	customerror.HasError(err, "抱歉未找到相关信息", 500)
	response.OK(c, result, "")
}
