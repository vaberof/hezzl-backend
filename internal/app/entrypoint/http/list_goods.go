package http

import (
	"encoding/json"
	"github.com/vaberof/hezzl-backend/internal/app/entrypoint/http/views"
	"github.com/vaberof/hezzl-backend/internal/domain/good"
	"github.com/vaberof/hezzl-backend/pkg/http/protocols/apiv1"
	"net/http"
	"strconv"
	"time"
)

const (
	defaultLimit  = 10
	defaultOffset = 1
)

type listGoodsResponseBody struct {
	Meta  metaPayload        `json:"meta"`
	Goods []*listGoodPayload `json:"goods"`
}

type metaPayload struct {
	Total   int `json:"total"`
	Removed int `json:"removed"`
	Limit   int `json:"limit"`
	Offset  int `json:"offset"`
}

type listGoodPayload struct {
	Id          int64     `json:"id"`
	ProjectId   int64     `json:"projectId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Priority    int       `json:"priority"`
	Removed     bool      `json:"removed"`
	CreatedAt   time.Time `json:"createdAt"`
}

func (h *Handler) ListGoodsHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, request *http.Request) {
		var limit, offset int
		var err error

		limitStr := request.URL.Query().Get("limit")
		offsetStr := request.URL.Query().Get("offset")

		if limitStr == "" {
			limit = defaultLimit
		} else {
			limit, err = strconv.Atoi(limitStr)
			if err != nil {
				views.RenderJSON(rw, request, http.StatusInternalServerError, apiv1.Error(CodeInternalError, ErrMessageInternalServerError, apiv1.ErrorDescription{"details": "Failed to convert limit to int"}))

				return
			}
			if limit < 0 {
				views.RenderJSON(rw, request, http.StatusBadRequest, apiv1.Error(CodeBadRequest, ErrMessageInvalidRequestBody, apiv1.ErrorDescription{"details": "'limit' must not be negative"}))

				return
			}
		}

		if offsetStr == "" {
			offset = defaultOffset
		} else {
			offset, err = strconv.Atoi(offsetStr)
			if err != nil {
				views.RenderJSON(rw, request, http.StatusInternalServerError, apiv1.Error(CodeInternalError, ErrMessageInternalServerError, apiv1.ErrorDescription{"details": "Failed to convert offset to int"}))

				return
			}
			if offset < 0 {
				views.RenderJSON(rw, request, http.StatusBadRequest, apiv1.Error(CodeBadRequest, ErrMessageInvalidRequestBody, apiv1.ErrorDescription{"details": "'offset' must not be negative"}))

				return
			}
		}

		domainGoods, err := h.goodService.List(limit, offset)
		if err != nil {
			views.RenderJSON(rw, request, http.StatusInternalServerError, apiv1.Error(CodeInternalError, ErrMessageInternalServerError, apiv1.ErrorDescription{"details": "Failed to list goods"}))

			return
		}

		payload, _ := json.Marshal(&listGoodsResponseBody{
			Meta:  h.buildMetaPayload(domainGoods, limit, offset),
			Goods: h.buildListGoodPayloads(domainGoods),
		})

		views.RenderJSON(rw, request, http.StatusOK, apiv1.Success(payload))
	}
}

func (h *Handler) buildMetaPayload(domainGoods []*good.Good, limit, offset int) metaPayload {
	var meta metaPayload

	meta.Total = len(domainGoods)
	meta.Removed = h.countRemovedGoods(domainGoods)
	meta.Limit = limit
	meta.Offset = offset

	return meta
}

func (h *Handler) countRemovedGoods(domainGoods []*good.Good) (count int) {
	for _, domainGood := range domainGoods {
		if domainGood.Removed {
			count++
		}
	}
	return count
}

func (h *Handler) buildListGoodPayloads(domainGoods []*good.Good) []*listGoodPayload {
	goodPayloads := make([]*listGoodPayload, len(domainGoods))
	for i := range domainGoods {
		goodPayloads[i] = h.buildListGoodPayload(domainGoods[i])
	}
	return goodPayloads
}

func (h *Handler) buildListGoodPayload(domainGood *good.Good) *listGoodPayload {
	var goodPayload listGoodPayload

	goodPayload.Id = domainGood.Id.Int64()
	goodPayload.ProjectId = domainGood.ProjectId.Int64()
	goodPayload.Name = domainGood.Name.String()
	goodPayload.Description = domainGood.Description.String()
	goodPayload.Priority = domainGood.Priority.Int()
	goodPayload.Removed = domainGood.Removed.Bool()
	goodPayload.CreatedAt = domainGood.CreatedAt.Time()

	return &goodPayload
}
