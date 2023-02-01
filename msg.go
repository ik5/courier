package courier

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/null/v2"
)

// ErrMsgNotFound is returned when trying to queue the status for a Msg that doesn't exit
var ErrMsgNotFound = errors.New("message not found")

// ErrWrongIncomingMsgStatus use do ignore the status update if the DB raise this
var ErrWrongIncomingMsgStatus = errors.New("incoming messages can only be PENDING or HANDLED")

// MsgID is our typing of the db int type
type MsgID null.Int64

// NewMsgID creates a new MsgID for the passed in int64
func NewMsgID(id int64) MsgID {
	return MsgID(id)
}

// String satisfies the Stringer interface
func (i MsgID) String() string {
	if i != NilMsgID {
		return strconv.FormatInt(int64(i), 10)
	}
	return "null"
}

func (i *MsgID) Scan(value any) error         { return null.ScanInt(value, i) }
func (i MsgID) Value() (driver.Value, error)  { return null.IntValue(i) }
func (i *MsgID) UnmarshalJSON(b []byte) error { return null.UnmarshalInt(b, i) }
func (i MsgID) MarshalJSON() ([]byte, error)  { return null.MarshalInt(i) }

// NilMsgID is our nil value for MsgID
var NilMsgID = MsgID(0)

// MsgUUID is the UUID of a message which has been received
type MsgUUID struct {
	uuid.UUID
}

// NilMsgUUID is a "zero value" message UUID
var NilMsgUUID = MsgUUID{uuid.Nil}

// NewMsgUUID creates a new unique message UUID
func NewMsgUUID() MsgUUID {
	u, _ := uuid.NewV4()
	return MsgUUID{u}
}

// NewMsgUUIDFromString creates a new message UUID for the passed in string
func NewMsgUUIDFromString(uuidString string) MsgUUID {
	uuid, _ := uuid.FromString(uuidString)
	return MsgUUID{uuid}
}

type FlowReference struct {
	UUID string `json:"uuid" validate:"uuid4"`
	Name string `json:"name"`
}

//-----------------------------------------------------------------------------
// Locale
//-----------------------------------------------------------------------------

// Locale is the combination of a language and optional country, e.g. US English, Brazilian Portuguese, encoded as the
// language code followed by the country code, e.g. eng-US, por-BR
type Locale string

func (l Locale) ToParts() (string, string) {
	if l == NilLocale || len(l) < 3 {
		return "", ""
	}

	parts := strings.SplitN(string(l), "-", 2)
	lang := parts[0]
	country := ""
	if len(parts) > 1 {
		country = parts[1]
	}

	return lang, country
}

var NilLocale = Locale("")

//-----------------------------------------------------------------------------
// Msg interface
//-----------------------------------------------------------------------------

// Msg is our interface to represent an incoming or outgoing message
type Msg interface {
	ID() MsgID
	UUID() MsgUUID
	Text() string
	Attachments() []string
	Locale() Locale
	ExternalID() string
	URN() urns.URN
	URNAuth() string
	ContactName() string
	QuickReplies() []string
	Topic() string
	Metadata() json.RawMessage
	ResponseToExternalID() string
	IsResend() bool

	Flow() *FlowReference
	FlowName() string
	FlowUUID() string

	Channel() Channel

	ReceivedOn() *time.Time
	SentOn() *time.Time

	HighPriority() bool

	WithContactName(name string) Msg
	WithReceivedOn(date time.Time) Msg
	WithExternalID(id string) Msg
	WithID(id MsgID) Msg
	WithUUID(uuid MsgUUID) Msg
	WithAttachment(url string) Msg
	WithLocale(Locale) Msg
	WithURNAuth(auth string) Msg
	WithMetadata(metadata json.RawMessage) Msg
	WithFlow(flow *FlowReference) Msg

	EventID() int64
	SessionStatus() string
}
