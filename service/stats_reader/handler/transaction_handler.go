package handler

import (
	"context"

	"github.com/MamangRust/microservice-ecommerce-grpc-stats-reader/repository"
	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
pbtxstats "github.com/MamangRust/microservice-ecommerce-grpc/pb/transaction"
	"go.uber.org/zap"
)

// TransactionStatsHandler serves TransactionStatsService and
// TransactionStatsByMerchantService from ClickHouse.
type TransactionStatsHandler struct {
	pbtxstats.UnimplementedTransactionStatsServiceServer
	pbtxstats.UnimplementedTransactionStatsByMerchantServiceServer
	repo repository.Repository
	log  logger.LoggerInterface
}

func NewTransactionStatsHandler(repo repository.Repository, log logger.LoggerInterface) *TransactionStatsHandler {
	return &TransactionStatsHandler{repo: repo, log: log}
}

// --- TransactionStatsService ---

func (h *TransactionStatsHandler) GetMonthlyAmountSuccess(ctx context.Context, req *pbtxstats.MonthAmountTransactionRequest) (*pbtxstats.ApiResponseTransactionMonthAmountSuccess, error) {
	data, err := h.repo.GetMonthlyAmount(ctx, int(req.GetYear()), int(req.GetMonth()), "success", 0)
	if err != nil {
		h.log.Error("GetMonthlyAmountSuccess failed", zap.Error(err))
		return nil, err
	}
	return &pbtxstats.ApiResponseTransactionMonthAmountSuccess{
		Status:  "success",
		Message: "Successfully fetched monthly amount success stats",
		Data:    mapMonthlyAmountSuccess(data),
	}, nil
}

func (h *TransactionStatsHandler) GetYearlyAmountSuccess(ctx context.Context, req *pbtxstats.YearAmountTransactionRequest) (*pbtxstats.ApiResponseTransactionYearAmountSuccess, error) {
	data, err := h.repo.GetYearlyAmount(ctx, int(req.GetYear()), "success", 0)
	if err != nil {
		h.log.Error("GetYearlyAmountSuccess failed", zap.Error(err))
		return nil, err
	}
	return &pbtxstats.ApiResponseTransactionYearAmountSuccess{
		Status:  "success",
		Message: "Successfully fetched yearly amount success stats",
		Data:    mapYearlyAmountSuccess(data),
	}, nil
}

func (h *TransactionStatsHandler) GetMonthlyAmountFailed(ctx context.Context, req *pbtxstats.MonthAmountTransactionRequest) (*pbtxstats.ApiResponseTransactionMonthAmountFailed, error) {
	data, err := h.repo.GetMonthlyAmount(ctx, int(req.GetYear()), int(req.GetMonth()), "failed", 0)
	if err != nil {
		h.log.Error("GetMonthlyAmountFailed failed", zap.Error(err))
		return nil, err
	}
	return &pbtxstats.ApiResponseTransactionMonthAmountFailed{
		Status:  "success",
		Message: "Successfully fetched monthly amount failed stats",
		Data:    mapMonthlyAmountFailed(data),
	}, nil
}

func (h *TransactionStatsHandler) GetYearlyAmountFailed(ctx context.Context, req *pbtxstats.YearAmountTransactionRequest) (*pbtxstats.ApiResponseTransactionYearAmountFailed, error) {
	data, err := h.repo.GetYearlyAmount(ctx, int(req.GetYear()), "failed", 0)
	if err != nil {
		h.log.Error("GetYearlyAmountFailed failed", zap.Error(err))
		return nil, err
	}
	return &pbtxstats.ApiResponseTransactionYearAmountFailed{
		Status:  "success",
		Message: "Successfully fetched yearly amount failed stats",
		Data:    mapYearlyAmountFailed(data),
	}, nil
}

