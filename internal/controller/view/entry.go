package view

import "vibecheck/internal/domain"

// ModalView is passed to the log-modal template.
// Open controls whether the overlay is visible; Error shows inline validation feedback.
type ModalView struct {
	Open  bool
	Error string
}

// HomeView is passed to the home page and entry-success templates.
type HomeView struct {
	ModalView
	TodayEntry *domain.Entry
}
