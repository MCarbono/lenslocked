package entity

type Gallery struct {
	ID     string
	UserID string
	Title  string
	Images []*Image
}

func NewGallery(ID, userID, title string) *Gallery {
	return &Gallery{
		ID:     ID,
		UserID: userID,
		Title:  title,
	}
}

func (g *Gallery) Update(title string) {
	if title != "" {
		g.Title = title
	}
}

func (g *Gallery) IsOwnedBy(userID string) bool {
	return g.UserID == userID
}

func (g *Gallery) AddImages(images []*Image) {
	g.Images = images
}
