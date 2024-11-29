package main

import (
	"context"
	"dbConnect-go-demo/config"
	"dbConnect-go-demo/database"
	"dbConnect-go-demo/global"
	"dbConnect-go-demo/service"
	"fmt"
	"log"
)

func main() {
	// 加载配置
	dbcfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化数据库连接
	err = database.InitDB(dbcfg)
	if err != nil {
		return
	}

	// 使用全局管理器获取数据库实例
	dbManager := database.GetManager()
	// 关闭所有数据库连接
	defer func(dbManager *database.DBManager) {
		err := dbManager.CloseConnections()
		if err != nil {
			log.Fatalf("Failed to close database connections: %v", err)
		}
		fmt.Println("All-zuihou db connections closed successfully")
	}(dbManager)

	//使用 user 测试
	userService := service.NewUserService()

	//redisDB-client
	fmt.Println("=== redisDB Demo start===")
	redisDB, err := dbManager.GetRedisDB()
	if err != nil {
		return
	}
	ctx := context.Background()
	redisDB.Set(ctx, "name", "Alice-redis", 0)
	val, err := global.RedisDb.Get(ctx, "name").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("name", val)
	fmt.Println("=== redisDB Demo end===")

	//mysqlDB-client
	fmt.Println("=== gormMysqlDB Demo start===")
	gormMysqlDB, err := dbManager.GetGormDB("mysql")
	if err != nil {
		log.Fatalf("Mysql Failed to get Mysql database connection: %v", err)
	}
	// 演示增删改查操作--用户
	fmt.Println("===Mysql User CRUD Demo ===")
	err = userService.CreateUser(gormMysqlDB, "Alice-mysql", "alice@mysql.com", "password123")
	if err != nil {
		log.Fatalf("Mysql Failed to create user: %v", err)
	}
	fmt.Println("Mysql User created successfully")

	userMysql, err := userService.GetUser(gormMysqlDB, 1)
	if err != nil {
		log.Fatalf("Mysql Failed to get user: %v", err)
	}
	fmt.Printf("Mysql User fetched: %+v\n", userMysql)
	fmt.Println("=== gormMysqlDB Demo end===")

	//postgresDB-client
	fmt.Println("=== gormPostgresDB Demo start===")
	gormPostgresDB, err := dbManager.GetGormDB("postgres")
	if err != nil {
		log.Fatalf("Postgres Failed to get Postgres database connection: %v", err)
	}
	// 演示增删改查操作--用户
	fmt.Println("===Postgres User CRUD Demo ===")
	err = userService.CreateUser(gormPostgresDB, "Alice-postgres", "alice@postgres.com", "password123")
	if err != nil {
		log.Fatalf("Postgres Failed to create user: %v", err)
	}
	fmt.Println("Postgres User created successfully")

	userPostgres, err := userService.GetUser(gormPostgresDB, 1)
	if err != nil {
		log.Fatalf("Postgres Failed to get user: %v", err)
	}
	fmt.Printf("Postgres User fetched: %+v\n", userPostgres)
	fmt.Println("=== gormPostgresDB Demo end===")

	// gormDmDB-client
	fmt.Println("=== gormDmDB Demo start===")
	gormDmDB, err := dbManager.GetGormDB("dm")
	if err != nil {
		log.Fatalf("DM Failed to get DM database connection: %v", err)
	}
	// 演示增删改查操作--用户
	fmt.Println("=== DM User CRUD Demo ===")
	err = userService.CreateUser(gormDmDB, "Alice-dm", "alice@dm.com", "password123")
	if err != nil {
		log.Fatalf("DM Failed to create user: %v", err)
	}
	fmt.Println("DM User created successfully")

	userDm, err := userService.GetUser(gormDmDB, 2)
	if err != nil {
		log.Fatalf("DM Failed to get user: %v", err)
	}
	fmt.Printf("DM User fetched: %+v\n", userDm)
	fmt.Println("=== gormDmDB Demo end===")

	// gormKingBaseDB-client
	fmt.Println("=== gormKingBaseDB Demo start===")
	gormKingBaseDB, err := dbManager.GetGormDB("kingbase")
	if err != nil {
		log.Fatalf("KingBase Failed to get KingBase database connection: %v", err)
	}
	// 演示增删改查操作--用户
	fmt.Println("===KingBase User CRUD Demo ===")
	err = userService.CreateUser(gormKingBaseDB, "Alice-kingbase", "alice@kingbase.com", "password123")
	if err != nil {
		log.Fatalf("KingBase Failed to create user: %v", err)
	}
	fmt.Println("KingBase User created successfully")

	userKingBase, err := userService.GetUser(gormKingBaseDB, 1)
	if err != nil {
		log.Fatalf("KingBase Failed to get user: %v", err)
	}
	fmt.Printf("KingBase User fetched: %+v\n", userKingBase)
	fmt.Println("=== gormKingBaseDB Demo end===")
}
