package models

type Author struct {
	Id   uint64 `json:"id"`
	Name string `json:"name"`
}

func GetAuthors() ([]Author, error) {
	var authors []Author

	query := `SELECT id, name FROM authors;`
	rows, err := execQuery(query)
	if err != nil {
		return authors, err
	}
	defer rows.Close()

	for rows.Next() {
		var id uint64
		var name string

		err = rows.Scan(&id, &name)
		if err != nil {
			return authors, err
		}

		author := Author{
			Id:   id,
			Name: name,
		}
		authors = append(authors, author)

	}

	return authors, nil
}
