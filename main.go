package main

import (
	"HydraServer/controller"
	"HydraServer/util/conf"
	"github.com/gin-gonic/gin"
	"math/rand"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}
func main() {
	//1.创建路由
	if conf.RunMode == "dev" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	//2.绑定路由规则，执行的函数
	//gin.Context，封装了request和response

	r.POST("/hydra/create", controller.Hydra)
	r.POST("/hydra/upload/username", controller.SaveUsername)
	r.POST("/hydra/upload/password", controller.SavePassword)
	//3.监听端口，默认在8080
	//Run("里面不指定端口号默认为8080")
	if err := r.SetTrustedProxies([]string{"0.0.0.0"}); err != nil {
		return
	}
	port := conf.GetVal("server", "port")
	addr := strings.Join([]string{":", port}, "")
	err := r.Run(addr)
	if err != nil {
		panic(err)
	}
}
