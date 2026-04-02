import { apiGet, apiPostJson } from './client';
import type { Order, Product, SeckillActivity } from './types';

/** dpc / 密码 123 的联调演示账号，仅此用户使用前端种子/默认数据 */
export const MOCK_DPC_USER_ID = 10001;
const MOCK_DPC_TOKEN = 'mock-token-dpc-123';
const MOCK_PRODUCT_IMAGE_URL = 'https://ddragon.leagueoflegends.com/cdn/img/champion/splash/Katarina_0.jpg';

const LS_OWNED_PRODUCTS = `taobao_demo_owned_products_${MOCK_DPC_USER_ID}`;
const LS_SECKILL_ACTIVITIES = `taobao_demo_seckill_activities_${MOCK_DPC_USER_ID}`;
const LS_ORDERS = `taobao_demo_orders_${MOCK_DPC_USER_ID}`;

function isMockDpcToken(token: string | null | undefined): boolean {
  return token === MOCK_DPC_TOKEN;
}

function readJson<T>(key: string, fallback: T): T {
  try {
    const raw = localStorage.getItem(key);
    if (!raw) return fallback;
    return JSON.parse(raw) as T;
  } catch {
    return fallback;
  }
}

function writeJson<T>(key: string, value: T) {
  localStorage.setItem(key, JSON.stringify(value));
}

function ensureSeed() {
  // products（SellerPage 的本地 owned 列表会直接使用同一个 key）
  const seededProducts: Product[] = [
    {
      id: '101',
      name: '模拟商品 A',
      description: '用于验证前端界面的演示商品。',
      price: 19900,
      stock: 20,
      image_url: MOCK_PRODUCT_IMAGE_URL,
      category: '数码',
      seller_id: String(MOCK_DPC_USER_ID),
    },
    {
      id: '102',
      name: '模拟商品 B',
      description: '当后端不可用或无数据时，这些数据会显示。',
      price: 9900,
      stock: 80,
      image_url: MOCK_PRODUCT_IMAGE_URL,
      category: '家居',
      seller_id: String(MOCK_DPC_USER_ID),
    },
  ];

  const seededActivities: SeckillActivity[] = [
    {
      id: '1000001',
      product_id: '101',
      product_name: '模拟商品 A',
      seckill_price: 5900,
      total_stock: 100,
      available_stock: 100,
      start_time: '2026-03-10 10:00:00',
      end_time: '2026-03-10 11:00:00',
      status: 1,
    },
    {
      id: '1000002',
      product_id: '102',
      product_name: '模拟商品 B',
      seckill_price: 4900,
      total_stock: 200,
      available_stock: 200,
      start_time: '2026-03-10 10:00:00',
      end_time: '2026-03-10 11:00:00',
      status: 1,
    },
  ];

  const seededOrders: Order[] = [
    {
      id: 1,
      order_no: 'SK10001',
      user_id: MOCK_DPC_USER_ID,
      product_id: 101,
      product_name: '模拟商品 A',
      activity_id: 1000001,
      amount: 19900,
      status: 0,
      created_at: '2026-01-01 12:00:00',
    },
    {
      id: 2,
      order_no: 'SK10002',
      user_id: MOCK_DPC_USER_ID,
      product_id: 102,
      product_name: '模拟商品 B',
      activity_id: 1000002,
      amount: 9900,
      status: 1,
      created_at: '2026-01-01 13:00:00',
    },
  ];

  const existingProducts = readJson<Product[]>(LS_OWNED_PRODUCTS, []);
  if (existingProducts.length === 0) {
    writeJson(LS_OWNED_PRODUCTS, seededProducts);
  } else {
    // Backfill older mock records that may not have image_url.
    const normalized = existingProducts.map((p) => ({
      ...p,
      image_url: p.image_url || MOCK_PRODUCT_IMAGE_URL,
    }));
    writeJson(LS_OWNED_PRODUCTS, normalized);
  }
  if (!localStorage.getItem(LS_SECKILL_ACTIVITIES)) {
    writeJson(LS_SECKILL_ACTIVITIES, seededActivities);
  }
  if (!localStorage.getItem(LS_ORDERS)) {
    writeJson(LS_ORDERS, seededOrders);
  }
}

