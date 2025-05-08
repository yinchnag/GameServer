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

package network

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"io"
	"strings"
	"unicode/utf8"

	"wgame_server/libray/core"

	"github.com/gorilla/websocket"
)

const (
	PACK_SIZE  = 4 // 包头长度
	MSGID_SIZE = 2 // 消息ID
)

// 解密PB包
func HF_DecodeMsgPB(msgData []byte) (int, int, error) {
	// 丢弃无效消息
	headerSize := PACK_SIZE + MSGID_SIZE
	if len(msgData) < headerSize {
		return 0, 0, fmt.Errorf("invalid message size=%d ", len(msgData))
	}

	// 读取消息头
	reader := bytes.NewReader(msgData[:headerSize])
	msglen := int32(0)
	msgid := uint16(0)
	errRead := binary.Read(reader, binary.LittleEndian, &msglen)
	if errRead != nil {
		return 0, 0, errRead
	}
	errRead = binary.Read(reader, binary.LittleEndian, &msgid)
	if errRead != nil {
		return 0, 0, errRead
	}
	if int(msglen) != len(msgData) {
		return 0, 0, fmt.Errorf("invalid message, len=%d msgid=%d", msglen, msgid)
	}
	return int(msglen), int(msgid), nil
}

// 加密PB消息
func HF_EncodeMsgPB(msgId uint16, body []byte) []byte {
	var buffer bytes.Buffer
	msglen := uint32(len(body) + PACK_SIZE + MSGID_SIZE)
	binary.Write(&buffer, binary.LittleEndian, &msglen)
	binary.Write(&buffer, binary.LittleEndian, &msgId)
	binary.Write(&buffer, binary.LittleEndian, body)
	return buffer.Bytes()
}

// 转JSON字符串
func HF_JtoA(v interface{}) string {
	s, err := core.Marshal(v)
	if err != nil {
		core.Logger.Error("HF_JtoA err:")
	}
	return string(s)
}

// 转JSON字符串
func HF_AtoJ(src string, proto interface{}) bool {
	err := core.Unmarshal([]byte(src), proto)
	if err != nil {
		core.Logger.Error("HF_AtoJ err:")
		return false
	}
	return true
}

// 转JSON二进制
func HF_JtoB(v interface{}) []byte {
	s, err := core.Marshal(v)
	if err != nil {
		core.Logger.Error("HF_JtoB err:")
	}
	return s
}

// 转JSON二进制
func HF_BtoJ(src interface{}, proto interface{}) bool {
	if src == nil {
		core.Logger.Error("HF_JtoB err: nil")
		return false
	}
	buff, ok := src.([]byte)
	if !ok {
		core.Logger.Error("HF_JtoB err: not []byte")
		return false
	}
	err := core.Unmarshal(buff, proto)
	if err != nil {
		core.Logger.Error("HF_JtoB err: unmarshal failed")
		return false
	}
	return true
}

// 压缩并转码Base64
func HF_CompressAndBase64(data []byte) string {
	var buf bytes.Buffer
	compressor := zlib.NewWriter(&buf)
	compressor.Write(data)
	compressor.Close()

	return base64.StdEncoding.EncodeToString(buf.Bytes())
}

// 解码Base64并解压缩
func HF_Base64AndDecompress(data string) []byte {
	defer func() {
		x := recover()
		if x != nil {
			core.Logger.Error("HF_Base64AndDecompress:")
		}
	}()

	// 对上面的编码结果进行base64解码
	decodeBytes, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		core.Logger.Error("base64 decode err:", err)
	}
	buff := bytes.NewReader(decodeBytes)
	var out bytes.Buffer
	reader, _ := zlib.NewReader(buff)
	io.Copy(&out, reader)
	return out.Bytes()
}

// 克隆对象 dst为指针(gob)
func HF_DeepCopy(dst, src interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}

// 克隆对象 dst为指针(json)
func HF_DeepCopy_Json(dst, src interface{}) error {
	buf, err := core.Marshal(src)
	if err != nil {
		return err
	}
	return core.Unmarshal(buf, dst)
}

// 过滤 emoji 表情
func HF_FilterEmoji(content string) string {
	new_content := ""
	for _, value := range content {
		_, size := utf8.DecodeRuneInString(string(value))
		if size <= 3 {
			new_content += string(value)
		}
	}
	return new_content
}

// 是否合法
func HF_IsLicitName(name []byte) bool {
	for i := 0; i < len(name); i++ {
		switch name[i] {
		case '\r', '\'', '\n', ' ', '	', '"', '\\':
			return false
		default:
		}
	}
	return true
}

// 得到客户端ip
func HF_GetWsConnIP(req *websocket.Conn) string {
	ip := req.RemoteAddr().String()
	return strings.Split(ip, ":")[0]
}
