package vcr

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/shufanhao/go-example/vcr/testclient"
	"gopkg.in/dnaeon/go-vcr.v3/cassette"
	"gopkg.in/dnaeon/go-vcr.v3/recorder"
)

var testRSACertificate = fromHex("3082024b308201b4a003020102020900e8f09d3fe25beaa6300d06092a864886f70d01010b0500301f310b3009060355040a1302476f3110300e06035504031307476f20526f6f74301e170d3136303130313030303030305a170d3235303130313030303030305a301a310b3009060355040a1302476f310b300906035504031302476f30819f300d06092a864886f70d010101050003818d0030818902818100db467d932e12270648bc062821ab7ec4b6a25dfe1e5245887a3647a5080d92425bc281c0be97799840fb4f6d14fd2b138bc2a52e67d8d4099ed62238b74a0b74732bc234f1d193e596d9747bf3589f6c613cc0b041d4d92b2b2423775b1c3bbd755dce2054cfa163871d1e24c4f31d1a508baab61443ed97a77562f414c852d70203010001a38193308190300e0603551d0f0101ff0404030205a0301d0603551d250416301406082b0601050507030106082b06010505070302300c0603551d130101ff0402300030190603551d0e041204109f91161f43433e49a6de6db680d79f60301b0603551d230414301280104813494d137e1631bba301d5acab6e7b30190603551d1104123010820e6578616d706c652e676f6c616e67300d06092a864886f70d01010b0500038181009d30cc402b5b50a061cbbae55358e1ed8328a9581aa938a495a1ac315a1a84663d43d32dd90bf297dfd320643892243a00bccf9c7db74020015faad3166109a276fd13c3cce10c5ceeb18782f16c04ed73bbb343778d0c1cf10fa1d8408361c94c722b9daedb4606064df4c1b33ec0d1bd42d4dbfe3d1360845c21d33be9fae7")
var testRSAPrivateKey, _ = x509.ParsePKCS1PrivateKey(fromHex("3082025b02010002818100db467d932e12270648bc062821ab7ec4b6a25dfe1e5245887a3647a5080d92425bc281c0be97799840fb4f6d14fd2b138bc2a52e67d8d4099ed62238b74a0b74732bc234f1d193e596d9747bf3589f6c613cc0b041d4d92b2b2423775b1c3bbd755dce2054cfa163871d1e24c4f31d1a508baab61443ed97a77562f414c852d702030100010281800b07fbcf48b50f1388db34b016298b8217f2092a7c9a04f77db6775a3d1279b62ee9951f7e371e9de33f015aea80660760b3951dc589a9f925ed7de13e8f520e1ccbc7498ce78e7fab6d59582c2386cc07ed688212a576ff37833bd5943483b5554d15a0b9b4010ed9bf09f207e7e9805f649240ed6c1256ed75ab7cd56d9671024100fded810da442775f5923debae4ac758390a032a16598d62f059bb2e781a9c2f41bfa015c209f966513fe3bf5a58717cbdb385100de914f88d649b7d15309fa49024100dd10978c623463a1802c52f012cfa72ff5d901f25a2292446552c2568b1840e49a312e127217c2186615aae4fb6602a4f6ebf3f3d160f3b3ad04c592f65ae41f02400c69062ca781841a09de41ed7a6d9f54adc5d693a2c6847949d9e1358555c9ac6a8d9e71653ac77beb2d3abaf7bb1183aa14278956575dbebf525d0482fd72d90240560fe1900ba36dae3022115fd952f2399fb28e2975a1c3e3d0b679660bdcb356cc189d611cfdd6d87cd5aea45aa30a2082e8b51e94c2f3dd5d5c6036a8a615ed0240143993d80ece56f877cb80048335701eb0e608cc0c1ca8c2227b52edf8f1ac99c562f2541b5ce81f0515af1c5b4770dba53383964b4b725ff46fdec3d08907df"))

func fromHex(s string) []byte {
	b, _ := hex.DecodeString(s)
	return b
}

func TestHTTPIBasicInteraction(t *testing.T) {
	// create a recorder
	rec, err := recorder.NewWithOptions(&recorder.Options{
		CassetteName: fmt.Sprintf("cassettes/%s", t.Name()),
		// only record once, so next time if you disconnect network, it should work.
		Mode:               recorder.ModeRecordOnce,
		SkipRequestLatency: true,
	})

	rec.SetMatcher(func(req *http.Request, c cassette.Request) bool {
		// ignore hostname prefixes and URI paths on replays
		return req.Method == c.Method
	})

	if err != nil {
		t.Fatalf("Failed to create recorder: %v", err)
	}

	defer func() {
		if err := rec.Stop(); err != nil {
			log.Printf("Failed to stop recorder: %v", err)
		}
	}()

	client := &http.Client{
		Transport: rec,
	}

	resp, err := client.Get("https://baidu.com")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	log.Println("Request made successfully")
}

func TestRetryHttpIntercept(t *testing.T) {
	httpClient := cleanhttp.DefaultPooledClient()
	transport := httpClient.Transport.(*http.Transport)
	transport.MaxIdleConnsPerHost = 10
	tlsConfig := transport.TLSClientConfig

	c0 := tls.Certificate{
		Certificate: [][]byte{testRSACertificate},
		PrivateKey:  testRSAPrivateKey,
	}
	// config tlsConfig
	if tlsConfig == nil {
		transport.TLSClientConfig = &tls.Config{
			Certificates: []tls.Certificate{c0},
			MinVersion:   tls.VersionTLS12,
			// RootCAs:      "",
		}
	}

	rec, err := recorder.NewWithOptions(&recorder.Options{
		CassetteName:       fmt.Sprintf("cassettes/%s", t.Name()),
		Mode:               recorder.ModeRecordOnly,
		SkipRequestLatency: true,
		RealTransport:      httpClient.Transport,
	})

	rec.SetMatcher(func(req *http.Request, c cassette.Request) bool {
		// ignore hostname prefixes and URI paths on replays
		return req.Method == c.Method
	})

	if err != nil {
		t.Fatalf("Failed to create recorder: %v", err)
	}

	defer func() {
		if err := rec.Stop(); err != nil {
			log.Printf("Failed to stop recorder: %v", err)
		}
	}()

	// without this code, it doesn't work
	httpClient.Transport = rec

	client, _ := testclient.NewTestClient(testclient.Config{
		TLSCert: nil,
		URL:     "https://google.com",
	})

	client.SetHttpClient(httpClient)
	_, err = client.Get()

	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
}