function nowStr() {
  const d = new Date();
  const iso = d.toISOString(); // 2026-03-10T12:34:56.789Z
  return iso.slice(0, 19).replace('T', ' ');
}

export async function userRegister(username: string, password: string) {
  // dpc/123 mock：直接返回一个可用登录态，避免依赖后端
  if (username === 'dpc' && password === '123') {
    return { user_id: MOCK_DPC_USER_ID, token: MOCK_DPC_TOKEN };
  }
  return apiPostJson<{
    user_id: number;
    token: string;
    status_code?: number;
    status_msg?: string;
  }>('/user/register/', { username, password });
}

export async function userLogin(username: string, password: string) {
  // dpc/123 mock：兼容登录页“真实调用接口不影响”的要求
  if (username === 'dpc' && password === '123') {
    return { user_id: MOCK_DPC_USER_ID, token: MOCK_DPC_TOKEN };
  }
  return apiPostJson<{
    user_id: number;
    token: string;
    status_code?: number;
    status_msg?: string;
  }>('/user/login/', { username, password });
}

export async function userInfo(userId: number, token: string) {
  if (isMockDpcToken(token) && userId === MOCK_DPC_USER_ID) {
    return {
      user: { id: MOCK_DPC_USER_ID, name: 'dpc', avatar: undefined, signature: undefined },
    };
  }
  return apiGet<{
    user?: { id: number; name: string; avatar?: string; signature?: string };
    status_code?: number;
    status_msg?: string;
  }>(`/user/?user_id=${userId}`, token);
}

export async function productList(page: number, size: number, token?: string | null, category?: string) {
  if (isMockDpcToken(token ?? null)) {
    ensureSeed();
    const all = readJson<Product[]>(LS_OWNED_PRODUCTS, []).slice();
    const filtered = category ? all.filter((p) => p.category === category) : all;
    const total = filtered.length;
    const offset = (page - 1) * size;
    const list = filtered.slice(offset, offset + size);
    return { product_list: list, total };
  }
  const q = new URLSearchParams({ page: String(page), size: String(size) });
  if (category) q.set('category', category);
  return apiGet<{
    product_list: Product[];
    total: number;
  }>(`/product/list/?${q}`, token ?? null);
}

export async function productDetail(productId: string, token?: string | null) {
  if (isMockDpcToken(token ?? null)) {
    ensureSeed();
    const all = readJson<Product[]>(LS_OWNED_PRODUCTS, []);
    const p = all.find((x) => x.id === productId);
    if (!p) throw new Error('商品不存在');
    return { product: p };
  }
  const q = new URLSearchParams({ product_id: productId });
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
  if (isMockDpcToken(token)) {
    ensureSeed();
    const current = readJson<Product[]>(LS_OWNED_PRODUCTS, []);
    const nextId = current.reduce((m, p) => Math.max(m, Number(p.id) || 0), 0) + 1;
    const product: Product = {
      id: String(nextId),
      name: body.name,
      description: body.description,
      price: body.price,
      stock: body.stock,
      image_url: body.image_url,
      category: body.category,
      seller_id: String(MOCK_DPC_USER_ID),
    };
    const merged = [product, ...current];
    writeJson(LS_OWNED_PRODUCTS, merged);
    return { product };
  }
  return apiPostJson<{ product: Product }>(
    '/product/create/',
    { token, ...body },
    token,
  );
}

