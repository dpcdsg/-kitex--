import http from "k6/http";
import { check, group, sleep } from "k6";

/**
 * 全量「已注册 HTTP 路由」巡检：每个接口在每次迭代内至少请求 1 次（HTTP 层）。
 * 不含：代码里有但路由未挂载的接口（如 /seckill/product/publish/）。
 * 含：GET /metrics（Prometheus 抓取端点）。
 *
 * 业务上许多 POST 会失败（权限、库存、状态机），但本项目错误多为 HTTP 200 + JSON errno；
 * 本脚本默认只校验 HTTP 状态码，用于压测网关、RPC、DB 链路与 Sentinel。
 */
const BASE_URL = __ENV.BASE_URL || "http://127.0.0.1:10001";
const USERNAME = __ENV.USERNAME || "dpc";
const PASSWORD = __ENV.PASSWORD || "123";

const jsonHeaders = { headers: { "Content-Type": "application/json" } };

export const options = {
  vus: Number(__ENV.VUS || 5),
  duration: __ENV.DURATION || "60s",
  thresholds: {
    http_req_failed: ["rate<0.3"],
    checks: ["rate>0.85"],
  },
};

function j(res) {
  try {
    return res.json();
  } catch (e) {
    return null;
  }
}

function firstProductIdFromList(res) {
  const b = j(res);
  const list = b && b.product_list;
  if (list && list.length && list[0].id) return String(list[0].id);
  return "";
}

function firstActivityIdFromList(res) {
  const b = j(res);
  const list = b && b.activity_list;
  if (list && list.length && list[0].id) return String(list[0].id);
  return "";
}

function firstOrderNoFromList(res) {
  const b = j(res);
  const list = b && b.order_list;
  if (list && list.length && list[0].order_no) return String(list[0].order_no);
  return "";
}

function productIdFromCreateBody(res) {
  const b = j(res);
  if (b && b.product && b.product.id) return String(b.product.id);
  return "";
}

function activityIdFromCreateBody(res) {
  const b = j(res);
  if (b && b.activity && b.activity.id) return String(b.activity.id);
  return "";
}

function orderNoFromBuyBody(res) {
  const b = j(res);
  if (b && b.order_no) return String(b.order_no);
  return "";
}

export function setup() {
  http.post(
    `${BASE_URL}/seckill/user/register/`,
    JSON.stringify({ username: USERNAME, password: PASSWORD }),
    jsonHeaders
  );

  const loginRes = http.post(
    `${BASE_URL}/seckill/user/login/`,
    JSON.stringify({ username: USERNAME, password: PASSWORD }),
    jsonHeaders
  );
  const lb = j(loginRes);
  if (!lb || Number(lb.status_code) !== 0 || !lb.token) {
    throw new Error("setup: login failed (need valid USERNAME/PASSWORD and user service)");
  }

  const token = lb.token;
  const userId = String(lb.user_id);
  const t = encodeURIComponent(token);

  const pl = http.get(
    `${BASE_URL}/seckill/product/list/?token=${t}&page=1&size=10`
  );
  const al = http.get(
    `${BASE_URL}/seckill/activity/list/?page=1&size=10&token=${t}`
  );
  const ol = http.get(
    `${BASE_URL}/seckill/order/list/?token=${t}&page=1&size=10`
  );

  const baseProductId =
    __ENV.PRODUCT_ID || firstProductIdFromList(pl) || "1";
  const baseActivityId =
    __ENV.ACTIVITY_ID || firstActivityIdFromList(al) || "1000001";
  const baseOrderNo = __ENV.ORDER_NO || firstOrderNoFromList(ol) || "";

  return { token, userId, baseProductId, baseActivityId, baseOrderNo };
}

function ok200(res) {
  return res.status === 200;
}

