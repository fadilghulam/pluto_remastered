package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"

	db "pluto_remastered/config"
	"pluto_remastered/controllers"
	mobile "pluto_remastered/controllers/mobile"
	"pluto_remastered/helpers"

	"github.com/zishang520/socket.io/v2/socket"
)

func Setup(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Landing Page local pluto remastered!")
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
	app.Get("/completeSales", mobile.CompleteSales)

	app.Get("/getPermission", mobile.GetPermission)

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

	mobileRoute.Get("getCheckSO", mobile.GetCheckSO)
	mobileRoute.Post("setStock", mobile.SetStock)

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

type NotificationPayload struct {
	UserId        *string `json:"user_id"`
	ReferenceName *string `json:"reference_name"`
	ReferenceId   *string `json:"reference_id"`
}

func SocketIoSetup(app *fiber.App) {
	// Socket.IO setup
	fmt.Println("socket io startup")

	socketio := socket.NewServer(nil, nil)

	// socketio.On("connection", func(clients ...interface{}) {

	// socketio.Of("/", nil).On("connection", func(clients ...interface{}) {

	// })

	//send to specific user
	app.Get("/send/:id", func(c *fiber.Ctx) error {
		clientID := c.Params("id")

		socketio.To(socket.Room(clientID)).Emit("returnData", "this is return data to specific client")

		return c.SendString("Message sent to client: " + clientID)
	})

	app.Get("/checkRoom/:room", func(c *fiber.Ctx) error {
		room := c.Params("room")

		temp := 0
		socketio.In(socket.Room(room)).FetchSockets()(func(sockets []*socket.RemoteSocket, err error) {
			if err != nil {
				// handle error
				fmt.Println("Error fetching sockets:", err)
			}

			// for _, _ := range sockets {
			// 	temp++
			// }
			temp = len(sockets)
		})

		socketio.To(socket.Room(room)).Emit("returnData", "this is return data for room")

		return c.SendString("Room checked for: " + room + " sum sockets : " + strconv.Itoa(temp))
	})

	// channel := "testnotify"
	_, err := db.DBPGX.Exec(context.Background(), "LISTEN testnotify")
	if err != nil {
		fmt.Println("Error setting up listener: %v\n", err.Error())
	}
	_, err = db.DBPGX.Exec(context.Background(), "LISTEN testnotify2")
	if err != nil {
		fmt.Println("Error setting up listener: %v\n", err.Error())
	}
	_, err = db.DBPGX.Exec(context.Background(), "LISTEN testnotify3")
	if err != nil {
		fmt.Println("Error setting up listener: %v\n", err.Error())
	}

	go func() {
		for {
			notification, err := db.DBPGX.WaitForNotification(context.Background())
			if err != nil {
				log.Fatalf("Error waiting for notification: %v\n", err)
			}
			fmt.Printf("Received notification on channel %s: %s\n", notification.Channel, notification.Payload)

			var payload NotificationPayload
			err = json.Unmarshal([]byte(notification.Payload), &payload)
			if err != nil {
				log.Fatalf("Error unmarshalling JSON: %v\n", err)
			}

			fmt.Println(&payload)

			// Emit notification to all connected Socket.IO clients
			// socket.BroadcastToNamespace("/", "notification", notification.Payload)
			socketio.Local().Emit("notification", "test message from pg notify with data : "+*payload.ReferenceName+" = "+*payload.ReferenceId)
		}
	}()

	socketio.On("connection", func(clients ...any) {
		fmt.Println("socket io connect")
		client := clients[0].(*socket.Socket)
		fmt.Println("Client connected: ", client.Id())

		//all connected clients
		allClients2 := socketio.Of("/", nil).Sockets()
		fmt.Println("all clients 2 : ", allClients2.Len())

		//broadcast
		socketio.Local().Emit("test", "test message from "+client.Id())

		client.Join(socket.Room("pluto"))
		client.Join(socket.Room("malang"))

		client.On("/pluto/requestData", func(args ...interface{}) {
			fmt.Println(args)
			fmt.Println(args[0])
			fmt.Println(args[1])
		})

		client.On("disconnect", func(args ...interface{}) {
			// client.Disconnect(true)
			fmt.Println("Client disconnected: ", client.Id())
			client.Disconnect(true)

			allClients2 := socketio.Of("/", nil).Sockets()
			fmt.Println("all clients 2 : ", allClients2.Len())
		})
	})

	socketio.On("disconnect", func(clients ...interface{}) {
		client := clients[0].(*socket.Socket)
		fmt.Println("Client disconnected: ", client.Id())
		client.Disconnect(true)
	})

	socketio.Of("/connectCustomer", nil).On("connection", func(clients ...interface{}) {
		client := clients[0].(*socket.Socket)
		client.On("userId", func(args ...interface{}) {
			userIdClient := args[0].(string)
			result, err := helpers.RefreshUser(userIdClient)
			if err != nil {
				fmt.Println(err)
				client.Emit("returnData", "Failed to get user data")
			}
			client.Emit("returnData", result)
		})
	})
	app.Get("/socket.io", adaptor.HTTPHandler(socketio.ServeHandler(nil)))
	app.Post("/socket.io", adaptor.HTTPHandler(socketio.ServeHandler(nil)))

	app.Use("/ws/*", func(c *fiber.Ctx) error {
		socketio.ServeClient()
		return nil
	})
}
