//go:build !((linux || darwin) && purego) && !ios && !android

package application

// navigationPolicySupported reports whether this backend calls OnBeforeNavigate.
const navigationPolicySupported = true
