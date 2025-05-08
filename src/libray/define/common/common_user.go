package common

// 无法修改数据，或者该数据只能由玩家本人修改
// 该数据在玩家不在线的时候也会加载到服务器缓存中，以方便数据查找
type IUser interface {
	GetUid() int64
	GetUName() string
}
type Player_ struct {
	Uid   int64  // 玩家唯一id
	UName string // 玩家昵称

}
