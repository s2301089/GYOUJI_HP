# GYOUJI_HP
仙台高専広瀬行事委員会広報用HP

# アプリ基本設計書
## 1. アプリ概要
### 1.1. アプリ名
行事委員会紹介＋秋季スポーツ大会トーナメント結果速報アプリ

### 1.2. 目的
学校のスポーツ大会において、学生がリアルタイムで試合結果を把握できるようにすることで、大会参加への関心と一体感を高める。

### 1.3. ターゲットユーザー
管理者 (superroot, admin): 行事委員会のメンバー

一般ユーザー (student): 大会に参加・観戦する全学生

### 1.4. 技術スタック
フロントエンド: SvelteKit

バックエンド: Gin (Go)

データベース: MySQL

その他: Dockerによる環境構築を推奨

## 2. 機能要件
### 2.1. ユーザー種別と権限

|ユーザー識別 | 認証方法 | 主な権限 |
| --- | --- | --- |
| superroot | 個別ID/パスワード | 全競技の試合結果登録・編集・管理者アカウントの管理 |
| admin | 競技ごとに発行された個別ID/パスワード | 担当競技の試合結果の登録・編集 |
| student | 全員共通のゲストパスワード | 全トーナメントの閲覧 |

### 2.2. 主な機能一覧
全ユーザー共通
- トップページ（行事委員会紹介）の閲覧
- ログイン機能

管理者 (superroot, admin)
- 管理者用ダッシュボードへのアクセス
- 担当競技のトーナメント一覧表示
- 試合ごとのスコア入力・更新機能
- （superrootのみ）adminアカウントの管理機能
- 全競技のトーナメント表の閲覧
- 試合結果のリアルタイム確認

一般ユーザー (student)
- 全競技のトーナメント表の閲覧
- 試合結果のリアルタイム確認

### 2.3. 特記事項（天候）
- 雨天時: 8人制サッカーは中止。卓球は「雨天時トーナメント」に切り替わる。
- 晴天時: 全てのトーナメントが実施される。

## 3. データベース設計 (ER図)
```
+-------------+      +----------------+      +-----------+      +-----------+
|    users    |      |  tournaments   |      |   teams   |      |  matches  |
+-------------+      +----------------+      +-----------+      +-----------+
| id (PK)     |      | id (PK)        |      | id (PK)   |      | id (PK)   |
| username    |--+   | name           |      | name      |      | tournament_id (FK) |
| hashed_pass |  |   | sport          |      | tournament_id (FK) |-->| round      |
| role        |  |   | weather_cond   |--+   +-----------+      | match_num   |
| assigned_sport|  |   +----------------+  |                      | team1_id (FK)    |--+
+-------------+  |                         |                      | team2_id (FK)    |--+
                 |                         |                      | team1_score  |  |
                 +-------------------------+--------------------->| team2_score  |  |
                                           |                      | winner_team_id (FK)|--+
                                           |                      | next_match_id (FK) |--> (self)
                                           +----------------------+-----------+--+

```

### 3.2. テーブル定義
users テーブル
ユーザー情報を管理する。
| カラム名 | データ型 | 説明 |
| :--- | :--- | :--- |
| id | INT (PK) | ユーザーID |
| username | VARCHAR | ユーザー名（ログインID） |
| hashed_password | VARCHAR | ハッシュ化されたパスワード |
| role | ENUM('superroot', 'admin', 'student') | ユーザーの役割 |
| assigned_sport | VARCHAR | adminの場合、担当する競技名を格納 (例: volleyball) |

tournaments テーブル
大会そのものを管理する。
| カラム名 | データ型 | 説明 |
| :--- | :--- | :--- |
| id | INT (PK) | 大会ID |
| name | VARCHAR | 大会名 (例: "令和7年度 秋季大会 バレーボール") |
| sport | VARCHAR | 競技名 (例: "volleyball", "soccer") |
| weather_condition | ENUM('sunny', 'rainy', 'any') | 開催天候条件 |

teams テーブル
参加チーム（クラス）を管理する。
| カラム名 | データ型 | 説明 |
| :--- | :--- | :--- |
| id | INT (PK) | チームID |
| name | VARCHAR | チーム名 (例: "3年A組") |
| tournament_id | INT (FK) | 参加している大会ID (tournaments.id) |

matches テーブル
各試合の情報を管理する。
| カラム名 | データ型 | 説明 |
| :--- | :--- | :--- |
| id | INT (PK) | 試合ID |
| tournament_id | INT (FK) | どの大会の試合か (tournaments.id) |
| round | INT | 回戦 (例: 1, 2, 3) |
| match_number_in_round | INT | 回戦内での試合番号 (例: 1回戦の第2試合) |
| team1_id | INT (FK, NULL可) | 対戦チーム1 (teams.id) |
| team2_id | INT (FK, NULL可) | 対戦チーム2 (teams.id) |
| team1_score | INT | チーム1のスコア |
| team2_score | INT | チーム2のスコア |
| winner_team_id | INT (FK, NULL可) | 勝者チーム (teams.id) |
| next_match_id | INT (FK, NULL可) | 勝者が進む次の試合ID (matches.id) |

## 4. APIエンドポイント設計 (主要なもの)
- POST /api/auth/login: ログイン処理
- POST /api/auth/logout: ログアウト処理
- GET /api/tournaments: 全トーナメント情報（試合情報含む）を取得
- GET /api/tournaments/:id: 指定したIDのトーナメント情報を取得
- PUT /api/matches/:id: (要認証: admin/superroot) 試合結果（スコア）を更新

## 5. 開発手順案
1. 環境構築: DockerでGo, Node.js, MySQLのコンテナを準備。
2. DBセットアップ: 上記設計に基づき、MySQLにテーブルを作成。初期データ（ユーザーアカウント、トーナメント枠）を投入。
3. バックエンド開発 (Gin):
    1. DB接続処理を実装。
    2. ユーザー認証APIを実装。
    3. トーナメント情報取得APIを実装。
    4. 試合結果更新APIを実装。

4. フロントエンド開発 (SvelteKit):
    1. ルーティング（トップページ, ログイン, トーナメント表示, 管理者ダッシュボード）を設定。
    2. 各ページのUIコンポーネントを作成。
    3. バックエンドAPIと連携し、データの表示・更新処理を実装。
    4. 結合・テスト: 全体の動作確認とデバッグ。

5. デプロイ: サーバーにアプリケーションを配置し公開。