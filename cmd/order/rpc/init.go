package rpc

import (
	"context"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/loadbalance"
	"github.com/cloudwego/kitex/pkg/retry"
	etcd "github.com/kitex-contrib/registry-etcd"
	trace "github.com/kitex-contrib/tracer-opentracing"
	"github.com/ozline/tiktok/config"
	"github.com/ozline/tiktok/kitex_gen/product"
	"github.com/ozline/tiktok/kitex_gen/product/productservice"
	"github.com/ozline/tiktok/pkg/constants"
	"github.com/ozline/tiktok/pkg/errno"
	"github.com/ozline/tiktok/pkg/middleware"
)

var productClient productservice.Client

func Init() {
	initProductRPC()
}

func initProductRPC() {
	r, err := etcd.NewEtcdResolver([]string{config.Etcd.Addr})
	if err != nil {
		panic(err)
	}

	c, err := productservice.NewClient(
		constants.ProductServiceName,
		client.WithMiddleware(middleware.CommonMiddleware),
		client.WithMuxConnection(constants.MuxConnection),
		client.WithRPCTimeout(constants.RPCTimeout),
		client.WithConnectTimeout(constants.ConnectTimeout),
		client.WithFailureRetry(retry.NewFailurePolicy()),
		client.WithResolver(r),
		client.WithSuite(trace.NewDefaultClientSuite()),
		client.WithLoadBalancer(loadbalance.NewWeightedRoundRobinBalancer()),
	)
	if err != nil {
		panic(err)
	}

	productClient = c
}

func GetSeckillActivity(ctx context.Context, activityId int64, token string) (*product.SeckillActivity, error) {
	resp, err := productClient.GetSeckill(ctx, &product.GetSeckillRequest{
		ActivityId: activityId,
		Token:      token,
	})
	if err != nil {
		return nil, err
	}
	if resp.Base.Code != errno.SuccessCode {
		return nil, errno.NewErrNo(resp.Base.Code, resp.Base.Msg)
	}
	return resp.Activity, nil
}

func DeductStock(ctx context.Context, activityId int64, token string) error {
	resp, err := productClient.DeductStock(ctx, &product.DeductStockRequest{
		ActivityId: activityId,
		Token:      token,
	})
	if err != nil {
		return err
	}
	if resp.Base.Code != errno.SuccessCode {
		return errno.NewErrNo(resp.Base.Code, resp.Base.Msg)
	}
	return nil
}
