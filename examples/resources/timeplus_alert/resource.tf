resource "timeplus_stream" "machine" {
  name        = "cpu_temperature"
  description = "A stream for alert example."
  column {
    name = "machine"
    type = "string"
  }
  column {
    name = "temp"
    type = "uint8"
  }
}

resource "timeplus_alert" "example" {
  name        = "cpu_too_hot"
  description = "Alarm for when the CPU gets too hot"
  severity    = 1
  action      = "slack"
  properties = jsonencode({
    url              = "https://hooks.slack.com/services/my/slack/web_hook"
    trigger_template = "Machine {{ .machine }} is getting too hot: {{ .temp }}"
    resolve_template = "Machine {{ .machine }} is cooling down: {{ .temp }}"
  })
  trigger_sql = "select machine, temp from ${resource.timeplus_stream.machine.name} where temp >= 98"
  resolve_sql = "select machine, temp from ${resource.timeplus_stream.machine.name} where temp < 95"
}
