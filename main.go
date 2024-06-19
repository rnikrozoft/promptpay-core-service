package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"github.com/golang-jwt/jwt/v5"
	"github.com/samber/lo"
	"github.com/spf13/viper"

	_ "github.com/rnikrozoft/promptpay-core-service/docs"
	"github.com/rnikrozoft/promptpay-core-service/handler"
	"github.com/rnikrozoft/promptpay-core-service/pkg/errs"
	"github.com/rnikrozoft/promptpay-core-service/pkg/logs"
	"github.com/rnikrozoft/promptpay-core-service/repository"
	"github.com/rnikrozoft/promptpay-core-service/service/user"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"

	jwtware "github.com/gofiber/contrib/jwt"
)

var (
	db   *gorm.DB
	conf *Configs
	ctx  = context.Background()
)

func init() {
	//read config from config.yml
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	lo.Must0(viper.ReadInConfig())
	lo.Must0(viper.Unmarshal(&conf))

	//connect to database
	db = lo.Must1(gorm.Open(sqlserver.Open(conf.Database.Address), &gorm.Config{}))
	getSecret(conf.App.Azure)
}

// get secret from azure key vault
func getSecret(az Azure) {
	cred := lo.Must1(azidentity.NewDefaultAzureCredential(nil))
	client := lo.Must1(azsecrets.NewClient(az.KeyVaultURL, cred, nil))
	resp := lo.Must1(client.GetSecret(ctx, az.SecretName, "", nil))
	conf.App.Azure.SecretVault = *resp.Value
}

func main() {
	app := fiber.New()
	app.Use(
		recover.New(),
		cors.New(),
	)

	app.Get("/", monitor.New())
	app.Get("/swagger/*", swagger.HandlerDefault)

	userRepository := repository.NewUser(db)
	userService := user.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userService)
	app.Get("/users", userHandler.GetUser)
	app.Post("/login", login)

	//middleware
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(conf.App.Azure.SecretVault)},

		ErrorHandler: func(_ *fiber.Ctx, _ error) error {
			return errs.NewUnaurhenticatedError()
		},
	}))

	//sadas
	app.Get("/restricted", restricted)

	go func() {
		if err := app.Listen(fmt.Sprintf(":%d", conf.App.Port)); err != nil {
			logs.Error(fmt.Sprintf("Error starting server: %v", err))
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	_ = <-c // This blocks the main thread until an interrupt is received
	fmt.Println("Gracefully shutting down...")
	_ = app.Shutdown()

	fmt.Println("Running cleanup tasks...")

	// Your cleanup tasks go here
	// db.Close()
	// redisConn.Close()
	fmt.Println("Fiber was successful shutdown.")
}

func login(c *fiber.Ctx) error {
	user := c.FormValue("user")
	pass := c.FormValue("pass")

	// Throws Unauthorized error
	if user != "john" || pass != "doe" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	// Create the Claims
	claims := jwt.MapClaims{
		"name":  "John Doe",
		"admin": true,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(conf.App.Azure.SecretVault))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"token": t})
}

func restricted(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	return c.SendString("Welcome " + name)
}
