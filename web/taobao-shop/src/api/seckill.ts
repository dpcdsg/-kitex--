import { apiGet, apiPostJson } from './client';
import type { Order, Product, SeckillActivity } from './types';

export async function userRegister(username: string, password: string) {
  return apiPostJson<{
    user_id: number;
    token: string;
    status_code?: number;
    status_msg?: string;
  }>('/user/register/', { username, password });
}

export async function userLogin(username: string, password: string) {
  return apiPostJson<{
    user_id: number;
    token: string;
    status_code?: number;
    status_msg?: string;
  }>('/user/login/', { username, password });
}

export async function userInfo(userId: number, token: string) {
  return apiGet<{
    user?: { id: number; name: string; avatar?: string; signature?: string };
    status_code?: number;
    status_msg?: string;
  }>(`/user/?user_id=${userId}`, token);
}

export async function productList(page: number, size: number, token?: string | null, category?: string) {
  const q = new URLSearchParams({ page: String(page), size: String(size) });
  if (category) q.set('category', category);
  return apiGet<{
    product_list: Product[];
    total: number;
  }>(`/product/list/?${q}`, token ?? null);
}

export async function productDetail(productId: number, token?: string | null) {
  const q = new URLSearchParams({ product_id: String(productId) });
  return apiGet<{ product: Product }>(`/product/detail/?${q}`, token ?? null);
}

export async function productCreate(
  token: string,
  body: {
    name: string;
    description: string;
    price: number;
    stock: number;
    image_url: string;
    category: string;
  },
) {
  return apiPostJson<{ product: Product }>(
    '/product/create/',
    { token, ...body },
    token,
  );
}

export async function seckillList(page: number, size: number, token?: string | null, status = -1) {
  const q = new URLSearchParams({
    page: String(page),
    size: String(size),
    status: String(status),
  });
  return apiGet<{ activity_list: SeckillActivity[]; total: number }>(
    `/activity/list/?${q}`,
    token ?? null,
  );
}

export async function seckillDetail(activityId: number, token?: string | null) {
  const q = new URLSearchParams({ activity_id: String(activityId) });
  return apiGet<{ activity: SeckillActivity }>(`/activity/detail/?${q}`, token ?? null);
}

export async function seckillCreate(
  token: string,
  body: {
    product_id: number;
    seckill_price: number;
    total_stock: number;
    start_time: string;
    end_time: string;
  },
) {
  return apiPostJson<{ activity: SeckillActivity }>(
    '/activity/create/',
    { token, ...body },
    token,
  );
}

export async function seckillAction(token: string, activityId: number) {
  return apiPostJson<{ order_no: string }>(
    '/action/',
    { token, activity_id: activityId },
    token,
  );
}

export async function orderList(token: string, page: number, size: number, status = -1) {
  const q = new URLSearchParams({
    page: String(page),
    size: String(size),
    status: String(status),
  });
  return apiGet<{ order_list: Order[]; total: number }>(`/order/list/?${q}`, token);
}

export async function orderDetail(token: string, orderNo: string) {
  const q = new URLSearchParams({ order_no: orderNo });
  return apiGet<{ order: Order }>(`/order/detail/?${q}`, token);
}

export async function orderPay(token: string, orderNo: string) {
  return apiPostJson<Record<string, unknown>>('/order/pay/', { token, order_no: orderNo }, token);
}

export async function orderCancel(token: string, orderNo: string) {
  return apiPostJson<Record<string, unknown>>('/order/cancel/', { token, order_no: orderNo }, token);
}
