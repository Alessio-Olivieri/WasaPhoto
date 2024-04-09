package schemas

type Stream struct {
	TotalPages  int    `json:"number_of_pages"`
	Page_number int    `json:"page_number"`
	Posts       []Post `json:"posts"`
}

func (request *Stream) IsValid() bool {
	return (request.Page_number >= 0)
}
