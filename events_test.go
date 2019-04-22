package goboxer

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
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
		SessionID         string
		SourceRaw         json.RawMessage
		source            BoxResource
		AdditionalDetails map[string]interface{}
		CreatedAt         time.Time
		RecordedAt        time.Time
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
				RecordedAt: *setTime("2012-12-12T10:53:48-08:00"),
				EventType:  UE_ITEM_CREATE,
				SessionID:  "70090280850c8d2a1933c1",
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
