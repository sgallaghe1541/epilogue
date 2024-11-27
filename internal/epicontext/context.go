package epicontext

type ContextKey string

const (
	IsAuthenticatedContextKey = ContextKey("authenticatedUserID")
	EpilogueUser              = ContextKey("epilogueUser")
)
