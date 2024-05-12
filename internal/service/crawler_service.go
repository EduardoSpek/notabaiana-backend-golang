package service

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"

	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/entity"
)

type CrawlerService struct {}

func NewCrawler() *CrawlerService {
	return &CrawlerService{}
}

func (c *CrawlerService) GetRSS(url string) entity.RSS {
	resp, err := http.Get(url)
    if err != nil {
        fmt.Println("Erro ao obter o feed:", err)
        return entity.RSS{}
    }
    defer resp.Body.Close()

    data, err := io.ReadAll(resp.Body)
    if err != nil {
        fmt.Println("Erro ao ler o feed:", err)
        return entity.RSS{}
    }

    //data := XML

    var rss entity.RSS
    if err := xml.Unmarshal(data, &rss); err != nil {
        fmt.Println("Erro ao decodificar o XML:", err)
        return entity.RSS{}
    }

	return rss
}