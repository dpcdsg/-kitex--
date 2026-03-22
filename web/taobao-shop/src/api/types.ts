export type User = {
  id: number;
  name: string;
  avatar?: string;
  signature?: string;
};

export type Product = {
  id: number;
  name: string;
  description: string;
  price: number;
  stock: number;
  image_url: string;
  category: string;
};

export type SeckillActivity = {
  id: number;
  product_id: number;
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
