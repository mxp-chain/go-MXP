package goMXP

import (
	"encoding/json"

	"github.com/pkg/errors"
)

/*
UserActivatedProtocolOverrides represents user activated protocl overrides on the MXP network.

RPC:
	/config/network/user_activated_protocol_overrides (GET)

Link:
	https://MXP.gitlab.io/api/rpc.html#get-config-network-user-activated-protocol-overrides
*/
type UserActivatedProtocolOverrides struct {
	ReplacedProtocol    string `json:"replaced_protocol"`
	ReplacementProtocol string `json:"replacement_protocol"`
}

/*
UserActivatedProtocolOverrides list of protocols which replace other protocols.

Path:
	/config/network/user_activated_protocol_overrides (GET)

Link:
	https://MXP.gitlab.io/api/rpc.html#get-config-network-user-activated-protocol-overrides
*/
func (t *GoMXP) UserActivatedProtocolOverrides() (UserActivatedProtocolOverrides, error) {
	resp, err := t.get("/config/network/user_activated_protocol_overrides")
	if err != nil {
		return UserActivatedProtocolOverrides{}, errors.Wrap(err, "failed to get blocks")
	}

	var userActivatedProtocolOverride UserActivatedProtocolOverrides
	err = json.Unmarshal(resp, &userActivatedProtocolOverride)
	if err != nil {
		return userActivatedProtocolOverride, errors.Wrap(err, "failed to unmarshal blocks")
	}

	return userActivatedProtocolOverride, nil
}
