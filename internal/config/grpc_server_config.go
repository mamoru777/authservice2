package config

type GrpcServerConfig struct {
	HTTPAddr string `env:"GATEWAY_HTTP_ADDR" envDefault:":13998"`
	GRPCAddr string `env:"GATEWAY_GRPC_ADDR" envDefault:":13999"`
}
