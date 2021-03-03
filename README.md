# Home Automation Utils

Go rewrite of the [ha-utils](https://github.com/klajbard/ha-utils) package

## ENV variables
- SG_SESSID: Steamgifts session id (`PHPSESSID`)
- HVA_ID: Hardverapro session id (`identifier`)
- NCORE_USERNAME and NCORE_PASSWORD: nCore credentials to obtain session cookie (uses session cookies)
- FIXERAPI: fixer.io API key to query currencies
- HASS_URL: Home Assistant server address
- HASS_TOKEN: Home Assistant access token
- various SLACK_API hook urls from the `https://hooks.slack.com/services/${SLACK_API}` for Slack notifications


## Config

Example config

```yaml
marketplace:
  enabled: true
  jofogas: 
    - name: "ikea+tradfri"
  hardverapro:
    - name: "ikea+tradfri"
habump:
  - name: "hva_item_2"
    id: "93h7z3n4vt7z34ntvi7834z7vtzw374znr73zn475zvn37w4z575vnz2384z5nvi"
```