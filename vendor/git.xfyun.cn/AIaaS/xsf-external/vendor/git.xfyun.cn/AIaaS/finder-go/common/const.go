package finder

const DefaultCacheDir = "findercache"

const (
	ConfigEventPrefix          = "config_"
	ServiceEventPrefix         = "service_"
	ServiceConfEventPrefix     = "service_conf_"
	ServiceProviderEventPrefix = "service_provider_"
	ServiceConsumerEventPrefix = "service_consumer_"
)

type InstanceChangedEventType string

const (
	INSTANCEADDED  InstanceChangedEventType = "INSTANCEADDED"
	INSTANCEREMOVE InstanceChangedEventType = "INSTANCEREMOVE"
)

const ZK_NODE_DOSE_NOT_EXIST  ="zk: node does not exist"