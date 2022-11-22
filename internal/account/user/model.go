package user

import (
	"errors"
	"log"
	"strings"
	"time"

	db "github.com/colynn/pontus/internal/db"
	"github.com/colynn/pontus/internal/pkg/customerror"
	"github.com/colynn/pontus/internal/pkg/password"
)

// UserName ..
type UserName struct {
	Username string `gorm:"type:varchar(64)" json:"username"`
}

type PassWord struct {
	// 密码
	Password string `gorm:"type:varchar(128)" json:"password"`
}

type LoginM struct {
	Username string `gorm:"type:varchar(64)" json:"username"`
	Password string `gorm:"type:varchar(128)" json:"password"`
}

type SysUserId struct {
	ID int `gorm:"primary_key;AUTO_INCREMENT"  json:"userId"` // 编码
}

// SysUserB ..
type SysUserB struct {
	RealName string `gorm:"type:varchar(128)" json:"realName"` // 昵称
	Rule     string `gorm:"type:varchar(16); null" json:"rule"`
	Phone    string `gorm:"type:varchar(18)" json:"phone"`   // 昵称
	RoleID   int    `gorm:"type:int(11)" json:"roleId"`      // 角色编码
	Salt     string `gorm:"type:varchar(255)" json:"salt"`   //盐
	Avatar   string `gorm:"type:varchar(255)" json:"avatar"` //头像
	Email    string `gorm:"type:varchar(128)" json:"email"`  //邮箱
	Remark   string `gorm:"type:varchar(255)" json:"remark"` //备注
	db.BaseModel
}

type SysUser struct {
	SysUserId
	SysUserB
	LoginM
}

func (SysUser) TableName() string {
	return "sys_user"
}

// Login model defined
type Login struct {
	Username string `form:"UserName" json:"username" binding:"required"`
	Password string `form:"Password" json:"password" binding:"required"`
	// TODO: comment code tmp
	// Code     string `form:"Code" json:"code" binding:"required"`
	Code string `form:"Code" json:"code"`
	UUID string `form:"UUID" json:"uuid" binding:"required"`
}

// LoginLog ..
type LoginLog struct {
	InfoId        int       `json:"infoId" gorm:"primary_key;AUTO_INCREMENT"` //主键
	Username      string    `json:"username" gorm:"type:varchar(128);"`       //用户名
	Status        string    `json:"status" gorm:"type:int(1);"`               //状态
	Ipaddr        string    `json:"ipaddr" gorm:"type:varchar(255);"`         //ip地址
	LoginLocation string    `json:"loginLocation" gorm:"type:varchar(255);"`  //归属地
	Browser       string    `json:"browser" gorm:"type:varchar(255);"`        //浏览器
	Os            string    `json:"os" gorm:"type:varchar(255);"`             //系统
	Platform      string    `json:"platform" gorm:"type:varchar(255);"`       // 固件
	LoginTime     time.Time `json:"loginTime" gorm:"type:timestamp;"`         //登录时间
	CreateBy      string    `json:"createBy" gorm:"type:varchar(128);"`       //创建人
	UpdateBy      string    `json:"updateBy" gorm:"type:varchar(128);"`       //更新者
	DataScope     string    `json:"dataScope" gorm:"-"`                       //数据
	Params        string    `json:"params" gorm:"-"`                          //
	Remark        string    `json:"remark" gorm:"type:varchar(255);"`         //备注
	Msg           string    `json:"msg" gorm:"type:varchar(255);"`
	db.BaseModel
}

func (LoginLog) TableName() string {
	return "sys_loginlog"
}

func (e *LoginLog) Get() (LoginLog, error) {
	var doc LoginLog

	table := db.Eloquent.Table(e.TableName())
	if e.Ipaddr != "" {
		table = table.Where("ipaddr = ?", e.Ipaddr)
	}
	if e.InfoId != 0 {
		table = table.Where("info_id = ?", e.InfoId)
	}

	if err := table.First(&doc).Error; err != nil {
		return doc, err
	}
	return doc, nil
}

