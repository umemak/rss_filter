# rss_filter

## flow

- Read config yaml
- Fetch RSS
  - [mmcdole/gofeed: Parse RSS, Atom and JSON feeds in Go](https://github.com/mmcdole/gofeed)
- Filter
- Write result RSS

## config yaml

```yaml
- sites:
  - url: https://rss.example.com
  - filters:
    - target: title
    - condition: include
    - text: hoge
```
