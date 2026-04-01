import { useEffect, useState } from 'react';
import { Link, useParams } from 'react-router-dom';
import { productDetail } from '@/api/seckill';
import { useAuth } from '@/context/AuthContext';
import type { Product } from '@/api/types';
import { fenToYuan } from '@/util/money';

const PLACEHOLDER = 'https://via.placeholder.com/400x400/f5f5f5/ff5000?text=陶宝';

export function ProductDetailPage() {
  const { id } = useParams();
  const { token } = useAuth();
  const [p, setP] = useState<Product | null>(null);
  const [err, setErr] = useState<string | null>(null);

  useEffect(() => {
    if (!id) return;
    let cancelled = false;
    (async () => {
      try {
        const r = await productDetail(id, token);
        if (!cancelled) {
          setP(r.product);
          setErr(null);
        }
      } catch (e) {
        if (!cancelled) setErr(e instanceof Error ? e.message : '加载失败');
      }
    })();
    return () => {
      cancelled = true;
    };
  }, [id, token]);

  if (err) {
    return (
      <div className="container">
        <p className="err">{err}</p>
        <Link to="/">返回首页</Link>
      </div>
    );
  }

  if (!p) {
    return (
      <div className="container">
        <p className="muted">加载中…</p>
      </div>
    );
  }

  return (
    <div className="container">
      <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 24, alignItems: 'start' }}>
        <div className="card" style={{ padding: 0 }}>
          <img
            src={p.image_url || PLACEHOLDER}
            alt={p.name}
            style={{ width: '100%', display: 'block' }}
            onError={(e) => {
              (e.target as HTMLImageElement).src = PLACEHOLDER;
            }}
          />
        </div>
        <div>
          <h1 style={{ fontSize: '1.35rem', marginTop: 0 }}>{p.name}</h1>
          <div className="price" style={{ fontSize: '1.75rem' }}>
            {fenToYuan(p.price)}
          </div>
          <p className="muted">库存 {p.stock}</p>
          {p.category && <p className="muted">分类：{p.category}</p>}
          <div style={{ marginTop: 16, padding: 12, background: '#fafafa', borderRadius: 8 }}>
            {p.description || '暂无描述'}
          </div>
          <p style={{ marginTop: 24 }}>
            <Link to="/">返回首页</Link>
          </p>
        </div>
      </div>
    </div>
  );
}
