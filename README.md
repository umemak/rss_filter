# RSS Filter

RSS フィルターサービス - 複数のRSSフィードを集約・フィルタリングするGoアプリケーション

## 概要

このサービスは、設定されたRSSソースからフィードを取得し、指定した条件でフィルタリングして統合されたRSSフィードを生成します。Cloud Runで動作し、HTTPリクエストでRSSフィードを配信します。

## フロー

1. 設定ファイルからRSSソースと フィルター条件を読み込み
2. 各RSSフィードを取得・解析
3. フィルター条件に基づいてコンテンツをフィルタリング
4. 統合されたRSSフィードを生成・配信
5. (New) 配信した記事のURLをFirestoreに保存し、次回以降は重複して配信しない

## GCPセットアップ

本アプリケーションは、配信済みの記事URLを永続化するためにFirestoreを使用します。

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
