package service

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"strconv"

	"github.com/wafiydev/okefin-service/internal/dto"
	"github.com/wafiydev/okefin-service/internal/models"
	"github.com/wafiydev/okefin-service/internal/repository"

)

type ProdukService interface {
	CreateProduk(userID uint, req *dto.CreateProdukRequest) (*dto.ProdukResponse, error)
	GetAllProdukByTokoID(userID uint, page, limit int) ([]dto.ProdukResponse, int, int64, error)
	GetProdukByID(userID uint, id uint) (*dto.ProdukResponse, error)
	UpdateProduk(userID uint, id uint, req *dto.UpdateProdukRequest) (*dto.ProdukResponse, error)
	DeleteProduk(userID uint, id uint) error
}

type produkService struct {
	produkRepo   repository.ProdukRepository
	tokoRepo     repository.TokoRepository
	fotoRepo     repository.FotoProdukRepository
	categoryRepo repository.CategoryRepository
}

func NewProdukService(produkRepo repository.ProdukRepository, tokoRepo repository.TokoRepository, fotoRepo repository.FotoProdukRepository, categoryRepo repository.CategoryRepository) ProdukService {
	return &produkService{produkRepo: produkRepo, tokoRepo: tokoRepo, fotoRepo: fotoRepo, categoryRepo: categoryRepo}
}

