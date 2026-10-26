package quasar

type EnvelopeResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

type CheckUpdatesResponse struct {
	HasUpdate   bool   `json:"hasUpdate"`
	Version     string `json:"version,omitempty"`
	DownloadUrl string `json:"downloadUrl,omitempty"`
	Crc32       int    `json:"crc32,omitempty"`
	IsCritical  bool   `json:"critical,omitempty"`
	FallbackUrl string `json:"fallbackUrl,omitempty"`
	FallbackIp  string `json:"fallbackIp,omitempty"`
	UpdateId    string `json:"updateId,omitempty"`
}

type GetTimezoneResponse struct {
	Status          string `json:"status"`
	Timezone        string `json:"timezone"`
	OffsetInSeconds int    `json:"offset_sec"`
}

type RegisterResponse struct {
	Status            string `json:"status"`
	HasCriticalUpdate bool   `json:"hasCriticalUpdate"`
}
