package service

import (
	"errors"
	"time"

	"github.com/wafiydev/okefin-service/internal/dto"
	"github.com/wafiydev/okefin-service/internal/models"
	"github.com/wafiydev/okefin-service/internal/repository"
)

type AlamatService interface {
	CreateAlamat(userID uint, req *dto.CreateAlamatRequest) (*dto.AlamatResponse, error)
	GetMyAlamat(userID uint, page, limit int) ([]dto.AlamatResponse, int, int64, error)
	GetAlamatByID(alamatID, userID uint) (*dto.AlamatResponse, error)
	UpdateAlamat(alamatID, userID uint, req *dto.UpdateAlamatRequest) (*dto.AlamatResponse, error)
	DeleteAlamat(alamatID, userID uint) error
}

type alamatService struct {
	alamatRepo repository.AlamatRepository
}

func NewAlamatService(alamatRepo repository.AlamatRepository) AlamatService {
	return &alamatService{alamatRepo: alamatRepo}
}

func (s *alamatService) CreateAlamat(userID uint, req *dto.CreateAlamatRequest) (*dto.AlamatResponse, error) {
	alamat := &models.Alamat{
		IDUser:       userID,
		JudulAlamat:  req.JudulAlamat,
		NamaPenerima: req.NamaPenerima,
		NoTelp:       req.NoTelp,
		DetailAlamat: req.DetailAlamat,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err := s.alamatRepo.CreateAlamat(alamat)
	if err != nil {
		return nil, err
	}

	response := &dto.AlamatResponse{
		ID:           alamat.ID,
		JudulAlamat:  alamat.JudulAlamat,
		NamaPenerima: alamat.NamaPenerima,
		NoTelp:       alamat.NoTelp,
		DetailAlamat: alamat.DetailAlamat,
	}

	return response, nil
}

func (s *alamatService) GetMyAlamat(userID uint, page, limit int) ([]dto.AlamatResponse, int, int64, error) {
	offset := (page - 1) * limit
	alamat, total, err := s.alamatRepo.GetAlamatByUserID(userID, limit, offset)
	if err != nil {
		return nil, 0, 0, err
	}

	var responses []dto.AlamatResponse
	for _, a := range alamat {
		response := dto.AlamatResponse{
			ID:           a.ID,
			JudulAlamat:  a.JudulAlamat,
			NamaPenerima: a.NamaPenerima,
			NoTelp:       a.NoTelp,
			DetailAlamat: a.DetailAlamat,
		}
		responses = append(responses, response)
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))
	return responses, totalPages, total, nil
}

func (s *alamatService) GetAlamatByID(alamatID, userID uint) (*dto.AlamatResponse, error) {
	// Check ownership
	if !s.alamatRepo.CheckAlamatOwnership(alamatID, userID) {
		return nil, errors.New("alamat not found or access denied")
	}

	alamat, err := s.alamatRepo.GetAlamatByID(alamatID)
	if err != nil {
		return nil, errors.New("alamat not found")
	}

	response := &dto.AlamatResponse{
		ID:           alamat.ID,
		JudulAlamat:  alamat.JudulAlamat,
		NamaPenerima: alamat.NamaPenerima,
		NoTelp:       alamat.NoTelp,
		DetailAlamat: alamat.DetailAlamat,
	}

	return response, nil
}

func (s *alamatService) UpdateAlamat(alamatID, userID uint, req *dto.UpdateAlamatRequest) (*dto.AlamatResponse, error) {
	// Check ownership
	if !s.alamatRepo.CheckAlamatOwnership(alamatID, userID) {
		return nil, errors.New("alamat not found or access denied")
	}

	alamat, err := s.alamatRepo.GetAlamatByID(alamatID)
	if err != nil {
		return nil, errors.New("alamat not found")
	}

	// Update fields
	if req.JudulAlamat != "" {
		alamat.JudulAlamat = req.JudulAlamat
	}
	if req.NamaPenerima != "" {
		alamat.NamaPenerima = req.NamaPenerima
	}
	if req.NoTelp != "" {
		alamat.NoTelp = req.NoTelp
	}
	if req.DetailAlamat != "" {
		alamat.DetailAlamat = req.DetailAlamat
	}

	alamat.UpdatedAt = time.Now()

	err = s.alamatRepo.UpdateAlamat(alamat)
	if err != nil {
		return nil, err
	}

	response := &dto.AlamatResponse{
		ID:           alamat.ID,
		JudulAlamat:  alamat.JudulAlamat,
		NamaPenerima: alamat.NamaPenerima,
		NoTelp:       alamat.NoTelp,
		DetailAlamat: alamat.DetailAlamat,
	}

	return response, nil
}

func (s *alamatService) DeleteAlamat(alamatID, userID uint) error {
	// Check ownership
	if !s.alamatRepo.CheckAlamatOwnership(alamatID, userID) {
		return errors.New("alamat not found or access denied")
	}

	return s.alamatRepo.DeleteAlamat(alamatID)
}
