package menu

import (
	"errors"
	"fmt"

	"github.com/colynn/pontus/internal/db"
)

// InitPaths ..
func InitPaths(menu *Menu) (err error) {
	parentMenu := new(Menu)
	if int(menu.ParentID) != 0 {
		db.Eloquent.Table("sys_menu").Where("menu_id = ?", menu.ParentID).First(parentMenu)
		if parentMenu.Paths == "" {
			err = errors.New("父级paths异常，请尝试对当前节点父级菜单进行更新操作！")
			return
		}
		menu.Paths = parentMenu.Paths + fmt.Sprintf("/%v", menu.MenuID)
	} else {
		menu.Paths = fmt.Sprintf("/0/%v", menu.MenuID)
	}
	db.Eloquent.Table("sys_menu").Where("menu_id = ?", menu.MenuID).Update("paths", menu.Paths)
	return
}
