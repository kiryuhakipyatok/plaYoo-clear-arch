package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"playoo/internal/domain/entity"
	"playoo/internal/domain/repository"
	"time"
)

type NewsService interface {
	CreateNews(c context.Context, title, body, link string, picture *multipart.FileHeader) (*entity.News, error)
	GetById(c context.Context, id string) (*entity.News, error)
	GetByAmount(c context.Context, amount int) ([]entity.News, error)
}

type newsService struct {
	NewsRepository    repository.NewsRepository
	CommentRepository repository.CommentRepository
}

func NewNewsService(newsRepository repository.NewsRepository) NewsService {
	return &newsService{
		NewsRepository: newsRepository,
	}
}

func (nr newsService) CreateNews(c context.Context, title, body, link string, picture *multipart.FileHeader) (*entity.News, error) {
	news := entity.News{
		Id:    uuid.New(),
		Title: title,
		Body:  body,
		Time:  time.Now(),
		Link:  link,
	}
	uploadDir := "../../files/news-pictures"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return nil, err
	}
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("upload directory does not exist: %s", uploadDir)
	}
	fileName := fmt.Sprintf("%s-news-picture%s", news.Id, filepath.Ext(picture.Filename))
	filepath := filepath.Join(uploadDir, fileName)
	dst, err := os.Create(filepath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()
	src, err := picture.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()
	_, err = io.Copy(dst, src)
	if err != nil {
		return nil, err
	}

	var (
		host = os.Getenv("HOST")
		port = os.Getenv("PORT")
	)

	fileURL := fmt.Sprintf("http://%s:%s/files/news-pictures/%s", host, port, fileName)

	news.Picture = fileURL
	if err := nr.NewsRepository.Create(c, news); err != nil {
		return nil, err
	}
	return &news, nil
}

func (nr newsService) GetById(c context.Context, id string) (*entity.News, error) {
	news, err := nr.NewsRepository.FindById(c, id)
	if err != nil {
		return nil, err
	}
	return news, nil
}

func (nr newsService) GetByAmount(c context.Context, amount int) ([]entity.News, error) {
	somenews, err := nr.NewsRepository.FindByAmount(c, amount)
	if err != nil {
		return nil, err
	}
	return somenews, nil
}
