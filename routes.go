package main

import (
	"fmt"
	"reflect"

	"github.com/gofiber/fiber/v2"
)

func StartHttp() {
	InitWriteAheadLog(UpdateState)

	app := fiber.New()
	app.Get("/", StateHandler)

	app.Post("/move", HttpCommandHandler[MoveCommand]())
	app.Post("/move-all", HttpCommandHandler[MoveAllCommand]())
	app.Post("/increment", HttpCommandHandler[IncrementCommand]())
	app.Post("/decrement", HttpCommandHandler[DecrementCommand]())

	fmt.Println("Starting server on :8080")
	app.Listen(":8080")
}

func StateHandler(c *fiber.Ctx) error {
	return c.SendString("Count: " + fmt.Sprint(countCommands) + "\n" + "State: " + fmt.Sprint(counters) + "\n")
}

func HttpCommandHandler[K MoveCommand | MoveAllCommand | IncrementCommand | DecrementCommand]() fiber.Handler {
	var commandType string = reflect.TypeOf(*new(K)).Name()
	return func(c *fiber.Ctx) error {
		var payload K

		if err := c.BodyParser(&payload); err != nil {
			return err
		}

		responseChannel := ProposeCommandToWAL(commandType, payload)
		result := <-responseChannel
		return c.SendString("Result: " + fmt.Sprint(result))
	}
}
