package main

import (
	"log"
	"time"

	"github.com/billing-engine/boostrap"
	"github.com/billing-engine/internal/commons"
	"github.com/billing-engine/internal/controller"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New(fiber.Config{})
	app.Use(cors.New())

	appConfig := boostrap.Boostrap()

	api := app.Group("/api")
	v1 := api.Group("/v1")

	controller := controller.NewController(appConfig)
	v1.Get("/get-outstanding", controller.GetOutstanding) // ✅
	v1.Get("/is-delinquent", controller.IsDelinquent)     // ✅
	v1.Post("/make-payment", controller.MakePayment)      // ✅
	v1.Post("/create-loan", controller.CreateLoan)        // ✅

	// schedule apps for checking loan from borrower
	go func() {
		ticker := time.NewTicker(time.Duration(commons.DifferentTime))
		defer ticker.Stop()

		for range ticker.C {
			err := controller.ProcessScheduleTask()
			if err != nil {
				log.Println("some error from scheduler", err)
			}
		}
	}()

	log.Fatal(app.Listen(":9005"))
}
