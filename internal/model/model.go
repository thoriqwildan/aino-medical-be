package model

type WebResponse[T any] struct {
	Code        int             `json:"code"`
	Message     string          `json:"message"`
	AccessToken string          `json:"access_token,omitempty"`
	Meta        *PaginationPage `json:"meta,omitempty"`
	Data        *T              `json:"data,omitempty"`
	Errors      any             `json:"errors,omitempty"`
}

type PaginationPage struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
}

type SearchPagingQuery struct {
	SearchValue string `query:"search_value" json:"search_value"`
	Page        int    `json:"page,omitempty" validate:"omitempty,numeric"`
	Limit       int    `json:"limit,omitempty" validate:"omitempty,numeric"`
}

type PagingQuery struct {
	Page  int `json:"page,omitempty" query:"page" validate:"omitempty,numeric"`
	Limit int `json:"limit,omitempty" query:"limit" validate:"omitempty,numeric"`
}
