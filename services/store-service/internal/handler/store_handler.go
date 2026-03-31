package handler

import (
	"github.com/labstack/echo/v4"
	"jualan-online/services/store-service/internal/service"
	"net/http"
)

type StoreHandler struct {
	svc service.StoreService
}

func NewStoreHandler(svc service.StoreService) *StoreHandler {
	return &StoreHandler{svc: svc}
}

func (h *StoreHandler) Create(c echo.Context) error {
	var req struct {
		Name         string `json:"name"`
		PlatformType string `json:"platform_type"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	store, err := h.svc.CreateStore(c.Request().Context(), req.Name, req.PlatformType)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, store)
}

func (h *StoreHandler) GetByID(c echo.Context) error {
	id := c.Param("id")
	store, err := h.svc.GetStoreByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "store not found"})
	}

	return c.JSON(http.StatusOK, store)
}

func (h *StoreHandler) GetAll(c echo.Context) error {
	stores, err := h.svc.GetAllStores(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, stores)
}

func (h *StoreHandler) Update(c echo.Context) error {
	id := c.Param("id")
	var req struct {
		Name         string `json:"name"`
		PlatformType string `json:"platform_type"`
		IsActive     bool   `json:"is_active"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	store, err := h.svc.UpdateStore(c.Request().Context(), id, req.Name, req.PlatformType, req.IsActive)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, store)
}

func (h *StoreHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if err := h.svc.DeleteStore(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *StoreHandler) Health(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "OK"})
}
