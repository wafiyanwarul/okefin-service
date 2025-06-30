package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/wafiydev/okefin-service/internal/dto"
	"github.com/wafiydev/okefin-service/internal/models"
	"github.com/wafiydev/okefin-service/internal/repository"
)

type TransaksiService interface {
	CreateTransaksi(userID uint, req *dto.CreateTransaksiRequest, details []dto.DetailTransaksiRequest) (*dto.TransaksiResponse, error)
	GetAllTransaksiByUserID(userID uint, page, limit int) ([]dto.TransaksiResponse, int, int64, error)
	GetTransaksiByID(userID uint, id uint) (*dto.TransaksiResponse, error)
	UpdateTransaksi(userID uint, id uint, req *dto.UpdateTransaksiRequest) (*dto.TransaksiResponse, error)
}

type transaksiService struct {
	transaksiRepo repository.TransaksiRepository
	produkRepo    repository.ProdukRepository
}

func NewTransaksiService(transaksiRepo repository.TransaksiRepository, produkRepo repository.ProdukRepository) TransaksiService {
	return &transaksiService{transaksiRepo: transaksiRepo, produkRepo: produkRepo}
}

func (s *transaksiService) CreateTransaksi(userID uint, req *dto.CreateTransaksiRequest, details []dto.DetailTransaksiRequest) (*dto.TransaksiResponse, error) {
	// Validate total harga matches sum of detail harga
	var calculatedTotal float64
	for _, d := range details {
		calculatedTotal += d.Harga * float64(d.Jumlah)
	}
	if calculatedTotal != req.TotalHarga {
		return nil, errors.New("total harga does not match sum of detail harga")
	}

	// Create transaction
	transaksi := &models.Trx{
		IDUser:           userID,
		AlamatPengiriman: req.AlamatID,
		HargaTotal:       int(req.TotalHarga), // Convert float64 to int
		KodeInvoice:      generateInvoiceCode(),
		// Status field not present, using MethodBayar as placeholder
		MethodBayar: "pending", // Default status
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := s.transaksiRepo.CreateTransaksi(transaksi)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaksi: %v", err)
	}

	// Create detail and log for each product
	var logProduks []models.LogProduk
	for _, detailReq := range details {
		produk, err := s.produkRepo.GetProdukByID(detailReq.ProdukID, userID)
		if err != nil {
			return nil, errors.New("invalid produk ID")
		}
		if produk.Stok < detailReq.Jumlah {
			return nil, errors.New("insufficient stock for produk ID " + fmt.Sprint(detailReq.ProdukID))
		}

		// Create LogProduk for each item
		logProduk := &models.LogProduk{
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
		err = s.produkRepo.CreateLogProduk(logProduk)
		if err != nil {
			return nil, fmt.Errorf("failed to create log produk: %v", err)
		}
		logProduks = append(logProduks, *logProduk)

		// Create DetailTrx linked to the transaction and log
		detail := &models.DetailTrx{
			IDTrx:       transaksi.ID,
			IDLogProduk: logProduk.ID,
			IDToko:      produk.IDToko,
			Kuantitas:   detailReq.Jumlah,
			HargaTotal:  int(detailReq.Harga * float64(detailReq.Jumlah)), // Convert to int
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		err = s.transaksiRepo.CreateDetailTransaksi(detail)
		if err != nil {
			return nil, fmt.Errorf("failed to create detail transaksi: %v", err)
		}
	}

	// Update product stock
	for _, detailReq := range details {
		produk, err := s.produkRepo.GetProdukByID(detailReq.ProdukID, userID)
		if err != nil {
			return nil, err
		}
		produk.Stok -= detailReq.Jumlah
		err = s.produkRepo.UpdateProduk(produk)
		if err != nil {
			return nil, fmt.Errorf("failed to update produk stock: %v", err)
		}
	}

	// Fetch details to build response
	detailsResp, err := s.transaksiRepo.GetDetailsByTransaksiID(transaksi.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch details: %v", err)
	}

	var detailRespList []dto.DetailTransaksiResp
	for _, d := range detailsResp {
		log, err := s.produkRepo.GetLogProdukByID(d.IDLogProduk)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch log produk: %v", err)
		}
		detailRespList = append(detailRespList, dto.DetailTransaksiResp{
			ID:       d.ID,
			ProdukID: log.IDProduk,
			Jumlah:   d.Kuantitas,
			Harga:    float64(d.HargaTotal) / float64(d.Kuantitas), // Convert int to float64
		})
	}

	response := &dto.TransaksiResponse{
		ID:         transaksi.ID,
		AlamatID:   transaksi.AlamatPengiriman,
		UserID:     transaksi.IDUser,
		TotalHarga: float64(transaksi.HargaTotal), // Convert int to float64
		// Status field not present, using MethodBayar as placeholder
		Status:    transaksi.MethodBayar,
		CreatedAt: transaksi.CreatedAt.Format(time.RFC3339),
		UpdatedAt: transaksi.UpdatedAt.Format(time.RFC3339),
		Details:   detailRespList,
	}
	return response, nil
}

func (s *transaksiService) GetAllTransaksiByUserID(userID uint, page, limit int) ([]dto.TransaksiResponse, int, int64, error) {
	offset := (page - 1) * limit
	transaksiList, total, err := s.transaksiRepo.GetAllTransaksiByUserID(userID, limit, offset)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to fetch transaksi list: %v", err)
	}

	var responses []dto.TransaksiResponse
	for _, t := range transaksiList {
		details, err := s.transaksiRepo.GetDetailsByTransaksiID(t.ID)
		if err != nil {
			return nil, 0, 0, fmt.Errorf("failed to fetch details: %v", err)
		}

		var detailRespList []dto.DetailTransaksiResp
		for _, d := range details {
			log, err := s.produkRepo.GetLogProdukByID(d.IDLogProduk)
			if err != nil {
				return nil, 0, 0, fmt.Errorf("failed to fetch log produk: %v", err)
			}
			detailRespList = append(detailRespList, dto.DetailTransaksiResp{
				ID:       d.ID,
				ProdukID: log.IDProduk,
				Jumlah:   d.Kuantitas,
				Harga:    float64(d.HargaTotal) / float64(d.Kuantitas), // Convert int to float64
			})
		}

		response := dto.TransaksiResponse{
			ID:         t.ID,
			AlamatID:   t.AlamatPengiriman,
			UserID:     t.IDUser,
			TotalHarga: float64(t.HargaTotal), // Convert int to float64
			// Status field not present, using MethodBayar as placeholder
			Status:    t.MethodBayar,
			CreatedAt: t.CreatedAt.Format(time.RFC3339),
			UpdatedAt: t.UpdatedAt.Format(time.RFC3339),
			Details:   detailRespList,
		}
		responses = append(responses, response)
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))
	return responses, totalPages, total, nil
}

