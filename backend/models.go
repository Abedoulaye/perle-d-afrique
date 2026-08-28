package main
import "time"

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
	PriceInCents int64 `json:"price"`
	Stock int `json:"stock"`
	Image string `json:"image"`
}

type Cart struct {
	ID int `json:"id"`
	UserID int `json:"user_id"`,
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