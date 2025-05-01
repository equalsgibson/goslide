package slide

type VirtualMachineRestore struct {
	AgentID      string `json:"agent_id"`
	CPUCount     string `json:"cpu_count"`
	CreatedAt    string `json:"created_at"`
	DeviceID     string `json:"device_id"`
	DiskBus      string `json:"disk_bus"`
	ExpiresAt    string `json:"expires_at"`
	MemoryInMB   string `json:"memory_in_mb"`
	NetworkModel string `json:"network_model"`
	NetworkType  string `json:"network_type"`
	SnapshotID   string `json:"snapshot_id"`
	State        string `json:"state"`
	VirtID       string `json:"virt_id"`
	VNC          string `json:"vnc"`
	VNCPassword  string `json:"vnc_password"`
}

type VirtualMachineRestoreService struct {
	baseEndpoint  string
	requestClient *requestClient
}

type VirtualMachineVNC struct {
	Host         string `json:"host"`
	Port         string `json:"port"`
	Type         string `json:"type"`
	WebsocketURI string `json:"websocket_uri"`
}
