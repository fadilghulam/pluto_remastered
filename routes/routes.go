package routes

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"pluto_remastered/controllers"
	mobile "pluto_remastered/controllers/mobile"
)

func Setup(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Landing Page2!")
	})

	app.Post("login", controllers.Login)
	app.Post("sendOtp", controllers.SendOtp)

	app.Get("/getUtcTime", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "success",
			"data":    time.Now().UTC(),
		})
	})

	app.Get("/cronGenerateUserId", controllers.GenerateTransactionsUserId)
	app.Get("/testCronGenerateUserId", controllers.TestGenerateUserId)
	app.Get("/cronGenerateUserLog", controllers.GenerateUserLog)

	app.Get("/generateFlag", controllers.GenerateFlag)
	app.Get("/getData", controllers.GetData)
	app.Get("/getDataToday", controllers.GetDataToday)

	app.Post("/insertTransactions", mobile.InsertTransactions)

	serviceRoute := app.Group("service")
	serviceRoute.Post("doUpload", controllers.DoUpload)

	officeRoute := app.Group("office")
	// officeRoute.Use(AuthMiddleware)

	officeRoute.Get("/getProductTrends", controllers.GetProductTrends)
	officeRoute.Get("/TestQuery", controllers.TestQuery)
	officeRoute.Get("/getSalesmanDaily", controllers.GetSalesmanDailySales)
	officeRoute.Get("/getUserBranch", controllers.GetUserBranch)

	mobileRoute := app.Group("pluto-mobile")
	mobileRoute.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Landing Page Pluto Mobile!")
	})

	mobileRoute.Get("getAppVersioning", mobile.GetAppVersioning)

	mobileRoute.Get("getGudang", mobile.GetGudang)
	mobileRoute.Get("getProdukGudang", mobile.GetProdukByGudang)
	mobileRoute.Get("getItemGudang", mobile.GetItemByGudang)
	mobileRoute.Post("confirmOrder", mobile.ConfirmOrder)

	mobileRoute.Get("getListPengajuan", mobile.GetDataRequests)

	//sales
	mobileRoute.Get("getStokProduk", mobile.GetStokProduk)
	mobileRoute.Get("getListOrder", mobile.GetListOrder)
	mobileRoute.Post("postOrder", mobile.PostOrder)

	//md
	mobileRoute.Get("getStokItem", mobile.GetStokItem)
	mobileRoute.Get("getListOrderItem", mobile.GetListOrderMD)
	mobileRoute.Post("postOrderItem", mobile.PostOrderMD)

	mobileRoute.Use(AuthMiddleware)
	mobileRoute.Post("getRefreshUser", controllers.RefreshDataUser)
}
