package main

import (
	"context"
	"log"
	"net/http"
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
	r := api.NewRouter(db)

	// シグナルを受け取るためのコンテキストを作成
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt, os.Kill)
	defer stop()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// errgroupを作成
	g, ctx := errgroup.WithContext(ctx)

	// シグナルを受け取り、サーバーをシャットダウンするゴルーチンをerrgroupで実行
	g.Go(func() error {
		// シグナルを受け取るまで待機
		<-ctx.Done()

		// 5秒のタイムアウト付きコンテキストを作成
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		// サーバーをシャットダウン(新しい接続の受け付けを停止し、contextがキャンセルされたら終了する)
		if err := srv.Shutdown(ctx); err != nil {
			return err
		}
		return nil
	})

	log.Println("server start at :8080")
	// 正常にシャットダウンした場合はhttp.ErrServerClosedが返る
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("ListenAndServe(): %v", err)
	}

	// サーバーがシャットダウンされるゴルーチンが終了するまで待機
	if err := g.Wait(); err != nil {
		log.Fatalf("g.Wait(): %v", err)
	}

}
