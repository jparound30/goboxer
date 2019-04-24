package goboxer

import (
	"reflect"
	"testing"
	"time"
)

func TestCollaborationStatus_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		cs      *CollaborationStatus
		want    []byte
		wantErr bool
	}{
		{"nil", nil, []byte(`null`), false},
		{"accepted", setCollaborationStatus(COLLABORATION_STATUS_ACCEPTED), []byte(`"accepted"`), false},
		{"pending", setCollaborationStatus(COLLABORATION_STATUS_PENDING), []byte(`"pending"`), false},
		{"rejected", setCollaborationStatus(COLLABORATION_STATUS_REJECTED), []byte(`"rejected"`), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cs.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("CollaborationStatus.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CollaborationStatus.MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCollaborationStatus_String(t *testing.T) {
	tests := []struct {
		name string
		cs   *CollaborationStatus
		want string
	}{
		{"nil", nil, "<nil>"},
		{"accepted", setCollaborationStatus(COLLABORATION_STATUS_ACCEPTED), "accepted"},
		{"pending", setCollaborationStatus(COLLABORATION_STATUS_PENDING), "pending"},
		{"rejected", setCollaborationStatus(COLLABORATION_STATUS_REJECTED), "rejected"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cs.String(); got != tt.want {
				t.Errorf("CollaborationStatus.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCollaboration_String(t *testing.T) {
	type fields struct {
		apiInfo        *apiInfo
		Type           *string
		ID             *string
		CreatedBy      *UserGroupMini
		CreatedAt      *time.Time
		ModifiedAt     *time.Time
		ExpiresAt      *time.Time
		Status         *CollaborationStatus
		AccessibleBy   *UserGroupMini
		InviteEmail    *string
		Role           *Role
		AcknowledgedAt *time.Time
		Item           *ItemMini
		CanViewPath    *bool
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"nil",
			fields{
				Type: nil,
			},
			"<nil>",
		},
		{"normal",
			fields{
				Type: setStringPtr("collaboration"),
				ID:   setStringPtr("791293"),
				CreatedBy: &UserGroupMini{
					Type:  setUserType(TYPE_USER),
					ID:    setStringPtr("17738362"),
					Name:  setStringPtr("sean rose"),
					Login: setStringPtr("sean@box.com"),
				},
				CreatedAt:  setTime("2012-12-12T10:54:37-08:00"),
				ModifiedAt: setTime("2012-12-12T11:30:43-08:00"),
				ExpiresAt:  nil,
				Status:     setCollaborationStatus(COLLABORATION_STATUS_ACCEPTED),
				AccessibleBy: &UserGroupMini{
					Type:  setUserType(TYPE_USER),
					ID:    setStringPtr("18203124"),
					Name:  setStringPtr("sean"),
					Login: setStringPtr("sean+test@box.com"),
				},
				Role:           setRole(EDITOR),
				AcknowledgedAt: setTime("2012-12-12T11:30:43-08:00"),
				Item: &ItemMini{
					Type:       setItemTypePtr(TYPE_FOLDER),
					ID:         setStringPtr("11446500"),
					SequenceId: setStringPtr("0"),
					ETag:       setStringPtr("0"),
					Name:       setStringPtr("Shared Pictures"),
				},
			},
			"{ Type:collaboration, ID:791293, CreatedBy:{ Type:user, ID:17738362, Name:sean rose, Login:sean@box.com }, " +
				"CreatedAt:2012-12-12T10:54:37-08:00, Modified:2012-12-12T11:30:43-08:00, ExpiresAt:<nil>, Status:accepted," +
				" AccessibleBy:{ Type:user, ID:18203124, Name:sean, Login:sean+test@box.com }, InviteEmail:<nil>, Role:editor, " +
				"AcknowledgedAt:2012-12-12T11:30:43-08:00, Item:{ Type:folder, ID:11446500, Name:Shared Pictures, SequenceId:0, ETag:0 }, " +
				"CanViewPath:<nil> }",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var c *Collaboration
			if tt.fields.Type == nil {
				c = nil
			} else {
				c = &Collaboration{
					apiInfo:        tt.fields.apiInfo,
					Type:           tt.fields.Type,
					ID:             tt.fields.ID,
					CreatedBy:      tt.fields.CreatedBy,
					CreatedAt:      tt.fields.CreatedAt,
					ModifiedAt:     tt.fields.ModifiedAt,
					ExpiresAt:      tt.fields.ExpiresAt,
					Status:         tt.fields.Status,
					AccessibleBy:   tt.fields.AccessibleBy,
					InviteEmail:    tt.fields.InviteEmail,
					Role:           tt.fields.Role,
					AcknowledgedAt: tt.fields.AcknowledgedAt,
					Item:           tt.fields.Item,
					CanViewPath:    tt.fields.CanViewPath,
				}
			}
			if got := c.String(); got != tt.want {
				t.Errorf("Collaboration.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCollaboration_ResourceType(t *testing.T) {
	type fields struct {
		apiInfo        *apiInfo
		Type           *string
		ID             *string
		CreatedBy      *UserGroupMini
		CreatedAt      *time.Time
		ModifiedAt     *time.Time
		ExpiresAt      *time.Time
		Status         *CollaborationStatus
		AccessibleBy   *UserGroupMini
		InviteEmail    *string
		Role           *Role
		AcknowledgedAt *time.Time
		Item           *ItemMini
		CanViewPath    *bool
	}
	tests := []struct {
		name   string
		fields fields
		want   BoxResourceType
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Collaboration{
				apiInfo:        tt.fields.apiInfo,
				Type:           tt.fields.Type,
				ID:             tt.fields.ID,
				CreatedBy:      tt.fields.CreatedBy,
				CreatedAt:      tt.fields.CreatedAt,
				ModifiedAt:     tt.fields.ModifiedAt,
				ExpiresAt:      tt.fields.ExpiresAt,
				Status:         tt.fields.Status,
				AccessibleBy:   tt.fields.AccessibleBy,
				InviteEmail:    tt.fields.InviteEmail,
				Role:           tt.fields.Role,
				AcknowledgedAt: tt.fields.AcknowledgedAt,
				Item:           tt.fields.Item,
				CanViewPath:    tt.fields.CanViewPath,
			}
			if got := c.ResourceType(); got != tt.want {
				t.Errorf("Collaboration.ResourceType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewCollaboration(t *testing.T) {
	type args struct {
		api *ApiConn
	}
	tests := []struct {
		name string
		args args
		want *Collaboration
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCollaboration(tt.args.api); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCollaboration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCollaboration_GetInfoReq(t *testing.T) {
	type fields struct {
		apiInfo        *apiInfo
		Type           *string
		ID             *string
		CreatedBy      *UserGroupMini
		CreatedAt      *time.Time
		ModifiedAt     *time.Time
		ExpiresAt      *time.Time
		Status         *CollaborationStatus
		AccessibleBy   *UserGroupMini
		InviteEmail    *string
		Role           *Role
		AcknowledgedAt *time.Time
		Item           *ItemMini
		CanViewPath    *bool
	}
	type args struct {
		collaborationId string
		fields          []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Request
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Collaboration{
				apiInfo:        tt.fields.apiInfo,
				Type:           tt.fields.Type,
				ID:             tt.fields.ID,
				CreatedBy:      tt.fields.CreatedBy,
				CreatedAt:      tt.fields.CreatedAt,
				ModifiedAt:     tt.fields.ModifiedAt,
				ExpiresAt:      tt.fields.ExpiresAt,
				Status:         tt.fields.Status,
				AccessibleBy:   tt.fields.AccessibleBy,
				InviteEmail:    tt.fields.InviteEmail,
				Role:           tt.fields.Role,
				AcknowledgedAt: tt.fields.AcknowledgedAt,
				Item:           tt.fields.Item,
				CanViewPath:    tt.fields.CanViewPath,
			}
			if got := c.GetInfoReq(tt.args.collaborationId, tt.args.fields); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Collaboration.GetInfoReq() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCollaboration_GetInfo(t *testing.T) {
	type fields struct {
		apiInfo        *apiInfo
		Type           *string
		ID             *string
		CreatedBy      *UserGroupMini
		CreatedAt      *time.Time
		ModifiedAt     *time.Time
		ExpiresAt      *time.Time
		Status         *CollaborationStatus
		AccessibleBy   *UserGroupMini
		InviteEmail    *string
		Role           *Role
		AcknowledgedAt *time.Time
		Item           *ItemMini
		CanViewPath    *bool
	}
	type args struct {
		collaborationId string
		fields          []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Collaboration
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Collaboration{
				apiInfo:        tt.fields.apiInfo,
				Type:           tt.fields.Type,
				ID:             tt.fields.ID,
				CreatedBy:      tt.fields.CreatedBy,
				CreatedAt:      tt.fields.CreatedAt,
				ModifiedAt:     tt.fields.ModifiedAt,
				ExpiresAt:      tt.fields.ExpiresAt,
				Status:         tt.fields.Status,
				AccessibleBy:   tt.fields.AccessibleBy,
				InviteEmail:    tt.fields.InviteEmail,
				Role:           tt.fields.Role,
				AcknowledgedAt: tt.fields.AcknowledgedAt,
				Item:           tt.fields.Item,
				CanViewPath:    tt.fields.CanViewPath,
			}
			got, err := c.GetInfo(tt.args.collaborationId, tt.args.fields)
			if (err != nil) != tt.wantErr {
				t.Errorf("Collaboration.GetInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Collaboration.GetInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCollaboration_SetItem(t *testing.T) {
	type fields struct {
		apiInfo        *apiInfo
		Type           *string
		ID             *string
		CreatedBy      *UserGroupMini
		CreatedAt      *time.Time
		ModifiedAt     *time.Time
		ExpiresAt      *time.Time
		Status         *CollaborationStatus
		AccessibleBy   *UserGroupMini
		InviteEmail    *string
		Role           *Role
		AcknowledgedAt *time.Time
		Item           *ItemMini
		CanViewPath    *bool
	}
	type args struct {
		typ ItemType
		id  string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Collaboration
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Collaboration{
				apiInfo:        tt.fields.apiInfo,
				Type:           tt.fields.Type,
				ID:             tt.fields.ID,
				CreatedBy:      tt.fields.CreatedBy,
				CreatedAt:      tt.fields.CreatedAt,
				ModifiedAt:     tt.fields.ModifiedAt,
				ExpiresAt:      tt.fields.ExpiresAt,
				Status:         tt.fields.Status,
				AccessibleBy:   tt.fields.AccessibleBy,
				InviteEmail:    tt.fields.InviteEmail,
				Role:           tt.fields.Role,
				AcknowledgedAt: tt.fields.AcknowledgedAt,
				Item:           tt.fields.Item,
				CanViewPath:    tt.fields.CanViewPath,
			}
			if got := c.SetItem(tt.args.typ, tt.args.id); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Collaboration.SetItem() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCollaboration_SetAccessibleById(t *testing.T) {
	type fields struct {
		apiInfo        *apiInfo
		Type           *string
		ID             *string
		CreatedBy      *UserGroupMini
		CreatedAt      *time.Time
		ModifiedAt     *time.Time
		ExpiresAt      *time.Time
		Status         *CollaborationStatus
		AccessibleBy   *UserGroupMini
		InviteEmail    *string
		Role           *Role
		AcknowledgedAt *time.Time
		Item           *ItemMini
		CanViewPath    *bool
	}
	type args struct {
		typ UserGroupType
		id  string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Collaboration
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Collaboration{
				apiInfo:        tt.fields.apiInfo,
				Type:           tt.fields.Type,
				ID:             tt.fields.ID,
				CreatedBy:      tt.fields.CreatedBy,
				CreatedAt:      tt.fields.CreatedAt,
				ModifiedAt:     tt.fields.ModifiedAt,
				ExpiresAt:      tt.fields.ExpiresAt,
				Status:         tt.fields.Status,
				AccessibleBy:   tt.fields.AccessibleBy,
				InviteEmail:    tt.fields.InviteEmail,
				Role:           tt.fields.Role,
				AcknowledgedAt: tt.fields.AcknowledgedAt,
				Item:           tt.fields.Item,
				CanViewPath:    tt.fields.CanViewPath,
			}
			if got := c.SetAccessibleById(tt.args.typ, tt.args.id); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Collaboration.SetAccessibleById() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCollaboration_SetAccessibleByEmailForNewUser(t *testing.T) {
	type fields struct {
		apiInfo        *apiInfo
		Type           *string
		ID             *string
		CreatedBy      *UserGroupMini
		CreatedAt      *time.Time
		ModifiedAt     *time.Time
		ExpiresAt      *time.Time
		Status         *CollaborationStatus
		AccessibleBy   *UserGroupMini
		InviteEmail    *string
		Role           *Role
		AcknowledgedAt *time.Time
		Item           *ItemMini
		CanViewPath    *bool
	}
	type args struct {
		typ   UserGroupType
		login string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Collaboration
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Collaboration{
				apiInfo:        tt.fields.apiInfo,
				Type:           tt.fields.Type,
				ID:             tt.fields.ID,
				CreatedBy:      tt.fields.CreatedBy,
				CreatedAt:      tt.fields.CreatedAt,
				ModifiedAt:     tt.fields.ModifiedAt,
				ExpiresAt:      tt.fields.ExpiresAt,
				Status:         tt.fields.Status,
				AccessibleBy:   tt.fields.AccessibleBy,
				InviteEmail:    tt.fields.InviteEmail,
				Role:           tt.fields.Role,
				AcknowledgedAt: tt.fields.AcknowledgedAt,
				Item:           tt.fields.Item,
				CanViewPath:    tt.fields.CanViewPath,
			}
			if got := c.SetAccessibleByEmailForNewUser(tt.args.typ, tt.args.login); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Collaboration.SetAccessibleByEmailForNewUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCollaboration_SetRole(t *testing.T) {
	type fields struct {
		apiInfo        *apiInfo
		Type           *string
		ID             *string
		CreatedBy      *UserGroupMini
		CreatedAt      *time.Time
		ModifiedAt     *time.Time
		ExpiresAt      *time.Time
		Status         *CollaborationStatus
		AccessibleBy   *UserGroupMini
		InviteEmail    *string
		Role           *Role
		AcknowledgedAt *time.Time
		Item           *ItemMini
		CanViewPath    *bool
	}
	type args struct {
		role Role
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Collaboration
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Collaboration{
				apiInfo:        tt.fields.apiInfo,
				Type:           tt.fields.Type,
				ID:             tt.fields.ID,
				CreatedBy:      tt.fields.CreatedBy,
				CreatedAt:      tt.fields.CreatedAt,
				ModifiedAt:     tt.fields.ModifiedAt,
				ExpiresAt:      tt.fields.ExpiresAt,
				Status:         tt.fields.Status,
				AccessibleBy:   tt.fields.AccessibleBy,
				InviteEmail:    tt.fields.InviteEmail,
				Role:           tt.fields.Role,
				AcknowledgedAt: tt.fields.AcknowledgedAt,
				Item:           tt.fields.Item,
				CanViewPath:    tt.fields.CanViewPath,
			}
			if got := c.SetRole(tt.args.role); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Collaboration.SetRole() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCollaboration_SetCanViewPath(t *testing.T) {
	type fields struct {
		apiInfo        *apiInfo
		Type           *string
		ID             *string
		CreatedBy      *UserGroupMini
		CreatedAt      *time.Time
		ModifiedAt     *time.Time
		ExpiresAt      *time.Time
		Status         *CollaborationStatus
		AccessibleBy   *UserGroupMini
		InviteEmail    *string
		Role           *Role
		AcknowledgedAt *time.Time
		Item           *ItemMini
		CanViewPath    *bool
	}
	type args struct {
		canViewPath bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Collaboration
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Collaboration{
				apiInfo:        tt.fields.apiInfo,
				Type:           tt.fields.Type,
				ID:             tt.fields.ID,
				CreatedBy:      tt.fields.CreatedBy,
				CreatedAt:      tt.fields.CreatedAt,
				ModifiedAt:     tt.fields.ModifiedAt,
				ExpiresAt:      tt.fields.ExpiresAt,
				Status:         tt.fields.Status,
				AccessibleBy:   tt.fields.AccessibleBy,
				InviteEmail:    tt.fields.InviteEmail,
				Role:           tt.fields.Role,
				AcknowledgedAt: tt.fields.AcknowledgedAt,
				Item:           tt.fields.Item,
				CanViewPath:    tt.fields.CanViewPath,
			}
			if got := c.SetCanViewPath(tt.args.canViewPath); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Collaboration.SetCanViewPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCollaboration_CreateReq(t *testing.T) {
	type fields struct {
		apiInfo        *apiInfo
		Type           *string
		ID             *string
		CreatedBy      *UserGroupMini
		CreatedAt      *time.Time
		ModifiedAt     *time.Time
		ExpiresAt      *time.Time
		Status         *CollaborationStatus
		AccessibleBy   *UserGroupMini
		InviteEmail    *string
		Role           *Role
		AcknowledgedAt *time.Time
		Item           *ItemMini
		CanViewPath    *bool
	}
	type args struct {
		fields []string
		notify bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Request
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Collaboration{
				apiInfo:        tt.fields.apiInfo,
				Type:           tt.fields.Type,
				ID:             tt.fields.ID,
				CreatedBy:      tt.fields.CreatedBy,
				CreatedAt:      tt.fields.CreatedAt,
				ModifiedAt:     tt.fields.ModifiedAt,
				ExpiresAt:      tt.fields.ExpiresAt,
				Status:         tt.fields.Status,
				AccessibleBy:   tt.fields.AccessibleBy,
				InviteEmail:    tt.fields.InviteEmail,
				Role:           tt.fields.Role,
				AcknowledgedAt: tt.fields.AcknowledgedAt,
				Item:           tt.fields.Item,
				CanViewPath:    tt.fields.CanViewPath,
			}
			if got := c.CreateReq(tt.args.fields, tt.args.notify); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Collaboration.CreateReq() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCollaboration_Create(t *testing.T) {
	type fields struct {
		apiInfo        *apiInfo
		Type           *string
		ID             *string
		CreatedBy      *UserGroupMini
		CreatedAt      *time.Time
		ModifiedAt     *time.Time
		ExpiresAt      *time.Time
		Status         *CollaborationStatus
		AccessibleBy   *UserGroupMini
		InviteEmail    *string
		Role           *Role
		AcknowledgedAt *time.Time
		Item           *ItemMini
		CanViewPath    *bool
	}
	type args struct {
		fields []string
		notify bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Collaboration
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Collaboration{
				apiInfo:        tt.fields.apiInfo,
				Type:           tt.fields.Type,
				ID:             tt.fields.ID,
				CreatedBy:      tt.fields.CreatedBy,
				CreatedAt:      tt.fields.CreatedAt,
				ModifiedAt:     tt.fields.ModifiedAt,
				ExpiresAt:      tt.fields.ExpiresAt,
				Status:         tt.fields.Status,
				AccessibleBy:   tt.fields.AccessibleBy,
				InviteEmail:    tt.fields.InviteEmail,
				Role:           tt.fields.Role,
				AcknowledgedAt: tt.fields.AcknowledgedAt,
				Item:           tt.fields.Item,
				CanViewPath:    tt.fields.CanViewPath,
			}
			got, err := c.Create(tt.args.fields, tt.args.notify)
			if (err != nil) != tt.wantErr {
				t.Errorf("Collaboration.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Collaboration.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCollaboration_UpdateReq(t *testing.T) {
	type fields struct {
		apiInfo        *apiInfo
		Type           *string
		ID             *string
		CreatedBy      *UserGroupMini
		CreatedAt      *time.Time
		ModifiedAt     *time.Time
		ExpiresAt      *time.Time
		Status         *CollaborationStatus
		AccessibleBy   *UserGroupMini
		InviteEmail    *string
		Role           *Role
		AcknowledgedAt *time.Time
		Item           *ItemMini
		CanViewPath    *bool
	}
	type args struct {
		collaborationId string
		role            Role
		status          *CollaborationStatus
		canViewPath     *bool
		fields          []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Request
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Collaboration{
				apiInfo:        tt.fields.apiInfo,
				Type:           tt.fields.Type,
				ID:             tt.fields.ID,
				CreatedBy:      tt.fields.CreatedBy,
				CreatedAt:      tt.fields.CreatedAt,
				ModifiedAt:     tt.fields.ModifiedAt,
				ExpiresAt:      tt.fields.ExpiresAt,
				Status:         tt.fields.Status,
				AccessibleBy:   tt.fields.AccessibleBy,
				InviteEmail:    tt.fields.InviteEmail,
				Role:           tt.fields.Role,
				AcknowledgedAt: tt.fields.AcknowledgedAt,
				Item:           tt.fields.Item,
				CanViewPath:    tt.fields.CanViewPath,
			}
			if got := c.UpdateReq(tt.args.collaborationId, tt.args.role, tt.args.status, tt.args.canViewPath, tt.args.fields); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Collaboration.UpdateReq() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCollaboration_Update(t *testing.T) {
	type fields struct {
		apiInfo        *apiInfo
		Type           *string
		ID             *string
		CreatedBy      *UserGroupMini
		CreatedAt      *time.Time
		ModifiedAt     *time.Time
		ExpiresAt      *time.Time
		Status         *CollaborationStatus
		AccessibleBy   *UserGroupMini
		InviteEmail    *string
		Role           *Role
		AcknowledgedAt *time.Time
		Item           *ItemMini
		CanViewPath    *bool
	}
	type args struct {
		collaborationId string
		role            Role
		status          *CollaborationStatus
		canViewPath     *bool
		fields          []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Collaboration
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Collaboration{
				apiInfo:        tt.fields.apiInfo,
				Type:           tt.fields.Type,
				ID:             tt.fields.ID,
				CreatedBy:      tt.fields.CreatedBy,
				CreatedAt:      tt.fields.CreatedAt,
				ModifiedAt:     tt.fields.ModifiedAt,
				ExpiresAt:      tt.fields.ExpiresAt,
				Status:         tt.fields.Status,
				AccessibleBy:   tt.fields.AccessibleBy,
				InviteEmail:    tt.fields.InviteEmail,
				Role:           tt.fields.Role,
				AcknowledgedAt: tt.fields.AcknowledgedAt,
				Item:           tt.fields.Item,
				CanViewPath:    tt.fields.CanViewPath,
			}
			got, err := c.Update(tt.args.collaborationId, tt.args.role, tt.args.status, tt.args.canViewPath, tt.args.fields)
			if (err != nil) != tt.wantErr {
				t.Errorf("Collaboration.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Collaboration.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCollaboration_DeleteReq(t *testing.T) {
	type fields struct {
		apiInfo        *apiInfo
		Type           *string
		ID             *string
		CreatedBy      *UserGroupMini
		CreatedAt      *time.Time
		ModifiedAt     *time.Time
		ExpiresAt      *time.Time
		Status         *CollaborationStatus
		AccessibleBy   *UserGroupMini
		InviteEmail    *string
		Role           *Role
		AcknowledgedAt *time.Time
		Item           *ItemMini
		CanViewPath    *bool
	}
	type args struct {
		collaborationId string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Request
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Collaboration{
				apiInfo:        tt.fields.apiInfo,
				Type:           tt.fields.Type,
				ID:             tt.fields.ID,
				CreatedBy:      tt.fields.CreatedBy,
				CreatedAt:      tt.fields.CreatedAt,
				ModifiedAt:     tt.fields.ModifiedAt,
				ExpiresAt:      tt.fields.ExpiresAt,
				Status:         tt.fields.Status,
				AccessibleBy:   tt.fields.AccessibleBy,
				InviteEmail:    tt.fields.InviteEmail,
				Role:           tt.fields.Role,
				AcknowledgedAt: tt.fields.AcknowledgedAt,
				Item:           tt.fields.Item,
				CanViewPath:    tt.fields.CanViewPath,
			}
			if got := c.DeleteReq(tt.args.collaborationId); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Collaboration.DeleteReq() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCollaboration_Delete(t *testing.T) {
	type fields struct {
		apiInfo        *apiInfo
		Type           *string
		ID             *string
		CreatedBy      *UserGroupMini
		CreatedAt      *time.Time
		ModifiedAt     *time.Time
		ExpiresAt      *time.Time
		Status         *CollaborationStatus
		AccessibleBy   *UserGroupMini
		InviteEmail    *string
		Role           *Role
		AcknowledgedAt *time.Time
		Item           *ItemMini
		CanViewPath    *bool
	}
	type args struct {
		collaborationId string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Collaboration{
				apiInfo:        tt.fields.apiInfo,
				Type:           tt.fields.Type,
				ID:             tt.fields.ID,
				CreatedBy:      tt.fields.CreatedBy,
				CreatedAt:      tt.fields.CreatedAt,
				ModifiedAt:     tt.fields.ModifiedAt,
				ExpiresAt:      tt.fields.ExpiresAt,
				Status:         tt.fields.Status,
				AccessibleBy:   tt.fields.AccessibleBy,
				InviteEmail:    tt.fields.InviteEmail,
				Role:           tt.fields.Role,
				AcknowledgedAt: tt.fields.AcknowledgedAt,
				Item:           tt.fields.Item,
				CanViewPath:    tt.fields.CanViewPath,
			}
			if err := c.Delete(tt.args.collaborationId); (err != nil) != tt.wantErr {
				t.Errorf("Collaboration.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCollaboration_PendingCollaborationsReq(t *testing.T) {
	type fields struct {
		apiInfo        *apiInfo
		Type           *string
		ID             *string
		CreatedBy      *UserGroupMini
		CreatedAt      *time.Time
		ModifiedAt     *time.Time
		ExpiresAt      *time.Time
		Status         *CollaborationStatus
		AccessibleBy   *UserGroupMini
		InviteEmail    *string
		Role           *Role
		AcknowledgedAt *time.Time
		Item           *ItemMini
		CanViewPath    *bool
	}
	type args struct {
		offset int
		limit  int
		fields []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Request
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Collaboration{
				apiInfo:        tt.fields.apiInfo,
				Type:           tt.fields.Type,
				ID:             tt.fields.ID,
				CreatedBy:      tt.fields.CreatedBy,
				CreatedAt:      tt.fields.CreatedAt,
				ModifiedAt:     tt.fields.ModifiedAt,
				ExpiresAt:      tt.fields.ExpiresAt,
				Status:         tt.fields.Status,
				AccessibleBy:   tt.fields.AccessibleBy,
				InviteEmail:    tt.fields.InviteEmail,
				Role:           tt.fields.Role,
				AcknowledgedAt: tt.fields.AcknowledgedAt,
				Item:           tt.fields.Item,
				CanViewPath:    tt.fields.CanViewPath,
			}
			if got := c.PendingCollaborationsReq(tt.args.offset, tt.args.limit, tt.args.fields); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Collaboration.PendingCollaborationsReq() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCollaboration_PendingCollaborations(t *testing.T) {
	type fields struct {
		apiInfo        *apiInfo
		Type           *string
		ID             *string
		CreatedBy      *UserGroupMini
		CreatedAt      *time.Time
		ModifiedAt     *time.Time
		ExpiresAt      *time.Time
		Status         *CollaborationStatus
		AccessibleBy   *UserGroupMini
		InviteEmail    *string
		Role           *Role
		AcknowledgedAt *time.Time
		Item           *ItemMini
		CanViewPath    *bool
	}
	type args struct {
		offset int
		limit  int
		fields []string
	}
	tests := []struct {
		name              string
		fields            fields
		args              args
		wantPendingList   []*Collaboration
		wantOutOffset     int
		wantOutLimit      int
		wantOutTotalCount int
		wantErr           bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Collaboration{
				apiInfo:        tt.fields.apiInfo,
				Type:           tt.fields.Type,
				ID:             tt.fields.ID,
				CreatedBy:      tt.fields.CreatedBy,
				CreatedAt:      tt.fields.CreatedAt,
				ModifiedAt:     tt.fields.ModifiedAt,
				ExpiresAt:      tt.fields.ExpiresAt,
				Status:         tt.fields.Status,
				AccessibleBy:   tt.fields.AccessibleBy,
				InviteEmail:    tt.fields.InviteEmail,
				Role:           tt.fields.Role,
				AcknowledgedAt: tt.fields.AcknowledgedAt,
				Item:           tt.fields.Item,
				CanViewPath:    tt.fields.CanViewPath,
			}
			gotPendingList, gotOutOffset, gotOutLimit, gotOutTotalCount, err := c.PendingCollaborations(tt.args.offset, tt.args.limit, tt.args.fields)
			if (err != nil) != tt.wantErr {
				t.Errorf("Collaboration.PendingCollaborations() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotPendingList, tt.wantPendingList) {
				t.Errorf("Collaboration.PendingCollaborations() gotPendingList = %v, want %v", gotPendingList, tt.wantPendingList)
			}
			if gotOutOffset != tt.wantOutOffset {
				t.Errorf("Collaboration.PendingCollaborations() gotOutOffset = %v, want %v", gotOutOffset, tt.wantOutOffset)
			}
			if gotOutLimit != tt.wantOutLimit {
				t.Errorf("Collaboration.PendingCollaborations() gotOutLimit = %v, want %v", gotOutLimit, tt.wantOutLimit)
			}
			if gotOutTotalCount != tt.wantOutTotalCount {
				t.Errorf("Collaboration.PendingCollaborations() gotOutTotalCount = %v, want %v", gotOutTotalCount, tt.wantOutTotalCount)
			}
		})
	}
}
