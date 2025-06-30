package dto

type CreateCategoryRequest struct {
	NamaCategory string `json:"nama_category" validate:"required"`
}

type UpdateCategoryRequest struct {
	NamaCategory string `json:"nama_category"`
}

type CategoryResponse struct {
	ID           uint   `json:"id"`
	NamaCategory string `json:"nama_category"`
}