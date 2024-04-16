package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/eduardospek/bn-api/internal/domain/entity"
	"github.com/eduardospek/bn-api/internal/utils"
	_ "github.com/mattn/go-sqlite3"
)

var (
	ErrNewsExists = errors.New("notícia já cadastrada com este título")	
    ErrNewsNotExistsWithID = errors.New("não existe notícia com este ID")
)

type NewsSQLiteRepository struct {}

func NewNewsSQLiteRepository() *NewsSQLiteRepository {
	repo := &NewsSQLiteRepository{}
	repo.CreateNewsTable()
	return &NewsSQLiteRepository{}
}

func (repo *NewsSQLiteRepository) CreateNewsTable() error {    
    db, err := conn.Connect()
	
	if err != nil {
        fmt.Println(err)
		return err
	}

    defer db.Close()
    _, err = db.Exec(`CREATE TABLE IF NOT EXISTS news (
        id VARCHAR(36) PRIMARY KEY NOT NULL,
        title VARCHAR(250) NOT NULL,
        text TEXT NOT NULL,
        link VARCHAR(250) NOT NULL,
        image VARCHAR(250) NOT NULL,
        slug VARCHAR(250) NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
    CREATE TRIGGER update_news
	AFTER UPDATE ON news
	FOR EACH ROW
	BEGIN
		UPDATE news SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
	END;`)
    return err
}

// insertNews insere um novo usuário no banco de dados
func (repo *NewsSQLiteRepository) Create(news entity.News) (entity.News, error) {    
    db, _ := conn.Connect()
	defer db.Close()

    insertQuery := "INSERT INTO news (id, title, text, link, image, slug) VALUES (?, ?, ?, ?, ?, ?)"
    _, err := db.Exec(insertQuery, news.ID, news.Title, news.Text, news.Link, news.Image, news.Slug)

    if err != nil {
		return entity.News{}, err
	}     
    
    return news, err
}

func (repo *NewsSQLiteRepository) Update(news entity.News) (entity.News, error)  {    
    db, _ := conn.Connect()
	defer db.Close()
    
    _, err := repo.GetById(news.ID)
    if err != nil {
        fmt.Println(err)
		return entity.News{}, err
	}    

    query := "UPDATE news SET"
    if news.Title != "" {
		query += " name = '" + news.Title + "'"
	}
    if news.Text != "" {
		if news.Title != "" {
			query += ","
		}
		query += " text = '" + news.Text + "'"
	}
    if news.Link != "" {
		if news.Title != "" {
			query += ","
		}
		query += " text = '" + news.Link + "'"
	}
    if news.Image != "" {
		if news.Title != "" {
			query += ","
		}
		query += " text = '" + news.Image + "'"
	}
    if news.Slug != "" {
		if news.Title != "" {
			query += ","
		}
		query += " text = '" + news.Slug + "'"
	}
	query += " WHERE id = '" + fmt.Sprint(news.ID) + "'"    

    if news.Title != "" || news.Text != "" || news.Link != "" || news.Image != "" || news.Slug != "" {
		_, err := db.Exec(query)
		if err != nil {
			fmt.Println(err)
			return entity.News{}, err
		}	
	}

    updatenews, err := repo.GetById(news.ID)

    if err != nil {
        fmt.Println(err)
		return entity.News{}, err
	}

    return updatenews, err
}

func (repo *NewsSQLiteRepository) GetById(id string) (entity.News, error) {
	db, err := conn.Connect()	
	
	if err != nil {
        fmt.Println("Erro ao conectar ao DB")
		return entity.News{}, err
	}   
    
    defer db.Close()

    newsQuery := "SELECT * FROM news WHERE id = ?"
    row := db.QueryRow(newsQuery, id)    

    // Variáveis para armazenar os dados do usuário
    var title, text, link, image, slug string
    var created_at, updated_at time.Time

    // Recuperando os valores do banco de dados
    err = row.Scan(&id, &title, &text, &link, &image, &slug, &created_at,  &updated_at)
    if err != nil {        
        // Se não houver usuário correspondente ao ID fornecido, retornar nil
        if err == sql.ErrNoRows {            
            return entity.News{}, ErrNewsNotExistsWithID
        }
        // Se ocorrer outro erro, retornar o erro        
        return entity.News{}, err
    }

    // Criando um objeto models.News com os dados recuperados
    news := &entity.News{
        ID: id,
        Title: title,
        Text:    text,
        Link: link,
        Image: image,
        Slug: slug,
        CreatedAt: created_at.Local(),
        UpdatedAt: updated_at.Local(),
    }
    
    return *news, err
}

