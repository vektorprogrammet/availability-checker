## Availability Checker
I `config.json` kan man definere en liste med websites:
```json
"websites": [
    "https://example.com",
    "https://google.com"
  ]
```

Disse vil bli pinget hvert minutt. Hvis Availability Checker ikke får `200 OK` tilbake som svar vil det sendes en melding på slack til alle kanalene som er spesifisert i konfigen:

```json
"slackWorkspaces": [
    {
      "endpoint": "https://hooks.slack.com/services/XXXXX/XXXXX/XXXXX",
      "channel": "#log_channel",
      "username": "bot_name",
      "icon_emoji": ":robot_face:"
    },
    {
      "endpoint": "https://hooks.slack.com/services/YYYYY/YYYYY/YYYYY",
      "channel": "#general",
      "username": "bot_name",
      "icon_emoji": ":robot_face:"
    }
],
```
