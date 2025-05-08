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
	"reflect"

	"wgame_server/libray/core"
	"wgame_server/libray/network"
)

// 解密存档包
func SQL_Decode(dst interface{}, src interface{}) {
	defer func() {
		if err := recover(); err != nil {
			core.Logger.Errorln("SQL_Decode 拦截到错误:", err)
		}
	}()
	srcType := reflect.TypeOf(src).Elem()
	srcData := reflect.ValueOf(src).Elem()
	compare := map[string][]int{}
	for i := 0; i < srcType.NumField(); i++ {
		name := srcType.Field(i).Name
		tag := srcType.Field(i).Tag.Get("json")
		if tag == "" {
			continue
		}
		ignore := srcType.Field(i).Tag.Get("ignore")
		if ignore == "1" {
			continue
		}
		_, ok := compare[name]
		if ok {
			core.Logger.Warn("重复映射", name)
			continue
		}
		compare[name] = []int{i, -1}
	}

	dstType := reflect.TypeOf(dst).Elem()
	dstData := reflect.ValueOf(dst).Elem()
	for i := 0; i < dstType.NumField(); i++ {
		name := dstType.Field(i).Name
		val, ok := compare["T_"+name]
		if !ok {
			// core.Logger.Debug("未映射", name)
			continue
		}
		val[1] = i
	}

	for k, v := range compare {
		if v[0] < 0 || v[1] < 0 {
			core.Logger.Info("无效未映射", k)
			continue
		}
		srcField := srcData.Field(v[0])
		dstField := dstData.Field(v[1])
		if srcField.Kind() == reflect.String && !CheckDataKind(dstField.Kind()) {
			typField := dstType.Field(v[1])
			srcFieldVal := srcField.Interface().(string)
			if srcFieldVal == "" {
				continue
			}
			if dstField.Kind() == reflect.Array {
				continue // 不支持数组类型存档
			}
			swapData := reflect.New(typField.Type)
			swapInter := swapData.Interface()
			err := core.Base64Decode(srcFieldVal, &swapInter)
			if err != nil {
				core.Logger.Errorf("SQL_Decode %s err: %v", k, err.Error())
			} else {
				dstField.Set(swapData.Elem())
			}
		} else {
			dstField.Set(srcField)
		}
	}
}

// 加密存档包
func SQL_Encode(dst interface{}, src interface{}) {
	defer func() {
		if err := recover(); err != nil {
			core.Logger.Errorln("SQL_Encode 拦截到错误:", err)
		}
	}()
	dstType := reflect.TypeOf(dst).Elem()
	dstData := reflect.ValueOf(dst).Elem()
	compare := map[string][]int{}
	for i := 0; i < dstType.NumField(); i++ {
		tag := dstType.Field(i).Tag.Get("json")
		name := dstType.Field(i).Name
		if tag == "" {
			continue
		}
		ignore := dstType.Field(i).Tag.Get("ignore")
		if ignore == "1" {
			continue
		}
		_, ok := compare[name]
		if ok {
			core.Logger.Warn("重复映射", name)
			continue
		}
		compare[name] = []int{-1, i}
	}

	srcType := reflect.TypeOf(src).Elem()
	srcData := reflect.ValueOf(src).Elem()
	for i := 0; i < srcType.NumField(); i++ {
		name := srcType.Field(i).Name
		val, ok := compare["T_"+name]
		if !ok {
			// core.Logger.Warn("未映射", name)
			continue
		}
		val[0] = i
	}

	for k, v := range compare {
		if v[0] < 0 || v[1] < 0 {
			core.Logger.Info("无效未映射", k)
			continue
		}
		src := srcData.Field(v[0])
		dst := dstData.Field(v[1])
		if dst.Kind() == reflect.String && !CheckDataKind(src.Kind()) {
			if src.Kind() == reflect.Array {
				continue // 不支持数组类型存档
			}
			srcVal := src.Interface()
			encodeVal := core.Base64Encode(srcVal)
			dst.Set(reflect.ValueOf(encodeVal))
		} else {
			srcVal := src.Addr().Interface()
			dstVal := dst.Addr().Interface()
			network.HF_DeepCopy_Json(dstVal, srcVal)
		}
	}
}
