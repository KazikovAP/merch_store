package integration_test

// import (
// 	"bytes"
// 	"database/sql"
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"net/http"
// 	"net/http/httptest"
// 	"os"
// 	"testing"

// 	"github.com/KazikovAP/merch_store/internal/handlers"
// 	"github.com/KazikovAP/merch_store/internal/middleware"
// 	"github.com/KazikovAP/merch_store/internal/repository"
// 	"github.com/KazikovAP/merch_store/internal/service"
// 	"github.com/gorilla/mux"
// )

// func setupTestServer() *mux.Router {
// 	// Получаем настройки из переменных окружения (для тестовой БД)
// 	dbHost := os.Getenv("DB_HOST")
// 	dbPort := os.Getenv("DB_PORT")
// 	dbUser := os.Getenv("DB_USER")
// 	dbPassword := os.Getenv("DB_PASSWORD")
// 	dbName := os.Getenv("DB_NAME")
// 	jwtSecret := os.Getenv("JWT_SECRET")

// 	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
// 		dbHost, dbPort, dbUser, dbPassword, dbName)

// 	db, err := sql.Open("postgres", connStr)
// 	if err != nil {
// 		panic(err)
// 	}
// 	// Для интеграционных тестов желательно использовать отдельную тестовую БД или сбрасывать данные

// 	userRepo := repository.NewUserRepository(db)
// 	txnRepo := repository.NewTransactionRepository(db)
// 	invRepo := repository.NewInventoryRepository(db)
// 	authService := service.NewAuthService(userRepo, jwtSecret)
// 	userService := service.NewUserService(userRepo, txnRepo, invRepo)
// 	coinService := service.NewCoinService(userRepo, txnRepo)
// 	purchaseService := service.NewPurchaseService(userRepo, invRepo)
// 	h := handlers.NewHandler(authService, userService, coinService, purchaseService)

// 	r := mux.NewRouter()
// 	r.HandleFunc("/api/auth", h.Auth).Methods("POST")
// 	api := r.PathPrefix("/api").Subrouter()
// 	api.Use(middleware.JwtMiddleware(jwtSecret))
// 	api.HandleFunc("/info", h.Info).Methods("GET")
// 	api.HandleFunc("/sendCoin", h.SendCoin).Methods("POST")
// 	api.HandleFunc("/buy/{item}", h.Buy).Methods("GET")

// 	return r
// }

// func TestIntegration_PurchaseMerch(t *testing.T) {
// 	router := setupTestServer()
// 	ts := httptest.NewServer(router)
// 	defer ts.Close()

// 	// Аутентифицируем пользователя (регистрация происходит автоматически)
// 	authPayload := map[string]string{
// 		"username": "integrationUser",
// 		"password": "password123",
// 	}

// 	authBody, _ := json.Marshal(authPayload)

// 	resp, err := http.Post(ts.URL+"/api/auth", "application/json", bytes.NewBuffer(authBody))
// 	if err != nil {
// 		t.Fatalf("auth request failed: %v", err)
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		t.Fatalf("auth failed, status: %d", resp.StatusCode)
// 	}

// 	var authResp map[string]string
// 	bodyBytes, _ := ioutil.ReadAll(resp.Body)
// 	json.Unmarshal(bodyBytes, &authResp)
// 	token := authResp["token"]

// 	// Покупка товара
// 	client := &http.Client{}
// 	req, _ := http.NewRequest("GET", ts.URL+"/api/buy/t-shirt", http.NoBody)
// 	req.Header.Set("Authorization", "Bearer "+token)

// 	resp, err = client.Do(req)
// 	if err != nil {
// 		t.Fatalf("buy request failed: %v", err)
// 	}

// 	if resp.StatusCode != http.StatusOK {
// 		t.Fatalf("buy request failed, status: %d", resp.StatusCode)
// 	}
// }

// func TestIntegration_SendCoin(t *testing.T) {
// 	router := setupTestServer()
// 	ts := httptest.NewServer(router)
// 	defer ts.Close()

// 	// Аутентификация двух пользователей
// 	user1 := map[string]string{
// 		"username": "user1",
// 		"password": "password1",
// 	}
// 	user2 := map[string]string{
// 		"username": "user2",
// 		"password": "password2",
// 	}

// 	// Аутентифицируем user1
// 	body1, _ := json.Marshal(user1)

// 	resp, err := http.Post(ts.URL+"/api/auth", "application/json", bytes.NewBuffer(body1))
// 	if err != nil {
// 		t.Fatalf("auth user1 failed: %v", err)
// 	}

// 	defer resp.Body.Close()

// 	var authResp1 map[string]string

// 	bodyBytes, _ := ioutil.ReadAll(resp.Body)

// 	json.Unmarshal(bodyBytes, &authResp1)

// 	token1 := authResp1["token"]

// 	// Аутентифицируем user2
// 	body2, _ := json.Marshal(user2)

// 	resp, err = http.Post(ts.URL+"/api/auth", "application/json", bytes.NewBuffer(body2))
// 	if err != nil {
// 		t.Fatalf("auth user2 failed: %v", err)
// 	}

// 	defer resp.Body.Close()

// 	var authResp2 map[string]string

// 	bodyBytes, _ = ioutil.ReadAll(resp.Body)

// 	json.Unmarshal(bodyBytes, &authResp2)

// 	// user1 передаёт монеты user2
// 	sendPayload := map[string]interface{}{
// 		"toUser": "user2",
// 		"amount": 100,
// 	}
// 	sendBody, _ := json.Marshal(sendPayload)
// 	client := &http.Client{}
// 	req, _ := http.NewRequest("POST", ts.URL+"/api/sendCoin", bytes.NewBuffer(sendBody))
// 	req.Header.Set("Content-Type", "application/json")
// 	req.Header.Set("Authorization", "Bearer "+token1)
// 	resp, err = client.Do(req)
// 	if err != nil {
// 		t.Fatalf("sendCoin request failed: %v", err)
// 	}
// 	if resp.StatusCode != http.StatusOK {
// 		t.Fatalf("sendCoin failed, status: %d", resp.StatusCode)
// 	}
// }
