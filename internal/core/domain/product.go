package domain

// {
// 	"product_id": "987",
// 	"name": "Wireless Mouse",
// 	"description": "A smooth and reliable wireless mouse",
// 	"brand": "TechBrand",
// 	"category": "Electronics",
// 	"price": 25.00,
// 	"stock": 100,
// 	"specifications": {
// 	  "color": "Black",
// 	  "weight": "150g",
// 	  "dimensions": "10x5x3 cm"
// 	},
// 	"review_ids": ["rev001", "rev002"],
// 	"rating": 4.5,
// 	"images": [
// 	  "image_url_1.jpg",
// 	  "image_url_2.jpg"
// 	]
//   }

type Product struct {
	Id             string            `json:"id"`
	Name           string            `json:"name"`
	Description    string            `json:"description"`
	Brand          string            `json:"brand"`
	Category       string            `json:"category"`
	SubCategory    string            `json:"sub_category"`
	BasePrice      float64           `json:"price"`
	Variations     []Variation       `json:"variations"`
	Specifications map[string]string `json:"specifications"`
	ReviewIDs      []string          `json:"review_ids"`
	Rating         float64           `json:"rating"`
	Images         []string          `json:"images"`
}

type Variation struct {
	Sku   string
	Stock int
	Size  int
	Color string
	Price float64
}
