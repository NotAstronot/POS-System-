package rest

import (
	"net/http"
	"pos-system/internal/repository/postgres"
	"pos-system/internal/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

type EmployeeHandler struct {
	usecase *usecase.EmployeeUsecase
}

func NewEmployeeHandler(u *usecase.EmployeeUsecase) *EmployeeHandler {
	return &EmployeeHandler{usecase: u}
}

func (h *EmployeeHandler) List(c *gin.Context) {
	data, err := h.usecase.ListAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *EmployeeHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	item, err := h.usecase.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "karyawan tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}

func (h *EmployeeHandler) Create(c *gin.Context) {
	var req struct {
		NIK          string  `json:"nik" binding:"required"`
		Name         string  `json:"name" binding:"required"`
		Email        string  `json:"email"`
		Phone        string  `json:"phone"`
		DepartmentID *int64  `json:"department_id"`
		Position     string  `json:"position"`
		HireDate     string  `json:"hire_date"`
		SalaryBase   float64 `json:"salary_base"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hireDate := req.HireDate
	if hireDate == "" {
		hireDate = "CURRENT_DATE"
	}
	item := &postgres.Employee{
		NIK: req.NIK, Name: req.Name, Email: req.Email, Phone: req.Phone,
		DepartmentID: req.DepartmentID, Position: req.Position,
		HireDate: hireDate, SalaryBase: req.SalaryBase,
	}
	if err := h.usecase.Create(c.Request.Context(), item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}

func (h *EmployeeHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req struct {
		Name         string  `json:"name"`
		Email        string  `json:"email"`
		Phone        string  `json:"phone"`
		DepartmentID *int64  `json:"department_id"`
		Position     string  `json:"position"`
		HireDate     string  `json:"hire_date"`
		SalaryBase   float64 `json:"salary_base"`
		IsActive     *bool   `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item, err := h.usecase.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "karyawan tidak ditemukan"})
		return
	}
	if req.Name != "" {
		item.Name = req.Name
	}
	if req.Email != "" {
		item.Email = req.Email
	}
	if req.Phone != "" {
		item.Phone = req.Phone
	}
	if req.DepartmentID != nil {
		item.DepartmentID = req.DepartmentID
	}
	if req.Position != "" {
		item.Position = req.Position
	}
	if req.HireDate != "" {
		item.HireDate = req.HireDate
	}
	if req.SalaryBase > 0 {
		item.SalaryBase = req.SalaryBase
	}
	if req.IsActive != nil {
		item.IsActive = *req.IsActive
	}
	if err := h.usecase.Update(c.Request.Context(), item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}

func (h *EmployeeHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.usecase.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "karyawan dihapus"})
}
