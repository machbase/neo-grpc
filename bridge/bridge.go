package bridge

import (
	"fmt"
	"net"
	"time"
)

func ConvertToDatum(arr ...any) ([]*Datum, error) {
	ret := make([]*Datum, len(arr))
	for i := range arr {
		switch v := arr[i].(type) {
		case int32:
			ret[i] = &Datum{Value: &Datum_VInt32{VInt32: v}}
		case uint32:
			ret[i] = &Datum{Value: &Datum_VUint32{VUint32: v}}
		case int64:
			ret[i] = &Datum{Value: &Datum_VInt64{VInt64: v}}
		case uint64:
			ret[i] = &Datum{Value: &Datum_VUint64{VUint64: v}}
		case float32:
			ret[i] = &Datum{Value: &Datum_VFloat{VFloat: v}}
		case float64:
			ret[i] = &Datum{Value: &Datum_VDouble{VDouble: v}}
		case string:
			ret[i] = &Datum{Value: &Datum_VString{VString: v}}
		case bool:
			ret[i] = &Datum{Value: &Datum_VBool{VBool: v}}
		case []byte:
			ret[i] = &Datum{Value: &Datum_VBytes{VBytes: v}}
		case net.IP:
			ret[i] = &Datum{Value: &Datum_VIp{VIp: v.String()}}
		case time.Time:
			ret[i] = &Datum{Value: &Datum_VTime{VTime: v.UnixNano()}}
		default:
			if v == nil {
				ret[i] = &Datum{Value: &Datum_VNull{VNull: true}}
			} else {
				return nil, fmt.Errorf("ConvertToDatum() does not support %T", v)
			}
		}
	}
	return ret, nil
}

func ConvertFromDatum(arr ...*Datum) ([]any, error) {
	ret := make([]any, len(arr))
	for i := range arr {
		switch v := arr[i].Value.(type) {
		case *Datum_VInt32:
			ret[i] = v.VInt32
		case *Datum_VUint32:
			ret[i] = v.VUint32
		case *Datum_VInt64:
			ret[i] = v.VInt64
		case *Datum_VUint64:
			ret[i] = v.VUint64
		case *Datum_VFloat:
			ret[i] = v.VFloat
		case *Datum_VDouble:
			ret[i] = v.VDouble
		case *Datum_VString:
			ret[i] = v.VString
		case *Datum_VBool:
			ret[i] = v.VBool
		case *Datum_VBytes:
			ret[i] = v.VBytes
		case *Datum_VIp:
			ret[i] = net.ParseIP(v.VIp)
		case *Datum_VTime:
			ret[i] = time.Unix(0, v.VTime)
		case *Datum_VNull:
			ret[i] = nil
		default:
			return nil, fmt.Errorf("ConvertFromDatum() does not support %T", v)
		}
	}
	return ret, nil
}
