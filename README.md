# grafana-presto-datasource
Presto datasource plugin for Grafana. Support below data types(other data types will response as strings): Boolean, Integer, Floating-Point, Fixed-Precision, String, Date and Time.


# How to build
1. make dep
2. make build-js
3. make build-go
4. make package

# How to install plugin
1. Copy the plugin to Grafana plugin folder.
2. Add plugin name `grafana-presto-datasource` to field `allow_loading_unsigned_plugins` in config.