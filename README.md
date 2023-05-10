# connect-tasklist

これは Connect を利用したタスクリストのサンプルアプリです。

- バックエンド: connect-go(Go 言語)
- フロントエンド: connect-query(TypeScript)

## 起動方法

最初に API サーバーを起動する必要があります。
backend ディレクトリへ移動し VS Code の devcontainer でプロジェクトを開始または`docker-compose up`でコンテナを起動します。
その後ターミナルで以下のコマンドを実行するとサーバーが起動します。

```bash
make
```

次に frontend ディレクトリへ移動します。
別のターミナルで以下のコマンドを実行すると Vite の開発サーバーが起動します。

```bash
pnpm install
pnpm dev
```

WEB ブラウザから`localhost:5173`にアクセスします。
データベースにはテストユーザーが登録済みなので「test@example.com / pass」でログインできます。
ログインに成功するとタスク一覧ページが表示されます。
