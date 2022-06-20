package main

import (
	"context"
	"echelon/thumbnail/pkg/api"
	"echelon/thumbnail/pkg/cache"
	"echelon/thumbnail/pkg/config"
	"echelon/thumbnail/pkg/server"
	"fmt"
	"log"
	"net"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"google.golang.org/grpc"
)

func main() {
	c, err := config.GetConfig()
	if err != nil {
		log.Fatal(err.Error())
	}

	ytSrv, err := youtube.NewService(context.Background(), option.WithAPIKey(c.GoogleApiKey))
	if err != nil {
		log.Fatal(err.Error())
	}

	db, err := cache.NewDB()
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()

	s := grpc.NewServer()
	srv := server.NewServer(ytSrv, db)

	api.RegisterThumbnailServer(s, srv)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", c.Port))
	if err != nil {
		log.Fatal(err.Error())
	}
	defer l.Close()

	log.Printf("listening on port %d", c.Port)

	if err := s.Serve(l); err != nil {
		log.Fatal(err.Error())
	}
}
