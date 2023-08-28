resource "timeplus_dashboard" "example" {
  name        = "example"
  description = "A dashboard example with control, markdown and chart panels."
  panels      = <<JSON
[
  {
    "id": "d421e9bd-1b63-4862-9bcb-e9bd7a2b594b",
    "title": "Title",
    "description": "",
    "position": {
      "h": 1,
      "nextX": 12,
      "nextY": 2,
      "w": 12,
      "x": 0,
      "y": 1
    },
    "viz_type": "markdown",
    "viz_content": "",
    "viz_config": {
      "mdString": "Error Monitoring"
    }
  },
  {
    "id": "310bf39f-b218-423d-bf4d-41794d63715a",
    "title": "App Name",
    "description": "",
    "position": {
      "h": 1,
      "nextX": 3,
      "nextY": 1,
      "w": 3,
      "x": 0,
      "y": 0
    },
    "viz_type": "control",
    "viz_content": "",
    "viz_config": {
      "chartType": "text",
      "defaultValue": "my_app",
      "inlineValues": "",
      "label": "App Name",
      "target": "app_name"
    }
  },
  {
    "id": "a71775b7-0444-4074-a8ad-4dc5ee05d43f",
    "title": "Error count in last 30 min",
    "description": "",
    "position": {
      "h": 4,
      "nextX": 6,
      "nextY": 6,
      "w": 6,
      "x": 0,
      "y": 2
    },
    "viz_type": "chart",
    "viz_content": "select window_start, window_end, sum(value) from tumble(metrics, 30m) where name = 'error' group by window_start, window_end",
    "viz_config": {
      "chartType": "line"
    }
  }
]
JSON
}
