package main

import (
	"github.com/HanYin666/goweb/gateway"
	"github.com/HanYin666/goweb/msgo"
	"github.com/HanYin666/goweb/register"
	"net/http"
	"time"
)

func main() {
	engine := msgo.Default()
	engine.OpenGateway = true
	var configs []gateway.GWConfig
	configs = append(configs, gateway.GWConfig{
		Name: "order",
		Path: "/order/**",
		Header: func(req *http.Request) {
			req.Header.Set("my", "mszlu")
		},
		ServiceName: "orderCenter",
	}, gateway.GWConfig{
		Name: "goods",
		Path: "/goods/**",
		Header: func(req *http.Request) {
			req.Header.Set("my", "mszlu")
		},
		ServiceName: "goodsCenter",
	})
	engine.SetGatewayConfig(configs)
	engine.RegisterType = "etcd"
	engine.RegisterOption = register.Option{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	}
	engine.Run(":80")
}
