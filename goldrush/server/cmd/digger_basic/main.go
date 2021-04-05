// Command digger_basic implements the most trivial and dumb client which
// can be implemented in an hour.
package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/go-openapi/swag"

	"github.com/Djarvur/allcups-itrally-2020-task/api/openapi/client"
	"github.com/Djarvur/allcups-itrally-2020-task/api/openapi/client/op"
	"github.com/Djarvur/allcups-itrally-2020-task/api/openapi/model"
)

const codeWrongCoord = 1000

func errCode(err error) int {
	var errDefault interface{ GetPayload() *model.Error }
	switch ok := errors.As(err, &errDefault); true {
	case ok:
		return int(*errDefault.GetPayload().Code)
	case err == nil:
		return 0
	default:
		return -1
	}
}

//nolint:gocognit // TODO.
func main() {
	cfg := client.DefaultTransportConfig()
	if len(os.Args) == 1+1 {
		cfg.Host = os.Args[1] // host:port
	}
	c := client.NewHTTPClientWithConfig(nil, cfg)

	balance := 0
	license := &model.License{}
	for x := int64(0); true; x++ {
	ROW:
		for y := int64(0); true; y++ {
			explore, err := c.Op.ExploreArea(op.NewExploreAreaParams().WithArgs(&model.Area{
				PosX:  swag.Int64(x),
				PosY:  swag.Int64(y),
				SizeX: 1,
				SizeY: 1,
			}))
			switch code := errCode(err); true {
			case code == codeWrongCoord:
				break ROW
			case code == -1:
				fmt.Println(err)
				os.Exit(0)
			case code != 0:
				continue ROW
			}
			for depth, left := 1, explore.Payload.Amount; depth <= 10 && left > 0; depth++ {
				for license.DigAllowed <= license.DigUsed {
					res, err := c.Op.IssueLicense(op.NewIssueLicenseParams().WithArgs(model.Wallet{}))
					if err == nil {
						license = res.Payload
					}
				}
				license.DigUsed++
				dig, err := c.Op.Dig(op.NewDigParams().WithArgs(&model.Dig{
					LicenseID: license.ID,
					PosX:      swag.Int64(x),
					PosY:      swag.Int64(y),
					Depth:     swag.Int64(int64(depth)),
				}))
				if err == nil {
					for _, treasure := range dig.Payload {
						res, err := c.Op.Cash(op.NewCashParams().WithArgs(treasure))
						if err == nil {
							left--
							balance += len(res.Payload)
							fmt.Println("balance:", balance)
						}
					}
				}
			}
		}
	}
}
