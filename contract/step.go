package contract

type Step interface {
	Execute(context ExecutionContext) error
}
