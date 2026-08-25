package main

import (
	"fmt"
	"log"
	"net/http"

	"api-gateway/internal/client"
	"api-gateway/internal/config"
	httpHandler "api-gateway/internal/handler/http"
	"api-gateway/internal/router"
)

func main() {
	cfg := config.LoadConfig()

	// Microservice client stubs
	authStub := &client.AuthClientDirectStub{}
	userStub := &client.UserClientDirectStub{}
	attendanceStub := &client.AttendanceClientDirectStub{}

	authHTTP := httpHandler.NewAuthHTTPHandler(authStub)
	userHTTP := httpHandler.NewUserHTTPHandler(userStub)
	attendanceHTTP := httpHandler.NewAttendanceHTTPHandler(attendanceStub, userStub)

	r := router.NewRouter(authHTTP, userHTTP, attendanceHTTP)

	fmt.Printf("===================================================\n")
	fmt.Printf("🚀 LOPI-Q API GATEWAY RUNNING ON http://localhost:%s\n", cfg.Port)
	fmt.Printf("📍 Posko 112 Geofence Radius Limit: 2.0 Meters\n")
	fmt.Printf("===================================================\n")

	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}
