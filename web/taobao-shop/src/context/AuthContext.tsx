import {
  createContext,
  useCallback,
  useContext,
  useMemo,
  useState,
  type ReactNode,
} from 'react';
import { userInfo as fetchUserInfo } from '@/api/seckill';

type AuthState = {
  userId: number | null;
  token: string | null;
  username: string | null;
};

const STORAGE_KEY = 'taobao_demo_auth';

function load(): AuthState {
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    if (!raw) return { userId: null, token: null, username: null };
    const j = JSON.parse(raw) as AuthState;
    return {
      userId: j.userId ?? null,
      token: j.token ?? null,
      username: j.username ?? null,
    };
  } catch {
    return { userId: null, token: null, username: null };
  }
}

type Ctx = AuthState & {
  setSession: (userId: number, token: string, username?: string | null) => void;
  logout: () => void;
  refreshProfile: () => Promise<void>;
};

const AuthContext = createContext<Ctx | null>(null);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [state, setState] = useState<AuthState>(() =>
    typeof window !== 'undefined' ? load() : { userId: null, token: null, username: null },
  );

  const persist = useCallback((s: AuthState) => {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(s));
    setState(s);
  }, []);

  const setSession = useCallback(
    (userId: number, token: string, username?: string | null) => {
      persist({ userId, token, username: username ?? null });
    },
    [persist],
  );

  const logout = useCallback(() => {
    localStorage.removeItem(STORAGE_KEY);
    setState({ userId: null, token: null, username: null });
  }, []);

  const refreshProfile = useCallback(async () => {
    if (!state.token || !state.userId) return;
    const r = await fetchUserInfo(state.userId, state.token);
    if (r.user?.name) {
      persist({
        userId: state.userId,
        token: state.token,
        username: r.user.name,
      });
    }
  }, [state.token, state.userId, persist]);

  const value = useMemo(
    () => ({
      ...state,
      setSession,
      logout,
      refreshProfile,
    }),
    [state, setSession, logout, refreshProfile],
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const c = useContext(AuthContext);
  if (!c) throw new Error('useAuth outside AuthProvider');
  return c;
}
