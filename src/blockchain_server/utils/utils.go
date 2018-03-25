package utils

import (
	"fmt"
	"math/big"
	"path/filepath"
	"os"
	"strings"
	l4g "github.com/alecthomas/log4go"
)

func Faltal_error(err error) {
	if err==nil {
		return
	}
	l4g.Error(err)
	os.Exit(1)
}


func string_has_prefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[0:len(prefix)] == prefix
}

func String_cat_prefix(s string, prefix string) (string) {
	if string_has_prefix(s, prefix) {
		return s[len(prefix):]
	}
	return s
}

func Hex_string_to_big_int(s string) (*big.Int, error) {
	var (
		err        error
		bigint     *big.Int
		string_hex string
	)
	string_hex = String_cat_prefix(s, "0x")

	bigint, isok := new(big.Int).SetString(string_hex, 16)
	if !isok {
		err = fmt.Errorf("can not convert hex string:%s, to big int", string_hex)
	}

	return bigint, err
}


func CurrentRuningFileDir() string {
	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	//path, err := filepath.Abs("./")
	if err != nil {
		Faltal_error(err)
	}
	return strings.Replace(path, "\\", "/", -1)
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

