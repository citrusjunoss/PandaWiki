package share

import (
	"context"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sbzhu/weworkapi_golang/wxbizmsgcrypt"

	"github.com/chaitin/panda-wiki/domain"
	"github.com/chaitin/panda-wiki/handler"
	"github.com/chaitin/panda-wiki/log"
	"github.com/chaitin/panda-wiki/pkg/bot/wechatservice"
	"github.com/chaitin/panda-wiki/usecase"
)

type ShareAppHandler struct {
	*handler.BaseHandler
	logger  *log.Logger
	usecase *usecase.AppUsecase
}

func NewShareAppHandler(
	e *echo.Echo,
	baseHandler *handler.BaseHandler,
	logger *log.Logger,
	usecase *usecase.AppUsecase,
) *ShareAppHandler {
	h := &ShareAppHandler{
		BaseHandler: baseHandler,
		logger:      logger.WithModule("handler.share.app"),
		usecase:     usecase,
	}

	share := e.Group("share/v1/app",
		func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				c.Response().Header().Set("Access-Control-Allow-Origin", "*")
				c.Response().Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
				c.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type, Origin, Accept")
				if c.Request().Method == "OPTIONS" {
					return c.NoContent(http.StatusOK)
				}
				return next(c)
			}
		})
	share.GET("/web/info", h.GetWebAppInfo)
	share.GET("/widget/info", h.GetWidgetAppInfo)

	share.GET("/wechat/app", h.VerifyUrlWechatApp)
	share.POST("/wechat/app", h.WechatHandlerApp)

	// 微信客服
	share.GET("/wechat/service", h.VerifyUrlWechatService)
	share.POST("/wechat/service", h.WechatHandlerService)

	return h
}

// GetAppInfo
//
//	@Summary		GetAppInfo
//	@Description	GetAppInfo
//	@Tags			share_app
//	@Accept			json
//	@Produce		json
//	@Param			X-KB-ID	header		string	true	"kb id"
//	@Success		200		{object}	domain.Response
//	@Router			/share/v1/app/web/info [get]
func (h *ShareAppHandler) GetWebAppInfo(c echo.Context) error {
	kbID := c.Request().Header.Get("X-KB-ID")
	if kbID == "" {
		return h.NewResponseWithError(c, "kb_id is required", nil)
	}
	appInfo, err := h.usecase.GetWebAppInfo(c.Request().Context(), kbID)
	if err != nil {
		return h.NewResponseWithError(c, err.Error(), err)
	}
	return h.NewResponseWithData(c, appInfo)
}

// GetWidgetAppInfo
//
//	@Summary		GetWidgetAppInfo
//	@Description	GetWidgetAppInfo
//	@Tags			share_app
//	@Accept			json
//	@Produce		json
//	@Param			X-KB-ID	header		string	true	"kb id"
//	@Success		200		{object}	domain.Response
//	@Router			/share/v1/app/widget/info [get]
func (h *ShareAppHandler) GetWidgetAppInfo(c echo.Context) error {
	kbID := c.Request().Header.Get("X-KB-ID")
	if kbID == "" {
		return h.NewResponseWithError(c, "kb_id is required", nil)
	}
	appInfo, err := h.usecase.GetWidgetAppInfo(c.Request().Context(), kbID)
	if err != nil {
		return h.NewResponseWithError(c, err.Error(), err)
	}
	return h.NewResponseWithData(c, appInfo)
}

func (h *ShareAppHandler) VerifyUrlWechatApp(c echo.Context) error {
	signature := c.QueryParam("msg_signature")
	timestamp := c.QueryParam("timestamp")
	nonce := c.QueryParam("nonce")
	echoStr := c.QueryParam("echostr")

	kbID := c.Request().Header.Get("X-KB-ID")

	if kbID == "" {
		return h.NewResponseWithError(c, "kb_id is required", nil)
	}

	if signature == "" || timestamp == "" || nonce == "" || echoStr == "" {
		return h.NewResponseWithError(
			c, "verify wechat params failed", nil,
		)
	}

	ctx := c.Request().Context()

	req, err := h.usecase.VerifyUrlWechatAPP(ctx, signature, timestamp, nonce, echoStr, kbID)
	if err != nil {
		return h.NewResponseWithError(c, "VerifyURL failed", err)
	}

	// success
	return c.String(http.StatusOK, string(req))
}

func (h *ShareAppHandler) WechatHandlerApp(c echo.Context) error {
	signature := c.QueryParam("msg_signature")
	timestamp := c.QueryParam("timestamp")
	nonce := c.QueryParam("nonce")

	kbID := c.Request().Header.Get("X-KB-ID")
	if kbID == "" {
		return h.NewResponseWithError(c, "kb_id is required", nil)
	}

	remoteIP := ""
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		h.logger.Error("get request failed", log.Error(err))
		return h.NewResponseWithError(c, "Internal Server Error", err)
	}
	defer c.Request().Body.Close()

	ctx := c.Request().Context()

	immediateResponse, err := h.usecase.SendImmediateResponse(ctx, signature, timestamp, nonce, body, kbID)
	if err != nil {
		h.logger.Error("send response failed", log.Error(err))
		return h.NewResponseWithError(c, "Failed to send immediate response", err)
	}

	go func(signature, timestamp, nonce string, body []byte, KbId string, remoteIP string) {
		ctx := context.Background()
		err := h.usecase.Wechat(ctx, signature, timestamp, nonce, body, KbId, remoteIP)
		if err != nil {
			h.logger.Error("wechat async failed")
		}
	}(signature, timestamp, nonce, body, kbID, remoteIP)

	return c.XMLBlob(http.StatusOK, []byte(immediateResponse))
}

