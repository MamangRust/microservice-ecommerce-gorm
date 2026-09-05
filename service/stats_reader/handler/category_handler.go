package handler

import (
	"context"

	"github.com/MamangRust/microservice-ecommerce-grpc-stats-reader/repository"
	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
pbcatestats "github.com/MamangRust/microservice-ecommerce-grpc/pb/category"
	"go.uber.org/zap"
)

// CategoryStatsHandler serves CategoryStatsService, CategoryStatsByIdService
// and CategoryStatsByMerchantService from ClickHouse.
type CategoryStatsHandler struct {
	pbcatestats.UnimplementedCategoryStatsServiceServer
	pbcatestats.UnimplementedCategoryStatsByIdServiceServer
	pbcatestats.UnimplementedCategoryStatsByMerchantServiceServer
	repo repository.Repository
	log  logger.LoggerInterface
}

func NewCategoryStatsHandler(repo repository.Repository, log logger.LoggerInterface) *CategoryStatsHandler {
	return &CategoryStatsHandler{repo: repo, log: log}
}

// --- CategoryStatsService ---

func (h *CategoryStatsHandler) FindMonthlyTotalPrices(ctx context.Context, req *pbcatestats.FindYearMonthTotalPrices) (*pbcatestats.ApiResponseCategoryMonthlyTotalPrice, error) {
	data, err := h.repo.GetMonthlyTotalPricing(ctx, int(req.GetYear()), int(req.GetMonth()), "", 0)
	if err != nil {
		h.log.Error("FindMonthlyTotalPrices failed", zap.Error(err))
		return nil, err
	}
	return &pbcatestats.ApiResponseCategoryMonthlyTotalPrice{
		Status:  "success",
		Message: "Monthly sales retrieved successfully",
		Data:    mapMonthlyPricing(data),
	}, nil
}

func (h *CategoryStatsHandler) FindYearlyTotalPrices(ctx context.Context, req *pbcatestats.FindYearTotalPrices) (*pbcatestats.ApiResponseCategoryYearlyTotalPrice, error) {
	data, err := h.repo.GetYearlyTotalPricing(ctx, int(req.GetYear()), "", 0)
	if err != nil {
		h.log.Error("FindYearlyTotalPrices failed", zap.Error(err))
		return nil, err
	}
	return &pbcatestats.ApiResponseCategoryYearlyTotalPrice{
		Status:  "success",
		Message: "Yearly payment methods retrieved successfully",
		Data:    mapYearlyPricing(data),
	}, nil
}

func (h *CategoryStatsHandler) FindMonthPrice(ctx context.Context, req *pbcatestats.FindYearCategory) (*pbcatestats.ApiResponseCategoryMonthPrice, error) {
	data, err := h.repo.GetMonthlyCategoryStats(ctx, int(req.GetYear()), "", 0)
	if err != nil {
		h.log.Error("FindMonthPrice failed", zap.Error(err))
		return nil, err
	}
	return &pbcatestats.ApiResponseCategoryMonthPrice{
		Status:  "success",
		Message: "Monthly payment methods retrieved successfully",
		Data:    mapMonthlyCategory(data),
	}, nil
}

func (h *CategoryStatsHandler) FindYearPrice(ctx context.Context, req *pbcatestats.FindYearCategory) (*pbcatestats.ApiResponseCategoryYearPrice, error) {
	data, err := h.repo.GetYearlyCategoryStats(ctx, int(req.GetYear()), "", 0)
	if err != nil {
		h.log.Error("FindYearPrice failed", zap.Error(err))
		return nil, err
	}
	return &pbcatestats.ApiResponseCategoryYearPrice{
		Status:  "success",
		Message: "Yearly payment methods retrieved successfully",
		Data:    mapYearlyCategory(data),
	}, nil
}

// --- CategoryStatsByIdService ---

