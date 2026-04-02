import http from "k6/http";
import { check, sleep } from "k6";

const BASE_URL = __ENV.BASE_URL || "http://127.0.0.1:8080";
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
    token = body.token || "";
  } catch (e) {
    token = "";
  }

  if (!token) {
    throw new Error("login failed: token is empty");
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
