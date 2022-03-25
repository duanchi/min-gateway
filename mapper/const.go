package mapper

type Constant struct {
	URL_TYPE []string
	LOADBALANCE_TYPE []string
	REQUEST_METHOD []string
	INTEGRATION_PROTOCOL_TYPE []string
}

var CONSTANT = Constant{
	URL_TYPE: []string {"path", "fnmatch", "regex", "service"},
	LOADBALANCE_TYPE: []string {"polling", "weighted_polling", "random", "ip_hash", "full_request"},
	REQUEST_METHOD: []string{"ALL", "GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS", "CONNECT", "WEBSOCKET"},
	INTEGRATION_PROTOCOL_TYPE: []string{"TRANSFER_PROTOCOL", "DATA_PROTOCOL"},
}
