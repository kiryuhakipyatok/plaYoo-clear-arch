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
		"error": fmt.Sprintf("%s not found: %s", what, err),
	})
}

func (eh *ErrorHandler) ErrorParse(c *fiber.Ctx, what string, err error) error {
	c.Status(fiber.StatusInternalServerError)
	eh.Logger.WithError(err).Error(fmt.Sprintf("error parse %s", what))
	return c.JSON(fiber.Map{
		"error": fmt.Sprintf("error parse %s: %s", what, err),
	})
}

func (eh *ErrorHandler) ErrorFetching(c *fiber.Ctx, what string, err error) error {
	c.Status(fiber.StatusInternalServerError)
	eh.Logger.WithError(err).Error(fmt.Sprintf("error fetching %s", what))
	return c.JSON(fiber.Map{
		"error": fmt.Sprintf("error fetching %s: %s", what, err),
	})
}

func (eh *ErrorHandler) FailedToValidate(c *fiber.Ctx, err error) error {
	c.Status(fiber.StatusInternalServerError)
	eh.Logger.WithError(err).Error("failed to validate request")
	return c.JSON(fiber.Map{
		"error": fmt.Sprintf("failed to validate request: %s", err),
	})
}

func (eh *ErrorHandler) FailedToCreate(c *fiber.Ctx, what string, err error) error {
	c.Status(fiber.StatusInternalServerError)
	eh.Logger.WithError(err).Error((fmt.Sprintf("failed to create %s", what)))
	return c.JSON(fiber.Map{
		"error": fmt.Sprintf("failed to create %s", what),
	})
}

func (eh *ErrorHandler) FailedToDelete(c *fiber.Ctx, what string, err error) error {
	c.Status(fiber.StatusInternalServerError)
	eh.Logger.WithError(err).Error((fmt.Sprintf("failed to delete %s", what)))
	return c.JSON(fiber.Map{
		"error": fmt.Sprintf("failed to delte %s", what),
	})
}
