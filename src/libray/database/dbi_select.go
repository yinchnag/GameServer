// JoysGames copyrights this specification. No part of this specification may be
// reproduced in any form or means, without the prior written consent of JoysGames.
//
// This specification is preliminary and is subject to change at any time without notice.
// JoysGames assumes no responsibility for any errors contained herein.
//
// The above copyright notice and this permission notice shall be included
// in all copies or substantial portions of the Software.
// @package JGServer
// @copyright joysgames.cn All rights reserved.
// @version v1.0

package database

import (
	"fmt"
	"reflect"
	"strings"

	"wgame_server/libray/core"
)

// 统计数量
const FROM_TOTAL string = "COUNT(1) AS total"

// 选择器类型
type SelectType string

const (
	SelectType_Insert  SelectType = "insert"  // 插入
	SelectType_Replace SelectType = "replace" // 覆盖
	SelectType_Select  SelectType = "select"  // 查询
	SelectType_Update  SelectType = "update"  // 更新
	SelectType_Delete  SelectType = "delete"  // 删除
	SelectType_Create  SelectType = "create"  // 创建
)

// 检查数据类型
func CheckDataKind(valType reflect.Kind) bool {
	switch valType {
	case reflect.Int8:
		return true
	case reflect.Int16:
		return true
	case reflect.Int32:
		return true
	case reflect.Int64:
		return true
	case reflect.Int:
		return true
	case reflect.Uint8:
		return true
	case reflect.Uint16:
		return true
	case reflect.Uint32:
		return true
	case reflect.Uint64:
		return true
	case reflect.Float32:
		return true
	case reflect.Float64:
		return true
	case reflect.String:
		return true
	default:
		return false
	}
}

// 选择器参数
type IOptConf interface {
	ToString() string
	GetType() SelectType
}

// 插入参数
type InsertConf struct {
	Table string                 // 表名称
	Data  map[string]interface{} // 数据
}

// 转义
func (that *InsertConf) ToString() string {
	sql := "INSERT INTO "
	keys := []string{}
	values := []string{}
	for key, value := range that.Data {
		if !CheckDataKind(reflect.TypeOf(value).Kind()) {
			continue
		}
		keys = append(keys, fmt.Sprintf("`%s`", key))
		values = append(values, fmt.Sprintf("'%v'", value))
	}
	sql += that.Table + " (" + strings.Join(keys, ",") + ") "
	sql += "VALUES (" + strings.Join(values, ",") + ") "
	return sql
}

// 获取类型
func (that *InsertConf) GetType() SelectType {
	return SelectType_Insert
}

// 替换参数
type ReplaceConf struct {
	Table string                 // 表名称
	Data  map[string]interface{} // 数据
}

// 转义
func (that *ReplaceConf) ToString() string {
	sql := "REPLACE INTO "
	keys := []string{}
	values := []string{}
	for key, value := range that.Data {
		if !CheckDataKind(reflect.TypeOf(value).Kind()) {
			continue
		}
		keys = append(keys, fmt.Sprintf("`%s`", key))
		values = append(values, fmt.Sprintf("'%v'", value))
	}
	sql += that.Table + " (" + strings.Join(keys, ",") + ") "
	sql += "VALUES (" + strings.Join(values, ",") + ") "
	return sql
}

// 获取类型
func (that *ReplaceConf) GetType() SelectType {
	return SelectType_Replace
}

// 查询参数
type SelectConf struct {
	Table   string // 表名称
	Where   string // 条件
	From    string // 选择项
	Groupby string // 分组
	Having  string // 过滤项
	Order   string // 排序
	Limit   string // 限制
	Join    string // 联合查询
}

// 转义
func (that *SelectConf) ToString() string {
	sql := "SELECT "
	if that.From != "" {
		sql += fmt.Sprintf("%s FROM %s ", that.From, that.Table)
	} else {
		sql += fmt.Sprintf("%s.* FROM %s ", that.Table, that.Table)
	}
	if that.Join != "" {
		sql += fmt.Sprintf("%s ", that.Join)
	}
	if that.Where != "" {
		sql += fmt.Sprintf("WHERE %s ", that.Where)
	} else {
		sql += "WHERE 0 "
	}
	if that.Groupby != "" {
		sql += fmt.Sprintf("GROUP BY %s ", that.Groupby)
	}
	if that.Having != "" {
		sql += fmt.Sprintf("HAVING %s ", that.Having)
	}
	if that.Order != "" {
		sql += fmt.Sprintf("ORDER BY %s ", that.Order)
	}
	if that.Limit != "" {
		sql += fmt.Sprintf("LIMIT %s ", that.Limit)
	}
	return sql
}

