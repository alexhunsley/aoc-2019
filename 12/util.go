package main

import (
	"bytes"
	"crypto/sha256"
	"strconv"
)

func computeHashForList(list []int, delim string) [32]byte {
	var buffer bytes.Buffer
	for i, _ := range list {
		buffer.WriteString(strconv.Itoa(list[i]))
		buffer.WriteString(delim)
	}
	return (sha256.Sum256([]byte(buffer.String())))
}

func computeStringHashForListWithDelim(list []int, delim string) string {
	hashData := computeHashForList(list, delim)
	return string(hashData[:32])
}

func computeStringHashForList(list []int) string {
	hashData := computeHashForList(list, ",")
	return string(hashData[:32])
}
