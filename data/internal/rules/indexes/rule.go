package indexes

import (
	"github.com/WLM1ke/gomoex"
	"github.com/WLM1ke/poptimizer/data/internal/domain"
	"github.com/WLM1ke/poptimizer/data/internal/repo"
	"github.com/WLM1ke/poptimizer/data/internal/rules/template"
	"github.com/WLM1ke/poptimizer/data/pkg/lgr"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func New(logger *lgr.Logger, db *mongo.Database, iss *gomoex.ISSClient, timeout time.Duration) domain.Rule {
	return template.NewRule[gomoex.Quote](
		"IndexesRule",
		logger,
		repo.NewMongo[gomoex.Quote](db),
		selector{},
		gateway{iss: iss},
		validator,
		true,
		template.EventCtxFuncWithTimeout(timeout),
	)
}
