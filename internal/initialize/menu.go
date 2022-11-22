package initialize

import (
	"github.com/colynn/pontus/internal/menu"

	log "unknwon.dev/clog/v2"
)

// InitMenu  init permission for admin role
func InitMenu() error {
	err := MenuRegister()
	if err != nil {
		log.Error("sub menu register error: %s", err.Error())
		return err
	}

	return nil
}

// MenuRegister ..
func MenuRegister() error {

	//
	MenuItems = append(MenuItems, CMDBMenuItems...)
	// 角色
	for _, item := range MenuItems {
		menuModel := menu.Menu{}
		menuModel.Title = item.menu.Title
		menuModel.MenuType = item.menu.MenuType
		menuModel.Redirect = item.menu.Redirect
		menuModel.Sort = item.menu.Sort
		menuItem, err := menuModel.GetMenuByArgs()
		var menuParentID int
		if err != nil {
			log.Warn("get menu item error: %s, then recreate ", err.Error())
			menuModel.ParentID = item.menu.ParentID
			menuModel.Visible = item.menu.Visible
			menuModel.IsFrame = item.menu.IsFrame
			menuModel.Path = item.menu.Path
			menuModel.Redirect = item.menu.Redirect
			menuModel.MenuName = item.menu.MenuName
			menuModel.Icon = item.menu.Icon
			menuItem.Sort = item.menu.Sort
			menuModel.Component = item.menu.Component
			menuID, err := menuModel.Create()
			if err != nil {
				log.Error("create menu %s error: %s", menuModel.MenuName, err.Error())
				return err
			}
			menuParentID = menuID

		} else {
			menuParentID = menuItem.MenuID
		}

		subMenuItems := item.Items
		// generate SubMentItems
		createSubItem(menuParentID, subMenuItems)
	}
	return nil
}

func createSubItem(menuParentID int, items []MenuSubItem) {
	// generate SubMentItems
	for _, apiItem := range items {
		itemModel := menu.Menu{
			Title:      apiItem.Title,
			Icon:       apiItem.Icon,
			Visible:    apiItem.Visible,
			MenuType:   apiItem.MenuType,
			MenuName:   apiItem.MenuName,
			Permission: apiItem.Permission,
			Component:  apiItem.Component,
			Path:       apiItem.Path,
			Redirect:   apiItem.Redirect,
			ParentID:   menuParentID,
			Sort:       apiItem.Sort,
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
			log.Trace("%s, %s re-create success", itemModel.MenuName, itemModel.Path)
		}
		if len(apiItem.Items) > 0 {
			createSubItem(subParentID, apiItem.Items)
		}
	}
}

// MenuItems ..
var MenuItems = []MenuItem{
	{
		menu: menu.Menu{
			MenuType:  menu.Directory,
			ParentID:  0,
			Title:     "系统管理",
			Visible:   "0",
			Sort:      90,
			IsFrame:   "1",
			MenuName:  "sysadmin",
			Icon:      "password",
			Component: "Layout",
			Path:      "/sysadmin",
		},
		Items: []MenuSubItem{
			{
				MenuType:   menu.MenuType,
				Title:      "角色管理",
				Visible:    "0",
				Sort:       3,
				Permission: "system:role:list",
				Component:  "/sysadmin/role/index",
				Path:       "role",
				MenuName:   "role",
				Items: []MenuSubItem{
					{
						MenuType:   menu.Button,
						Title:      "修改角色",
						Permission: "system:role:edit",
					},
					{
						MenuType:   menu.Button,
						Title:      "删除角色",
						Permission: "system:role:remove",
					},
					{
						MenuType:   menu.Button,
						Title:      "新增角色",
						Permission: "system:role:add",
					},
				},
			},
			{
				MenuType:   menu.MenuType,
				Title:      "用户管理",
				Visible:    "0",
				Sort:       2,
				Permission: "system:user:list",
				Component:  "/sysadmin/user/index",
				Path:       "user",
				MenuName:   "user",
				Items: []MenuSubItem{
					{
						MenuType:   menu.Button,
						Title:      "修改用户",
						Permission: "system:user:edit",
					},
					{
						MenuType:   menu.Button,
						Title:      "用户导入",
						Permission: "system:user:import",
					},
					{
						MenuType:   menu.Button,
						Title:      "删除用户",
						Permission: "system:user:remove",
					},
					{
						MenuType:   menu.Button,
						Title:      "新增用户",
						Permission: "system:user:add",
					},
				},
			},
			{
				MenuType:   menu.MenuType,
				Title:      "菜单管理",
				Visible:    "0",
				Sort:       4,
				Permission: "system:menu:list",
				Component:  "/sysadmin/menu/index",
				Path:       "menu",
				Items: []MenuSubItem{
					{
						MenuType:   menu.Button,
						Title:      "修改菜单",
						Permission: "system:menu:edit",
					},
					{
						MenuType:   menu.Button,
						Title:      "删除菜单",
						Permission: "system:menu:remove",
					},
					{
						MenuType:   menu.Button,
						Title:      "新增菜单",
						Permission: "system:menu:add",
					},
				},
			},
		},
	},
}

