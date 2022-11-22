package role

import (
	"errors"

	"github.com/colynn/pontus/internal/db"
)

// SysRole ..
type SysRole struct {
	RoleID      int    `json:"roleID" gorm:"primary_key;AUTO_INCREMENT"` // 角色编码
	RoleName    string `json:"roleName" gorm:"type:varchar(128);"`       // 角色名称
	Status      string `json:"status" gorm:"type:varchar(1);"`           //
	RoleKey     string `json:"roleKey" gorm:"type:varchar(128);"`        //角色代码
	RoleSort    int    `json:"roleSort" gorm:"type:int(4);"`             //角色排序
	Flag        string `json:"flag" gorm:"type:varchar(128);"`           //
	CreateBy    string `json:"createBy" gorm:"type:varchar(128);"`       //
	UpdateBy    string `json:"updateBy" gorm:"type:varchar(128);"`       //
	Description string `json:"description" gorm:"type:varchar(255);"`    //备注
	Admin       bool   `json:"admin" gorm:"type:char(1);"`
	DataScope   string `json:"dataScope" gorm:"-"`
	Params      string `json:"params" gorm:"-"`
	MenuIds     []int  `json:"menuIds" gorm:"-"`
	DeptIds     []int  `json:"deptIds" gorm:"-"`
	db.BaseModel
}

// TableName ..
func (SysRole) TableName() string {
	return "sys_role"
}

// GetPage ..
func (r *SysRole) GetPage(pageSize int, pageIndex int) ([]SysRole, int, error) {
	var doc []SysRole

	table := db.Eloquent.Select("*").Table("sys_role").Where("deleted_at is NULL")
	if r.RoleID != 0 {
		table = table.Where("role_id = ?", r.RoleID)
	}
	if r.RoleName != "" {
		table = table.Where("role_name like ?", "%"+r.RoleName+"%")
	}
	if r.Status != "" {
		table = table.Where("status = ?", r.Status)
	}
	if r.RoleKey != "" {
		table = table.Where("role_key = ?", r.RoleKey)
	}

	var count int

	if err := table.Order("role_sort").Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&doc).Error; err != nil {
		return nil, 0, err
	}
	table.Count(&count)
	return doc, count, nil
}

// Get ..
func (r *SysRole) Get() (SysRole SysRole, err error) {
	table := db.Eloquent.Table("sys_role")
	if r.RoleID != 0 {
		table = table.Where("role_id = ?", r.RoleID)
	}
	if r.RoleName != "" {
		table = table.Where("role_name = ?", r.RoleName)
	}

	if r.RoleKey != "" {
		table = table.Where("role_key = ?", r.RoleKey)
	}
	if err = table.First(&SysRole).Error; err != nil {
		return
	}

	return
}

// GetList ..
func (r *SysRole) GetList() (SysRole []SysRole, err error) {
	table := db.Eloquent.Table("sys_role")
	if r.RoleID != 0 {
		table = table.Where("role_id = ?", r.RoleID)
	}
	if r.RoleName != "" {
		table = table.Where("role_name = ?", r.RoleName)
	}
	if err = table.Order("role_sort").Find(&SysRole).Error; err != nil {
		return
	}

	return
}

// Insert ..
func (r *SysRole) Insert() (id int, err error) {
	r.UpdateBy = ""
	existRole := SysRole{}
	existRole.RoleKey = r.RoleKey
	_, existErr := existRole.Get()
	if existErr == nil {
		err = errors.New("权限字符不允许重复，请修改后重试")
		return
	}
	result := db.Eloquent.Table("sys_role").Create(&r)
	if result.Error != nil {
		err = result.Error
		return
	}
	id = r.RoleID
	return
}

// Update ..修改
func (r *SysRole) Update(id int) (update SysRole, err error) {
	if err = db.Eloquent.Table("sys_role").First(&update, id).Error; err != nil {
		return
	}

	//参数1:是要修改的数据
	//参数2:是修改的数据
	if err = db.Eloquent.Table("sys_role").Model(&update).Updates(&r).Error; err != nil {
		return
	}
	return
}

// BatchDelete 批量删除
func (r *SysRole) BatchDelete(id []int) (Result bool, err error) {
	if err = db.Eloquent.Table("sys_role").Where("role_id in (?)", id).Delete(&SysRole{}).Error; err != nil {
		return
	}
	Result = true
	return
}

// MenuIDList  ..
type MenuIDList struct {
	MenuID int `json:"menuId"`
}

// GetRoleMeunID 获取角色对应的菜单ids
func (r *SysRole) GetRoleMeunID() ([]int, error) {
	menuIds := make([]int, 0)
	menuList := make([]MenuIDList, 0)
	if err := db.Eloquent.Table("sys_role_menu").Select("sys_role_menu.menu_id").Joins("LEFT JOIN sys_menu on sys_menu.menu_id=sys_role_menu.menu_id").Where("role_id = ? ", r.RoleID).Where(" sys_role_menu.menu_id not in(select sys_menu.parent_id from sys_role_menu LEFT JOIN sys_menu on sys_menu.menu_id=sys_role_menu.menu_id where role_id =? )", r.RoleID).Find(&menuList).Error; err != nil {
		return nil, err
	}

	for i := 0; i < len(menuList); i++ {
		menuIds = append(menuIds, menuList[i].MenuID)
	}
	return menuIds, nil
}
