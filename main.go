package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/api"
	"github.com/ShinnosukeSuzuki/techtrain-mission-ca-tech-dojo-golang/cache"
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

	// キャッシュを初期化
	characterCache, err := cache.NewCharacterProbabilityCache(
		os.Getenv("REGION"),
		os.Getenv("BUCKET_NAME"),
		os.Getenv("FILE_PATH"),
	)
	if err != nil {
		log.Fatalf("Failed to initialize cache: %v", err)
	}

	// キャッシュの定期更新を開始
	err = characterCache.StartCron()
	if err != nil {
		log.Fatalf("Failed to start cron cache update: %v", err)
	}

	// ルーターを作成
	e := api.NewRouter(db, characterCache)

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

		// キャッシュの定期更新を停止
		if err := characterCache.StopCron(); err != nil {
			log.Printf("Failed to stop cron cache update: %v", err)
		}

		// サーバーをシャットダウン
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
