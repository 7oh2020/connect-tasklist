# connect-list(backend)

これは connect-go を利用したタスクリストのサンプルです。
レイヤードアーキテクチャを採用しており各種ロジックを分離しています。
データベースは PostgreSQL + SQLC + PGX でアクセスしています。
認証後は JWX を使用して JWT を生成します。

## 起動方法

VS Code 拡張の devcontainer でプロジェクトを開くか、`docker-compose up`でコンテナを起動します。
その後以下のコマンドを実行します。すると RSA 秘密鍵の生成や依存関係のインストールからサーバー起動までが一括で実行されます。

```bash
make
```
