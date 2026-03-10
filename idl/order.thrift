namespace go order

struct BaseResp {
    1: i64 code,
    2: string msg,
}

struct Order {
    1: i64 id,
    2: string order_no,
    3: i64 user_id,
    4: i64 product_id,
    5: string product_name,
    6: i64 activity_id,
    7: i64 amount,
    8: i64 status,
    9: string created_at,
    10: string updated_at,
}

struct SeckillRequest {
    1: string token,
    2: i64 activity_id,
}

struct SeckillResponse {
    1: BaseResp base,
    2: string order_no,
}

struct GetOrderRequest {
    1: string token,
    2: string order_no,
}

struct GetOrderResponse {
    1: BaseResp base,
    2: Order order,
}

struct ListOrdersRequest {
    1: string token,
    2: i64 status,
    3: i64 page,
    4: i64 size,
}

struct ListOrdersResponse {
    1: BaseResp base,
    2: list<Order> order_list,
    3: i64 total,
}

struct PayOrderRequest {
    1: string token,
    2: string order_no,
}

struct PayOrderResponse {
    1: BaseResp base,
}

struct CancelOrderRequest {
    1: string token,
    2: string order_no,
}

struct CancelOrderResponse {
    1: BaseResp base,
}

service OrderService {
    SeckillResponse Seckill(1: SeckillRequest req),
    GetOrderResponse GetOrder(1: GetOrderRequest req),
    ListOrdersResponse ListOrders(1: ListOrdersRequest req),
    PayOrderResponse PayOrder(1: PayOrderRequest req),
    CancelOrderResponse CancelOrder(1: CancelOrderRequest req),
}
