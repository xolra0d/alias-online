# proto

```go
import "github.com/xolra0d/alias-online/shared/proto/auth"
```

## Index

- [Constants](<#constants>)
- [Variables](<#variables>)
- [func RegisterAuthServiceServer\(s grpc.ServiceRegistrar, srv AuthServiceServer\)](<#RegisterAuthServiceServer>)
- [type AuthServiceClient](<#AuthServiceClient>)
  - [func NewAuthServiceClient\(cc grpc.ClientConnInterface\) AuthServiceClient](<#NewAuthServiceClient>)
- [type AuthServiceServer](<#AuthServiceServer>)
- [type LoginRequest](<#LoginRequest>)
  - [func \(\*LoginRequest\) Descriptor\(\) \(\[\]byte, \[\]int\)](<#LoginRequest.Descriptor>)
  - [func \(x \*LoginRequest\) GetLogin\(\) string](<#LoginRequest.GetLogin>)
  - [func \(x \*LoginRequest\) GetPassword\(\) string](<#LoginRequest.GetPassword>)
  - [func \(\*LoginRequest\) ProtoMessage\(\)](<#LoginRequest.ProtoMessage>)
  - [func \(x \*LoginRequest\) ProtoReflect\(\) protoreflect.Message](<#LoginRequest.ProtoReflect>)
  - [func \(x \*LoginRequest\) Reset\(\)](<#LoginRequest.Reset>)
  - [func \(x \*LoginRequest\) String\(\) string](<#LoginRequest.String>)
- [type LoginResponse](<#LoginResponse>)
  - [func \(\*LoginResponse\) Descriptor\(\) \(\[\]byte, \[\]int\)](<#LoginResponse.Descriptor>)
  - [func \(x \*LoginResponse\) GetExp\(\) int64](<#LoginResponse.GetExp>)
  - [func \(x \*LoginResponse\) GetToken\(\) string](<#LoginResponse.GetToken>)
  - [func \(\*LoginResponse\) ProtoMessage\(\)](<#LoginResponse.ProtoMessage>)
  - [func \(x \*LoginResponse\) ProtoReflect\(\) protoreflect.Message](<#LoginResponse.ProtoReflect>)
  - [func \(x \*LoginResponse\) Reset\(\)](<#LoginResponse.Reset>)
  - [func \(x \*LoginResponse\) String\(\) string](<#LoginResponse.String>)
- [type PingResponse](<#PingResponse>)
  - [func \(\*PingResponse\) Descriptor\(\) \(\[\]byte, \[\]int\)](<#PingResponse.Descriptor>)
  - [func \(x \*PingResponse\) GetOk\(\) bool](<#PingResponse.GetOk>)
  - [func \(\*PingResponse\) ProtoMessage\(\)](<#PingResponse.ProtoMessage>)
  - [func \(x \*PingResponse\) ProtoReflect\(\) protoreflect.Message](<#PingResponse.ProtoReflect>)
  - [func \(x \*PingResponse\) Reset\(\)](<#PingResponse.Reset>)
  - [func \(x \*PingResponse\) String\(\) string](<#PingResponse.String>)
- [type RegisterRequest](<#RegisterRequest>)
  - [func \(\*RegisterRequest\) Descriptor\(\) \(\[\]byte, \[\]int\)](<#RegisterRequest.Descriptor>)
  - [func \(x \*RegisterRequest\) GetLogin\(\) string](<#RegisterRequest.GetLogin>)
  - [func \(x \*RegisterRequest\) GetName\(\) string](<#RegisterRequest.GetName>)
  - [func \(x \*RegisterRequest\) GetPassword\(\) string](<#RegisterRequest.GetPassword>)
  - [func \(\*RegisterRequest\) ProtoMessage\(\)](<#RegisterRequest.ProtoMessage>)
  - [func \(x \*RegisterRequest\) ProtoReflect\(\) protoreflect.Message](<#RegisterRequest.ProtoReflect>)
  - [func \(x \*RegisterRequest\) Reset\(\)](<#RegisterRequest.Reset>)
  - [func \(x \*RegisterRequest\) String\(\) string](<#RegisterRequest.String>)
- [type RegisterResponse](<#RegisterResponse>)
  - [func \(\*RegisterResponse\) Descriptor\(\) \(\[\]byte, \[\]int\)](<#RegisterResponse.Descriptor>)
  - [func \(x \*RegisterResponse\) GetExp\(\) int64](<#RegisterResponse.GetExp>)
  - [func \(x \*RegisterResponse\) GetToken\(\) string](<#RegisterResponse.GetToken>)
  - [func \(\*RegisterResponse\) ProtoMessage\(\)](<#RegisterResponse.ProtoMessage>)
  - [func \(x \*RegisterResponse\) ProtoReflect\(\) protoreflect.Message](<#RegisterResponse.ProtoReflect>)
  - [func \(x \*RegisterResponse\) Reset\(\)](<#RegisterResponse.Reset>)
  - [func \(x \*RegisterResponse\) String\(\) string](<#RegisterResponse.String>)
- [type UnimplementedAuthServiceServer](<#UnimplementedAuthServiceServer>)
  - [func \(UnimplementedAuthServiceServer\) Login\(context.Context, \*LoginRequest\) \(\*LoginResponse, error\)](<#UnimplementedAuthServiceServer.Login>)
  - [func \(UnimplementedAuthServiceServer\) Ping\(context.Context, \*emptypb.Empty\) \(\*PingResponse, error\)](<#UnimplementedAuthServiceServer.Ping>)
  - [func \(UnimplementedAuthServiceServer\) Register\(context.Context, \*RegisterRequest\) \(\*RegisterResponse, error\)](<#UnimplementedAuthServiceServer.Register>)
- [type UnsafeAuthServiceServer](<#UnsafeAuthServiceServer>)


## Constants

<a name="AuthService_Ping_FullMethodName"></a>

```go
const (
    AuthService_Ping_FullMethodName     = "/auth.AuthService/Ping"
    AuthService_Register_FullMethodName = "/auth.AuthService/Register"
    AuthService_Login_FullMethodName    = "/auth.AuthService/Login"
)
```

## Variables

<a name="AuthService_ServiceDesc"></a>AuthService\_ServiceDesc is the grpc.ServiceDesc for AuthService service. It's only intended for direct use with grpc.RegisterService, and not to be introspected or modified \(even as a copy\)

```go
var AuthService_ServiceDesc = grpc.ServiceDesc{
    ServiceName: "auth.AuthService",
    HandlerType: (*AuthServiceServer)(nil),
    Methods: []grpc.MethodDesc{
        {
            MethodName: "Ping",
            Handler:    _AuthService_Ping_Handler,
        },
        {
            MethodName: "Register",
            Handler:    _AuthService_Register_Handler,
        },
        {
            MethodName: "Login",
            Handler:    _AuthService_Login_Handler,
        },
    },
    Streams:  []grpc.StreamDesc{},
    Metadata: "auth/auth.proto",
}
```

<a name="File_auth_auth_proto"></a>

```go
var File_auth_auth_proto protoreflect.FileDescriptor
```

<a name="RegisterAuthServiceServer"></a>
## func RegisterAuthServiceServer

```go
func RegisterAuthServiceServer(s grpc.ServiceRegistrar, srv AuthServiceServer)
```



<a name="AuthServiceClient"></a>
## type AuthServiceClient

AuthServiceClient is the client API for AuthService service.

For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.

```go
type AuthServiceClient interface {
    Ping(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*PingResponse, error)
    Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error)
    Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error)
}
```

<a name="NewAuthServiceClient"></a>
### func NewAuthServiceClient

```go
func NewAuthServiceClient(cc grpc.ClientConnInterface) AuthServiceClient
```



<a name="AuthServiceServer"></a>
## type AuthServiceServer

AuthServiceServer is the server API for AuthService service. All implementations must embed UnimplementedAuthServiceServer for forward compatibility.

```go
type AuthServiceServer interface {
    Ping(context.Context, *emptypb.Empty) (*PingResponse, error)
    Register(context.Context, *RegisterRequest) (*RegisterResponse, error)
    Login(context.Context, *LoginRequest) (*LoginResponse, error)
    // contains filtered or unexported methods
}
```

<a name="LoginRequest"></a>
## type LoginRequest



```go
type LoginRequest struct {
    Login    string `protobuf:"bytes,1,opt,name=login,proto3" json:"login,omitempty"`
    Password string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
    // contains filtered or unexported fields
}
```

<a name="LoginRequest.Descriptor"></a>
### func \(\*LoginRequest\) Descriptor

```go
func (*LoginRequest) Descriptor() ([]byte, []int)
```

Deprecated: Use LoginRequest.ProtoReflect.Descriptor instead.

<a name="LoginRequest.GetLogin"></a>
### func \(\*LoginRequest\) GetLogin

```go
func (x *LoginRequest) GetLogin() string
```



<a name="LoginRequest.GetPassword"></a>
### func \(\*LoginRequest\) GetPassword

```go
func (x *LoginRequest) GetPassword() string
```



<a name="LoginRequest.ProtoMessage"></a>
### func \(\*LoginRequest\) ProtoMessage

```go
func (*LoginRequest) ProtoMessage()
```



<a name="LoginRequest.ProtoReflect"></a>
### func \(\*LoginRequest\) ProtoReflect

```go
func (x *LoginRequest) ProtoReflect() protoreflect.Message
```



<a name="LoginRequest.Reset"></a>
### func \(\*LoginRequest\) Reset

```go
func (x *LoginRequest) Reset()
```



<a name="LoginRequest.String"></a>
### func \(\*LoginRequest\) String

```go
func (x *LoginRequest) String() string
```



<a name="LoginResponse"></a>
## type LoginResponse



```go
type LoginResponse struct {
    Token string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
    Exp   int64  `protobuf:"varint,2,opt,name=exp,proto3" json:"exp,omitempty"`
    // contains filtered or unexported fields
}
```

<a name="LoginResponse.Descriptor"></a>
### func \(\*LoginResponse\) Descriptor

```go
func (*LoginResponse) Descriptor() ([]byte, []int)
```

Deprecated: Use LoginResponse.ProtoReflect.Descriptor instead.

<a name="LoginResponse.GetExp"></a>
### func \(\*LoginResponse\) GetExp

```go
func (x *LoginResponse) GetExp() int64
```



<a name="LoginResponse.GetToken"></a>
### func \(\*LoginResponse\) GetToken

```go
func (x *LoginResponse) GetToken() string
```



<a name="LoginResponse.ProtoMessage"></a>
### func \(\*LoginResponse\) ProtoMessage

```go
func (*LoginResponse) ProtoMessage()
```



<a name="LoginResponse.ProtoReflect"></a>
### func \(\*LoginResponse\) ProtoReflect

```go
func (x *LoginResponse) ProtoReflect() protoreflect.Message
```



<a name="LoginResponse.Reset"></a>
### func \(\*LoginResponse\) Reset

```go
func (x *LoginResponse) Reset()
```



<a name="LoginResponse.String"></a>
### func \(\*LoginResponse\) String

```go
func (x *LoginResponse) String() string
```



<a name="PingResponse"></a>
## type PingResponse



```go
type PingResponse struct {
    Ok bool `protobuf:"varint,1,opt,name=ok,proto3" json:"ok,omitempty"`
    // contains filtered or unexported fields
}
```

<a name="PingResponse.Descriptor"></a>
### func \(\*PingResponse\) Descriptor

```go
func (*PingResponse) Descriptor() ([]byte, []int)
```

Deprecated: Use PingResponse.ProtoReflect.Descriptor instead.

<a name="PingResponse.GetOk"></a>
### func \(\*PingResponse\) GetOk

```go
func (x *PingResponse) GetOk() bool
```



<a name="PingResponse.ProtoMessage"></a>
### func \(\*PingResponse\) ProtoMessage

```go
func (*PingResponse) ProtoMessage()
```



<a name="PingResponse.ProtoReflect"></a>
### func \(\*PingResponse\) ProtoReflect

```go
func (x *PingResponse) ProtoReflect() protoreflect.Message
```



<a name="PingResponse.Reset"></a>
### func \(\*PingResponse\) Reset

```go
func (x *PingResponse) Reset()
```



<a name="PingResponse.String"></a>
### func \(\*PingResponse\) String

```go
func (x *PingResponse) String() string
```



<a name="RegisterRequest"></a>
## type RegisterRequest



```go
type RegisterRequest struct {
    Name     string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
    Login    string `protobuf:"bytes,2,opt,name=login,proto3" json:"login,omitempty"`
    Password string `protobuf:"bytes,3,opt,name=password,proto3" json:"password,omitempty"`
    // contains filtered or unexported fields
}
```

<a name="RegisterRequest.Descriptor"></a>
### func \(\*RegisterRequest\) Descriptor

```go
func (*RegisterRequest) Descriptor() ([]byte, []int)
```

Deprecated: Use RegisterRequest.ProtoReflect.Descriptor instead.

<a name="RegisterRequest.GetLogin"></a>
### func \(\*RegisterRequest\) GetLogin

```go
func (x *RegisterRequest) GetLogin() string
```



<a name="RegisterRequest.GetName"></a>
### func \(\*RegisterRequest\) GetName

```go
func (x *RegisterRequest) GetName() string
```



<a name="RegisterRequest.GetPassword"></a>
### func \(\*RegisterRequest\) GetPassword

```go
func (x *RegisterRequest) GetPassword() string
```



<a name="RegisterRequest.ProtoMessage"></a>
### func \(\*RegisterRequest\) ProtoMessage

```go
func (*RegisterRequest) ProtoMessage()
```



<a name="RegisterRequest.ProtoReflect"></a>
### func \(\*RegisterRequest\) ProtoReflect

```go
func (x *RegisterRequest) ProtoReflect() protoreflect.Message
```



<a name="RegisterRequest.Reset"></a>
### func \(\*RegisterRequest\) Reset

```go
func (x *RegisterRequest) Reset()
```



<a name="RegisterRequest.String"></a>
### func \(\*RegisterRequest\) String

```go
func (x *RegisterRequest) String() string
```



<a name="RegisterResponse"></a>
## type RegisterResponse



```go
type RegisterResponse struct {
    Token string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
    Exp   int64  `protobuf:"varint,2,opt,name=exp,proto3" json:"exp,omitempty"`
    // contains filtered or unexported fields
}
```

<a name="RegisterResponse.Descriptor"></a>
### func \(\*RegisterResponse\) Descriptor

```go
func (*RegisterResponse) Descriptor() ([]byte, []int)
```

Deprecated: Use RegisterResponse.ProtoReflect.Descriptor instead.

<a name="RegisterResponse.GetExp"></a>
### func \(\*RegisterResponse\) GetExp

```go
func (x *RegisterResponse) GetExp() int64
```



<a name="RegisterResponse.GetToken"></a>
### func \(\*RegisterResponse\) GetToken

```go
func (x *RegisterResponse) GetToken() string
```



<a name="RegisterResponse.ProtoMessage"></a>
### func \(\*RegisterResponse\) ProtoMessage

```go
func (*RegisterResponse) ProtoMessage()
```



<a name="RegisterResponse.ProtoReflect"></a>
### func \(\*RegisterResponse\) ProtoReflect

```go
func (x *RegisterResponse) ProtoReflect() protoreflect.Message
```



<a name="RegisterResponse.Reset"></a>
### func \(\*RegisterResponse\) Reset

```go
func (x *RegisterResponse) Reset()
```



<a name="RegisterResponse.String"></a>
### func \(\*RegisterResponse\) String

```go
func (x *RegisterResponse) String() string
```



<a name="UnimplementedAuthServiceServer"></a>
## type UnimplementedAuthServiceServer

UnimplementedAuthServiceServer must be embedded to have forward compatible implementations.

NOTE: this should be embedded by value instead of pointer to avoid a nil pointer dereference when methods are called.

```go
type UnimplementedAuthServiceServer struct{}
```

<a name="UnimplementedAuthServiceServer.Login"></a>
### func \(UnimplementedAuthServiceServer\) Login

```go
func (UnimplementedAuthServiceServer) Login(context.Context, *LoginRequest) (*LoginResponse, error)
```



<a name="UnimplementedAuthServiceServer.Ping"></a>
### func \(UnimplementedAuthServiceServer\) Ping

```go
func (UnimplementedAuthServiceServer) Ping(context.Context, *emptypb.Empty) (*PingResponse, error)
```



<a name="UnimplementedAuthServiceServer.Register"></a>
### func \(UnimplementedAuthServiceServer\) Register

```go
func (UnimplementedAuthServiceServer) Register(context.Context, *RegisterRequest) (*RegisterResponse, error)
```



<a name="UnsafeAuthServiceServer"></a>
## type UnsafeAuthServiceServer

UnsafeAuthServiceServer may be embedded to opt out of forward compatibility for this service. Use of this interface is not recommended, as added methods to AuthServiceServer will result in compilation errors.

```go
type UnsafeAuthServiceServer interface {
    // contains filtered or unexported methods
}
```

# proto

```go
import "github.com/xolra0d/alias-online/shared/proto/room_manager"
```

## Index

- [Constants](<#constants>)
- [Variables](<#variables>)
- [func RegisterRoomManagerServiceServer\(s grpc.ServiceRegistrar, srv RoomManagerServiceServer\)](<#RegisterRoomManagerServiceServer>)
- [type GetRoomWorkerRequest](<#GetRoomWorkerRequest>)
  - [func \(\*GetRoomWorkerRequest\) Descriptor\(\) \(\[\]byte, \[\]int\)](<#GetRoomWorkerRequest.Descriptor>)
  - [func \(x \*GetRoomWorkerRequest\) GetRoomId\(\) string](<#GetRoomWorkerRequest.GetRoomId>)
  - [func \(\*GetRoomWorkerRequest\) ProtoMessage\(\)](<#GetRoomWorkerRequest.ProtoMessage>)
  - [func \(x \*GetRoomWorkerRequest\) ProtoReflect\(\) protoreflect.Message](<#GetRoomWorkerRequest.ProtoReflect>)
  - [func \(x \*GetRoomWorkerRequest\) Reset\(\)](<#GetRoomWorkerRequest.Reset>)
  - [func \(x \*GetRoomWorkerRequest\) String\(\) string](<#GetRoomWorkerRequest.String>)
- [type GetRoomWorkerResponse](<#GetRoomWorkerResponse>)
  - [func \(\*GetRoomWorkerResponse\) Descriptor\(\) \(\[\]byte, \[\]int\)](<#GetRoomWorkerResponse.Descriptor>)
  - [func \(x \*GetRoomWorkerResponse\) GetWorker\(\) string](<#GetRoomWorkerResponse.GetWorker>)
  - [func \(\*GetRoomWorkerResponse\) ProtoMessage\(\)](<#GetRoomWorkerResponse.ProtoMessage>)
  - [func \(x \*GetRoomWorkerResponse\) ProtoReflect\(\) protoreflect.Message](<#GetRoomWorkerResponse.ProtoReflect>)
  - [func \(x \*GetRoomWorkerResponse\) Reset\(\)](<#GetRoomWorkerResponse.Reset>)
  - [func \(x \*GetRoomWorkerResponse\) String\(\) string](<#GetRoomWorkerResponse.String>)
- [type PingResponse](<#PingResponse>)
  - [func \(\*PingResponse\) Descriptor\(\) \(\[\]byte, \[\]int\)](<#PingResponse.Descriptor>)
  - [func \(x \*PingResponse\) GetOk\(\) bool](<#PingResponse.GetOk>)
  - [func \(\*PingResponse\) ProtoMessage\(\)](<#PingResponse.ProtoMessage>)
  - [func \(x \*PingResponse\) ProtoReflect\(\) protoreflect.Message](<#PingResponse.ProtoReflect>)
  - [func \(x \*PingResponse\) Reset\(\)](<#PingResponse.Reset>)
  - [func \(x \*PingResponse\) String\(\) string](<#PingResponse.String>)
- [type PingWorkerRequest](<#PingWorkerRequest>)
  - [func \(\*PingWorkerRequest\) Descriptor\(\) \(\[\]byte, \[\]int\)](<#PingWorkerRequest.Descriptor>)
  - [func \(x \*PingWorkerRequest\) GetLoadedRooms\(\) \[\]string](<#PingWorkerRequest.GetLoadedRooms>)
  - [func \(x \*PingWorkerRequest\) GetWorker\(\) string](<#PingWorkerRequest.GetWorker>)
  - [func \(\*PingWorkerRequest\) ProtoMessage\(\)](<#PingWorkerRequest.ProtoMessage>)
  - [func \(x \*PingWorkerRequest\) ProtoReflect\(\) protoreflect.Message](<#PingWorkerRequest.ProtoReflect>)
  - [func \(x \*PingWorkerRequest\) Reset\(\)](<#PingWorkerRequest.Reset>)
  - [func \(x \*PingWorkerRequest\) String\(\) string](<#PingWorkerRequest.String>)
- [type ProlongRoomRequest](<#ProlongRoomRequest>)
  - [func \(\*ProlongRoomRequest\) Descriptor\(\) \(\[\]byte, \[\]int\)](<#ProlongRoomRequest.Descriptor>)
  - [func \(x \*ProlongRoomRequest\) GetRoomId\(\) string](<#ProlongRoomRequest.GetRoomId>)
  - [func \(x \*ProlongRoomRequest\) GetWorker\(\) string](<#ProlongRoomRequest.GetWorker>)
  - [func \(\*ProlongRoomRequest\) ProtoMessage\(\)](<#ProlongRoomRequest.ProtoMessage>)
  - [func \(x \*ProlongRoomRequest\) ProtoReflect\(\) protoreflect.Message](<#ProlongRoomRequest.ProtoReflect>)
  - [func \(x \*ProlongRoomRequest\) Reset\(\)](<#ProlongRoomRequest.Reset>)
  - [func \(x \*ProlongRoomRequest\) String\(\) string](<#ProlongRoomRequest.String>)
- [type RegisterRoomRequest](<#RegisterRoomRequest>)
  - [func \(\*RegisterRoomRequest\) Descriptor\(\) \(\[\]byte, \[\]int\)](<#RegisterRoomRequest.Descriptor>)
  - [func \(x \*RegisterRoomRequest\) GetRoomId\(\) string](<#RegisterRoomRequest.GetRoomId>)
  - [func \(x \*RegisterRoomRequest\) GetWorker\(\) string](<#RegisterRoomRequest.GetWorker>)
  - [func \(\*RegisterRoomRequest\) ProtoMessage\(\)](<#RegisterRoomRequest.ProtoMessage>)
  - [func \(x \*RegisterRoomRequest\) ProtoReflect\(\) protoreflect.Message](<#RegisterRoomRequest.ProtoReflect>)
  - [func \(x \*RegisterRoomRequest\) Reset\(\)](<#RegisterRoomRequest.Reset>)
  - [func \(x \*RegisterRoomRequest\) String\(\) string](<#RegisterRoomRequest.String>)
- [type RegisterRoomResponse](<#RegisterRoomResponse>)
  - [func \(\*RegisterRoomResponse\) Descriptor\(\) \(\[\]byte, \[\]int\)](<#RegisterRoomResponse.Descriptor>)
  - [func \(x \*RegisterRoomResponse\) GetWorker\(\) string](<#RegisterRoomResponse.GetWorker>)
  - [func \(\*RegisterRoomResponse\) ProtoMessage\(\)](<#RegisterRoomResponse.ProtoMessage>)
  - [func \(x \*RegisterRoomResponse\) ProtoReflect\(\) protoreflect.Message](<#RegisterRoomResponse.ProtoReflect>)
  - [func \(x \*RegisterRoomResponse\) Reset\(\)](<#RegisterRoomResponse.Reset>)
  - [func \(x \*RegisterRoomResponse\) String\(\) string](<#RegisterRoomResponse.String>)
- [type ReleaseRoomRequest](<#ReleaseRoomRequest>)
  - [func \(\*ReleaseRoomRequest\) Descriptor\(\) \(\[\]byte, \[\]int\)](<#ReleaseRoomRequest.Descriptor>)
  - [func \(x \*ReleaseRoomRequest\) GetRoomId\(\) string](<#ReleaseRoomRequest.GetRoomId>)
  - [func \(x \*ReleaseRoomRequest\) GetWorker\(\) string](<#ReleaseRoomRequest.GetWorker>)
  - [func \(\*ReleaseRoomRequest\) ProtoMessage\(\)](<#ReleaseRoomRequest.ProtoMessage>)
  - [func \(x \*ReleaseRoomRequest\) ProtoReflect\(\) protoreflect.Message](<#ReleaseRoomRequest.ProtoReflect>)
  - [func \(x \*ReleaseRoomRequest\) Reset\(\)](<#ReleaseRoomRequest.Reset>)
  - [func \(x \*ReleaseRoomRequest\) String\(\) string](<#ReleaseRoomRequest.String>)
- [type RoomManagerServiceClient](<#RoomManagerServiceClient>)
  - [func NewRoomManagerServiceClient\(cc grpc.ClientConnInterface\) RoomManagerServiceClient](<#NewRoomManagerServiceClient>)
- [type RoomManagerServiceServer](<#RoomManagerServiceServer>)
- [type UnimplementedRoomManagerServiceServer](<#UnimplementedRoomManagerServiceServer>)
  - [func \(UnimplementedRoomManagerServiceServer\) GetRoomWorker\(context.Context, \*GetRoomWorkerRequest\) \(\*GetRoomWorkerResponse, error\)](<#UnimplementedRoomManagerServiceServer.GetRoomWorker>)
  - [func \(UnimplementedRoomManagerServiceServer\) Ping\(context.Context, \*emptypb.Empty\) \(\*PingResponse, error\)](<#UnimplementedRoomManagerServiceServer.Ping>)
  - [func \(UnimplementedRoomManagerServiceServer\) PingWorker\(context.Context, \*PingWorkerRequest\) \(\*emptypb.Empty, error\)](<#UnimplementedRoomManagerServiceServer.PingWorker>)
  - [func \(UnimplementedRoomManagerServiceServer\) ProlongRoom\(context.Context, \*ProlongRoomRequest\) \(\*emptypb.Empty, error\)](<#UnimplementedRoomManagerServiceServer.ProlongRoom>)
  - [func \(UnimplementedRoomManagerServiceServer\) RegisterRoom\(context.Context, \*RegisterRoomRequest\) \(\*RegisterRoomResponse, error\)](<#UnimplementedRoomManagerServiceServer.RegisterRoom>)
  - [func \(UnimplementedRoomManagerServiceServer\) ReleaseRoom\(context.Context, \*ReleaseRoomRequest\) \(\*emptypb.Empty, error\)](<#UnimplementedRoomManagerServiceServer.ReleaseRoom>)
- [type UnsafeRoomManagerServiceServer](<#UnsafeRoomManagerServiceServer>)


## Constants

<a name="RoomManagerService_Ping_FullMethodName"></a>

```go
const (
    RoomManagerService_Ping_FullMethodName          = "/room_manager.RoomManagerService/Ping"
    RoomManagerService_GetRoomWorker_FullMethodName = "/room_manager.RoomManagerService/GetRoomWorker"
    RoomManagerService_PingWorker_FullMethodName    = "/room_manager.RoomManagerService/PingWorker"
    RoomManagerService_RegisterRoom_FullMethodName  = "/room_manager.RoomManagerService/RegisterRoom"
    RoomManagerService_ProlongRoom_FullMethodName   = "/room_manager.RoomManagerService/ProlongRoom"
    RoomManagerService_ReleaseRoom_FullMethodName   = "/room_manager.RoomManagerService/ReleaseRoom"
)
```

## Variables

<a name="File_room_manager_room_manager_proto"></a>

```go
var File_room_manager_room_manager_proto protoreflect.FileDescriptor
```

<a name="RoomManagerService_ServiceDesc"></a>RoomManagerService\_ServiceDesc is the grpc.ServiceDesc for RoomManagerService service. It's only intended for direct use with grpc.RegisterService, and not to be introspected or modified \(even as a copy\)

```go
var RoomManagerService_ServiceDesc = grpc.ServiceDesc{
    ServiceName: "room_manager.RoomManagerService",
    HandlerType: (*RoomManagerServiceServer)(nil),
    Methods: []grpc.MethodDesc{
        {
            MethodName: "Ping",
            Handler:    _RoomManagerService_Ping_Handler,
        },
        {
            MethodName: "GetRoomWorker",
            Handler:    _RoomManagerService_GetRoomWorker_Handler,
        },
        {
            MethodName: "PingWorker",
            Handler:    _RoomManagerService_PingWorker_Handler,
        },
        {
            MethodName: "RegisterRoom",
            Handler:    _RoomManagerService_RegisterRoom_Handler,
        },
        {
            MethodName: "ProlongRoom",
            Handler:    _RoomManagerService_ProlongRoom_Handler,
        },
        {
            MethodName: "ReleaseRoom",
            Handler:    _RoomManagerService_ReleaseRoom_Handler,
        },
    },
    Streams:  []grpc.StreamDesc{},
    Metadata: "room_manager/room_manager.proto",
}
```

<a name="RegisterRoomManagerServiceServer"></a>
## func RegisterRoomManagerServiceServer

```go
func RegisterRoomManagerServiceServer(s grpc.ServiceRegistrar, srv RoomManagerServiceServer)
```



<a name="GetRoomWorkerRequest"></a>
## type GetRoomWorkerRequest



```go
type GetRoomWorkerRequest struct {
    RoomId string `protobuf:"bytes,1,opt,name=room_id,json=roomId,proto3" json:"room_id,omitempty"`
    // contains filtered or unexported fields
}
```

<a name="GetRoomWorkerRequest.Descriptor"></a>
### func \(\*GetRoomWorkerRequest\) Descriptor

```go
func (*GetRoomWorkerRequest) Descriptor() ([]byte, []int)
```

Deprecated: Use GetRoomWorkerRequest.ProtoReflect.Descriptor instead.

<a name="GetRoomWorkerRequest.GetRoomId"></a>
### func \(\*GetRoomWorkerRequest\) GetRoomId

```go
func (x *GetRoomWorkerRequest) GetRoomId() string
```



<a name="GetRoomWorkerRequest.ProtoMessage"></a>
### func \(\*GetRoomWorkerRequest\) ProtoMessage

```go
func (*GetRoomWorkerRequest) ProtoMessage()
```



<a name="GetRoomWorkerRequest.ProtoReflect"></a>
### func \(\*GetRoomWorkerRequest\) ProtoReflect

```go
func (x *GetRoomWorkerRequest) ProtoReflect() protoreflect.Message
```



<a name="GetRoomWorkerRequest.Reset"></a>
### func \(\*GetRoomWorkerRequest\) Reset

```go
func (x *GetRoomWorkerRequest) Reset()
```



<a name="GetRoomWorkerRequest.String"></a>
### func \(\*GetRoomWorkerRequest\) String

```go
func (x *GetRoomWorkerRequest) String() string
```



<a name="GetRoomWorkerResponse"></a>
## type GetRoomWorkerResponse



```go
type GetRoomWorkerResponse struct {
    Worker string `protobuf:"bytes,1,opt,name=worker,proto3" json:"worker,omitempty"`
    // contains filtered or unexported fields
}
```

<a name="GetRoomWorkerResponse.Descriptor"></a>
### func \(\*GetRoomWorkerResponse\) Descriptor

```go
func (*GetRoomWorkerResponse) Descriptor() ([]byte, []int)
```

Deprecated: Use GetRoomWorkerResponse.ProtoReflect.Descriptor instead.

<a name="GetRoomWorkerResponse.GetWorker"></a>
### func \(\*GetRoomWorkerResponse\) GetWorker

```go
func (x *GetRoomWorkerResponse) GetWorker() string
```



<a name="GetRoomWorkerResponse.ProtoMessage"></a>
### func \(\*GetRoomWorkerResponse\) ProtoMessage

```go
func (*GetRoomWorkerResponse) ProtoMessage()
```



<a name="GetRoomWorkerResponse.ProtoReflect"></a>
### func \(\*GetRoomWorkerResponse\) ProtoReflect

```go
func (x *GetRoomWorkerResponse) ProtoReflect() protoreflect.Message
```



<a name="GetRoomWorkerResponse.Reset"></a>
### func \(\*GetRoomWorkerResponse\) Reset

```go
func (x *GetRoomWorkerResponse) Reset()
```



<a name="GetRoomWorkerResponse.String"></a>
### func \(\*GetRoomWorkerResponse\) String

```go
func (x *GetRoomWorkerResponse) String() string
```



<a name="PingResponse"></a>
## type PingResponse



```go
type PingResponse struct {
    Ok bool `protobuf:"varint,1,opt,name=ok,proto3" json:"ok,omitempty"`
    // contains filtered or unexported fields
}
```

<a name="PingResponse.Descriptor"></a>
### func \(\*PingResponse\) Descriptor

```go
func (*PingResponse) Descriptor() ([]byte, []int)
```

Deprecated: Use PingResponse.ProtoReflect.Descriptor instead.

<a name="PingResponse.GetOk"></a>
### func \(\*PingResponse\) GetOk

```go
func (x *PingResponse) GetOk() bool
```



<a name="PingResponse.ProtoMessage"></a>
### func \(\*PingResponse\) ProtoMessage

```go
func (*PingResponse) ProtoMessage()
```



<a name="PingResponse.ProtoReflect"></a>
### func \(\*PingResponse\) ProtoReflect

```go
func (x *PingResponse) ProtoReflect() protoreflect.Message
```



<a name="PingResponse.Reset"></a>
### func \(\*PingResponse\) Reset

```go
func (x *PingResponse) Reset()
```



<a name="PingResponse.String"></a>
### func \(\*PingResponse\) String

```go
func (x *PingResponse) String() string
```



<a name="PingWorkerRequest"></a>
## type PingWorkerRequest



```go
type PingWorkerRequest struct {
    Worker      string   `protobuf:"bytes,1,opt,name=worker,proto3" json:"worker,omitempty"`
    LoadedRooms []string `protobuf:"bytes,2,rep,name=loaded_rooms,json=loadedRooms,proto3" json:"loaded_rooms,omitempty"`
    // contains filtered or unexported fields
}
```

<a name="PingWorkerRequest.Descriptor"></a>
### func \(\*PingWorkerRequest\) Descriptor

```go
func (*PingWorkerRequest) Descriptor() ([]byte, []int)
```

Deprecated: Use PingWorkerRequest.ProtoReflect.Descriptor instead.

<a name="PingWorkerRequest.GetLoadedRooms"></a>
### func \(\*PingWorkerRequest\) GetLoadedRooms

```go
func (x *PingWorkerRequest) GetLoadedRooms() []string
```



<a name="PingWorkerRequest.GetWorker"></a>
### func \(\*PingWorkerRequest\) GetWorker

```go
func (x *PingWorkerRequest) GetWorker() string
```



<a name="PingWorkerRequest.ProtoMessage"></a>
### func \(\*PingWorkerRequest\) ProtoMessage

```go
func (*PingWorkerRequest) ProtoMessage()
```



<a name="PingWorkerRequest.ProtoReflect"></a>
### func \(\*PingWorkerRequest\) ProtoReflect

```go
func (x *PingWorkerRequest) ProtoReflect() protoreflect.Message
```



<a name="PingWorkerRequest.Reset"></a>
### func \(\*PingWorkerRequest\) Reset

```go
func (x *PingWorkerRequest) Reset()
```



<a name="PingWorkerRequest.String"></a>
### func \(\*PingWorkerRequest\) String

```go
func (x *PingWorkerRequest) String() string
```



<a name="ProlongRoomRequest"></a>
## type ProlongRoomRequest



```go
type ProlongRoomRequest struct {
    RoomId string `protobuf:"bytes,1,opt,name=room_id,json=roomId,proto3" json:"room_id,omitempty"`
    Worker string `protobuf:"bytes,2,opt,name=worker,proto3" json:"worker,omitempty"`
    // contains filtered or unexported fields
}
```

<a name="ProlongRoomRequest.Descriptor"></a>
### func \(\*ProlongRoomRequest\) Descriptor

```go
func (*ProlongRoomRequest) Descriptor() ([]byte, []int)
```

Deprecated: Use ProlongRoomRequest.ProtoReflect.Descriptor instead.

<a name="ProlongRoomRequest.GetRoomId"></a>
### func \(\*ProlongRoomRequest\) GetRoomId

```go
func (x *ProlongRoomRequest) GetRoomId() string
```



<a name="ProlongRoomRequest.GetWorker"></a>
### func \(\*ProlongRoomRequest\) GetWorker

```go
func (x *ProlongRoomRequest) GetWorker() string
```



<a name="ProlongRoomRequest.ProtoMessage"></a>
### func \(\*ProlongRoomRequest\) ProtoMessage

```go
func (*ProlongRoomRequest) ProtoMessage()
```



<a name="ProlongRoomRequest.ProtoReflect"></a>
### func \(\*ProlongRoomRequest\) ProtoReflect

```go
func (x *ProlongRoomRequest) ProtoReflect() protoreflect.Message
```



<a name="ProlongRoomRequest.Reset"></a>
### func \(\*ProlongRoomRequest\) Reset

```go
func (x *ProlongRoomRequest) Reset()
```



<a name="ProlongRoomRequest.String"></a>
### func \(\*ProlongRoomRequest\) String

```go
func (x *ProlongRoomRequest) String() string
```



<a name="RegisterRoomRequest"></a>
## type RegisterRoomRequest



```go
type RegisterRoomRequest struct {
    RoomId string `protobuf:"bytes,1,opt,name=room_id,json=roomId,proto3" json:"room_id,omitempty"`
    Worker string `protobuf:"bytes,2,opt,name=worker,proto3" json:"worker,omitempty"`
    // contains filtered or unexported fields
}
```

<a name="RegisterRoomRequest.Descriptor"></a>
### func \(\*RegisterRoomRequest\) Descriptor

```go
func (*RegisterRoomRequest) Descriptor() ([]byte, []int)
```

Deprecated: Use RegisterRoomRequest.ProtoReflect.Descriptor instead.

<a name="RegisterRoomRequest.GetRoomId"></a>
### func \(\*RegisterRoomRequest\) GetRoomId

```go
func (x *RegisterRoomRequest) GetRoomId() string
```



<a name="RegisterRoomRequest.GetWorker"></a>
### func \(\*RegisterRoomRequest\) GetWorker

```go
func (x *RegisterRoomRequest) GetWorker() string
```



<a name="RegisterRoomRequest.ProtoMessage"></a>
### func \(\*RegisterRoomRequest\) ProtoMessage

```go
func (*RegisterRoomRequest) ProtoMessage()
```



<a name="RegisterRoomRequest.ProtoReflect"></a>
### func \(\*RegisterRoomRequest\) ProtoReflect

```go
func (x *RegisterRoomRequest) ProtoReflect() protoreflect.Message
```



<a name="RegisterRoomRequest.Reset"></a>
### func \(\*RegisterRoomRequest\) Reset

```go
func (x *RegisterRoomRequest) Reset()
```



<a name="RegisterRoomRequest.String"></a>
### func \(\*RegisterRoomRequest\) String

```go
func (x *RegisterRoomRequest) String() string
```



<a name="RegisterRoomResponse"></a>
## type RegisterRoomResponse



```go
type RegisterRoomResponse struct {
    Worker string `protobuf:"bytes,1,opt,name=worker,proto3" json:"worker,omitempty"`
    // contains filtered or unexported fields
}
```

<a name="RegisterRoomResponse.Descriptor"></a>
### func \(\*RegisterRoomResponse\) Descriptor

```go
func (*RegisterRoomResponse) Descriptor() ([]byte, []int)
```

Deprecated: Use RegisterRoomResponse.ProtoReflect.Descriptor instead.

<a name="RegisterRoomResponse.GetWorker"></a>
### func \(\*RegisterRoomResponse\) GetWorker

```go
func (x *RegisterRoomResponse) GetWorker() string
```



<a name="RegisterRoomResponse.ProtoMessage"></a>
### func \(\*RegisterRoomResponse\) ProtoMessage

```go
func (*RegisterRoomResponse) ProtoMessage()
```



<a name="RegisterRoomResponse.ProtoReflect"></a>
### func \(\*RegisterRoomResponse\) ProtoReflect

```go
func (x *RegisterRoomResponse) ProtoReflect() protoreflect.Message
```



<a name="RegisterRoomResponse.Reset"></a>
### func \(\*RegisterRoomResponse\) Reset

```go
func (x *RegisterRoomResponse) Reset()
```



<a name="RegisterRoomResponse.String"></a>
### func \(\*RegisterRoomResponse\) String

```go
func (x *RegisterRoomResponse) String() string
```



<a name="ReleaseRoomRequest"></a>
## type ReleaseRoomRequest



```go
type ReleaseRoomRequest struct {
    RoomId string `protobuf:"bytes,1,opt,name=room_id,json=roomId,proto3" json:"room_id,omitempty"`
    Worker string `protobuf:"bytes,2,opt,name=worker,proto3" json:"worker,omitempty"`
    // contains filtered or unexported fields
}
```

<a name="ReleaseRoomRequest.Descriptor"></a>
### func \(\*ReleaseRoomRequest\) Descriptor

```go
func (*ReleaseRoomRequest) Descriptor() ([]byte, []int)
```

Deprecated: Use ReleaseRoomRequest.ProtoReflect.Descriptor instead.

<a name="ReleaseRoomRequest.GetRoomId"></a>
### func \(\*ReleaseRoomRequest\) GetRoomId

```go
func (x *ReleaseRoomRequest) GetRoomId() string
```



<a name="ReleaseRoomRequest.GetWorker"></a>
### func \(\*ReleaseRoomRequest\) GetWorker

```go
func (x *ReleaseRoomRequest) GetWorker() string
```



<a name="ReleaseRoomRequest.ProtoMessage"></a>
### func \(\*ReleaseRoomRequest\) ProtoMessage

```go
func (*ReleaseRoomRequest) ProtoMessage()
```



<a name="ReleaseRoomRequest.ProtoReflect"></a>
### func \(\*ReleaseRoomRequest\) ProtoReflect

```go
func (x *ReleaseRoomRequest) ProtoReflect() protoreflect.Message
```



<a name="ReleaseRoomRequest.Reset"></a>
### func \(\*ReleaseRoomRequest\) Reset

```go
func (x *ReleaseRoomRequest) Reset()
```



<a name="ReleaseRoomRequest.String"></a>
### func \(\*ReleaseRoomRequest\) String

```go
func (x *ReleaseRoomRequest) String() string
```



<a name="RoomManagerServiceClient"></a>
## type RoomManagerServiceClient

RoomManagerServiceClient is the client API for RoomManagerService service.

For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.

```go
type RoomManagerServiceClient interface {
    Ping(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*PingResponse, error)
    // for main service.
    GetRoomWorker(ctx context.Context, in *GetRoomWorkerRequest, opts ...grpc.CallOption) (*GetRoomWorkerResponse, error)
    // for workers.
    PingWorker(ctx context.Context, in *PingWorkerRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
    RegisterRoom(ctx context.Context, in *RegisterRoomRequest, opts ...grpc.CallOption) (*RegisterRoomResponse, error)
    ProlongRoom(ctx context.Context, in *ProlongRoomRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
    ReleaseRoom(ctx context.Context, in *ReleaseRoomRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}
```

<a name="NewRoomManagerServiceClient"></a>
### func NewRoomManagerServiceClient

```go
func NewRoomManagerServiceClient(cc grpc.ClientConnInterface) RoomManagerServiceClient
```



<a name="RoomManagerServiceServer"></a>
## type RoomManagerServiceServer

RoomManagerServiceServer is the server API for RoomManagerService service. All implementations must embed UnimplementedRoomManagerServiceServer for forward compatibility.

```go
type RoomManagerServiceServer interface {
    Ping(context.Context, *emptypb.Empty) (*PingResponse, error)
    // for main service.
    GetRoomWorker(context.Context, *GetRoomWorkerRequest) (*GetRoomWorkerResponse, error)
    // for workers.
    PingWorker(context.Context, *PingWorkerRequest) (*emptypb.Empty, error)
    RegisterRoom(context.Context, *RegisterRoomRequest) (*RegisterRoomResponse, error)
    ProlongRoom(context.Context, *ProlongRoomRequest) (*emptypb.Empty, error)
    ReleaseRoom(context.Context, *ReleaseRoomRequest) (*emptypb.Empty, error)
    // contains filtered or unexported methods
}
```

<a name="UnimplementedRoomManagerServiceServer"></a>
## type UnimplementedRoomManagerServiceServer

UnimplementedRoomManagerServiceServer must be embedded to have forward compatible implementations.

NOTE: this should be embedded by value instead of pointer to avoid a nil pointer dereference when methods are called.

```go
type UnimplementedRoomManagerServiceServer struct{}
```

<a name="UnimplementedRoomManagerServiceServer.GetRoomWorker"></a>
### func \(UnimplementedRoomManagerServiceServer\) GetRoomWorker

```go
func (UnimplementedRoomManagerServiceServer) GetRoomWorker(context.Context, *GetRoomWorkerRequest) (*GetRoomWorkerResponse, error)
```



<a name="UnimplementedRoomManagerServiceServer.Ping"></a>
### func \(UnimplementedRoomManagerServiceServer\) Ping

```go
func (UnimplementedRoomManagerServiceServer) Ping(context.Context, *emptypb.Empty) (*PingResponse, error)
```



<a name="UnimplementedRoomManagerServiceServer.PingWorker"></a>
### func \(UnimplementedRoomManagerServiceServer\) PingWorker

```go
func (UnimplementedRoomManagerServiceServer) PingWorker(context.Context, *PingWorkerRequest) (*emptypb.Empty, error)
```



<a name="UnimplementedRoomManagerServiceServer.ProlongRoom"></a>
### func \(UnimplementedRoomManagerServiceServer\) ProlongRoom

```go
func (UnimplementedRoomManagerServiceServer) ProlongRoom(context.Context, *ProlongRoomRequest) (*emptypb.Empty, error)
```



<a name="UnimplementedRoomManagerServiceServer.RegisterRoom"></a>
### func \(UnimplementedRoomManagerServiceServer\) RegisterRoom

```go
func (UnimplementedRoomManagerServiceServer) RegisterRoom(context.Context, *RegisterRoomRequest) (*RegisterRoomResponse, error)
```



<a name="UnimplementedRoomManagerServiceServer.ReleaseRoom"></a>
### func \(UnimplementedRoomManagerServiceServer\) ReleaseRoom

```go
func (UnimplementedRoomManagerServiceServer) ReleaseRoom(context.Context, *ReleaseRoomRequest) (*emptypb.Empty, error)
```



<a name="UnsafeRoomManagerServiceServer"></a>
## type UnsafeRoomManagerServiceServer

UnsafeRoomManagerServiceServer may be embedded to opt out of forward compatibility for this service. Use of this interface is not recommended, as added methods to RoomManagerServiceServer will result in compilation errors.

```go
type UnsafeRoomManagerServiceServer interface {
    // contains filtered or unexported methods
}
```

# proto

```go
import "github.com/xolra0d/alias-online/shared/proto/vocab_manager"
```

## Index

- [Constants](<#constants>)
- [Variables](<#variables>)
- [func RegisterVocabManagerServiceServer\(s grpc.ServiceRegistrar, srv VocabManagerServiceServer\)](<#RegisterVocabManagerServiceServer>)
- [type GetAvailableVocabsResponse](<#GetAvailableVocabsResponse>)
  - [func \(\*GetAvailableVocabsResponse\) Descriptor\(\) \(\[\]byte, \[\]int\)](<#GetAvailableVocabsResponse.Descriptor>)
  - [func \(x \*GetAvailableVocabsResponse\) GetNames\(\) \[\]string](<#GetAvailableVocabsResponse.GetNames>)
  - [func \(\*GetAvailableVocabsResponse\) ProtoMessage\(\)](<#GetAvailableVocabsResponse.ProtoMessage>)
  - [func \(x \*GetAvailableVocabsResponse\) ProtoReflect\(\) protoreflect.Message](<#GetAvailableVocabsResponse.ProtoReflect>)
  - [func \(x \*GetAvailableVocabsResponse\) Reset\(\)](<#GetAvailableVocabsResponse.Reset>)
  - [func \(x \*GetAvailableVocabsResponse\) String\(\) string](<#GetAvailableVocabsResponse.String>)
- [type GetVocabRequest](<#GetVocabRequest>)
  - [func \(\*GetVocabRequest\) Descriptor\(\) \(\[\]byte, \[\]int\)](<#GetVocabRequest.Descriptor>)
  - [func \(x \*GetVocabRequest\) GetName\(\) string](<#GetVocabRequest.GetName>)
  - [func \(\*GetVocabRequest\) ProtoMessage\(\)](<#GetVocabRequest.ProtoMessage>)
  - [func \(x \*GetVocabRequest\) ProtoReflect\(\) protoreflect.Message](<#GetVocabRequest.ProtoReflect>)
  - [func \(x \*GetVocabRequest\) Reset\(\)](<#GetVocabRequest.Reset>)
  - [func \(x \*GetVocabRequest\) String\(\) string](<#GetVocabRequest.String>)
- [type GetVocabResponse](<#GetVocabResponse>)
  - [func \(\*GetVocabResponse\) Descriptor\(\) \(\[\]byte, \[\]int\)](<#GetVocabResponse.Descriptor>)
  - [func \(x \*GetVocabResponse\) GetPrimaryWords\(\) \[\]string](<#GetVocabResponse.GetPrimaryWords>)
  - [func \(x \*GetVocabResponse\) GetRudeWords\(\) \[\]string](<#GetVocabResponse.GetRudeWords>)
  - [func \(\*GetVocabResponse\) ProtoMessage\(\)](<#GetVocabResponse.ProtoMessage>)
  - [func \(x \*GetVocabResponse\) ProtoReflect\(\) protoreflect.Message](<#GetVocabResponse.ProtoReflect>)
  - [func \(x \*GetVocabResponse\) Reset\(\)](<#GetVocabResponse.Reset>)
  - [func \(x \*GetVocabResponse\) String\(\) string](<#GetVocabResponse.String>)
- [type PingResponse](<#PingResponse>)
  - [func \(\*PingResponse\) Descriptor\(\) \(\[\]byte, \[\]int\)](<#PingResponse.Descriptor>)
  - [func \(x \*PingResponse\) GetOk\(\) bool](<#PingResponse.GetOk>)
  - [func \(\*PingResponse\) ProtoMessage\(\)](<#PingResponse.ProtoMessage>)
  - [func \(x \*PingResponse\) ProtoReflect\(\) protoreflect.Message](<#PingResponse.ProtoReflect>)
  - [func \(x \*PingResponse\) Reset\(\)](<#PingResponse.Reset>)
  - [func \(x \*PingResponse\) String\(\) string](<#PingResponse.String>)
- [type UnimplementedVocabManagerServiceServer](<#UnimplementedVocabManagerServiceServer>)
  - [func \(UnimplementedVocabManagerServiceServer\) GetAvailableVocabs\(context.Context, \*emptypb.Empty\) \(\*GetAvailableVocabsResponse, error\)](<#UnimplementedVocabManagerServiceServer.GetAvailableVocabs>)
  - [func \(UnimplementedVocabManagerServiceServer\) GetVocab\(context.Context, \*GetVocabRequest\) \(\*GetVocabResponse, error\)](<#UnimplementedVocabManagerServiceServer.GetVocab>)
  - [func \(UnimplementedVocabManagerServiceServer\) Ping\(context.Context, \*emptypb.Empty\) \(\*PingResponse, error\)](<#UnimplementedVocabManagerServiceServer.Ping>)
- [type UnsafeVocabManagerServiceServer](<#UnsafeVocabManagerServiceServer>)
- [type VocabManagerServiceClient](<#VocabManagerServiceClient>)
  - [func NewVocabManagerServiceClient\(cc grpc.ClientConnInterface\) VocabManagerServiceClient](<#NewVocabManagerServiceClient>)
- [type VocabManagerServiceServer](<#VocabManagerServiceServer>)


## Constants

<a name="VocabManagerService_Ping_FullMethodName"></a>

```go
const (
    VocabManagerService_Ping_FullMethodName               = "/vocab_manager.VocabManagerService/Ping"
    VocabManagerService_GetAvailableVocabs_FullMethodName = "/vocab_manager.VocabManagerService/GetAvailableVocabs"
    VocabManagerService_GetVocab_FullMethodName           = "/vocab_manager.VocabManagerService/GetVocab"
)
```

## Variables

<a name="File_vocab_manager_vocab_manager_proto"></a>

```go
var File_vocab_manager_vocab_manager_proto protoreflect.FileDescriptor
```

<a name="VocabManagerService_ServiceDesc"></a>VocabManagerService\_ServiceDesc is the grpc.ServiceDesc for VocabManagerService service. It's only intended for direct use with grpc.RegisterService, and not to be introspected or modified \(even as a copy\)

```go
var VocabManagerService_ServiceDesc = grpc.ServiceDesc{
    ServiceName: "vocab_manager.VocabManagerService",
    HandlerType: (*VocabManagerServiceServer)(nil),
    Methods: []grpc.MethodDesc{
        {
            MethodName: "Ping",
            Handler:    _VocabManagerService_Ping_Handler,
        },
        {
            MethodName: "GetAvailableVocabs",
            Handler:    _VocabManagerService_GetAvailableVocabs_Handler,
        },
        {
            MethodName: "GetVocab",
            Handler:    _VocabManagerService_GetVocab_Handler,
        },
    },
    Streams:  []grpc.StreamDesc{},
    Metadata: "vocab_manager/vocab_manager.proto",
}
```

<a name="RegisterVocabManagerServiceServer"></a>
## func RegisterVocabManagerServiceServer

```go
func RegisterVocabManagerServiceServer(s grpc.ServiceRegistrar, srv VocabManagerServiceServer)
```



<a name="GetAvailableVocabsResponse"></a>
## type GetAvailableVocabsResponse



```go
type GetAvailableVocabsResponse struct {
    Names []string `protobuf:"bytes,1,rep,name=names,proto3" json:"names,omitempty"`
    // contains filtered or unexported fields
}
```

<a name="GetAvailableVocabsResponse.Descriptor"></a>
### func \(\*GetAvailableVocabsResponse\) Descriptor

```go
func (*GetAvailableVocabsResponse) Descriptor() ([]byte, []int)
```

Deprecated: Use GetAvailableVocabsResponse.ProtoReflect.Descriptor instead.

<a name="GetAvailableVocabsResponse.GetNames"></a>
### func \(\*GetAvailableVocabsResponse\) GetNames

```go
func (x *GetAvailableVocabsResponse) GetNames() []string
```



<a name="GetAvailableVocabsResponse.ProtoMessage"></a>
### func \(\*GetAvailableVocabsResponse\) ProtoMessage

```go
func (*GetAvailableVocabsResponse) ProtoMessage()
```



<a name="GetAvailableVocabsResponse.ProtoReflect"></a>
### func \(\*GetAvailableVocabsResponse\) ProtoReflect

```go
func (x *GetAvailableVocabsResponse) ProtoReflect() protoreflect.Message
```



<a name="GetAvailableVocabsResponse.Reset"></a>
### func \(\*GetAvailableVocabsResponse\) Reset

```go
func (x *GetAvailableVocabsResponse) Reset()
```



<a name="GetAvailableVocabsResponse.String"></a>
### func \(\*GetAvailableVocabsResponse\) String

```go
func (x *GetAvailableVocabsResponse) String() string
```



<a name="GetVocabRequest"></a>
## type GetVocabRequest



```go
type GetVocabRequest struct {
    Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
    // contains filtered or unexported fields
}
```

<a name="GetVocabRequest.Descriptor"></a>
### func \(\*GetVocabRequest\) Descriptor

```go
func (*GetVocabRequest) Descriptor() ([]byte, []int)
```

Deprecated: Use GetVocabRequest.ProtoReflect.Descriptor instead.

<a name="GetVocabRequest.GetName"></a>
### func \(\*GetVocabRequest\) GetName

```go
func (x *GetVocabRequest) GetName() string
```



<a name="GetVocabRequest.ProtoMessage"></a>
### func \(\*GetVocabRequest\) ProtoMessage

```go
func (*GetVocabRequest) ProtoMessage()
```



<a name="GetVocabRequest.ProtoReflect"></a>
### func \(\*GetVocabRequest\) ProtoReflect

```go
func (x *GetVocabRequest) ProtoReflect() protoreflect.Message
```



<a name="GetVocabRequest.Reset"></a>
### func \(\*GetVocabRequest\) Reset

```go
func (x *GetVocabRequest) Reset()
```



<a name="GetVocabRequest.String"></a>
### func \(\*GetVocabRequest\) String

```go
func (x *GetVocabRequest) String() string
```



<a name="GetVocabResponse"></a>
## type GetVocabResponse



```go
type GetVocabResponse struct {
    PrimaryWords []string `protobuf:"bytes,1,rep,name=primary_words,json=primaryWords,proto3" json:"primary_words,omitempty"`
    RudeWords    []string `protobuf:"bytes,2,rep,name=rude_words,json=rudeWords,proto3" json:"rude_words,omitempty"`
    // contains filtered or unexported fields
}
```

<a name="GetVocabResponse.Descriptor"></a>
### func \(\*GetVocabResponse\) Descriptor

```go
func (*GetVocabResponse) Descriptor() ([]byte, []int)
```

Deprecated: Use GetVocabResponse.ProtoReflect.Descriptor instead.

<a name="GetVocabResponse.GetPrimaryWords"></a>
### func \(\*GetVocabResponse\) GetPrimaryWords

```go
func (x *GetVocabResponse) GetPrimaryWords() []string
```



<a name="GetVocabResponse.GetRudeWords"></a>
### func \(\*GetVocabResponse\) GetRudeWords

```go
func (x *GetVocabResponse) GetRudeWords() []string
```



<a name="GetVocabResponse.ProtoMessage"></a>
### func \(\*GetVocabResponse\) ProtoMessage

```go
func (*GetVocabResponse) ProtoMessage()
```



<a name="GetVocabResponse.ProtoReflect"></a>
### func \(\*GetVocabResponse\) ProtoReflect

```go
func (x *GetVocabResponse) ProtoReflect() protoreflect.Message
```



<a name="GetVocabResponse.Reset"></a>
### func \(\*GetVocabResponse\) Reset

```go
func (x *GetVocabResponse) Reset()
```



<a name="GetVocabResponse.String"></a>
### func \(\*GetVocabResponse\) String

```go
func (x *GetVocabResponse) String() string
```



<a name="PingResponse"></a>
## type PingResponse



```go
type PingResponse struct {
    Ok bool `protobuf:"varint,1,opt,name=ok,proto3" json:"ok,omitempty"`
    // contains filtered or unexported fields
}
```

<a name="PingResponse.Descriptor"></a>
### func \(\*PingResponse\) Descriptor

```go
func (*PingResponse) Descriptor() ([]byte, []int)
```

Deprecated: Use PingResponse.ProtoReflect.Descriptor instead.

<a name="PingResponse.GetOk"></a>
### func \(\*PingResponse\) GetOk

```go
func (x *PingResponse) GetOk() bool
```



<a name="PingResponse.ProtoMessage"></a>
### func \(\*PingResponse\) ProtoMessage

```go
func (*PingResponse) ProtoMessage()
```



<a name="PingResponse.ProtoReflect"></a>
### func \(\*PingResponse\) ProtoReflect

```go
func (x *PingResponse) ProtoReflect() protoreflect.Message
```



<a name="PingResponse.Reset"></a>
### func \(\*PingResponse\) Reset

```go
func (x *PingResponse) Reset()
```



<a name="PingResponse.String"></a>
### func \(\*PingResponse\) String

```go
func (x *PingResponse) String() string
```



<a name="UnimplementedVocabManagerServiceServer"></a>
## type UnimplementedVocabManagerServiceServer

UnimplementedVocabManagerServiceServer must be embedded to have forward compatible implementations.

NOTE: this should be embedded by value instead of pointer to avoid a nil pointer dereference when methods are called.

```go
type UnimplementedVocabManagerServiceServer struct{}
```

<a name="UnimplementedVocabManagerServiceServer.GetAvailableVocabs"></a>
### func \(UnimplementedVocabManagerServiceServer\) GetAvailableVocabs

```go
func (UnimplementedVocabManagerServiceServer) GetAvailableVocabs(context.Context, *emptypb.Empty) (*GetAvailableVocabsResponse, error)
```



<a name="UnimplementedVocabManagerServiceServer.GetVocab"></a>
### func \(UnimplementedVocabManagerServiceServer\) GetVocab

```go
func (UnimplementedVocabManagerServiceServer) GetVocab(context.Context, *GetVocabRequest) (*GetVocabResponse, error)
```



<a name="UnimplementedVocabManagerServiceServer.Ping"></a>
### func \(UnimplementedVocabManagerServiceServer\) Ping

```go
func (UnimplementedVocabManagerServiceServer) Ping(context.Context, *emptypb.Empty) (*PingResponse, error)
```



<a name="UnsafeVocabManagerServiceServer"></a>
## type UnsafeVocabManagerServiceServer

UnsafeVocabManagerServiceServer may be embedded to opt out of forward compatibility for this service. Use of this interface is not recommended, as added methods to VocabManagerServiceServer will result in compilation errors.

```go
type UnsafeVocabManagerServiceServer interface {
    // contains filtered or unexported methods
}
```

<a name="VocabManagerServiceClient"></a>
## type VocabManagerServiceClient

VocabManagerServiceClient is the client API for VocabManagerService service.

For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.

```go
type VocabManagerServiceClient interface {
    Ping(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*PingResponse, error)
    GetAvailableVocabs(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetAvailableVocabsResponse, error)
    GetVocab(ctx context.Context, in *GetVocabRequest, opts ...grpc.CallOption) (*GetVocabResponse, error)
}
```

<a name="NewVocabManagerServiceClient"></a>
### func NewVocabManagerServiceClient

```go
func NewVocabManagerServiceClient(cc grpc.ClientConnInterface) VocabManagerServiceClient
```



<a name="VocabManagerServiceServer"></a>
## type VocabManagerServiceServer

VocabManagerServiceServer is the server API for VocabManagerService service. All implementations must embed UnimplementedVocabManagerServiceServer for forward compatibility.

```go
type VocabManagerServiceServer interface {
    Ping(context.Context, *emptypb.Empty) (*PingResponse, error)
    GetAvailableVocabs(context.Context, *emptypb.Empty) (*GetAvailableVocabsResponse, error)
    GetVocab(context.Context, *GetVocabRequest) (*GetVocabResponse, error)
    // contains filtered or unexported methods
}
```

Generated by [gomarkdoc](<https://github.com/princjef/gomarkdoc>)
