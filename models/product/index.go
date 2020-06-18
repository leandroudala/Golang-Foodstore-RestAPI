package product

// Product stores data about product
type Product struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description",omitempty`
	Price       float32 `json:"price"`
	Image		string	`json:"image",omitempty`
}
