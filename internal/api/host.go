package api

type HostApi interface {
	func GetHosts() []*Host
	func GetHost(name string) Host
}
