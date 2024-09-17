package routes

import (
	"pluto_remastered/controllers"
	mobile "pluto_remastered/controllers/mobile"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Landing Page!")
	})

	app.Post("login", controllers.Login)
	app.Post("sendOtp", controllers.SendOtp)

	// app.Get("/cronGenerateUserId", controllers.GenerateTransactionsUserId)
	// app.Get("/testCronGenerateUserId", controllers.TestGenerateUserId)
	// app.Get("/cronGenerateUserLog", controllers.GenerateUserLog)

	officeRoute := app.Group("office")
	// officeRoute.Use(AuthMiddleware)

	officeRoute.Get("/getProductTrends", controllers.GetProductTrends)
	officeRoute.Get("/TestQuery", controllers.TestQuery)
	officeRoute.Get("/getSalesmanDaily", controllers.GetSalesmanDailySales)

	mobileRoute := app.Group("pluto-mobile")
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
	mobileRoute.Get("getListOrderMD", mobile.GetListOrderMD)
	mobileRoute.Post("postOrderMD", mobile.PostOrderMD)

	mobileRoute.Use(AuthMiddleware)
	mobileRoute.Post("getRefreshUser", controllers.RefreshDataUser)
}