func (s *transaksiService) GetTransaksiByID(userID uint, id uint) (*dto.TransaksiResponse, error) {
	transaksi, err := s.transaksiRepo.GetTransaksiByID(id, userID)
	if err != nil {
		return nil, errors.New("transaksi not found or access denied")
	}

	details, err := s.transaksiRepo.GetDetailsByTransaksiID(transaksi.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch details: %v", err)
	}

	var detailRespList []dto.DetailTransaksiResp
	for _, d := range details {
		log, err := s.produkRepo.GetLogProdukByID(d.IDLogProduk)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch log produk: %v", err)
		}
		detailRespList = append(detailRespList, dto.DetailTransaksiResp{
			ID:       d.ID,
			ProdukID: log.IDProduk,
			Jumlah:   d.Kuantitas,
			Harga:    float64(d.HargaTotal) / float64(d.Kuantitas), // Convert int to float64
		})
	}

	response := &dto.TransaksiResponse{
		ID:         transaksi.ID,
		AlamatID:   transaksi.AlamatPengiriman,
		UserID:     transaksi.IDUser,
		TotalHarga: float64(transaksi.HargaTotal), // Convert int to float64
		// Status field not present, using MethodBayar as placeholder
		Status:    transaksi.MethodBayar,
		CreatedAt: transaksi.CreatedAt.Format(time.RFC3339),
		UpdatedAt: transaksi.UpdatedAt.Format(time.RFC3339),
		Details:   detailRespList,
	}
	return response, nil
}

func (s *transaksiService) UpdateTransaksi(userID uint, id uint, req *dto.UpdateTransaksiRequest) (*dto.TransaksiResponse, error) {
	transaksi, err := s.transaksiRepo.GetTransaksiByID(id, userID)
	if err != nil {
		return nil, errors.New("transaksi not found or access denied")
	}

	if req.Status != "" {
		validStatuses := map[string]bool{"pending": true, "completed": true, "cancelled": true}
		if !validStatuses[req.Status] {
			return nil, errors.New("invalid status value")
		}
		// Use MethodBayar as placeholder for status
		transaksi.MethodBayar = req.Status
	}
	transaksi.UpdatedAt = time.Now()

	err = s.transaksiRepo.UpdateTransaksi(transaksi)
	if err != nil {
		return nil, fmt.Errorf("failed to update transaksi: %v", err)
	}

	details, err := s.transaksiRepo.GetDetailsByTransaksiID(transaksi.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch details: %v", err)
	}

	var detailRespList []dto.DetailTransaksiResp
	for _, d := range details {
		log, err := s.produkRepo.GetLogProdukByID(d.IDLogProduk)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch log produk: %v", err)
		}
		detailRespList = append(detailRespList, dto.DetailTransaksiResp{
			ID:       d.ID,
			ProdukID: log.IDProduk,
			Jumlah:   d.Kuantitas,
			Harga:    float64(d.HargaTotal) / float64(d.Kuantitas), // Convert int to float64
		})
	}

	response := &dto.TransaksiResponse{
		ID:         transaksi.ID,
		AlamatID:   transaksi.AlamatPengiriman,
		UserID:     transaksi.IDUser,
		TotalHarga: float64(transaksi.HargaTotal), // Convert int to float64
		// Status field not present, using MethodBayar as placeholder
		Status:    transaksi.MethodBayar,
		CreatedAt: transaksi.CreatedAt.Format(time.RFC3339),
		UpdatedAt: transaksi.UpdatedAt.Format(time.RFC3339),
		Details:   detailRespList,
	}
	return response, nil
}

func generateInvoiceCode() string {
	return fmt.Sprintf("INV-%d-%s", time.Now().Unix(), fmt.Sprintf("%06d", time.Now().Nanosecond())[0:6])
}
