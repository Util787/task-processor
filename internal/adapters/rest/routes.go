package rest

import "net/http"

func (h *Handler) InitRoutes() http.Handler {
	mux := http.NewServeMux()

	v1 := "/api/v1"
	mux.Handle(v1+"/health", basicMiddleware(h.log, validateMethodMiddleware(http.MethodGet, http.HandlerFunc(h.getHealth))))
	mux.Handle(v1+"/enqueue", basicMiddleware(h.log, validateMethodMiddleware(http.MethodPost, http.HandlerFunc(h.enqueueTask))))
	mux.Handle(v1+"/tasks/", basicMiddleware(h.log, validateMethodMiddleware(http.MethodGet, http.HandlerFunc(h.getTaskState))))

	return mux
}
