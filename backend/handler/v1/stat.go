package v1

import (
	"github.com/labstack/echo/v4"

	"github.com/chaitin/panda-wiki/consts"
	"github.com/chaitin/panda-wiki/handler"
	"github.com/chaitin/panda-wiki/log"
	"github.com/chaitin/panda-wiki/middleware"
	"github.com/chaitin/panda-wiki/usecase"
)

type StatHandler struct {
	*handler.BaseHandler
	usecase *usecase.StatUseCase
	auth    middleware.AuthMiddleware
	logger  *log.Logger
}

func NewStatHandler(baseHandler *handler.BaseHandler, echo *echo.Echo, usecase *usecase.StatUseCase, logger *log.Logger, auth middleware.AuthMiddleware) *StatHandler {
	h := &StatHandler{
		BaseHandler: baseHandler,
		usecase:     usecase,
		auth:        auth,
		logger:      logger.WithModule("handler.v1.stat"),
	}

	group := echo.Group("/api/v1/stat", h.auth.Authorize, auth.ValidateKBUserPerm(consts.UserKBPermissionDataOperate))
	group.GET("/hot_pages", h.GetHotPages)
	group.GET("/referer_hosts", h.GetRefererHosts)
	group.GET("/browsers", h.GetBrowsers)
	group.GET("/count", h.GetCount)
	// instant count (30min, every 1min)
	group.GET("/instant_count", h.GetInstantCount)
	// instant pages (latest 10 pages)
	group.GET("/instant_pages", h.GetInstantPages)
	// geo (24h)
	group.GET("/geo_count", h.GetGeoCount)
	// conversation (24h)
	group.GET("/conversation_distribution", h.GetConversationDistribution)
	return h
}

// GetHotPages get hot pages
//
//	@Summary		GetHotPages
//	@Description	GetHotPages
//	@Tags			stat
//	@Accept			json
//	@Produce		json
//	@Param			kb_id	query		string	true	"kb_id"
//	@Success		200		{object}	domain.Response
//	@Router			/api/v1/stat/hot_pages [get]
func (h *StatHandler) GetHotPages(c echo.Context) error {
	kbID := c.QueryParam("kb_id")
	if kbID == "" {
		return h.NewResponseWithError(c, "kb_id is required", nil)
	}
	hotPages, err := h.usecase.GetHotPages(c.Request().Context(), kbID)
	if err != nil {
		return h.NewResponseWithError(c, "get hot pages failed", err)
	}
	return h.NewResponseWithData(c, hotPages)
}

// GetRefererHosts get hot referer hosts
//
//	@Summary		GetRefererHosts
//	@Description	GetRefererHosts
//	@Tags			stat
//	@Accept			json
//	@Produce		json
//	@Param			kb_id	query		string	true	"kb_id"
//	@Success		200		{object}	domain.Response
//	@Router			/api/v1/stat/referer_hosts [get]
func (h *StatHandler) GetRefererHosts(c echo.Context) error {
	kbID := c.QueryParam("kb_id")
	if kbID == "" {
		return h.NewResponseWithError(c, "kb_id is required", nil)
	}
	refererHosts, err := h.usecase.GetHotRefererHosts(c.Request().Context(), kbID)
	if err != nil {
		return h.NewResponseWithError(c, "get hot referer hosts failed", err)
	}
	return h.NewResponseWithData(c, refererHosts)
}

// GetBrowsers get hot browsers
//
//	@Summary		GetBrowsers
//	@Description	GetBrowsers
//	@Tags			stat
//	@Accept			json
//	@Produce		json
//	@Param			kb_id	query		string	true	"kb_id"
//	@Success		200		{object}	domain.Response(data=domain.HotBrowserResp)
//	@Router			/api/v1/stat/browsers [get]
func (h *StatHandler) GetBrowsers(c echo.Context) error {
	kbID := c.QueryParam("kb_id")
	if kbID == "" {
		return h.NewResponseWithError(c, "kb_id is required", nil)
	}
	hotBrowsers, err := h.usecase.GetHotBrowsers(c.Request().Context(), kbID)
	if err != nil {
		return h.NewResponseWithError(c, "get hot browsers failed", err)
	}
	return h.NewResponseWithData(c, hotBrowsers)
}

