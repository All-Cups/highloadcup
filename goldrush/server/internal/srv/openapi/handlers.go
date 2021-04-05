package openapi

import (
	"context"
	"errors"

	"github.com/Djarvur/allcups-itrally-2020-task/api/openapi/model"
	"github.com/Djarvur/allcups-itrally-2020-task/api/openapi/restapi/op"
	"github.com/Djarvur/allcups-itrally-2020-task/internal/app/game"
	"github.com/Djarvur/allcups-itrally-2020-task/internal/app/resource"
)

func (srv *server) HealthCheck(params op.HealthCheckParams) op.HealthCheckResponder {
	ctx, log := fromRequest(params.HTTPRequest)
	status, err := srv.app.HealthCheck(ctx)
	switch {
	default:
		return errHealthCheck(log, err, codeInternal)
	case err == nil:
		return op.NewHealthCheckOK().WithPayload(status)
	}
}

func (srv *server) GetBalance(params op.GetBalanceParams) op.GetBalanceResponder {
	ctx, log := fromRequest(params.HTTPRequest)
	if srv.cfg.OpGetBalanceRate != 0 && !srv.limitGetBalance.Allow() {
		return errGetBalance(log, errTooManyRequests, codeTooManyRequests)
	}

	balance, wallet, err := srv.app.Balance(ctx)
	switch {
	default:
		return errGetBalance(log, err, codeInternal)
	case err == nil:
		return op.NewGetBalanceOK().WithPayload(apiBalance(balance, wallet))
	}
}

func (srv *server) ListLicenses(params op.ListLicensesParams) op.ListLicensesResponder {
	ctx, log := fromRequest(params.HTTPRequest)
	if srv.cfg.OpListLicensesRate != 0 && !srv.limitListLicenses.Allow() {
		return errListLicenses(log, errTooManyRequests, codeTooManyRequests)
	}

	licenses, err := srv.app.Licenses(ctx)
	switch {
	default:
		return errListLicenses(log, err, codeInternal)
	case errors.Is(err, resource.ErrRPCInternal):
		return errListLicenses(log, err, codeBadGateway)
	case errors.Is(err, resource.ErrRPCTimeout):
		return errListLicenses(log, err, codeGatewayTimeout)
	case err == nil:
		return op.NewListLicensesOK().WithPayload(apiLicenseList(licenses))
	}
}

func (srv *server) IssueLicense(params op.IssueLicenseParams) op.IssueLicenseResponder {
	ctx, log := fromRequest(params.HTTPRequest)
	if srv.cfg.OpIssueLicenseRate != 0 && !srv.limitIssueLicense.Allow() {
		return errIssueLicense(log, errTooManyRequests, codeTooManyRequests)
	}

	license, err := srv.app.IssueLicense(ctx, appWallet(params.Args))
	switch {
	default:
		return errIssueLicense(log, err, codeInternal)
	case errors.Is(err, game.ErrActiveLicenseLimit):
		return errIssueLicense(log, err, codeActiveLicenseLimit)
	case errors.Is(err, game.ErrBogusCoin):
		return errIssueLicense(log, err, codePaymentRequired)
	case errors.Is(err, resource.ErrRPCInternal):
		return errIssueLicense(log, err, codeBadGateway)
	case errors.Is(err, resource.ErrRPCTimeout):
		return errIssueLicense(log, err, codeGatewayTimeout)
	case err == nil:
		return op.NewIssueLicenseOK().WithPayload(apiLicense(license))
	}
}

func (srv *server) ExploreArea(params op.ExploreAreaParams) op.ExploreAreaResponder {
	ctx, log := fromRequest(params.HTTPRequest)
	if srv.cfg.OpExploreAreaRate != 0 && !srv.limitExploreArea.Allow() {
		return errExploreArea(log, errTooManyRequests, codeTooManyRequests)
	}
	ctx, cancel := context.WithTimeout(ctx, srv.cfg.OpExploreAreaTimeout)
	defer cancel()

	count, err := srv.app.ExploreArea(ctx, appArea(params.Args))
	switch {
	default:
		return errExploreArea(log, err, codeInternal)
	case ctx.Err() != nil:
		return errExploreArea(log, errServiceUnavailable, codeServiceUnavailable)
	case errors.Is(err, game.ErrWrongCoord):
		return errExploreArea(log, err, codeWrongCoord)
	case err == nil:
		return op.NewExploreAreaOK().WithPayload(&model.Report{
			Area:   params.Args,
			Amount: model.Amount(count),
		})
	}
}

func (srv *server) Dig(params op.DigParams) op.DigResponder {
	ctx, log := fromRequest(params.HTTPRequest)
	if srv.cfg.OpDigRate != 0 && !srv.limitDig.Allow() {
		return errDig(log, errTooManyRequests, codeTooManyRequests)
	}
	ctx, cancel := context.WithTimeout(ctx, srv.cfg.OpDigTimeout)
	defer cancel()

	coord := appCoord(params.Args)
	treasure, err := srv.app.Dig(ctx, int(*params.Args.LicenseID), coord)
	switch {
	default:
		return errDig(log, err, codeInternal)
	case ctx.Err() != nil:
		return errDig(log, errServiceUnavailable, codeServiceUnavailable)
	case errors.Is(err, game.ErrNoSuchLicense):
		return errDig(log, err, codeForbidden)
	case errors.Is(err, game.ErrWrongCoord):
		return errDig(log, err, codeWrongCoord)
	case errors.Is(err, game.ErrWrongDepth):
		return errDig(log, err, codeWrongDepth)
	case err == nil && treasure == "":
		return errDig(log, game.ErrNoThreasure, codeNotFound)
	case err == nil:
		return op.NewDigOK().WithPayload(apiTreasureList(treasure))
	}
}

func (srv *server) Cash(params op.CashParams) op.CashResponder {
	ctx, log := fromRequest(params.HTTPRequest)
	if srv.cfg.OpCashRate != 0 && !srv.limitCash.Allow() {
		return errCash(log, errTooManyRequests, codeTooManyRequests)
	}
	if srv.inPercent(srv.cfg.OpCashPercentFail) {
		return errCash(log, errServiceUnavailable, codeServiceUnavailable)
	}

	wallet, err := srv.app.Cash(ctx, string(params.Args))
	switch {
	default:
		return errCash(log, err, codeInternal)
	case errors.Is(err, game.ErrWrongCoord):
		return errCash(log, err, codeWrongCoord)
	case errors.Is(err, game.ErrNotDigged):
		return errCash(log, err, codeNotDigged)
	case errors.Is(err, game.ErrNoThreasure):
		return errCash(log, err, codeNotFound)
	case err == nil:
		return op.NewCashOK().WithPayload(apiWallet(wallet))
	}
}
