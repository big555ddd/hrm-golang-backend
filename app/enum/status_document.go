package enum

type StatusDocument string

const (
	STATUS_DOCUMENT_PENDING  StatusDocument = "pending"
	STATUS_DOCUMENT_APPROVED StatusDocument = "approved"
	STATUS_DOCUMENT_REJECTED StatusDocument = "rejected"
)

func GetStatusDocument(t StatusDocument) StatusDocument {
	switch t {
	case STATUS_DOCUMENT_PENDING:
		return STATUS_DOCUMENT_PENDING
	case STATUS_DOCUMENT_APPROVED:
		return STATUS_DOCUMENT_APPROVED
	case STATUS_DOCUMENT_REJECTED:
		return STATUS_DOCUMENT_REJECTED
	default:
		return STATUS_DOCUMENT_PENDING
	}
}
