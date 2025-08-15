package handler

import (
	"io"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"market/internal/middleware"
	"market/internal/service"
)

type PictureHandler struct {
	svc *service.PictureService
}

func NewPictureHandler(svc *service.PictureService) *PictureHandler {
	return &PictureHandler{svc: svc}
}

// POST /api/v1/products/:id/pictures (multipart form-data: file)
func (h *PictureHandler) Upload(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid product id")
	}
	file, err := c.FormFile("file")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "missing file")
	}
	f, err := file.Open()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "cannot open file")
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "cannot read file")
	}
	mime := file.Header.Get("Content-Type")
	sellerID := c.Locals(middleware.CtxUserID).(int64)
	pic, err := h.svc.UploadAndAttach(c.Context(), sellerID, id, data, mime)
	if err != nil {
		if err.Error() == "forbidden: not owner" {
			return fiber.NewError(fiber.StatusForbidden, err.Error())
		}
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.Status(fiber.StatusCreated).JSON(pic)
}

// GET /api/v1/products/:id/pictures (public)
func (h *PictureHandler) List(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid product id")
	}
	items, err := h.svc.List(c.Context(), id)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.JSON(items)
}

// GET /api/v1/pictures/:id (public) - отдает бинарник
func (h *PictureHandler) Download(c *fiber.Ctx) error {
	pid, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid picture id")
	}
	data, mime, err := h.svc.Download(c.Context(), pid)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}
	c.Set("Content-Type", mime)
	return c.Send(data)
}

// DELETE /api/v1/products/:id/pictures/:pid?hard=1
func (h *PictureHandler) Delete(c *fiber.Ctx) error {
	productID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid product id")
	}
	pictureID, err := strconv.ParseInt(c.Params("pid"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid picture id")
	}
	hard := c.Query("hard") == "1"
	sellerID := c.Locals(middleware.CtxUserID).(int64)
	if err := h.svc.Detach(c.Context(), sellerID, productID, pictureID, hard); err != nil {
		if err.Error() == "forbidden: not owner" {
			return fiber.NewError(fiber.StatusForbidden, err.Error())
		}
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// PUT /api/v1/products/:id/cover/:pid
func (h *PictureHandler) SetCover(c *fiber.Ctx) error {
	productID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid product id")
	}
	pictureID, err := strconv.ParseInt(c.Params("pid"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid picture id")
	}
	sellerID := c.Locals(middleware.CtxUserID).(int64)
	if err := h.svc.SetCover(c.Context(), sellerID, productID, pictureID); err != nil {
		if err.Error() == "forbidden: not owner" {
			return fiber.NewError(fiber.StatusForbidden, err.Error())
		}
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.SendStatus(fiber.StatusNoContent)
}
