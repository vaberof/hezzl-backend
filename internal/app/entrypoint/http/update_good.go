package http

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/render"
	"github.com/vaberof/hezzl-backend/internal/app/entrypoint/http/views"
	"github.com/vaberof/hezzl-backend/internal/domain/good"
	"github.com/vaberof/hezzl-backend/pkg/domain"
	"github.com/vaberof/hezzl-backend/pkg/http/protocols/apiv1"
	"net/http"
	"strconv"
	"time"
)

type updateGoodRequestBody struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

func (u *updateGoodRequestBody) Bind(req *http.Request) error {
	return nil
}

type updateGoodResponseBody struct {
	Id          int64     `json:"id"`
	ProjectId   int64     `json:"projectId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Priority    int       `json:"priority"`
	Removed     bool      `json:"removed"`
	CreatedAt   time.Time `json:"createdAt"`
}

func (h *Handler) UpdateGoodHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, request *http.Request) {
		goodIdStr := request.URL.Query().Get("id")
		if goodIdStr == "" {
			views.RenderJSON(rw, request, http.StatusBadRequest, apiv1.Error(CodeBadRequest, ErrMessageInvalidRequestBody, apiv1.ErrorDescription{"details": "Missing required query parameter 'id'"}))

			return
		}

		projectIdStr := request.URL.Query().Get("projectId")
		if projectIdStr == "" {
			views.RenderJSON(rw, request, http.StatusBadRequest, apiv1.Error(CodeBadRequest, ErrMessageInvalidRequestBody, apiv1.ErrorDescription{"details": "Missing required query parameter 'projectId'"}))

			return
		}

		updateGoodReqBody := &updateGoodRequestBody{}
		if err := render.Bind(request, updateGoodReqBody); err != nil {

			views.RenderJSON(rw, request, http.StatusBadRequest, apiv1.Error(CodeBadRequest, ErrMessageInvalidRequestBody, apiv1.ErrorDescription{"details": "Invalid request body"}))

			return
		}

		goodId, err := strconv.ParseInt(goodIdStr, 10, 64)
		if err != nil {
			views.RenderJSON(rw, request, http.StatusInternalServerError, apiv1.Error(CodeInternalError, ErrMessageInternalServerError, apiv1.ErrorDescription{"details": "Failed to convert id to int"}))

			return
		}

		projectId, err := strconv.ParseInt(projectIdStr, 10, 64)
		if err != nil {
			views.RenderJSON(rw, request, http.StatusInternalServerError, apiv1.Error(CodeInternalError, ErrMessageInternalServerError, apiv1.ErrorDescription{"details": "Failed to convert projectId to int"}))

			return
		}

		var goodDescription *domain.GoodDescription
		if updateGoodReqBody.Description != nil {
			domainDescription := domain.GoodDescription(*updateGoodReqBody.Description)
			goodDescription = &domainDescription
		}

		domainGood, err := h.goodService.Update(domain.GoodId(goodId), domain.ProjectId(projectId), domain.GoodName(updateGoodReqBody.Name), goodDescription)
		if err != nil {
			if errors.Is(err, good.ErrGoodNotFound) {
				views.RenderJSON(rw, request, http.StatusNotFound, apiv1.Error(CodeNotFound, ErrMessageGoodNotFound, apiv1.ErrorDescription{"details": "Good is not found"}))
			} else {
				views.RenderJSON(rw, request, http.StatusInternalServerError, apiv1.Error(CodeInternalError, ErrMessageInternalServerError, apiv1.ErrorDescription{"details": "Failed to update a good"}))
			}

			return
		}

		payload, _ := json.Marshal(&updateGoodResponseBody{
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
