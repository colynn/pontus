package menu

import (
	"github.com/colynn/pontus/internal/db"

	"github.com/jinzhu/gorm"
	log "unknwon.dev/clog/v2"
)

// MenuType
const (
	Button       = "F" // 按钮
	MenuType     = "C" // 菜单
	Directory    = "M" // 目录
	APIInterface = "A" // 接口
)

// Menu ..
type Menu struct {
	MenuID     int     `json:"menuID" gorm:"primary_key;AUTO_INCREMENT"`
	MenuName   string  `json:"menuName" gorm:"type:varchar(64);"`
	Title      string  `json:"title" gorm:"type:varchar(64);"`
	Icon       string  `json:"icon" gorm:"type:varchar(128);"`
	Path       string  `json:"path" gorm:"type:varchar(128);"`
	Redirect   string  `json:"redirect" gorm:"type:varchar(128);"`
	Paths      string  `json:"paths" gorm:"type:varchar(128);"`
	MenuType   string  `json:"menuType" gorm:"type:varchar(1);"`
	Action     string  `json:"action" gorm:"type:varchar(16);"`
	Permission string  `json:"permission" gorm:"type:varchar(32);"`
	ParentID   int     `json:"parentID" gorm:"type:int(11);"`
	NoCache    bool    `json:"noCache" gorm:"type:char(1);"`
	Breadcrumb string  `json:"breadcrumb" gorm:"type:varchar(255);"`
	Component  string  `json:"component" gorm:"type:varchar(255);"`
	Sort       int     `json:"sort" gorm:"type:int(4);"`
	Visible    string  `json:"visible" gorm:"type:char(1);"`
	CreateBy   string  `json:"createBy" gorm:"type:varchar(128);"`
	UpdateBy   string  `json:"updateBy" gorm:"type:varchar(128);"`
	IsFrame    string  `json:"isFrame" gorm:"type:int(1);DEFAULT:0;"`
	DataScope  string  `json:"dataScope" gorm:"-"`
	Params     string  `json:"params" gorm:"-"`
	Children   []*Menu `json:"children" gorm:"-"`
	IsSelect   bool    `json:"is_select" gorm:"-"`
	db.BaseModel
}

// TableName ..
func (Menu) TableName() string {
	return "sys_menu"
}

// Lable ..
type Lable struct {
	ID       int     `json:"id" gorm:"-"`
	Label    string  `json:"label" gorm:"-"`
	Children []Lable `json:"children" gorm:"-"`
}

// Role ..
type Role struct {
	Menu
	IsSelect bool `json:"is_select" gorm:"-"`
}

// SetMenuLable ..
func (e *Menu) SetMenuLable() (m []Lable, err error) {
	menulist, err := e.Get()

	m = make([]Lable, 0)
	for i := 0; i < len(menulist); i++ {
		if menulist[i].ParentID != 0 {
			continue
		}
		e := Lable{}
		e.ID = menulist[i].MenuID
		e.Label = menulist[i].Title
		menusInfo := renderMenuLableTreeData(menulist, e)

		m = append(m, menusInfo)
	}
	return
}

// Get ..
func (e *Menu) Get() (Menus []*Menu, err error) {
	table := e.GetFilter()
	if err = table.Order("sort").Find(&Menus).Error; err != nil {
		return nil, err
	}
	return
}

// GetFilter ..
func (e *Menu) GetFilter() (table *gorm.DB) {
	table = db.Eloquent.Table(e.TableName()).Where("deleted_at is NULL")
	if e.MenuName != "" {
		table = table.Where("menu_name = ?", e.MenuName)
	}
	if e.Title != "" {
		table = table.Where("title = ?", e.Title)
	}
	if e.Path != "" {
		table = table.Where("path = ?", e.Path)
	}
	if e.Action != "" {
		table = table.Where("action = ?", e.Action)
	}
	if e.MenuType != "" {
		table = table.Where("menu_type = ?", e.MenuType)
	}

	if e.Visible != "" {
		table = table.Where("visible = ?", e.Visible)
	}
	// comment tmp
	// if e.Redirect != "" {
	// 	table = table.Where("redirect = ?", e.Redirect)
	// }
	return
}

// GetMenuList ..
func (e *Menu) GetMenuList() (menus []*Menu, err error) {
	menulist, err := e.Get()

	m := make([]*Menu, 0, len(menulist))
	for i := 0; i < len(menulist); i++ {
		if menulist[i].ParentID != 0 {
			continue
		}
		menusInfo := renderMenuTree(menulist, menulist[i])

		m = append(m, menusInfo)
	}
	return m, nil
}

// GetMenuByArgs ..
func (e *Menu) GetMenuByArgs() (menu Menu, err error) {
	log.Trace("GetMenuByTitleAndType title: %s, type: %s", e.Title, e.MenuType)
	table := e.GetFilter()
	if err = table.Find(&menu).Error; err != nil {
		return
	}
	return
}

