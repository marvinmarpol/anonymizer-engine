package app

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/marvinmarpol/golang-boilerplate/internal/pkg/utils/cryptho"
	"github.com/marvinmarpol/golang-boilerplate/internal/protection/domain/mask"
	"github.com/marvinmarpol/golang-boilerplate/internal/protection/transport/pubsub"
	"github.com/marvinmarpol/golang-boilerplate/internal/protection/transport/web"

	watermillMiddleware "github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/go-chi/chi/middleware"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"

	"github.com/marvinmarpol/golang-boilerplate/internal/common/messaging"
	nrMiddleware "github.com/marvinmarpol/golang-boilerplate/internal/common/middlewares/monitoring"

	"github.com/newrelic/go-agent/v3/newrelic"
)

var (
	db       *pg.DB
	nrapp    *newrelic.Application
	dbModels = []interface{}{
		(*mask.Mask)(nil),
	}
)

func Run() {
	fx.New(fx.Provide(httpServer),
		fx.Invoke(func(*http.Server) {})).Run()
}

func httpServer(lc fx.Lifecycle) *http.Server {
	// load config
	config, err := loadConfig()
	if err != nil {
		logrus.WithField("err", err).Fatal("Failed to load config")
		return nil
	}

	// set log level
	logrus.SetLevel(logrus.Level(config.LogLevel))

	// init server
	srv := http.Server{
		Addr: config.ServiceAddress,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// load public and private key for KEK
			publicKey, err := cryptho.LoadRSAPublicKeyFromFile(config.PublicKeyPath)
			if err != nil {
				logrus.WithContext(ctx).WithField("err", err).Fatal("Failed to get public key")
				return err
			}
			privateKey, err := cryptho.LoadRSAPrivateKeyFromFile(config.PrivateKeyPath)
			if err != nil {
				logrus.WithContext(ctx).WithField("err", err).Fatal("Failed to get private key")
				return err
			}

			// connect to db
			db = pg.Connect(&pg.Options{
				Addr:         config.DBAddress,
				User:         config.DBUser,
				Password:     config.DBPassword,
				Database:     config.DBName,
				PoolSize:     config.DBPoolSize,
				MaxRetries:   config.DBMaxRetries,
				ReadTimeout:  time.Millisecond * time.Duration(config.DBReadTimeout),
				WriteTimeout: time.Millisecond * time.Duration(config.DBWriteTimeout),
				IdleTimeout:  time.Millisecond * time.Duration(config.DBIdleTimeout),
			})
			if err := db.Ping(ctx); err != nil {
				logrus.WithContext(ctx).WithField("err", err).Fatal("Failed to load database")
				return err
			}

			// create db models
			for _, v := range dbModels {
				err := db.Model(v).CreateTable(&orm.CreateTableOptions{
					IfNotExists: true,
				})
				if err != nil {
					logrus.WithContext(ctx).WithField("err", err).Fatal("Failed to create database model")
					return err
				}
			}

			// init newrelic app
			nrapp, err = newrelic.NewApplication(
				newrelic.ConfigLicense(config.NewrelicAPIKey),
				newrelic.ConfigAppName(config.NewrelicAPMName),
				newrelic.ConfigEnabled(true),
				newrelic.ConfigDistributedTracerEnabled(true),
			)
			if err != nil {
				logrus.WithContext(ctx).WithField("err", err).Fatal("Failed to initialize newrelic application")
				return err
			}

			// init http routes
			serverRoutes, err := InitializeWebServer(db, publicKey, privateKey)
			if err != nil {
				logrus.WithContext(ctx).WithField("err", err).Fatal("Failed to initialize web server")
				return err
			}

			// register the routes
			srv.Handler = web.RegisterRoutes(serverRoutes, []web.MiddlewareFunc{
				middleware.Logger,
				nrMiddleware.NewrelicAPM(nrapp),
				middleware.BasicAuth(config.Realm, map[string]string{
					config.BasicUsername: config.BasicPassword,
				}),
			})

			// listen to network
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				logrus.WithContext(ctx).WithField("err", err).Fatal("Failed to listen to http network")
				return err
			}

			// init pubsub handlers
			pubsubRoutes, err := InitializePubsubServer(db, publicKey, privateKey)
			if err != nil {
				logrus.WithContext(ctx).WithField("err", err).Fatal("Failed to initialize pubsub server")
				return err
			}

			// init pubsub router
			messaging, err := messaging.NewGooglePubSub(config.GoogleProjectID, config.GoogleCredential, messaging.DefaultLogger())
			if err != nil {
				logrus.WithContext(ctx).WithField("err", err).Fatal("Failed to initialize messaging module: ", config.GoogleCredential)
				return err
			}

			// register pubsub route or handler
			pubsub.RegisterRoutes(&messaging.Router, pubsubRoutes, messaging.Subscriber, []pubsub.MiddlewareFunc{
				watermillMiddleware.Recoverer,
				watermillMiddleware.Retry{
					MaxRetries:      config.PubsubMaxRetry,
					InitialInterval: time.Millisecond * time.Duration(config.PubsubRetryMS),
				}.Middleware,
			})

			// serve http
			go srv.Serve(ln)

			// serve pubsub
			/* go func() error {
				err = messaging.Router.Run(context.Background())
				if err != nil {
					logrus.WithContext(ctx).WithField("err", err).Fatal("Failed to listen to messaging network")
				}
				return err
			}() */

			return nil
		},

		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})

	return &srv
}