func (h *CategoryStatsHandler) FindMonthlyTotalPricesById(ctx context.Context, req *pbcatestats.FindYearMonthTotalPriceById) (*pbcatestats.ApiResponseCategoryMonthlyTotalPrice, error) {
	data, err := h.repo.GetMonthlyTotalPricing(ctx, int(req.GetYear()), int(req.GetMonth()), "category_id", req.GetCategoryId())
	if err != nil {
		h.log.Error("FindMonthlyTotalPricesById failed", zap.Error(err))
		return nil, err
	}
	return &pbcatestats.ApiResponseCategoryMonthlyTotalPrice{
		Status:  "success",
		Message: "Monthly sales retrieved successfully",
		Data:    mapMonthlyPricing(data),
	}, nil
}

func (h *CategoryStatsHandler) FindYearlyTotalPricesById(ctx context.Context, req *pbcatestats.FindYearTotalPriceById) (*pbcatestats.ApiResponseCategoryYearlyTotalPrice, error) {
	data, err := h.repo.GetYearlyTotalPricing(ctx, int(req.GetYear()), "category_id", req.GetCategoryId())
	if err != nil {
		h.log.Error("FindYearlyTotalPricesById failed", zap.Error(err))
		return nil, err
	}
	return &pbcatestats.ApiResponseCategoryYearlyTotalPrice{
		Status:  "success",
		Message: "Yearly payment methods retrieved successfully",
		Data:    mapYearlyPricing(data),
	}, nil
}

func (h *CategoryStatsHandler) FindMonthPriceById(ctx context.Context, req *pbcatestats.FindYearCategoryById) (*pbcatestats.ApiResponseCategoryMonthPrice, error) {
	data, err := h.repo.GetMonthlyCategoryStats(ctx, int(req.GetYear()), "category_id", req.GetCategoryId())
	if err != nil {
		h.log.Error("FindMonthPriceById failed", zap.Error(err))
		return nil, err
	}
	return &pbcatestats.ApiResponseCategoryMonthPrice{
		Status:  "success",
		Message: "Monthly payment methods retrieved successfully",
		Data:    mapMonthlyCategory(data),
	}, nil
}

func (h *CategoryStatsHandler) FindYearPriceById(ctx context.Context, req *pbcatestats.FindYearCategoryById) (*pbcatestats.ApiResponseCategoryYearPrice, error) {
	data, err := h.repo.GetYearlyCategoryStats(ctx, int(req.GetYear()), "category_id", req.GetCategoryId())
	if err != nil {
		h.log.Error("FindYearPriceById failed", zap.Error(err))
		return nil, err
	}
	return &pbcatestats.ApiResponseCategoryYearPrice{
		Status:  "success",
		Message: "Yearly payment methods retrieved successfully",
		Data:    mapYearlyCategory(data),
	}, nil
}

// --- CategoryStatsByMerchantService ---

func (h *CategoryStatsHandler) FindMonthlyTotalPricesByMerchant(ctx context.Context, req *pbcatestats.FindYearMonthTotalPriceByMerchant) (*pbcatestats.ApiResponseCategoryMonthlyTotalPrice, error) {
	data, err := h.repo.GetMonthlyTotalPricing(ctx, int(req.GetYear()), int(req.GetMonth()), "merchant_id", req.GetMerchantId())
	if err != nil {
		h.log.Error("FindMonthlyTotalPricesByMerchant failed", zap.Error(err))
		return nil, err
	}
	return &pbcatestats.ApiResponseCategoryMonthlyTotalPrice{
		Status:  "success",
		Message: "Monthly sales retrieved successfully",
		Data:    mapMonthlyPricing(data),
	}, nil
}

func (h *CategoryStatsHandler) FindYearlyTotalPricesByMerchant(ctx context.Context, req *pbcatestats.FindYearTotalPriceByMerchant) (*pbcatestats.ApiResponseCategoryYearlyTotalPrice, error) {
	data, err := h.repo.GetYearlyTotalPricing(ctx, int(req.GetYear()), "merchant_id", req.GetMerchantId())
	if err != nil {
		h.log.Error("FindYearlyTotalPricesByMerchant failed", zap.Error(err))
		return nil, err
	}
	return &pbcatestats.ApiResponseCategoryYearlyTotalPrice{
		Status:  "success",
		Message: "Yearly payment methods retrieved successfully",
		Data:    mapYearlyPricing(data),
	}, nil
}

