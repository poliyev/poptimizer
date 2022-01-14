package dates

import (
	"context"
	"github.com/WLM1ke/gomoex"
	"github.com/WLM1ke/poptimizer/data/internal/domain"
	"time"
)

type gateway struct {
	iss *gomoex.ISSClient
}

func (g gateway) Get(ctx context.Context, table domain.Table[gomoex.Date], _ time.Time) ([]gomoex.Date, error) {
	rows, err := g.iss.MarketDates(ctx, gomoex.EngineStock, gomoex.MarketShares)

	switch {
	case err != nil:
		return nil, err
	case table.IsEmpty():
		return rows, nil
	case table.LastRow().Till.Before(rows[0].Till):
		return rows, nil
	default:
		return nil, nil
	}
}
