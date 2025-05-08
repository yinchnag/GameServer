package backstage

import (
	"wgame_server/libray/module"

	"github.com/gin-gonic/gin"
)

// 后台管理器
type BackstageMgr struct {
	module.ModObj
}

func (that *BackstageMgr) Init(mod interface{}) module.IModule {
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.Run() // 监听并在 0.0.0.0:8080 上启动服务
	return that
}
