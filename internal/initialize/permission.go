package initialize

import (
	"bufio"
	"os"
	"strings"

	"github.com/colynn/pontus/internal/account/role"
	"github.com/colynn/pontus/internal/menu"
	"github.com/colynn/pontus/tools"

	log "unknwon.dev/clog/v2"
)

// InitPermission  init permission for admin role
func InitPermission() error {
	err := APIInterfaceRegister()
	if err != nil {
		log.Error("api interface register error: %s", err.Error())
		return err
	}

	// Get admin role ID
	roleModel := role.SysRole{}
	roleModel.RoleKey = "admin"
	roleItem, err := roleModel.Get()
	if err != nil {
		return err
	}

	roleID := roleItem.RoleID
	// get all api interface menu ids
	roleMenu := menu.RoleMenu{}
	menuIDs, err := roleMenu.GetIDS()
	if err != nil {
		log.Error("get menuid error: %s", err.Error())
		return err
	}

	// TODO: change to support idempotent
	if _, err = roleMenu.DeleteRoleMenu(1); err != nil {
		return err
	}
	if _, err = roleMenu.Insert(roleID, menuIDs); err != nil {
		return err
	}

	return nil
}

// MenuSubItemForAPI ..
type MenuSubItemForAPI struct {
	menu  menu.Menu
	Items []menu.Menu
}

// MenuSubItem ..
type MenuSubItem struct {
	MenuType   string
	Title      string
	Visible    string
	Permission string
	Component  string
	Icon       string
	Path       string
	Redirect   string
	MenuName   string
	Sort       int
	Items      []MenuSubItem
}

// MenuItem ..
type MenuItem struct {
	menu  menu.Menu
	Items []MenuSubItem
}

// APIInterfaceRegister ..
func APIInterfaceRegister() error {
	// 角色
	menuModel := menu.Menu{}
	menuModel.Title = "接口权限"
	menuModel.MenuType = menu.Directory
	menuItem, err := menuModel.GetMenuByArgs()
	var menuParentID int
	if err != nil {
		log.Error("get menu item error: %s, then recreate ", err.Error())
		menuModel.ParentID = 0
		menuModel.Visible = "1"
		menuModel.IsFrame = "1"
		menuID, err := menuModel.Create()
		if err != nil {
			log.Error("create menu error: %s", err.Error())
			return err
		}
		menuParentID = menuID

	} else {
		menuParentID = menuItem.MenuID
	}

	// generate APIItems
	for _, apiItem := range APIItems {
		itemModel := menu.Menu{
			Title:      apiItem.menu.Title,
			Visible:    apiItem.menu.Visible,
			MenuType:   apiItem.menu.MenuType,
			Permission: apiItem.menu.Permission,
			ParentID:   menuParentID,
		}
		var subParentID int
		existItem, err := itemModel.GetMenuByArgs()
		if err == nil {
			subParentID = existItem.MenuID
		} else {
			subParentID, err = itemModel.Create()
			if err != nil {
				log.Error("create item: %v error: %s", itemModel, err.Error())
				continue
			}
		}

		for _, subItem := range apiItem.Items {
			subItemModel := menu.Menu{
				Title:      subItem.Title,
				Visible:    subItem.Visible,
				MenuType:   subItem.MenuType,
				Permission: subItem.Permission,
				ParentID:   subParentID,
				Path:       subItem.Path,
				Action:     subItem.Action,
			}
			_, err := subItemModel.GetMenuByArgs()
			if err != nil {
				log.Error("titel: %v get error: %s", subItem.Title, err.Error())
				// create
				_, err := subItemModel.Create()
				if err != nil {
					log.Error("create item: %v error: %s", subItemModel, err.Error())
					continue
				}
			}
		}
	}
	return nil
}

// generateAPIItems based on router.go generate APIItems
func generateAPIItems(routerGroup, urlPrefix string) []MenuItem {
	// TODO: need migrate router into docker run-env
	routerPath := tools.EnsureAbs("internal/router/router.go")
	filePr, err := os.Open(routerPath)
	if err != nil {
		log.Error("oper 'router/router.go' error: %s", err.Error())
		return nil
	}
	scanner := bufio.NewScanner(filePr)
	for scanner.Scan() {
		line := scanner.Text()
		lineItems := strings.Split(line, "// ")
		log.Trace("%v, %d", lineItems, len(lineItems))
	}

	return nil
}

