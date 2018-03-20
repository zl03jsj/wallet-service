package config

import (
	"blockchain_server/utils"
	"encoding/json"
	"io/ioutil"
	l4g "github.com/alecthomas/log4go"
	"fmt"
	"os"
)

var (
	configer Configer
	fordebug = true
)

type ClientConfig struct {
	RPC_url string				   `json:"rpc_url"`
	Start_scan_Blocknumber uint64  `json:"start_sacn_blocknumber,string,omitempty"`
}

type Configer struct {
	Cryptofile  	string	`json:"crypto_file"`
	Log_conf_file	string  `json:"log_conf_file"`
	Log_path		string  `json:"log_path"`

	Clientconfig    map[string]*ClientConfig
}


func (self *ClientConfig)SaveConfig() error {
	return nil
}

func check(err error) {
    if err != nil {
    	fmt.Printf("error :%s\n", err)
    	os.Exit(1)
    }
}

func init () {
	configfile := getConfigFilePath()
	dat, err := ioutil.ReadFile(configfile)
	check(err)
	err = json.Unmarshal(dat, &configer)
	check(err)
	l4g.LoadConfiguration(configer.Log_conf_file)
}

func GetConfiger() (*Configer) {
	return &configer
}

func getConfigFilePath() string {
	if fordebug {
		return "/Users/cengliang/code/wallet-service/src/blockchain_server/res/app_debug.config"
	}
	return utils.CurrentRuningFileDir() + "/../res/app.config"
}

//func CryptoKeyFilePath() string {
//	return utils.CurrentRuningFileDir() + "/cryptokey.json"
//}
//func NewConfiger(config_file string) {
//    dat, err := ioutil.ReadFile(config_file)
//    check(err)
//    err = json.Unmarshal(dat, configer)
//    check(err)
//fmt.Print(string(dat))
////你经常会想对于一个文件是怎么读并且读取到哪一部分进行更多的控制。对于这个任务，从使用 os.Open打开一个文件获取一个 os.File 值开始。
//f, err := os.Open("/tmp/dat")
//check(err)
////从文件开始位置读取一些字节。这里最多读取 5 个字节，并且这也是我们实际读取的字节数。
//b1 := make([]byte, 5)
//n1, err := f.Read(b1)
//check(err)
//fmt.Printf("%d bytes: %s\n", n1, string(b1))
////你也可以 Seek 到一个文件中已知的位置并从这个位置开始进行读取。
//o2, err := f.Seek(6, 0)
//check(err)
//b2 := make([]byte, 2)
//n2, err := f.Read(b2)
//check(err)
//fmt.Printf("%d bytes @ %d: %s\n", n2, o2, string(b2))
////io 包提供了一些可以帮助我们进行文件读取的函数。例如，上面的读取可以使用 ReadAtLeast 得到一个更健壮的实现。
//o3, err := f.Seek(6, 0)
//check(err)
//b3 := make([]byte, 2)
//n3, err := io.ReadAtLeast(f, b3, 2)
//check(err)
//fmt.Printf("%d bytes @ %d: %s\n", n3, o3, string(b3))
////没有内置的回转支持，但是使用 Seek(0, 0) 实现。
//_, err = f.Seek(0, 0)
//check(err)
////bufio 包实现了带缓冲的读取，这不仅对有很多小的读取操作的能提升性能，也提供了很多附加的读取函数。
//r4 := bufio.NewReader(f)
//b4, err := r4.Peek(5)
//check(err)
//fmt.Printf("5 bytes: %s\n", string(b4))
////任务结束后要关闭这个文件（通常这个操作应该在 Open操作后立即使用 defer 来完成）。
//f.Stop()
//}
