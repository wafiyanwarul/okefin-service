package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/wafiydev/okefin-service/internal/dto"
	"github.com/wafiydev/okefin-service/internal/models"
	"github.com/wafiydev/okefin-service/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(req *dto.RegisterRequest) error
	Login(req *dto.LoginRequest) (*dto.LoginResponse, error)
	GetProvinsiByID(id string) (*dto.ProvinsiResponse, error)
	GetKotaByID(id string) (*dto.KotaResponse, error)
}

type authService struct {
	authRepo repository.AuthRepository
}

func NewAuthService(authRepo repository.AuthRepository) AuthService {
	return &authService{authRepo: authRepo}
}

func (s *authService) Register(req *dto.RegisterRequest) error {
	// Check if email already exists
	_, err := s.authRepo.GetUserByEmail(req.Email)
	if err == nil {
		return errors.New("email already exists")
	}

	// Check if no_telp already exists
	_, err = s.authRepo.GetUserByNoTelp(req.NoTelp)
	if err == nil {
		return errors.New("no_telp already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.KataSandi), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Create user
	user := &models.User{
		Nama:         req.Nama,
		KataSandi:    string(hashedPassword),
		NoTelp:       req.NoTelp,
		TanggalLahir: req.TanggalLahir,
		Pekerjaan:    req.Pekerjaan,
		Email:        req.Email,
		IDProvinsi:   req.IDProvinsi,
		IDKota:       req.IDKota,
		IsAdmin:      false,
	}

	err = s.authRepo.CreateUser(user)
	if err != nil {
		return err
	}

	// Create toko automatically
	toko := &models.Toko{
		IDUser:   user.ID,
		NamaToko: fmt.Sprintf("Toko %s", user.Nama),
		UrlFoto:  "",
	}

	return s.authRepo.CreateToko(toko)
}

func (s *authService) Login(req *dto.LoginRequest) (*dto.LoginResponse, error) {
	// Get user by no_telp
	user, err := s.authRepo.GetUserByNoTelp(req.NoTelp)
	if err != nil {
		return nil, errors.New("No Telp atau kata sandi salah")
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.KataSandi), []byte(req.KataSandi))
	if err != nil {
		return nil, errors.New("No Telp atau kata sandi salah")
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    strconv.Itoa(int(user.ID)),
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, err
	}

	// Get provinsi and kota data
	provinsi, _ := s.GetProvinsiByID(user.IDProvinsi)
	kota, _ := s.GetKotaByID(user.IDKota)

	response := &dto.LoginResponse{
		Nama:         user.Nama,
		NoTelp:       user.NoTelp,
		TanggalLahir: user.TanggalLahir,
		Tentang:      user.Tentang,
		Pekerjaan:    user.Pekerjaan,
		Email:        user.Email,
		IDProvinsi:   provinsi,
		IDKota:       kota,
		Token:        tokenString,
	}

	return response, nil
}

func (s *authService) GetProvinsiByID(id string) (*dto.ProvinsiResponse, error) {
	url := fmt.Sprintf("https://www.emsifa.com/api-wilayah-indonesia/api/province/%s.json", id)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var provinsi dto.ProvinsiResponse
	err = json.NewDecoder(resp.Body).Decode(&provinsi)
	if err != nil {
		return nil, err
	}

	return &provinsi, nil
}

func (s *authService) GetKotaByID(id string) (*dto.KotaResponse, error) {
	url := fmt.Sprintf("https://www.emsifa.com/api-wilayah-indonesia/api/regency/%s.json", id)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var kota dto.KotaResponse
	err = json.NewDecoder(resp.Body).Decode(&kota)
	if err != nil {
		return nil, err
	}

	return &kota, nil
}
