package com

// Context interface for providing a dynamically enabled components
type Context interface {
	ComponentEnabled(name string) bool
}
