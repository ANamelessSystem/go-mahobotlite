package botlogic

type GroupMessageOut struct {
	GroupID    int                      `json:"group_id"`
	Message    []map[string]interface{} `json:"message"`
	AutoEscape bool                     `json:"auto_escape"`
}

type Stat struct {
	PacketReceived  int `json:"packet_received"`
	PacketSent      int `json:"packet_sent"`
	PacketLost      int `json:"packet_lost"`
	MessageReceived int `json:"message_received"`
	MessageSent     int `json:"message_sent"`
	DisconnectTimes int `json:"disconnect_times"`
	LostTimes       int `json:"lost_times"`
	LastMessageTime int `json:"last_message_time"`
}

type Status struct {
	AppEnabled     bool        `json:"app_enabled"`
	AppGood        bool        `json:"app_good"`
	AppInitialized bool        `json:"app_initialized"`
	Good           bool        `json:"good"`
	Online         bool        `json:"online"`
	PluginsGood    interface{} `json:"plugins_good"`
	Stat           Stat        `json:"stat"`
}

type MetaEvent struct {
	PostType      string `json:"post_type"`
	MetaEventType string `json:"meta_event_type"`
	Time          int64  `json:"time"`
	SelfID        int    `json:"self_id"`
	Status        Status `json:"status"`
	Interval      int    `json:"interval"`
}

type Sender struct {
	Age      int    `json:"age"`
	Area     string `json:"area"`
	Card     string `json:"card"`
	Level    string `json:"level"`
	Nickname string `json:"nickname"`
	Role     string `json:"role"`
	Sex      string `json:"sex"`
	Title    string `json:"title"`
	UserID   int    `json:"user_id"`
}

type GroupMessageIn struct {
	PostType    string      `json:"post_type"`
	MessageType string      `json:"message_type"`
	Time        int64       `json:"time"`
	SelfID      int         `json:"self_id"`
	SubType     string      `json:"sub_type"`
	Anonymous   interface{} `json:"anonymous"`
	GroupID     int         `json:"group_id"`
	Message     string      `json:"message"`
	MessageID   int64       `json:"message_id"`
	Font        int64       `json:"font"`
	MessageSeq  int64       `json:"message_seq"`
	RawMessage  string      `json:"raw_message"`
	UserID      int         `json:"user_id"`
	Sender      Sender      `json:"sender"`
}

type LoginInfoResponse struct {
	Data struct {
		Nickname string `json:"nickname"`
		UserID   int    `json:"user_id"`
	} `json:"data"`
}

var HeartbeatReceived = make(chan struct{}, 1)

type botCmdParams struct {
	IsAll        bool
	EventID      int // rec_id in database
	UserID       int
	AttackBoss   int
	AttackRound  int
	AttackDamage int
	AttackType   int
	AttackLost   bool
}

type CQServerResponse struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Msg     string      `json:"msg"`
	RetCode int         `json:"retcode"`
	Status  string      `json:"status"`
	Wording string      `json:"wording"`
}
