import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { userRegister, userInfo } from '@/api/seckill';
import { useAuth } from '@/context/AuthContext';

export function RegisterPage() {
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
    try {
      const r = await userRegister(username, password);
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
      setErr(e instanceof Error ? e.message : '注册失败');
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="container" style={{ maxWidth: 480 }}>
      <h1 style={{ fontSize: '1.5rem' }}>注册陶宝账号</h1>
      <p className="muted">连接本机 Kitex API（/seckill/user/register/）</p>
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
            autoComplete="new-password"
            required
          />
        </div>
        {err && <p className="err">{err}</p>}
        <button type="submit" className="btn-primary" disabled={loading}>
          {loading ? '提交中…' : '同意协议并注册'}
        </button>
      </form>
      <p style={{ marginTop: 16 }}>
        已有账号？ <Link to="/login">登录</Link>
      </p>
    </div>
  );
}
