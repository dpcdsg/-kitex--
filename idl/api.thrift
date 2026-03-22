namespace go api

// ==================== Models ====================

struct User {
    1: required i64 id,
    2: required string name,
    3: optional string avatar,
    4: optional string signature,
}

struct Product {
    1: required i64 id,
    2: required string name,
    3: required string description,
    4: required i64 price,
    5: required i64 stock,
    6: required string image_url,
    7: required string category,
    8: required i64 seller_id,
}

struct SeckillActivity {
    1: required i64 id,
    2: required i64 product_id,
    3: required string product_name,
    4: required i64 seckill_price,
    5: required i64 total_stock,
    6: required i64 available_stock,
    7: required string start_time,
    8: required string end_time,
    9: required i64 status,
}

struct Order {
    1: required i64 id,
    2: required string order_no,
    3: required i64 user_id,
    4: required i64 product_id,
    5: required string product_name,
    6: required i64 activity_id,
    7: required i64 amount,
    8: required i64 status,
    9: required string created_at,
}

// ==================== User ====================

struct UserRegisterRequest {
    1: required string username,
    2: required string password,
}

struct UserRegisterResponse {
    1: required i64 status_code = 0,
    2: optional string status_msg,
    3: required i64 user_id,
    4: required string token,
}

struct UserLoginRequest {
    1: required string username,
    2: required string password,
}

struct UserLoginResponse {
    1: required i64 status_code = 0,
    2: optional string status_msg,
    3: required i64 user_id,
    4: required string token,
}

struct UserRequest {
    1: required i64 user_id,
    2: required string token,
}

struct UserResponse {
    1: required i64 status_code = 0,
    2: optional string status_msg,
    3: required User user,
}

// ==================== Product ====================

struct ProductCreateRequest {
    1: required string token,
    2: required string name,
    3: required string description,
    4: required i64 price,
    5: required i64 stock,
    6: required string image_url,
    7: required string category,
}

struct ProductCreateResponse {
    1: required i64 status_code = 0,
    2: optional string status_msg,
    3: required Product product,
}

struct ProductUpdateRequest {
    1: required string token,
    2: required i64 product_id,
    3: required string name,
    4: required string description,
    5: required i64 price,
    6: required i64 stock,
    7: required string image_url,
    8: required string category,
}

struct ProductUpdateResponse {
    1: required i64 status_code = 0,
    2: optional string status_msg,
    3: required Product product,
}

struct ProductDetailRequest {
    1: required i64 product_id,
    2: optional string token,
}

struct ProductDetailResponse {
    1: required i64 status_code = 0,
    2: optional string status_msg,
    3: required Product product,
}

struct ProductListRequest {
    1: optional string token,
    2: optional i64 page,
    3: optional i64 size,
    4: optional string category,
}

struct ProductListResponse {
    1: required i64 status_code = 0,
    2: optional string status_msg,
    3: required list<Product> product_list,
    4: required i64 total,
}

// ==================== Seckill ====================

struct SeckillCreateRequest {
    1: required string token,
    2: required i64 product_id,
    3: required i64 seckill_price,
    4: required i64 total_stock,
    5: required string start_time,
    6: required string end_time,
}

struct SeckillCreateResponse {
    1: required i64 status_code = 0,
    2: optional string status_msg,
    3: required SeckillActivity activity,
}

struct SeckillListRequest {
    1: optional string token,
    2: optional i64 status,
    3: optional i64 page,
    4: optional i64 size,
}

struct SeckillListResponse {
    1: required i64 status_code = 0,
    2: optional string status_msg,
    3: required list<SeckillActivity> activity_list,
    4: required i64 total,
}

struct SeckillDetailRequest {
    1: required i64 activity_id,
    2: optional string token,
}

struct SeckillDetailResponse {
    1: required i64 status_code = 0,
    2: optional string status_msg,
    3: required SeckillActivity activity,
}

struct SeckillActionRequest {
    1: required string token,
    2: required i64 activity_id,
}

struct SeckillActionResponse {
    1: required i64 status_code = 0,
    2: optional string status_msg,
    3: required string order_no,
}

// ==================== Order ====================

struct OrderListRequest {
    1: required string token,
    2: optional i64 status,
    3: optional i64 page,
    4: optional i64 size,
}

struct OrderListResponse {
    1: required i64 status_code = 0,
    2: optional string status_msg,
    3: required list<Order> order_list,
    4: required i64 total,
}

struct OrderDetailRequest {
    1: required string token,
    2: required string order_no,
}

struct OrderDetailResponse {
    1: required i64 status_code = 0,
    2: optional string status_msg,
    3: required Order order,
}

struct OrderPayRequest {
    1: required string token,
    2: required string order_no,
}

struct OrderPayResponse {
    1: required i64 status_code = 0,
    2: optional string status_msg,
}

struct OrderCancelRequest {
    1: required string token,
    2: required string order_no,
}

struct OrderCancelResponse {
    1: required i64 status_code = 0,
    2: optional string status_msg,
}

// ==================== Services ====================

service UserService {
    UserRegisterResponse UserRegister(1: UserRegisterRequest req) (api.post="/seckill/user/register/")
    UserLoginResponse UserLogin(1: UserLoginRequest req) (api.post="/seckill/user/login/")
    UserResponse UserInfo(1: UserRequest req) (api.get="/seckill/user/")
}

service ProductService {
    ProductCreateResponse ProductCreate(1: ProductCreateRequest req) (api.post="/seckill/product/create/")
    ProductUpdateResponse ProductUpdate(1: ProductUpdateRequest req) (api.post="/seckill/product/update/")
    ProductDetailResponse ProductDetail(1: ProductDetailRequest req) (api.get="/seckill/product/detail/")
    ProductListResponse ProductList(1: ProductListRequest req) (api.get="/seckill/product/list/")

    SeckillCreateResponse SeckillCreate(1: SeckillCreateRequest req) (api.post="/seckill/activity/create/")
    SeckillListResponse SeckillList(1: SeckillListRequest req) (api.get="/seckill/activity/list/")
    SeckillDetailResponse SeckillDetail(1: SeckillDetailRequest req) (api.get="/seckill/activity/detail/")
}

service OrderService {
    SeckillActionResponse SeckillAction(1: SeckillActionRequest req) (api.post="/seckill/action/")

    OrderListResponse OrderList(1: OrderListRequest req) (api.get="/seckill/order/list/")
    OrderDetailResponse OrderDetail(1: OrderDetailRequest req) (api.get="/seckill/order/detail/")
    OrderPayResponse OrderPay(1: OrderPayRequest req) (api.post="/seckill/order/pay/")
    OrderCancelResponse OrderCancel(1: OrderCancelRequest req) (api.post="/seckill/order/cancel/")
}