func (s *produkService) CreateProduk(userID uint, req *dto.CreateProdukRequest) (*dto.ProdukResponse, error) {
	toko, err := s.tokoRepo.GetTokoByUserID(userID)
	if err != nil {
		return nil, errors.New("toko not found")
	}

	if _, err := s.categoryRepo.GetCategoryByID(req.IDCategory); err != nil {
		return nil, errors.New("invalid category ID")
	}

	hargaStr := fmt.Sprintf("%.2f", req.Harga)

	produk := &models.Produk{
		NamaProduk:    req.NamaProduk,
		Slug:          generateSlug(req.NamaProduk),
		HargaReseller: hargaStr,
		HargaKonsumen: hargaStr,
		Stok:          req.Stok,
		Deskripsi:     req.Deskripsi,
		IDToko:        toko.ID,
		IDCategory:    req.IDCategory,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	err = s.produkRepo.CreateProduk(produk)
	if err != nil {
		return nil, fmt.Errorf("failed to create produk: %v", err)
	}

	for _, url := range req.UrlFotos {
		foto := &models.FotoProduk{
			IDProduk:  produk.ID,
			Url:       url,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := s.fotoRepo.CreateFotoProduk(foto); err != nil {
			return nil, fmt.Errorf("failed to create foto produk: %v", err)
		}
	}

	log := &models.LogProduk{
		IDProduk:      produk.ID,
		NamaProduk:    produk.NamaProduk,
		Slug:          produk.Slug,
		HargaReseller: produk.HargaReseller,
		HargaKonsumen: produk.HargaKonsumen,
		Deskripsi:     produk.Deskripsi,
		IDToko:        produk.IDToko,
		IDCategory:    produk.IDCategory,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	if err := s.produkRepo.CreateLogProduk(log); err != nil {
		return nil, fmt.Errorf("failed to log produk creation: %v", err)
	}

	fotos, err := s.fotoRepo.GetFotoProdukByProdukID(produk.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch fotos: %v", err)
	}
	var fotoUrls []string
	for _, f := range fotos {
		fotoUrls = append(fotoUrls, f.Url)
	}

	response := &dto.ProdukResponse{
		ID:         produk.ID,
		NamaProduk: produk.NamaProduk,
		Harga:      req.Harga,
		Stok:       produk.Stok,
		IDCategory: produk.IDCategory,
		UrlFotos:   fotoUrls,
		Deskripsi:  produk.Deskripsi,
		TokoID:     produk.IDToko,
		CreatedAt:  produk.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  produk.UpdatedAt.Format(time.RFC3339),
	}
	return response, nil
}

func (s *produkService) GetAllProdukByTokoID(userID uint, page, limit int) ([]dto.ProdukResponse, int, int64, error) {
	toko, err := s.tokoRepo.GetTokoByUserID(userID)
	if err != nil {
		return nil, 0, 0, errors.New("toko not found")
	}

	offset := (page - 1) * limit
	produkList, total, err := s.produkRepo.GetAllProdukByTokoID(toko.ID, limit, offset)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to fetch produk list: %v", err)
	}

	var responses []dto.ProdukResponse
	for _, p := range produkList {
		fotos, err := s.fotoRepo.GetFotoProdukByProdukID(p.ID)
		if err != nil {
			return nil, 0, 0, fmt.Errorf("failed to fetch fotos: %v", err)
		}
		var fotoUrls []string
		for _, f := range fotos {
			fotoUrls = append(fotoUrls, f.Url)
		}

		response := dto.ProdukResponse{
			ID:         p.ID,
			NamaProduk: p.NamaProduk,
			Harga:      parsePrice(p.HargaKonsumen),
			Stok:       p.Stok,
			IDCategory: p.IDCategory,
			UrlFotos:   fotoUrls,
			Deskripsi:  p.Deskripsi,
			TokoID:     p.IDToko,
			CreatedAt:  p.CreatedAt.Format(time.RFC3339),
			UpdatedAt:  p.UpdatedAt.Format(time.RFC3339),
		}
		responses = append(responses, response)
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))
	return responses, totalPages, total, nil
}

func (s *produkService) GetProdukByID(userID uint, id uint) (*dto.ProdukResponse, error) {
	produk, err := s.produkRepo.GetProdukByID(id, userID)
	if err != nil {
		return nil, errors.New("produk not found or access denied")
	}

	fotos, err := s.fotoRepo.GetFotoProdukByProdukID(produk.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch fotos: %v", err)
	}
	var fotoUrls []string
	for _, f := range fotos {
		fotoUrls = append(fotoUrls, f.Url)
	}

	response := &dto.ProdukResponse{
		ID:         produk.ID,
		NamaProduk: produk.NamaProduk,
		Harga:      parsePrice(produk.HargaKonsumen),
		Stok:       produk.Stok,
		IDCategory: produk.IDCategory,
		UrlFotos:   fotoUrls,
		Deskripsi:  produk.Deskripsi,
		TokoID:     produk.IDToko,
		CreatedAt:  produk.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  produk.UpdatedAt.Format(time.RFC3339),
	}
	return response, nil
}

func (s *produkService) UpdateProduk(userID uint, id uint, req *dto.UpdateProdukRequest) (*dto.ProdukResponse, error) {
	produk, err := s.produkRepo.GetProdukByID(id, userID)
	if err != nil {
		return nil, errors.New("produk not found or access denied")
	}

	if req.NamaProduk != "" {
		produk.NamaProduk = req.NamaProduk
		produk.Slug = generateSlug(req.NamaProduk)
	}
	if req.Harga > 0 {
		produk.HargaReseller = fmt.Sprintf("%.2f", req.Harga)
		produk.HargaKonsumen = fmt.Sprintf("%.2f", req.Harga)
	}
	if req.Stok > 0 {
		produk.Stok = req.Stok
	}
	if req.IDCategory > 0 {
		if _, err := s.categoryRepo.GetCategoryByID(req.IDCategory); err != nil {
			return nil, errors.New("invalid category ID")
		}
		produk.IDCategory = req.IDCategory
	}
	if len(req.UrlFotos) > 0 {
		if err := s.fotoRepo.DeleteFotoProdukByProdukID(produk.ID); err != nil {
			return nil, fmt.Errorf("failed to delete existing fotos: %v", err)
		}
		for _, url := range req.UrlFotos {
			foto := &models.FotoProduk{
				IDProduk:  produk.ID,
				Url:       url,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			if err := s.fotoRepo.CreateFotoProduk(foto); err != nil {
				return nil, fmt.Errorf("failed to update foto produk: %v", err)
			}
		}
	}
	if req.Deskripsi != "" {
		produk.Deskripsi = req.Deskripsi
	}
	produk.UpdatedAt = time.Now()

	err = s.produkRepo.UpdateProduk(produk)
	if err != nil {
		return nil, fmt.Errorf("failed to update produk: %v", err)
	}

	log := &models.LogProduk{
		IDProduk:      produk.ID,
		NamaProduk:    produk.NamaProduk,
		Slug:          produk.Slug,
		HargaReseller: produk.HargaReseller,
		HargaKonsumen: produk.HargaKonsumen,
		Deskripsi:     produk.Deskripsi,
		IDToko:        produk.IDToko,
		IDCategory:    produk.IDCategory,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	if err := s.produkRepo.CreateLogProduk(log); err != nil {
		return nil, fmt.Errorf("failed to log produk update: %v", err)
	}

	fotos, err := s.fotoRepo.GetFotoProdukByProdukID(produk.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch fotos: %v", err)
	}
	var fotoUrls []string
	for _, f := range fotos {
		fotoUrls = append(fotoUrls, f.Url)
	}

	response := &dto.ProdukResponse{
		ID:         produk.ID,
		NamaProduk: produk.NamaProduk,
		Harga:      req.Harga,
		Stok:       produk.Stok,
		IDCategory: produk.IDCategory,
		UrlFotos:   fotoUrls,
		Deskripsi:  produk.Deskripsi,
		TokoID:     produk.IDToko,
		CreatedAt:  produk.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  produk.UpdatedAt.Format(time.RFC3339),
	}
	return response, nil
}

func (s *produkService) DeleteProduk(userID uint, id uint) error {
	if !s.produkRepo.CheckProdukOwnership(id, userID) {
		return errors.New("access denied")
	}

	produk, err := s.produkRepo.GetProdukByID(id, userID)
	if err == nil {
		log := &models.LogProduk{
			IDProduk:      produk.ID,
			NamaProduk:    produk.NamaProduk,
			Slug:          produk.Slug,
			HargaReseller: produk.HargaReseller,
			HargaKonsumen: produk.HargaKonsumen,
			Deskripsi:     produk.Deskripsi,
			IDToko:        produk.IDToko,
			IDCategory:    produk.IDCategory,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
		if err := s.produkRepo.CreateLogProduk(log); err != nil {
			return fmt.Errorf("failed to log produk deletion: %v", err)
		}
	}

	if err := s.fotoRepo.DeleteFotoProdukByProdukID(id); err != nil {
		return fmt.Errorf("failed to delete fotos: %v", err)
	}

	return s.produkRepo.DeleteProduk(id)
}

func generateSlug(name string) string {
	return strings.ToLower(strings.ReplaceAll(name, " ", "-"))
}

func parsePrice(priceStr string) float64 {
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return 0.0
	}
	return price
}