// 验证微信客服消息
func (h *ShareAppHandler) VerifyUrlWechatService(c echo.Context) error {
	signature := c.QueryParam("msg_signature")
	timestamp := c.QueryParam("timestamp")
	nonce := c.QueryParam("nonce")
	echoStr := c.QueryParam("echostr")

	kbID := c.Request().Header.Get("X-KB-ID")

	if kbID == "" {
		return h.NewResponseWithError(c, "kb_id is required", nil)
	}

	if signature == "" || timestamp == "" || nonce == "" || echoStr == "" {
		return h.NewResponseWithError(
			c, "verify wechat service params failed", nil,
		)
	}

	ctx := c.Request().Context()

	req, err := h.usecase.VerifyUrlWechatService(ctx, signature, timestamp, nonce, echoStr, kbID)
	if err != nil {
		return h.NewResponseWithError(c, "VerifyURL_Service failed", err)
	}

	// success
	return c.String(http.StatusOK, string(req))
}

// 处理用户的请求，客服回复消息
func (h *ShareAppHandler) WechatHandlerService(c echo.Context) error {
	signature := c.QueryParam("msg_signature")
	timestamp := c.QueryParam("timestamp")
	nonce := c.QueryParam("nonce")

	kbID := c.Request().Header.Get("X-KB-ID")

	if kbID == "" {
		return h.NewResponseWithError(c, "kb_id is required", nil)
	}

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		h.logger.Error("get request failed", log.Error(err))
		return h.NewResponseWithError(c, "Internal Server Error", err)
	}
	defer c.Request().Body.Close()

	ctx := c.Request().Context()

	// 查找数据库，找到对应的app配置
	appInfo, err := h.usecase.GetAppDetailByKBIDAndAppType(ctx, kbID, domain.AppTypeWechatServiceBot)

	if err != nil {
		h.logger.Error("GetAppDetailByKBIDAndAppType failed", log.Error(err))
		return h.NewResponseWithError(c, "GetAppDetailByKBIDAndAppType failed", err)
	}
	if appInfo.Settings.WeChatServiceIsEnabled != nil && !*appInfo.Settings.WeChatServiceIsEnabled {
		return h.NewResponseWithError(c, "wechat service bot is not enabled", nil)
	}

	// 创建一个wechat service对象
	wechatServiceConf, err := h.usecase.NewWechatServiceConfig(context.Background(), appInfo, kbID)

	h.logger.Info("wechat service config", log.Any("wechat service config", wechatServiceConf))

	if err != nil {
		return h.NewResponseWithError(c, "New WechatServiceConfig failed", err)
	}

	// 解密消息
	wxCrypt := wxbizmsgcrypt.NewWXBizMsgCrypt(wechatServiceConf.Token, wechatServiceConf.EncodingAESKey, wechatServiceConf.CorpID, wxbizmsgcrypt.XmlType)
	decryptMsg, errCode := wxCrypt.DecryptMsg(signature, timestamp, nonce, body)
	if errCode != nil {
		h.logger.Error("DecryptMsg failed", log.Any("decryptMsg err", errCode))
		return h.NewResponseWithError(c, "DecryptMsg failed", nil)
	}
	// 解密成功了
	h.logger.Info("收到了用户发送或者事件的原始消息（解密之后）", log.Any("decryptMsg", decryptMsg))

	// 反序列化
	msg, err := wechatServiceConf.UnmarshalMsg(decryptMsg)
	if err != nil {
		h.logger.Error("UnmarshalMsg failed", log.Error(err))
		return h.NewResponseWithError(c, "UnmarshalMsg failed", err)
	}

	// 这个是用户发送的消息
	h.logger.Info("收到用户消息（提问的问题）用户,事件类型", msg.ToUserName, msg.Event)

	// 没有问题则开启异步拉去用户发给微信客服的消息，并且处理用户的消息
	go func(WechatServiceConf *wechatservice.WechatServiceConfig, msg *wechatservice.WeixinUserAskMsg, kbID string) {
		ctx := context.Background()
		err := h.usecase.WechatService(ctx, msg, kbID, WechatServiceConf)
		if err != nil {
			h.logger.Error("wechat async failed", log.Any("Wechat_Service", err))
		}
	}(wechatServiceConf, msg, kbID)

	// 先响应
	return c.JSON(http.StatusOK, "success") // 不会发送给用户，会被微信服务器忽略
}
