# Frontend

## 技術スタック (Technology Stack)

- **フレームワーク:** SvelteKit
- **UIライブラリ:** Tailwind CSS
- **データ可視化:** Bracketry (トーナメント表の描画に使用)

brackertryの基本的な使い方
```
import { createBracket } from 'bracketry'

const wrapper = document.querySelector('#your-wrapper-element')
const data = { ... } // data of a specific shape

createBracket(data, wrapper)
```

## ディレクトリ構成案 (Planned)

```
src/
├── lib/
│   ├── components/              # UIコンポーネント
│   │   ├── Tournament.svelte    # D3.jsを使ったトーナメント描画コンポーネント
│   │   ├── Match.svelte         # 個別の試合コンポーネント
│   │   └── Header.svelte        # ヘッダー
│   └── services/                # 外部との通信
│       └── api.js               # バックエンドAPIクライアント
├── routes/
│   ├── +page.svelte             # トップページ (行事委員会紹介)
│   ├── login/
│   │   └── +page.svelte         # ログインページ
│   ├── tournaments/
│   │   └── [sport]/             # 各競技のトーナメントページ (動的ルーティング)
│   │       └── +page.svelte
│   └── admin/                     # 管理者向けページ
│       ├── +layout.svelte       # 管理者ページの共通レイアウト
│       ├── +page.svelte         # 管理者ダッシュボード
│       └── users/               # (superroot用) ユーザー管理ページ
│           └── +page.svelte
│
├── App.svelte                      # ルートコンポーネント(行事委員会宣伝ページ。ここにログインボタンを含む。)
└── app.css                        # Tailwind CSSのベーススタイル
```