func (h *TransactionStatsHandler) GetMonthlyTransactionMethodSuccess(ctx context.Context, req *pbtxstats.MonthMethodTransactionRequest) (*pbtxstats.ApiResponseTransactionMonthPaymentMethod, error) {
	data, err := h.repo.GetMonthlyMethod(ctx, int(req.GetYear()), int(req.GetMonth()), "success", 0)
	if err != nil {
		h.log.Error("GetMonthlyTransactionMethodSuccess failed", zap.Error(err))
		return nil, err
	}
	return &pbtxstats.ApiResponseTransactionMonthPaymentMethod{
		Status:  "success",
		Message: "Successfully fetched monthly transaction method success stats",
		Data:    mapMonthlyMethod(data),
	}, nil
}

func (h *TransactionStatsHandler) GetYearlyTransactionMethodSuccess(ctx context.Context, req *pbtxstats.YearMethodTransactionRequest) (*pbtxstats.ApiResponseTransactionYearPaymentmethod, error) {
	data, err := h.repo.GetYearlyMethod(ctx, int(req.GetYear()), "success", 0)
	if err != nil {
		h.log.Error("GetYearlyTransactionMethodSuccess failed", zap.Error(err))
		return nil, err
	}
	return &pbtxstats.ApiResponseTransactionYearPaymentmethod{
		Status:  "success",
		Message: "Successfully fetched yearly transaction method success stats",
		Data:    mapYearlyMethod(data),
	}, nil
}

func (h *TransactionStatsHandler) GetMonthlyTransactionMethodFailed(ctx context.Context, req *pbtxstats.MonthMethodTransactionRequest) (*pbtxstats.ApiResponseTransactionMonthPaymentMethod, error) {
	data, err := h.repo.GetMonthlyMethod(ctx, int(req.GetYear()), int(req.GetMonth()), "failed", 0)
	if err != nil {
		h.log.Error("GetMonthlyTransactionMethodFailed failed", zap.Error(err))
		return nil, err
	}
	return &pbtxstats.ApiResponseTransactionMonthPaymentMethod{
		Status:  "success",
		Message: "Successfully fetched monthly transaction method failed stats",
		Data:    mapMonthlyMethod(data),
	}, nil
}

func (h *TransactionStatsHandler) GetYearlyTransactionMethodFailed(ctx context.Context, req *pbtxstats.YearMethodTransactionRequest) (*pbtxstats.ApiResponseTransactionYearPaymentmethod, error) {
	data, err := h.repo.GetYearlyMethod(ctx, int(req.GetYear()), "failed", 0)
	if err != nil {
		h.log.Error("GetYearlyTransactionMethodFailed failed", zap.Error(err))
		return nil, err
	}
	return &pbtxstats.ApiResponseTransactionYearPaymentmethod{
		Status:  "success",
		Message: "Successfully fetched yearly transaction method failed stats",
		Data:    mapYearlyMethod(data),
	}, nil
}

// --- TransactionStatsByMerchantService ---

func (h *TransactionStatsHandler) GetMonthlyAmountSuccessByMerchant(ctx context.Context, req *pbtxstats.MonthAmountTransactionMerchantRequest) (*pbtxstats.ApiResponseTransactionMonthAmountSuccess, error) {
	data, err := h.repo.GetMonthlyAmount(ctx, int(req.GetYear()), int(req.GetMonth()), "success", req.GetMerchantId())
	if err != nil {
		h.log.Error("GetMonthlyAmountSuccessByMerchant failed", zap.Error(err))
		return nil, err
	}
	return &pbtxstats.ApiResponseTransactionMonthAmountSuccess{
		Status:  "success",
		Message: "Successfully fetched monthly amount success stats by merchant",
		Data:    mapMonthlyAmountSuccess(data),
	}, nil
}

func (h *TransactionStatsHandler) GetYearlyAmountSuccessByMerchant(ctx context.Context, req *pbtxstats.YearAmountTransactionMerchantRequest) (*pbtxstats.ApiResponseTransactionYearAmountSuccess, error) {
	data, err := h.repo.GetYearlyAmount(ctx, int(req.GetYear()), "success", req.GetMerchantId())
	if err != nil {
		h.log.Error("GetYearlyAmountSuccessByMerchant failed", zap.Error(err))
		return nil, err
	}
	return &pbtxstats.ApiResponseTransactionYearAmountSuccess{
		Status:  "success",
		Message: "Successfully fetched yearly amount success stats by merchant",
		Data:    mapYearlyAmountSuccess(data),
	}, nil
}

