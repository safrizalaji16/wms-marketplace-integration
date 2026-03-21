package dto

type Response[T any] struct {
	Message string `json:"message"`
	Data    T      `json:"data"`
}

func CreateResponseError(message string) Response[any] {
	return Response[any]{
		Message: message,
		Data:    "",
	}
}

func CreateResponseSuccess[T any](data T) Response[T] {
	return Response[T]{
		Message: "Success",
		Data:    data,
	}
}

func CreateResponseErrorData(message string, data map[string]string) Response[map[string]string] {
	return Response[map[string]string]{
		Message: message,
		Data:    data,
	}
}
