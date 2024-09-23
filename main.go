package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/api"
	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/db"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/sync/errgroup"
)

func main() {
	// DB接続
	db, err := db.NewDB()
	if err != nil {
		log.Println("failed to connect database", err)
		return
	}
	defer db.Close()

	// DB接続を確認し、接続が確認できない場合はサーバーを停止
	if err := db.Ping(); err != nil {
		log.Fatalf("server shutdown because db connection failed: %v", err)
	}

	// ルーターを作成
	e := api.NewRouter(db)

	// シグナルを受け取るためのコンテキストを作成
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt, os.Kill)
	defer stop()

	// errgroupを作成
	g, ctx := errgroup.WithContext(ctx)

	// シグナルを受け取り、サーバーをシャットダウンするゴルーチンをerrgroupで実行
	g.Go(func() error {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		return e.Shutdown(shutdownCtx)
	})

	// サーバーを起動
	g.Go(func() error {
		return e.Start(":8080")
	})

	log.Println("server start at :8080")

	// エラーが発生するまで待機
	if err := g.Wait(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
