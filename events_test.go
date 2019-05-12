package goboxer

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestBoxEvent_ResourceType(t *testing.T) {
	boxEvent := BoxEvent{}
	if boxEvent.ResourceType() != EventResource {
		t.Errorf("not match")
	}
}

func TestBoxEvent_Source(t *testing.T) {
	json001 := `
{
	"type": "folder",
	"id": "11446498",
	"sequence_id": "0",
	"etag": "0",
	"name": "Pictures",
	"created_at": "2012-12-12T10:53:43-08:00",
	"modified_at": "2012-12-12T10:53:43-08:00",
	"description": null,
	"size": 0,
	"created_by": {
		"type": "user",
		"id": "17738362",
		"name": "sean rose",
		"login": "sean@box.com"
	},
	"modified_by": {
		"type": "user",
		"id": "17738362",
		"name": "sean rose",
		"login": "sean@box.com"
	},
	"owned_by": {
		"type": "user",
		"id": "17738362",
		"name": "sean rose",
		"login": "sean@box.com"
	},
	"shared_link": null,
	"parent": {
		"type": "folder",
		"id": "0",
		"sequence_id": null,
		"etag": null,
		"name": "All Files"
	},
	"item_status": "active",
	"synced": false
	}
`
	var folder Folder
	_ = json.Unmarshal([]byte(json001), &folder)
	type fields struct {
		Type              string
		EventID           string
		CreatedBy         *UserGroupMini
		EventType         EventType
		SessionID         *string
		SourceRaw         json.RawMessage
		source            BoxResource
		AdditionalDetails map[string]interface{}
		CreatedAt         time.Time
		RecordedAt        *time.Time
		ActionBy          *UserGroupMini
	}
	tests := []struct {
		name   string
		fields fields
		want   BoxResource
	}{
		{"normal",
			fields{
				Type:    "event",
				EventID: "EVENTID001",
				CreatedBy: &UserGroupMini{
					Type:  setUserType(TYPE_USER),
					ID:    setStringPtr("17738362"),
					Name:  setStringPtr("sean rose"),
					Login: setStringPtr("sean@box.com"),
				},
				CreatedAt:  *setTime("2012-12-12T10:53:43-08:00"),
				RecordedAt: setTime("2012-12-12T10:53:48-08:00"),
				EventType:  UE_ITEM_CREATE,
				SessionID:  setStringPtr("70090280850c8d2a1933c1"),
				SourceRaw:  []byte(json001),
			},
			&folder,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			be := &BoxEvent{
				Type:              tt.fields.Type,
				EventID:           tt.fields.EventID,
				CreatedBy:         tt.fields.CreatedBy,
				EventType:         tt.fields.EventType,
				SessionID:         tt.fields.SessionID,
				SourceRaw:         tt.fields.SourceRaw,
				source:            tt.fields.source,
				AdditionalDetails: tt.fields.AdditionalDetails,
				CreatedAt:         tt.fields.CreatedAt,
				RecordedAt:        tt.fields.RecordedAt,
				ActionBy:          tt.fields.ActionBy,
			}
			if got := be.Source(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BoxEvent.Source() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewEvent(t *testing.T) {
	url := "https://example.com"
	apiConn := commonInit(url)
	type args struct {
		api *APIConn
	}
	tests := []struct {
		name string
		args args
		want *Event
	}{
		{"normal", args{api: apiConn}, &Event{apiInfo: &apiInfo{api: apiConn}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewEvent(tt.args.api); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEvent_UserEvent(t *testing.T) {
	// test server (dummy box api)
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// for request.Send() return error (auth failed)
			if strings.HasPrefix(r.URL.Path, "/oauth2/token") {
				w.WriteHeader(401)
				return
			}
			// URL check
			if !strings.HasPrefix(r.URL.Path, "/2.0/events") {
				t.Errorf("invalid access url %s : %s", r.URL.Path, "/2.0/events")
			}
			// Method check
			if r.Method != http.MethodGet {
				t.Fatalf("invalid http method")
			}
			// Header check
			if r.Header.Get("Authorization") == "" {
				t.Fatalf("not exists access token")
			}
			// ok, return some response
			streamPosition := r.URL.Query().Get("stream_position")

			switch streamPosition {
			case "500":
				w.WriteHeader(500)
			case "400":
				w.Header().Set("content-Type", "application/json")
				w.WriteHeader(400)
				resp, _ := ioutil.ReadFile("testdata/genericerror/400_notempty.json")
				_, _ = w.Write(resp)
			case "999":
				w.Header().Set("content-Type", "application/json")
				_, _ = w.Write([]byte("invalid json"))
			default:
				w.Header().Set("content-Type", "application/json")
				resp, _ := ioutil.ReadFile("testdata/events/events_user_normal.json")
				_, _ = w.Write(resp)
			}
			return
		},
	))
	defer ts.Close()

	apiConn := commonInit(ts.URL)
	api := &apiInfo{api: apiConn}

	json001 := `
{
	"type": "folder",
	"id": "11446498",
	"sequence_id": "0",
	"etag": "0",
	"name": "Pictures",
	"created_at": "2012-12-12T10:53:43-08:00",
	"modified_at": "2012-12-12T10:53:43-08:00",
	"description": null,
	"size": 0,
	"created_by": {
		"type": "user",
		"id": "17738362",
		"name": "sean rose",
		"login": "sean@box.com"
	},
	"modified_by": {
		"type": "user",
		"id": "17738362",
		"name": "sean rose",
		"login": "sean@box.com"
	},
	"owned_by": {
		"type": "user",
		"id": "17738362",
		"name": "sean rose",
		"login": "sean@box.com"
	},
	"shared_link": null,
	"parent": {
		"type": "folder",
		"id": "0",
		"sequence_id": null,
		"etag": null,
		"name": "All Files"
	},
	"item_status": "active",
	"synced": false
	}
`
	var folder Folder
	_ = json.Unmarshal([]byte(json001), &folder)
	folder.apiInfo = api

	type fields struct {
		apiInfo *apiInfo
	}
	type args struct {
		streamType     StreamType
		streamPosition string
		limit          int
	}
	tests := []struct {
		name                   string
		fields                 fields
		args                   args
		wantEvents             []*BoxEvent
		wantNextStreamPosition string
		wantErr                bool
		errType                interface{}
		status                 int
	}{
		{"normal",
			fields{api},
			args{Changes, "now", 1000},
			[]*BoxEvent{
				{
					Type:    "event",
					EventID: "f82c3ba03e41f7e8a7608363cc6c0390183c3f83",
					CreatedBy: &UserGroupMini{
						Type:  setUserType(TYPE_USER),
						ID:    setStringPtr("17738362"),
						Name:  setStringPtr("sean rose"),
						Login: setStringPtr("sean@box.com"),
					},
					EventType:         UE_ITEM_CREATE,
					SessionID:         setStringPtr("70090280850c8d2a1933c1"),
					SourceRaw:         nil,
					source:            &folder,
					AdditionalDetails: nil,
					CreatedAt:         *setTime("2012-12-12T10:53:43-08:00"),
					RecordedAt:        setTime("2012-12-12T10:53:48-08:00"),
					ActionBy:          nil,
				},
			},
			"1348790499819",
			false,
			nil,
			200,
		},
		{"http_error/400",
			fields{api},
			args{All, "400", 10},
			nil,
			"",
			true,
			&ApiStatusError{},
			400,
		},
		{"returned_invalid_json/999",
			fields{api},
			args{All, "999", 10},
			nil,
			"",
			true,
			&ApiOtherError{},
			0,
		},
		{"senderror",
			fields{api},
			args{All, "999", 10},
			nil,
			"",
			true,
			&ApiOtherError{},
			0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "senderror" {
				apiConn.Expires = 0
			} else {
				apiConn.Expires = 6000
			}
			e := &Event{
				apiInfo: tt.fields.apiInfo,
			}
			gotEvents, gotNextStreamPosition, err := e.UserEvent(tt.args.streamType, tt.args.streamPosition, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("Event.UserEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errType != nil {
				if reflect.TypeOf(err).String() != reflect.TypeOf(tt.errType).String() {
					t.Errorf("got err = %v, wanted errorType %v", err, tt.errType)
					return
				}
				if reflect.TypeOf(tt.errType) == reflect.TypeOf(&ApiStatusError{}) {
					apiStatusError := err.(*ApiStatusError)
					if tt.status != apiStatusError.Status {
						t.Errorf("status code may be not corrected [%d]", apiStatusError.Status)
						return
					}
					return
				} else {
					return
				}
			} else if err != nil {
				return
			}

			opts := diffCompOptions(Folder{}, BoxEvent{})
			ignore := cmpopts.IgnoreUnexported(apiInfo{})
			opts = append(opts, ignore)
			if diff := cmp.Diff(&gotEvents, &tt.wantEvents, opts...); diff != "" {
				t.Errorf("Marshal/Unmarshal differs: (-got +want)\n%s", diff)
				return
			}
			if gotNextStreamPosition != tt.wantNextStreamPosition {
				t.Errorf("Event.UserEvent() gotNextStreamPosition = %v, want %v", gotNextStreamPosition, tt.wantNextStreamPosition)
			}
		})
	}
}

func TestEvent_EnterpriseEvent(t *testing.T) {
	// test server (dummy box api)
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// for request.Send() return error (auth failed)
			if strings.HasPrefix(r.URL.Path, "/oauth2/token") {
				w.WriteHeader(401)
				return
			}
			// URL check
			if !strings.HasPrefix(r.URL.Path, "/2.0/events") {
				t.Errorf("invalid access url %s : %s", r.URL.Path, "/2.0/events")
			}
			// Method check
			if r.Method != http.MethodGet {
				t.Fatalf("invalid http method")
			}
			// Header check
			if r.Header.Get("Authorization") == "" {
				t.Fatalf("not exists access token")
			}
			// ok, return some response
			streamPosition := r.URL.Query().Get("stream_position")

			switch streamPosition {
			case "500":
				w.WriteHeader(500)
			case "400":
				w.Header().Set("content-Type", "application/json")
				w.WriteHeader(400)
				resp, _ := ioutil.ReadFile("testdata/genericerror/400_notempty.json")
				_, _ = w.Write(resp)
			case "999":
				w.Header().Set("content-Type", "application/json")
				_, _ = w.Write([]byte("invalid json"))
			default:
				w.Header().Set("content-Type", "application/json")
				resp, _ := ioutil.ReadFile("testdata/events/events_enterprise_normal.json")
				_, _ = w.Write(resp)
			}
			return
		},
	))
	defer ts.Close()

	apiConn := commonInit(ts.URL)
	api := &apiInfo{api: apiConn}

	userJson001 := `
{
  "type":"user",
  "id":"222853849",
  "name":"Nick Lee",
  "login":"nlee+demo4@box.com"
}`
	var user001 User
	_ = json.Unmarshal([]byte(userJson001), &user001)
	user001.apiInfo = api

	fileJson001 := `
{
  "type":"file",
  "id":"2385897476",
  "name":"Word Doc.docx",
  "parent":{
    "type":"folder",
    "name":"Folder 1",
    "id":"8430655483"
  },
  "owned_by":{
    "type":"user",
    "id":"222853849",
    "name":"sean rose",
    "login":"sean+awesome@box.com"
  }
}`
	var file001 File
	_ = json.Unmarshal([]byte(fileJson001), &file001)
	file001.apiInfo = api

	wantBes := []*BoxEvent{
		{
			Type:    "event",
			EventID: "b9a2393a-20cf-4307-90f5-004110dec209",
			CreatedBy: &UserGroupMini{
				Type:  setUserType(TYPE_USER),
				ID:    setStringPtr("222853849"),
				Name:  setStringPtr("sean rose"),
				Login: setStringPtr("sean+awesome@box.com"),
			},
			EventType:         EE_ADD_LOGIN_ACTIVITY_DEVICE,
			SessionID:         nil,
			SourceRaw:         nil,
			source:            &user001,
			AdditionalDetails: nil,
			CreatedAt:         *setTime("2015-12-02T09:41:31-08:00"),
			RecordedAt:        nil,
			ActionBy:          nil,
		},
		{
			Type:    "event",
			EventID: "1a4ade15-b1ff-4cc3-89a8-955e1522557c",
			CreatedBy: &UserGroupMini{
				Type:  setUserType(TYPE_USER),
				ID:    setStringPtr("222853849"),
				Name:  setStringPtr("Nick Lee"),
				Login: setStringPtr("nlee+demo4@box.com"),
			},
			EventType:         EE_DOWNLOAD,
			SessionID:         nil,
			SourceRaw:         nil,
			source:            &file001,
			AdditionalDetails: map[string]interface{}{"size": 21696.0, "ekm_id": "8afe49be-4ced-42cb-a3b0-8342a1cfe111", "version_id": "25012733090"},
			CreatedAt:         *setTime("2018-08-21T14:25:38-07:00"),
			RecordedAt:        nil,
			ActionBy:          nil,
		},
	}
	type fields struct {
		apiInfo *apiInfo
	}
	type args struct {
		streamPosition string
		eventTypes     []EventType
		createdAfter   *time.Time
		createdBefore  *time.Time
		limit          int
	}
	tests := []struct {
		name                   string
		fields                 fields
		args                   args
		wantEvents             []*BoxEvent
		wantNextStreamPosition string
		wantErr                bool
		errType                interface{}
		status                 int
	}{
		{"normal",
			fields{api},
			args{"0", []EventType{EE_COPY}, nil, nil, 10},
			wantBes,
			"1152922976252290886",
			false,
			nil,
			200,
		},
		{"normal/createdAfter",
			fields{api},
			args{"", nil, setTime("2018-08-21T14:25:38-07:00"), nil, 10},
			wantBes,
			"1152922976252290886",
			false,
			nil,
			200,
		},
		{"normal/createdBefore",
			fields{api},
			args{"", nil, nil, setTime("2018-08-21T14:25:38-07:00"), 10},
			wantBes,
			"1152922976252290886",
			false,
			nil,
			200,
		},
		{"normal/limit",
			fields{api},
			args{"", nil, nil, setTime("2018-08-21T14:25:38-07:00"), 1000},
			wantBes,
			"1152922976252290886",
			false,
			nil,
			200,
		},
		{"http_error/400",
			fields{api},
			args{"400", nil, nil, nil, 10},
			nil,
			"",
			true,
			&ApiStatusError{},
			400,
		},
		{"returned_invalid_json/999",
			fields{api},
			args{"999", nil, nil, nil, 10},
			nil,
			"",
			true,
			&ApiOtherError{},
			0,
		},
		{"senderror",
			fields{api},
			args{"999", nil, nil, nil, 10},
			nil,
			"",
			true,
			&ApiOtherError{},
			0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "senderror" {
				apiConn.Expires = 0
			} else {
				apiConn.Expires = 6000
			}
			e := &Event{
				apiInfo: tt.fields.apiInfo,
			}
			gotEvents, gotNextStreamPosition, err := e.EnterpriseEvent(tt.args.streamPosition, tt.args.eventTypes, tt.args.createdAfter, tt.args.createdBefore, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("Event.EnterpriseEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errType != nil {
				if reflect.TypeOf(err).String() != reflect.TypeOf(tt.errType).String() {
					t.Errorf("got err = %v, wanted errorType %v", err, tt.errType)
					return
				}
				if reflect.TypeOf(tt.errType) == reflect.TypeOf(&ApiStatusError{}) {
					apiStatusError := err.(*ApiStatusError)
					if tt.status != apiStatusError.Status {
						t.Errorf("status code may be not corrected [%d]", apiStatusError.Status)
						return
					}
					return
				} else {
					return
				}
			} else if err != nil {
				return
			}

			opts := diffCompOptions(Folder{}, BoxEvent{}, File{}, User{})
			ignore := cmpopts.IgnoreUnexported(apiInfo{})
			opts = append(opts, ignore)
			if diff := cmp.Diff(&gotEvents, &tt.wantEvents, opts...); diff != "" {
				t.Errorf("Marshal/Unmarshal differs: (-got +want)\n%s", diff)
				return
			}
			if gotNextStreamPosition != tt.wantNextStreamPosition {
				t.Errorf("Event.EnterpriseEvent() gotNextStreamPosition = %v, want %v", gotNextStreamPosition, tt.wantNextStreamPosition)
			}
		})
	}
}

func Test_buildEventTypesQueryParams(t *testing.T) {
	type args struct {
		eventTypes []EventType
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"nil", args{nil}, ""},
		{"empty list", args{[]EventType{}}, ""},
		{"empty list", args{[]EventType{EE_ACCESS_GRANTED, EE_COPY}}, "event_type=ACCESS_GRANTED,COPY"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildEventTypesQueryParams(tt.args.eventTypes); got != tt.want {
				t.Errorf("buildEventTypesQueryParams() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBoxEvent_String(t *testing.T) {
	fileJson001 := `
{
  "type":"file",
  "id":"2385897476",
  "name":"Word Doc.docx",
  "parent":{
    "type":"folder",
    "name":"Folder 1",
    "id":"8430655483"
  },
  "owned_by":{
    "type":"user",
    "id":"222853849",
    "name":"sean rose",
    "login":"sean+awesome@box.com"
  }
}`

	type fields struct {
		Type              string
		EventID           string
		CreatedBy         *UserGroupMini
		EventType         EventType
		SessionID         *string
		SourceRaw         json.RawMessage
		source            BoxResource
		AdditionalDetails map[string]interface{}
		CreatedAt         time.Time
		RecordedAt        *time.Time
		ActionBy          *UserGroupMini
	}
	tests := []struct {
		name   string
		fields *fields
	}{
		{"normal",
			&fields{
				Type:              "event",
				EventID:           "eventId",
				CreatedBy:         &UserGroupMini{Type: setUserType(TYPE_USER), ID: setStringPtr("111"), Name: setStringPtr("name"), Login: setStringPtr("aaa@example.com")},
				EventType:         EE_COPY,
				SessionID:         setStringPtr("SESSIONID"),
				SourceRaw:         []byte(fileJson001),
				AdditionalDetails: map[string]interface{}{"additional": "details"},
				CreatedAt:         time.Now(),
				RecordedAt:        nil,
				ActionBy:          nil,
			},
		},
		{"normal sessionId nil",
			&fields{
				Type:              "event",
				EventID:           "eventId",
				CreatedBy:         &UserGroupMini{Type: setUserType(TYPE_USER), ID: setStringPtr("111"), Name: setStringPtr("name"), Login: setStringPtr("aaa@example.com")},
				EventType:         EE_COPY,
				SessionID:         nil,
				SourceRaw:         []byte(fileJson001),
				AdditionalDetails: map[string]interface{}{"additional": "details"},
				CreatedAt:         time.Now(),
				RecordedAt:        nil,
				ActionBy:          nil,
			},
		},
		{"nil",
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var be *BoxEvent
			if tt.fields == nil {
				be = nil
			} else {
				be = &BoxEvent{
					Type:              tt.fields.Type,
					EventID:           tt.fields.EventID,
					CreatedBy:         tt.fields.CreatedBy,
					EventType:         tt.fields.EventType,
					SessionID:         tt.fields.SessionID,
					SourceRaw:         tt.fields.SourceRaw,
					source:            tt.fields.source,
					AdditionalDetails: tt.fields.AdditionalDetails,
					CreatedAt:         tt.fields.CreatedAt,
					RecordedAt:        tt.fields.RecordedAt,
					ActionBy:          tt.fields.ActionBy,
				}
			}
			if got := be.String(); got == "" {
				t.Errorf("BoxEvent.String() = \"\"")
			}
		})
	}
}