export default function (data) {
  const { token, userId, baseProductId, baseActivityId, baseOrderNo } = data;
  const t = encodeURIComponent(token);
  const uniq = `${__VU}-${__ITER}-${Date.now()}`;

  group("user", () => {
    let r = http.post(
      `${BASE_URL}/seckill/user/register/`,
      JSON.stringify({ username: USERNAME, password: PASSWORD }),
      jsonHeaders
    );
    check(r, { "POST /seckill/user/register/": ok200 });

    r = http.post(
      `${BASE_URL}/seckill/user/login/`,
      JSON.stringify({ username: USERNAME, password: PASSWORD }),
      jsonHeaders
    );
    check(r, { "POST /seckill/user/login/": ok200 });

    r = http.get(
      `${BASE_URL}/seckill/user/?user_id=${userId}&token=${t}`
    );
    check(r, { "GET /seckill/user/": ok200 });
  });

  sleep(0.02);

  group("product_read", () => {
    let r = http.get(
      `${BASE_URL}/seckill/product/list/?token=${t}&page=1&size=10`
    );
    check(r, { "GET /seckill/product/list/": ok200 });

    r = http.get(
      `${BASE_URL}/seckill/product/detail/?token=${t}&product_id=${baseProductId}`
    );
    check(r, { "GET /seckill/product/detail/": ok200 });
  });

  let newProductId = "";
  group("product_write", () => {
    const r = http.post(
      `${BASE_URL}/seckill/product/create/`,
      JSON.stringify({
        token: token,
        name: `k6_${uniq}`,
        description: "loadtest all_api_k6",
        price: 100,
        stock: 500,
        image_url: "https://example.com/p.png",
        category: "test",
      }),
      jsonHeaders
    );
    check(r, { "POST /seckill/product/create/": ok200 });
    newProductId = productIdFromCreateBody(r) || baseProductId;
  });

  sleep(0.02);

  group("activity_read", () => {
    let r = http.get(
      `${BASE_URL}/seckill/activity/list/?page=1&size=10&token=${t}`
    );
    check(r, { "GET /seckill/activity/list/": ok200 });

    r = http.get(
      `${BASE_URL}/seckill/activity/detail/?token=${t}&activity_id=${baseActivityId}`
    );
    check(r, { "GET /seckill/activity/detail/": ok200 });
  });

  let newActivityId = "";
  group("activity_write", () => {
    const cr = http.post(
      `${BASE_URL}/seckill/activity/create/`,
      JSON.stringify({
        token: token,
        product_id: newProductId,
        seckill_price: 1,
        total_stock: 50,
        start_time: "2026-01-01 00:00:00",
        end_time: "2030-12-31 23:59:59",
      }),
      jsonHeaders
    );
    check(cr, { "POST /seckill/activity/create/": ok200 });
    newActivityId = activityIdFromCreateBody(cr) || baseActivityId;

    const ur = http.post(
      `${BASE_URL}/seckill/activity/update/`,
      JSON.stringify({
        token: token,
        activity_id: newActivityId,
        product_id: newProductId,
        seckill_price: 2,
        total_stock: 48,
        start_time: "2026-01-01 00:00:00",
        end_time: "2030-12-31 23:59:59",
      }),
      jsonHeaders
    );
    check(ur, { "POST /seckill/activity/update/": ok200 });
  });

  sleep(0.02);

  group("seckill_action", () => {
    const r = http.post(
      `${BASE_URL}/seckill/action/`,
      JSON.stringify({
        token: token,
        activity_id: newActivityId,
      }),
      jsonHeaders
    );
    check(r, { "POST /seckill/action/": ok200 });
  });

  sleep(0.02);

  let orderNoPay = "";
  let orderNoCancel = "";
  group("order", () => {
    let r = http.get(
      `${BASE_URL}/seckill/order/list/?token=${t}&page=1&size=10`
    );
    check(r, { "GET /seckill/order/list/": ok200 });

    const detailNo = baseOrderNo || "0";
    r = http.get(
      `${BASE_URL}/seckill/order/detail/?token=${t}&order_no=${encodeURIComponent(detailNo)}`
    );
    check(r, { "GET /seckill/order/detail/": ok200 });

    const b1 = http.post(
      `${BASE_URL}/seckill/order/buy/`,
      JSON.stringify({ token: token, product_id: newProductId }),
      jsonHeaders
    );
    check(b1, { "POST /seckill/order/buy/ (cancel path)": ok200 });
    orderNoCancel = orderNoFromBuyBody(b1);

    if (orderNoCancel) {
      r = http.post(
        `${BASE_URL}/seckill/order/cancel/`,
        JSON.stringify({ token: token, order_no: orderNoCancel }),
        jsonHeaders
      );
      check(r, { "POST /seckill/order/cancel/": ok200 });
    } else {
      check(true, { "POST /seckill/order/cancel/ (skipped: no order_no)": () => true });
    }

    const b2 = http.post(
      `${BASE_URL}/seckill/order/buy/`,
      JSON.stringify({ token: token, product_id: newProductId }),
      jsonHeaders
    );
    check(b2, { "POST /seckill/order/buy/ (pay path)": ok200 });
    orderNoPay = orderNoFromBuyBody(b2);

    if (orderNoPay) {
      r = http.post(
        `${BASE_URL}/seckill/order/pay/`,
        JSON.stringify({ token: token, order_no: orderNoPay }),
        jsonHeaders
      );
      check(r, { "POST /seckill/order/pay/": ok200 });
    } else {
      check(true, { "POST /seckill/order/pay/ (skipped: no order_no)": () => true });
    }
  });

  sleep(0.02);

  group("product_update_then_activity_delete_then_product_delete", () => {
    const r = http.post(
      `${BASE_URL}/seckill/product/update/`,
      JSON.stringify({
        token: token,
        product_id: newProductId,
        name: `k6_${uniq}_u`,
        description: "updated",
        price: 101,
        stock: 400,
        image_url: "https://example.com/p2.png",
        category: "test",
      }),
      jsonHeaders
    );
    check(r, { "POST /seckill/product/update/": ok200 });

    const ad = http.post(
      `${BASE_URL}/seckill/activity/delete/`,
      JSON.stringify({ token: token, activity_id: newActivityId }),
      jsonHeaders
    );
    check(ad, { "POST /seckill/activity/delete/": ok200 });

    const pd = http.post(
      `${BASE_URL}/seckill/product/delete/`,
      JSON.stringify({ token: token, product_id: newProductId }),
      jsonHeaders
    );
    check(pd, { "POST /seckill/product/delete/": ok200 });
  });

  sleep(0.02);

  group("metrics", () => {
    const r = http.get(`${BASE_URL}/metrics`);
    check(r, { "GET /metrics": ok200 });
  });
}