func (h *TransactionStatsHandler) GetMonthlyAmountFailedByMerchant(ctx context.Context, req *pbtxstats.MonthAmountTransactionMerchantRequest) (*pbtxstats.ApiResponseTransactionMonthAmountFailed, error) {
	data, err := h.repo.GetMonthlyAmount(ctx, int(req.GetYear()), int(req.GetMonth()), "failed", req.GetMerchantId())
	if err != nil {
		h.log.Error("GetMonthlyAmountFailedByMerchant failed", zap.Error(err))
		return nil, err
	}
	return &pbtxstats.ApiResponseTransactionMonthAmountFailed{
		Status:  "success",
		Message: "Successfully fetched monthly amount failed stats by merchant",
		Data:    mapMonthlyAmountFailed(data),
	}, nil
}

func (h *TransactionStatsHandler) GetYearlyAmountFailedByMerchant(ctx context.Context, req *pbtxstats.YearAmountTransactionMerchantRequest) (*pbtxstats.ApiResponseTransactionYearAmountFailed, error) {
	data, err := h.repo.GetYearlyAmount(ctx, int(req.GetYear()), "failed", req.GetMerchantId())
	if err != nil {
		h.log.Error("GetYearlyAmountFailedByMerchant failed", zap.Error(err))
		return nil, err
	}
	return &pbtxstats.ApiResponseTransactionYearAmountFailed{
		Status:  "success",
		Message: "Successfully fetched yearly amount failed stats by merchant",
		Data:    mapYearlyAmountFailed(data),
	}, nil
}

func (h *TransactionStatsHandler) GetMonthlyTransactionMethodByMerchantSuccess(ctx context.Context, req *pbtxstats.MonthMethodTransactionMerchantRequest) (*pbtxstats.ApiResponseTransactionMonthPaymentMethod, error) {
	data, err := h.repo.GetMonthlyMethod(ctx, int(req.GetYear()), int(req.GetMonth()), "success", req.GetMerchantId())
	if err != nil {
		h.log.Error("GetMonthlyTransactionMethodByMerchantSuccess failed", zap.Error(err))
		return nil, err
	}
	return &pbtxstats.ApiResponseTransactionMonthPaymentMethod{
		Status:  "success",
		Message: "Successfully fetched monthly transaction method success stats by merchant",
		Data:    mapMonthlyMethod(data),
	}, nil
}

func (h *TransactionStatsHandler) GetYearlyTransactionMethodByMerchantSuccess(ctx context.Context, req *pbtxstats.YearMethodTransactionMerchantRequest) (*pbtxstats.ApiResponseTransactionYearPaymentmethod, error) {
	data, err := h.repo.GetYearlyMethod(ctx, int(req.GetYear()), "success", req.GetMerchantId())
	if err != nil {
		h.log.Error("GetYearlyTransactionMethodByMerchantSuccess failed", zap.Error(err))
		return nil, err
	}
	return &pbtxstats.ApiResponseTransactionYearPaymentmethod{
		Status:  "success",
		Message: "Successfully fetched yearly transaction method success stats by merchant",
		Data:    mapYearlyMethod(data),
	}, nil
}

func (h *TransactionStatsHandler) GetMonthlyTransactionMethodByMerchantFailed(ctx context.Context, req *pbtxstats.MonthMethodTransactionMerchantRequest) (*pbtxstats.ApiResponseTransactionMonthPaymentMethod, error) {
	data, err := h.repo.GetMonthlyMethod(ctx, int(req.GetYear()), int(req.GetMonth()), "failed", req.GetMerchantId())
	if err != nil {
		h.log.Error("GetMonthlyTransactionMethodByMerchantFailed failed", zap.Error(err))
		return nil, err
	}
	return &pbtxstats.ApiResponseTransactionMonthPaymentMethod{
		Status:  "success",
		Message: "Successfully fetched monthly transaction method failed stats by merchant",
		Data:    mapMonthlyMethod(data),
	}, nil
}

