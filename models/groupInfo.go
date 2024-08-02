package models

import "gorm.io/gorm"

type GroupInfo struct {
	gorm.Model
	Name    string //群名称
	OwnerId uint   //群拥有者
	Type    int    //群类型
	Icon    string //头像
	Desc    string //描述
}

func (g *GroupInfo) GroupTableName() string {
	return "groupInfo"
}
