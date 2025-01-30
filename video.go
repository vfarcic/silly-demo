package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/jackc/pgx/v5"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

type Video struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

func getConn() *pgx.Conn {
	url := os.Getenv("DB_URI")
	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		slog.Error("Failed to parse DB_URL", "error", err)
		return nil
	}
	return conn
}

func videosGetHandler(ctx *gin.Context) {
	slog.Debug("Handling request", "URI", ctx.Request.RequestURI)
	var videos []Video
	if strings.ToLower(os.Getenv("DB")) == "fs" {
		var err error
		videos, err = getVideosFromFile()
		if err != nil {
			httpErrorInternalServerError(err, ctx)
			return
		}
	} else {
		conn := getConn()
		if conn == nil {
			return
		}
		defer conn.Close(context.Background())
		rows, err := conn.Query(context.Background(), "SELECT id, title FROM videos")
		if err != nil {
			httpErrorInternalServerError(err, ctx)
			return
		}
		defer rows.Close()
		for rows.Next() {
			var video Video
			err := rows.Scan(&video.ID, &video.Title)
			if err != nil {
				httpErrorInternalServerError(err, ctx)
				return
			}
			videos = append(videos, video)
		}
	}
	ctx.JSON(http.StatusOK, videos)
}

func getVideosFromFile() ([]Video, error) {
	dir := os.Getenv("FS_DIR")
	if len(dir) == 0 {
		dir = "/cache"
	}
	path := fmt.Sprintf("%s/videos.yaml", dir)
	// Create a file if it doesn't exist
	if _, err := os.Stat(path); os.IsNotExist(err) {
		_, err := os.Create(path)
		if err != nil {
			return nil, err
		}
	}
	var videos []Video
	yamlData, err := os.ReadFile(path)
	if err != nil {
		return videos, err
	}
	err = yaml.Unmarshal(yamlData, &videos)
	return videos, err
}

func videoPostHandler(ctx *gin.Context) {
	slog.Debug("Handling request", "URI", ctx.Request.RequestURI)
	id := ctx.Query("id")
	if len(id) == 0 {
		httpErrorBadRequest(errors.New("id is empty"), ctx)
		return
	}
	title := ctx.Query("title")
	if len(title) == 0 {
		httpErrorBadRequest(errors.New("title is empty"), ctx)
		return
	}
	video := &Video{
		ID:    id,
		Title: title,
	}
	if strings.ToLower(os.Getenv("DB")) == "fs" {
		videos, err := getVideosFromFile()
		if err != nil {
			httpErrorInternalServerError(err, ctx)
			return
		}
		videos = append(videos, *video)
		dir := os.Getenv("FS_DIR")
		if len(dir) == 0 {
			dir = "/cache"
		}
		path := fmt.Sprintf("%s/videos.yaml", dir)
		yamlData, err := yaml.Marshal(videos)
		if err != nil {
			httpErrorInternalServerError(err, ctx)
			return
		}
		err = os.WriteFile(path, yamlData, 0644)
		if err != nil {
			httpErrorInternalServerError(err, ctx)
		}
	} else {
		conn := getConn()
		if conn == nil {
			return
		}
		defer conn.Close(context.Background())
		_, err := conn.Exec(context.Background(), "INSERT INTO videos(id, title) VALUES ($1, $2)", video.ID, video.Title)
		if err != nil {
			httpErrorInternalServerError(err, ctx)
			return
		}
	}
}
