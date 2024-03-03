package http

import "github.com/go-chi/chi/v5"

type Handler struct {
	goodService GoodService
}

func NewHandler(goodService GoodService) *Handler {
	return &Handler{goodService: goodService}
}

func (h *Handler) InitRoutes(router chi.Router) chi.Router {
	router.Route("/api/v1", func(apiV1 chi.Router) {

		apiV1.Route("/good", func(good chi.Router) {
			good.Post("/create", h.CreateGoodHandler())
			good.Patch("/update", h.UpdateGoodHandler())
			good.Patch("/reprioritize", h.UpdateGoodPriorityHandler())
			good.Delete("/remove", h.DeleteGoodHandler())
		})

		apiV1.Route("/goods", func(goods chi.Router) {
			goods.Get("/list", h.ListGoodsHandler())
		})
	})

	return router
}
