package internal

type ActorConfigOption func(config *actorConfig)

type actorConfig struct {
	Capacity   uint64 // 列表容量
	Dropping   bool   // 是否丢弃溢出消息
	Throughput int    // 吞吐量
}

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
