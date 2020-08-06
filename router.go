package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func newRouter(svc *gatewayService) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(httpLogger)
	r.Use(middleware.StripSlashes)
	r.Use(svc.authn)

	r.Get("/signin", signinHandler)

	r.Route("/finding", func(r chi.Router) {
		r.Use(svc.authzWithProject)
		r.Get("/list-finding", svc.listFindingHandler)
		r.Get("/get-finding", svc.getFindingHandler)
		r.Get("/list-finding-tag", svc.listFindingTagHandler)
		r.Get("/list-resource/", svc.listResourceHandler)
		r.Get("/get-resource/", svc.getResourceHandler)
		r.Get("/list-resource-tag/", svc.listResourceTagHandler)
		r.Group(func(r chi.Router) {
			r.Use(middleware.AllowContentType("application/json"))
			r.Post("/put-finding", svc.putFindingHandler)
			r.Post("/delete-finding", svc.deleteFindingHandler)
			r.Post("/tag-finding", svc.tagFindingHandler)
			r.Post("/untag-finding", svc.untagFindingHandler)
			r.Post("/put-resource/", svc.putResourceHandler)
			r.Post("/delete-resource/", svc.deleteResourceHandler)
			r.Post("/tag-resource/", svc.tagResourceHandler)
			r.Post("/untag-resource/", svc.untagResourceHandler)
		})
	})

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
	return r
}
