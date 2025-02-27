package cli

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/oleoneto/mock-http/pkg/schema"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var MockServerCmd = &cobra.Command{
	Use:    "server",
	Short:  "Run a mock HTTP server",
	PreRun: state.Flags.File.StdinHook("file"),
	Run: func(cmd *cobra.Command, args []string) {

		err := yaml.Unmarshal(state.Flags.File.Data, &config)
		if err != nil {
			log.Fatalln(err)
		}

		app := fiber.New(fiber.Config{
			DisableStartupMessage: false,
			EnablePrintRoutes:     showRoutes,
		})

		app.Use(recover.New(recover.Config{EnableStackTrace: false}))
		app.Use(func(c *fiber.Ctx) error {
			t := time.Now() // request start time

			err := c.Next()

			logrus.WithFields(logrus.Fields{
				"duration": time.Since(t),
				"protocol": c.Protocol(),
				"method":   c.Method(),
				"ip":       c.IP(),
				"path":     c.Path(),
				"status":   c.Response().StatusCode(),
			}).Infoln(c.Path())

			return err
		})

		for _, r := range config.Routes {
			app.Add(r.Method, r.Path, func(c *fiber.Ctx) error {
				b, err := json.Marshal(r.Body)
				if err != nil {
					return err
				}

				c.Status(r.StatusCode)
				return c.Send(b)
			})
		}

		// app.All("/*", func(c *fiber.Ctx) error { return c.SendStatus(http.StatusNotImplemented) })

		app.Listen(fmt.Sprintf("127.0.0.1:%d", serverPort))
	},
}

func init() {
	MockServerCmd.Flags().VarP(&state.Flags.File, "file", "f", "")
	MockServerCmd.Flags().IntVarP(&serverPort, "port", "p", serverPort, "")
	MockServerCmd.Flags().BoolVarP(&showRoutes, "show-routes", "r", showRoutes, "")
}

var (
	serverPort = 3333
	showRoutes bool
	config     struct {
		Routes []schema.Route `yaml:"routes" json:"routes"`
	}
)
