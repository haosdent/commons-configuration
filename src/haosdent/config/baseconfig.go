package config

type Configer interface {
	Get(k string) (val string, err error)
	// GetBool(k string) (val bool, err error)
	// GetByte(k string) (val byte, err error)
	// GetFloat(k string) (val float32, err error)
	// GetDouble(k string) (val float64, err error)
	// GetInt(k string) (val int, err error)
	// GetUint(k string) (val uint, err error)
	// GetLong(k string) (val int64, err error)
	// GetStrArr(k string) (val []string, err error)

	// ContainsKey(k string) (ret bool)
	// IsEmpty(k string) (ret bool)
	// GetKeys() (keys []string)

	AddProp(k string, v interface{})
	// ClearProp(k string)
	// Clear()
	// SetProp(k string, v interface{})
	Save() error
}
