package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/minhmannh2001/sms/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type SMSDatabase interface {
	CreateServer(server *entity.Server) error
	ViewServer(id int) (entity.Server, error)
	ViewServers(from int, to int, perpage int, sortby string, order string, filter string) ([]entity.Server, error)
	UpdateServer(server *entity.Server) error
	DeleteServer(id int) error
	CheckServerExistence(ip string) bool
	CheckServerName(name string) bool
	GetServerByIp(ip string) entity.Server
	ReplaceConnection(dbName string)
	CloseDBConnection()
}

type smsDatabase struct {
	connection *gorm.DB
}

func NewSMSDatabase() SMSDatabase {
	err := godotenv.Load()
	if err != nil {
		godotenv.Load("./../.env")
		// panic("Failed to load env file")
	}

	var dbName, dbUser, dbPassword, dbHost, dbPort string

	dbName = os.Getenv("DBNAME")
	dbUser = os.Getenv("DB_USERNAME")
	dbPassword = os.Getenv("DB_PASSWORD")
	dbHost = os.Getenv("DB_HOST")
	dbPort = os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", dbHost, dbUser, dbPassword, dbName, dbPort)
	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln("Failed to connect to the database")
		return nil
	}
	log.Println("Connected successfully to the database")
	conn.AutoMigrate(&entity.Server{})
	return &smsDatabase{
		connection: conn,
	}
}

func (db *smsDatabase) CreateServer(server *entity.Server) error {
	err := db.connection.Create(server).Error
	if err != nil {
		return err
	}
	fmt.Println("Execute CreateServer() function from sms-database.go file")
	return nil
}

func (db *smsDatabase) ViewServer(id int) (entity.Server, error) {
	var server entity.Server

	if err := db.connection.Where("id = ?", id).First(&server).Error; err != nil {
		return server, err
	}
	fmt.Println("Execute ViewServer() function from sms-database.go file")
	return server, nil
}

func (db *smsDatabase) ViewServers(from int, to int, perpage int, sortby string, order string, filter string) ([]entity.Server, error) {
	var servers []entity.Server

	sql := "SELECT * FROM servers"

	if filter != "" {
		sql = fmt.Sprintf("%s WHERE status='%s'", sql, filter)
	}

	if !(from == 0 && to == 0 && perpage == 0) {
		if from == 0 && to != 0 {
			from = to
		}
		if from != 0 && to == 0 {
			to = from
		}
		sql = fmt.Sprintf("%s LIMIT %d OFFSET %d", sql, perpage*(to-from+1), (from-1)*perpage)
	}

	if order == "" {
		order = "asc"
	}
	if sortby != "" {
		sql = fmt.Sprintf("%s ORDER BY %s %s", sql, sortby, order)
	}

	if err := db.connection.Raw(sql).Scan(&servers).Error; err != nil {
		return nil, err
	}

	fmt.Println("Execute ViewServers() function from sms-database.go file")
	return servers, nil
}

func (db *smsDatabase) UpdateServer(server *entity.Server) error {
	id := server.Id
	err := db.connection.Where("id = ?", id).Updates(server).Error
	if err != nil {
		return err
	}
	fmt.Println("Execute UpdateServer() function from sms-database.go file")
	return nil
}

func (db *smsDatabase) DeleteServer(id int) error {
	err := db.connection.Table(entity.Server{}.TableName()).Where("id = ?", id).Delete(nil).Error
	if err != nil {
		return err
	}
	fmt.Println("Execute DeleteServer() function from sms-database.go file")
	return nil
}

func (db *smsDatabase) CheckServerExistence(ip string) bool {
	var server entity.Server
	if err := db.connection.Where("ipv4 = ?", ip).First(&server).Error; err != nil {
		return false
	}
	return true
}

func (db *smsDatabase) CheckServerName(name string) bool {
	var server entity.Server
	if err := db.connection.Where("name = ?", name).First(&server).Error; err != nil {
		return false
	}
	return true
}

func (db *smsDatabase) GetServerByIp(ip string) entity.Server {
	var server entity.Server
	db.connection.Where("ipv4 = ?", ip).First(&server)
	return server
}

func (db *smsDatabase) ReplaceConnection(dbName string) {
	err := godotenv.Load()
	if err != nil {
		godotenv.Load("./../.env")
	}

	dbUser := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", dbHost, dbUser, dbPassword, dbName, dbPort)
	new_conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln("Failed to connect to the database")
	}
	log.Println("Connected successfully to the database")
	new_conn.AutoMigrate(&entity.Server{})

	db.connection = new_conn
}

func (db *smsDatabase) CloseDBConnection() {
	dbSQL, err := db.connection.DB()
	if err != nil {
		panic("Failed to close connection from database")
	}
	dbSQL.Close()
}
