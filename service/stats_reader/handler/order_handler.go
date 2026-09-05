package handler

import (
	"context"

	"github.com/MamangRust/microservice-ecommerce-grpc-stats-reader/repository"
	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	pborderstats "github.com/MamangRust/microservice-ecommerce-grpc/pb/order"
	"go.uber.org/zap"
)

// OrderStatsHandler serves OrderStatsService (revenue + order aggregates) from
// ClickHouse.
type OrderStatsHandler struct {
	pborderstats.UnimplementedOrderStatsServiceServer
	repo repository.Repository
	log  logger.LoggerInterface
}

func NewOrderStatsHandler(repo repository.Repository, log logger.LoggerInterface) *OrderStatsHandler {
	return &OrderStatsHandler{repo: repo, log: log}
}

func (h *OrderStatsHandler) FindMonthlyTotalRevenue(ctx context.Context, req *pborderstats.FindYearMonthTotalRevenue) (*pborderstats.ApiResponseOrderMonthlyTotalRevenue, error) {
	data, err := h.repo.GetMonthlyTotalRevenue(ctx, int(req.GetYear()), int(req.GetMonth()), 0)
	if err != nil {
		h.log.Error("FindMonthlyTotalRevenue failed", zap.Error(err))
		return nil, err
	}
	return &pborderstats.ApiResponseOrderMonthlyTotalRevenue{
		Status:  "success",
		Message: "Monthly sales retrieved successfully",
		Data:    mapMonthlyRevenue(data),
	}, nil
}

func (h *OrderStatsHandler) FindYearlyTotalRevenue(ctx context.Context, req *pborderstats.FindYearTotalRevenue) (*pborderstats.ApiResponseOrderYearlyTotalRevenue, error) {
	data, err := h.repo.GetYearlyTotalRevenue(ctx, int(req.GetYear()), 0)
	if err != nil {
		h.log.Error("FindYearlyTotalRevenue failed", zap.Error(err))
		return nil, err
	}
	return &pborderstats.ApiResponseOrderYearlyTotalRevenue{
		Status:  "success",
		Message: "Yearly payment methods retrieved successfully",
		Data:    mapYearlyRevenue(data),
	}, nil
}

func (h *OrderStatsHandler) FindMonthlyRevenue(ctx context.Context, req *pborderstats.FindYearOrder) (*pborderstats.ApiResponseOrderMonthly, error) {
	data, err := h.repo.GetMonthlyOrderStats(ctx, int(req.GetYear()), 0)
	if err != nil {
		h.log.Error("FindMonthlyRevenue failed", zap.Error(err))
		return nil, err
	}
	return &pborderstats.ApiResponseOrderMonthly{
		Status:  "success",
		Message: "Monthly revenue data retrieved",
		Data:    mapMonthlyOrder(data),
	}, nil
}

func (h *OrderStatsHandler) FindYearlyRevenue(ctx context.Context, req *pborderstats.FindYearOrder) (*pborderstats.ApiResponseOrderYearly, error) {
	data, err := h.repo.GetYearlyOrderStats(ctx, int(req.GetYear()), 0)
	if err != nil {
		h.log.Error("FindYearlyRevenue failed", zap.Error(err))
		return nil, err
	}
	return &pborderstats.ApiResponseOrderYearly{
		Status:  "success",
		Message: "Yearly revenue data retrieved",
		Data:    mapYearlyOrder(data),
	}, nil
}

func (h *OrderStatsHandler) FindMonthlyTotalRevenueByMerchant(ctx context.Context, req *pborderstats.FindYearMonthTotalRevenueByMerchant) (*pborderstats.ApiResponseOrderMonthlyTotalRevenue, error) {
	data, err := h.repo.GetMonthlyTotalRevenue(ctx, int(req.GetYear()), int(req.GetMonth()), req.GetMerchantId())
	if err != nil {
		h.log.Error("FindMonthlyTotalRevenueByMerchant failed", zap.Error(err))
		return nil, err
	}
	return &pborderstats.ApiResponseOrderMonthlyTotalRevenue{
		Status:  "success",
		Message: "Monthly sales retrieved successfully",
		Data:    mapMonthlyRevenue(data),
	}, nil
}

