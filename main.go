package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func main() {
	InitWriteAheadLog(UpdateState)

	app := fiber.New()
	app.Get("/", StateHandler)

	app.Post("/move", MoveHandler)
	app.Post("/increment", IncrementHandler)
	app.Post("/decrement", DecrementHandler)

	fmt.Println("Starting server on :8080")
	app.Listen(":8080")
}

func StateHandler(c *fiber.Ctx) error {
	return c.SendString("Count: " + fmt.Sprint(countCommands) + "\n" + "State: " + fmt.Sprint(counters) + "\n")
}

func MoveHandler(c *fiber.Ctx) error {
	payload := MoveCommand{}

	if err := c.BodyParser(&payload); err != nil {
		return err
	}

	responseChannel := ProposeCommandToWAL("MoveCommand", payload)
	result := <-responseChannel
	return c.SendString("Result: " + fmt.Sprint(result))
}

func IncrementHandler(c *fiber.Ctx) error {
	payload := IncrementCommand{}

	if err := c.BodyParser(&payload); err != nil {
		return err
	}

	responseChannel := ProposeCommandToWAL("IncrementCommand", payload)
	result := <-responseChannel
	return c.SendString("Result: " + fmt.Sprint(result))
}

func DecrementHandler(c *fiber.Ctx) error {
	payload := DecrementCommand{}

	if err := c.BodyParser(&payload); err != nil {
		return err
	}

	responseChannel := ProposeCommandToWAL("DecrementCommand", payload)
	result := <-responseChannel
	return c.SendString("Result: " + fmt.Sprint(result))
}
