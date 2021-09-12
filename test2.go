// https://git.sr.ht/~kiwec/reversi.network/tree/25e4bcca7182790e3df525fbb198a3e6893436a0/item/cmd/reversinet/main.go
// https://blog.csdn.net/qq_40530622/article/details/109212295
// https://github.com/gofiber/websocket
// https://github.com/spxvszero/go_shell_socket/blob/main/webpage.go
// https://github.com/ramnode/LookingGlass
// https://gist.github.com/tmichel/7390690
// websocat ws://127.0.0.1:4000/ws/ping4?host=1.2.3.4
package main

import (
    "log"
    "fmt"
    "errors"
    "encoding/json"
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/websocket/v2"
    "github.com/gofiber/fiber/v2/middleware/favicon"
    "bufio"
    "io"
    "os/exec"
)

type Payload struct {
    Cmd string `json:"cmd"`
    Host  string `json:"host"`
}

func main() {

    app := fiber.New(fiber.Config{
        ServerHeader:           "fiber",
        AppName:                "mrvm lookingglass",
        DisableStartupMessage:  true,
        ErrorHandler:           fiberErrorHandler,
    })
    app.Use(favicon.New(favicon.Config{File: "./public/favicon.ico"}))

    app.Static("/", "./public")

    app.Get("/remoteIp", func(c *fiber.Ctx) error {
        //queryValue := c.Query("callback")
        //return c.JSON(fiber.Map{"code": 200, "ip": c.IP()})
        msg := "remoteIp='" + c.IP() + "'";
        //c.Type("text/javascript") 
        return c.SendString(msg)
    })
    app.Use("/ws", func(c *fiber.Ctx) error {
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

                bytes := []byte(msg)

                // sanatize the bytes
                var p Payload
                err = json.Unmarshal(bytes, &p)
                if err != nil {
                    log.Println("read:", err)
                    break
                }

                // execute and get a pipe
                app := "echo"
                args := []string{"'unknown cmd'"}

                switch p.Cmd {
                    case "host":
                        app = "host"
                        args = []string{p.Host}
                    case "mtr4":
                        app = "mtr"
                        args = []string{"-4", "-r", "-w", p.Host}
                    case "mtr6":
                        // 2a00:c70:1:213:246:58:d8c:1
                        app = "mtr"
                        args = []string{"-6", "-r", "-w", p.Host}
                    case "ping4":
                        app = "ping"
                        args = []string{"-4", "-c", "5", p.Host}
                    case "ping6":
                        app = "ping"
                        args = []string{"-6" , "-c", "5", p.Host}
                    case "traceroute4":
                        app = "traceroute"
                        args = []string{"-4", "-w2", p.Host}
                    case "traceroute6":
                        app = "traceroute"
                        args = []string{"-6", "-w2", p.Host}
                }

                cmd := exec.Command(app, args...)

                stdout, err := cmd.StdoutPipe()
                if err != nil {
                    log.Println(err)
                    return
                }

                stderr, err := cmd.StderrPipe()
                if err != nil {
                    log.Println(err)
                    return
                }

                if err := cmd.Start(); err != nil {
                    log.Println(err)
                    return
                }

                s := bufio.NewScanner(io.MultiReader(stdout, stderr))
                for s.Scan() {
                    if err = c.WriteMessage(mt, s.Bytes()); err != nil {
                        log.Println("write:", err)
                        break
                    }
                }

                if err = c.WriteMessage(mt, []byte("done")); err != nil {
                    log.Println("write:", err)
                    break
                }

                if err := cmd.Wait(); err != nil {
                    log.Println(err)
                    return
                }
            }
        }))

    log.Fatal(app.Listen(":3000"))
}

var fiberErrorHandler = func(c *fiber.Ctx, err error) error {
    fmt.Println(err, "http server handler error")
    c.JSON(map[string]string{"status": "fail", "result": "serverErr"})
    return errors.New("err")
}