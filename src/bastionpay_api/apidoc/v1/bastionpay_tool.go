package v1

import (
	"bastionpay_api/apidoc"
	"bastionpay_api/api/v1"
)

var ApiDocRecharge = apidoc.ApiDoc{
	VerName:"v1",
	SrvName:"bastionpay",
	FuncName:"recharge",
	Level:0,
	Comment:"模拟充值",
	Path:"/v1/bastionpay/recharge",
	Input:&v1.ReqRecharge{},
	Output:new(string),
}

var ApiDocGenerate = apidoc.ApiDoc{
	VerName:"v1",
	SrvName:"bastionpay",
	FuncName:"generate",
	Level:0,
	Comment:"模拟挖矿",
	Path:"/v1/bastionpay/generate",
	Input:&v1.ReqGenerate{},
	Output:new(string),
}