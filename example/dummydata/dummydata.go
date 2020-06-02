package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"log"
	_const "polestar/auth/const"
	"polestar/auth/server/store"
	"polestar/common/db"
	"polestar/common/model"
	"polestar/common/service"
	"polestar/common/util"
	"time"
)

func main() {
	dbConn := db.NewMySQLConnectionInstance()
	defer func() {
		if dbConn != nil {
			dbConn.Close()
		}
	}()

	// 添加Client
	clientStore := store.NewMySQLClientStoreInstance()

	secret, err := bcrypt.GenerateFromPassword([]byte(util.ApplicationConfig().Auth.Jwt.Secret), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
	} else {
		cli := &store.ClientInfo{
			ClientId:             "democlient",
			ClientSecret:         string(secret),
			GrantType:            []_const.GrantType{_const.GrantTypePasswordCredentials, _const.GrantTypeRefreshToken},
			Scope:                []string{"ums"},
			AccessTokenDuration:  time.Minute * 10,
			RefreshTokenDuration: time.Hour * 24 * 2,
		}
		err = clientStore.AddClient(cli)
		if err != nil {
			log.Println(err.Error())
		} else {
			log.Printf("客户端[%s]插入成功。", cli.ClientId)
		}
	}

	// 添加用户
	user, err := service.NewSysUserService().GetUserByUserName("admin")
	if gorm.IsRecordNotFoundError(err) {
		pwd, err := bcrypt.GenerateFromPassword([]byte("123"), bcrypt.DefaultCost)
		if err != nil {
			log.Println(err)
		} else {
			user := &model.SysUser{
				UserName:    "admin",
				Name:        "管理员",
				Password:    string(pwd),
				HeaderImage: "cat.jpg",
				Comment:     "none",
				Status:      "00",
			}
			createUser, err := service.NewSysUserService().CreateUser(user)
			if err != nil {
				log.Println(err)
			} else {
				fmt.Printf("用户[%s]数据插入成功。\n", createUser.UserName)
			}
		}
	} else if err != nil {
		log.Println(err)
	} else {
		log.Printf("用户[%s]已存在！\n", user.UserName)
	}

	// 添加角色
	role, err := service.NewSysRoleService().GetRoleByEnName("admin")
	if gorm.IsRecordNotFoundError(err) {
		role := &model.SysRole{
			EnName:  "admin",
			Name:    "管理员",
			Comment: "",
		}
		createRole, err := service.NewSysRoleService().CreateRole(role)
		if err != nil {
			log.Println(err)
		} else {
			fmt.Printf("角色[%s]数据插入成功。\n", createRole.EnName)
		}
	} else if err != nil {
		log.Println(err)
	} else {
		log.Printf("角色[%s]已存在！\n", role.EnName)
	}

	role, err = service.NewSysRoleService().GetRoleByEnName("user")
	if gorm.IsRecordNotFoundError(err) {
		role := &model.SysRole{
			EnName:  "user",
			Name:    "普通用户",
			Comment: "",
		}
		createRole, err := service.NewSysRoleService().CreateRole(role)
		if err != nil {
			log.Println(err)
		} else {
			fmt.Printf("角色[%s]数据插入成功。\n", createRole.EnName)
		}
	} else if err != nil {
		log.Println(err)
	} else {
		log.Printf("角色[%s]已存在！\n", role.EnName)
	}

	// 添加权限
	permission, err := service.NewSysPermissionService().GetPermissionByEnName("UMS_VIEW")
	if gorm.IsRecordNotFoundError(err) {
		permission := &model.SysPermission{
			EnName:  "UMS_VIEW",
			Name:    "UMS系统查询权限",
			Comment: "",
		}
		createPermission, err := service.NewSysPermissionService().CreatePermission(permission)
		if err != nil {
			log.Println(err)
		} else {
			fmt.Printf("权限[%s]数据插入成功。\n", createPermission.EnName)
		}
	} else if err != nil {
		log.Println(err)
	} else {
		log.Printf("权限[%s]已存在！\n", permission.EnName)
	}

	permission, err = service.NewSysPermissionService().GetPermissionByEnName("UMS_EDIT")
	if gorm.IsRecordNotFoundError(err) {
		permission := &model.SysPermission{
			EnName:  "UMS_EDIT",
			Name:    "UMS系统编辑权限",
			Comment: "",
		}
		createPermission, err := service.NewSysPermissionService().CreatePermission(permission)
		if err != nil {
			log.Println(err)
		} else {
			fmt.Printf("权限[%s]数据插入成功。\n", createPermission.EnName)
		}
	} else if err != nil {
		log.Println(err)
	} else {
		log.Printf("权限[%s]已存在！\n", permission.EnName)
	}

	// 用户角色信息
	service.NewSysUserRoleService().DeleteUserRole(1, 1)
	service.NewSysUserRoleService().DeleteUserRole(1, 2)

	userRole := &model.SysUserRole{
		UserId: 1,
		RoleId: 1,
	}
	createUserRole, err := service.NewSysUserRoleService().CreateUserRole(userRole)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Printf("用户-角色[%v - %v]数据插入成功。\n", createUserRole.UserId, createUserRole.RoleId)
	}

	userRole = &model.SysUserRole{
		UserId: 1,
		RoleId: 2,
	}
	createUserRole, err = service.NewSysUserRoleService().CreateUserRole(userRole)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Printf("用户-角色[%v - %v]数据插入成功。\n", createUserRole.UserId, createUserRole.RoleId)
	}

	// 角色权限
	rolePermissionService := service.NewSysRolePermissionService()
	rolePermissionService.DeleteRolePermission(1, 1)
	rolePermissionService.DeleteRolePermission(1, 2)

	rp := &model.SysRolePermission{
		RoleId:       1,
		PermissionId: 1,
	}
	rp, err = rolePermissionService.CreateRolePermission(rp)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Printf("角色-权限[%v - %v]数据插入成功。\n", rp.RoleId, rp.PermissionId)
	}

	rp = &model.SysRolePermission{
		RoleId:       1,
		PermissionId: 2,
	}
	rp, err = rolePermissionService.CreateRolePermission(rp)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Printf("角色-权限[%v - %v]数据插入成功。\n", rp.RoleId, rp.PermissionId)
	}
}
