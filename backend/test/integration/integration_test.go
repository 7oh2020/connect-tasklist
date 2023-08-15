package integration

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/7oh2020/connect-tasklist/backend/infrastructure/persistence/model/db"
	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	qry     db.Querier
	issuer  string
	keyPath string
	timeout time.Duration
)

type testResponse struct {
	status int
	body   string
	header http.Header
}

type testServer struct {
	*httptest.Server
}

func newTestServer(t *testing.T, hdr http.Handler) *testServer {
	ts := httptest.NewTLSServer(hdr)
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &testServer{ts}
}

func (ts *testServer) sendPostRequest(t *testing.T, token string, urlPath string, input string) (*testResponse, error) {
	req, err := http.NewRequest(http.MethodPost, ts.URL+urlPath, bytes.NewBuffer([]byte(input)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := ts.Client().Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return &testResponse{
		status: resp.StatusCode,
		body:   string(body),
		header: resp.Header,
	}, nil
}

func TestMain(m *testing.M) {
	// 環境変数を取得する
	var envIssuer, envKeyPath, dbURL string
	var ok bool

	if envIssuer, ok = os.LookupEnv("REPO_NAME"); !ok {
		panic(fmt.Errorf("issuer not set: %s", issuer))
	}
	issuer = envIssuer

	if envKeyPath, ok = os.LookupEnv("PRIVATE_KEY_PATH"); !ok {
		panic(fmt.Errorf("private-key-path not set: %s", keyPath))
	}
	keyPath = envKeyPath

	if dbURL, ok = os.LookupEnv("DATABASE_URL"); !ok {
		panic(fmt.Errorf("database-url not set: %s", dbURL))
	}

	// PostgreSQLに接続する
	poolCfg, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		panic(fmt.Errorf("unable to parse pool-config: %s", dbURL))
	}
	poolCfg.MaxConns = 10
	poolCfg.MinConns = 1
	pool, err := pgxpool.ConnectConfig(context.Background(), poolCfg)
	if err != nil {
		panic(fmt.Errorf("unable to connect pool: %s", dbURL))
	}
	defer pool.Close()

	// SQLCのクライアントを作成する
	qry = db.New(pool)

	// タイムアウトの設定
	timeout = time.Hour

	// テストの開始
	code := m.Run()

	os.Exit(code)
}
