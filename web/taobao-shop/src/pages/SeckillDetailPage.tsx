import { useEffect, useState } from 'react';
import { Link, useNavigate, useParams } from 'react-router-dom';
import { seckillAction, seckillDetail } from '@/api/seckill';
import { useAuth } from '@/context/AuthContext';
import type { SeckillActivity } from '@/api/types';
import { fenToYuan } from '@/util/money';

export function SeckillDetailPage() {
  const { id } = useParams();
  const nav = useNavigate();
  const { token } = useAuth();
  const [a, setA] = useState<SeckillActivity | null>(null);
  const [err, setErr] = useState<string | null>(null);
  const [msg, setMsg] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (!id) return;
    let cancelled = false;
    (async () => {
      try {
        const r = await seckillDetail(id, token);
        if (!cancelled) {
          setA(r.activity);
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

  async function onSeckill() {
    if (!token || !id) {
      nav('/login');
      return;
    }
    setMsg(null);
    setLoading(true);
    try {
      const r = await seckillAction(token, id);
      setMsg(`下单成功，订单号：${r.order_no}`);
    } catch (e) {
      setMsg(e instanceof Error ? e.message : '秒杀失败');
    } finally {
      setLoading(false);
    }
  }

  if (err) {
    return (
      <div className="container">
        <p className="err">{err}</p>
        <Link to="/">返回首页</Link>
      </div>
    );
  }

  if (!a) {
    return (
      <div className="container">
        <p className="muted">加载中…</p>
      </div>
    );
  }

  return (
    <div className="container" style={{ maxWidth: 640 }}>
      <h1 style={{ fontSize: '1.35rem' }}>{a.product_name}</h1>
      <div className="price" style={{ fontSize: '1.75rem' }}>
        秒杀价 {fenToYuan(a.seckill_price)}
      </div>
      <p className="muted">
        活动库存 {a.available_stock} / {a.total_stock}
      </p>
      <p className="muted">
        {a.start_time} ~ {a.end_time} · 状态 {a.status}
      </p>
      <p>
        <Link to={`/product/${a.product_id}`}>查看原商品</Link>
      </p>
      {msg && <p style={{ color: msg.startsWith('下单') ? '#080' : '#c00' }}>{msg}</p>}
      <button type="button" className="btn-primary" style={{ marginTop: 12 }} onClick={onSeckill} disabled={loading}>
        {loading ? '提交中…' : '立即秒杀'}
      </button>
      <p style={{ marginTop: 16 }}>
        <Link to="/">返回首页</Link>
      </p>
    </div>
  );
}
