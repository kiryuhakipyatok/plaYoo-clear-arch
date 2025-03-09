package errors

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func NotFound(c *fiber.Ctx, log *logrus.Logger, what string,err error) error{
	c.Status(fiber.StatusNotFound)
	log.WithError(err).Error(fmt.Sprintf("%s not found",what))
	return c.JSON(fiber.Map{
		"error": fmt.Sprintf("%s not found: %s",what,err),
	})
}

func ErrorParse(c *fiber.Ctx, log *logrus.Logger, what string, err error) error{
	c.Status(fiber.StatusInternalServerError)
	log.WithError(err).Error(fmt.Sprintf("error parse %s",what))
	return c.JSON(fiber.Map{
		"error": fmt.Sprintf("error parse %s: %s",what,err),
	})
}

func ErrorFetching(c *fiber.Ctx, log *logrus.Logger,what string, err error) error{
	c.Status(fiber.StatusInternalServerError)
	log.WithError(err).Error(fmt.Sprintf("error fetching %s",what))
	return c.JSON(fiber.Map{
		"error": fmt.Sprintf("error fetching %s: %s",what,err),
	})
}

func FailedToValidate(c *fiber.Ctx, log *logrus.Logger,err error) error{
	c.Status(fiber.StatusInternalServerError)
	log.WithError(err).Error("failed to validate request")
	return c.JSON(fiber.Map{
		"error": fmt.Sprintf("failed to validate request: %s",err),
	})
}

func FailedToCreate(c *fiber.Ctx, log *logrus.Logger,what string,err error) error{
	c.Status(fiber.StatusInternalServerError)
	log.WithError(err).Error((fmt.Sprintf("failed to create %s",what)))
	return c.JSON(fiber.Map{
		"error": fmt.Sprintf("failed to create %s",what),
	})
}

func FailedToDelete(c *fiber.Ctx, log *logrus.Logger,what string,err error) error{
	c.Status(fiber.StatusInternalServerError)
	log.WithError(err).Error((fmt.Sprintf("failed to delete %s",what)))
	return c.JSON(fiber.Map{
		"error": fmt.Sprintf("failed to delte %s",what),
	})
}