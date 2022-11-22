package audit

import (
	app "github.com/colynn/pontus/internal"
	"github.com/colynn/pontus/internal/audit"
	"github.com/colynn/pontus/internal/pkg/customerror"
	"github.com/colynn/pontus/internal/pkg/response"

	log "unknwon.dev/clog/v2"

	"github.com/gin-gonic/gin"
)

// GetTangibleAuditLog ..
func GetTangibleAuditLog(c *gin.Context) {
	req := app.PaginationReq{}
	audiModel := audit.SysAudit{}
	err := c.Bind(&req)
	if req.PageIndex == 0 {
		// reset to default value
		req.PageIndex = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}
	instanceID := c.Param("instanceID")
	log.Trace("params: %+v", req)
	result, count, err := audiModel.GetPage(&req, instanceID)
	customerror.HasError(err, "", 500)
	response.PageOK(c, result, count, req.PageIndex, req.PageSize, "success")
}
