package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/eduardospek/notabaiana-backend-golang/config"
	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/entity"
	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/port"
	usecase "github.com/eduardospek/notabaiana-backend-golang/internal/domain/usecase/download"
	"github.com/eduardospek/notabaiana-backend-golang/internal/utils"
)

type Response struct {
	PageProps PageProps `json:"pageProps"`
}

type PageProps struct {
	Top            []*Album       `json:"top"`
	Album          *Album         `json:"album"`
	RecommendedCds []*Album       `json:"recommendedCds"`
	AlbumsResponse AlbumsResponse `json:"albumsResponse"`
}

type AlbumsResponse struct {
	Albums []*Album `json:"albums"`
}

type Album struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Slug     string `json:"slug"`
	Cover    string `json:"cover"`
	BigCover string `json:"bigCover"`
	Username string `json:"username"`
	Name     string `json:"name"`
	CatName  string `json:"catName"`
	File     string `json:"file"`
	Files    []File `json:"files"`
}

type File struct {
	File     string `json:"file"`
	Path     string `json:"path"`
	Position int    `json:"position"`
}

type AlbumChan struct {
	Album *Album
	Error error
}

type CopierDownloadService struct {
	DownloadRepository port.CopierDownload
	ImgDownloader      port.ImageDownloader
}

func NewCopierDownload(DownloadRepository port.CopierDownload, ImgDownloader port.ImageDownloader) *CopierDownloadService {
	return &CopierDownloadService{DownloadRepository: DownloadRepository, ImgDownloader: ImgDownloader}
}

func (c *CopierDownloadService) Start(rss *[]string, minutes time.Duration) {

	go c.Run(rss)

	ticker := time.NewTicker(minutes * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		go c.Run(rss)
	}
}

func (c *CopierDownloadService) Run(list_downloads *[]string) {

	cwd, err := os.Getwd()

	if err != nil {
		fmt.Println("CopierDownload: Erro ao obter o caminho do executável:", err)
	}

	diretorio := strings.Replace(cwd, "test", "", -1) + "/images/downloads/"

	lista := c.Copier(list_downloads)

	for _, n := range lista {

		go func() {

			getByLinkDownloadUsecase := usecase.NewGetByLinkDownloadUsecase(c.DownloadRepository)
			downloadGet, _ := getByLinkDownloadUsecase.GetByLink(n.Link)

			if downloadGet.ID != "" {
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

func (s *CopierDownloadService) Copier(list_downloads *[]string) []*entity.Download {

	if list_downloads != nil {

		var response *Response
		var lista []*entity.Download
		var lista_albuns []*Album
		var item_atual string
		var files []*entity.Music

		var cover, category string
		var periodo = []string{"dia", "semana", "mes", "geral"}

		for _, item := range *list_downloads {

			//fmt.Printf("%d - %s\n", ii, item)

			if !strings.Contains(item, "recomendados") && !strings.Contains(item, "estourados") {
				periodo = []string{"dia", "semana", "mes", "geral"}
			} else {
				periodo = []string{"dia"}
			}

			for _, prd := range periodo {

				if item_atual == item {
					continue
				}

				url := strings.Replace(item, "dia", prd, -1)

				//fmt.Printf("%d - %s\n", iii, url)

				if !strings.Contains(url, "recomendados") && !strings.Contains(url, "estourados") && !strings.Contains(url, "recentes") {
					partes := strings.Split(item, "=")

					if partes[4] != "" {
						category = partes[4]
					}
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

				body = bytes.TrimPrefix(body, []byte("\xef\xbb\xbf"))

				// Decodificando o JSON
				err = json.Unmarshal(body, &response)

				if err != nil {

					version, err := strconv.Atoi(config.Suamusica_api_version)

					if err != nil {
						fmt.Println("Erro ao converter version", err)
					}
					version++
					config.Suamusica_api_version = strconv.Itoa(version)
					return nil
				}

				if strings.Contains(url, "recomendados") {
					lista_albuns = response.PageProps.RecommendedCds
				} else if strings.Contains(url, "recentes") {
					lista_albuns = response.PageProps.AlbumsResponse.Albums
				} else {
					lista_albuns = response.PageProps.Top
				}

				for _, album := range lista_albuns {

					done := make(chan *AlbumChan)

					go s.GetDataAlbum(album.Username, album.Slug, done)
					album_data := <-done

					if album_data.Error != nil {
						fmt.Println("Erro ao obter dados completos do album:", err)
						continue
					}

					if album_data.Album != nil {
						category = strings.ToLower(album_data.Album.CatName)

						if album_data.Album.BigCover != "" {
							cover = album_data.Album.BigCover
						} else if album_data.Album.Cover != "" {
							cover = album_data.Album.Cover
						} else {
							cover = album.Cover
						}

						if album_data.Album.Files != nil {
							files = nil
							for _, f := range album_data.Album.Files {
								files = append(files, &entity.Music{
									File:     f.File,
									Path:     f.Path,
									Position: f.Position,
								})
							}
						}
					}

					download := &entity.Download{
						Category: category,
						Title:    album.Title,
						Link:     config.Suamusica_UrlSite + "/" + album.Username + "/" + album.Slug,
						Image:    cover,
						Musics:   files,
					}

					lista = append(lista, download)

				}

				item_atual = item

			}

		}

		return lista

	} else {
		return nil
	}

}

func (s *CopierDownloadService) GetDataAlbum(username, slug string, done chan<- *AlbumChan) {

	var response Response
	var album *Album

	url := config.Suamusica_UrlSite + "/_next/data/webid-" + config.Suamusica_api_version + "/pt-BR/" + username + "/" + slug + ".json?slug=" + username

	// Fazendo a requisição GET
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Erro ao fazer a requisição:", err)
		done <- &AlbumChan{
			Album: nil,
			Error: err,
		}
	}
	defer resp.Body.Close()

	// Lendo o corpo da resposta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Erro ao ler o corpo da resposta:", err)
		fmt.Println("URL:", url)
		done <- &AlbumChan{
			Album: nil,
			Error: err,
		}
	}

	body = bytes.TrimPrefix(body, []byte("\xef\xbb\xbf"))

	// Decodificando o JSON
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("Erro ao decodificar o JSON:", err)
		done <- &AlbumChan{
			Album: nil,
			Error: err,
		}
	}

	album = response.PageProps.Album

	if album == nil {
		done <- &AlbumChan{
			Album: nil,
			Error: err,
		}
	}

	done <- &AlbumChan{
		Album: album,
		Error: nil,
	}
}
