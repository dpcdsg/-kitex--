import { useEffect, useMemo, useState } from 'react';
import { Link } from 'react-router-dom';
import { productCreate, productList, productUpdate, seckillCreate, seckillList } from '@/api/seckill';
import { useAuth } from '@/context/AuthContext';
import type { Product, SeckillActivity } from '@/api/types';

export function SellerPage() {
  const { token, userId } = useAuth();
  const [err, setErr] = useState<string | null>(null);
  const [ok, setOk] = useState<string | null>(null);
  const [products, setProducts] = useState<Product[]>([]);
  const [productsLoading, setProductsLoading] = useState(false);
  const [productsIsMock, setProductsIsMock] = useState(false);

  const mockProducts: Product[] = useMemo(
    () => [
      {
        id: 101,
        name: '模拟商品 A',
        description: '用于验证前端界面的演示商品。',
        price: 19900,
        stock: 20,
        image_url: '',
        category: '数码',
      },
      {
        id: 102,
        name: '模拟商品 B',
        description: '当后端不可用或无数据时，这些数据会显示。',
        price: 9900,
        stock: 80,
        image_url: '',
        category: '家居',
      },
    ],
    [],
  );

  const productsStorageKey = useMemo(() => {
    if (!userId) return null;
    return `taobao_demo_owned_products_${userId}`;
  }, [userId]);

  function loadOwnedProducts(): Product[] {
    if (!productsStorageKey) return [];
    try {
      const raw = localStorage.getItem(productsStorageKey);
      if (!raw) return [];
      const arr = JSON.parse(raw) as Product[];
      return Array.isArray(arr) ? arr : [];
    } catch {
      return [];
    }
  }

  function saveOwnedProducts(ps: Product[]) {
    if (!productsStorageKey) return;
    localStorage.setItem(productsStorageKey, JSON.stringify(ps));
  }

  useEffect(() => {
    let cancelled = false;

    async function loadProducts() {
      if (!token || !userId) {
        setProductsIsMock(true);
        setProducts(mockProducts);
        return;
      }

      setProductsLoading(true);
      setProductsIsMock(false);
      setErr(null);

      const ownedLocal = loadOwnedProducts();
      try {
        const r = await productList(1, 50, token);
        const list = r.product_list ?? [];

        const sellerFiltered = (list as unknown as Array<Product & { seller_id?: number }>).filter((p) => {
          const sid = p.seller_id ?? 0;
          return sid === userId && sid !== 0;
        });

        // 如果后端没有返回 seller_id 或过滤结果为空，回退到本地已创建商品
        const next = sellerFiltered.length > 0 ? sellerFiltered : ownedLocal;
        if (!cancelled) {
          if (next.length > 0) {
            setProducts(next);
            setProductsIsMock(false);
          } else {
            const mockedForMe = mockProducts.map((p) => ({ ...p, seller_id: userId }));
            setProducts(mockedForMe);
            setProductsIsMock(true);
          }
        }
      } catch {
        if (!cancelled) {
          if (ownedLocal.length > 0) {
            setProducts(ownedLocal);
            setProductsIsMock(false);
          } else {
            const mockedForMe = mockProducts.map((p) => ({ ...p, seller_id: userId }));
            setProducts(mockedForMe);
            setProductsIsMock(true);
          }
        }
      } finally {
        if (!cancelled) setProductsLoading(false);
      }
    }

    loadProducts();
    return () => {
      cancelled = true;
    };
  }, [token, userId, mockProducts, productsStorageKey]);

  const [activities, setActivities] = useState<SeckillActivity[]>([]);
  const [activitiesLoading, setActivitiesLoading] = useState(false);
  const [activitiesIsMock, setActivitiesIsMock] = useState(false);

  useEffect(() => {
    let cancelled = false;

    async function loadActivities() {
      if (!products.length) {
        setActivities([]);
        setActivitiesIsMock(true);
        return;
      }

      const myProductIds = new Set(products.map((p) => p.id));
      setActivitiesLoading(true);
      setActivitiesIsMock(false);

      try {
        const r = await seckillList(1, 50, token);
        const list = (r.activity_list ?? []).filter((a) => myProductIds.has(a.product_id));
        if (!cancelled) setActivities(list);
      } catch {
        if (cancelled) return;
        // 后端不可用时，用“商品列表”生成一个最小活动展示
        const fallback = products.slice(0, 2).map((p, idx) => ({
          id: 900000 + idx,
          product_id: p.id,
          product_name: p.name,
          seckill_price: 5900 + idx * 100,
          total_stock: p.stock,
          available_stock: p.stock,
          start_time: '2026-03-10 10:00:00',
          end_time: '2026-03-10 11:00:00',
          status: 1,
        }));
        setActivities(fallback);
        setActivitiesIsMock(true);
      } finally {
        if (!cancelled) setActivitiesLoading(false);
      }
    }

    if (!token) {
      setActivitiesIsMock(true);
      setActivities([]);
      return;
    }

    loadActivities();
    return () => {
      cancelled = true;
    };
  }, [token, products]);

  const [showCreateProductModal, setShowCreateProductModal] = useState(false);
  const [showCreateSeckillModal, setShowCreateSeckillModal] = useState(false);

  const [editingId, setEditingId] = useState<number | null>(null);
  const editingProduct = useMemo(() => products.find((p) => p.id === editingId) ?? null, [products, editingId]);

  const [eName, setEName] = useState('');
  const [eDesc, setEDesc] = useState('');
  const [ePrice, setEPrice] = useState('');
  const [eStock, setEStock] = useState('');
  const [eImg, setEImg] = useState('');
  const [eCat, setECat] = useState('');

  function startEdit(p: Product) {
    setEditingId(p.id);
    setEName(p.name);
    setEDesc(p.description);
    setEPrice(String(p.price));
    setEStock(String(p.stock));
    setEImg(p.image_url);
    setECat(p.category);
    setErr(null);
    setOk(null);
  }

  function cancelEdit() {
    setEditingId(null);
  }

  const canEdit = !!token && !!userId;
  const canEditProducts = canEdit;

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
      setShowCreateProductModal(false);

      // 同步本地“我的商品”列表：用于后续修改接口只操作自己发布的商品
      const next: Product = {
        ...(r.product as unknown as Product),
        id: r.product.id,
        seller_id: userId ?? undefined,
      };
      setProducts((prev) => {
        const filtered = prev.filter((x) => x.id !== next.id);
        const merged = [next, ...filtered];
        saveOwnedProducts(merged);
        setProductsIsMock(false);
        return merged;
      });
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
      setShowCreateSeckillModal(false);
    } catch (e) {
      setErr(e instanceof Error ? e.message : '创建失败');
    }
  }

  const now = new Date();
  const defaultStart = new Date(now.getTime() + 60_000).toISOString().slice(0, 16);
  const defaultEnd = new Date(now.getTime() + 3600_000).toISOString().slice(0, 16);

  async function onSaveEdit(e: React.FormEvent) {
    e.preventDefault();
    if (!editingProduct) return;

    setErr(null);
    setOk(null);

    const price = parseInt(ePrice, 10);
    const stock = parseInt(eStock, 10);
    if (Number.isNaN(price) || Number.isNaN(stock)) {
      setErr('价格、库存应为整数（单位：分）');
      return;
    }

    const next: Product = {
      ...editingProduct,
      name: eName,
      description: eDesc,
      price,
      stock,
      image_url: eImg,
      category: eCat,
    };

    const ownedNext = products.map((p) => (p.id === next.id ? next : p));

    // 先更新本地展示，保证前端效果不被后端影响
    setProducts((prev) => prev.map((p) => (p.id === next.id ? next : p)));
    saveOwnedProducts(ownedNext);

    if (!token) {
      setOk('已在本地更新（未登录，未调用后端接口）');
      cancelEdit();
      return;
    }

    try {
      const r = await productUpdate(token, {
        product_id: next.id,
        name: next.name,
        description: next.description,
        price: next.price,
        stock: next.stock,
        image_url: next.image_url,
        category: next.category,
      });

      const updated = r.product as unknown as Product;
      const ownedUpdated = ownedNext.map((p) => (p.id === next.id ? updated : p));
      setProducts((prev) => prev.map((p) => (p.id === next.id ? updated : p)));
      saveOwnedProducts(ownedUpdated);
      setOk('修改成功');
      cancelEdit();
    } catch (ex) {
      setErr(ex instanceof Error ? ex.message : '修改失败，已保留本地展示数据');
      cancelEdit();
    }
  }

  return (
    <div className="container" style={{ maxWidth: 560 }}>
      <h1 style={{ fontSize: '1.35rem' }}>卖家中心</h1>
      <p className="muted">创建商品与秒杀活动，供首页展示与下单联调。</p>
      {!canEdit && <p className="muted">未登录时也会展示模拟商品，用于验证页面效果；修改功能将被禁用。</p>}
      {err && <p className="err">{err}</p>}
      {ok && <p style={{ color: '#080' }}>{ok}</p>}

      <section className="card" style={{ padding: 20, marginBottom: 16 }}>
        <h2 style={{ fontSize: '1.1rem', marginTop: 0 }}>发布</h2>
        <div style={{ display: 'flex', gap: 12, flexWrap: 'wrap' }}>
          <button
            type="button"
            className="btn-primary"
            disabled={!token}
            onClick={() => setShowCreateProductModal(true)}
          >
            发布商品
          </button>
          <button
            type="button"
            className="btn-primary"
            disabled={!token}
            onClick={() => {
              if (!sProductId && products.length > 0) setSProductId(String(products[0].id));
              setShowCreateSeckillModal(true);
            }}
          >
            发布秒杀活动
          </button>
        </div>
        <p className="muted" style={{ marginTop: 10 }}>
          主界面用于展示你已发布的商品与秒杀活动；点击按钮弹窗填写发布信息。
        </p>
      </section>

      {showCreateProductModal && (
        <div
          style={{
            position: 'fixed',
            inset: 0,
            background: 'rgba(0,0,0,0.35)',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            zIndex: 1000,
          }}
          onClick={() => setShowCreateProductModal(false)}
        >
          <div className="card" style={{ width: 'min(720px, 92vw)', padding: 20 }} onClick={(e) => e.stopPropagation()}>
            <h2 style={{ fontSize: '1.1rem', marginTop: 0 }}>发布商品</h2>
            <form onSubmit={onCreateProduct}>
              <table style={{ width: '100%', borderCollapse: 'collapse' }}>
                <tbody>
                  <tr>
                    <td style={{ padding: 8, width: 140 }}>名称</td>
                    <td style={{ padding: 8 }}>
                      <input value={pName} onChange={(e) => setPName(e.target.value)} required />
                    </td>
                  </tr>
                  <tr>
                    <td style={{ padding: 8 }}>描述</td>
                    <td style={{ padding: 8 }}>
                      <textarea value={pDesc} onChange={(e) => setPDesc(e.target.value)} rows={2} />
                    </td>
                  </tr>
                  <tr>
                    <td style={{ padding: 8 }}>价格（分）</td>
                    <td style={{ padding: 8 }}>
                      <input value={pPrice} onChange={(e) => setPPrice(e.target.value)} required />
                    </td>
                  </tr>
                  <tr>
                    <td style={{ padding: 8 }}>库存</td>
                    <td style={{ padding: 8 }}>
                      <input value={pStock} onChange={(e) => setPStock(e.target.value)} required />
                    </td>
                  </tr>
                  <tr>
                    <td style={{ padding: 8 }}>图片 URL</td>
                    <td style={{ padding: 8 }}>
                      <input value={pImg} onChange={(e) => setPImg(e.target.value)} placeholder="https://..." />
                    </td>
                  </tr>
                  <tr>
                    <td style={{ padding: 8 }}>分类</td>
                    <td style={{ padding: 8 }}>
                      <input value={pCat} onChange={(e) => setPCat(e.target.value)} />
                    </td>
                  </tr>
                </tbody>
              </table>
              <div style={{ display: 'flex', gap: 12, marginTop: 12 }}>
                <button type="submit" className="btn-primary" disabled={!token}>
                  创建商品
                </button>
                <button type="button" className="btn" onClick={() => setShowCreateProductModal(false)}>
                  关闭
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {showCreateSeckillModal && (
        <div
          style={{
            position: 'fixed',
            inset: 0,
            background: 'rgba(0,0,0,0.35)',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            zIndex: 1000,
          }}
          onClick={() => setShowCreateSeckillModal(false)}
        >
          <div className="card" style={{ width: 'min(720px, 92vw)', padding: 20 }} onClick={(e) => e.stopPropagation()}>
            <h2 style={{ fontSize: '1.1rem', marginTop: 0 }}>发布秒杀活动</h2>
            <form onSubmit={onCreateSeckill}>
              <table style={{ width: '100%', borderCollapse: 'collapse' }}>
                <tbody>
                  <tr>
                    <td style={{ padding: 8, width: 140 }}>商品</td>
                    <td style={{ padding: 8 }}>
                      <select
                        value={sProductId}
                        onChange={(e) => setSProductId(e.target.value)}
                        style={{ width: '100%' }}
                        required
                      >
                        <option value="" disabled>
                          请选择商品
                        </option>
                        {products.map((p) => (
                          <option key={p.id} value={p.id}>
                            {p.name}（{p.id}）
                          </option>
                        ))}
                      </select>
                    </td>
                  </tr>
                  <tr>
                    <td style={{ padding: 8 }}>秒杀价（分）</td>
                    <td style={{ padding: 8 }}>
                      <input value={sPrice} onChange={(e) => setSPrice(e.target.value)} required />
                    </td>
                  </tr>
                  <tr>
                    <td style={{ padding: 8 }}>活动库存</td>
                    <td style={{ padding: 8 }}>
                      <input value={sStock} onChange={(e) => setSStock(e.target.value)} required />
                    </td>
                  </tr>
                  <tr>
                    <td style={{ padding: 8 }}>开始时间</td>
                    <td style={{ padding: 8 }}>
                      <input
                        type="datetime-local"
                        value={sStart || defaultStart}
                        onChange={(e) => setSStart(e.target.value)}
                        required
                      />
                    </td>
                  </tr>
                  <tr>
                    <td style={{ padding: 8 }}>结束时间</td>
                    <td style={{ padding: 8 }}>
                      <input
                        type="datetime-local"
                        value={sEnd || defaultEnd}
                        onChange={(e) => setSEnd(e.target.value)}
                        required
                      />
                    </td>
                  </tr>
                </tbody>
              </table>
              <div style={{ display: 'flex', gap: 12, marginTop: 12 }}>
                <button type="submit" className="btn-primary" disabled={!token}>
                  创建活动
                </button>
                <button type="button" className="btn" onClick={() => setShowCreateSeckillModal(false)}>
                  关闭
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      <section className="card" style={{ padding: 20, marginTop: 16 }}>
        <h2 style={{ fontSize: '1.1rem', marginTop: 0 }}>商品列表</h2>
        {productsLoading ? (
          <p className="muted">加载中…</p>
        ) : (
          <p className="muted">
            {canEdit ? '仅允许编辑你创建的商品。后端不可用时会保持本地展示。' : productsIsMock ? '未登录：展示模拟数据' : '未登录：展示数据'}
          </p>
        )}

        {products.length === 0 ? (
          <p className="muted">暂无商品</p>
        ) : (
          <div style={{ display: 'grid', gap: 10, marginTop: 12 }}>
            {products.map((p) => (
              <div key={p.id} style={{ border: '1px solid #eee', borderRadius: 8, padding: 12 }}>
                <div style={{ display: 'flex', justifyContent: 'space-between', gap: 12 }}>
                  <div>
                    <div style={{ fontWeight: 700 }}>
                      {p.name}（ID：{p.id}）
                    </div>
                    <div className="muted" style={{ marginTop: 4 }}>
                      分类：{p.category || '-'} · 价格：{p.price} · 库存：{p.stock}
                    </div>
                  </div>
                  <div>
                    <button
                      type="button"
                      className="btn-primary"
              disabled={!canEditProducts}
                      onClick={() => startEdit(p)}
                    >
                      修改
                    </button>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </section>

      <section className="card" style={{ padding: 20, marginTop: 16 }}>
        <h2 style={{ fontSize: '1.1rem', marginTop: 0 }}>秒杀活动</h2>
        {activitiesLoading ? (
          <p className="muted">加载中…</p>
        ) : (
          <p className="muted">{activitiesIsMock ? '未从后端获取活动，展示本地演示数据。' : '展示你已发布的秒杀活动。'}</p>
        )}
        {activities.length === 0 ? (
          <p className="muted">暂无秒杀活动</p>
        ) : (
          <div style={{ display: 'grid', gap: 10, marginTop: 12 }}>
            {activities.map((a) => (
              <div key={a.id} style={{ border: '1px solid #eee', borderRadius: 8, padding: 12 }}>
                <div style={{ fontWeight: 700 }}>
                  活动 ID：{a.id} · 商品：{a.product_name}（{a.product_id}）
                </div>
                <div className="muted" style={{ marginTop: 4 }}>
                  秒杀价：{a.seckill_price} · 库存：{a.available_stock} / {a.total_stock} · 状态：{a.status}
                </div>
                <div className="muted" style={{ marginTop: 4 }}>
                  {a.start_time} ~ {a.end_time}
                </div>
              </div>
            ))}
          </div>
        )}
      </section>

      {editingProduct && (
        <section className="card" style={{ padding: 20, marginTop: 16 }}>
          <h2 style={{ fontSize: '1.1rem', marginTop: 0 }}>修改商品</h2>
          <form onSubmit={onSaveEdit}>
            <div className="form-row">
              <label>名称</label>
              <input value={eName} onChange={(e) => setEName(e.target.value)} required />
            </div>
            <div className="form-row">
              <label>描述</label>
              <textarea value={eDesc} onChange={(e) => setEDesc(e.target.value)} rows={2} />
            </div>
            <div className="form-row">
              <label>价格（分）</label>
              <input value={ePrice} onChange={(e) => setEPrice(e.target.value)} required />
            </div>
            <div className="form-row">
              <label>库存</label>
              <input value={eStock} onChange={(e) => setEStock(e.target.value)} required />
            </div>
            <div className="form-row">
              <label>图片 URL（可空）</label>
              <input value={eImg} onChange={(e) => setEImg(e.target.value)} placeholder="https://..." />
            </div>
            <div className="form-row">
              <label>分类</label>
              <input value={eCat} onChange={(e) => setECat(e.target.value)} />
            </div>
            <div style={{ display: 'flex', gap: 12, marginTop: 12 }}>
              <button type="submit" className="btn-primary" disabled={!canEditProducts}>
                保存修改
              </button>
              <button type="button" className="btn" onClick={cancelEdit}>
                取消
              </button>
            </div>
          </form>
        </section>
      )}

      <p style={{ marginTop: 16 }}>
        <Link to="/">返回首页</Link>
      </p>
    </div>
  );
}
