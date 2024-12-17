// wire.go
//go:build wireinject
// +build wireinject

package app

import (
	"crypto/rsa"

	"github.com/marvinmarpol/golang-boilerplate/internal/protection/command"
	"github.com/marvinmarpol/golang-boilerplate/internal/protection/domain/mask"
	"github.com/marvinmarpol/golang-boilerplate/internal/protection/query"
	"github.com/marvinmarpol/golang-boilerplate/internal/protection/service"
	"github.com/marvinmarpol/golang-boilerplate/internal/protection/transport/pubsub"
	"github.com/marvinmarpol/golang-boilerplate/internal/protection/transport/web"

	"github.com/go-pg/pg/v10"
	"github.com/google/wire"
)

var moduleSet = wire.NewSet(
	// command handlers
	command.NewCreateMaskHandler,
	wire.Bind(new(command.CommandHandler[command.CreateMaskCommand]), new(*command.CreateMaskHandler)),

	command.NewUpdateTokenHandler,
	wire.Bind(new(command.CommandHandler[command.UpdateTokenCommand]), new(*command.UpdateTokenHandler)),

	command.NewUpdateMaskHandler,
	wire.Bind(new(command.CommandHandler[command.UpdateMaskCommand]), new(*command.UpdateMaskHandler)),

	// query handlers
	query.NewGetCypherHandler,
	wire.Bind(new(query.QueryHandler[query.GetCypherQuery]), new(*query.GetCypherHandler)),

	query.NewGetMaskHandler,
	wire.Bind(new(query.QueryHandler[query.GetMaskQuery]), new(*query.GetMaskHandler)),

	query.NewGetTokenHandler,
	wire.Bind(new(query.QueryHandler[query.GetTokenQuery]), new(*query.GetTokenHandler)),

	query.NewGetRotateCandidateHandler,
	wire.Bind(new(query.QueryHandler[query.GetRotateCandidateQuery]), new(*query.GetRotateCandidateHandler)),

	// repositories
	mask.NewPostgresRepository,
	wire.Bind(new(mask.Repository), new(*mask.PostgresRepository)),

	// services
	service.NewServiceImpl,
	wire.Bind(new(service.Services), new(*service.ServiceImpl)),

	// command and query list
	wire.Struct(new(command.Commands), "*"),
	wire.Struct(new(query.Queries), "*"),
)

func InitializeWebServer(db *pg.DB, pubKey *rsa.PublicKey, priKey *rsa.PrivateKey) (web.Route, error) {
	panic(wire.Build(
		moduleSet,
		web.NewServer,
		wire.Bind(new(web.Route), new(*web.Server)),
	))
}

func InitializePubsubServer(db *pg.DB, pubKey *rsa.PublicKey, priKey *rsa.PrivateKey) (pubsub.Route, error) {
	panic(wire.Build(
		moduleSet,
		pubsub.NewServer,
		wire.Bind(new(pubsub.Route), new(*pubsub.Server)),
	))
}
