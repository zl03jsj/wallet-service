package main

import (
	//"api_router/base/service"
	service "api_router/base/service2"
	"api_router/base/data"
	"api_router/account_srv/handler"
	"fmt"
	"context"
	"time"
	"strings"
	"api_router/account_srv/install"
	"encoding/json"
	"api_router/base/utils"
	"errors"
	"api_router/account_srv/user"
	"os"
	l4g "github.com/alecthomas/log4go"
	"api_router/account_srv/db"
)

const AccountSrvConfig = "account.json"

func installWallet(dir string) error {
	var err error

	fi, err := os.Open(dir+"/wallet.install")
	if err == nil {
		defer fi.Close()
		l4g.Info("Super wallet is installed, have a fun...")
		return nil
	}

	l4g.Info("First time to use Super wallet, need to install step by step...")
	newRsa := false
	l4g.Info("1. Create wallet rsa key...")
	_, err = os.Open(dir+"/private.pem")
	if err != nil {
		newRsa = true
	}
	_, err = os.Open(dir+"/public.pem")
	if err != nil {
		newRsa = true
	}
	if newRsa{
		l4g.Info("Create new wallet rsa key in %s", dir)
		pri := fmt.Sprintf(dir+"/private.pem")
		pub := fmt.Sprintf(dir+"/public.pem")
		err = utils.RsaGen(2048, pri, pub)
		if err != nil {
			return err
		}
	}else{
		l4g.Info("A wallet rsa key is exist...")
	}

	l4g.Info("2. Create wallet genesis admin...")

	err = func() error {
		uc, err := install.AddUser(true)
		if err != nil {
			return err
		}
		b, _ := json.Marshal(*uc)

		var req data.SrvRequestData
		var res data.SrvResponseData
		req.Data.Argv.Message = string(b)
		handler.AccountInstance().Create(&req, &res)
		if res.Data.Err != data.NoErr {
			return errors.New(res.Data.ErrMsg)
		}

		uca := user.AckUserCreate{}
		err = json.Unmarshal([]byte(res.Data.Value.Message), &uca)

		l4g.Info("3. Record genesis admin user key: %s", uca.UserKey)
		return nil
	}()
	if err != nil {
		return err
	}


	err = func() error {
		l4g.Info("4. Update wallet genesis admin key...")

		up, err := install.UpdateKey()
		if err != nil {
			return err
		}
		b, _ := json.Marshal(*up)

		var req data.SrvRequestData
		var res data.SrvResponseData
		req.Data.Argv.Message = string(b)

		handler.AccountInstance().UpdateKey(&req, &res)
		if res.Data.Err != data.NoErr {
			return errors.New(res.Data.ErrMsg)
		}

		upa := user.AckUserUpdateKey{}
		err = json.Unmarshal([]byte(res.Data.Value.Message), &upa)

		l4g.Info(upa)
		return nil
	}()
	if err != nil {
		return err
	}

	// write a tag file
	fo, err := os.Create(dir+"/wallet.install")
	if err != nil {
		return err
	}
	defer fo.Close()
	return nil
}

func main() {
	appDir, _:= utils.GetAppDir()
	appDir += "/SuperWallet"

	l4g.LoadConfiguration(appDir + "/log.xml")
	defer l4g.Close()

	cfgPath := appDir + "/" + AccountSrvConfig
	db.Init(cfgPath)

	accountDir := appDir + "/account"
	err := os.MkdirAll(accountDir, os.ModePerm)
	if err != nil && os.IsExist(err) == false {
		l4g.Error("Create dir failed：%s - %s", accountDir, err.Error())
		return
	}

	err = installWallet(accountDir)
	if err != nil {
		l4g.Error("Install super wallet failed: %s", err.Error())
		return
	}

	// init
	handler.AccountInstance().Init(accountDir)

	// create service node
	l4g.Info("config path: %s", cfgPath)
	nodeInstance, err := service.NewServiceNode(cfgPath)
	if nodeInstance == nil || err != nil{
		l4g.Error("Create service node failed: %s", err.Error())
		return
	}

	// register APIs
	service.RegisterNodeApi(nodeInstance, handler.AccountInstance())

	// start service node
	ctx, cancel := context.WithCancel(context.Background())
	service.StartNode(ctx, nodeInstance)

	time.Sleep(time.Second*2)
	for ; ;  {
		fmt.Println("Input 'q' to quit...")
		var input string
		input = utils.ScanLine()

		argv := strings.Split(input, " ")

		if argv[0] == "q" {
			cancel()
			break;
		}
	}

	l4g.Info("Waiting all routine quit...")
	service.StopNode(nodeInstance)
	l4g.Info("All routine is quit...")
}