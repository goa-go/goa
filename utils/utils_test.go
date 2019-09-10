package utils

import (
	"testing"
)

func TestStr2Bytes(t *testing.T) {
	str := "abcABC123(+-*/)[啊，。！]"
	bytes := Str2Bytes(str)

	for i, v := range bytes {
		if v != []byte(str)[i] {
			t.Errorf("Str2Bytes failed: %s -> %v", str, bytes)
		}
	}
}

func TestBytes2Str(t *testing.T) {
	bytes := []byte("abcABC123(+-*/)[啊，。！]")
	str := Bytes2Str(bytes)

	if str != "abcABC123(+-*/)[啊，。！]" {
		t.Errorf("Bytes2Str failed: %v -> %s", bytes, str)
	}
}