// 获取类型
func (that *SelectConf) GetType() SelectType {
	return SelectType_Select
}

// 更新参数
type UpdateConf struct {
	Table    string                 // 表名称
	Data     map[string]interface{} // 数据
	Where    string                 // 条件
	Concat   map[string]interface{} // 连接
	Increase map[string]interface{} // 增加
	Decrease map[string]interface{} // 减少
}

// 转义
func (that *UpdateConf) ToString() string {
	sql := fmt.Sprintf("UPDATE %s ", that.Table)
	updates := []string{}
	for key, value := range that.Data {
		if !CheckDataKind(reflect.TypeOf(value).Kind()) {
			continue
		}
		updates = append(updates, fmt.Sprintf("`%s` = '%v' ", key, value))
	}
	for key, value := range that.Concat {
		if !CheckDataKind(reflect.TypeOf(value).Kind()) {
			continue
		}
		updates = append(updates, fmt.Sprintf("`%s` = concat(`%s`,'%v') ", key, key, value))
	}
	for key, value := range that.Increase {
		if !CheckDataKind(reflect.TypeOf(value).Kind()) {
			continue
		}
		updates = append(updates, fmt.Sprintf("`%s` = `%s` + '%v'", key, key, value))
	}
	for key, value := range that.Decrease {
		if !CheckDataKind(reflect.TypeOf(value).Kind()) {
			continue
		}
		updates = append(updates, fmt.Sprintf("`%s` = `%s` - '%v'", key, key, value))
	}
	sql += fmt.Sprintf("SET %s", strings.Join(updates, ","))
	if that.Where != "" {
		sql += fmt.Sprintf("WHERE %s ", that.Where)
	} else {
		sql += "WHERE 0 "
	}
	return sql
}

// 获取类型
func (that *UpdateConf) GetType() SelectType {
	return SelectType_Update
}

// 删除参数
type DeleteConf struct {
	Table string // 表名称
	Where string // 条件
	Limit string // 限制
}

// 转义
func (that *DeleteConf) ToString() string {
	sql := fmt.Sprintf("DELETE FROM %s ", that.Table)
	if that.Where != "" {
		sql += fmt.Sprintf("WHERE %s ", that.Where)
	} else {
		sql += "WHERE 0 "
	}
	if that.Limit != "" {
		sql += fmt.Sprintf("LIMIT %s ", that.Limit)
	}
	return sql
}

// 获取类型
func (that *DeleteConf) GetType() SelectType {
	return SelectType_Delete
}

// 表列项
type TableField struct {
	Name string      // 名称
	Type string      // 类型
	Desc string      // 描述
	Def  interface{} // 默认值
}

// 创建参数
type CreateConf struct {
	Type    SelectType   // 类型
	Table   string       // 表名称
	Columns []TableField // 列项
	Primary []string     // 主键
}

// 转义
func (that *CreateConf) ToString() string {
	sql := fmt.Sprintf("CREATE TABLE `%s` (", that.Table)
	for _, column := range that.Columns {
		idx := core.FindSlice[string](that.Primary, func(value string, key int) bool {
			return value == column.Name
		})
		if idx != -1 {
			sql += fmt.Sprintf("`%s` %s NOT NULL COMMENT '%s' ,", column.Name, column.Type, column.Desc)
		} else {
			sql += fmt.Sprintf("`%s` %s NULL COMMENT '%s' ,", column.Name, column.Type, column.Desc)
		}
	}
	sql += fmt.Sprintf("PRIMARY KEY (`%s`) );", strings.Join(that.Primary, "`, `"))
	return sql
}

// 获取类型
func (that *CreateConf) GetType() SelectType {
	return SelectType_Create
}

// 标准SQL选择器
type DbiSelect struct {
	Conf IOptConf // 参数
}

// 转义
func (that *DbiSelect) ToString() string {
	return that.Conf.ToString()
}

// 获取类型
func (that *DbiSelect) GetType() SelectType {
	return that.Conf.GetType()
}
