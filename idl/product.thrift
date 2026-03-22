namespace go product

struct BaseResp {
    1: i64 code,
    2: string msg,
}

struct Product {
    1: i64 id,
    2: string name,
    3: string description,
    4: i64 price,
    5: i64 stock,
    6: string image_url,
    7: string category,
    8: i64 status,
}

struct SeckillActivity {
    1: i64 id,
    2: i64 product_id,
    3: string product_name,
    4: i64 seckill_price,
    5: i64 total_stock,
    6: i64 available_stock,
    7: string start_time,
    8: string end_time,
    9: i64 status,
}

struct CreateProductRequest {
    1: string token,
    2: string name,
    3: string description,
    4: i64 price,
    5: i64 stock,
    6: string image_url,
    7: string category,
}

struct CreateProductResponse {
    1: BaseResp base,
    2: Product product,
    3: i64 seller_id,
}

struct UpdateProductRequest {
    1: string token,
    2: i64 product_id,
    3: string name,
    4: string description,
    5: i64 price,
    6: i64 stock,
    7: string image_url,
    8: string category,
}

struct UpdateProductResponse {
    1: BaseResp base,
    2: Product product,
    3: i64 seller_id,
}

struct GetProductRequest {
    1: i64 product_id,
    2: string token,
}

struct GetProductResponse {
    1: BaseResp base,
    2: Product product,
    3: i64 seller_id,
}

struct ListProductsRequest {
    1: string token,
    2: i64 page,
    3: i64 size,
    4: string category,
}

struct ListProductsResponse {
    1: BaseResp base,
    2: list<Product> product_list,
    3: i64 total,
}

struct CreateSeckillRequest {
    1: string token,
    2: i64 product_id,
    3: i64 seckill_price,
    4: i64 total_stock,
    5: string start_time,
    6: string end_time,
}

struct CreateSeckillResponse {
    1: BaseResp base,
    2: SeckillActivity activity,
}

struct GetSeckillRequest {
    1: i64 activity_id,
    2: string token,
}

struct GetSeckillResponse {
    1: BaseResp base,
    2: SeckillActivity activity,
}

struct ListSeckillRequest {
    1: string token,
    2: i64 status,
    3: i64 page,
    4: i64 size,
}

struct ListSeckillResponse {
    1: BaseResp base,
    2: list<SeckillActivity> activity_list,
    3: i64 total,
}

struct DeductStockRequest {
    1: i64 activity_id,
    2: string token,
}

struct DeductStockResponse {
    1: BaseResp base,
}

service ProductService {
    CreateProductResponse CreateProduct(1: CreateProductRequest req),
    UpdateProductResponse UpdateProduct(1: UpdateProductRequest req),
    GetProductResponse GetProduct(1: GetProductRequest req),
    ListProductsResponse ListProducts(1: ListProductsRequest req),
    CreateSeckillResponse CreateSeckill(1: CreateSeckillRequest req),
    GetSeckillResponse GetSeckill(1: GetSeckillRequest req),
    ListSeckillResponse ListSeckill(1: ListSeckillRequest req),
    DeductStockResponse DeductStock(1: DeductStockRequest req),
}
