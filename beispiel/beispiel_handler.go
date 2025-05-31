package beispiel

type BeispielAnlegen_Request struct {
	Name string `json:"name" validate:"required,min=3,max=100"`
}

type BeispielAnlegen_Response struct {
	ID string `json:"id"`
}

func BeispielAnlegen(args *BeispielAnlegen_Request) (BeispielAnlegen_Response, error) {
	return BeispielAnlegen_Response{
		ID: "12345",
	}, nil
}
