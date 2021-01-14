package utils

import (
	"crypto/md5"
	"encoding/hex"
	"math"
	"math/rand"
	"strconv"
	"time"
	"unsafe"
)

func HexMd5(str string) string {
	h := md5.New()
	strPtr := (*[2]uintptr)(unsafe.Pointer(&str))
	strBs := *(*[]byte)(unsafe.Pointer(&[3]uintptr{strPtr[0], strPtr[1], strPtr[1]}))
	h.Write(strBs)
	ciphers := h.Sum(nil)
	res := make([]byte, 32)
	hex.Encode(res, ciphers)
	return *(*string)(unsafe.Pointer(&res))
}

//Int64ToStr int64转string
func Int64ToStr(i int64) string {
	return strconv.FormatInt(i, 10)
}

func Rand(min, max int) int {
	if min > max {
		panic("min: min cannot be greater than max")
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	n := r.Intn(math.MaxInt16)
	return n/((math.MaxInt16)/(max-min+1)) + min
}

func IsInStringArray(arr []string, item string) bool {
	for i := range arr {
		if arr[i] == item {
			return true
		}
	}
	return false
}
