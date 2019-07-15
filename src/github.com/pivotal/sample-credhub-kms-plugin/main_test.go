package main_test

import (
	"context"
	"encoding/base64"
	"io/ioutil"
	"os"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "github.com/pivotal/sample-credhub-kms-plugin/v1beta1"
)

var _ = Describe("Main", func() {
	var startServer func(args ...string) *gexec.Session
	var socketAddr string
	var publicCert string
	var privateKey string
	BeforeEach(func() {
		executablePath, err := gexec.Build("github.com/pivotal/sample-credhub-kms-plugin")
		Expect(err).NotTo(HaveOccurred())

		publicCert = "./grpc-kms-certs/grpc_kms_server_cert.pem"
		privateKey = "./grpc-kms-certs/grpc_kms_server_key.pem"

		startServer = func(args ...string) *gexec.Session {
			cmd := exec.Command(executablePath, args...)
			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			return session
		}

		tempFile, err := ioutil.TempFile("", "sample-credhub-kms-plugin")
		Expect(err).NotTo(HaveOccurred())
		socketAddr = tempFile.Name()
	})

	AfterEach(func() {
		gexec.Kill()
		os.Remove(socketAddr)
	})

	It("Starts the CLI", func() {
		session := startServer(socketAddr, publicCert, privateKey)
		Eventually(session.Err).Should(gbytes.Say("Listening on unix domain socket"))

		creds, err := credentials.NewClientTLSFromFile(publicCert, "localhost")
		Expect(err).NotTo(HaveOccurred())

		channel, err := grpc.Dial("unix://"+socketAddr, grpc.WithTransportCredentials(creds))
		Expect(err).NotTo(HaveOccurred())

		client := pb.NewKeyManagementServiceClient(channel)

		response, err := client.Encrypt(context.Background(), &pb.EncryptRequest{Plain: []byte("Test")})
		Expect(err).NotTo(HaveOccurred())

		verifyResp := base64.StdEncoding.EncodeToString([]byte("Test"))
		Expect(string(response.Cipher)).To(Equal(verifyResp))
	})

	It("exits gracefully when it receives a terminate signal", func() {
		session := startServer(socketAddr, publicCert, privateKey)
		Eventually(session.Err).Should(gbytes.Say("Listening on unix domain socket"))

		session.Terminate()
		Eventually(session).Should(gexec.Exit(0))
		Expect(session.Err).To(gbytes.Say("Stopped gRPC server"))
	})

	It("exits gracefully when it receives an interrupt signal", func() {
		session := startServer(socketAddr, publicCert, privateKey)
		Eventually(session.Err).Should(gbytes.Say("Listening on unix domain socket"))

		session.Interrupt()
		Eventually(session).Should(gexec.Exit(0))
		Expect(session.Err).To(gbytes.Say("Stopped gRPC server"))
	})

	Context("when no argument is provided", func() {
		It("prints a usage message", func() {
			session := startServer()
			Eventually(session).Should(gexec.Exit(1))
			Expect(session.Err).To(gbytes.Say("Usage: .* <path-to-unix-socket> <public-key> <private-key>"))
		})
	})
})
