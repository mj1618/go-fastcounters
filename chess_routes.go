package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/mj1618/go-fastcounters/wal"
)

func ChessStartHttp() {
	app := fiber.New()
	app.Get("/", ChessStateHandler)

	app.Post("/move", ChessHttpCommandHandler())

	fmt.Println("Starting server on :8080")
	app.Listen(":8080")
}

func ChessStateHandler(c *fiber.Ctx) error {
	return c.SendString("State: " + fmt.Sprint(GetChessState()) + "\n")
}

func ChessHttpCommandHandler() fiber.Handler {

	return func(c *fiber.Ctx) error {
		payload := MoveCommand{}
		responseChannel := wal.ProposeCommandToWAL("ChessMoveCommand", payload)
		result := <-responseChannel
		return c.SendString("Result: " + fmt.Sprint(result) + "\n")
	}
}
