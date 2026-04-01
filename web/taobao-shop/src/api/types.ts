export type User = {
  id: number;
  name: string;
  avatar?: string;
  signature?: string;
};

export type Product = {
  id: string;
  name: string;
  description: string;
  price: number;
  stock: number;
  image_url: string;
  category: string;
  // 后端返回时可能携带 seller_id；前端在“仅编辑自己商品”的场景下可用。
  // 若后端未返回或为 0，前端会回退到本地 mock/本地持久化逻辑。
  seller_id?: string;
};

export type SeckillActivity = {
  id: string;
  product_id: string;
  product_name: string;
  seckill_price: number;
  total_stock: number;
  available_stock: number;
  start_time: string;
  end_time: string;
  status: number;
};

export type Order = {
  id: number;
  order_no: string;
  user_id: number;
  product_id: number;
  product_name: string;
  activity_id: number;
  amount: number;
  status: number;
  created_at: string;
};
