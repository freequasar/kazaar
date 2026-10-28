package quasar

type SyncInfoResponse struct {
	Status       string             `json:"status"`
	Config       SyncInfoConfig     `json:"config"`
	Subscription SubscriptionObject `json:"subscription"`
}

type SubscriptionObject struct {
	Mode string `json:"mode"`
	// for basic speakers it´s just "transaction"
	// but for subscription based it's "subscription"
	// with this additional fields are added
	BrickEnabled bool `json:"enabled,omitempty"`
	TimeToLock   int  `json:"ttl,omitempty"` // NOTE probably unix time, experiments required
}

type SyncInfoConfig struct {
	AccountConfig AccountConfig `json:"account_config"`
	DeviceConfig  DeviceConfig  `json:"device_config"`
	SystemConfig  SystemConfig  `json:"system_config"`
}

type AccountConfig struct {
	// stab
	Jingle      bool        `json:"jingle"`  // signal after user sez "Алиса" or "ياسمينة"
	Spotter     string      `json:"spotter"` // alisa, yandex for ru-RU, yasmina for ar-SA
	CmdSpotters CmdSpotters `json:"enabledCommandSpotters"`
}

type DeviceConfig struct {
	Name      string     `json:"name,omitempty"`
	Locale    string     `json:"locale,omitempty"`
	Equalizer *Equalizer `json:"equalizer"`
}

type SystemConfig struct {
	ADBDaemon ADBDaemon `json:"adbd,omitempty"`
}

// #region account config structs

type CmdSpotters struct {
	Call      CallSpotters      `json:"call"`
	Music     MusicSpotters     `json:"music"`
	SmartHome SmartHomeSpotters `json:"smartHome"`
	TV        TVSpotters        `json:"tv"`
}

type CallSpotters struct {
	Answer bool `json:"answer"`
}

type MusicSpotters struct {
	Bluetooth    bool `json:"bluetooth"`
	Feedback     bool `json:"feedback"`
	Navigation   bool `json:"navigation"`
	PlayAndPause bool `json:"playAndPause"`
	Volume       bool `json:"volume"`
}

type SmartHomeSpotters struct {
	Light bool `json:"light"`
	TV    bool `json:"tv"`
}

type TVSpotters struct { // prob for speakers with hdmi?
	BackToHome bool `json:"backToHome"`
	Navigation bool `json:"navigation"`
}

// #endregion

// #region device config structs

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

// #endregion

// #region system config structs

type DaemonObject struct {
	Enabled bool `json:"enabled"`
}

type ADBDaemon struct {
	DaemonObject
	RemoteActivation bool `json:"remoteActivation"`
	Network          bool `json:"network"`
}

// #endregion
