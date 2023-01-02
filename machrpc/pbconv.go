package machrpc

import (
	"fmt"
	"net"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func ConvertAnyToPb(params []any) ([]*anypb.Any, error) {
	pbparams := make([]*anypb.Any, len(params))
	var err error
	for i, p := range params {
		if p == nil {
			pbparams[i] = nil
			continue
		}
		switch v := p.(type) {
		case *int:
			pbparams[i], err = anypb.New(wrapperspb.Int32(int32(*v)))
		case int:
			pbparams[i], err = anypb.New(wrapperspb.Int32(int32(v)))
		case *int32:
			pbparams[i], err = anypb.New(wrapperspb.Int32(*v))
		case int32:
			pbparams[i], err = anypb.New(wrapperspb.Int32(v))
		case *int64:
			pbparams[i], err = anypb.New(wrapperspb.Int64(*v))
		case int64:
			pbparams[i], err = anypb.New(wrapperspb.Int64(v))
		case *float32:
			pbparams[i], err = anypb.New(wrapperspb.Float(*v))
		case float32:
			pbparams[i], err = anypb.New(wrapperspb.Float(v))
		case *float64:
			pbparams[i], err = anypb.New(wrapperspb.Double(*v))
		case float64:
			pbparams[i], err = anypb.New(wrapperspb.Double(v))
		case *string:
			pbparams[i], err = anypb.New(wrapperspb.String(*v))
		case string:
			pbparams[i], err = anypb.New(wrapperspb.String(v))
		case []byte:
			pbparams[i], err = anypb.New(wrapperspb.Bytes(v))
		case *net.IP:
			pbparams[i], err = anypb.New(wrapperspb.String(v.String()))
		case net.IP:
			pbparams[i], err = anypb.New(wrapperspb.String(v.String()))
		case *time.Time:
			pbparams[i], err = anypb.New(wrapperspb.Int64(v.UnixNano()))
		case time.Time:
			pbparams[i], err = anypb.New(wrapperspb.Int64(v.UnixNano()))
		default:
			return nil, fmt.Errorf("unsupported params[%d] type %T", i, p)
		}
		if err != nil {
			return nil, errors.Wrapf(err, "convert params[%d]", i)
		}
	}
	return pbparams, nil
}

func ConvertPbToAny(pbvals []*anypb.Any) []any {
	vals := make([]any, len(pbvals))
	for i, pbval := range pbvals {
		var value any
		switch pbval.TypeUrl {
		case "type.googleapis.com/google.protobuf.StringValue":
			var v wrapperspb.StringValue
			pbval.UnmarshalTo(&v)
			value = v.Value
		case "type.googleapis.com/google.protobuf.BoolValue":
			var v wrapperspb.BoolValue
			pbval.UnmarshalTo(&v)
			value = v.Value
		case "type.googleapis.com/google.protobuf.BytesValue":
			var v wrapperspb.BytesValue
			pbval.UnmarshalTo(&v)
			value = v.Value
		case "type.googleapis.com/google.protobuf.DoubleValue":
			var v wrapperspb.DoubleValue
			pbval.UnmarshalTo(&v)
			value = v.Value
		case "type.googleapis.com/google.protobuf.FloatValue":
			var v wrapperspb.FloatValue
			pbval.UnmarshalTo(&v)
			value = v.Value
		case "type.googleapis.com/google.protobuf.Int32Value":
			var v wrapperspb.Int32Value
			pbval.UnmarshalTo(&v)
			value = v.Value
		case "type.googleapis.com/google.protobuf.UInt32Value":
			var v wrapperspb.UInt32Value
			pbval.UnmarshalTo(&v)
			value = v.Value
		case "type.googleapis.com/google.protobuf.Int64Value":
			var v wrapperspb.Int64Value
			pbval.UnmarshalTo(&v)
			value = v.Value
		case "type.googleapis.com/google.protobuf.UInt64Value":
			var v wrapperspb.UInt64Value
			pbval.UnmarshalTo(&v)
			value = v.Value
		case "":
			value = nil
		default:
			value = pbval
		}
		vals[i] = value
	}
	return vals
}
