package role

import (
	log "unknwon.dev/clog/v2"
)

// Service ..
type Service struct {
	RoleDao *SysRole
}

// NewService ..
func NewService() *Service {
	return &Service{
		RoleDao: &SysRole{},
	}
}

// GetRoleItemByKey ..
func (r *Service) GetRoleItemByKey(role string) (item SysRole, err error) {
	r.RoleDao.RoleKey = role
	item, err = r.RoleDao.Get()
	if err != nil {
		log.Error("when create user get role by role key error: %s", err.Error())
		return
	}
	return
}
