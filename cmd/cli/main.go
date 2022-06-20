package main

import (
	"bytes"
	"context"
	"echelon/thumbnail/pkg/api"
	"echelon/thumbnail/pkg/config"
	"flag"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	var isAsync bool
	flag.BoolVar(&isAsync, "async", false, "is async mode")
	flag.Parse()

	c, err := config.GetConfig()
	if err != nil {
		log.Fatal(err.Error())
	}

	conn, err := grpc.Dial(fmt.Sprintf(":%d", c.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err.Error())
	}
	defer conn.Close()

	client := api.NewThumbnailClient(conn)

	args := os.Args[1:]
	for _, arg := range args {
		if isAsync {
			go downloadImage(arg, client)
		} else {
			downloadImage(arg, client)
		}
	}
}

func downloadImage(ytUrl string, client api.ThumbnailClient) error {
	u, err := url.Parse(ytUrl)
	if err != nil {
		return err
	}

	id := u.Query().Get("v")
	req := &api.ThumbnailRequest{
		VideoId: id,
	}

	res, err := client.GetThumbnail(context.Background(), req)
	if err != nil {
		log.Fatal(err.Error())
	}

	r := bytes.NewReader(res.Thumbnail)

	out, err := os.Create(fmt.Sprintf("./%s.jpg", id))
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, r)
	if err != nil {
		return err
	}

	return nil
}
