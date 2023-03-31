package id

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"sync"
	"unsafe"
)

// ID 标准数字唯一键
type ID uint64

// String 返回 ID 为 hex 字符串
func (id ID) String() string {
	return fmt.Sprintf("%016x", uint64(id))
}

// MarshalJSON 编码 ID 为 JSON encoding
func (id ID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

// UnmarshalJSON 解码 给定数据为 Json encoding 或 Json 数字
func (id *ID) UnmarshalJSON(data []byte) error {
	i, err := parseJSONString(data)
	if err == nil {
		*id = i
		return nil
	}
	i, err = parseJSONInt(data)
	if err == nil {
		*id = i
		return nil
	}
	return fmt.Errorf("%s is not a valid ID", data)
}

// ParseID 转译给定字符串为 ID
func ParseID(s string) (ID, error) {
	i, err := strconv.ParseUint(s, 16, 64)
	if err != nil {
		return 0, err
	}
	return ID(i), nil
}

// parseJSONString 把 Json encoding 转化为 ID
func parseJSONString(data []byte) (ID, error) {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return 0, err
	}
	i, err := ParseID(s)
	if err != nil {
		return 0, err
	}
	return i, nil
}

// parseJSONString 把 Json 数字 转化为 ID
func parseJSONInt(data []byte) (ID, error) {
	var i uint64
	if err := json.Unmarshal(data, &i); err != nil {
		return 0, err
	}
	return ID(i), nil
}

const (
	idSize  = aes.BlockSize / 2 // 64 bits
	keySize = aes.BlockSize     // 128 bits
)

var (
	ctr []byte
	n   int
	b   []byte
	c   cipher.Block
	m   sync.Mutex
)

func init() {
	buf := make([]byte, keySize+aes.BlockSize)
	_, err := io.ReadFull(rand.Reader, buf)
	if err != nil {
		panic(err) // /dev/urandom had better work
	}
	c, err = aes.NewCipher(buf[:keySize])
	if err != nil {
		panic(err) // AES had better work
	}
	n = aes.BlockSize
	ctr = buf[keySize:]
	b = make([]byte, aes.BlockSize)
}

// GenerateID 返回一个随机生成的 64 字节 ID。这个函数进程安全
// 这个方法大概花费 13.29 ns, 不产出任何内存垃圾
func GenerateID() ID {
	m.Lock()
	if n == aes.BlockSize {
		c.Encrypt(b, ctr)
		for i := aes.BlockSize - 1; i >= 0; i-- { // increment ctr
			ctr[i]++
			if ctr[i] != 0 {
				break
			}
		}
		n = 0
	}
	id := *(*ID)(unsafe.Pointer(&b[n])) // zero-copy b/c we're arch-neutral
	n += idSize
	m.Unlock()
	return id
}
