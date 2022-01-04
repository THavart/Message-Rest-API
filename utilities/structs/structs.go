package structs

// Message struct (Model)
type Message struct {
	ID        string  `json:"id"`
	Content   string  `json:"content"`
	Author    *Author `json:"author"`
	Timestamp string  `json:""`
}

// Author struct
type Author struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}
