package service

import (
	"errors"
	"time"

	"github.com/wafiydev/okefin-service/internal/dto"
	"github.com/wafiydev/okefin-service/internal/models"
	"github.com/wafiydev/okefin-service/internal/repository"
)

type TokoService interface {
	CreateToko(userID uint, req *dto.CreateTokoRequest) (*dto.TokoWithUserResponse, error)
	GetTokoByID(id uint, userID uint) (*dto.TokoWithUserResponse, error)
	GetTokoByUserID(userID uint) (*dto.TokoWithUserResponse, error)
	UpdateToko(id uint, userID uint, req *dto.UpdateTokoRequest) (*dto.TokoWithUserResponse, error)
	DeleteToko(id uint, userID uint) error
}

type tokoService struct {
	tokoRepo repository.TokoRepository
}

func NewTokoService(tokoRepo repository.TokoRepository) TokoService {
	return &tokoService{tokoRepo: tokoRepo}
}

func (s *tokoService) CreateToko(userID uint, req *dto.CreateTokoRequest) (*dto.TokoWithUserResponse, error) {
	toko := &models.Toko{
		IDUser:    userID,
		NamaToko:  req.NamaToko,
		UrlFoto:   req.UrlFoto,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := s.tokoRepo.CreateToko(toko)
	if err != nil {
		return nil, err
	}

	// Fetch the created toko with user details
	fetchedToko, err := s.tokoRepo.GetTokoByID(toko.ID)
	if err != nil {
		return nil, err
	}

	response := &dto.TokoWithUserResponse{
		ID:       fetchedToko.ID,
		NamaToko: fetchedToko.NamaToko,
		UrlFoto:  fetchedToko.UrlFoto,
		User:     dto.UserResponse{}, // Assuming UserResponse is defined; adjust if needed
	}

	return response, nil
}

func (s *tokoService) GetTokoByID(id uint, userID uint) (*dto.TokoWithUserResponse, error) {
	toko, err := s.tokoRepo.GetTokoByID(id)
	if err != nil {
		return nil, err
	}
	if toko.IDUser != userID {
		return nil, errors.New("access denied")
	}

	response := &dto.TokoWithUserResponse{
		ID:       toko.ID,
		NamaToko: toko.NamaToko,
		UrlFoto:  toko.UrlFoto,
		User:     dto.UserResponse{}, // Adjust if UserResponse is defined
	}

	return response, nil
}

func (s *tokoService) GetTokoByUserID(userID uint) (*dto.TokoWithUserResponse, error) {
	toko, err := s.tokoRepo.GetTokoByUserID(userID)
	if err != nil {
		return nil, err
	}

	response := &dto.TokoWithUserResponse{
		ID:       toko.ID,
		NamaToko: toko.NamaToko,
		UrlFoto:  toko.UrlFoto,
		User:     dto.UserResponse{}, // Adjust if UserResponse is defined
	}

	return response, nil
}

func (s *tokoService) UpdateToko(id uint, userID uint, req *dto.UpdateTokoRequest) (*dto.TokoWithUserResponse, error) {
	toko, err := s.tokoRepo.GetTokoByID(id)
	if err != nil {
		return nil, err
	}
	if toko.IDUser != userID {
		return nil, errors.New("access denied")
	}

	if req.NamaToko != "" {
		toko.NamaToko = req.NamaToko
	}
	if req.UrlFoto != "" {
		toko.UrlFoto = req.UrlFoto
	}
	toko.UpdatedAt = time.Now()

	err = s.tokoRepo.UpdateToko(toko)
	if err != nil {
		return nil, err
	}

	fetchedToko, err := s.tokoRepo.GetTokoByID(id)
	if err != nil {
		return nil, err
	}

	response := &dto.TokoWithUserResponse{
		ID:       fetchedToko.ID,
		NamaToko: fetchedToko.NamaToko,
		UrlFoto:  fetchedToko.UrlFoto,
		User:     dto.UserResponse{}, // Adjust if UserResponse is defined
	}

	return response, nil
}

func (s *tokoService) DeleteToko(id uint, userID uint) error {
	toko, err := s.tokoRepo.GetTokoByID(id)
	if err != nil {
		return err
	}
	if toko.IDUser != userID {
		return errors.New("access denied")
	}

	return s.tokoRepo.DeleteToko(id)
}
