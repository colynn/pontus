package cmdb

import (
	"github.com/colynn/pontus/internal/account/user"
	"github.com/colynn/pontus/internal/db"

	log "unknwon.dev/clog/v2"
)

// GetUserRealName ..
func GetUserRealName(userID int) string {
	if userID == 0 {
		return ""
	}
	value := user.SysUser{}
	if err := db.Eloquent.Select("*").Table("sys_user").Where("id = ?", userID).First(&value).Error; err != nil {
		log.Warn("get sys user by id error: %v", err.Error())
	}
	if value.RealName != "" {
		return value.RealName
	}
	return value.Username
}

type userID struct {
	ID int
}

// AssetRecipientUserID ..
type AssetRecipientUserID struct {
	UserID int
}

// getCurrentRecipientID ..
func getCurrentRecipientID(username string) (userIDs []AssetRecipientUserID, err error) {
	sysUsers := []userID{}
	err = db.Eloquent.Select("id").Table("sys_user").Where("deleted_at is NULL").Where("username like ?", "%"+username+"%").Or("real_name like ?", "%"+username+"%").Find(&sysUsers).Error
	if err != nil {
		return
	}
	ids := make([]int, 0, len(sysUsers))
	for _, idItem := range sysUsers {
		ids = append(ids, idItem.ID)
	}
	log.Trace("sys user ids: %v", ids)
	err = db.Eloquent.Select("user_id").Table("cmdb_asset_user").Where("deleted_at is NULL and ended_at is NULL").Where("user_id in (?)", ids).Find(&userIDs).Error
	return
}

func getNeedDeleteTags(origin, new []AssetTag, instanceID string) []AssetTag {
	needDeteteTags := make([]AssetTag, 0, len(origin))
	for _, t := range origin {
		deleteTag := true
		for _, new := range new {
			if t.TagKey == new.TagKey && t.TagValue == new.TagValue {
				deleteTag = false
				break
			}
		}
		if deleteTag {
			// log.Trace("tag-> id: %d, instance id: %s", t.ID, instanceID)
			needDeteteTags = append(needDeteteTags, t)
		}
	}
	return needDeteteTags
}
