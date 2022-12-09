/*
Copyright 2022 Ciena Corporation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package gns3

const (
	TypeNat              = "nat"
	TypeEthernetSwitch   = "eithernet_switch"
	TypeVpcs             = "vpcs"
	TypeRouter           = "router"
	TypeFirewall         = "firewall"
	TypeMultilayerSwitch = "multilayer_switch"

	CategoryMultilayerSwitch = "multilayer_switch"
	CategorySwitch           = "switch"
	CategoryGuest            = "guest"
	CategoryRouter           = "router"
	CategoryFirewall         = "firewall"

	SymbolDockerGuest      = ":/symbols/classic/docker_guest.svg"
	SymbolQemuGuest        = ":/symbols/classic/qemu_guest.svg"
	SymbolRouter           = ":/symbols/classic/router.svg"
	SymbolEthernetSwitch   = ":/symbols/classic/ethernet_switch.svg"
	SymbolMultilayerSwitch = ":/symbols/classic/multilayer_switch.svg"
	SymbolFirewall         = ":/symbols/classic_firewall.svg"
	SymbolCloud            = ":/symbols/classic/cloud.svg"
	SymbolVpcs             = ":/symbols/classic/vpcs_guest.svg"

	TemplateTypeQemu     = "qemu"
	TemplateTypeIou      = "iou"
	TemplateTypeDynamips = "dynamips"
	TemplateTypeDocker   = "docker"
)
