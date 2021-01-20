aria2 Exporter for Prometheus
=============================

This exporter exports statistics metrics about torrents downloading and seeding
in [aria2](https://aria2.github.io/) via the [rpc
interface](https://aria2.github.io/manual/en/html/aria2c.html#rpc-interface).

Usage
-----

```
docker run -d --name aria2_exporter -p 9578:9578 \
  -e ARIA2_URL=http://aria2.example.com:6800 \
  -e ARIA2_RPC_SECRET=aria2-rpc-secret-token \
  sbruder/aria2_exporter
```

Replace `aria2.example.com:6800` with the host and port of your aria2 instance.

Replace `aria2-rpc-secret-token` with the RPC secret authorization token if aria2
is configured to use it or leave it blank otherwise.

Metrics are available on http://localhost:9578/metrics or on the endpoint set
in `ARIA2_EXPORTER_LISTEN_ADDRESS`.
