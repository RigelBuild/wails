//go:build purego || ios || android

package application

// navigationPolicySupported is false where the backend has no navigation hook:
// iOS, Android and the purego backends.
const navigationPolicySupported = false