export async function productUpdate(
  token: string,
  body: {
    product_id: string;
    name: string;
    description: string;
    price: number;
    stock: number;
    image_url: string;
    category: string;
  },
) {
  if (isMockDpcToken(token)) {
    ensureSeed();
    const current = readJson<Product[]>(LS_OWNED_PRODUCTS, []);
    const idx = current.findIndex((p) => p.id === body.product_id);
    if (idx < 0) throw new Error('商品不存在');

    const updated: Product = {
      ...current[idx],
      name: body.name,
      description: body.description,
      price: body.price,
      stock: body.stock,
      image_url: body.image_url,
      category: body.category,
      seller_id: String(MOCK_DPC_USER_ID),
    };

    const next = current.slice();
    next[idx] = updated;
    writeJson(LS_OWNED_PRODUCTS, next);
    return { product: updated };
  }
  return apiPostJson<{ product: Product }>(
    '/product/update/',
    { token, ...body },
    token,
  );
}

export async function productPublish(token: string, productId: string) {
  if (isMockDpcToken(token)) {
    // mock 模式不区分发布/删除，仅用于接口连通性
    return {};
  }
  return apiPostJson<Record<string, unknown>>('/product/publish/', { token, product_id: productId }, token);
}

export async function productDelete(token: string, productId: string) {
  if (isMockDpcToken(token)) {
    // mock 下：从本地 owned 列表移除（模拟软删除效果）
    ensureSeed();
    const products = readJson<Product[]>(LS_OWNED_PRODUCTS, []);
    const idx = products.findIndex((p) => p.id === productId);
    if (idx >= 0) {
      const next = products.slice();
      next.splice(idx, 1);
      writeJson(LS_OWNED_PRODUCTS, next);
    }
    return {};
  }
  return apiPostJson<Record<string, unknown>>('/product/delete/', { token, product_id: productId }, token);
}

export async function seckillDelete(token: string, activityId: string) {
  if (isMockDpcToken(token)) {
    ensureSeed();
    const activities = readJson<SeckillActivity[]>(LS_SECKILL_ACTIVITIES, []);
    const next = activities.filter((a) => a.id !== activityId);
    writeJson(LS_SECKILL_ACTIVITIES, next);
    return {};
  }

  return apiPostJson<Record<string, unknown>>(
    '/activity/delete/',
    { token, activity_id: activityId },
    token,
  );
}

