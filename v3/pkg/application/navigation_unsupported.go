//go:build ((linux || darwin) && purego) || ios || android

package application

// navigationPolicySupported is false where the backend has no navigation hook:
// iOS, Android and the Linux and macOS purego backends.
const navigationPolicySupported = false
