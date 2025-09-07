# RSS Filter

RSS フィルターサービス - 複数のRSSフィードを集約・フィルタリングするGoアプリケーション

## 概要

このサービスは、設定されたRSSソースからフィードを取得し、指定した条件でフィルタリングして統合されたRSSフィードを生成します。Cloud Runで動作し、HTTPリクエストでRSSフィードを配信します。

## フロー

1. Firestoreから設定を読み込み
2. 各RSSフィードを取得・解析
3. フィルター条件に基づいてコンテンツをフィルタリング
4. 統合されたRSSフィードを生成・配信
5. 配信した記事のURLをFirestoreに保存し、次回以降は重複して配信しない

## GCPセットアップ

本アプリケーションは、設定と配信済みの記事URLを永続化するためにFirestoreを使用します。

### 1. Firestore APIの有効化

Cloud RunをデプロイするGoogle Cloudプロジェクトで、Firestore APIを有効にする必要があります。

```sh
gcloud services enable firestore.googleapis.com
```

### 2. Firestoreデータベースの作成

APIを有効化した後、GCPコンソールからFirestoreデータベースを作成してください。ネイティブモード(Native mode)を選択し、ロケーションを選択します（例: `asia-northeast1`）。

### 3. IAM権限の設定

Cloud Runサービスが使用するサービスアカウントに、Firestoreへの読み書き権限を付与する必要があります。デフォルトでは、Cloud Runは `PROJECT_NUMBER-compute@developer.gserviceaccount.com` というサービスアカウントを使用します。

```sh
# プロジェクトIDとプロジェクト番号を取得
PROJECT_ID=$(gcloud config get-value project)
PROJECT_NUMBER=$(gcloud projects describe $PROJECT_ID --format="value(projectNumber)")

# サービスアカウントにロールを付与
gcloud projects add-iam-policy-binding $PROJECT_ID \
    --member="serviceAccount:$PROJECT_NUMBER-compute@developer.gserviceaccount.com" \
    --role="roles/datastore.user"
```

### 4. 環境変数の設定

Cloud Runサービスに、使用するGoogle CloudプロジェクトIDを環境変数として設定する必要があります。

デプロイコマンドに `--set-env-vars` フラグを追加してください。

```sh
gcloud run deploy rss-filter --source . --platform managed --region us-central1 --allow-unauthenticated --set-env-vars="GCP_PROJECT_ID=$PROJECT_ID"
```

## デプロイメント

### Cloud Runへのデプロイ

```sh
# Cloud Run APIを有効化
gcloud services enable run.googleapis.com

# デプロイ (環境変数の設定を含む)
PROJECT_ID=$(gcloud config get-value project)
gcloud run deploy rss-filter --source . --platform managed --region us-central1 --allow-unauthenticated --set-env-vars="GCP_PROJECT_ID=$PROJECT_ID"
```

### ローカル実行

ローカルで実行する場合は、Google Cloudへの認証と環境変数の設定が必要です。

```sh
# 認証
gcloud auth application-default login

# 環境変数を設定して実行
export GCP_PROJECT_ID=$(gcloud config get-value project)
go run .
```

## API使用方法

### エンドポイント

- **RSSフィード**: `/?name=<config_name>` - 設定された名前のRSSフィードを取得
- **設定UI**: `/config` - 設定を管理するためのWeb UI

## 設定

設定はFirestoreの `rss_configurations` コレクションに保存されます。設定は `/config` エンドポイントのUIから管理できます。

- **コレクション**: `rss_configurations`
- **ドキュメントID**: 設定名 (例: `my_feed`)

### ドキュメントの構造

```json
{
  "name": "my_feed",
  "rsss": [
    {
      "title": "Tech News",
      "url": "https://example.com/tech_news.rss"
    },
    {
      "title": "Go Blog",
      "url": "https://dev.to/feed/tag/go"
    }
  ],
  "ignores": [
    {
      "condition": 1, // TITLE_INCLUDES
      "text": "[PR]"
    },
    {
      "condition": 5, // URL_INCLUDES
      "text": "/sponsored/"
    }
  ]
}
```

### `condition` の値

| 値 | 説明 |
| :--- | :--- |
| 1 | タイトルに含まれる (TITLE_INCLUDES) |
| 2 | タイトルまたは概要に含まれる (TITLE_OR_DESC_INCLUDES) |
| 3 | タイトルの先頭 (TITLE_BEGINS) |
| 4 | タイトルの終わり (TITLE_ENDS) |
| 5 | URLに含まれる (URL_INCLUDES) |
| 6 | 著者に含まれる (AUTHOR_INCLUDES) |

## 技術仕様

- **言語**: Go
- **プラットフォーム**: Google Cloud Run
- **データベース**: Google Cloud Firestore
- **コンテナ化**: Docker
- **依存関係**:
  - `github.com/gorilla/feeds` - RSS生成
  - `github.com/mmcdole/gofeed` - RSS解析
  - `cloud.google.com/go/firestore` - データベースアクセス
  - `github.com/gorilla/mux` - HTTPルーティング
