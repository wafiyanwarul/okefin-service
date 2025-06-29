package service

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/wafiydev/okefin-service/internal/dto"
	"github.com/wafiydev/okefin-service/internal/repository"
)

type UserService interface {
	GetUserProfile(userID uint) (*dto.UserResponse, error)
	UpdateUserProfile(userID uint, req *dto.UpdateUserRequest) (*dto.UserResponse, error)
	UploadFile(file *multipart.FileHeader) (*dto.UploadResponse, error)
	GetAllUsers(page, limit int) ([]dto.UserResponse, int, int64, error)
}

type userService struct {
	userRepo    repository.UserRepository
	authService AuthService
}

func NewUserService(userRepo repository.UserRepository, authService AuthService) UserService {
	return &userService{
		userRepo:    userRepo,
		authService: authService,
	}
}

func (s *userService) GetUserProfile(userID uint) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Get provinsi and kota data
	provinsi, _ := s.authService.GetProvinsiByID(user.IDProvinsi)
	kota, _ := s.authService.GetKotaByID(user.IDKota)

	response := &dto.UserResponse{
		ID:           user.ID,
		Nama:         user.Nama,
		NoTelp:       user.NoTelp,
		TanggalLahir: user.TanggalLahir,
		JenisKelamin: user.JenisKelamin,
		Tentang:      user.Tentang,
		Pekerjaan:    user.Pekerjaan,
		Email:        user.Email,
		IDProvinsi:   provinsi,
		IDKota:       kota,
		IsAdmin:      user.IsAdmin,
	}

	return response, nil
}

func (s *userService) UpdateUserProfile(userID uint, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Check if email is being changed and if it already exists
	if req.Email != "" && req.Email != user.Email {
		if s.userRepo.CheckEmailExists(req.Email, userID) {
			return nil, errors.New("email already exists")
		}
		user.Email = req.Email
	}

	// Update fields
	if req.Nama != "" {
		user.Nama = req.Nama
	}
	if req.TanggalLahir != "" {
		user.TanggalLahir = req.TanggalLahir
	}
	if req.JenisKelamin != "" {
		user.JenisKelamin = req.JenisKelamin
	}
	if req.Tentang != "" {
		user.Tentang = req.Tentang
	}
	if req.Pekerjaan != "" {
		user.Pekerjaan = req.Pekerjaan
	}
	if req.IDProvinsi != "" {
		user.IDProvinsi = req.IDProvinsi
	}
	if req.IDKota != "" {
		user.IDKota = req.IDKota
	}

	err = s.userRepo.UpdateUser(user)
	if err != nil {
		return nil, err
	}

	return s.GetUserProfile(userID)
}

func (s *userService) UploadFile(file *multipart.FileHeader) (*dto.UploadResponse, error) {
	// Create uploads directory if it doesn't exist
	uploadDir := "uploads"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.MkdirAll(uploadDir, 0755)
	}

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	filepath := filepath.Join(uploadDir, filename)

	// Open uploaded file
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	// Create destination file
	dst, err := os.Create(filepath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	// Copy file content
	_, err = io.Copy(dst, src)
	if err != nil {
		return nil, err
	}

	// Return file URL
	fileURL := fmt.Sprintf("/uploads/%s", filename)
	return &dto.UploadResponse{URL: fileURL}, nil
}

func (s *userService) GetAllUsers(page, limit int) ([]dto.UserResponse, int, int64, error) {
	offset := (page - 1) * limit
	users, total, err := s.userRepo.GetAllUsers(limit, offset)
	if err != nil {
		return nil, 0, 0, err
	}

	var responses []dto.UserResponse
	for _, user := range users {
		// Get provinsi and kota data
		provinsi, _ := s.authService.GetProvinsiByID(user.IDProvinsi)
		kota, _ := s.authService.GetKotaByID(user.IDKota)

		response := dto.UserResponse{
			ID:           user.ID,
			Nama:         user.Nama,
			NoTelp:       user.NoTelp,
			TanggalLahir: user.TanggalLahir,
			JenisKelamin: user.JenisKelamin,
			Tentang:      user.Tentang,
			Pekerjaan:    user.Pekerjaan,
			Email:        user.Email,
			IDProvinsi:   provinsi,
			IDKota:       kota,
			IsAdmin:      user.IsAdmin,
		}
		responses = append(responses, response)
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))
	return responses, totalPages, total, nil
}