// GetByMenuID ..
func (e *Menu) GetByMenuID() (menu Menu, err error) {
	query := db.Eloquent.Table(e.TableName())
	if err = query.Where("menu_id = ?", e.MenuID).Find(&menu).Error; err != nil {
		return
	}
	return
}

// Create ..
func (e *Menu) Create() (id int, err error) {
	result := db.Eloquent.Table(e.TableName()).Create(&e)
	if result.Error != nil {
		err = result.Error
		return
	}
	err = InitPaths(e)
	if err != nil {
		return
	}
	id = e.MenuID
	return
}

// Update ..
func (e *Menu) Update(id int) (update Menu, err error) {
	if err = db.Eloquent.Table("sys_menu").First(&update, id).Error; err != nil {
		return
	}

	//参数1:是要修改的数据
	//参数2:是修改的数据
	log.Trace("permission: " + e.Permission)
	if err = db.Eloquent.Table("sys_menu").Model(&update).Update(&e).Error; err != nil {
		log.Error("udpate menu item error: %s", err.Error())
		return
	}
	// gorm can not update 0 value
	if e.Redirect == "" || e.Path == "" || e.Permission == "" {
		if err = db.Eloquent.Table(e.TableName()).Save(&e).Error; err != nil {
			return
		}
	}
	return
}

// Delete ..
func (e *Menu) Delete(id int) (success bool, err error) {
	if err = db.Eloquent.Table(e.TableName()).Where("menu_id = ?", id).Delete(&Menu{}).Error; err != nil {
		success = false
		return
	}
	success = true
	return
}

// GetMenuListByRole ..
func (e *Menu) GetMenuListByRole(roleName string) ([]*Menu, error) {
	menulist, err := e.GetByRoleName(roleName)
	if err != nil {
		log.Error("get menu list by role occur errror: %s", err.Error())
	}
	m := make([]*Menu, 0)
	for i := 0; i < len(menulist); i++ {
		if menulist[i].ParentID != 0 {
			continue
		}
		menusInfo := renderMenuTree(menulist, menulist[i])

		m = append(m, menusInfo)
	}
	return m, nil
}

// GetByRoleName ..
func (e *Menu) GetByRoleName(rolename string) (Menus []*Menu, err error) {
	table := db.Eloquent.Table(e.TableName()).Select("sys_menu.*").Joins("left join sys_role_menu on sys_role_menu.menu_id=sys_menu.menu_id")
	table = table.Where("sys_role_menu.role_name=? and menu_type in (?,?)", rolename, Directory, MenuType)
	if err = table.Order("sort").Find(&Menus).Error; err != nil {
		return
	}
	log.Trace("menus length: %v", len(Menus))
	return
}

// renderMenuLableTreeData ..
func renderMenuLableTreeData(menulist []*Menu, menu Lable) Lable {
	list := menulist

	min := make([]Lable, 0)
	for j := 0; j < len(list); j++ {

		if menu.ID != list[j].ParentID {
			continue
		}
		mi := Lable{}
		mi.ID = list[j].MenuID
		mi.Label = list[j].Title
		mi.Children = []Lable{}
		if list[j].MenuType != "F" {
			ms := renderMenuLableTreeData(menulist, mi)
			min = append(min, ms)
		} else {
			min = append(min, mi)
		}

	}
	menu.Children = min
	return menu
}

func renderMenuTree(menulist []*Menu, menu *Menu) *Menu {
	list := menulist

	min := make([]*Menu, 0)
	for j := 0; j < len(list); j++ {

		if menu.MenuID != list[j].ParentID {
			continue
		}
		mi := &Menu{}
		mi.CreatedAt = list[j].CreatedAt
		mi.UpdatedAt = list[j].UpdatedAt
		mi.MenuID = list[j].MenuID
		mi.MenuName = list[j].MenuName
		mi.Title = list[j].Title
		mi.Icon = list[j].Icon
		mi.Path = list[j].Path
		mi.Redirect = list[j].Redirect
		mi.MenuType = list[j].MenuType
		mi.Action = list[j].Action
		mi.Permission = list[j].Permission
		mi.ParentID = list[j].ParentID
		mi.NoCache = list[j].NoCache
		mi.Breadcrumb = list[j].Breadcrumb
		mi.Component = list[j].Component
		mi.Sort = list[j].Sort
		mi.Visible = list[j].Visible
		mi.Children = []*Menu{}

		if mi.MenuType != "F" {
			ms := renderMenuTree(menulist, mi)
			min = append(min, ms)

		} else {
			min = append(min, mi)
		}

	}
	menu.Children = min
	return menu
}
