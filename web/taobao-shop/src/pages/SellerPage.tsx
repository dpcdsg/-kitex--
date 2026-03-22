import { useState } from 'react';
import { Link } from 'react-router-dom';
import { productCreate, seckillCreate } from '@/api/seckill';
import { useAuth } from '@/context/AuthContext';

export function SellerPage() {
  const { token } = useAuth();
  const [err, setErr] = useState<string | null>(null);
  const [ok, setOk] = useState<string | null>(null);

  const [pName, setPName] = useState('演示商品');
  const [pDesc, setPDesc] = useState('本地陶宝演示');
  const [pPrice, setPPrice] = useState('9900');
  const [pStock, setPStock] = useState('100');
  const [pImg, setPImg] = useState('');
  const [pCat, setPCat] = useState('数码');

  const [sProductId, setSProductId] = useState('');
  const [sPrice, setSPrice] = useState('5900');
  const [sStock, setSStock] = useState('50');
  const [sStart, setSStart] = useState('');
  const [sEnd, setSEnd] = useState('');

  if (!token) {
    return (
      <div className="container">
        <p>请先登录后发布商品。</p>
        <Link to="/login">去登录</Link>
      </div>
    );
  }

  async function onCreateProduct(e: React.FormEvent) {
    e.preventDefault();
    setErr(null);
    setOk(null);
    try {
      const price = parseInt(pPrice, 10);
      const stock = parseInt(pStock, 10);
      if (Number.isNaN(price) || Number.isNaN(stock)) {
        setErr('价格、库存应为整数（单位：分）');
        return;
      }
      const r = await productCreate(token!, {
        name: pName,
        description: pDesc,
        price,
        stock,
        image_url: pImg,
        category: pCat,
      });
      setOk(`商品创建成功，ID：${r.product.id}`);
      setSProductId(String(r.product.id));
    } catch (e) {
      setErr(e instanceof Error ? e.message : '创建失败');
    }
  }

  async function onCreateSeckill(e: React.FormEvent) {
    e.preventDefault();
    setErr(null);
    setOk(null);
    try {
      const pid = parseInt(sProductId, 10);
      const sp = parseInt(sPrice, 10);
      const ts = parseInt(sStock, 10);
      if (Number.isNaN(pid) || pid <= 0) {
        setErr('请填写有效的商品 ID');
        return;
      }
      if (!sStart || !sEnd) {
        setErr('请选择开始与结束时间');
        return;
      }
      const start = new Date(sStart).toISOString();
      const end = new Date(sEnd).toISOString();
      const r = await seckillCreate(token!, {
        product_id: pid,
        seckill_price: sp,
        total_stock: ts,
        start_time: start,
        end_time: end,
      });
      setOk(`秒杀活动创建成功，活动 ID：${r.activity.id}`);
    } catch (e) {
      setErr(e instanceof Error ? e.message : '创建失败');
    }
  }

  const now = new Date();
  const defaultStart = new Date(now.getTime() + 60_000).toISOString().slice(0, 16);
  const defaultEnd = new Date(now.getTime() + 3600_000).toISOString().slice(0, 16);

  return (
    <div className="container" style={{ maxWidth: 560 }}>
      <h1 style={{ fontSize: '1.35rem' }}>卖家中心</h1>
      <p className="muted">创建商品与秒杀活动，供首页展示与下单联调。</p>
      {err && <p className="err">{err}</p>}
      {ok && <p style={{ color: '#080' }}>{ok}</p>}

      <section className="card" style={{ padding: 20, marginBottom: 16 }}>
        <h2 style={{ fontSize: '1.1rem', marginTop: 0 }}>发布商品</h2>
        <form onSubmit={onCreateProduct}>
          <div className="form-row">
            <label>名称</label>
            <input value={pName} onChange={(e) => setPName(e.target.value)} required />
          </div>
          <div className="form-row">
            <label>描述</label>
            <textarea value={pDesc} onChange={(e) => setPDesc(e.target.value)} rows={2} />
          </div>
          <div className="form-row">
            <label>价格（分，如 9900 = ¥99.00）</label>
            <input value={pPrice} onChange={(e) => setPPrice(e.target.value)} required />
          </div>
          <div className="form-row">
            <label>库存</label>
            <input value={pStock} onChange={(e) => setPStock(e.target.value)} required />
          </div>
          <div className="form-row">
            <label>图片 URL（可空）</label>
            <input value={pImg} onChange={(e) => setPImg(e.target.value)} placeholder="https://..." />
          </div>
          <div className="form-row">
            <label>分类</label>
            <input value={pCat} onChange={(e) => setPCat(e.target.value)} />
          </div>
          <button type="submit" className="btn-primary">
            创建商品
          </button>
        </form>
      </section>

      <section className="card" style={{ padding: 20 }}>
        <h2 style={{ fontSize: '1.1rem', marginTop: 0 }}>创建秒杀活动</h2>
        <form onSubmit={onCreateSeckill}>
          <div className="form-row">
            <label>商品 ID</label>
            <input value={sProductId} onChange={(e) => setSProductId(e.target.value)} placeholder="先创建商品" required />
          </div>
          <div className="form-row">
            <label>秒杀价（分）</label>
            <input value={sPrice} onChange={(e) => setSPrice(e.target.value)} required />
          </div>
          <div className="form-row">
            <label>活动库存</label>
            <input value={sStock} onChange={(e) => setSStock(e.target.value)} required />
          </div>
          <div className="form-row">
            <label>开始时间</label>
            <input
              type="datetime-local"
              value={sStart || defaultStart}
              onChange={(e) => setSStart(e.target.value)}
              required
            />
          </div>
          <div className="form-row">
            <label>结束时间</label>
            <input
              type="datetime-local"
              value={sEnd || defaultEnd}
              onChange={(e) => setSEnd(e.target.value)}
              required
            />
          </div>
          <button type="submit" className="btn-primary">
            创建活动
          </button>
        </form>
      </section>

      <p style={{ marginTop: 16 }}>
        <Link to="/">返回首页</Link>
      </p>
    </div>
  );
}
