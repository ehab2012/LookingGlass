package main

import (
    "log"
    "fmt"
    "errors"
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/websocket/v2"
    "github.com/gofiber/fiber/v2/middleware/favicon"
)

type Payload struct {
    Hostbox string `form:"hostbox"`
    Cmd string `form:"cmd"`
}

func main() {

    app := fiber.New(fiber.Config{
        ServerHeader:          "Looking glass",
        BodyLimit:             1 * 1024 * 1024 * 1024,
        Concurrency:           1 * 1024 * 1024 * 1024,
        DisableStartupMessage: false,
        ErrorHandler:          fiberErrorHandler,
    })

    app.Static("/", "./public")
    app.Use(favicon.New(favicon.Config{File: "./public/favicon.ico"}))

    app.Use("/ws", func(c *fiber.Ctx) error {
        // IsWebSocketUpgrade returns true if the client
        // requested upgrade to the WebSocket protocol.
        if websocket.IsWebSocketUpgrade(c) {
            c.Locals("allowed", true)
            return c.Next()
        }
        return fiber.ErrUpgradeRequired
    })

    app.Get("/ws", websocket.New(func(c *websocket.Conn) {
            var (
                mt  int
                msg []byte
                err error
            )
            for {
                if mt, msg, err = c.ReadMessage(); err != nil {
                    log.Println("read:", err)
                    break
                }
                log.Printf("recv: %s", msg)

                if err = c.WriteMessage(mt, msg); err != nil {
                    log.Println("write:", err)
                    break
                }
            }

        }))

/*
    app.Post("/", func(c *fiber.Ctx) error {
        payload := new(Payload)
        if err := c.BodyParser(payload); err != nil {
            return err
        }
        fmt.Println(payload.Hostbox)
        fmt.Println(payload.Cmd)
        return c.Send([]byte(payload.Hostbox)) // []byte("user=john")
    })    

    app.Get("/other", func(c *fiber.Ctx) error {
        return c.SendString("Hello, World!")
    })
*/
    log.Fatal(app.Listen(":3000"))
}

var fiberErrorHandler = func(c *fiber.Ctx, err error) error {
    fmt.Println(err, "http server handler error")
    c.JSON(map[string]string{"status": "fail", "result": "serverErr"})
    return errors.New("err")
}
