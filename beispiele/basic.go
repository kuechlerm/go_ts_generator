package beispiele

const Basic_Path = "/basic"

type Basic_Request struct {
	RequiredString string `json:"requiredString" validate:"required"`
	OptionalString string `json:"optionalString"`
	RequiredInt    int    `json:"requiredInt" validate:"required"`
	OptionalInt    int    `json:"optionalInt"`
	RequiredBool   bool   `json:"requiredBool" validate:"required"`
	OptionalBool   bool   `json:"optionalBool"`
}

type Basic_Response struct {
	ResponseString string `json:"responseString" validate:"required"`
}

func IgnoreMe() {
	//
}
