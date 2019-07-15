package smoke_tests

import (
	"context"
	"crypto/x509"
	"os"
	pb "smoke-tests/v1beta1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var _ = Describe("SmokeTests", func() {
	It("can encrypt and decrypt using sample plugin", func() {
		certString, ok := os.LookupEnv("KMS_PLUGIN_CA_CERT")
		if !ok {
			Fail("Must set KMS_PLUGIN_CA_CERT")
		}
		certHost, ok := os.LookupEnv("KMS_PLUGIN_HOST")
		if !ok {
			Fail("Must set KMS_PLUGIN_HOST")
		}

		certPool := x509.NewCertPool()
		certPool.AppendCertsFromPEM([]byte(certString))
		creds := credentials.NewClientTLSFromCert(certPool, certHost)

		socket := "/var/vcap/sys/run/kms-plugin/kms-plugin.sock"
		conn, err := grpc.Dial("unix://"+socket, grpc.WithTransportCredentials(creds))
		Expect(err).NotTo(HaveOccurred())

		client := pb.NewKeyManagementServiceClient(conn)
		encryptResponse, err := client.Encrypt(context.Background(), &pb.EncryptRequest{Plain: []byte("Test")})
		Expect(err).NotTo(HaveOccurred())

		decryptResponse, err := client.Decrypt(context.Background(), &pb.DecryptRequest{Cipher: encryptResponse.Cipher})
		Expect(err).NotTo(HaveOccurred())

		Expect(string(decryptResponse.Plain)).To(Equal("Test"))
	})

})
