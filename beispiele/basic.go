package beispiele

const Eins_Path = "/eins"

type Eins_Request struct {
	RequiredString string `json:"requiredString" validate:"required"`
	OptionalString string `json:"optionalString"`
	RequiredInt    int    `json:"requiredInt" validate:"required"`
	OptionalInt    int    `json:"optionalInt"`
	RequiredBool   bool   `json:"requiredBool" validate:"required"`
	OptionalBool   bool   `json:"optionalBool"`
}

type Eins_Response struct {
	ResponseString string `json:"responseString" validate:"required"`
}

const Zwei_Path = "/zwei"

type Zwei_Request struct {
	OptionalString string `json:"optionalString"`
}

type Zwei_Response struct {
	ResponseString string `json:"responseString" validate:"required"`
}

func IgnoreMe() {
	//
}
