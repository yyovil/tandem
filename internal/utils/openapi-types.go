package utils

type Type string

const (
	TypeUnspecified Type = "TYPE_UNSPECIFIED"
	TypeString Type = "STRING"
	TypeNumber Type = "NUMBER"
	TypeInteger Type = "INTEGER"
	TypeBoolean Type = "BOOLEAN"
	TypeArray Type = "ARRAY"
	TypeObject Type = "OBJECT"
	TypeNULL Type = "NULL"
)