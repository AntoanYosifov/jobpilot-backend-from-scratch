package model

type Job struct {
	ID      int    `gorm:"primaryKey" json:"id"`
	Title   string `json:"title"`
	Company string `json:"company"`
	Status  string `json:"status"`
	Date    string `json:"date"`
}
