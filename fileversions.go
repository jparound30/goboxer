package goboxer

type FileVersion struct {
	apiInfo *apiInfo
	Type    string `json:"type,omitempty"`
	ID      string `json:"id,omitempty"`
	Sha1    string `json:"sha1,omitempty"`
}

func (fv *FileVersion) ResourceType() BoxResourceType {
	return FileVersionResource
}
