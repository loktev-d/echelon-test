package server

import (
	"context"
	"database/sql"
	"echelon/thumbnail/pkg/api"
	"echelon/thumbnail/pkg/cache"
	"errors"
	"io"
	"log"
	"net/http"

	"google.golang.org/api/youtube/v3"
)

type Server struct {
	api.UnimplementedThumbnailServer

	ytSrv *youtube.Service
	db    *sql.DB
}

func (s *Server) GetThumbnail(ctx context.Context, req *api.ThumbnailRequest) (*api.ThumbnailResponse, error) {
	log.Printf("req videoId: %s", req.VideoId)

	fileBytes, err := cache.GetCache(s.db, req.GetVideoId())
	if err == nil {
		log.Printf("videoId '%s' is cached", req.VideoId)
		return &api.ThumbnailResponse{
			Thumbnail: fileBytes,
		}, nil
	} else if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	part := []string{"snippet"}
	ytRes, err := s.ytSrv.Videos.List(part).Id(req.GetVideoId()).Do()
	if err != nil {
		return nil, err
	}

	thumbnails := ytRes.Items[0].Snippet.Thumbnails
	var thumbnailUrl string

	if thumbnails.Maxres != nil {
		thumbnailUrl = thumbnails.Maxres.Url
	} else if thumbnails.Standard != nil {
		thumbnailUrl = thumbnails.Standard.Url
	} else if thumbnails.High != nil {
		thumbnailUrl = thumbnails.High.Url
	} else if thumbnails.Medium != nil {
		thumbnailUrl = thumbnails.Medium.Url
	} else if thumbnails.Default != nil {
		thumbnailUrl = thumbnails.Default.Url
	}

	httpRes, err := http.Get(thumbnailUrl)
	if err != nil {
		return nil, err
	}
	defer httpRes.Body.Close()

	bodyBytes, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return nil, err
	}

	if err = cache.InsertCache(s.db, req.GetVideoId(), bodyBytes); err != nil {
		return nil, err
	}

	return &api.ThumbnailResponse{
		Thumbnail: bodyBytes,
	}, nil
}

func NewServer(ytSrv *youtube.Service, db *sql.DB) *Server {
	return &Server{
		ytSrv: ytSrv,
		db:    db,
	}
}
