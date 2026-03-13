package types

// Environment keys shared across runner/gateway/plugins.
const (
	EnvNATSURL       = "NATS_URL"
	EnvPluginDataDir = "PLUGIN_DATA_DIR"
	EnvAPIHost       = "API_HOST"
	EnvAPIPort       = "API_PORT"
	EnvPluginRPCSbj  = "PLUGIN_RPC_SUBJECT"
	EnvRuntimeFile   = "RUNTIME_FILE"
)

// Internal RPC/API method names.
const (
	RPCMethodHealthCheck       = "/_internal/health"
	RPCMethodInitialize        = "initialize"
	RPCMethodLoggingGetLevel   = "logging/get_level"
	RPCMethodLoggingSetLevel   = "logging/set_level"
	RPCMethodStorageUpdate     = "storage/update"
	RPCMethodPluginReset       = "plugin/reset"
	RPCMethodDevicesCreate     = "devices/create"
	RPCMethodDevicesUpdate     = "devices/update"
	RPCMethodDevicesDelete     = "devices/delete"
	RPCMethodDevicesList       = "devices/list"
	RPCMethodDevicesRefresh    = "devices/refresh"
	RPCMethodEntitiesCreate    = "entities/create"
	RPCMethodEntitiesUpdate    = "entities/update"
	RPCMethodEntitiesDelete    = "entities/delete"
	RPCMethodEntitiesList      = "entities/list"
	RPCMethodEntitiesRefresh   = "entities/refresh"
	RPCMethodSnapshotsSave     = "entities/snapshots/save"
	RPCMethodSnapshotsList     = "entities/snapshots/list"
	RPCMethodSnapshotsDelete   = "entities/snapshots/delete"
	RPCMethodSnapshotsRestore  = "entities/snapshots/restore"
	RPCMethodCommandsCreate    = "entities/commands/create"
	RPCMethodCommandsStatusGet = "commands/status/get"
	RPCMethodEventsIngest      = "entities/events/ingest"
	RPCMethodStorageFlush      = "storage/flush"
	RPCMethodScriptsGet        = "scripts/get"
	RPCMethodScriptsPut        = "scripts/put"
	RPCMethodScriptsDelete     = "scripts/delete"
	RPCMethodScriptStateGet    = "scripts/state/get"
	RPCMethodScriptStatePut    = "scripts/state/put"
	RPCMethodScriptStateDelete = "scripts/state/delete"
)

// NATS subjects shared across runtime components.
const (
	SubjectRPCPrefix      = "slidebolt.rpc."
	SubjectRegistration   = "slidebolt.registration"
	SubjectDiscoveryProbe = "slidebolt.discovery.probe"
	SubjectSearchPlugins  = "slidebolt.search.plugins"
	SubjectSearchDevices  = "slidebolt.search.devices"
	SubjectSearchEntities = "slidebolt.search.entities"
	SubjectEntityEvents   = "slidebolt.entity.events"
	SubjectCommandStatus  = "slidebolt.command.status"

	SubjectDeviceCreated    = "slidebolt.device.created"
	SubjectDeviceRead       = "slidebolt.device.read"
	SubjectDeviceUpdated    = "slidebolt.device.updated"
	SubjectDeviceDeleted    = "slidebolt.device.deleted"
	SubjectEntityCreated    = "slidebolt.entity.created"
	SubjectEntityRead       = "slidebolt.entity.read"
	SubjectEntityUpdated    = "slidebolt.entity.updated"
	SubjectEntityDeleted    = "slidebolt.entity.deleted"
	SubjectGatewayDiscovery = "slidebolt.gateway.discovery"
)
