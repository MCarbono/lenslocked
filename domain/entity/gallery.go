package entity

type Gallery struct {
	ID     string
	UserID string
	Title  string
}

func NewGallery(ID, userID, title string) *Gallery {
	return &Gallery{
		ID:     ID,
		UserID: userID,
		Title:  title,
	}
}
