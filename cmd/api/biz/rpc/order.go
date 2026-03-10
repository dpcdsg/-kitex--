package rpc

import (
	"context"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/loadbalance"
	"github.com/cloudwego/kitex/pkg/retry"
	etcd "github.com/kitex-contrib/registry-etcd"
	trace "github.com/kitex-contrib/tracer-opentracing"
	"github.com/ozline/tiktok/config"
	"github.com/ozline/tiktok/kitex_gen/order"
	"github.com/ozline/tiktok/kitex_gen/order/orderservice"
	"github.com/ozline/tiktok/pkg/constants"
	"github.com/ozline/tiktok/pkg/errno"
	"github.com/ozline/tiktok/pkg/middleware"
)

func InitOrderRPC() {
	r, err := etcd.NewEtcdResolver([]string{config.Etcd.Addr})
	if err != nil {
		panic(err)
	}

	c, err := orderservice.NewClient(
		constants.OrderServiceName,
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

	orderClient = c
}

func SeckillAction(ctx context.Context, req *order.SeckillRequest) (string, error) {
	resp, err := orderClient.Seckill(ctx, req)
	if err != nil {
		return "", err
	}
	if resp.Base.Code != errno.SuccessCode {
		return "", errno.NewErrNo(resp.Base.Code, resp.Base.Msg)
	}
	return resp.OrderNo, nil
}

func OrderList(ctx context.Context, req *order.ListOrdersRequest) ([]*order.Order, int64, error) {
	resp, err := orderClient.ListOrders(ctx, req)
	if err != nil {
		return nil, 0, err
	}
	if resp.Base.Code != errno.SuccessCode {
		return nil, 0, errno.NewErrNo(resp.Base.Code, resp.Base.Msg)
	}
	return resp.OrderList, resp.Total, nil
}

func OrderDetail(ctx context.Context, req *order.GetOrderRequest) (*order.Order, error) {
	resp, err := orderClient.GetOrder(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.Base.Code != errno.SuccessCode {
		return nil, errno.NewErrNo(resp.Base.Code, resp.Base.Msg)
	}
	return resp.Order, nil
}

func OrderPay(ctx context.Context, req *order.PayOrderRequest) error {
	resp, err := orderClient.PayOrder(ctx, req)
	if err != nil {
		return err
	}
	if resp.Base.Code != errno.SuccessCode {
		return errno.NewErrNo(resp.Base.Code, resp.Base.Msg)
	}
	return nil
}

func OrderCancel(ctx context.Context, req *order.CancelOrderRequest) error {
	resp, err := orderClient.CancelOrder(ctx, req)
	if err != nil {
		return err
	}
	if resp.Base.Code != errno.SuccessCode {
		return errno.NewErrNo(resp.Base.Code, resp.Base.Msg)
	}
	return nil
}
