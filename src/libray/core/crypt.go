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

package core

import (
	"encoding/base64"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

// 转码常量
var UUID_CHARS = []rune{
	'a', 'b', 'c', 'd', 'e', 'f',
	'g', 'h', 'i', 'j', 'k', 'l',
	'm', 'n', 'o', 'p', 'q', 'r',
	's', 't', 'u', 'v', 'w', 'x',
	'y', 'z', '1', '2', '3', '4',
	'5', '6', '7', '8', '9', '0',
	'A', 'B', 'C', 'D', 'E', 'F',
	'G', 'H', 'I', 'J', 'K', 'L',
	'M', 'N', 'O', 'P', 'Q', 'R',
	'S', 'T', 'U', 'V', 'W', 'X',
	'Y', 'Z', '_', '-',
}

// Base64加密
func Base64Encode(src interface{}) string {
	tmp, err := Marshal(src)
	if err != nil {
		Logger.Error("Base64Encode err:")
	}
	return base64.StdEncoding.EncodeToString(tmp)
}

// Base64解密
func Base64Decode(src string, proto interface{}) error {
	tmp, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return err
	}
	err = Unmarshal(tmp, proto)
	if err != nil {
		return err
	}
	return nil
}

// 8位UUID
func UUID8() string {
	uuidVal := uuid.New()
	code := strings.Replace(uuidVal.String(), "-", "", -1)
	buffer := ""
	for i := 0; i < 8; i++ {
		str := code[i*4 : i*4+4]
		num, _ := strconv.ParseUint(str, 16, 16)
		buffer += string(UUID_CHARS[num%0x3e])
	}
	return buffer
}

// 生成一个自增的uuid
func CreateItemAutoUUID(itemID uint32, autoNum *uint32) uint64 {
	var uuid uint64 = uint64(*autoNum)<<32 | uint64(itemID)
	*autoNum++
	return uuid
}

// 从物品UUID中获得背包模块ID，物品大类型，物品子类型
func ByUUIDGetItemID(uuid uint64) uint32 {
	ITEM_ID := uint64(0xffffffff00000000)

	return uint32(uuid | ITEM_ID)
}
