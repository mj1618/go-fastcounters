package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func main() {
	InitWriteAheadLog(UpdateState)

	app := fiber.New()
	app.Get("/", RootHandler)
	app.Get("/state", StateHandler)
	fmt.Println("Starting server on :8080")
	app.Listen(":8080")
}

func RootHandler(c *fiber.Ctx) error {
	responseChannel := ProposeCommandToWAL("MoveCommand", MoveCommand{FromAddress: 1, ToAddress: 2, Amount: 10})
	result := <-responseChannel
	return c.SendString("Result: " + fmt.Sprint(result))
}

func StateHandler(c *fiber.Ctx) error {
	return c.SendString("State: " + fmt.Sprint(GetCommandCounts()))
}
