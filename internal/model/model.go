package model

type Book struct {
	Id      string
	Name    string
	Authors []Author
}

type Author struct {
	Id      string
	Name    string
	Surname string
}