// CMDBMenuItems ..
var CMDBMenuItems = []MenuItem{
	{
		menu: menu.Menu{
			MenuType:  menu.Directory,
			ParentID:  0,
			Title:     "资产管理",
			Visible:   "0",
			IsFrame:   "1",
			Sort:      1,
			Redirect:  "/asset/index",
			MenuName:  "asset",
			Icon:      "asset",
			Component: "Layout",
			Path:      "/asset",
		},
		Items: []MenuSubItem{
			{
				MenuType:   menu.MenuType,
				Title:      "资产汇总",
				Visible:    "0",
				Sort:       1,
				Permission: "asset:overview:info",
				Component:  "/asset/index",
				Path:       "index",
				MenuName:   "assetInfo",
				Icon:       "form",
				Items:      []MenuSubItem{},
			},
			{
				MenuType:   menu.MenuType,
				Title:      "资产列表",
				Visible:    "0",
				Sort:       3,
				Permission: "asset:overview:list",
				Component:  "/asset/list",
				Redirect:   "/asset/list/physical",
				Path:       "list",
				MenuName:   "资产列表",
				Icon:       "example",
				Items: []MenuSubItem{
					{
						MenuType:   menu.MenuType,
						Title:      "物理设备",
						Visible:    "0",
						Sort:       2,
						Permission: "asset:physical:list",
						Component:  "/asset/physical/index",
						Path:       "physical",
						MenuName:   "物理设备",
						Items:      []MenuSubItem{},
					},
					// -- physical start
					{
						MenuType:   menu.MenuType,
						Title:      "物理机详情",
						Visible:    "1",
						Sort:       2,
						Permission: "asset:physical:detail",
						Component:  "/asset/physical/detail",
						Path:       "physical/detail/:instanceID",
						MenuName:   "physicalDetail",
						Items:      []MenuSubItem{},
					},
					// -- physical end
					{
						MenuType:   menu.MenuType,
						Title:      "终端设备",
						Visible:    "0",
						Sort:       3,
						Permission: "asset:pc:list",
						Component:  "/asset/pc/list",
						Path:       "pc",
						MenuName:   "终端列表",
						Items:      []MenuSubItem{},
					},
					// -- pc detail start
					{
						MenuType:   menu.Button,
						Title:      "终端导入",
						Permission: "asset:pc:import",
						MenuName:   "终端导入",
						Items:      []MenuSubItem{},
					},
					{
						MenuType:   menu.MenuType,
						Title:      "终端详情",
						Visible:    "1",
						Sort:       2,
						Permission: "asset:pc:detail",
						Component:  "/asset/pc/detail",
						Path:       "pc/detail/:instanceID",
						MenuName:   "pcDetail",
						Items:      []MenuSubItem{},
					},
					{
						MenuType:   menu.Button,
						Title:      "修改终端信息",
						Permission: "asset:pc:edit",
						MenuName:   "修改终端信息",
						Items:      []MenuSubItem{},
					},
					{
						MenuType:   menu.Button,
						Title:      "删除终端信息",
						MenuName:   "删除终端信息",
						Permission: "asset:pc:remove",
						Items:      []MenuSubItem{},
					},
					// -- pc detail end
				},
			},
		},
	},
}
