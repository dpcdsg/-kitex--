/**
 * 后端 Hertz JWT 中间件只从 query 或 x-www-form-urlencoded 取 token，
 * JSON body 里的 token 不会被识别，因此所有需登录的请求都把 token 挂在 URL 上。
 */
const API_PREFIX = '/seckill';

export type ApiError = { status_code: number; status_msg?: string };

function isApiError(j: unknown): j is ApiError {
  return (
    typeof j === 'object' &&
    j !== null &&
    'status_code' in j &&
    typeof (j as ApiError).status_code === 'number' &&
    (j as ApiError).status_code === -1
  );
}

export function withToken(pathWithQuery: string, token: string | null): string {
  if (!token) return pathWithQuery;
  const sep = pathWithQuery.includes('?') ? '&' : '?';
  return `${pathWithQuery}${sep}token=${encodeURIComponent(token)}`;
}

async function parseJson<T>(res: Response): Promise<T> {
  const text = await res.text();
  let data: unknown = {};
  try {
    data = text ? JSON.parse(text) : {};
  } catch {
    throw new Error(`无效 JSON 响应: ${text.slice(0, 200)}`);
  }
  return data as T;
}

export async function apiGet<T>(path: string, token?: string | null): Promise<T> {
  const url = withToken(`${API_PREFIX}${path}`, token ?? null);
  const res = await fetch(url, { credentials: 'omit' });
  const data = await parseJson<T & ApiError>(res);
  if (isApiError(data)) {
    throw new Error(data.status_msg || '请求失败');
  }
  return data;
}

export async function apiPostJson<T>(
  path: string,
  body: Record<string, unknown>,
  token?: string | null,
): Promise<T> {
  const url = withToken(`${API_PREFIX}${path}`, token ?? null);
  const res = await fetch(url, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  });
  const data = await parseJson<T & ApiError>(res);
  if (isApiError(data)) {
    throw new Error(data.status_msg || '请求失败');
  }
  return data;
}