func (e *LoginLog) Create() (LoginLog, error) {
	var doc LoginLog
	e.CreateBy = "0"
	e.UpdateBy = "0"
	result := db.Eloquent.Table(e.TableName()).Create(&e)
	if result.Error != nil {
		err := result.Error
		return doc, err
	}
	doc = *e
	return doc, nil
}

// GetUser ..
func (u *Login) GetUser() (user SysUser, e error) {

	e = db.Eloquent.Table("sys_user").Where("username = ? ", u.Username).Find(&user).Error
	if e != nil {
		if strings.Contains(e.Error(), "record not found") {
			customerror.HasError(e, "账号或密码错误(代码204)", 500)
		}
		log.Print(e)
		return
	}
	_, e = password.CompareHashAndPassword(user.Password, u.Password)
	if e != nil {
		if strings.Contains(e.Error(), "hashedPassword is not the hash of the given password") {
			customerror.HasError(e, "账号或密码错误(代码201)", 500)
		}
		log.Print(e)
		return
	}
	return user, nil
}

// Get 获取用户数据
func (e *SysUser) Get(userID int) (user *SysUser, err error) {
	item := SysUser{}
	table := db.Eloquent.Table(e.TableName()).Select("sys_user.*")
	table = table.Joins("left join sys_role on sys_user.role_id=sys_role.role_id")
	if userID != 0 {
		table = table.Where("id = ?", userID)
	}

	if e.Username != "" {
		table = table.Where("username = ?", e.Username)
	}

	if e.RoleID != 0 {
		table = table.Where("role_id = ?", e.RoleID)
	}

	if err = table.First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// GetUserIDByRealName ..
func (e *SysUser) GetUserIDByRealName(realName string) (id int, err error) {
	item := SysUser{}
	table := db.Eloquent.Table(e.TableName()).Select("sys_user.*")
	if realName != "" {
		table = table.Where("real_name = ?", realName)
	} else {
		return
	}
	err = table.First(&item).Error
	id = item.ID
	return
}

// GetByRoleID 获取用户数据
func (e *SysUser) GetByRoleID(roleID int) (user []SysUser, err error) {
	item := []SysUser{}
	table := db.Eloquent.Table(e.TableName()).Select("sys_user.*")

	if roleID != 0 {
		table = table.Where("role_id = ?", roleID)
	}

	if err = table.Find(&item).Error; err != nil {
		return nil, err
	}
	return item, nil
}

// GetPage ..
func (e *SysUser) GetPage(pageSize int, pageIndex int) ([]SysUserPage, int, error) {
	var doc []SysUserPage

	table := db.Eloquent.Select("sys_user.*").Table(e.TableName()).Where("deleted_at is NULL")

	if e.Username != "" {
		table = table.Where("username like ?", "%"+e.Username+"%")
	}

	var count int

	if err := table.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&doc).Error; err != nil {
		return nil, 0, err
	}
	table.Count(&count)
	return doc, count, nil
}

// Insert .. 添加
func (e *SysUser) Insert(user *SysUser) (id int, err error) {

	// check 用户名
	var count int
	db.Eloquent.Table(e.TableName()).Where("username = ? and deleted_at is NULL", user.Username).Count(&count)
	if count > 0 {
		err = errors.New("账户已存在！")
		return
	}

	if err = db.Eloquent.Table(e.TableName()).Create(&user).Error; err != nil {
		return
	}
	id = user.ID
	return
}

// Update ..修改
func (e *SysUser) Update(id int) (update SysUser, err error) {
	if err = db.Eloquent.Table(e.TableName()).First(&update, id).Error; err != nil {
		return
	}

	// args1:是要修改的数据
	// args2:是修改的数据
	if err = db.Eloquent.Table(e.TableName()).Model(&update).Updates(&e).Error; err != nil {
		return
	}
	return
}

// BatchDelete 批量删除
func (e *SysUser) BatchDelete(id []int) (Result bool, err error) {
	if err = db.Eloquent.Table(e.TableName()).Where("id in (?)", id).Delete(&SysUser{}).Error; err != nil {
		return
	}
	Result = true
	return
}

// GetList ..
func (e *SysUser) GetList() (users []SysUser, err error) {
	err = db.Eloquent.Table(e.TableName()).Find(&users).Error
	return
}
