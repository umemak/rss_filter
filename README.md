# RSS Filter

RSS フィルターサービス - 複数のRSSフィードを集約・フィルタリングするGoアプリケーション

## 概要

このサービスは、設定されたRSSソースからフィードを取得し、指定した条件でフィルタリングして統合されたRSSフィードを生成します。Cloud Runで動作し、HTTPリクエストでRSSフィードを配信します。

## フロー

1. 設定ファイルからRSSソースと フィルター条件を読み込み
2. 各RSSフィードを取得・解析
3. フィルター条件に基づいてコンテンツをフィルタリング
4. 統合されたRSSフィードを生成・配信

## デプロイメント

### Cloud Runへのデプロイ

```sh
# Cloud Run APIを有効化
gcloud services enable run.googleapis.com

# デプロイ
gcloud run deploy rss-filter --source . --platform managed --region us-central1 --allow-unauthenticated
```

### ローカル実行

```sh
go run .
```

## API使用方法

### エンドポイント

- **設定名指定**: `?name=<config_name>` - 設定された名前のRSSフィードを取得
- **URL指定**: `?url=<rss_url>` - 指定したURLのRSSフィードをフィルタリングして取得

## 設定

- `config.go` - RSSソースとフィルター条件を定義
- `config.go.sample` - 設定ファイルのサンプル
- 各種フィルター配列(`filter1`〜`filter6`)でフィルタリング戦略を定義

## 技術仕様

- **言語**: Go 1.25.0
- **プラットフォーム**: Google Cloud Run
- **コンテナ化**: Docker
- **依存関係**:
  - `github.com/gorilla/feeds` - RSS生成
  - `github.com/mmcdole/gofeed` - RSS解析
