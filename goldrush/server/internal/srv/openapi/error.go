//go:generate gobin -m -run github.com/cheekybits/genny -in=$GOFILE -out=gen.$GOFILE gen "HealthCheck=Cash,Dig,ExploreArea,GetBalance,IssueLicense,ListLicenses"
//go:generate sed -i -e "\\,^//go:generate,d" gen.$GOFILE

package openapi

import (
	"net/http"

	"github.com/go-openapi/swag"

	"github.com/Djarvur/allcups-itrally-2020-task/api/openapi/model"
	"github.com/Djarvur/allcups-itrally-2020-task/api/openapi/restapi/op"
	"github.com/Djarvur/allcups-itrally-2020-task/pkg/def"
)

func errHealthCheck(log Log, err error, code errCode) op.HealthCheckResponder {
	if code.status < http.StatusInternalServerError {
		log.Info("client error", def.LogHTTPStatus, code.status, "code", code.extra, "err", err)
	} else {
		log.PrintErr("server error", def.LogHTTPStatus, code.status, "code", code.extra, "err", err)
	}

	msg := err.Error()
	if code.status == http.StatusInternalServerError { // Do no expose details about internal errors.
		msg = "internal error" //nolint:goconst // Duplicated by go:generate.
	}

	return op.NewHealthCheckDefault(code.status).WithPayload(&model.Error{
		Code:    swag.Int32(code.extra),
		Message: swag.String(msg),
	})
}
