package application

import (
	"errors"
	"fmt"
)

// NavigationRequest describes a navigation that a window's OnBeforeNavigate
// is being asked to allow or refuse.
type NavigationRequest struct {
	// URL is the absolute URL the frame is about to load.
	URL string

	// IsMainFrame reports whether the navigation replaces the window's
	// top-level document rather than the content of an iframe. A request to
	// open a new window (target="_blank", window.open) is a new top-level
	// document too, and is reported as the main frame. WebKitGTK does not say
	// which frame a navigation targets, so on Linux every navigation is
	// reported as the main frame.
	IsMainFrame bool
}

// shouldNavigate applies the window's OnBeforeNavigate to request. A window
// without one allows everything, as the webview does by default.
func (w *WebviewWindow) shouldNavigate(request NavigationRequest) bool {
	policy := w.options.OnBeforeNavigate
	if policy == nil {
		return true
	}
	return policy(w, request)
}

// interceptsNavigation reports whether the window has an OnBeforeNavigate,
// so the platform layer can leave navigations to the webview's default
// handling when it does not.
func (w *WebviewWindow) interceptsNavigation() bool {
	return w.options.OnBeforeNavigate != nil
}

// allowNavigation reports whether request may go ahead in window. The
// platform layer calls it on the main thread while the webview holds the
// navigation. A policy that panics refuses: the panic is passed to the
// application's PanicHandler and the zero result stands.
func allowNavigation(window *WebviewWindow, request NavigationRequest) (allow bool) {
	defer handlePanic()
	return window.shouldNavigate(request)
}

// newWindowRequest is the part of WebView2's NewWindowRequested event args
// that refuseNewWindow uses.
type newWindowRequest interface {
	GetUri() (string, error)
	PutHandled(bool) error
}

// refuseNewWindow reports a new-window request to the window's policy as a
// main-frame navigation, then always marks it handled without a window: a
// default WebView2 popup would not carry the policy.
func refuseNewWindow(window *WebviewWindow, request newWindowRequest) error {
	uri, uriErr := request.GetUri()
	if uriErr == nil {
		allowNavigation(window, NavigationRequest{URL: uri, IsMainFrame: true})
	} else {
		uriErr = fmt.Errorf("reading new window URL: %w", uriErr)
	}
	if err := request.PutHandled(true); err != nil {
		return errors.Join(uriErr, fmt.Errorf("refusing new window: %w", err))
	}
	return uriErr
}

// windowInterceptsNavigationByID is interceptsNavigation for platform
// callbacks that only carry a window ID. A window that has gone away is
// intercepted, so decideNavigation refuses it.
func windowInterceptsNavigationByID(windowID uint) bool {
	window := navigationWindow(windowID)
	return window == nil || window.interceptsNavigation()
}

// decideNavigation is allowNavigation for platform callbacks that only carry
// a window ID. A window that has gone away refuses.
func decideNavigation(windowID uint, request NavigationRequest) bool {
	window := navigationWindow(windowID)
	if window == nil {
		return false
	}
	return allowNavigation(window, request)
}

func navigationWindow(windowID uint) *WebviewWindow {
	if globalApplication == nil {
		return nil
	}
	window, ok := globalApplication.Window.GetByID(windowID)
	if !ok {
		return nil
	}
	webviewWindow, _ := window.(*WebviewWindow)
	return webviewWindow
}
