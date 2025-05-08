package rom

type Meta struct {
	Id int32 `rom:"primary_key"`
}

type IMeta interface {
	// 存储数据,参数为是否强制保存
	// 在为false的情况下会存档中间件，后续再找时机落地mysql
	Save(bool)
	Load(bool)
}
