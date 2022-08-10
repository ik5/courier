package burstsms

import (
	"net/http/httptest"
	"testing"

	"github.com/nyaruka/courier"
	. "github.com/nyaruka/courier/handlers"
)

var testChannels = []courier.Channel{
	courier.NewMockChannel("8eb23e93-5ecb-45ba-b726-3b064e0c56ab", "BS", "2020", "US", nil),
}

var (
	receiveURL = "/c/bs/8eb23e93-5ecb-45ba-b726-3b064e0c56ab/receive/"

	validReceive  = "response=Msg&mobile=254791541111"
	missingNumber = "response=Msg"

	statusURL = "/c/bs/8eb23e93-5ecb-45ba-b726-3b064e0c56ab/status/"

	validStatus   = "message_id=12345&status=pending"
	invalidStatus = "message_id=12345&status=unknown"
)

var testCases = []ChannelHandleTestCase{
	{Label: "Receive Valid", URL: receiveURL + "?" + validReceive, Status: 200, Response: "Message Accepted",
		Text: Sp("Msg"), URN: Sp("tel:+254791541111")},
	{Label: "Receive Missing Number", URL: receiveURL + "?" + missingNumber, Status: 400, Response: "required field 'mobile'"},

	{Label: "Status Valid", URL: statusURL + "?" + validStatus, Status: 200, Response: "Status Update Accepted",
		ExternalID: Sp("12345"), MsgStatus: Sp("S")},
	{Label: "Receive Invalid Status", URL: statusURL + "?" + invalidStatus, Status: 400, Response: "unknown status value"},
}

func TestHandler(t *testing.T) {
	RunChannelTestCases(t, testChannels, newHandler(), testCases)
}

func BenchmarkHandler(b *testing.B) {
	RunChannelBenchmarks(b, testChannels, newHandler(), testCases)
}

func setSendURL(s *httptest.Server, h courier.ChannelHandler, c courier.Channel, m courier.Msg) {
	sendURL = s.URL
}

var defaultSendTestCases = []ChannelSendTestCase{
	{Label: "Plain Send",
		MsgText: "Simple Message ☺", MsgURN: "tel:+250788383383", MsgAttachments: []string{"image/jpeg:https://foo.bar/image.jpg"},
		ExpectedStatus: "W", ExpectedExternalID: "19835",
		MockResponseBody: `{ "message_id": 19835, "recipients": 3, "cost": 1.000 }`, MockResponseStatus: 200,
		ExpectedPostParams: map[string]string{
			"to":      "250788383383",
			"message": "Simple Message ☺\nhttps://foo.bar/image.jpg",
			"from":    "2020",
		},
		SendPrep: setSendURL},
	{Label: "Invalid JSON",
		MsgText: "Invalid JSON", MsgURN: "tel:+250788383383",
		ExpectedStatus:   "E",
		MockResponseBody: `not json`, MockResponseStatus: 200,
		SendPrep: setSendURL},
	{Label: "Error Response",
		MsgText: "Error Response", MsgURN: "tel:+250788383383",
		ExpectedStatus:   "F",
		MockResponseBody: `{ "message_id": 0 }`, MockResponseStatus: 200,
		SendPrep: setSendURL},
	{Label: "Error Sending",
		MsgText: "Error Message", MsgURN: "tel:+250788383383",
		ExpectedStatus:   "E",
		MockResponseBody: `Bad Gateway`, MockResponseStatus: 501,
		SendPrep: setSendURL},
}

func TestSending(t *testing.T) {
	var defaultChannel = courier.NewMockChannel("8eb23e93-5ecb-45ba-b726-3b064e0c56ab", "BS", "2020", "US",
		map[string]interface{}{
			courier.ConfigUsername: "user1",
			courier.ConfigPassword: "pass1",
		})
	RunChannelSendTestCases(t, defaultChannel, newHandler(), defaultSendTestCases, nil)
}
