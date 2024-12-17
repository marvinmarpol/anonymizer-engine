package newrelic

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/newrelic/go-agent/v3/newrelic"
)

const (
	routeParamAttributeKeyPrefix = "routeparam."
)

func attachTxnAttributes(rctx *chi.Context, txn *newrelic.Transaction) {
	for i := range rctx.URLParams.Keys {
		key, value := rctx.URLParams.Keys[i], rctx.URLParams.Values[i]

		// The key must contain fewer than or equal to 255 bytes.
		prefixedKey := routeParamAttributeKeyPrefix + key
		if len(prefixedKey) > 255 {
			continue
		}

		txn.AddAttribute(prefixedKey, value)
	}
}

func NewrelicAPM(nrapp *newrelic.Application) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rctx := chi.RouteContext(r.Context())

			routePath := r.URL.RawPath
			if r.URL.RawPath == "" {
				routePath = r.URL.Path
			}

			router := rctx.Routes
			newRctx := chi.NewRouteContext()
			if !router.Match(newRctx, r.Method, routePath) {
				next.ServeHTTP(w, r)
				return
			}

			path := newRctx.RoutePattern()
			name := fmt.Sprintf("%s %s", r.Method, path)

			txn := newrelic.FromContext(r.Context())
			if txn == nil {
				txn = nrapp.StartTransaction(name)
			}
			defer txn.End()

			txn.SetWebRequestHTTP(r)
			w = txn.SetWebResponse(w)

			// Carry the transaction throughout the context.
			r = newrelic.RequestWithTransactionContext(r, txn)

			attachTxnAttributes(rctx, txn)

			next.ServeHTTP(w, r)
		})
	}
}
