package server

import (
	"bytes"
	"crypto/tls"
	"errors"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

// Generated cert and key via the crypto/tls package:
//   go run generate_cert.go  --rsa-bits 1024 --host 127.0.0.1,::1,example.com --ca --start-date "Jan 1 00:00:00 1970" --duration=1000000h
var tlsCert = []byte(`-----BEGIN CERTIFICATE-----
MIICNTCCAZ6gAwIBAgIRAPc89REgYR2GEXgKW7Ebd/QwDQYJKoZIhvcNAQELBQAw
EjEQMA4GA1UEChMHQWNtZSBDbzAgFw03MDAxMDEwMDAwMDBaGA8yMDg0MDEyOTE2
MDAwMFowEjEQMA4GA1UEChMHQWNtZSBDbzCBnzANBgkqhkiG9w0BAQEFAAOBjQAw
gYkCgYEAmv2AL8b6AUo70zOto2P8rj3mowcAKQt2ZWnMypDgO+ST7FX4eCTMZVhm
2378Z6ukTtQCPNgV+WW71s5YhE8ipEmrVddO/9TseuDEFdsMVuAUYGGjz0ObCpmd
Ex/zMyZ6+nniZEUMSaogwXnggmh+wQQu7Pj1Lp6rdLtmE0IZEd8CAwEAAaOBiDCB
hTAOBgNVHQ8BAf8EBAMCAqQwEwYDVR0lBAwwCgYIKwYBBQUHAwEwDwYDVR0TAQH/
BAUwAwEB/zAdBgNVHQ4EFgQU50ZBQcexLwwYcRcKMQg4zNnn3/gwLgYDVR0RBCcw
JYILZXhhbXBsZS5jb22HBH8AAAGHEAAAAAAAAAAAAAAAAAAAAAEwDQYJKoZIhvcN
AQELBQADgYEANKux/Z//CHjzwSS01gyi3dbN9k9YIo0ineyvLWxgfoJ7rdBk3yi2
fc50pyot4l07A62+axbq/jdpfjp2WrrYcHZ2Zte/hTmu5cHdVJ/ATwGlHnuX9luN
T+VAAkUtkNNZF0ENGYyqozoroiSY1YavXZRYjpAdDM/QZE3dtznMnag=
-----END CERTIFICATE-----`)

var tlsKey = []byte(`-----BEGIN PRIVATE KEY-----
MIICdwIBADANBgkqhkiG9w0BAQEFAASCAmEwggJdAgEAAoGBAJr9gC/G+gFKO9Mz
raNj/K495qMHACkLdmVpzMqQ4Dvkk+xV+HgkzGVYZtt+/GerpE7UAjzYFfllu9bO
WIRPIqRJq1XXTv/U7HrgxBXbDFbgFGBho89DmwqZnRMf8zMmevp54mRFDEmqIMF5
4IJofsEELuz49S6eq3S7ZhNCGRHfAgMBAAECgYEAjYsTJQ7XRJRRvFjGq5/lpA7d
8Wa0S8evEYCkyR2z2p6uRLhimP4qOCeaj7wRsG+9N0xA2uYQc1noSIzbo8FNhUXx
Se99ZkxhL0lcf9O1N+IXdKSAUwd8grOx4hePicImnTsfbIHApF11ia71gErWC6pC
N31SxNBkGjgk/XnjRSECQQDEYz0S/H07y97H7W3P0roTpuyyQ34LM//mqczjVOtH
sDy280C+HK332yTrbsBWhnc1XsJTVvmBVL43fg12rWnRAkEAyglhH1iJqkkqBg/n
IPHmtooJxEm9Tc7a5SnElxw1AA4J2e0Q3VHHrv71bGwJzIOJnpMLkwyXQ84DMXKn
zCT8rwJBAKJZb8ncVSMzvG4Gx7sAd7d7TT1pMi/MwwZK5Qgh0YcoGGkd2y3Ow7qc
vX6rxfFBpBbIyVjgE89o4P87A6JSZaECQC0fbSagEpcKoi0abieIf1uzdrz1Lydq
lK7r5jFNpnStrfuTG9oiQrPN06h0dLfNhPX6p3IfNMV7BHGNxvYsKJcCQEj7nzCT
+aUvLbPGtyNfQAAuTh39Pj6WEw88xs8bVZrFdiOIpTvbUvP4f3RTJh+mz5NFSIe6
RROnfM47FiKPQTE=
-----END PRIVATE KEY-----`)

func TestListenAndServeStarts(t *testing.T) {
	srv := Static{}

	go srv.ListenAndServe()
	defer srv.Shutdown(0)
	time.Sleep(100 * time.Millisecond)

	if srv.Port == 0 {
		t.Error("ListenAndServe did not update port")
	}

	resp, err := http.Get("http://localhost:" + strconv.Itoa(srv.Port) + "/camera")
	if err != nil {
		t.Error("Request to server failed:", err)
	}
	if resp.StatusCode != 200 {
		t.Error("Request to server failed with status code", resp.StatusCode)
	}
}

func TestListenAndServeServesFile(t *testing.T) {
	tempDir := t.TempDir()
	srv := Static{
		Directory: tempDir,
	}

	go srv.ListenAndServe()
	defer srv.Shutdown(0)
	time.Sleep(100 * time.Millisecond)

	// Set up a file that we can serve
	fileContent := []byte("get in the robot! 🤖")
	ioutil.WriteFile(filepath.Join(tempDir, "instructions.txt"), fileContent, 0644)

	resp, err := http.Get("http://localhost:" + strconv.Itoa(srv.Port) + "/camera/instructions.txt")
	if err != nil {
		t.Fatal("Request to server failed:", err)
	}
	if resp.StatusCode != 200 {
		t.Error("Request to server failed with status code", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	if bytes.Compare(fileContent, body) != 0 {
		t.Error("Response body did not match, given:", body)
	}
}

func TestListenAndServeServesFileWithTLS(t *testing.T) {
	tempDir := t.TempDir()

	// Set up certificates for TLS
	tlsCertFile := filepath.Join(tempDir, "cert.pem")
	ioutil.WriteFile(tlsCertFile, tlsCert, 0644)

	tlsKeyFile := filepath.Join(tempDir, "key.pem")
	ioutil.WriteFile(tlsKeyFile, tlsKey, 0644)

	srv := Static{
		Directory: tempDir,
		Cert:      tlsCertFile,
		Key:       tlsKeyFile,
	}

	go srv.ListenAndServe()
	defer srv.Shutdown(0)
	time.Sleep(100 * time.Millisecond)

	// Set up a file that we can serve
	fileContent := []byte("get in the robot! 🤖")
	ioutil.WriteFile(filepath.Join(tempDir, "instructions.txt"), fileContent, 0644)

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := http.Client{
		Transport: transport,
	}

	resp, err := client.Get("https://localhost:" + strconv.Itoa(srv.Port) + "/camera/instructions.txt")
	if err != nil {
		t.Fatal("Request to server failed:", err)
	}
	if resp.StatusCode != 200 {
		t.Error("Request to server failed with status code", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	if bytes.Compare(fileContent, body) != 0 {
		t.Error("Response body did not match, given:", body)
	}
}

func TestListenAndServeReturns404(t *testing.T) {
	srv := Static{}

	go srv.ListenAndServe()
	defer srv.Shutdown(0)
	time.Sleep(100 * time.Millisecond)

	resp, err := http.Get("http://localhost:" + strconv.Itoa(srv.Port) + "/camera/instructions.txt")
	if err != nil {
		t.Fatal("Request to server failed:", err)
	}

	if resp.StatusCode != 404 {
		t.Error("Request to server failed with status code", resp.StatusCode)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	bodyText := string(body)
	if bodyText != "404 page not found\n" {
		t.Error("Response body did not match, given:", bodyText)
	}
}

func TestListenAndServeReturnsErrForInvalidDirectory(t *testing.T) {
	srv := Static{
		Directory: "totallybaddirectory",
	}

	err := srv.ListenAndServe()
	defer srv.Shutdown(0)

	if !errors.Is(err, ErrInvalidDirectory) {
		t.Error("ListenAndServe failed to return an error for an invalid directory:", err)
	}
}

func TestShutdownReleasesPort(t *testing.T) {
	srv := Static{}

	go srv.ListenAndServe()
	time.Sleep(100 * time.Millisecond)

	if srv.Port == 0 {
		t.Error("ListenAndServe did not update port")
	}

	srv.Shutdown(0)

	conn, err := net.Dial("tcp", "localhost:"+strconv.Itoa(srv.Port))
	if conn != nil {
		conn.Close()
	}
	if err == nil {
		t.Error("Request to server failed:", err)
	}
}
