"""Emit tiktok-api-seckill.import.json from tiktok-api-seckill.json for Grafana UI import."""
import json
from pathlib import Path

root = Path(__file__).resolve().parents[1]
src = root / "config/grafana/dashboards/tiktok-api-seckill.json"
dst = root / "config/grafana/dashboards/tiktok-api-seckill.import.json"

raw = json.loads(src.read_text(encoding="utf-8"))
inputs = [
    {
        "name": "DS_PROMETHEUS",
        "label": "Prometheus",
        "description": "指向已抓取 API /metrics 的 Prometheus",
        "type": "datasource",
        "pluginId": "prometheus",
        "pluginName": "Prometheus",
    }
]
d = {"__inputs": inputs, **raw}


def swap_ds_uid(obj):
    if isinstance(obj, dict):
        if obj.get("type") == "prometheus" and obj.get("uid") == "prometheus":
            obj["uid"] = "${DS_PROMETHEUS}"
        for v in obj.values():
            swap_ds_uid(v)
    elif isinstance(obj, list):
        for item in obj:
            swap_ds_uid(item)


swap_ds_uid(d)
dst.write_text(json.dumps(d, ensure_ascii=False, indent=2), encoding="utf-8")
print("Wrote", dst)
