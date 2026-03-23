import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { userLogin, userInfo } from '@/api/seckill';
import { useAuth } from '@/context/AuthContext';

export function LoginPage() {
  const nav = useNavigate();
  const { setSession } = useAuth();
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [err, setErr] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  async function onSubmit(e: React.FormEvent) {
    e.preventDefault();
    setErr(null);
    setLoading(true);

    // 前端联调 mock：用于在后端不可用时验证登录态页面展示效果。
    // 约定：用户名 dpc，密码 123 则直接写入本地登录态（token 为任意非空字符串）。
    if (username === 'dpc' && password === '123') {
      setSession(10001, 'mock-token-dpc-123', 'dpc');
      nav('/seller');
      setLoading(false);
      return;
    }

    try {
      const r = await userLogin(username, password);
      let displayName: string | null = null;
      try {
        const info = await userInfo(r.user_id, r.token);
        displayName = info.user?.name ?? null;
      } catch {
        /* ignore */
      }
      setSession(r.user_id, r.token, displayName);
      nav('/');
    } catch (e) {
      setErr(e instanceof Error ? e.message : '登录失败');
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="container" style={{ maxWidth: 480 }}>
      <h1 style={{ fontSize: '1.5rem' }}>登录陶宝</h1>
      <p className="muted">连接本机 Kitex API（/seckill/user/login/）</p>
      <form onSubmit={onSubmit} style={{ marginTop: 20 }}>
        <div className="form-row">
          <label>用户名</label>
          <input value={username} onChange={(e) => setUsername(e.target.value)} autoComplete="username" required />
        </div>
        <div className="form-row">
          <label>密码</label>
          <input
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            autoComplete="current-password"
            required
          />
        </div>
        {err && <p className="err">{err}</p>}
        <button type="submit" className="btn-primary" disabled={loading}>
          {loading ? '登录中…' : '登录'}
        </button>
      </form>
      <p style={{ marginTop: 16 }}>
        没有账号？ <Link to="/register">免费注册</Link>
      </p>
    </div>
  );
}
