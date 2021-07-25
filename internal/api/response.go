package api

import "fmt"

type APIResponseError struct {
	Status string `json:"status"`
	Code   string `json:"code"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

type APIResponseBody struct {
	Meta   map[string]interface{} `json:"meta"`
	Data   interface{}            `json:"data"`
	Errors []APIResponseError     `json:"errors"`
}

type APIResponse struct {
	StatusCode int
	Headers    map[string]string
	Body       *APIResponseBody
}

type APIResponseBuilder struct {
	Resp *APIResponse
}

func NewAPIResponseBuilder() *APIResponseBuilder {
	return &APIResponseBuilder{
		Resp: &APIResponse{
			Headers: map[string]string{},
			Body: &APIResponseBody{
				Meta:   map[string]interface{}{},
				Data:   nil,
				Errors: []APIResponseError{},
			},
		},
	}
}

func (builder *APIResponseBuilder) StatusCode(code int) *APIResponseBuilder {
	builder.Resp.StatusCode = code
	return builder
}

func (builder *APIResponseBuilder) Headers(headers map[string]string) *APIResponseBuilder {
	for k, v := range headers {
		builder.Resp.Headers[k] = v
	}
	return builder
}

func (builder *APIResponseBuilder) Header(name string, value string) *APIResponseBuilder {
	builder.Resp.Headers[name] = value
	return builder
}

func (builder *APIResponseBuilder) Data(data interface{}) *APIResponseBuilder {
	builder.Resp.Body.Data = data
	return builder
}

func (builder *APIResponseBuilder) Meta(key string, meta interface{}) *APIResponseBuilder {
	builder.Resp.Body.Meta[key] = meta
	return builder
}

func (builder *APIResponseBuilder) Error(err Errorer) *APIResponseBuilder {
	apiErrors := err.APIErrors()
	for _, apiError := range apiErrors {
		// Set status code if higher than existing HTTP status
		if int(apiError.Status()) > builder.Resp.StatusCode {
			builder.StatusCode(int(apiError.Status()))
		}
		// add Error to error array
		builder.Resp.Body.Errors = append(builder.Resp.Body.Errors, APIResponseError{
			Status: fmt.Sprint(apiError.Status()),
			Code:   string(apiError.Code()),
			Title:  apiError.Title(),
			Detail: apiError.Detail(),
		})
	}
	return builder
}

func (builder *APIResponseBuilder) Errors(errors []Errorer) *APIResponseBuilder {
	for _, err := range errors {
		builder.Error(err)
	}
	return builder
}

func (builder *APIResponseBuilder) Build() *APIResponse {
	return builder.Resp
}
