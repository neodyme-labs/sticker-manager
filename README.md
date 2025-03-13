# Sticker-Manager

This application allows you to configure a custom sticker ednpoint for multiple users by implementing a very small integration manager.
The alternative would be to use https://github.com/turt2live/matrix-dimension which is deprecated and a huge service.

To actually configure the sticker server use this https://github.com/maunium/stickerpicker

After configuring the service you can configure this as your integration manager either inside your homeserver or the Element installation using:
```
integrations_ui_url=https://example.com/
integrations_rest_url=https://example.com/api
```

# Development

Add this to your account data inside the `m.widget` state using the developer tools
```json
{
  "integration-manager": {
    "content": {
      "type": "m.integration_manager",
      "url": "http://127.0.0.1:8080/",
      "data": {
        "api_url": "http://127.0.0.1:8080/api"
      }
    },
    "type": "m.widget"
  }
}
```