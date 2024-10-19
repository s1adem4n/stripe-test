package main

import (
	"fmt"
	"log"
	"net/http"
	"net/mail"
	"os"
	"regexp"
	"strings"
	"time"

	"math/rand"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	"github.com/pocketbase/pocketbase/tools/mailer"
	stripe "github.com/stripe/stripe-go/v80"
	stripeClient "github.com/stripe/stripe-go/v80/client"
	"googlemaps.github.io/maps"

	"striper/frontend"
	_ "striper/migrations"
)

func SPAMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		if err != nil {
			if he, ok := err.(*echo.HTTPError); ok {
				if he.Code == http.StatusNotFound {
					c.Request().URL.Path = "/"
					return next(c)
				}
			}
		}
		return err
	}
}

type CheckoutParams struct {
	VerificationCode string  `json:"verificationCode"`
	Name             string  `json:"name"`
	Email            string  `json:"email"`
	Phone            string  `json:"phone"`
	Address          Address `json:"address"`
}
type Address struct {
	Line1      string `json:"line1"`
	Line2      string `json:"line2"`
	City       string `json:"city"`
	PostalCode string `json:"postalCode"`
}

type Config struct {
	StripeSecretKey             string `envconfig:"STRIPE_SECRET_KEY"`
	GoogleMapsAPIKey            string `envconfig:"GOOGLE_MAPS_API_KEY"`
	AlwaysValidVerificationCode string `envconfig:"ALWAYS_VALID_VERIFICATION_CODE"`
}

