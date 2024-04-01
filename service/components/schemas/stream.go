package schemas

type Stream struct {
	Posts       []Post `json:"posts"`
	TotalPages  int    `json:"number_of_pages"`
	Page_number int    `json:"page_number"`
}

func (request *Stream) IsValid() bool {
	return (request.Page_number >= 0)
}
