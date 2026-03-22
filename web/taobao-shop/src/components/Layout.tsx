import { Link, Outlet } from 'react-router-dom';
import { useAuth } from '@/context/AuthContext';

export function Layout() {
  const { username, token, logout } = useAuth();

  return (
    <div style={{ minHeight: '100vh', display: 'flex', flexDirection: 'column' }}>
      <header
        style={{
          background: '#fff',
          borderBottom: '2px solid #ff5000',
          boxShadow: '0 1px 4px rgba(0,0,0,0.06)',
        }}
      >
        <div className="container" style={{ display: 'flex', alignItems: 'center', gap: 16, padding: '10px 12px' }}>
          <Link to="/" style={{ display: 'flex', alignItems: 'baseline', gap: 6 }}>
            <span style={{ fontSize: '1.75rem', fontWeight: 800, color: '#ff5000', letterSpacing: 2 }}>陶宝</span>
            <span className="muted" style={{ fontSize: '0.85rem' }}>
              本地演示
            </span>
          </Link>
          <div style={{ flex: 1 }} />
          <nav style={{ display: 'flex', gap: 16, alignItems: 'center', fontSize: '0.95rem' }}>
            <Link to="/">首页</Link>
            <Link to="/seller">卖家中心</Link>
            {token ? (
              <>
                <Link to="/orders">我的订单</Link>
                <span className="muted">{username || '用户'}</span>
                <button type="button" className="btn-ghost" onClick={() => logout()}>
                  退出
                </button>
              </>
            ) : (
              <>
                <Link to="/login">登录</Link>
                <Link to="/register">
                  <span className="btn-primary" style={{ display: 'inline-block', padding: '6px 14px' }}>
                    免费注册
                  </span>
                </Link>
              </>
            )}
          </nav>
        </div>
        <div
          style={{
            background: 'linear-gradient(90deg, #ff9000 0%, #ff5000 100%)',
            color: '#fff',
            fontSize: '0.85rem',
            padding: '6px 12px',
          }}
        >
          <div className="container" style={{ display: 'flex', gap: 24, flexWrap: 'wrap' }}>
            <span>秒杀热卖</span>
            <span>陶宝直播</span>
            <span>聚划算</span>
            <span>天猫国际</span>
            <span className="muted" style={{ color: 'rgba(255,255,255,0.85)' }}>
              仅为课程/本地联调界面，与任何商业平台无关
            </span>
          </div>
        </div>
      </header>
      <main style={{ flex: 1, padding: '20px 0' }}>
        <Outlet />
      </main>
      <footer style={{ background: '#2c2c2c', color: '#999', padding: '20px', fontSize: '0.85rem', textAlign: 'center' }}>
        <p style={{ margin: 0 }}>陶宝 Demo · Kitex 秒杀电商后端联调前端</p>
      </footer>
    </div>
  );
}