// APIItems declare menu api item
var APIItems = []MenuSubItemForAPI{
	{
		menu: menu.Menu{
			MenuType:   menu.MenuType,
			Title:      "个人中心",
			Visible:    "1",
			Permission: "system:profile:index",
		},
		Items: []menu.Menu{
			{
				MenuType:   menu.APIInterface,
				Title:      "修改用户",
				Visible:    "1",
				Permission: "system:profile:edit",
				Path:       "/api/v1/user/profile",
				Action:     "PUT",
			},
		},
	},
	{
		menu: menu.Menu{
			MenuType:   menu.MenuType,
			Title:      "终端设备",
			Visible:    "1",
			Permission: "asset:pc:index",
		},
		Items: []menu.Menu{
			{
				MenuType:   menu.APIInterface,
				Title:      "终端导入",
				Visible:    "1",
				Permission: "asset:pc:import",
				Path:       "/api/v1/assets/pc/upload",
				Action:     "POST",
			},
			{
				MenuType:   menu.APIInterface,
				Title:      "编辑设备",
				Visible:    "1",
				Permission: "asset:pc:edit",
				Path:       "/api/v1/assets/tangibles/:instanceID",
				Action:     "PUT",
			},
			{
				MenuType:   menu.APIInterface,
				Title:      "删除设备",
				Visible:    "1",
				Permission: "asset:pc:remove",
				Path:       "/api/v1/assets/tangibles/:instanceID",
				Action:     "DELETE",
			},
		},
	},
	{
		menu: menu.Menu{
			MenuType:   menu.MenuType,
			Title:      "用户管理",
			Visible:    "1",
			Permission: "system:user:list",
		},
		Items: []menu.Menu{
			{
				MenuType:   menu.APIInterface,
				Title:      "修改用户",
				Visible:    "1",
				Permission: "system:user:edit",
				Path:       "/api/v1/users/:userID",
				Action:     "PUT",
			},
			{
				MenuType:   menu.APIInterface,
				Title:      "用户导入",
				Visible:    "1",
				Permission: "system:user:import",
				Path:       "/api/v1/user/upload",
				Action:     "POST",
			},
			{
				MenuType:   menu.APIInterface,
				Title:      "删除用户",
				Visible:    "1",
				Permission: "system:user:remove",
				Path:       "/api/v1/users/:userID",
				Action:     "DELETE",
			},
			{
				MenuType:   menu.APIInterface,
				Title:      "新增用户",
				Visible:    "1",
				Permission: "system:user:add",
				Path:       "/api/v1/users",
				Action:     "POST",
			},
		},
	},
	{
		menu: menu.Menu{
			MenuType:   menu.MenuType,
			Title:      "角色管理",
			Visible:    "1",
			Permission: "system:role:list",
		},
		Items: []menu.Menu{
			{
				MenuType:   menu.APIInterface,
				Title:      "新增角色",
				Visible:    "1",
				Permission: "system:role:add",
				Path:       "/api/v1/roles",
				Action:     "POST",
			},
			{
				MenuType:   menu.APIInterface,
				Title:      "更新角色",
				Visible:    "1",
				Permission: "system:role:edit",
				Path:       "/api/v1/roles/:roleID",
				Action:     "PUT",
			},
			{
				MenuType:   menu.APIInterface,
				Title:      "删除角色",
				Visible:    "1",
				Permission: "system:role:remove",
				Path:       "/api/v1/roles/:roleID",
				Action:     "DELETE",
			},
		},
	},
	{
		menu: menu.Menu{
			MenuType:   menu.MenuType,
			Title:      "菜单管理",
			Visible:    "1",
			Permission: "system:role:list",
		},
		Items: []menu.Menu{
			{
				MenuType:   menu.APIInterface,
				Title:      "新增菜单",
				Visible:    "1",
				Permission: "system:menu:add",
				Path:       "/api/v1/menus",
				Action:     "POST",
			},
			{
				MenuType:   menu.APIInterface,
				Title:      "更新菜单",
				Visible:    "1",
				Permission: "system:menu:edit",
				Path:       "/api/v1/menus/:id",
				Action:     "PUT",
			},
			{
				MenuType:   menu.APIInterface,
				Title:      "删除菜单",
				Visible:    "1",
				Permission: "system:menu:remove",
				Path:       "/api/v1/menus/:id",
				Action:     "DELETE",
			},
		},
	},
}
