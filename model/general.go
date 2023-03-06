package model

// GenericResponse defines the generic REST response of a REST call
type GenericResponse struct {
	Code    int         `json:"code" example:"200" swaggertype:"integer"`
	Message string      `json:"message" swaggertype:"string" example:"success"`
	Data    interface{} `json:"data,omitempty" swaggertype:"object"`
	Error   interface{} `json:"error,omitempty" swaggertype:"object"`
	Count   int64       `json:"count,omitempty" swaggertype:"integer" example:"10"`
}
