package slide

import (
	"context"
	"net/http"

	"golang.org/x/oauth2"
)

type serviceConfig struct {
	roundtripper http.RoundTripper
	apiURL       string
}

type configOption func(s *serviceConfig)

func WithCustomRoundtripper(roundtripper http.RoundTripper) configOption {
	return func(s *serviceConfig) {
		s.roundtripper = roundtripper
	}
}

func AddRequestPreProcessor(roundtripper http.RoundTripper) configOption {
	return func(s *serviceConfig) {
		s.roundtripper = roundtripper
	}
}

func WithSlogger(roundtripper http.RoundTripper) configOption {
	return func(s *serviceConfig) {
		s.roundtripper = roundtripper
	}
}

type Service struct {
	accounts               AccountService
	agents                 AgentService
	alerts                 AlertService
	backups                BackupService
	clients                ClientService
	devices                DeviceService
	fileRestores           FileRestoreService
	health                 HealthCheck
	imageExportRestores    ImageExportRestoreService
	networks               NetworkService
	snapshots              SnapshotService
	users                  UserService
	virtualMachineRestores VirtualMachineRestoreService
}

func NewService(
	apiToken string,
	options ...configOption,
) Service {
	config := &serviceConfig{
		apiURL: "api.slide.tech",
	}

	for _, option := range options {
		option(config)
	}

	oauthToken := &oauth2.Token{
		AccessToken: apiToken,
	}

	requestClient := &requestClient{
		token:  oauthToken,
		apiURL: config.apiURL,
		httpClient: &http.Client{
			Transport: config.roundtripper,
		},
	}

	return Service{
		accounts: AccountService{
			baseEndpoint:  "/v1/account",
			requestClient: requestClient,
		},
		agents: AgentService{
			baseEndpoint:  "/v1/agent",
			requestClient: requestClient,
		},
		alerts: AlertService{
			baseEndpoint:  "/v1/alert",
			requestClient: requestClient,
		},
		backups: BackupService{
			baseEndpoint:  "/v1/backup",
			requestClient: requestClient,
		},
		clients: ClientService{
			baseEndpoint:  "/v1/client",
			requestClient: requestClient,
		},
		devices: DeviceService{
			baseEndpoint:  "/v1/device",
			requestClient: requestClient,
		},
		fileRestores: FileRestoreService{
			baseEndpoint:  "/v1/restore/file",
			requestClient: requestClient,
		},
		health: HealthCheck{
			requestClient: requestClient,
		},
		imageExportRestores: ImageExportRestoreService{
			baseEndpoint:  "/v1/restore/image",
			requestClient: requestClient,
		},
		snapshots: SnapshotService{
			baseEndpoint:  "/v1/snapshot",
			requestClient: requestClient,
		},
		users: UserService{
			baseEndpoint:  "/v1/user",
			requestClient: requestClient,
		},
		virtualMachineRestores: VirtualMachineRestoreService{
			baseEndpoint:  "/v1/restore/virt",
			requestClient: requestClient,
		},
	}
}

// https://docs.slide.tech/api/#tag/accounts
func (s Service) Accounts() AccountService {
	return s.accounts
}

// https://docs.slide.tech/api/#tag/agents
func (s Service) Agents() AgentService {
	return s.agents
}

// https://docs.slide.tech/api/#tag/alerts
func (s Service) Alerts() AlertService {
	return s.alerts
}

// https://docs.slide.tech/api/#tag/backups
func (s Service) Backups() BackupService {
	return s.backups
}

// https://docs.slide.tech/api/#tag/clients
func (s Service) Clients() ClientService {
	return s.clients
}

// https://docs.slide.tech/api/#tag/devices
func (s Service) Devices() DeviceService {
	return s.devices
}

// https://docs.slide.tech/api/#tag/restores-file
func (s Service) FileRestores() FileRestoreService {
	return s.fileRestores
}

// https://docs.slide.tech/api/#tag/restores-image
func (s Service) ImageExportRestores() ImageExportRestoreService {
	return s.imageExportRestores
}

// https://docs.slide.tech/api/#tag/networks
func (s Service) Networks() NetworkService {
	return s.networks
}

// https://docs.slide.tech/api/#tag/snapshots
func (s Service) Snapshots() SnapshotService {
	return s.snapshots
}

// https://docs.slide.tech/api/#tag/users
func (s Service) Users() UserService {
	return s.users
}

// https://docs.slide.tech/api/#tag/restores-virtual-machine
func (s Service) VirtualMachineRestores() VirtualMachineRestoreService {
	return s.virtualMachineRestores
}

func (s Service) CheckAuthenticationToken(ctx context.Context) (bool, error) {
	return s.health.IsAuthenticated(ctx)
}
