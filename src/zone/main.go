package main

import (
	_ "net/http/pprof"
	"sync"

	"wgame_server/libray/actor"
	"wgame_server/module/activity"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func Test() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 迁移 schema
	db.AutoMigrate(&Product{})

	// Create
	db.Create(&Product{Code: "D42", Price: 100})

	// Read
	var product Product
	db.First(&product, 1)                 // 根据整型主键查找
	db.First(&product, "code = ?", "D42") // 查找 code 字段值为 D42 的记录

	// Update - 将 product 的 price 更新为 200
	db.Model(&product).Update("Price", 200)
	// Update - 更新多个字段
	db.Model(&product).Updates(Product{Price: 200, Code: "F42"}) // 仅更新非零值字段
	db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

	// Delete - 删除 product
	db.Delete(&product, 1)
}

func main() {
	actorSystem := actor.NewActorSystem()
	actorSystem.AllocActor(func() actor.IRceiver { return &activity.ActivityActor{} })
	actorSystem.Start()
	rets, err := actorSystem.ModInvoke(nil, 0, "ActivityActor", "GetInt", "1")
	ret, _ := rets[0].Interface().(int), err
	_ = ret
	wg := &sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}
