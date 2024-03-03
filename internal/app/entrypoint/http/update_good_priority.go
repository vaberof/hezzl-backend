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
)

type updateGoodPriorityRequestBody struct {
	NewPriority int `json:"newPriority"`
}

func (u *updateGoodPriorityRequestBody) Bind(req *http.Request) error {
	return nil
}

type updateGoodPriorityResponseBody struct {
	Priorities []*updateGoodPriorityPayload `json:"priorities"`
}

type updateGoodPriorityPayload struct {
	Id       int64 `json:"id"`
	Priority int   `json:"priority"`
}

func (h *Handler) UpdateGoodPriorityHandler() http.HandlerFunc {
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

		updateGoodPriorityReqBody := &updateGoodPriorityRequestBody{}
		if err := render.Bind(request, updateGoodPriorityReqBody); err != nil {
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

		domainGoods, err := h.goodService.ChangePriority(domain.GoodId(goodId), domain.ProjectId(projectId), domain.GoodPriority(updateGoodPriorityReqBody.NewPriority))
		if err != nil {
			if errors.Is(err, good.ErrGoodNotFound) {
				views.RenderJSON(rw, request, http.StatusNotFound, apiv1.Error(CodeNotFound, ErrMessageGoodNotFound, apiv1.ErrorDescription{"details": "Good is not found"}))
			} else {
				views.RenderJSON(rw, request, http.StatusInternalServerError, apiv1.Error(CodeInternalError, ErrMessageInternalServerError, apiv1.ErrorDescription{"details": "Failed to update a good"}))
			}

			return
		}

		payload, _ := json.Marshal(&updateGoodPriorityResponseBody{
			Priorities: h.buildUpdateGoodPriorityPayloads(domainGoods),
		})

		views.RenderJSON(rw, request, http.StatusOK, apiv1.Success(payload))
	}
}

func (h *Handler) buildUpdateGoodPriorityPayloads(domainGoods []*good.Good) []*updateGoodPriorityPayload {
	goodPriorityPayloads := make([]*updateGoodPriorityPayload, len(domainGoods))
	for i := range domainGoods {
		goodPriorityPayloads[i] = h.buildUpdateGoodPriorityPayload(domainGoods[i])
	}
	return goodPriorityPayloads
}

func (h *Handler) buildUpdateGoodPriorityPayload(domainGood *good.Good) *updateGoodPriorityPayload {
	var goodPriorityPayload updateGoodPriorityPayload

	goodPriorityPayload.Id = domainGood.Id.Int64()
	goodPriorityPayload.Priority = domainGood.Priority.Int()

	return &goodPriorityPayload
}
