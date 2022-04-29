package mapper

type Constant struct {
	URL_TYPE                          []string
	LOADBALANCE_TYPE                  []string
	REQUEST_METHOD                    []string
	INTEGRATION_PROTOCOL_TYPE         []string
	IS_AUTHORIZE                      []bool
	IS_CUSTOM_TOKEN                   []bool
	IS_EPHEMERAL                      []bool
	BOOLEAN_TYPE                      []bool
	IS_EPHEMERAL_REVERSE              map[bool]int
	URL_TYPE_REVERSE                  map[string]int
	LOADBALANCE_TYPE_REVERSE          map[string]int
	INTEGRATION_PROTOCOL_TYPE_REVERSE map[string]int
	IS_AUTHORIZE_REVERSE              map[bool]int
	IS_CUSTOM_TOKEN_REVERSE           map[bool]int
	BOOLEAN_TYPE_REVERSE              map[bool]int
}

var CONSTANT = Constant{
	URL_TYPE: []string{"path", "fnmatch", "regex", "service"},
	URL_TYPE_REVERSE: map[string]int{
		"path":    0,
		"fnmatch": 1,
		"regex":   2,
		"service": 3,
	},
	LOADBALANCE_TYPE: []string{"polling", "weighted_polling", "random", "ip_hash", "full_request"},
	LOADBALANCE_TYPE_REVERSE: map[string]int{
		"polling":          0,
		"weighted_polling": 1,
		"random":           2,
		"ip_hash":          3,
		"full_request":     4,
	},
	REQUEST_METHOD:            []string{"ALL", "GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS", "CONNECT", "WEBSOCKET"},
	INTEGRATION_PROTOCOL_TYPE: []string{"TRANSFER_PROTOCOL", "DATA_PROTOCOL"},
	INTEGRATION_PROTOCOL_TYPE_REVERSE: map[string]int{
		"TRANSFER_PROTOCOL": 0,
		"DATA_PROTOCOL":     1,
	},
	IS_AUTHORIZE: []bool{false, true},
	IS_AUTHORIZE_REVERSE: map[bool]int{
		false: 0,
		true:  1,
	},
	IS_CUSTOM_TOKEN: []bool{false, true},
	IS_CUSTOM_TOKEN_REVERSE: map[bool]int{
		false: 0,
		true:  1,
	},
	IS_EPHEMERAL: []bool{false, true},
	IS_EPHEMERAL_REVERSE: map[bool]int{
		false: 0,
		true:  1,
	},
	BOOLEAN_TYPE: []bool{false, true},
	BOOLEAN_TYPE_REVERSE: map[bool]int{
		false: 0,
		true:  1,
	},
}
