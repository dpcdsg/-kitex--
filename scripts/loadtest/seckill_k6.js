import http from "k6/http";
import { check, sleep } from "k6";

// 全栈 Docker（docker-compose.apps.yml）下 API 映射为宿主机 10001:10001
const BASE_URL = __ENV.BASE_URL || "http://127.0.0.1:10001";
const USERNAME = __ENV.USERNAME || "dpc";
const PASSWORD = __ENV.PASSWORD || "123";
const ACTIVITY_ID = __ENV.ACTIVITY_ID || "1000001";

export const options = {
  vus: Number(__ENV.VUS || 50),
  duration: __ENV.DURATION || "30s",
  thresholds: {
    http_req_failed: ["rate<0.2"],
    http_req_duration: ["p(95)<800"],
    checks: ["rate>0.8"],
  },
};

export function setup() {
  // 用户可能尚未存在：先注册（已存在时接口返回业务错误，忽略即可）
  http.post(
    `${BASE_URL}/seckill/user/register/`,
    JSON.stringify({ username: USERNAME, password: PASSWORD }),
    { headers: { "Content-Type": "application/json" } }
  );

  const loginResp = http.post(
    `${BASE_URL}/seckill/user/login/`,
    JSON.stringify({
      username: USERNAME,
      password: PASSWORD,
    }),
    {
      headers: { "Content-Type": "application/json" },
    }
  );

  check(loginResp, {
    "login http 200": (r) => r.status === 200,
  });

  let token = "";
  try {
    const body = loginResp.json();
    if (Number(body.status_code) === 0 && body.token) {
      token = body.token;
    }
  } catch (e) {
    token = "";
  }

  if (!token) {
    throw new Error(
      "login failed: no token (check USERNAME/PASSWORD and that user service is up)"
    );
  }

  return { token };
}

export default function (data) {
  const listResp = http.get(`${BASE_URL}/seckill/activity/list/?page=1&size=10&status=-1`);
  check(listResp, {
    "list activity ok": (r) => r.status === 200,
  });

  const seckillResp = http.post(
    `${BASE_URL}/seckill/action/`,
    JSON.stringify({
      token: data.token,
      activity_id: ACTIVITY_ID,
    }),
    {
      headers: { "Content-Type": "application/json" },
    }
  );

  check(seckillResp, {
    "seckill request returned": (r) => r.status === 200,
  });

  sleep(0.1);
}
