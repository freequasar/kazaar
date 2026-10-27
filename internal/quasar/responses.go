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

// get sync info

type SyncInfoResponse struct {
	Status       string             `json:"status"`
	Config       SyncInfoConfig     `json:"config"`
	Subscription SubscriptionObject `json:"subscription"`
}

type SubscriptionObject struct {
	Mode string `json:"mode"` // should be always "transaction"
}

type SyncInfoConfig struct {
	AccountConfig AccountConfig `json:"account_config"`
	DeviceConfig  DeviceConfig  `json:"device_config"`
	SystemConfig  SystemConfig  `json:"system_config"`
}

type AccountConfig struct {
	// stab
	Jingle  bool   `json:"jingle"`  // signal after user sez "Алиса" or "ياسمينة"
	Spotter string `json:"spotter"` // alisa, yandex or yasmina
}

type DeviceConfig struct {
	Name      string    `json:"name,omitempty"`
	Locale    string    `json:"locale,omitempty"`
	Equalizer Equalizer `json:"equalizer,omitempty"`
}

type SystemConfig struct {
	ADBDaemon ADBDaemon `json:"adbd,omitempty"`
}

// device config structs

type Equalizer struct {
	Enabled        bool            `json:"enabled"`
	ActivePresetID string          `json:"active_preset_id"`
	Bands          []EqualizerBand `json:"bands"`
}

type EqualizerBand struct {
	Frequency int     `json:"freq"`
	Gain      float32 `json:"gain"`
	Width     int     `json:"width"`
}

// system config structs

type DaemonObject struct {
	Enabled bool `json:"enabled"`
}

type ADBDaemon struct {
	DaemonObject
	RemoteActivation bool `json:"remoteActivation"`
	Network          bool `json:"network"`
}
