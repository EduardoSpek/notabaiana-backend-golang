package service

import (
	"fmt"
	"image"
	"io"
	"mime/multipart"
	"os"
	"strings"

	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/entity"
	"github.com/eduardospek/notabaiana-backend-golang/internal/domain/port"
)

var (
	image_dimensions = [2]int{1280, 1280}
)

type ContactService struct {
	ContactRepository port.ContactRepository
	imagedownloader   port.ImageDownloader
}

func NewContactService(contact_repository port.ContactRepository, downloader port.ImageDownloader) *ContactService {
	return &ContactService{ContactRepository: contact_repository, imagedownloader: downloader}
}

func (cs *ContactService) AdminCreate(contact entity.ContactDTO) (entity.ContactDTO, error) {

	newcontact := entity.NewContact(contact)
	_, err := newcontact.Validations()

	if err != nil {
		return entity.ContactDTO{}, err
	}

	contact, err = cs.ContactRepository.AdminCreate(*newcontact)

	if err != nil {
		return entity.ContactDTO{}, err
	}
	return contact, nil
}

func (cs *ContactService) AdminFindAll() ([]entity.ContactDTO, error) {

	lista, err := cs.ContactRepository.AdminFindAll()

	if err != nil {
		return []entity.ContactDTO{}, err
	}
	return lista, nil
}

func (cs *ContactService) AdminGetByID(id string) (entity.ContactDTO, error) {

	contact, err := cs.ContactRepository.AdminGetByID(id)

	if err != nil {
		return entity.ContactDTO{}, err
	}
	return contact, nil
}

func (cs *ContactService) AdminDelete(id string) error {

	err := cs.ContactRepository.AdminDelete(id)

	if err != nil {
		return err
	}
	return nil
}

func (cs *ContactService) AdminDeleteAll(contacts []entity.ContactDTO) error {

	err := cs.ContactRepository.AdminDeleteAll(contacts)

	if err != nil {
		return err
	}
	return nil
}

func (cs *ContactService) AdminCreateForm(image multipart.File, contact entity.ContactDTO) (entity.ContactDTO, error) {
	newcontact := entity.NewContact(contact)
	_, err := newcontact.Validations()

	if err != nil {
		return entity.ContactDTO{}, err
	}

	contactWithImages := cs.SaveImages(image, *newcontact)

	contactCreated, err := cs.ContactRepository.AdminCreate(contactWithImages)

	if err != nil {
		return entity.ContactDTO{}, err
	}

	return contactCreated, nil
}

func (cs *ContactService) SaveImages(image multipart.File, contact entity.Contact) entity.Contact {
	var file string
	var err error

	diretorio := "/images/contacts/"

	file = contact.ID + ".jpg"
	pathFile := diretorio + file

	err = cs.SaveImageForm(image, diretorio, file, image_dimensions[0], image_dimensions[1])

	if err != nil {
		pathFile = ""
	}

	if image != nil {
		contact.Image = pathFile
	}

	return contact
}

func (cs *ContactService) SaveImageForm(file multipart.File, diretorio, filename string, width, height int) error {

	if file == nil {
		return nil
	}

	defer file.Close()

	cwd, err := os.Getwd()

	if err != nil {
		fmt.Println("Erro ao obter o caminho do executável:", err)
	}

	pathImage := strings.Replace(cwd, "test", "", -1) + diretorio + filename

	f, err := os.Create(pathImage)
	if err != nil {
		return ErrParseForm
	}
	defer f.Close()
	io.Copy(f, file)

	f, err = os.Open(pathImage)

	if err != nil {
		return ErrParseForm
	}

	// Resize the image
	img, _, err := image.Decode(f)
	if err != nil {
		return ErrDecodeImage
	}

	err = cs.imagedownloader.SaveImage(img, width, height, pathImage)

	if err != nil {
		fmt.Println(err)
		return ErrDecodeImage
	}

	return nil

}
