package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"connectrpc.com/connect"
	"github.com/7oh2020/connect-tasklist/backend/infrastructure/persistence/model/db"
	"github.com/7oh2020/connect-tasklist/backend/interfaces/di"
	"github.com/7oh2020/connect-tasklist/backend/interfaces/interceptor"
	"github.com/7oh2020/connect-tasklist/backend/interfaces/rpc/auth/v1/auth_v1connect"
	"github.com/7oh2020/connect-tasklist/backend/interfaces/rpc/task/v1/task_v1connect"
	"github.com/7oh2020/connect-tasklist/backend/interfaces/rpc/user/v1/user_v1connect"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {

	// 環境変数の取得
	var issuer, keyPath, url string
	var ok bool
	if issuer, ok = os.LookupEnv("REPO_NAME"); !ok {
		return fmt.Errorf("issuer not set: %s", issuer)
	}
	if keyPath, ok = os.LookupEnv("PRIVATE_KEY_PATH"); !ok {
		return fmt.Errorf("private-key-path not set: %s", keyPath)
	}
	if url, ok = os.LookupEnv("DATABASE_URL"); !ok {
		return fmt.Errorf("database-url not set: %s", url)
	}

	// PostgreSQLに接続する
	poolCfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		return fmt.Errorf("unable to parse pool-config: %s", url)
	}
	poolCfg.MaxConns = 10
	poolCfg.MinConns = 1
	pool, err := pgxpool.ConnectConfig(context.Background(), poolCfg)
	if err != nil {
		return fmt.Errorf("unable to connect pool: %s", url)
	}
	defer pool.Close()

	// SQLCのクライアントを作成する
	qry := db.New(pool)

	// JWTの有効期限
	timeout := 1 * time.Hour

	// ハンドラを作成する
	authServer, err := di.InitAuth(issuer, keyPath, qry, timeout)
	if err != nil {
		return err
	}
	userServer := di.InitUser(qry)
	taskServer := di.InitTask(qry)

	// インターセプタを作成する
	authInterceptor := connect.WithInterceptors(interceptor.NewAuthInterceptor(issuer, keyPath))

	// サーバーの起動
	mux := http.NewServeMux()
	mux.Handle(auth_v1connect.NewAuthServiceHandler(authServer))
	mux.Handle(user_v1connect.NewUserServiceHandler(userServer))
	mux.Handle(task_v1connect.NewTaskServiceHandler(taskServer, authInterceptor))

	return http.ListenAndServe(
		"localhost:8080",
		// CORSハンドラでリクエストを許可する(本番では全ホスト許可ではなく特定のホストのみ許可するべき)
		cors.AllowAll().Handler(
			// HTTP1.1リクエストはHTTP/2にアップグレードされる
			h2c.NewHandler(mux, &http2.Server{}),
		),
	)

}
