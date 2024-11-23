package infra

import (
	"neuro-most/ai-service/config"
	"neuro-most/ai-service/internal/adapters/vllm"
	"neuro-most/ai-service/internal/infra/router"
)

type app struct {
	cfg    config.Config
	router router.Router
	vllm   vllm.Vllm
}

func Config(cfg config.Config) *app {
	return &app{cfg: cfg}
}

func (a *app) Vllm() *app {
	a.vllm = *vllm.NewVllm(a.cfg.OpenApiUrl, "token123")
	return a
}

func (a *app) Serve() *app {
	a.router = router.NewRouter(a.vllm)
	return a
}

func (a *app) Start() {
	a.router.Listen()
}
