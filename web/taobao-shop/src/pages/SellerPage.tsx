import { useEffect, useMemo, useState } from 'react';
import { Link } from 'react-router-dom';
import {
  MOCK_DPC_USER_ID,
  productCreate,
  productList,
  productUpdate,
  seckillCreate,
  seckillList,
  seckillUpdate,
} from '@/api/seckill';
import { useAuth } from '@/context/AuthContext';
import type { Product, SeckillActivity } from '@/api/types';

const MOCK_PRODUCT_IMAGE_URL = 'https://www.leagueoflegends.com/zh-tw/champions/katarina/';
const DEFAULT_PRODUCT_IMAGE_URL = 'https://via.placeholder.com/640x360/f5f5f5/999999?text=No+Image';

export function SellerPage() {
  const { token, userId } = useAuth();
  const [err, setErr] = useState<string | null>(null);
  const [ok, setOk] = useState<string | null>(null);
  const [products, setProducts] = useState<Product[]>([]);
  const [productsLoading, setProductsLoading] = useState(false);

  const mockProducts: Product[] = useMemo(
    () => [
      {
        id: '101',
        name: '模拟商品 A',
        description: '用于验证前端界面的演示商品。',
        price: 19900,
        stock: 20,
        image_url: MOCK_PRODUCT_IMAGE_URL,
        category: '数码',
      },
      {
        id: '102',
        name: '模拟商品 B',
        description: '当后端不可用或无数据时，这些数据会显示。',
        price: 9900,
        stock: 80,
        image_url: MOCK_PRODUCT_IMAGE_URL,
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
        setProducts([]);
        return;
      }

      setProductsLoading(true);
      setErr(null);

      const ownedLocal = loadOwnedProducts();
      const isDemoUser = userId === MOCK_DPC_USER_ID;
      try {
        const r = await productList(1, 50, token);
        const list = r.product_list ?? [];

        const sellerFiltered = (list as unknown as Array<Product & { seller_id?: string }>).filter((p) => {
          const sid = p.seller_id ?? '';
          return sid === String(userId) && sid !== '';
        });

        // 如果后端没有返回 seller_id 或过滤结果为空，回退到本地已创建商品
        const next = sellerFiltered.length > 0 ? sellerFiltered : ownedLocal;
        if (!cancelled) {
          if (next.length > 0) {
            setProducts(next);
          } else if (isDemoUser) {
            const mockedForMe = mockProducts.map((p) => ({ ...p, seller_id: String(userId) }));
            setProducts(mockedForMe);
          } else {
            setProducts([]);
          }
        }
      } catch {
        if (!cancelled) {
          if (ownedLocal.length > 0) {
            setProducts(ownedLocal);
          } else if (isDemoUser) {
            const mockedForMe = mockProducts.map((p) => ({ ...p, seller_id: String(userId) }));
            setProducts(mockedForMe);
          } else {
            setProducts([]);
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
        setActivitiesIsMock(false);
        return;
      }

      const myProductIds = new Set(products.map((p) => p.id));
      setActivitiesLoading(true);
      setActivitiesIsMock(false);

      try {
        const r = await seckillList(1, 50, token);
        const list = (r.activity_list ?? []).filter((a) => myProductIds.has(String(a.product_id)));
        if (!cancelled) setActivities(list);
      } catch {
        if (cancelled) return;
        if (userId === MOCK_DPC_USER_ID) {
          const fallback = products.slice(0, 2).map((p, idx) => ({
            id: String(900000 + idx),
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
        } else {
          setActivities([]);
          setActivitiesIsMock(false);
        }
      } finally {
        if (!cancelled) setActivitiesLoading(false);
      }
    }

    if (!token) {
      setActivitiesIsMock(false);
      setActivities([]);
      return;
    }

    loadActivities();
    return () => {
      cancelled = true;
    };
  }, [token, userId, products]);

  const [showCreateProductModal, setShowCreateProductModal] = useState(false);
  const [showCreateSeckillModal, setShowCreateSeckillModal] = useState(false);

  const [editingId, setEditingId] = useState<string | null>(null);
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
  const [pImg, setPImg] = useState(DEFAULT_PRODUCT_IMAGE_URL);
  const [pCat, setPCat] = useState('数码');

  const [sProductId, setSProductId] = useState('');
  const [sPrice, setSPrice] = useState('5900');
  const [sStock, setSStock] = useState('50');
  const [sStart, setSStart] = useState('');
  const [sEnd, setSEnd] = useState('');

  const [editingActivityId, setEditingActivityId] = useState<string | null>(null);
  const editingActivity = useMemo(
    () => activities.find((a) => a.id === editingActivityId) ?? null,
    [activities, editingActivityId],
  );
  const [aProductId, setAProductId] = useState('');
  const [aPrice, setAPrice] = useState('');
  const [aStock, setAStock] = useState('');
  const [aStart, setAStart] = useState('');
  const [aEnd, setAEnd] = useState('');

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
        seller_id: userId ? String(userId) : undefined,
      };
      setProducts((prev) => {
        const filtered = prev.filter((x) => x.id !== next.id);
        const merged = [next, ...filtered];
        saveOwnedProducts(merged);
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
      const sp = parseInt(sPrice, 10);
      const ts = parseInt(sStock, 10);
      if (!/^\d+$/.test(sProductId)) {
        setErr('请填写有效的商品 ID');
        return;
      }
      const startInput = sStart || defaultStart;
      const endInput = sEnd || defaultEnd;
      const timeErr = validateSeckillTimeRange(startInput, endInput);
      if (timeErr) {
        setErr(timeErr);
        return;
      }
      const start = new Date(startInput).toISOString();
      const end = new Date(endInput).toISOString();
      const r = await seckillCreate(token!, {
        product_id: sProductId,
        seckill_price: sp,
        total_stock: ts,
        start_time: start,
        end_time: end,
      });
      setOk(`秒杀活动创建成功，活动 ID：${r.activity.id}`);
      setActivities((prev) => [r.activity, ...prev.filter((x) => x.id !== r.activity.id)]);
      setShowCreateSeckillModal(false);
    } catch (e) {
      setErr(e instanceof Error ? e.message : '创建失败');
    }
  }

  const now = new Date();
  const defaultStart = new Date(now.getTime() + 60_000).toISOString().slice(0, 16);
  const defaultEnd = new Date(now.getTime() + 3600_000).toISOString().slice(0, 16);

  function validateSeckillTimeRange(startInput: string, endInput: string): string | null {
    const startDate = new Date(startInput);
    const endDate = new Date(endInput);
    if (Number.isNaN(startDate.getTime()) || Number.isNaN(endDate.getTime())) {
      return '时间格式无效，请重新选择';
    }

    const todayStart = new Date();
    todayStart.setHours(0, 0, 0, 0);
    if (startDate.getTime() < todayStart.getTime() || endDate.getTime() < todayStart.getTime()) {
      return '开始和结束时间不能早于今天';
    }

    if (endDate.getTime() - startDate.getTime() <= 1000) {
      return '结束时间必须晚于开始时间 1 秒以上';
    }
    return null;
  }

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

  function toLocalInputValue(raw: string): string {
    if (!raw) return '';
    const normalized = raw.includes('T') ? raw : raw.replace(' ', 'T');
    const d = new Date(normalized);
    if (Number.isNaN(d.getTime())) return '';
    return d.toISOString().slice(0, 16);
  }

  function startEditActivity(a: SeckillActivity) {
    setEditingActivityId(a.id);
    setAProductId(String(a.product_id));
    setAPrice(String(a.seckill_price));
    setAStock(String(a.total_stock));
    setAStart(toLocalInputValue(a.start_time));
    setAEnd(toLocalInputValue(a.end_time));
    setErr(null);
    setOk(null);
  }

  function cancelEditActivity() {
    setEditingActivityId(null);
  }

  async function onSaveActivityEdit(e: React.FormEvent) {
    e.preventDefault();
    if (!editingActivity || !token) return;
    setErr(null);
    setOk(null);

    const productId = aProductId;
    const seckillPrice = parseInt(aPrice, 10);
    const totalStock = parseInt(aStock, 10);
    if (!/^\d+$/.test(productId) || Number.isNaN(seckillPrice) || Number.isNaN(totalStock)) {
      setErr('活动参数必须为有效数字');
      return;
    }
    if (!aStart || !aEnd) {
      setErr('请填写活动开始和结束时间');
      return;
    }

    const startISO = new Date(aStart).toISOString();
    const endISO = new Date(aEnd).toISOString();
    try {
      const r = await seckillUpdate(token, {
        activity_id: editingActivity.id,
        product_id: productId,
        seckill_price: seckillPrice,
        total_stock: totalStock,
        start_time: startISO,
        end_time: endISO,
      });
      setActivities((prev) => prev.map((x) => (x.id === r.activity.id ? r.activity : x)));
      setOk('活动修改成功');
      cancelEditActivity();
    } catch (ex) {
      setErr(ex instanceof Error ? ex.message : '活动修改失败');
    }
  }

  return (
    <div className="container" style={{ maxWidth: 1120 }}>
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
            {canEdit
              ? '仅允许编辑你创建的商品。后端不可用时会保持本地展示。'
              : '请登录后管理商品；演示账号 dpc/123 在无数据时可看到示例商品。'}
          </p>
        )}

        {products.length === 0 ? (
          <p className="muted">暂无商品</p>
        ) : (
          <div style={{ display: 'flex', gap: 12, marginTop: 12, flexWrap: 'wrap' }}>
            {products.map((p) => (
              <div
                key={p.id}
                style={{
                  border: '1px solid #eee',
                  borderRadius: 8,
                  padding: 12,
                  flex: '1 1 calc(50% - 12px)',
                  minWidth: 300,
                }}
              >
                <div
                  style={{
                    width: '100%',
                    aspectRatio: '16 / 9',
                    background: '#fafafa',
                    borderRadius: 6,
                    overflow: 'hidden',
                    marginBottom: 10,
                  }}
                >
                  <img
                    src={p.image_url || (userId === MOCK_DPC_USER_ID ? MOCK_PRODUCT_IMAGE_URL : DEFAULT_PRODUCT_IMAGE_URL)}
                    alt={p.name}
                    style={{ width: '100%', height: '100%', objectFit: 'cover' }}
                    onError={(e) => {
                      (e.target as HTMLImageElement).src =
                        userId === MOCK_DPC_USER_ID ? MOCK_PRODUCT_IMAGE_URL : DEFAULT_PRODUCT_IMAGE_URL;
                    }}
                  />
                </div>
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
          <p className="muted">
            {activitiesIsMock
              ? '未从后端获取活动，展示本地演示数据（仅演示账号）。'
              : '展示你已发布的秒杀活动。'}
          </p>
        )}
        {activities.length === 0 ? (
          <p className="muted">暂无秒杀活动</p>
        ) : (
          <div style={{ display: 'flex', gap: 12, marginTop: 12, flexWrap: 'wrap' }}>
            {activities.map((a) => (
              <div
                key={a.id}
                style={{
                  border: '1px solid #eee',
                  borderRadius: 8,
                  padding: 12,
                  flex: '1 1 calc(50% - 12px)',
                  minWidth: 320,
                }}
              >
                <div style={{ fontWeight: 700 }}>
                  活动 ID：{a.id} · 商品：{a.product_name}（{a.product_id}）
                </div>
                <div className="muted" style={{ marginTop: 4 }}>
                  秒杀价：{a.seckill_price} · 库存：{a.available_stock} / {a.total_stock} · 状态：{a.status}
                </div>
                <div className="muted" style={{ marginTop: 4 }}>
                  {a.start_time} ~ {a.end_time}
                </div>
                <div style={{ marginTop: 10 }}>
                  <button type="button" className="btn-primary" disabled={!canEdit} onClick={() => startEditActivity(a)}>
                    修改活动
                  </button>
                </div>
              </div>
            ))}
          </div>
        )}
      </section>

      {editingProduct && (
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
          onClick={cancelEdit}
        >
          <div className="card" style={{ width: 'min(720px, 92vw)', padding: 20 }} onClick={(e) => e.stopPropagation()}>
            <h2 style={{ fontSize: '1.1rem', marginTop: 0 }}>修改商品</h2>
            <form onSubmit={onSaveEdit}>
              <table style={{ width: '100%', borderCollapse: 'collapse' }}>
                <tbody>
                  <tr>
                    <td style={{ padding: 8, width: 140 }}>名称</td>
                    <td style={{ padding: 8 }}>
                      <input value={eName} onChange={(e) => setEName(e.target.value)} required />
                    </td>
                  </tr>
                  <tr>
                    <td style={{ padding: 8 }}>描述</td>
                    <td style={{ padding: 8 }}>
                      <textarea value={eDesc} onChange={(e) => setEDesc(e.target.value)} rows={2} />
                    </td>
                  </tr>
                  <tr>
                    <td style={{ padding: 8 }}>价格（分）</td>
                    <td style={{ padding: 8 }}>
                      <input value={ePrice} onChange={(e) => setEPrice(e.target.value)} required />
                    </td>
                  </tr>
                  <tr>
                    <td style={{ padding: 8 }}>库存</td>
                    <td style={{ padding: 8 }}>
                      <input value={eStock} onChange={(e) => setEStock(e.target.value)} required />
                    </td>
                  </tr>
                  <tr>
                    <td style={{ padding: 8 }}>图片 URL（可空）</td>
                    <td style={{ padding: 8 }}>
                      <input value={eImg} onChange={(e) => setEImg(e.target.value)} placeholder="https://..." />
                    </td>
                  </tr>
                  <tr>
                    <td style={{ padding: 8 }}>分类</td>
                    <td style={{ padding: 8 }}>
                      <input value={eCat} onChange={(e) => setECat(e.target.value)} />
                    </td>
                  </tr>
                </tbody>
              </table>
              <div style={{ display: 'flex', gap: 12, marginTop: 12 }}>
                <button type="submit" className="btn-primary" disabled={!canEditProducts}>
                  保存修改
                </button>
                <button type="button" className="btn" onClick={cancelEdit}>
                  关闭
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {editingActivity && (
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
          onClick={cancelEditActivity}
        >
          <div className="card" style={{ width: 'min(720px, 92vw)', padding: 20 }} onClick={(e) => e.stopPropagation()}>
            <h2 style={{ fontSize: '1.1rem', marginTop: 0 }}>修改秒杀活动</h2>
            <form onSubmit={onSaveActivityEdit}>
              <table style={{ width: '100%', borderCollapse: 'collapse' }}>
                <tbody>
                  <tr>
                    <td style={{ padding: 8, width: 140 }}>商品 ID</td>
                    <td style={{ padding: 8 }}>
                      <input value={aProductId} onChange={(e) => setAProductId(e.target.value)} required />
                    </td>
                  </tr>
                  <tr>
                    <td style={{ padding: 8 }}>秒杀价（分）</td>
                    <td style={{ padding: 8 }}>
                      <input value={aPrice} onChange={(e) => setAPrice(e.target.value)} required />
                    </td>
                  </tr>
                  <tr>
                    <td style={{ padding: 8 }}>总库存</td>
                    <td style={{ padding: 8 }}>
                      <input value={aStock} onChange={(e) => setAStock(e.target.value)} required />
                    </td>
                  </tr>
                  <tr>
                    <td style={{ padding: 8 }}>开始时间</td>
                    <td style={{ padding: 8 }}>
                      <input type="datetime-local" value={aStart} onChange={(e) => setAStart(e.target.value)} required />
                    </td>
                  </tr>
                  <tr>
                    <td style={{ padding: 8 }}>结束时间</td>
                    <td style={{ padding: 8 }}>
                      <input type="datetime-local" value={aEnd} onChange={(e) => setAEnd(e.target.value)} required />
                    </td>
                  </tr>
                </tbody>
              </table>
              <div style={{ display: 'flex', gap: 12, marginTop: 12 }}>
                <button type="submit" className="btn-primary" disabled={!canEdit}>
                  保存活动修改
                </button>
                <button type="button" className="btn" onClick={cancelEditActivity}>
                  关闭
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      <p style={{ marginTop: 16 }}>
        <Link to="/">返回首页</Link>
      </p>
    </div>
  );
}
