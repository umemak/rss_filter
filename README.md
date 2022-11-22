# rss_filter

## flow

- Read config yaml
- Fetch RSS
- Filter
- Write result RSS

## deploy

```sh
gcloud functions deploy Handler --runtime go116 --trigger-http --allow-unauthenticated
```