func generateVerificationCode() string {
	// random number between 0 and 999999, padded with zeros
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func vaildateEmail(email string) bool {
	return emailRegex.MatchString(email)
}

const VerifyTemplate = `Gib folgenden Code ein, um deine Bestellung abzuschließen: %s
Er ist 5 Minuten lang gültig.

Dein %s Team
`

type AddressPrediction struct {
	PlaceID     string `json:"placeId"`
	Description string `json:"description"`
}

func main() {
	godotenv.Load()

	app := pocketbase.New()
	logger := app.Logger()

	var config Config
	if err := envconfig.Process("", &config); err != nil {
		logger.Error("failed to load the config", "error", err)
		return
	}
	fmt.Printf("Config: %+v\n", config)

	isGoRun := strings.HasPrefix(os.Args[0], os.TempDir())

	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		Automigrate: isGoRun,
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		sc := &stripeClient.API{}
		sc.Init(config.StripeSecretKey, nil)

		gc, err := maps.NewClient(maps.WithAPIKey(config.GoogleMapsAPIKey))
		if err != nil {
			logger.Error("failed to create google maps client", "error", err)
			return err
		}

		appURL := app.Settings().Meta.AppUrl

		mails := app.NewMailClient()
		dao := app.Dao()
		codes, err := dao.FindCollectionByNameOrId("codes")
		if err != nil {
			logger.Error("failed to find codes collection", "error", err)
			return err
		}

		e.Router.Use(SPAMiddleware)
		e.Router.Use(middleware.GzipWithConfig(middleware.GzipConfig{
			Skipper: func(c echo.Context) bool {
				return strings.HasPrefix(c.Request().URL.Path, "/_/")
			},
		}))

		subFS := echo.MustSubFS(frontend.Assets, "build")
		e.Router.StaticFS("/", subFS)

		rateLimited := e.Router.Group("")
		rateLimited.Use(middleware.RateLimiterWithConfig(middleware.RateLimiterConfig{
			Store: middleware.NewRateLimiterMemoryStore(1),
			IdentifierExtractor: func(c echo.Context) (string, error) {
				var ip string
				if c.RealIP() != "" {
					ip = c.RealIP()
				} else {
					remoteAddr := c.Request().RemoteAddr
					ip = remoteAddr[:strings.LastIndex(remoteAddr, ":")]
				}
				return ip, nil
			},
			Skipper: func(c echo.Context) bool {
				return c.Request().Method == http.MethodOptions
			},
		}))

		rateLimited.GET("/api/autocomplete-address", func(c echo.Context) error {
			ctx := c.Request().Context()
			res, err := gc.PlaceAutocomplete(ctx, &maps.PlaceAutocompleteRequest{
				Input:    c.QueryParam("input"),
				Language: "de",
				Components: map[maps.Component][]string{
					maps.ComponentCountry: {"DE"},
				},
				Types: maps.AutocompletePlaceTypeAddress,
			})
			if err != nil {
				logger.Error("failed to autocomplete address", "error", err)
				return c.String(http.StatusInternalServerError, "Failed to autocomplete address")
			}

			var predictions []AddressPrediction
			for _, prediction := range res.Predictions {
				predictions = append(predictions, AddressPrediction{
					PlaceID:     prediction.PlaceID,
					Description: prediction.Description,
				})
			}

			return c.JSON(http.StatusOK, predictions)
		})
		rateLimited.GET("/api/place-details", func(c echo.Context) error {
			ctx := c.Request().Context()
			res, err := gc.PlaceDetails(ctx, &maps.PlaceDetailsRequest{
				PlaceID:  c.QueryParam("placeId"),
				Language: "de",
				Fields: []maps.PlaceDetailsFieldMask{
					maps.PlaceDetailsFieldMaskAddressComponent,
				},
			})
			if err != nil {
				logger.Error("failed to get place details", "error", err)
				return c.String(http.StatusInternalServerError, "Failed to get place details")
			}

			var address Address
			var streetNumber string
			for _, component := range res.AddressComponents {
				switch component.Types[0] {
				case "route":
					address.Line1 = component.LongName
				case "street_number":
					streetNumber = component.LongName
				case "locality":
					address.City = component.LongName
				case "postal_code":
					address.PostalCode = component.LongName
				}
			}
			if streetNumber != "" {
				address.Line1 += " " + streetNumber
			}

			return c.JSON(http.StatusOK, address)
		})

		rateLimited.GET("/api/verify", func(c echo.Context) error {
			email := c.QueryParam("email")
			if !vaildateEmail(email) {
				return c.String(http.StatusBadRequest, "Invalid email")
			}

			// delete any existing records for this email
			records, err := dao.FindRecordsByFilter("codes", "email = {:email}", "", 100, 0, dbx.Params{
				"email": email,
			})
			if err != nil {
				return c.String(http.StatusInternalServerError, "Failed to find records")
			}
			if len(records) > 0 {
				for _, record := range records {
					if err := dao.DeleteRecord(record); err != nil {
						return c.String(http.StatusInternalServerError, "Failed to delete record")
					}
				}
			}

			code := generateVerificationCode()
			record := models.NewRecord(codes)
			record.Set("email", email)
			record.Set("code", code)
			expires := time.Now().Add(5 * time.Minute)
			record.Set("expires", expires)

			if err := dao.Save(record); err != nil {
				return c.String(http.StatusInternalServerError, "Failed to save record")
			}

			err = mails.Send(&mailer.Message{
				From: mail.Address{
					Address: app.Settings().Meta.SenderAddress,
					Name:    app.Settings().Meta.SenderName,
				},
				To: []mail.Address{
					{
						Address: email,
					},
				},
				Subject: fmt.Sprintf("Dein Verifizierungscode für %s", app.Settings().Meta.AppName),
				Text:    fmt.Sprintf(VerifyTemplate, code, app.Settings().Meta.AppName),
			})
			if err != nil {
				app.Logger().Error("failed to send email", "error", err, "email", email)
				app.Dao().DeleteRecord(record)

				return c.String(http.StatusInternalServerError, "Failed to send email")
			}

			return c.String(http.StatusOK, "Code sent")
		})

		e.Router.POST("/api/checkout", func(c echo.Context) error {
			var params CheckoutParams
			if err := c.Bind(&params); err != nil {
				return c.String(http.StatusBadRequest, err.Error())
			}

			if params.VerificationCode != config.AlwaysValidVerificationCode {
				record, err := dao.FindFirstRecordByFilter("codes", "email = {:email}", dbx.Params{
					"email": params.Email,
				})
				if err != nil {
					return c.String(http.StatusUnauthorized, "Verification failed")
				}
				if record.GetDateTime("expires").Time().Before(time.Now()) {
					return c.String(http.StatusUnauthorized, "Verification code expired")
				}
				code := record.GetString("code")
				if code != params.VerificationCode {
					return c.String(http.StatusUnauthorized, "Verification failed")
				}
				if err := dao.DeleteRecord(record); err != nil {
					return c.String(http.StatusInternalServerError, "Failed to delete record")
				}
			}

			var customerID string
			customers := sc.Customers.List(&stripe.CustomerListParams{
				Email: stripe.String(params.Email),
			})
			if customers.Err() != nil {
				app.Logger().Error("failed to list customers", "error", customers.Err(), "path", c.Path(), "email", params.Email)

				return c.String(http.StatusInternalServerError, "customer lookup failed")
			}
			if customers.Next() {
				customerID = customers.Customer().ID
			} else {
				customer, err := sc.Customers.New(&stripe.CustomerParams{
					Name:  stripe.String(params.Name),
					Email: stripe.String(params.Email),
					Phone: stripe.String(params.Phone),
					Address: &stripe.AddressParams{
						Line1:      stripe.String(params.Address.Line1),
						Line2:      stripe.String(params.Address.Line2),
						City:       stripe.String(params.Address.City),
						PostalCode: stripe.String(params.Address.PostalCode),
						Country:    stripe.String("DE"),
					},
				})
				if err != nil {
					app.Logger().Error("failed to create customer", "error", err, "path", c.Path(), "email", params.Email)

					return c.String(http.StatusInternalServerError, "customer creation failed")
				}
				customerID = customer.ID
			}

			sessionParams := stripe.CheckoutSessionParams{
				Customer: stripe.String(customerID),
				LineItems: []*stripe.CheckoutSessionLineItemParams{
					{
						PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
							Currency: stripe.String(string(stripe.CurrencyEUR)),
							ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
								Name:        stripe.String("Keks"),
								Description: stripe.String("lecker lecker"),
								TaxCode:     stripe.String("txcd_99999999"),
							},
							UnitAmount: stripe.Int64(1500),
						},
						Quantity: stripe.Int64(1),
					},
					{
						PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
							Currency: stripe.String(string(stripe.CurrencyEUR)),
							ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
								Name:        stripe.String("Kuchen"),
								Description: stripe.String("lecker lecker"),
								TaxCode:     stripe.String("txcd_99999999"),
							},
							UnitAmount: stripe.Int64(2500),
						},
						Quantity: stripe.Int64(1),
					},
				},
				Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
				SuccessURL: stripe.String(appURL + "/success"),
				CancelURL:  stripe.String(appURL + "/cancel"),
				AutomaticTax: &stripe.CheckoutSessionAutomaticTaxParams{
					Enabled: stripe.Bool(true),
				},
			}

			session, err := sc.CheckoutSessions.New(&sessionParams)
			if err != nil {
				return c.String(http.StatusInternalServerError, err.Error())
			}

			return c.JSON(http.StatusOK, echo.Map{
				"url": session.URL,
			})
		})

		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
