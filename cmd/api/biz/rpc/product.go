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

func InitProductRPC() {
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

func ProductCreate(ctx context.Context, req *product.CreateProductRequest) (*product.Product, error) {
	resp, err := productClient.CreateProduct(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.Base.Code != errno.SuccessCode {
		return nil, errno.NewErrNo(resp.Base.Code, resp.Base.Msg)
	}
	return resp.Product, nil
}

func ProductDetail(ctx context.Context, productId int64, token string) (*product.Product, error) {
	resp, err := productClient.GetProduct(ctx, &product.GetProductRequest{
		ProductId: productId,
		Token:     token,
	})
	if err != nil {
		return nil, err
	}
	if resp.Base.Code != errno.SuccessCode {
		return nil, errno.NewErrNo(resp.Base.Code, resp.Base.Msg)
	}
	return resp.Product, nil
}

func ProductList(ctx context.Context, req *product.ListProductsRequest) ([]*product.Product, int64, error) {
	resp, err := productClient.ListProducts(ctx, req)
	if err != nil {
		return nil, 0, err
	}
	if resp.Base.Code != errno.SuccessCode {
		return nil, 0, errno.NewErrNo(resp.Base.Code, resp.Base.Msg)
	}
	return resp.ProductList, resp.Total, nil
}

func SeckillCreate(ctx context.Context, req *product.CreateSeckillRequest) (*product.SeckillActivity, error) {
	resp, err := productClient.CreateSeckill(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.Base.Code != errno.SuccessCode {
		return nil, errno.NewErrNo(resp.Base.Code, resp.Base.Msg)
	}
	return resp.Activity, nil
}

func SeckillList(ctx context.Context, req *product.ListSeckillRequest) ([]*product.SeckillActivity, int64, error) {
	resp, err := productClient.ListSeckill(ctx, req)
	if err != nil {
		return nil, 0, err
	}
	if resp.Base.Code != errno.SuccessCode {
		return nil, 0, errno.NewErrNo(resp.Base.Code, resp.Base.Msg)
	}
	return resp.ActivityList, resp.Total, nil
}

func SeckillDetail(ctx context.Context, activityId int64, token string) (*product.SeckillActivity, error) {
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
