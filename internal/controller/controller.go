package controller

import (
	"context"

	"github.com/billing-engine/config"
	"github.com/billing-engine/internal/service"
	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	AppConfig *config.AppConfig
}

func NewController(appConfig *config.AppConfig) *Controller {
	return &Controller{
		AppConfig: appConfig,
	}
}

type GetOunstandingRequest struct {
	Username string `json:"username"`
}

type IsDelinquentRequest struct {
	Usernanme string `json:"username"`
}

type MakePaymentRequest struct {
	Username string  `json:"username"`
	Amount   float64 `json:"amount"`
}

type CreateLoanRequest struct {
	Username string  `json:"username"`
	Amount   float64 `json:"amount"`
}

type GetOunstandingResponse struct {
	Amount float64 `json:"amount"`
	Status string  `json:"status"`
}

func (ctrl *Controller) GetOutstanding(c *fiber.Ctx) error {
	input := new(GetOunstandingRequest)

	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"is_error": true,
			"message":  "failed to parsing data",
		})
	}

	amount, err := ctrl.AppConfig.Service.GetOutStanding(context.Background(), input.Username)
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"is_error": true,
			"message":  "failed create loan",
			"error":    err.Error(),
		})
	}

	response := GetOunstandingResponse{
		Amount: amount,
		Status: "still exist",
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"is_error": false,
		"success":  "success",
		"data":     response,
		"message":  "successfully created",
	})
}

func (ctrl *Controller) IsDelinquent(c *fiber.Ctx) error {
	input := new(IsDelinquentRequest)

	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"is_error": true,
			"message":  "failed to parsing data",
		})
	}

	IsDelinquent, err := ctrl.AppConfig.Service.IsDelinquent(context.Background(), input.Usernanme)
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"is_error": true,
			"message":  "failed create loan",
			"error":    err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"is_error":      false,
		"is_delinquent": IsDelinquent,
		"message":       "Status user delinquent",
	})
}

func (ctrl *Controller) MakePayment(c *fiber.Ctx) error {
	input := new(MakePaymentRequest)

	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"is_error": true,
			"message":  "failed to parsing data",
		})
	}

	message, err := ctrl.AppConfig.Service.MakePayment(context.Background(), service.MakePaymentEntity{
		Username: input.Username,
		Amount:   input.Amount,
	})
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"is_error": true,
			"message":  "failed make payment",
			"error":    err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"is_error": false,
		"success":  "success",
		"message":  message,
	})
}

func (ctrl *Controller) CreateLoan(c *fiber.Ctx) error {
	input := new(CreateLoanRequest)

	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"is_error": true,
			"message":  "failed to parsing data",
		})
	}

	err := ctrl.AppConfig.Service.CreateLoan(context.Background(), service.CreateLoanEntity{
		Username: input.Username,
		Amount:   input.Amount,
	})
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"is_error": true,
			"message":  "failed create loan",
			"error":    err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"is_error": false,
		"success":  "success",
		"message":  "successfully created",
	})
}

func (ctrl *Controller) ProcessScheduleTask() error {
	err := ctrl.AppConfig.Service.ScheduleTask(context.Background())
	if err != nil {
		return err
	}

	return nil
}