export async function seckillList(page: number, size: number, token?: string | null, status = -1) {
  if (isMockDpcToken(token ?? null)) {
    ensureSeed();
    const all = readJson<SeckillActivity[]>(LS_SECKILL_ACTIVITIES, []);
    const filtered = status >= 0 ? all.filter((a) => a.status === status) : all;
    const total = filtered.length;
    const offset = (page - 1) * size;
    const list = filtered.slice(offset, offset + size);
    return { activity_list: list, total };
  }
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

export async function seckillDetail(activityId: string, token?: string | null) {
  if (isMockDpcToken(token ?? null)) {
    ensureSeed();
    const all = readJson<SeckillActivity[]>(LS_SECKILL_ACTIVITIES, []);
    const a = all.find((x) => x.id === activityId);
    if (!a) throw new Error('秒杀活动不存在');
    return { activity: a };
  }
  const q = new URLSearchParams({ activity_id: activityId });
  return apiGet<{ activity: SeckillActivity }>(`/activity/detail/?${q}`, token ?? null);
}

export async function seckillCreate(
  token: string,
  body: {
    product_id: string;
    seckill_price: number;
    total_stock: number;
    start_time: string;
    end_time: string;
  },
) {
  if (isMockDpcToken(token)) {
    ensureSeed();
    const products = readJson<Product[]>(LS_OWNED_PRODUCTS, []);
    const product = products.find((p) => p.id === body.product_id);
    const activities = readJson<SeckillActivity[]>(LS_SECKILL_ACTIVITIES, []);
    const nextId = activities.reduce((m, a) => Math.max(m, Number(a.id) || 0), 0) + 1;
    const activity: SeckillActivity = {
      id: String(nextId),
      product_id: body.product_id,
      product_name: product?.name ?? `商品 #${body.product_id}`,
      seckill_price: body.seckill_price,
      total_stock: body.total_stock,
      available_stock: body.total_stock,
      start_time: body.start_time,
      end_time: body.end_time,
      status: 1,
    };
    writeJson(LS_SECKILL_ACTIVITIES, [activity, ...activities]);
    return { activity };
  }
  return apiPostJson<{ activity: SeckillActivity }>(
    '/activity/create/',
    { token, ...body },
    token,
  );
}

export async function seckillUpdate(
  token: string,
  body: {
    activity_id: string;
    product_id: string;
    seckill_price: number;
    total_stock: number;
    start_time: string;
    end_time: string;
  },
) {
  if (isMockDpcToken(token)) {
    ensureSeed();
    const products = readJson<Product[]>(LS_OWNED_PRODUCTS, []);
    const product = products.find((p) => p.id === body.product_id);
    const activities = readJson<SeckillActivity[]>(LS_SECKILL_ACTIVITIES, []);
    const idx = activities.findIndex((a) => a.id === body.activity_id);
    if (idx < 0) throw new Error('秒杀活动不存在');

    const old = activities[idx];
    const sold = Math.max(0, old.total_stock - old.available_stock);
    if (body.total_stock < sold) {
      throw new Error('活动总库存不能小于已售数量');
    }
    const available = body.total_stock - sold;

    const updated: SeckillActivity = {
      ...old,
      product_id: body.product_id,
      product_name: product?.name ?? old.product_name,
      seckill_price: body.seckill_price,
      total_stock: body.total_stock,
      available_stock: available,
      start_time: body.start_time,
      end_time: body.end_time,
      status: available <= 0 ? 2 : old.status,
    };
    const next = activities.slice();
    next[idx] = updated;
    writeJson(LS_SECKILL_ACTIVITIES, next);
    return { activity: updated };
  }

  return apiPostJson<{ activity: SeckillActivity }>(
    '/activity/update/',
    { token, ...body },
    token,
  );
}

export async function seckillAction(token: string, activityId: string) {
  // activityId 来自 URL params（string），后端 expects int64。
  if (isMockDpcToken(token)) {
    ensureSeed();
    const activities = readJson<SeckillActivity[]>(LS_SECKILL_ACTIVITIES, []);
    const idx = activities.findIndex((a) => a.id === activityId);
    if (idx < 0) throw new Error('秒杀活动不存在');
    if (activities[idx].available_stock <= 0) throw new Error('库存不足/已售罄');

    const activity = { ...activities[idx] };
    activity.available_stock = activity.available_stock - 1;
    if (activity.available_stock <= 0) activity.status = 2;

    const nextActivities = activities.slice();
    nextActivities[idx] = activity;
    writeJson(LS_SECKILL_ACTIVITIES, nextActivities);

    const orders = readJson<Order[]>(LS_ORDERS, []);
    const nextId = orders.reduce((m, o) => Math.max(m, o.id), 0) + 1;
    const orderNo = `SK${nextId}`;
    const order: Order = {
      id: nextId,
      order_no: orderNo,
      user_id: MOCK_DPC_USER_ID,
      product_id: Number(activity.product_id),
      product_name: activity.product_name,
      // mock 下订单列表不展示 activity_id；用 0 避免 JS number 精度问题
      activity_id: 0,
      amount: activity.seckill_price,
      status: 0,
      created_at: nowStr(),
    };
    writeJson(LS_ORDERS, [order, ...orders]);
    return { order_no: orderNo };
  }
  return apiPostJson<{ order_no: string }>(
    '/action/',
    { token, activity_id: activityId },
    token,
  );
}

/**
 * 普通商品下单（当前仅在 dpc/123 mock 模式可用）。
 * 后端未提供普通下单接口时会直接抛错。
 */
export async function productBuy(token: string, productId: string) {
  if (isMockDpcToken(token)) {
    ensureSeed();

    const products = readJson<Product[]>(LS_OWNED_PRODUCTS, []);
    const idx = products.findIndex((p) => p.id === productId);
    if (idx < 0) throw new Error('商品不存在');

    const product = products[idx];
    if (product.stock <= 0) throw new Error('库存不足/已售罄');

    // 扣减库存（仅影响 mock，本地“我的商品”列表会同步展示）
    const nextProducts = products.slice();
    const updatedProduct: Product = { ...product, stock: product.stock - 1 };
    nextProducts[idx] = updatedProduct;
    writeJson(LS_OWNED_PRODUCTS, nextProducts);

    const orders = readJson<Order[]>(LS_ORDERS, []);
    const nextId = orders.reduce((m, o) => Math.max(m, o.id), 0) + 1;
    const orderNo = `N${nextId}`;

    const productIdNum = Number(product.id);
    if (!Number.isFinite(productIdNum)) throw new Error('非法 product_id');

    const order: Order = {
      id: nextId,
      order_no: orderNo,
      user_id: MOCK_DPC_USER_ID,
      product_id: productIdNum,
      product_name: product.name,
      activity_id: 0,
      amount: product.price,
      status: 0,
      created_at: nowStr(),
    };

    writeJson(LS_ORDERS, [order, ...orders]);
    return { order_no: orderNo };
  }

  return apiPostJson<{ order_no: string }>(
    '/order/buy/',
    { token, product_id: productId },
    token,
  );
}

export async function orderList(token: string, page: number, size: number, status = -1) {
  if (isMockDpcToken(token)) {
    ensureSeed();
    const orders = readJson<Order[]>(LS_ORDERS, []);
    const filtered = status >= 0 ? orders.filter((o) => o.status === status) : orders;
    const total = filtered.length;
    const sorted = filtered.slice().sort((a, b) => b.id - a.id);
    const offset = (page - 1) * size;
    const list = sorted.slice(offset, offset + size);
    return { order_list: list, total };
  }
  const q = new URLSearchParams({
    page: String(page),
    size: String(size),
    status: String(status),
  });
  return apiGet<{ order_list: Order[]; total: number }>(`/order/list/?${q}`, token);
}

export async function orderDetail(token: string, orderNo: string) {
  if (isMockDpcToken(token)) {
    ensureSeed();
    const orders = readJson<Order[]>(LS_ORDERS, []);
    const o = orders.find((x) => x.order_no === orderNo);
    if (!o) throw new Error('订单不存在');
    return { order: o };
  }
  const q = new URLSearchParams({ order_no: orderNo });
  return apiGet<{ order: Order }>(`/order/detail/?${q}`, token);
}

export async function orderPay(token: string, orderNo: string) {
  if (isMockDpcToken(token)) {
    ensureSeed();
    const orders = readJson<Order[]>(LS_ORDERS, []);
    const idx = orders.findIndex((o) => o.order_no === orderNo);
    if (idx < 0) throw new Error('订单不存在');
    orders[idx] = { ...orders[idx], status: 1 };
    writeJson(LS_ORDERS, orders);
    return {};
  }
  return apiPostJson<Record<string, unknown>>('/order/pay/', { token, order_no: orderNo }, token);
}

export async function orderCancel(token: string, orderNo: string) {
  if (isMockDpcToken(token)) {
    ensureSeed();
    const orders = readJson<Order[]>(LS_ORDERS, []);
    const idx = orders.findIndex((o) => o.order_no === orderNo);
    if (idx < 0) throw new Error('订单不存在');
    orders[idx] = { ...orders[idx], status: 2 };
    writeJson(LS_ORDERS, orders);
    return {};
  }
  return apiPostJson<Record<string, unknown>>('/order/cancel/', { token, order_no: orderNo }, token);
}
