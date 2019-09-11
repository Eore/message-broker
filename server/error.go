package server

const (
	ErrorNull = iota
	ErrorClientExist
	ErrorJoinFirst
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e Error) New(id int) Error {
	err := Error{Code: id}
	switch id {
	case ErrorNull:
		err.Message = ""
	case ErrorClientExist:
		err.Message = "client exists"
	case ErrorJoinFirst:
		err.Message = "join first"
	}
	return err
}
