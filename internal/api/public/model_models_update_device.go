/*
Nexodus API

This is the Nexodus API Server.

API version: 1.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package public

// ModelsUpdateDevice struct for ModelsUpdateDevice
type ModelsUpdateDevice struct {
	AdvertiseCidrs  []string         `json:"advertise_cidrs,omitempty"`
	Endpoints       []ModelsEndpoint `json:"endpoints,omitempty"`
	Hostname        string           `json:"hostname,omitempty"`
	Revision        int32            `json:"revision,omitempty"`
	SecurityGroupId string           `json:"security_group_id,omitempty"`
	SymmetricNat    bool             `json:"symmetric_nat,omitempty"`
	VpcId           string           `json:"vpc_id,omitempty"`
}
