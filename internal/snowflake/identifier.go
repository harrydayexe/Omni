package snowflake

// Identifier is an interface for types that have a snowflake id.
type Identifier interface {
	// Id returns the snowflake id of the entity.
	Id() Snowflake
}
