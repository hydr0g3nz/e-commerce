package dto

type ProductListPage struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Category    string `json:"category"`
	VariantsNum int    `json:"variants_num"`
	Price       int    `json:"price"`
	Sale        int    `json:"sale"`
	Image1      string `json:"image1"`
	Image2      string `json:"image2"`
}
