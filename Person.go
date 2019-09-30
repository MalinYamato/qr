package main

type UserId = int

type Person struct {
	Nic           string            `json:"nic"`
	FirstName     string            `json:"firstName"`
	LastName      string            `json:"lastName"`
	Password      string            `json:"password"`
	Gender        string            `json:"gender"`
	Relationship  string            `json:"relationship"`
	Long          string            `json:"long,omitempty"`
	Lat           string            `json:"lat,omitempty"`
	PictureURL    string            `json:"pictureURL,omitempty"`
	BirthDate     Date              `json:"birthDate"`
	Languages     map[string]string `json:"Languages,omitempty"`
	LanguagesList []string          `json:"LanguagesList,omitempty"`
	Profession    string            `json:"profession"`
	Description   string            `json:"description,omitempty"`
	GoogleID      string            `json:"googleId,omitempty"`
	UserID        UserId            `json:"userId"`
	Email         string            `json:"rakumail"`
	LoggedIn      bool              `json:"loggedIn"`
}