// GetCount get count
//
//	@Summary		GetCount
//	@Description	GetCount
//	@Tags			stat
//	@Accept			json
//	@Produce		json
//	@Param			kb_id	query		string	true	"kb_id"
//	@Success		200		{object}	domain.Response
//	@Router			/api/v1/stat/count [get]
func (h *StatHandler) GetCount(c echo.Context) error {
	kbID := c.QueryParam("kb_id")
	if kbID == "" {
		return h.NewResponseWithError(c, "kb_id is required", nil)
	}
	count, err := h.usecase.GetCount(c.Request().Context(), kbID)
	if err != nil {
		return h.NewResponseWithError(c, "get count failed", err)
	}
	return h.NewResponseWithData(c, count)
}

// GetInstantCount get instant count
//
//	@Summary		GetInstantCount
//	@Description	GetInstantCount
//	@Tags			stat
//	@Accept			json
//	@Produce		json
//	@Param			kb_id	query		string	true	"kb_id"
//	@Success		200		{object}	domain.Response
//	@Router			/api/v1/stat/instant_count [get]
func (h *StatHandler) GetInstantCount(c echo.Context) error {
	kbID := c.QueryParam("kb_id")
	if kbID == "" {
		return h.NewResponseWithError(c, "kb_id is required", nil)
	}
	count, err := h.usecase.GetInstantCount(c.Request().Context(), kbID)
	if err != nil {
		return h.NewResponseWithError(c, "get instant count failed", err)
	}
	return h.NewResponseWithData(c, count)
}

// GetInstantPages get instant pages
//
//	@Summary		GetInstantPages
//	@Description	GetInstantPages
//	@Tags			stat
//	@Accept			json
//	@Produce		json
//	@Param			kb_id	query		string	true	"kb_id"
//	@Success		200		{object}	domain.Response
//	@Router			/api/v1/stat/instant_pages [get]
func (h *StatHandler) GetInstantPages(c echo.Context) error {
	kbID := c.QueryParam("kb_id")
	if kbID == "" {
		return h.NewResponseWithError(c, "kb_id is required", nil)
	}
	pages, err := h.usecase.GetInstantPages(c.Request().Context(), kbID)
	if err != nil {
		return h.NewResponseWithError(c, "get instant pages failed", err)
	}
	return h.NewResponseWithData(c, pages)
}

// GetGeoCount get geo count
//
//	@Summary		GetGeoCount
//	@Description	GetGeoCount
//	@Tags			stat
//	@Accept			json
//	@Produce		json
//	@Param			kb_id	query		string	true	"kb_id"
//	@Success		200		{object}	domain.Response
//	@Router			/api/v1/stat/geo_count [get]
func (h *StatHandler) GetGeoCount(c echo.Context) error {
	kbID := c.QueryParam("kb_id")
	if kbID == "" {
		return h.NewResponseWithError(c, "kb_id is required", nil)
	}
	geoCount, err := h.usecase.GetGeoCount(c.Request().Context(), kbID)
	if err != nil {
		return h.NewResponseWithError(c, "get geo count failed", err)
	}
	return h.NewResponseWithData(c, geoCount)
}

// GetConversationDistribution
//
//	@Summary		GetConversationDistribution
//	@Description	GetConversationDistribution
//	@Tags			stat
//	@Accept			json
//	@Produce		json
//	@Param			kb_id	query		string	true	"kb_id"
//	@Success		200		{object}	domain.Response
//	@Router			/api/v1/stat/conversation_distribution [get]
func (h *StatHandler) GetConversationDistribution(c echo.Context) error {
	kbID := c.QueryParam("kb_id")
	if kbID == "" {
		return h.NewResponseWithError(c, "kb_id is required", nil)
	}
	distribution, err := h.usecase.GetConversationDistribution(c.Request().Context(), kbID)
	if err != nil {
		return h.NewResponseWithError(c, "get conversation distribution failed", err)
	}
	return h.NewResponseWithData(c, distribution)
}
