package rest

import (
	"net/http"
	"pos-system/internal/repository/postgres"
	"pos-system/internal/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SalaryHandler struct {
	usecase *usecase.SalaryUsecase
}

func NewSalaryHandler(u *usecase.SalaryUsecase) *SalaryHandler {
	return &SalaryHandler{usecase: u}
}

func (h *SalaryHandler) List(c *gin.Context) {
	data, err := h.usecase.ListAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *SalaryHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	item, err := h.usecase.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "gaji tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}

func (h *SalaryHandler) Create(c *gin.Context) {
	var req struct {
		EmployeeID      int64   `json:"employee_id" binding:"required"`
		PayPeriod       string  `json:"pay_period" binding:"required"`
		BaseSalary      float64 `json:"base_salary" binding:"required"`
		CommissionTotal float64 `json:"commission_total"`
		Deduction       float64 `json:"deduction"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item := &postgres.Salary{
		EmployeeID: req.EmployeeID, PayPeriod: req.PayPeriod,
		BaseSalary: req.BaseSalary, CommissionTotal: req.CommissionTotal, Deduction: req.Deduction,
	}
	if err := h.usecase.Create(c.Request.Context(), item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}

func (h *SalaryHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req struct {
		EmployeeID      int64   `json:"employee_id"`
		PayPeriod       string  `json:"pay_period"`
		BaseSalary      float64 `json:"base_salary"`
		CommissionTotal float64 `json:"commission_total"`
		Deduction       float64 `json:"deduction"`
		Status          string  `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item, err := h.usecase.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "gaji tidak ditemukan"})
		return
	}
	if req.EmployeeID != 0 {
		item.EmployeeID = req.EmployeeID
	}
	if req.PayPeriod != "" {
		item.PayPeriod = req.PayPeriod
	}
	if req.BaseSalary > 0 {
		item.BaseSalary = req.BaseSalary
	}
	item.CommissionTotal = req.CommissionTotal
	item.Deduction = req.Deduction
	if req.Status != "" {
		item.Status = req.Status
	}
	if err := h.usecase.Update(c.Request.Context(), item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}

func (h *SalaryHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.usecase.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "gaji dihapus"})
}

func (h *SalaryHandler) MarkPaid(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.usecase.MarkPaid(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "gaji ditandai sudah dibayar"})
}
