package handler

import (
	"api_router/push_srv/db"
	"api_router/base/data"
	//"api_router/base/service"
	service "api_router/base/service2"
	"sync"
	l4g "github.com/alecthomas/log4go"
	"api_router/base/nethelper"
	"encoding/json"
	"bastionpay_api/api"
)

type Push struct{
	rwmu sync.RWMutex
	usersCallbackUrl map[string]string
}

var defaultPush = &Push{}

func PushInstance() *Push{
	return defaultPush
}

func (push * Push)Init() {
	push.usersCallbackUrl = make(map[string]string)
}

func (push * Push)getUserCallbackUrl(userKey string) (string, error)  {
	url := func() string{
		push.rwmu.RLock()
		defer push.rwmu.RUnlock()

		return push.usersCallbackUrl[userKey]
	}()
	if url != "" {
		return url,nil
	}

	return func() (string, error){
		push.rwmu.Lock()
		defer push.rwmu.Unlock()

		url := push.usersCallbackUrl[userKey]
		if url != "" {
			return url, nil
		}
		url, err := db.ReadUserCallbackUrl(userKey)
		if err != nil {
			return "", err
		}
		push.usersCallbackUrl[userKey] = url
		return url, nil
	}()
}

func (push * Push)GetApiGroup()(map[string]service.NodeApi){
	nam := make(map[string]service.NodeApi)

	func(){
		service.RegisterApi(&nam,
			"pushdata", data.APILevel_client, push.PushData)
	}()

	return nam
}

// 推送数据
func (push *Push)PushData(req *data.SrvRequestData, res *data.SrvResponseData) {
	url, err := push.getUserCallbackUrl(req.Data.Argv.UserKey)
	if err != nil {
		l4g.Error("(%s) no user callback: %s",req.Data.Argv.UserKey, err.Error())
		res.Data.Err = data.ErrPushSrvPushData
		return
	}

	func(){
		pushData := api.UserResponseData{}
		pushData.Value = req.Data.Argv

		// call url
		b, err := json.Marshal(pushData)
		if err != nil {
			l4g.Error("error json message: %s", err.Error())
			res.Data.Err = data.ErrDataCorrupted
			return
		}
		var ret string
		err = nethelper.CallToHttpServer(url, "", string(b), &ret)
		if err != nil {
			l4g.Error("push http: %s", err.Error())
			res.Data.Err = data.ErrPushSrvPushData
			return
		}

		res.Data.Value.Message = ret
	}()

	l4g.Info("push:", req.Data.Argv.Message)
}