func (h *TransactionStatsHandler) GetYearlyTransactionMethodByMerchantFailed(ctx context.Context, req *pbtxstats.YearMethodTransactionMerchantRequest) (*pbtxstats.ApiResponseTransactionYearPaymentmethod, error) {
	data, err := h.repo.GetYearlyMethod(ctx, int(req.GetYear()), "failed", req.GetMerchantId())
	if err != nil {
		h.log.Error("GetYearlyTransactionMethodByMerchantFailed failed", zap.Error(err))
		return nil, err
	}
	return &pbtxstats.ApiResponseTransactionYearPaymentmethod{
		Status:  "success",
		Message: "Successfully fetched yearly transaction method failed stats by merchant",
		Data:    mapYearlyMethod(data),
	}, nil
}

// --- Mappers ---

func mapMonthlyAmountSuccess(data []repository.MonthlyAmount) []*pbtxstats.TransactionMonthlyAmountSuccess {
	var out []*pbtxstats.TransactionMonthlyAmountSuccess
	for _, d := range data {
		out = append(out, &pbtxstats.TransactionMonthlyAmountSuccess{
			Year:         d.Year,
			Month:        d.Month,
			TotalSuccess: int32(d.TotalCount),
			TotalAmount:  int32(d.TotalAmount),
		})
	}
	return out
}

func mapYearlyAmountSuccess(data []repository.YearlyAmount) []*pbtxstats.TransactionYearlyAmountSuccess {
	var out []*pbtxstats.TransactionYearlyAmountSuccess
	for _, d := range data {
		out = append(out, &pbtxstats.TransactionYearlyAmountSuccess{
			Year:         d.Year,
			TotalSuccess: int32(d.TotalCount),
			TotalAmount:  int32(d.TotalAmount),
		})
	}
	return out
}

func mapMonthlyAmountFailed(data []repository.MonthlyAmount) []*pbtxstats.TransactionMonthlyAmountFailed {
	var out []*pbtxstats.TransactionMonthlyAmountFailed
	for _, d := range data {
		out = append(out, &pbtxstats.TransactionMonthlyAmountFailed{
			Year:        d.Year,
			Month:       d.Month,
			TotalFailed: int32(d.TotalCount),
			TotalAmount: int32(d.TotalAmount),
		})
	}
	return out
}

func mapYearlyAmountFailed(data []repository.YearlyAmount) []*pbtxstats.TransactionYearlyAmountFailed {
	var out []*pbtxstats.TransactionYearlyAmountFailed
	for _, d := range data {
		out = append(out, &pbtxstats.TransactionYearlyAmountFailed{
			Year:        d.Year,
			TotalFailed: int32(d.TotalCount),
			TotalAmount: int32(d.TotalAmount),
		})
	}
	return out
}

func mapMonthlyMethod(data []repository.MonthlyMethod) []*pbtxstats.TransactionMonthlyMethod {
	var out []*pbtxstats.TransactionMonthlyMethod
	for _, d := range data {
		out = append(out, &pbtxstats.TransactionMonthlyMethod{
			Month:             d.Month,
			PaymentMethod:     d.PaymentMethod,
			TotalTransactions: int32(d.TotalCount),
			TotalAmount:       int32(d.TotalAmount),
		})
	}
	return out
}

func mapYearlyMethod(data []repository.YearlyMethod) []*pbtxstats.TransactionYearlyMethod {
	var out []*pbtxstats.TransactionYearlyMethod
	for _, d := range data {
		out = append(out, &pbtxstats.TransactionYearlyMethod{
			Year:              d.Year,
			PaymentMethod:     d.PaymentMethod,
			TotalTransactions: int32(d.TotalCount),
			TotalAmount:       int32(d.TotalAmount),
		})
	}
	return out
}
