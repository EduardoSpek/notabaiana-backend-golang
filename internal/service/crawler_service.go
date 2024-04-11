package service

import (
	"encoding/xml"
	"fmt"

	"github.com/eduardospek/bn-api/internal/domain/entity"
)



type Crawler struct {
	List_news []entity.News
}

func NewCrawler(list_news []entity.News) *Crawler {
	return &Crawler{ List_news: list_news }
}

func (c *Crawler) GetRSS(url string) entity.RSS {
	// resp, err := http.Get(url)
    // if err != nil {
    //     fmt.Println("Erro ao obter o feed:", err)
    //     return entity.RSS{}
    // }
    // defer resp.Body.Close()

    // data, err := io.ReadAll(resp.Body)
    // if err != nil {
    //     fmt.Println("Erro ao ler o feed:", err)
    //     return entity.RSS{}
    // }

    data := XML

    var rss entity.RSS
    if err := xml.Unmarshal(data, &rss); err != nil {
        fmt.Println("Erro ao decodificar o XML:", err)
        return entity.RSS{}
    }

	return rss
}