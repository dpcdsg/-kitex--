import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { productList, seckillList } from '@/api/seckill';
import { useAuth } from '@/context/AuthContext';
import type { Product, SeckillActivity } from '@/api/types';
import { fenToYuan } from '@/util/money';

const PLACEHOLDER = 'https://via.placeholder.com/200x200/f5f5f5/ff5000?text=陶宝';

export function HomePage() {
  const { token } = useAuth();
  const [products, setProducts] = useState<Product[]>([]);
  const [activities, setActivities] = useState<SeckillActivity[]>([]);
  const [err, setErr] = useState<string | null>(null);

  useEffect(() => {
    let cancelled = false;
    (async () => {
      try {
        const [p, s] = await Promise.all([
          productList(1, 24, token),
          seckillList(1, 8, token, -1),
        ]);
        if (!cancelled) {
          setProducts(p.product_list || []);
          setActivities(s.activity_list || []);
          setErr(null);
        }
      } catch (e) {
        if (!cancelled) setErr(e instanceof Error ? e.message : '加载失败');
      }
    })();
    return () => {
      cancelled = true;
    };
  }, [token]);

  return (
    <div className="container">
      <section style={{ marginBottom: 28 }}>
        <h2 style={{ fontSize: '1.25rem', borderLeft: '4px solid #ff5000', paddingLeft: 10, marginBottom: 16 }}>
          限时秒杀
        </h2>
        {err && <p className="err">{err}</p>}
        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(260px, 1fr))', gap: 12 }}>
          {activities.length === 0 && !err && <p className="muted">暂无秒杀活动，请先在卖家中心创建商品与活动。</p>}
          {activities.map((a) => (
            <Link key={a.id} to={`/seckill/${a.id}`} className="card" style={{ padding: 12, display: 'block' }}>
              <div style={{ fontWeight: 600, marginBottom: 8 }}>{a.product_name || `商品 #${a.product_id}`}</div>
              <div className="price">{fenToYuan(a.seckill_price)}</div>
              <div className="muted" style={{ marginTop: 6 }}>
                库存 {a.available_stock} · {a.start_time} 起
              </div>
            </Link>
          ))}
        </div>
      </section>

      <section>
        <h2 style={{ fontSize: '1.25rem', borderLeft: '4px solid #ff5000', paddingLeft: 10, marginBottom: 16 }}>
          猜你喜欢
        </h2>
        <div className="grid-products">
          {products.map((p) => (
            <Link key={p.id} to={`/product/${p.id}`} className="card" style={{ textDecoration: 'none', color: 'inherit' }}>
              <div style={{ aspectRatio: '1', background: '#fafafa', overflow: 'hidden' }}>
                <img
                  src={p.image_url || PLACEHOLDER}
                  alt={p.name}
                  style={{ width: '100%', height: '100%', objectFit: 'cover' }}
                  onError={(e) => {
                    (e.target as HTMLImageElement).src = PLACEHOLDER;
                  }}
                />
              </div>
              <div style={{ padding: 10 }}>
                <div style={{ fontSize: '0.95rem', height: 40, overflow: 'hidden' }}>{p.name}</div>
                <div className="price">{fenToYuan(p.price)}</div>
                {p.category && (
                  <div className="muted" style={{ marginTop: 4 }}>
                    {p.category}
                  </div>
                )}
              </div>
            </Link>
          ))}
        </div>
        {products.length === 0 && !err && <p className="muted">暂无商品。</p>}
      </section>
    </div>
  );
}
