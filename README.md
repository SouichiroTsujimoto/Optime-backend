# Yotei Backend API

## セットアップ

### 前提条件

- Go 1.24以上
- PostgreSQL

### ローカル開発

1. リポジトリをクローン

```bash
git clone <repository-url>
cd yotei-backend
```

2. 依存パッケージをインストール

```bash
go mod download
```

3. 環境変数を設定

`.env.example`をコピーして`.env`を作成し、適切な値を設定してください。

```bash
cp .env.example .env
```

4. PostgreSQLデータベースを作成

```bash
createdb yotei_db
```

5. サーバーを起動

```bash
go run main.go
```

サーバーは `http://localhost:3000` で起動します。

## APIエンドポイント

### イベント作成

**POST** `/api/v1/events`

イベントと候補日を作成します。

リクエスト例：
```json
{
  "title": "新年会",
  "candidate_dates": [
    "2025-01-10T19:00:00Z",
    "2025-01-11T19:00:00Z",
    "2025-01-12T19:00:00Z"
  ]
}
```

レスポンス例：
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "title": "新年会",
  "candidate_dates": [
    {
      "id": 1,
      "date_time": "2025-01-10T19:00:00Z"
    },
    {
      "id": 2,
      "date_time": "2025-01-11T19:00:00Z"
    },
    {
      "id": 3,
      "date_time": "2025-01-12T19:00:00Z"
    }
  ],
  "created_at": "2024-12-01T10:00:00Z"
}
```

### イベント詳細取得

**GET** `/api/v1/events/:id`

特定のイベントの詳細を取得します。

## Renderへのデプロイ

1. Renderのダッシュボードで「New +」→「Blueprint」を選択
2. このリポジトリを接続
3. `render.yaml`が自動的に検出され、PostgreSQLデータベースとWebサービスが作成されます

または、手動でデプロイする場合：

1. Renderでデータベースを作成（PostgreSQL）
2. Renderで新しいWebサービスを作成
3. 環境変数`DATABASE_URL`にデータベース接続URLを設定
4. ビルドコマンド: `go build -o main .`
5. 起動コマンド: `./main`

## 開発

### データベーススキーマ

- **Event**: イベント情報
- **CandidateDate**: 候補日
- **Participant**: 参加者
- **Response**: 参加者の回答

### 新しいエンドポイントの追加

1. `handlers/`ディレクトリにハンドラー関数を作成
2. `main.go`でルートを追加

## ライセンス

MIT

