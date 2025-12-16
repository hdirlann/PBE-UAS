package utils

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
	"time"
)

type ErrorResponse struct {
	Code        int         `json:"code"`
	Message     string      `json:"message"`
	Description string      `json:"description"`
	Details     interface{} `json:"details,omitempty"`
	Timestamp   time.Time   `json:"timestamp"`
}

type SuccessResponse struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// =========================
// ERROR HELPERS (SRS STRICT)
// =========================

// 400
func BadRequest(c *fiber.Ctx, details interface{}) error {
	return sendError(c, http.StatusBadRequest, "Bad Request", "Invalid input data", details)
}

// 401
func Unauthorized(c *fiber.Ctx, details interface{}) error {
	return sendError(c, http.StatusUnauthorized, "Unauthorized", "Missing or invalid token", details)
}

// 403
func Forbidden(c *fiber.Ctx, details interface{}) error {
	return sendError(c, http.StatusForbidden, "Forbidden", "Insufficient permissions", details)
}

// 404
func NotFound(c *fiber.Ctx, details interface{}) error {
	return sendError(c, http.StatusNotFound, "Not Found", "Resource not found", details)
}

// 409
func Conflict(c *fiber.Ctx, details interface{}) error {
	return sendError(c, http.StatusConflict, "Conflict", "Duplicate entry", details)
}

// 422
func Unprocessable(c *fiber.Ctx, details interface{}) error {
	return sendError(c, 422, "Unprocessable Entity", "Validation error", details)
}

// 500
func InternalError(c *fiber.Ctx, details interface{}) error {
	return sendError(c, http.StatusInternalServerError, "Internal Server Error", "Server error", details)
}

func sendError(c *fiber.Ctx, code int, message, description string, details interface{}) error {
	res := ErrorResponse{
		Code:        code,
		Message:     message,
		Description: description,
		Details:     details,
		Timestamp:   time.Now(),
	}
	return c.Status(code).JSON(res)
}

// ================
// SUCCESS HELPERS
// ================
func OK(c *fiber.Ctx, data interface{}) error {
	return sendSuccess(c, http.StatusOK, "OK", data)
}

func Created(c *fiber.Ctx, data interface{}) error {
	return sendSuccess(c, http.StatusCreated, "Created", data)
}

func sendSuccess(c *fiber.Ctx, status int, message string, data interface{}) error {
	res := SuccessResponse{
		Code:      status,
		Message:   message,
		Data:      data,
		Timestamp: time.Now(),
	}
	return c.Status(status).JSON(res)
}
