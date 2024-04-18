package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/eduardospek/bn-api/internal/domain/entity"
	"github.com/eduardospek/bn-api/internal/utils"
	_ "github.com/lib/pq"
)

var (
	ErrNewsExists = errors.New("notícia já cadastrada com este título")	
    ErrNewsNotExistsWithID = errors.New("não existe notícia com este ID")
)

type NewsSupabaseRepository struct {}

func NewNewsSupabaseRepository() *NewsSupabaseRepository {
	repo := &NewsSupabaseRepository{}
	err := repo.CreateNewsTable()

    if err != nil {
        fmt.Println(err)
        panic(err)
    }
	return &NewsSupabaseRepository{}
}

func (repo *NewsSupabaseRepository) CreateNewsTable() error {    
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
    );`)
    if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Tabela 'news' criada com sucesso!")

	// Verificando se o trigger já existe
	var triggerExists bool
	err = db.QueryRow(`
		SELECT EXISTS (
			SELECT 1
			FROM   pg_trigger
			WHERE  tgname = 'news_update_trigger'
		);
	`).Scan(&triggerExists)
	
    if err != nil {
		log.Fatal(err)
	}

	// Criando o trigger somente se ele não existir
	if !triggerExists {
		_, err = db.Exec(`
			CREATE OR REPLACE FUNCTION update_updated_at()
			RETURNS TRIGGER AS $$
			BEGIN
				NEW.updated_at = now();
				RETURN NEW;
			END;
			$$ LANGUAGE plpgsql;
			
			CREATE TRIGGER news_update_trigger
			BEFORE UPDATE ON news
			FOR EACH ROW EXECUTE FUNCTION update_updated_at();
		`)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Trigger 'news_update_trigger' criada com sucesso!")
	} else {
		fmt.Println("Trigger 'news_update_trigger' já existe.")
	}
    
    return err
}

// insertNews insere um novo usuário no banco de dados
func (repo *NewsSupabaseRepository) Create(news entity.News) (entity.News, error) {    
    db, _ := conn.Connect()
	defer db.Close()    

    sql, err := db.Prepare(`INSERT INTO news (id, title, text, link, image, slug) VALUES ($1, $2, $3, $4, $5, $6)`)

    if err != nil {        
		return entity.News{}, err
	}  
    
    _, err = sql.Exec(news.ID, news.Title, news.Text, news.Link, news.Image, news.Slug)

    if err != nil {        
		return entity.News{}, err
	}     
    
    return news, err
}

func (repo *NewsSupabaseRepository) Update(news entity.News) (entity.News, error)  {    
    db, _ := conn.Connect()
	defer db.Close()
    
    _, err := repo.GetById(news.ID)
    if err != nil {
        fmt.Println(err)
		return entity.News{}, err
	}    

    query := "UPDATE news SET title = $1, text = $2, link = $3, image = $4, slug = $5 WHERE id = $6"

    sql, err := db.Prepare(query)
    
    if err != nil {
        fmt.Println(err)
        return entity.News{}, err
    }	


    if news.Title != "" || news.Text != "" || news.Link != "" || news.Image != "" || news.Slug != "" {
		_, err := sql.Exec(news.Title, news.Text, news.Link, news.Image, news.Slug, news.ID)
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

func (repo *NewsSupabaseRepository) GetById(id string) (entity.News, error) {
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

func (repo *NewsSupabaseRepository) GetBySlug(slug string) (entity.News, error) {
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

func (repo *NewsSupabaseRepository) FindAll(page, limit int) (interface{}, error) {
	
	db, err := conn.Connect()	

	if err != nil {        
		return nil, err
	}

    defer db.Close()

    offset := (page - 1) * limit

    sql, err := db.Prepare("SELECT * FROM news ORDER BY created_at DESC LIMIT $1 OFFSET $2")
    
    if err != nil {        
        return nil, err
    }  

    rows, err := sql.Query(limit, offset)
    if err != nil {        
        return nil, err
    }    

    defer sql.Close()

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

    sql, err = db.Prepare("SELECT COUNT(*) as total FROM news")

    if err != nil {
        return nil, err
    }

    row := sql.QueryRow()
    var total int

    err = row.Scan(&total)
    if err != nil {                
        return nil, err        
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

func (repo *NewsSupabaseRepository) Delete(id string) (error) {
	
	db, err := conn.Connect()	

    if err != nil {        
		return err
	}

    defer db.Close()

    _, err = repo.GetById(id)

    if err != nil {        
		return err
	}

    sql, err := db.Prepare("DELETE FROM news WHERE id = $1")

    if err != nil {
        return err
    }
    _ , err = sql.Exec(id)

    if err != nil {        
		return err
	}

    return nil

}

//VALIDATIONS
func (repo *NewsSupabaseRepository) NewsExists(title string) error {
    db, _ := conn.Connect()
	defer db.Close()    

    sql, err := db.Prepare("SELECT title FROM news WHERE title = $1")    

    if err != nil {
        return err
    }
    
    row := sql.QueryRow(title)    

    // Recuperando os valores do banco de dados
    err = row.Scan(&title)
    if err != nil {                   
        return nil
    }
  
    return ErrNewsExists
}