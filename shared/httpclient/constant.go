package httpclient

import (
	"net"
	"net/http"
	"time"
)

const (
	clientTimeout = 30
)

var (
	netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: clientTimeout * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
	}

	netClient = &http.Client{
		Timeout:   clientTimeout * time.Second,
		Transport: netTransport,
	}
)
