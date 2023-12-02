package model

type LatLngLiteral struct {
	Login          string `json:"login"`
	HashedPassword []byte `json:"-"`

	Email string `json:"email"`
}
