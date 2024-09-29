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
var (
	ErrInvalidProduct   = "invalid product"
	ErrInvalidVariation = "invalid variation"
)

type Product struct {
	ID          string `json:"product_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Brand       string `json:"brand"`
	Category    string `json:"category"`
	// SubCategory    string            `json:"sub_category"`
	Variations     []Variation       `json:"variations"`
	Specifications map[string]string `json:"specifications"`
	ReviewIDs      []string          `json:"review_ids"`
	Rating         float64           `json:"rating"`
	Images         []string          `json:"images"`
}

type Variation struct {
	Sku   string
	Stock int
	Size  string
	Color string
	Price float64
}

func (p *Product) IsCanCreate() bool {
	if p.Name == "" {
		return false
	}
	if p.Category == "" {
		return false
	}
	if p.Brand == "" {
		return false
	}
	if len(p.Variations) == 0 {
		return false
	}
	if len(p.Specifications) == 0 {
		return false
	}
	if p.Description == "" {
		return false
	}
	if len(p.Images) == 0 {
		return false
	}
	for _, v := range p.Variations {
		if !v.IsCanAdd() {
			return false
		}
	}
	return true
}

func (v *Variation) IsCanAdd() bool {
	if v.Sku == "" {
		return false
	}
	if v.Size == "" {
		return false
	}
	if v.Color == "" {
		return false
	}
	if v.Price == 0 {
		return false
	}
	return true
}
