package openapi

import (
	"github.com/go-openapi/swag"

	"github.com/Djarvur/allcups-itrally-2020-task/api/openapi/model"
	"github.com/Djarvur/allcups-itrally-2020-task/internal/app/game"
)

func appWallet(ms model.Wallet) []int {
	vs := make([]int, len(ms))
	for i := range ms {
		vs[i] = int(ms[i])
	}
	return vs
}

func apiWallet(vs []int) model.Wallet {
	ms := make(model.Wallet, len(vs))
	for i := range vs {
		ms[i] = uint32(vs[i])
	}
	return ms
}

func apiBalance(balance int, wallet []int) *model.Balance {
	return &model.Balance{
		Balance: swag.Uint32(uint32(balance)),
		Wallet:  apiWallet(wallet),
	}
}

func apiLicense(v game.License) *model.License {
	return &model.License{
		ID:         swag.Int64(int64(v.ID)),
		DigAllowed: model.Amount(v.DigAllowed),
		DigUsed:    model.Amount(v.DigUsed),
	}
}

func apiLicenseList(vs []game.License) model.LicenseList {
	ms := make(model.LicenseList, len(vs))
	for i := range vs {
		ms[i] = apiLicense(vs[i])
	}
	return ms
}

func appArea(m *model.Area) game.Area {
	return game.Area{
		X:     int(*m.PosX),
		Y:     int(*m.PosY),
		SizeX: int(m.SizeX),
		SizeY: int(m.SizeY),
	}
}

func appCoord(m *model.Dig) game.Coord {
	return game.Coord{
		X:     int(*m.PosX),
		Y:     int(*m.PosY),
		Depth: uint8(*m.Depth),
	}
}

func apiTreasureList(v string) model.TreasureList {
	return model.TreasureList{model.Treasure(v)}
}
