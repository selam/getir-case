package main

import (
	"context"
	"flag"
	"fmt"
	"getircase/databases"
	"getircase/handlers"
	"getircase/lib/config"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//

func main() {
	cfgFile := flag.String("config", "./config.json", "absolute or relative path of configuration file")
	flag.Parse()

	if cfgFile == nil || *cfgFile == "" {
		os.Exit(-1)
	}

	cfg, err := config.Parse(cfgFile)
	if err != nil {
		log.Fatalf("error while parsing configuration file: %s", err.Error())
	}
	var mongoConnection databases.MongoClient
	var redisConnection *databases.RedisConnection
	var inmemoryConnection databases.Inmemory
	// initialize database connections
	for idx := range cfg.Databases {
		switch cfg.Databases[idx].Type {
		case "mongodb":
			if mongoConnection, err = databases.InitializeMongodb(cfg.Databases[idx]); err != nil {
				log.Fatalf("can't connect to mongodb: %s", err.Error())
			}
		case "inmemory":
			if inmemoryConnection, err = databases.InitializeInmemory(cfg.Databases[idx]); err != nil {
				log.Fatalf("can't initialize in memorydb: %s", err.Error())
			}
		case "redis":
			if redisConnection, err = databases.InitializeRedis(cfg.Databases[idx]); err != nil {
				log.Fatalf("can't connect to redis: %s", err.Error())
			}
		}
	}
	fmt.Printf("%+v", inmemoryConnection)
	// parse application flags
	// create http mux from std lib of go
	mux := http.NewServeMux()
	mux.Handle("/mongodb/records", handlers.NewMongodbHandler(mongoConnection))
	mux.Handle("/redis", handlers.NewRedisHandler(redisConnection))

	//  inmemory term is not clear in case file
	//  as any in memory service like redis, memcache etc or in memory structure in application.
	//  so i use sync map for in memory local storage also implement same functionality with redis too
	mux.Handle("/inmemory", handlers.NewInmemoryHandler(inmemoryConnection))

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.Application.Host, cfg.Application.Port),
		Handler: mux,
	}
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		fmt.Printf("Http server running on %s!\n", server.Addr)
		if err := server.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				log.Println("[Shutdown] Listeners are closed, waiting to finish opened connections")
				return
			}
			log.Println(err.Error())
			os.Exit(-1)
		}

	}()

	<-stop
	fmt.Println("")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		if err == context.DeadlineExceeded {
			log.Println("[Shutdown] connections can't finish their job for given time")
			return
		}
		log.Printf("[Shutdown] Error while shutdown %s", err.Error())
		os.Exit(-1)
	}

	log.Println("[Shutdown] closed")

}
