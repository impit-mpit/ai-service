package main

import (
	"neuro-most/ai-service/config"
	"neuro-most/ai-service/internal/infra"
)

func main() {
	cfg, err := config.NewLoadConfig()
	if err != nil {
		panic(err)
	}
	infra.Config(cfg).Vllm().Serve().Start()
}
