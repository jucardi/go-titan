package logx

import "strings"

const (
	vendorLookup = "/vendor/"
)

var (
	CallerProcessConfig = &CallerCleanupConfig{
		RemoveFileInfo:       true,
		RemoveVendorPrefix:   true,
		RemovePrefixSegments: 2,
	}
	callerProcessor = callerHandler{}
)

type CallerCleanupConfig struct {
	RemovePrefixSegments int
	RemoveVendorPrefix   bool
	RemoveFileInfo       bool
}

type callerHandler struct {
}

func (p *callerHandler) Cleanup(caller string) string {
	ret := caller
	ret = p.vendor(ret)
	ret = p.prefix(ret)
	ret = p.fileInfo(ret)
	return ret
}

func (p *callerHandler) vendor(caller string) string {
	if !CallerProcessConfig.RemoveVendorPrefix || !strings.Contains(caller, vendorLookup) {
		return caller
	}
	split := strings.Split(caller, vendorLookup)
	return split[1]
}

func (p *callerHandler) prefix(caller string) string {
	split := strings.Split(caller, "/")
	if len(split) <= CallerProcessConfig.RemovePrefixSegments {
		return caller
	}
	return strings.Join(split[2:], "/")
}

func (p *callerHandler) fileInfo(caller string) string {
	if !CallerProcessConfig.RemoveFileInfo {
		return caller
	}
	split := strings.Split(caller, "/")
	if len(split) <= 1 {
		return caller
	}
	return strings.Join(split[:len(split)-1], "/")
}
