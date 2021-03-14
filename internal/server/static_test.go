package server

import (
	"bytes"
	"crypto/tls"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

// It's important to run these tests with a timeout since we don't want to deadlock

// Generated cert and key via:
//   go run generate_cert.go  --rsa-bits 1024 --host 127.0.0.1,::1,example.com --ca --start-date "Jan 1 00:00:00 1970" --duration=1000000h
var tlsCert = []byte(`-----BEGIN CERTIFICATE-----
MIIBwDCCASmgAwIBAgIJAKDwOOc421AyMA0GCSqGSIb3DQEBCwUAMBQxEjAQBgNV
BAMMCWxvY2FsaG9zdDAgFw0yMTAzMTQwMzE5MDhaGA8yMTE4MDkwMjAzMTkwOFow
FDESMBAGA1UEAwwJbG9jYWxob3N0MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKB
gQCcOiwhdWEn2G0Vz68PYq3Xa6GhyDN6DRbswdCkaoV7fptTuPME+wUDpqDcUGY3
w7jUs7LeIScQuu2jKqWXs38X+TazDJHIZqlNb23BQJtRwAorN8lb/bUqJdPhuCmi
Ht4CV2QNIENAwFZT7gi2E9dOQTpsvqQ2EhXXPkn+o2k17wIDAQABoxgwFjAUBgNV
HREEDTALgglsb2NhbGhvc3QwDQYJKoZIhvcNAQELBQADgYEAiHRL2HhKf0+5B4ib
6p+BQRpvtFlLQWE7sfS0zIQU3zuwcAvSt559o5/XnCtPXrpl/vyXJobDMhL1UWXi
9GxFaloCyyhCSX+3X3ZR0zykZ5v2wdtt1+QDdteSktYQs4enJjYGuBerl+cMDCVh
j0voFP26qanpsXcBBd9Ed6qrHZg=
-----END CERTIFICATE-----`)

var tlsKey = []byte(`-----BEGIN PRIVATE KEY-----
MIICdgIBADANBgkqhkiG9w0BAQEFAASCAmAwggJcAgEAAoGBAJw6LCF1YSfYbRXP
rw9irddroaHIM3oNFuzB0KRqhXt+m1O48wT7BQOmoNxQZjfDuNSzst4hJxC67aMq
pZezfxf5NrMMkchmqU1vbcFAm1HACis3yVv9tSol0+G4KaIe3gJXZA0gQ0DAVlPu
CLYT105BOmy+pDYSFdc+Sf6jaTXvAgMBAAECgYBKigLmT9/0J/IcNwQ6ngV9w+/R
hfjxoe8vNWY8HErl2kl4D8V7udzxmf4kQAQVVUAJ1FdiXoUKGXYqpL1vfQgFBF9x
gwIwo/h65S+DZgpE9/IudPu4Kir3Fg07/3vxOcpLXXs7elqyAjXIV8beati8LJqn
xxCopADRSbApdo1haQJBAMjEu7bRQMhkXa90fvbiRXQyR0vlHim6HUjvYs08ScIn
nzkNzBHzKmQihln7ChL/D71lUAJtS3xzQy3yad0Bh6UCQQDHNJbIHnh0XftSJkyV
LVpRGzmrxPUrTWjcPCdgS/kf7YA92Aj+m9So+2G54kdxxnY5kcM+wt4ig3tPcans
p/MDAkB7RpmActJVeZMw9dYz39IHvAudJW007+ulah//p0Ie7ldNIBSq/OWNoMlg
HM4dxfGzOK89HkEYhGm+n7ezFYplAkEAl8s/9mAZo3qV7qRWiPoVL2aSjIw50fRb
qi6ARsW9oRGmPfnn6LOv2dAsSKvfixgSsI2c/K8a+6u7A+9172qPJwJAUxiklMeB
am/1vWIAKgHRZoSJe6O/FGnWZY+axa64tOTd+HWhKFmWLa8TsG7Qo2TDqDPGhMSQ
JU3hq9/zrQddjg==
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

	if err.Error() != "directory does not exist" {
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
