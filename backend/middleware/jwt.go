package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	echoMiddleware "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"

	"github.com/chaitin/panda-wiki/config"
	"github.com/chaitin/panda-wiki/consts"
	"github.com/chaitin/panda-wiki/domain"
	"github.com/chaitin/panda-wiki/log"
	"github.com/chaitin/panda-wiki/repo/pg"
)

type JWTMiddleware struct {
	config         *config.Config
	jwtMiddleware  echo.MiddlewareFunc
	logger         *log.Logger
	userAccessRepo *pg.UserAccessRepository
}

func NewJWTMiddleware(config *config.Config, logger *log.Logger, userAccessRepo *pg.UserAccessRepository) *JWTMiddleware {
	jwtMiddleware := echoMiddleware.WithConfig(echoMiddleware.Config{
		SigningKey: []byte(config.Auth.JWT.Secret),
		ErrorHandler: func(c echo.Context, err error) error {
			logger.Error("jwt auth failed", log.Error(err))
			return c.JSON(http.StatusUnauthorized, domain.PWResponse{
				Success: false,
				Message: "Unauthorized",
			})
		},
	})
	return &JWTMiddleware{
		config:         config,
		jwtMiddleware:  jwtMiddleware,
		logger:         logger.WithModule("middleware.jwt"),
		userAccessRepo: userAccessRepo,
	}
}

func (m *JWTMiddleware) Authorize(next echo.HandlerFunc) echo.HandlerFunc {
	return m.jwtMiddleware(func(c echo.Context) error {
		// JWT authentication was successful, update access time
		if userID, ok := m.MustGetUserID(c); ok {
			c.Set("user_id", userID)
			m.userAccessRepo.UpdateAccessTime(userID)
		}
		return next(c)
	})
}

func (m *JWTMiddleware) ValidateUserRole(role consts.UserRole) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userID := c.Get("user_id").(string)

			valid, err := m.userAccessRepo.ValidateRole(userID, role)

			if err != nil || !valid {
				m.logger.Error("ValidateRole check", log.Any("user_id", userID), log.Any("valid", valid))
				return c.JSON(http.StatusForbidden, domain.PWResponse{
					Success: false,
					Message: "StatusForbidden ValidateRole",
				})
			}

			return next(c)
		}
	}
}

func (m *JWTMiddleware) ValidateKBUserPerm(perm consts.UserKBPermission) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			userID := c.Get("user_id").(string)

			kbId, _ := GetKbID(c)

			valid, err := m.userAccessRepo.ValidateKBPerm(kbId, userID, perm)
			if err != nil || !valid {
				if err != nil {
					m.logger.Error("ValidateKBUserPerm ValidateKBPerm failed", log.Error(err))
				} else {
					m.logger.Info("ValidateKBUserPerm ValidateKBPerm failed", log.String("kb_id", kbId), log.String("user_id", userID))
				}
				return c.JSON(http.StatusForbidden, domain.PWResponse{
					Success: false,
					Message: "Unauthorized ValidateKBPerm",
				})
			}

			return next(c)
		}
	}
}

func (m *JWTMiddleware) ValidateLicenseEdition(needEdition consts.LicenseEdition) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			edition, ok := c.Get("edition").(consts.LicenseEdition)
			if !ok {
				return c.JSON(http.StatusForbidden, domain.PWResponse{
					Success: false,
					Message: "Unauthorized ValidateLicenseEdition",
				})
			}

			if edition < needEdition {
				return c.JSON(http.StatusForbidden, domain.PWResponse{
					Success: false,
					Message: "Unauthorized ValidateLicenseEdition",
				})
			}

			return next(c)
		}
	}
}

func (m *JWTMiddleware) MustGetUserID(c echo.Context) (string, bool) {
	user, ok := c.Get("user").(*jwt.Token)
	if !ok || user == nil {
		return "", false
	}
	claims, ok := user.Claims.(jwt.MapClaims)
	if !ok {
		return "", false
	}
	id, ok := claims["id"].(string)
	return id, ok
}

func GetKbID(c echo.Context) (string, error) {
	switch c.Request().Method {
	case http.MethodGet, http.MethodDelete:

		kbId := c.QueryParam("kb_id")
		if kbId != "" {
			return kbId, nil
		}

		if strings.Contains(c.Request().URL.Path, "knowledge_base") {
			kbId = c.QueryParam("id")
			if kbId != "" {
				return kbId, nil
			}
		}
		return "", nil

	case http.MethodPost, http.MethodPatch, http.MethodPut:

		bodyBytes, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return "", err
		}

		c.Request().Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		var m map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &m); err == nil {
			if id, exists := m["kb_id"].(string); exists && id != "" {
				return id, nil
			}
			if strings.Contains(c.Request().URL.Path, "knowledge_base") {
				if id, exists := m["id"].(string); exists && id != "" {
					return id, nil
				}
			}
		}
		return "", nil
	default:
		return "", nil
	}
}
