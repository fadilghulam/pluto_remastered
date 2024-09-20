package db

import (
	// "devecode_app/models"
	"fmt"
	"log"

	// "go_sales_api/models"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() {
	godotenv.Load()

	username := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	database := os.Getenv("POSTGRES_DB")

	// dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, database)
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", host, username, password, database, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
		// Logger: logger.New(&logrusWriter{Logger: log}, logger.Config{
		// 	SlowThreshold:             time.Second * 5,
		// 	Colorful:                  false,
		// 	IgnoreRecordNotFoundError: true,
		// 	ParameterizedQueries:      true,
		// 	LogLevel:                  logger.Info,
		// }),
	})

	// DbHost := os.Getenv("MYSQL_HOST")
	// DbName := os.Getenv("MYSQL_DBNAME")
	// DbUsername := os.Getenv("MYSQL_USER")
	// DbPassword := os.Getenv("MYSQL_PASSWORD")

	// connection := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", DbUsername, DbPassword, DbHost, DbName)
	// dbConnection, err := gorm.Open(mysql.Open(connection), &gorm.Config{})

	if err != nil {
		panic("connection failed to the database ")
	}
	DB = db
	fmt.Println("db connected successfully")

	// go GenerateStruct(db)
	GenerateStruct(db)

	// AutoMigrate(db)
	//if err := DB.AutoMigrate(&models.Cashier{}, &models.Category{}, &models.Payment{}, &models.PaymentType{}, &models.Product{}, &models.Discount{}, &models.Order{}).Error; err != nil {
	//	log.Fatalf("Migration failed %v", err)
	//}

}

func InitOauth() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func GenerateStruct(db *gorm.DB) *gorm.DB {

	mydir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}

	g := gen.NewGenerator(gen.Config{
		OutPath: mydir + "/internal/generated",
		// OutPath: "../golang-api/internal/generated",
		Mode: gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
	})

	g.UseDB(db)

	g.ApplyBasic(
		// Generate struct `User` based on table `users`
		// g.GenerateModel("otp"),
		g.GenerateModel("customer_type_request"),
		g.GenerateModel("customer_access_visit_extra"),
		g.GenerateModel("rute_move_request"),
		g.GenerateModel("customer_move_request"),
		g.GenerateModel("customer_access"),
		g.GenerateModel("salesman_access"),
		g.GenerateModel("salesman_request_so"),
		g.GenerateModel("salesman_request"),
		g.GenerateModel("salesman_access_kunjungan"),

	// 	// Generate struct `Employee` based on table `users`
	//    g.GenerateModelAs("users", "Employee"),

	)

	// g.ApplyBasic(
	// 	g.GenerateAllTable()...,
	// )

	g.Execute()

	return nil
}
