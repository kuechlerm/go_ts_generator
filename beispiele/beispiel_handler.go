package beispiele

const BeispielAnlegen_Path = "/beispielanlegen"

type BeispielAnlegen_Request struct {
	Name string `json:"name" validate:"required,min=3,max=100"`
}

type BeispielAnlegen_Response struct {
	ID string `json:"id"`
}

// für Test, ob funcs ignoriert werden und so soll später die Dateistruktur sein
func BeispielAnlegen(args *BeispielAnlegen_Request) (BeispielAnlegen_Response, error) {
	return BeispielAnlegen_Response{
		ID: "12345",
	}, nil
}

const BeispielAendern_Path = "/beispielaendern"

type BeispielAendern_Request struct {
	ID   string `json:"id" validate:"required"`
	Name string `json:"name" validate:"required,min=3,max=100"`
}

type BeispielAendern_Response struct {
	ID string `json:"id"`
}
