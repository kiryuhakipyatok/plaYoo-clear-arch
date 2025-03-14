package errors

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type ErrorHandler struct {
	Logger *logrus.Logger
}

func NewErrorHandler(log *logrus.Logger) *ErrorHandler {
	return &ErrorHandler{
		Logger: log,
	}
}

func (eh *ErrorHandler) NotFound(c *fiber.Ctx, what string, err error) error {
	c.Status(fiber.StatusNotFound)
	eh.Logger.WithError(err).Error(fmt.Sprintf("%s not found", what))
	return c.JSON(fiber.Map{
		"error": fmt.Errorf("%s not found: %w", what, err),
	})
}

func (eh *ErrorHandler) ErrorParse(c *fiber.Ctx, what string, err error) error {
	c.Status(fiber.StatusInternalServerError)
	eh.Logger.WithError(err).Error(fmt.Sprintf("error parse %s", what))
	return c.JSON(fiber.Map{
		"error": fmt.Errorf("error parse %s: %w", what, err),
	})
}

func (eh *ErrorHandler) ErrorFetching(c *fiber.Ctx, what string, err error) error {
	c.Status(fiber.StatusInternalServerError)
	eh.Logger.WithError(err).Error(fmt.Sprintf("error fetching %s", what))
	return c.JSON(fiber.Map{
		"error": fmt.Errorf("error fetching %s: %w", what, err),
	})
}

func (eh *ErrorHandler) FailedToValidate(c *fiber.Ctx, err error) error {
	c.Status(fiber.StatusInternalServerError)
	eh.Logger.WithError(err).Error("failed to validate request")
	return c.JSON(fiber.Map{
		"error": fmt.Errorf("failed to validate request: %w", err),
	})
}

func (eh *ErrorHandler) FailedToCreate(c *fiber.Ctx, what string, err error) error {
	c.Status(fiber.StatusInternalServerError)
	eh.Logger.WithError(err).Error((fmt.Sprintf("failed to create %s", what)))
	return c.JSON(fiber.Map{
		"error": fmt.Errorf("failed to create %s: %w", what,err),
	})
}

func (eh *ErrorHandler) FailedToDelete(c *fiber.Ctx, what string, err error) error {
	c.Status(fiber.StatusInternalServerError)
	eh.Logger.WithError(err).Error((fmt.Sprintf("failed to delete %s", what)))
	return c.JSON(fiber.Map{
		"error": fmt.Errorf("failed to delte %s: %w", what,err),
	})
}
