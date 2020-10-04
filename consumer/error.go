package consumer

// Error interface for application error
type Error interface {
	error
	UserError() bool
}
