package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/entity"
	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/port"
	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/usecase"
	"github.com/eduardospek/notabaiana-backend-golang/internal/utils"
)

var (
	urlSite = "https://suamusica.com.br"
)

type Response struct {
	PageProps PageProps `json:"pageProps"`
}

type PageProps struct {
	Top []Album `json:"top"`
}

type Album struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Slug     string `json:"slug"`
	Cover    string `json:"cover"`
	BigCover string `json:"bigCover"`
	Username string `json:"username"`
	Name     string `json:"name"`
}

type CopierDownloadService struct {
	DownloadRepository port.DownloadRepository
	ImgDownloader      port.ImageDownloader
}

func NewCopierDownload(DownloadRepository port.DownloadRepository, ImgDownloader port.ImageDownloader) *CopierDownloadService {
	return &CopierDownloadService{DownloadRepository: DownloadRepository, ImgDownloader: ImgDownloader}
}

func (c *CopierDownloadService) Start(rss []string, minutes time.Duration) {

	go c.Run(rss)

	ticker := time.NewTicker(minutes * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		go c.Run(rss)
	}
}

func (c *CopierDownloadService) Run(list_downloads []string) {

	cwd, err := os.Getwd()

	if err != nil {
		fmt.Println("CopierDownload: Erro ao obter o caminho do executável:", err)
	}

	diretorio := strings.Replace(cwd, "test", "", -1) + "/images/downloads/"

	lista := c.Copier(list_downloads)

	for _, n := range lista {

		go func() {

			getByLinkDownloadUsecase := usecase.NewGetByLinkDownloadUsecase(c.DownloadRepository)
			downloadGet, err := getByLinkDownloadUsecase.GetByLink(n.Link)

			if err != nil {
				return
			}

			if downloadGet != nil {
				return
			}

			createDownloadUsecase := usecase.NewCreateDownloadUsecase(c.DownloadRepository)
			downloadCreated, err := createDownloadUsecase.Create(n)

			if err != nil {
				fmt.Println("CopierDownload: ", err)
			}

			img, err := utils.DownloadImage(downloadCreated.Image)

			if err != nil {
				downloadCreated.Image = ""
				updateDownloadUsecase := usecase.NewUpdateDownloadUsecase(c.DownloadRepository)
				updateDownloadUsecase.Update(downloadCreated)
				return
			}

			outputPath := diretorio + downloadCreated.ID + ".jpg"

			width := 300
			height := int(float64(img.Bounds().Dy()) * (float64(width) / float64(img.Bounds().Dx())))

			err = c.ImgDownloader.CropAndSaveImage(img, width, height, outputPath)

			if err != nil {
				downloadCreated.Image = ""
				updateDownloadUsecase := usecase.NewUpdateDownloadUsecase(c.DownloadRepository)
				updateDownloadUsecase.Update(downloadCreated)

				fmt.Println("CopierDownload: Erro ao salvar a imagem:", err)
			} else {
				downloadCreated.Image = downloadCreated.ID + ".jpg"
				updateDownloadUsecase := usecase.NewUpdateDownloadUsecase(c.DownloadRepository)
				updateDownloadUsecase.Update(downloadCreated)
			}

		}()
	}
}

func (s *CopierDownloadService) Copier(list_downloads []string) []*entity.Download {

	var response Response
	var lista []*entity.Download

	var category string
	var periodo = []string{"dia", "semana", "mes", "geral"}

	for _, item := range list_downloads {

		for _, prd := range periodo {
			url := strings.Replace(item, "dia", prd, -1)

			partes := strings.Split(item, "=")

			if partes[1] != "" {
				category = partes[4]
			}

			// Fazendo a requisição GET
			resp, err := http.Get(url)
			if err != nil {
				fmt.Println("Erro ao fazer a requisição:", err)
			}
			defer resp.Body.Close()

			// Lendo o corpo da resposta
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Erro ao ler o corpo da resposta:", err)
			}

			// Decodificando o JSON
			err = json.Unmarshal(body, &response)
			if err != nil {
				fmt.Println("Erro ao decodificar o JSON:", err)
			}

			for _, album := range response.PageProps.Top {
				download := &entity.Download{
					Category: category,
					Title:    album.Title,
					Link:     urlSite + "/" + album.Username + "/" + album.Slug,
					Image:    album.BigCover,
				}

				lista = append(lista, download)

			}
		}

	}

	return lista
}
