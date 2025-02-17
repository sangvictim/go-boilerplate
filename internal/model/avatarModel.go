package model

type Avatar struct {
	ID           int64  `json:"id,omitempty"`
	UploadedBY   string `json:"uploaded_by,omitempty"`
	OriginalName string `json:"original_name,omitempty"`
	Key          string `json:"key,omitempty"`
	Size         int64  `json:"size,omitempty"`
	MimeType     string `json:"mime_type,omitempty"`
	Visibility   string `json:"visibility,omitempty"`
}

type AvatarRequest struct {
	UploadedBY   string `json:"uploaded_by,omitempty" validate:"required"`
	OriginalName string `json:"original_name,omitempty" validate:"required"`
	Key          string `json:"key,omitempty" validate:"required"`
	Size         int64  `json:"size,omitempty" validate:"required"`
	MimeType     string `json:"mime_type,omitempty" validate:"required"`
	Visibility   string `json:"visibility,omitempty" validate:"required"`
}

type AvatarResponse struct {
	UploadedBY   string `json:"uploaded_by,omitempty"`
	OriginalName string `json:"original_name,omitempty"`
	Key          string `json:"key,omitempty"`
	Url          string `json:"url,omitempty"`
	Size         int64  `json:"size,omitempty"`
	MimeType     string `json:"mime_type,omitempty"`
	Visibility   string `json:"visibility,omitempty"`
}
