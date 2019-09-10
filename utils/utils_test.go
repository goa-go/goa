package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStr2Bytes(t *testing.T) {
	str := "abcABC123(+-*/)[啊，。！]"
	bytes := Str2Bytes(str)

	assert.Equal(t, []byte("abcABC123(+-*/)[啊，。！]"), bytes)
}

func TestBytes2Str(t *testing.T) {
	bytes := []byte("abcABC123(+-*/)[啊，。！]")
	str := Bytes2Str(bytes)

	assert.Equal(t, "abcABC123(+-*/)[啊，。！]", str)
}
