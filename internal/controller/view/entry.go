package view

// ModalView is passed to the log-modal template.
// Open controls whether the overlay is visible; Error shows inline validation feedback.
type ModalView struct {
	Open  bool
	Error string
}