func (repo *NewsSQLiteRepository) GetBySlug(slug string) (entity.News, error) {
	db, err := conn.Connect()	
	
	if err != nil {
        fmt.Println("Erro ao conectar ao DB")
		return entity.News{}, err
	}   
    
    defer db.Close()

    newsQuery := "SELECT * FROM news WHERE slug = ?"
    row := db.QueryRow(newsQuery, slug)    

    // Variáveis para armazenar os dados do usuário
    var id, title, text, link, image string
    var created_at, updated_at time.Time

    // Recuperando os valores do banco de dados
    err = row.Scan(&id, &title, &text, &link, &image, &slug, &created_at,  &updated_at)
    if err != nil {        
        // Se não houver usuário correspondente ao ID fornecido, retornar nil
        if err == sql.ErrNoRows {            
            return entity.News{}, ErrNewsNotExistsWithID
        }
        // Se ocorrer outro erro, retornar o erro        
        return entity.News{}, err
    }

    // Criando um objeto models.News com os dados recuperados
    news := &entity.News{
        ID: id,
        Title: title,
        Text:    text,
        Link: link,
        Image: image,
        Slug: slug,
        CreatedAt: created_at.Local(),
        UpdatedAt: updated_at.Local(),
    }
    
    return *news, err
}

func (repo *NewsSQLiteRepository) FindAll(page, limit int) (interface{}, error) {
	
	db, err := conn.Connect()	

	if err != nil {        
		return nil, err
	}

    defer db.Close()

    offset := (page - 1) * limit
    
    rows, err := db.Query("SELECT * FROM news ORDER BY id DESC LIMIT ? OFFSET ?", limit, offset)
    if err != nil {        
        return nil, err
    }    

    defer rows.Close()

    var list_news []entity.News
    list_news = []entity.News{}
    
    for rows.Next() {
        var news entity.News
        err := rows.Scan(&news.ID, &news.Title, &news.Text, &news.Link, &news.Image, &news.Slug, &news.CreatedAt, &news.UpdatedAt)
        if err != nil {            
            return nil, err
        }
        news.CreatedAt = news.CreatedAt.Local()
        news.UpdatedAt = news.UpdatedAt.Local()
        list_news = append(list_news, news)
    }

    countQuery := "SELECT COUNT(*) as total FROM news"
    row := db.QueryRow(countQuery)
    var total int

    err = row.Scan(&total)
    if err != nil {        
        if err == sql.ErrNoRows {            
            return nil, err
        }
    }

    pagination := utils.Pagination(page, total)

    result := struct{
        List_news []entity.News `json:"news"`
        Pagination map[string][]int `json:"pagination"`
    }{
        List_news: list_news,
        Pagination: pagination,
    }
    
    return result, nil
}

func (repo *NewsSQLiteRepository) Delete(id string) (error) {
	
	db, err := conn.Connect()	

    if err != nil {        
		return err
	}

    defer db.Close()

    _, err = repo.GetById(id)

    if err != nil {        
		return err
	}

    _ , err = db.Exec("DELETE FROM news WHERE id = ?", id)

    if err != nil {        
		return err
	}

    return nil

}

//VALIDATIONS
func (repo *NewsSQLiteRepository) NewsExists(title string) error {
    db, _ := conn.Connect()
	defer db.Close()

    title = strings.TrimSpace(title)

    newsQuery := "SELECT title FROM news WHERE title = ?"
    row := db.QueryRow(newsQuery, title)    

    // Recuperando os valores do banco de dados
    err := row.Scan(&title)
    if err != nil {        
        if err == sql.ErrNoRows {            
            return nil
        }
    }
  
    return errors.New("já existe notícia com este título")
}