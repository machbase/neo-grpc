// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.12
// source: schedule.proto

package schedule

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	Management_ListSchedule_FullMethodName  = "/schedule.Management/ListSchedule"
	Management_AddSchedule_FullMethodName   = "/schedule.Management/AddSchedule"
	Management_GetSchedule_FullMethodName   = "/schedule.Management/GetSchedule"
	Management_DelSchedule_FullMethodName   = "/schedule.Management/DelSchedule"
	Management_StartSchedule_FullMethodName = "/schedule.Management/StartSchedule"
	Management_StopSchedule_FullMethodName  = "/schedule.Management/StopSchedule"
)

// ManagementClient is the client API for Management service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ManagementClient interface {
	ListSchedule(ctx context.Context, in *ListScheduleRequest, opts ...grpc.CallOption) (*ListScheduleResponse, error)
	AddSchedule(ctx context.Context, in *AddScheduleRequest, opts ...grpc.CallOption) (*AddScheduleResponse, error)
	GetSchedule(ctx context.Context, in *GetScheduleRequest, opts ...grpc.CallOption) (*GetScheduleResponse, error)
	DelSchedule(ctx context.Context, in *DelScheduleRequest, opts ...grpc.CallOption) (*DelScheduleResponse, error)
	StartSchedule(ctx context.Context, in *StartScheduleRequest, opts ...grpc.CallOption) (*StartScheduleResponse, error)
	StopSchedule(ctx context.Context, in *StopScheduleRequest, opts ...grpc.CallOption) (*StopScheduleResponse, error)
}

type managementClient struct {
	cc grpc.ClientConnInterface
}

func NewManagementClient(cc grpc.ClientConnInterface) ManagementClient {
	return &managementClient{cc}
}

func (c *managementClient) ListSchedule(ctx context.Context, in *ListScheduleRequest, opts ...grpc.CallOption) (*ListScheduleResponse, error) {
	out := new(ListScheduleResponse)
	err := c.cc.Invoke(ctx, Management_ListSchedule_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *managementClient) AddSchedule(ctx context.Context, in *AddScheduleRequest, opts ...grpc.CallOption) (*AddScheduleResponse, error) {
	out := new(AddScheduleResponse)
	err := c.cc.Invoke(ctx, Management_AddSchedule_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *managementClient) GetSchedule(ctx context.Context, in *GetScheduleRequest, opts ...grpc.CallOption) (*GetScheduleResponse, error) {
	out := new(GetScheduleResponse)
	err := c.cc.Invoke(ctx, Management_GetSchedule_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *managementClient) DelSchedule(ctx context.Context, in *DelScheduleRequest, opts ...grpc.CallOption) (*DelScheduleResponse, error) {
	out := new(DelScheduleResponse)
	err := c.cc.Invoke(ctx, Management_DelSchedule_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *managementClient) StartSchedule(ctx context.Context, in *StartScheduleRequest, opts ...grpc.CallOption) (*StartScheduleResponse, error) {
	out := new(StartScheduleResponse)
	err := c.cc.Invoke(ctx, Management_StartSchedule_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *managementClient) StopSchedule(ctx context.Context, in *StopScheduleRequest, opts ...grpc.CallOption) (*StopScheduleResponse, error) {
	out := new(StopScheduleResponse)
	err := c.cc.Invoke(ctx, Management_StopSchedule_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ManagementServer is the server API for Management service.
// All implementations must embed UnimplementedManagementServer
// for forward compatibility
type ManagementServer interface {
	ListSchedule(context.Context, *ListScheduleRequest) (*ListScheduleResponse, error)
	AddSchedule(context.Context, *AddScheduleRequest) (*AddScheduleResponse, error)
	GetSchedule(context.Context, *GetScheduleRequest) (*GetScheduleResponse, error)
	DelSchedule(context.Context, *DelScheduleRequest) (*DelScheduleResponse, error)
	StartSchedule(context.Context, *StartScheduleRequest) (*StartScheduleResponse, error)
	StopSchedule(context.Context, *StopScheduleRequest) (*StopScheduleResponse, error)
	mustEmbedUnimplementedManagementServer()
}

// UnimplementedManagementServer must be embedded to have forward compatible implementations.
type UnimplementedManagementServer struct {
}

func (UnimplementedManagementServer) ListSchedule(context.Context, *ListScheduleRequest) (*ListScheduleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListSchedule not implemented")
}
func (UnimplementedManagementServer) AddSchedule(context.Context, *AddScheduleRequest) (*AddScheduleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddSchedule not implemented")
}
func (UnimplementedManagementServer) GetSchedule(context.Context, *GetScheduleRequest) (*GetScheduleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSchedule not implemented")
}
func (UnimplementedManagementServer) DelSchedule(context.Context, *DelScheduleRequest) (*DelScheduleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DelSchedule not implemented")
}
func (UnimplementedManagementServer) StartSchedule(context.Context, *StartScheduleRequest) (*StartScheduleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StartSchedule not implemented")
}
func (UnimplementedManagementServer) StopSchedule(context.Context, *StopScheduleRequest) (*StopScheduleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StopSchedule not implemented")
}
func (UnimplementedManagementServer) mustEmbedUnimplementedManagementServer() {}

// UnsafeManagementServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ManagementServer will
// result in compilation errors.
type UnsafeManagementServer interface {
	mustEmbedUnimplementedManagementServer()
}

func RegisterManagementServer(s grpc.ServiceRegistrar, srv ManagementServer) {
	s.RegisterService(&Management_ServiceDesc, srv)
}

func _Management_ListSchedule_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListScheduleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ManagementServer).ListSchedule(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Management_ListSchedule_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ManagementServer).ListSchedule(ctx, req.(*ListScheduleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Management_AddSchedule_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddScheduleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ManagementServer).AddSchedule(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Management_AddSchedule_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ManagementServer).AddSchedule(ctx, req.(*AddScheduleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Management_GetSchedule_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetScheduleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ManagementServer).GetSchedule(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Management_GetSchedule_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ManagementServer).GetSchedule(ctx, req.(*GetScheduleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Management_DelSchedule_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DelScheduleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ManagementServer).DelSchedule(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Management_DelSchedule_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ManagementServer).DelSchedule(ctx, req.(*DelScheduleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Management_StartSchedule_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StartScheduleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ManagementServer).StartSchedule(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Management_StartSchedule_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ManagementServer).StartSchedule(ctx, req.(*StartScheduleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Management_StopSchedule_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StopScheduleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ManagementServer).StopSchedule(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Management_StopSchedule_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ManagementServer).StopSchedule(ctx, req.(*StopScheduleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Management_ServiceDesc is the grpc.ServiceDesc for Management service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Management_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "schedule.Management",
	HandlerType: (*ManagementServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListSchedule",
			Handler:    _Management_ListSchedule_Handler,
		},
		{
			MethodName: "AddSchedule",
			Handler:    _Management_AddSchedule_Handler,
		},
		{
			MethodName: "GetSchedule",
			Handler:    _Management_GetSchedule_Handler,
		},
		{
			MethodName: "DelSchedule",
			Handler:    _Management_DelSchedule_Handler,
		},
		{
			MethodName: "StartSchedule",
			Handler:    _Management_StartSchedule_Handler,
		},
		{
			MethodName: "StopSchedule",
			Handler:    _Management_StopSchedule_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "schedule.proto",
}
