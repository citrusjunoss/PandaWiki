package consts

import (
	"math"

	"github.com/labstack/echo/v4"
)

type contextKey string

const ContextKeyEdition contextKey = "edition"

type LicenseEdition int32

const (
	LicenseEditionFree        LicenseEdition = 0
	LicenseEditionContributor LicenseEdition = 1
	LicenseEditionEnterprise  LicenseEdition = 2
)

func GetLicenseEdition(c echo.Context) LicenseEdition {
	edition, _ := c.Get("edition").(LicenseEdition)
	return edition
}

func (e LicenseEdition) GetMaxAuth() int {
	switch e {
	case LicenseEditionFree:
		return 0
	case LicenseEditionContributor:
		return 10
	case LicenseEditionEnterprise:
		return math.MaxInt
	default:
		return 0
	}
}