func (h *CategoryStatsHandler) FindMonthPriceByMerchant(ctx context.Context, req *pbcatestats.FindYearCategoryByMerchant) (*pbcatestats.ApiResponseCategoryMonthPrice, error) {
	data, err := h.repo.GetMonthlyCategoryStats(ctx, int(req.GetYear()), "merchant_id", req.GetMerchantId())
	if err != nil {
		h.log.Error("FindMonthPriceByMerchant failed", zap.Error(err))
		return nil, err
	}
	return &pbcatestats.ApiResponseCategoryMonthPrice{
		Status:  "success",
		Message: "Monthly payment methods retrieved successfully",
		Data:    mapMonthlyCategory(data),
	}, nil
}

func (h *CategoryStatsHandler) FindYearPriceByMerchant(ctx context.Context, req *pbcatestats.FindYearCategoryByMerchant) (*pbcatestats.ApiResponseCategoryYearPrice, error) {
	data, err := h.repo.GetYearlyCategoryStats(ctx, int(req.GetYear()), "merchant_id", req.GetMerchantId())
	if err != nil {
		h.log.Error("FindYearPriceByMerchant failed", zap.Error(err))
		return nil, err
	}
	return &pbcatestats.ApiResponseCategoryYearPrice{
		Status:  "success",
		Message: "Yearly payment methods retrieved successfully",
		Data:    mapYearlyCategory(data),
	}, nil
}

// --- Mappers ---

func mapMonthlyPricing(data []repository.MonthlyRevenue) []*pbcatestats.CategoriesMonthlyTotalPriceResponse {
	var out []*pbcatestats.CategoriesMonthlyTotalPriceResponse
	for _, d := range data {
		out = append(out, &pbcatestats.CategoriesMonthlyTotalPriceResponse{
			Year:         d.Year,
			Month:        d.Month,
			TotalRevenue: int32(d.TotalRevenue),
		})
	}
	return out
}

func mapYearlyPricing(data []repository.YearlyRevenue) []*pbcatestats.CategoriesYearlyTotalPriceResponse {
	var out []*pbcatestats.CategoriesYearlyTotalPriceResponse
	for _, d := range data {
		out = append(out, &pbcatestats.CategoriesYearlyTotalPriceResponse{
			Year:         d.Year,
			TotalRevenue: int32(d.TotalRevenue),
		})
	}
	return out
}

func mapMonthlyCategory(data []repository.MonthlyCategory) []*pbcatestats.CategoryMonthPriceResponse {
	var out []*pbcatestats.CategoryMonthPriceResponse
	for _, d := range data {
		out = append(out, &pbcatestats.CategoryMonthPriceResponse{
			Month:        d.Month,
			CategoryId:   d.CategoryID,
			CategoryName: d.CategoryName,
			OrderCount:   int32(d.OrderCount),
			ItemsSold:    int32(d.ItemsSold),
			TotalRevenue: int32(d.TotalRevenue),
		})
	}
	return out
}

func mapYearlyCategory(data []repository.YearlyCategory) []*pbcatestats.CategoryYearPriceResponse {
	var out []*pbcatestats.CategoryYearPriceResponse
	for _, d := range data {
		out = append(out, &pbcatestats.CategoryYearPriceResponse{
			Year:               d.Year,
			CategoryId:         d.CategoryID,
			CategoryName:       d.CategoryName,
			OrderCount:         int32(d.OrderCount),
			ItemsSold:          int32(d.ItemsSold),
			TotalRevenue:       int32(d.TotalRevenue),
			UniqueProductsSold: int32(d.UniqueProductsSold),
		})
	}
	return out
}
