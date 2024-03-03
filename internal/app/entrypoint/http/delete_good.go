package http

import (
	"encoding/json"
	"errors"
	"github.com/vaberof/hezzl-backend/internal/app/entrypoint/http/views"
	"github.com/vaberof/hezzl-backend/internal/domain/good"
	"github.com/vaberof/hezzl-backend/pkg/domain"
	"github.com/vaberof/hezzl-backend/pkg/http/protocols/apiv1"
	"net/http"
	"strconv"
)

type deleteGoodResponseBody struct {
	Id        int64 `json:"id"`
	ProjectId int64 `json:"projectId"`
	Removed   bool  `json:"removed"`
}

func (h *Handler) DeleteGoodHandler() http.HandlerFunc {
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

		domainGood, err := h.goodService.Delete(domain.GoodId(goodId), domain.ProjectId(projectId))
		if err != nil {
			if errors.Is(err, good.ErrGoodNotFound) {
				views.RenderJSON(rw, request, http.StatusNotFound, apiv1.Error(CodeNotFound, ErrMessageGoodNotFound, apiv1.ErrorDescription{"details": "Good is not found"}))
			} else {
				views.RenderJSON(rw, request, http.StatusInternalServerError, apiv1.Error(CodeInternalError, ErrMessageInternalServerError, apiv1.ErrorDescription{"details": "Failed to delete a good"}))
			}

			return
		}

		payload, _ := json.Marshal(&deleteGoodResponseBody{
			Id:        goodId,
			ProjectId: projectId,
			Removed:   domainGood.Removed.Bool(),
		})

		views.RenderJSON(rw, request, http.StatusOK, apiv1.Success(payload))
	}
}
