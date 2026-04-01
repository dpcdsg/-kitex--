import { useEffect, useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { orderCancel, orderList, orderPay } from '@/api/seckill';
import { useAuth } from '@/context/AuthContext';
import type { Order } from '@/api/types';
import { fenToYuan } from '@/util/money';

const STATUS: Record<number, string> = {
  0: '待支付',
  1: '已支付',
  2: '已取消',
  3: '已过期',
};

export function OrdersPage() {
  const nav = useNavigate();
  const { token } = useAuth();
  const [orders, setOrders] = useState<Order[]>([]);
  const [err, setErr] = useState<string | null>(null);
  const [busy, setBusy] = useState<string | null>(null);

  const mockOrders: Order[] = [
    {
      id: 1,
      order_no: 'SK10001',
      user_id: 10001,
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
      user_id: 10001,
      product_id: 102,
      product_name: '模拟商品 B',
      activity_id: 1000002,
      amount: 9900,
      status: 1,
      created_at: '2026-01-01 13:00:00',
    },
  ];

  useEffect(() => {
    if (!token) {
      nav('/login');
      return;
    }

    let cancelled = false;
    (async () => {
      try {
        const r = await orderList(token, 1, 50);
        if (!cancelled) {
          setOrders(r.order_list || []);
          setErr(null);
        }
      } catch (e) {
        // 前端联调 mock：如果后端不可用/鉴权失败，则展示静态订单列表验证页面效果。
        if (!cancelled) {
          setOrders(mockOrders);
          setErr(null);
        }
      }
    })();
    return () => {
      cancelled = true;
    };
  }, [token, nav]);

  async function pay(orderNo: string) {
    if (!token) return;
    setBusy(orderNo);
    try {
      await orderPay(token, orderNo);
      const r = await orderList(token, 1, 50);
      setOrders(r.order_list || []);
    } catch (e) {
      alert(e instanceof Error ? e.message : '支付失败');
    } finally {
      setBusy(null);
    }
  }

  async function cancel(orderNo: string) {
    if (!token) return;
    setBusy(orderNo);
    try {
      await orderCancel(token, orderNo);
      const r = await orderList(token, 1, 50);
      setOrders(r.order_list || []);
    } catch (e) {
      alert(e instanceof Error ? e.message : '取消失败');
    } finally {
      setBusy(null);
    }
  }

  if (!token) return null;

  return (
    <div className="container">
      <h1 style={{ fontSize: '1.35rem' }}>我的订单</h1>
      {err && <p className="err">{err}</p>}
      <table style={{ width: '100%', borderCollapse: 'collapse', background: '#fff', borderRadius: 8 }}>
        <thead>
          <tr style={{ borderBottom: '1px solid #eee', textAlign: 'left' }}>
            <th style={{ padding: 10 }}>订单号</th>
            <th style={{ padding: 10 }}>商品</th>
            <th style={{ padding: 10 }}>金额</th>
            <th style={{ padding: 10 }}>状态</th>
            <th style={{ padding: 10 }}>时间</th>
            <th style={{ padding: 10 }}>操作</th>
          </tr>
        </thead>
        <tbody>
          {orders.map((o) => (
            <tr key={o.order_no} style={{ borderBottom: '1px solid #f5f5f5' }}>
              <td style={{ padding: 10, fontSize: '0.9rem' }}>{o.order_no}</td>
              <td style={{ padding: 10 }}>{o.product_name}</td>
              <td style={{ padding: 10 }}>{fenToYuan(o.amount)}</td>
              <td style={{ padding: 10 }}>{STATUS[o.status] ?? o.status}</td>
              <td style={{ padding: 10, fontSize: '0.85rem' }} className="muted">
                {o.created_at}
              </td>
              <td style={{ padding: 10 }}>
                {o.status === 0 && (
                  <>
                    <button
                      type="button"
                      className="btn-primary"
                      style={{ padding: '4px 10px', fontSize: '0.85rem', marginRight: 6 }}
                      disabled={busy === o.order_no}
                      onClick={() => pay(o.order_no)}
                    >
                      支付
                    </button>
                    <button
                      type="button"
                      className="btn-ghost"
                      style={{ padding: '4px 10px', fontSize: '0.85rem' }}
                      disabled={busy === o.order_no}
                      onClick={() => cancel(o.order_no)}
                    >
                      取消
                    </button>
                  </>
                )}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
      {orders.length === 0 && !err && <p className="muted">暂无订单</p>}
      <p style={{ marginTop: 16 }}>
        <Link to="/">返回首页</Link>
      </p>
    </div>
  );
}
