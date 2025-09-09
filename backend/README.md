```
backend/
├── cmd/
│   └── server/
│       └── main.go        # エントリーポイント
├── config/                # 設定ファイル（DB, 環境変数など）
│   └── config.go
├── internal/
│   ├── handler/           # HTTPハンドラ (コントローラ)
│   │   └── user_handler.go
│   ├── service/           # ビジネスロジック
│   │   └── user_service.go
│   ├── repository/        # DBアクセス (SQL/ORM)
│   │   └── user_repository.go
│   ├── model/             # モデル（構造体定義）
│   │   └── user.go
│   └── router/            # ルーティング設定
│       └── router.go
├── pkg/                   # 再利用可能なユーティリティ (JWT, Logger, etc.)
│   └── logger/logger.go
├── go.mod
└── go.sum
```