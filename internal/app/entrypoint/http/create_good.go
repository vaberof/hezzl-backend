package http

import (
	"encoding/json"
	"github.com/go-chi/render"
	"github.com/vaberof/hezzl-backend/internal/app/entrypoint/http/views"
	"github.com/vaberof/hezzl-backend/pkg/domain"
	"github.com/vaberof/hezzl-backend/pkg/http/protocols/apiv1"
	"net/http"
	"strconv"
	"time"
)

type createGoodRequestBody struct {
	Name string `json:"name"`
}

func (c *createGoodRequestBody) Bind(req *http.Request) error {
	return nil
}

type createGoodResponseBody struct {
	Id          int64     `json:"id"`
	ProjectId   int64     `json:"projectId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Priority    int       `json:"priority"`
	Removed     bool      `json:"removed"`
	CreatedAt   time.Time `json:"createdAt"`
}

func (h *Handler) CreateGoodHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, request *http.Request) {
		projectIdStr := request.URL.Query().Get("projectId")
		if projectIdStr == "" {
			views.RenderJSON(rw, request, http.StatusBadRequest, apiv1.Error(CodeBadRequest, ErrMessageInvalidRequestBody, apiv1.ErrorDescription{"details": "Missing required query parameter 'projectId'"}))

			return
		}

		createGoodReqBody := &createGoodRequestBody{}
		if err := render.Bind(request, createGoodReqBody); err != nil {
			views.RenderJSON(rw, request, http.StatusBadRequest, apiv1.Error(CodeBadRequest, ErrMessageInvalidRequestBody, apiv1.ErrorDescription{"details": "Invalid request body"}))

			return
		}

		projectId, err := strconv.ParseInt(projectIdStr, 10, 64)
		if err != nil {
			views.RenderJSON(rw, request, http.StatusInternalServerError, apiv1.Error(CodeInternalError, ErrMessageInternalServerError, apiv1.ErrorDescription{"details": "Failed to convert projectId to int"}))

			return
		}

		domainGood, err := h.goodService.Create(domain.ProjectId(projectId), domain.GoodName(createGoodReqBody.Name))
		if err != nil {
			views.RenderJSON(rw, request, http.StatusInternalServerError, apiv1.Error(CodeInternalError, ErrMessageInternalServerError, apiv1.ErrorDescription{"details": "Failed to create a new good"}))

			return
		}

		payload, _ := json.Marshal(&createGoodResponseBody{
			Id:          domainGood.Id.Int64(),
			ProjectId:   domainGood.ProjectId.Int64(),
			Name:        domainGood.Name.String(),
			Description: domainGood.Description.String(),
			Priority:    domainGood.Priority.Int(),
			Removed:     domainGood.Removed.Bool(),
			CreatedAt:   domainGood.CreatedAt.Time(),
		})

		views.RenderJSON(rw, request, http.StatusOK, apiv1.Success(payload))
	}
}
