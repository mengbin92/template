package chain

import (
	"context"
	"explorer/internal/conf"
	"sync"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
)

var (
	httpClient      *ethclient.Client
	wsClient	  *ethclient.Client
	initHttpOnce sync.Once
	initWSOnce sync.Once
)

func InitEthereumHttpClient(ctx context.Context, cfg *conf.ChainConfig, logger log.Logger) error {
	var err error

	initHttpOnce.Do(func() {
		httpClient, err = ethclient.Dial(cfg.HttpEndpoint)
		if err != nil {
			logger.Log(log.LevelError, "InitEthereumHTTPClient", "err", err)
		}
	})

	if err != nil {
		return errors.Wrap(err, "InitEthereumHTTPClient")
	}
	return err
}

func GetEthereumHttpClient() *ethclient.Client {
	return httpClient
}

func InitEthereumWSClient(ctx context.Context, cfg *conf.ChainConfig, logger log.Logger) error {
	var err error

	initWSOnce.Do(func() {
		httpClient, err = ethclient.Dial(cfg.WsEndpoint)
		if err != nil {
			logger.Log(log.LevelError, "InitEthereumWSClient", "err", err)
		}
	})

	if err != nil {
		return errors.Wrap(err, "InitEthereumWSClient")
	}
	return err
}

func GetEthereumWSClient() *ethclient.Client {
	return wsClient
}