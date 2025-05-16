package actor

import "runtime"

type (
	ActorSystemConfigOption func(config *ActorSystemConfig)
	ActorSystemConfig       struct {
		Capacity   int    // 列表容量
		Throughput int    // 吞吐量
		Cluster    uint32 // 集群ID
		Profile    bool   // 是否调试
	}
)

func defaultActorSystemConfig() *ActorSystemConfig {
	return &ActorSystemConfig{
		Capacity:   1024,
		Throughput: runtime.NumCPU() * 2,
	}
}

func ActorSystemConfigure(options ...ActorSystemConfigOption) *ActorSystemConfig {
	config := defaultActorSystemConfig()
	for _, option := range options {
		option(config)
	}
	return config
}

type (
	ActorConfigOption func(config *actorConfig)
	actorConfig       struct {
		Capacity   uint64 // 列表容量
		Dropping   bool   // 是否丢弃溢出消息
		Throughput int    // 吞吐量
	}
)

func defaultActorConfig() *actorConfig {
	return &actorConfig{
		Capacity:   1024,
		Dropping:   false,
		Throughput: 16,
	}
}

func ActorConfigure(option ...ActorConfigOption) *actorConfig {
	config := defaultActorConfig()
	for _, option := range option {
		option(config)
	}
	return config
}
