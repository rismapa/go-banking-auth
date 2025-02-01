package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/rismapa/go-banking-auth/routes"
)

func main() {
	routes.StartServer()
}
