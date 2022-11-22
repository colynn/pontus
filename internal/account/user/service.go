package user

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/colynn/pontus/config"
	"github.com/colynn/pontus/internal/account/role"
	orm "github.com/colynn/pontus/internal/db"

	"github.com/colynn/go-ldap-client"
	log "unknwon.dev/clog/v2"
)

// Service ..
type Service struct {
	UserDao *SysUser
}

// NewService ..
func NewService() *Service {
	return &Service{
		UserDao: &SysUser{},
	}
}

// Authenticate ..
func (u *Service) Authenticate(loginreq *Login) (userInfo *UserInfo, e error) {
	c := config.GetConfig()
	client := &ldap.Client{
		Base:               c.GetString("ldap.baseDN"),
		Host:               c.GetString("ldap.host"),
		Port:               c.GetInt("ldap.port"),
		UseSSL:             false,
		BindDN:             c.GetString("ldap.bindDN"),
		BindPassword:       c.GetString("ldap.bindPassword"),
		UserFilter:         c.GetString("ldap.userFilter"),
		GroupFilter:        "(name=%s)",
		Attributes:         []string{"givenName", "sn", "mail", "uid"},
		SkipTLS:            true,
		InsecureSkipVerify: true,
	}
	defer client.Close()

	authVerify, userResp, err := client.Authenticate(loginreq.Username, loginreq.Password)
	if !authVerify {
		if err != nil {
			log.Error("authVerify error: %s", err.Error())
		}
		return nil, fmt.Errorf("域账号或密码错误")
	}
	bindDN := userResp["dn"]
	bindDNSplit := strings.Split(bindDN, ",OU=people")
	userName := ""
	if len(bindDNSplit) > 0 {
		cnSplit := strings.Split(bindDNSplit[0], "CN=")
		if len(cnSplit) == 2 {
			userName = cnSplit[1]
		}
	}
	if userName == "" {
		return nil, fmt.Errorf("域帐号服务器返回信息格式有误")
	}
	userInfo = &UserInfo{
		Name:  userName,
		Email: userResp["mail"],
	}
	return userInfo, nil
}

// Exists ..
func (u *Service) Exists(loginreq *Login) (exist bool, user *SysUser, e error) {
	item := SysUser{}
	e = orm.Eloquent.Table("sys_user").Where("username = ?", loginreq.Username).Find(&item).Error
	if e != nil {
		log.Warn("when check user occur error: %s", e.Error())
		if strings.Contains(e.Error(), "record not found") {
			return false, nil, nil
		}
		return false, nil, e
	}
	user = &item
	return true, user, nil
}

// InsertOrUpdateUserData ..
func (u *Service) InsertOrUpdateUserData(filePath string) error {
	items, err := parseUploadFile(filePath)
	if err != nil {
		log.Error("when inser update terminal device, parse xls error: %s", err.Error())
		return err
	}
	syncPool := make(chan bool, 50)
	defaultRoleID := 0
	if len(items) > 0 {
		roleSvc := role.NewService()
		// TODO: get defautl key from config
		roleItem, err := roleSvc.GetRoleItemByKey("default")
		if err != nil {
			log.Error("when create user get role by role key error: %s", err.Error())
			return err
		}
		defaultRoleID = roleItem.RoleID
	}
	for _, item := range items {
		go func(item *ImportItem, roleID int) {
			syncPool <- true
			err := u.InsertOrUpdateUserOneItem(item, defaultRoleID)
			if err != nil {
				log.Warn("inster or update item error: %s", err.Error())
			}
			<-syncPool
		}(item, defaultRoleID)
	}
	return nil
}

// InsertOrUpdateUserOneItem ..
func (u *Service) InsertOrUpdateUserOneItem(item *ImportItem, roleID int) error {
	user := SysUser{
		SysUserB: SysUserB{
			Email:    strings.TrimSpace(item.Email),
			RealName: strings.TrimSpace(item.DisplayName),
			RoleID:   roleID,
		},
		LoginM: LoginM{
			Username: strings.TrimSpace(item.Username),
		},
	}
	if _, err := u.UserDao.Insert(&user); err != nil && !strings.Contains(err.Error(), "账户已存在") {
		return err
	}
	return nil
}

// GetUserIDByRealName ..
func (u *Service) GetUserIDByRealName(displayName string) int {
	id, err := u.UserDao.GetUserIDByRealName(displayName)
	if err != nil {
		log.Error("when get use id by realname: %v error: %s, skip update userid", displayName, err.Error())
	}
	return id
}
func parseUploadFile(filepath string) ([]*ImportItem, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(file)
	var line string
	var items []*ImportItem
	for {
		line, err = reader.ReadString('\n')
		if err != nil && err != io.EOF {
			break
		}

		// Process the line here.
		lineSplit := strings.Split(line, ",")
		if len(lineSplit) < 2 {
			// log.Warn("user import line: %+v format parse error, skip", lineSplit)
			if err == io.EOF {
				break
			} else {
				continue
			}
		}
		item := ImportItem{}
		if !strings.Contains(lineSplit[1], ".") {
			if err == io.EOF {
				break
			} else {
				continue
			}
		}
		item.DisplayName = lineSplit[0]
		item.Username = lineSplit[1]
		if len(lineSplit) == 3 {
			item.Email = lineSplit[2]
		}
		items = append(items, &item)

		if err == io.EOF {
			break
		}
	}
	if err != io.EOF {
		fmt.Printf(" > Failed with error: %v\n", err)
		return nil, err
	}
	return items, nil
}
