resource "timeplus_materialized_view" "example" {
  name                = "example"
  description         = "this is a view example"
  query               = "select max(value) from tumble(foo, 5m)"
  retention_size      = 10 * 1024 * 1204 * 1024 // 10Gi in bytes
  retention_period    = 7 * 24 * 60 * 60 * 1000 // 7 days in ms 
  historical_data_ttl = "to_datetime(_tp_time) + INTERVAL 30 DAYS"
}
