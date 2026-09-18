package main
import (
	"time"
	"sync"
)

type User struct {
    ID       int    `json:"id"`
    Email    string `json:"email,omitempty"`
    Password string `json:"password,omitempty"` // omitempty so it's not returned in JSON
	CreatedAt time.Time `json:"created_at"`
	Role string `json:"role"`
}

/* a common programming pattern is known as "storing money as integers", floats have percision issues when in arithmetic situations due to how computers store them in binary. int64 over just int since int caps at about 2 billion or 21 million in this case after decimal points, 64 bits leave us in the quadrillions area.*/

type Product struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	PriceInCents int64 `json:"price_cents"`
	Stock int `json:"stock"`
	Image string `json:"image"`
}

type Cart struct {
	ID int `json:"id"`
	UserID int `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type CartItem struct {
	ID int `json:"id"` // id is here for uniquely identifying individual items quickly
	CartID int `json:"cart_id"`
	ProductID int `json:"product_id"`
	Quantity int `json:"quantity"`
}

type Order struct {
	ID int `json:"id"`
	UserID int `json:"user_id"`
	Status string `json:"status"`
	TotalInCents int64 `json:"total"`
	CreatedAt time.Time `json:"created_at"`
}

type OrderItem struct {
	ID int `json:"id"`
	OrderID int `json:"order_id"`
	ProductID int `json:"product_id"`
	Quantity int `json:"quantity"`
	PriceInCentsAtPurchase int64 `json:"price_at_purchase"`
}

type CartDetail struct {
	Item CartItem
	Prod Product
}

type rateLimiter struct {
    mu       sync.Mutex // used to protect shared data from being accessed by multiple goroutines at the same time. if you  try to read/write the visitors map at the same time, you get a data race — which can crash your program or corrupt data. A mutex lets you say "only one goroutine at a time can touch this."
    visitors map[string]*visitor // a map of ips and their count, the string key is the ip part
}

type visitor struct {
    count    int
    resetAt  time.Time
}