func (h *OrderStatsHandler) FindYearlyTotalRevenueByMerchant(ctx context.Context, req *pborderstats.FindYearTotalRevenueByMerchant) (*pborderstats.ApiResponseOrderYearlyTotalRevenue, error) {
	data, err := h.repo.GetYearlyTotalRevenue(ctx, int(req.GetYear()), req.GetMerchantId())
	if err != nil {
		h.log.Error("FindYearlyTotalRevenueByMerchant failed", zap.Error(err))
		return nil, err
	}
	return &pborderstats.ApiResponseOrderYearlyTotalRevenue{
		Status:  "success",
		Message: "Yearly payment methods retrieved successfully",
		Data:    mapYearlyRevenue(data),
	}, nil
}

func (h *OrderStatsHandler) FindMonthlyRevenueByMerchant(ctx context.Context, req *pborderstats.FindYearOrderByMerchant) (*pborderstats.ApiResponseOrderMonthly, error) {
	data, err := h.repo.GetMonthlyOrderStats(ctx, int(req.GetYear()), req.GetMerchantId())
	if err != nil {
		h.log.Error("FindMonthlyRevenueByMerchant failed", zap.Error(err))
		return nil, err
	}
	return &pborderstats.ApiResponseOrderMonthly{
		Status:  "success",
		Message: "Monthly revenue by merchant data retrieved",
		Data:    mapMonthlyOrder(data),
	}, nil
}

func (h *OrderStatsHandler) FindYearlyRevenueByMerchant(ctx context.Context, req *pborderstats.FindYearOrderByMerchant) (*pborderstats.ApiResponseOrderYearly, error) {
	data, err := h.repo.GetYearlyOrderStats(ctx, int(req.GetYear()), req.GetMerchantId())
	if err != nil {
		h.log.Error("FindYearlyRevenueByMerchant failed", zap.Error(err))
		return nil, err
	}
	return &pborderstats.ApiResponseOrderYearly{
		Status:  "success",
		Message: "Yearly revenue by merchant data retrieved",
		Data:    mapYearlyOrder(data),
	}, nil
}

// --- Mappers ---

func mapMonthlyRevenue(data []repository.MonthlyRevenue) []*pborderstats.OrderMonthlyTotalRevenueResponse {
	var out []*pborderstats.OrderMonthlyTotalRevenueResponse
	for _, d := range data {
		out = append(out, &pborderstats.OrderMonthlyTotalRevenueResponse{
			Year:         d.Year,
			Month:        d.Month,
			TotalRevenue: int32(d.TotalRevenue),
		})
	}
	return out
}

func mapYearlyRevenue(data []repository.YearlyRevenue) []*pborderstats.OrderYearlyTotalRevenueResponse {
	var out []*pborderstats.OrderYearlyTotalRevenueResponse
	for _, d := range data {
		out = append(out, &pborderstats.OrderYearlyTotalRevenueResponse{
			Year:         d.Year,
			TotalRevenue: int32(d.TotalRevenue),
		})
	}
	return out
}

func mapMonthlyOrder(data []repository.MonthlyOrder) []*pborderstats.OrderMonthlyResponse {
	var out []*pborderstats.OrderMonthlyResponse
	for _, d := range data {
		out = append(out, &pborderstats.OrderMonthlyResponse{
			Month:          d.Month,
			OrderCount:     int32(d.OrderCount),
			TotalRevenue:   int32(d.TotalRevenue),
			TotalItemsSold: int32(d.TotalItemsSold),
		})
	}
	return out
}

func mapYearlyOrder(data []repository.YearlyOrder) []*pborderstats.OrderYearlyResponse {
	var out []*pborderstats.OrderYearlyResponse
	for _, d := range data {
		out = append(out, &pborderstats.OrderYearlyResponse{
			Year:               d.Year,
			OrderCount:         int32(d.OrderCount),
			TotalRevenue:       int32(d.TotalRevenue),
			TotalItemsSold:     int32(d.TotalItemsSold),
			UniqueProductsSold: int32(d.UniqueProductsSold),
		})
	}
	return out
}
