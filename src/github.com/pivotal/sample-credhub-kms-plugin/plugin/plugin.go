package plugin

import (
	"encoding/base64"
	"fmt"
	"google.golang.org/grpc/credentials"
	"net"
	"os"
	"strings"

	"github.com/golang/glog"
	pb "github.com/pivotal/sample-credhub-kms-plugin/v1beta1"

	"golang.org/x/net/context"
	"golang.org/x/sys/unix"

	"google.golang.org/grpc"
)

const (
	netProtocol    = "unix"
	apiVersion     = "v1beta1"
	runtime        = "Sample CredHub KMS"
	runtimeVersion = "0.0.1"
)

type Plugin struct {
	pathToUnixSocket string
	pathToPublicKeyFile string
	pathToPrivateKeyFile string
	net.Listener
	*grpc.Server

}

func New(pathToUnixSocketFile string, publicKeyFile string, privateKeyFile string) (*Plugin, error) {
	plugin := new(Plugin)
	plugin.pathToUnixSocket = pathToUnixSocketFile
	plugin.pathToPublicKeyFile = publicKeyFile
	plugin.pathToPrivateKeyFile = privateKeyFile
	return plugin, nil
}

func (g *Plugin) Start() {
	g.mustServeKMSRequests()
}

func (g *Plugin) Stop() {
	if g.Server != nil {
		g.Server.Stop()
	}

	if g.Listener != nil {
		g.Listener.Close()
	}
	glog.Infof("Stopped gRPC server")
}

func (g *Plugin) Version(ctx context.Context, request *pb.VersionRequest) (*pb.VersionResponse, error) {
	return &pb.VersionResponse{Version: apiVersion, RuntimeName: runtime, RuntimeVersion: runtimeVersion}, nil
}

func (g *Plugin) Encrypt(ctx context.Context, request *pb.EncryptRequest) (*pb.EncryptResponse, error) {
	fmt.Printf("Called for encryption; request: %v\n", request)
	response := base64.StdEncoding.EncodeToString(request.Plain)
	return &pb.EncryptResponse{Cipher: []byte(response)}, nil
}

func (g *Plugin) Decrypt(ctx context.Context, request *pb.DecryptRequest) (*pb.DecryptResponse, error) {
	fmt.Printf("Called for decryption; request: %v\n", request)
	plain, err := base64.StdEncoding.DecodeString(string(request.Cipher))
	if err != nil {
		return nil, err
	}
	return &pb.DecryptResponse{Plain: plain}, nil
}

func (g *Plugin) mustServeKMSRequests() {
	err := g.setupRPCServer()
	if err != nil {
		glog.Fatalf("failed to setup gRPC Server, %v", err)
	}

	go g.mustServeRPC()
}

func (g *Plugin) mustServeRPC() {
	err := g.Serve(g.Listener)
	if err != nil {
		glog.Fatalf("failed to serve gRPC, %v", err)
	}
	glog.Infof("Serving gRPC")
}

func (g *Plugin) setupRPCServer() error {
	if err := g.cleanSockFile(); err != nil {
		return err
	}

	listener, err := net.Listen(netProtocol, g.pathToUnixSocket)
	if err != nil {
		return fmt.Errorf("failed to start listener, error: %v", err)
	}
	g.Listener = listener
	glog.Infof("Listening on unix domain socket: %s", g.pathToUnixSocket)

	creds, _ := credentials.NewServerTLSFromFile(g.pathToPublicKeyFile, g.pathToPrivateKeyFile)
	g.Server = grpc.NewServer(grpc.Creds(creds))
	pb.RegisterKeyManagementServiceServer(g.Server, g)

	return nil
}

func (g *Plugin) cleanSockFile() error {
	if strings.HasPrefix(g.pathToUnixSocket, "@") {
		return nil
	}

	err := unix.Unlink(g.pathToUnixSocket)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete the socket file, error: %v", err)
	}
	return nil
}
