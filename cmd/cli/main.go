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
	"path/filepath"
	"sync"

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

	var urls []string
	args := os.Args[1:]

	for _, arg := range args {
		if string(arg[0]) == "-" {
			continue
		}

		urls = append(urls, arg)
	}

	var wg sync.WaitGroup
	wg.Add(len(urls))

	for _, url := range urls {
		if isAsync {
			go downloadImage(url, client, &wg)
		} else {
			downloadImage(url, client, nil)
		}
	}

	wg.Wait()
}

func downloadImage(ytUrl string, client api.ThumbnailClient, wg *sync.WaitGroup) {
	fmt.Printf("downloading from source: %s\n", ytUrl)

	if wg != nil {
		defer wg.Done()
	}

	u, err := url.Parse(ytUrl)
	if err != nil {
		log.Fatal(err.Error())
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
		log.Fatal(err.Error())
	}
	defer out.Close()

	_, err = io.Copy(out, r)
	if err != nil {
		log.Fatal(err.Error())
	}

	abs, err := filepath.Abs(out.Name())
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Printf("succeeded '%s', image path: %s\n", ytUrl, abs)
}
