# Monitoring

## Metrics

Prometheus metrics are served on `localhost:2115/metrics`. Backend metrics include websocket connection counts, cached websocket player counts, websocket message send/receive counts, websocket send and receive errors by type and category, matchmaking queue depth, active matches, move and draw event totals, and match results by outcome and reason. The `grafana-dashboard.json` file contains a Grafana dashboard configuration that can be imported to visualize these metrics.

## Logs

The backend writes logs to stdout. The `config.alloy` file contains the required blocks that should be added to the Grafana Alloy configuration to collect these logs with Loki. The [Grafana docs](https://grafana.com/docs/alloy/latest/monitor/monitor-docker-containers/) contain more details on how to set up log collection for Docker containers. After applying the configuration changes, you need to restart the Grafana Alloy container for the changes to take effect. Finally, the Grafana Alloy needs to get permission to access Docker. You can do this by running the following command, followed by a restart of Grafana Alloy:

```bash
sudo usermod -aG docker grafana-alloy
```

The `grafana-dashboard.json` file also contains a panel that visualizes the logs.
