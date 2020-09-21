package config

type IChain interface {
	afterSetValue()
	afterSetStringValue()
	afterSetEmpty()
	trySetStringValue(string)
	isEmpty(value interface{}) bool
}
