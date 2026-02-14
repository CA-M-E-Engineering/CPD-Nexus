package sgbuildex

// PushRequest is the top-level request sent to SGBuildex
// It represents ONE submission event
type PushRequest struct {
	Participants []ParticipantWrapper `json:"participants"`
	Payload      []any                `json:"payload"`
	OnBehalfOf   []OnBehalfWrapper    `json:"on_behalf_of,omitempty"` // submission-level
}

// ParticipantWrapper represents an individual (e.g. worker)
// who may act on behalf of another entity (e.g. subcontractor)
type ParticipantWrapper struct {
	ID         string           `json:"id"`
	Name       string           `json:"name"`
	Meta       *ParticipantMeta `json:"meta,omitempty"`
	OnBehalfOf *OnBehalfWrapper `json:"on_behalf_of,omitempty"` // participant-level
}

// ParticipantMeta holds contextual identifiers
type ParticipantMeta struct {
	DataRefID string `json:"data_ref_id,omitempty"` // e.g. project ref
}

// OnBehalfWrapper represents an organisation or authority
// that a participant or submission is acting on behalf of
type OnBehalfWrapper struct {
	ID   string `json:"id"`
	Name string `json:"name,omitempty"`
}
