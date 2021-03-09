# Home Automation Utils

Go rewrite of the [ha-utils](https://github.com/klajbard/ha-utils) package

## ENV variables
- SG_SESSID: Steamgifts session id (`PHPSESSID`)
- HVA_ID: Hardverapro session id (`identifier`)
- NCORE_USERNAME and NCORE_PASSWORD: nCore credentials to obtain session cookie (uses session cookies)
- FIXERAPI: fixer.io API key to query currencies
- HASS_URL: Home Assistant server address
- HASS_TOKEN: Home Assistant access token
- SLACK_APP_TOKEN (`xapp-xxx`) and SLACK_BOT_TOKEN (`xoxb-xxx`): Slack notifications


## Config

Example config

```yaml
marketplace:
  jofogas: 
    - name: "ikea+tradfri"
  hardverapro:
    - name: "ikea+tradfri"
habump:
  - identifier: "ha643a7872ba113a2dc8c1433ad968cd55"
    items:
      - name: "hva_item_2"
        id: "93h7z3n4vt7z34ntvi7834z7vtzw374znr73zn475zvn37w4z575vnz2384z5nvi"
      - name: "hva_item_5"
        id: "c4t3bt467tr346f76364bf34v8f34vf873b4fh873h48fv348nfh83h48fh3847h"
channels:
  - name: "channel1"
    id: "C01MV9VRUCC"
  - name: "channel2"
    id: "C016R5Q1YTW"
enable:
  bestbuy: false
  stockwatcher: false
  marketplace: false
  steamgifts: false
  dht: false
  arukereso: false
  covid: false
  bumphva: false
  ncore: false
  fuel: false
  fixerio: false
  awscost: false
```