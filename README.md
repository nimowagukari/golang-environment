# golang-environment

このリポジトリは個人的な go 言語のローカル開発環境用のリポジトリです。

## 学習ネタ

- 標準ライブラリお試し
  - flag
  - log
  - os
  - path
  - testing
  - sql
  - net
  - net/http
  - html
  - html/template
  - text/template
  - text/template/parse
- 自作ツールの作成
  - コマンドライン
    - CSV と JSON の相互変換
    - 同一スキーマの複数 DB への一括クエリ実行
    - SPF Record の再帰チェック
    - サーバに登録された Cron の検索・分析
  - Web API
    - メソッド検証：GET, POST, PUT, PATCH, DELETE
    - パラメータ渡し：クエリ文字列、ヘッダ、データ
    - エンコーディング：圧縮リクエストの受取
    - 画像などのバイナリの処理
    - Basic/Digest 認証
    - CORS
    - テストコード実装
- CI/CD 実装
- パフォーマンス検証
