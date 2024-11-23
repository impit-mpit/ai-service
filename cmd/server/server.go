package main

import (
	"neuro-most/template-service/config"
	"neuro-most/template-service/internal/infra"
)

func main() {
	cfg, err := config.NewLoadConfig()
	if err != nil {
		panic(err)
	}
	infra.Config(cfg).Database().Serve().Start()
}
