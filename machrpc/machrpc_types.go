package machrpc

type TableType int

const (
	LogTableType      TableType = iota + 0
	FixedTableType              = 1
	VolatileTableType           = 3
	LookupTableType             = 4
	KeyValueTableType           = 5
	TagTableType                = 6
)

type ColumnType int

const (
	Int16ColumnType    ColumnType = iota + 4
	Uint16ColumnType              = 104
	Int32ColumnType               = 8
	Uint32ColumnType              = 108
	Int64ColumnType               = 12
	Uint64ColumnType              = 112
	Float32ColumnType             = 16
	Float64ColumnType             = 20
	VarcharColumnType             = 5
	TextColumnType                = 49
	ClobColumnType                = 53
	BlobColumnType                = 57
	BinaryColumnType              = 97
	DatetimeColumnType            = 6
	IpV4ColumnType                = 32
	IpV6ColumnType                = 36
)
