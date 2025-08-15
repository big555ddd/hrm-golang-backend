package enum

type DocumentType string

const (
	DOCUMENT_TYPE_LEAVE    DocumentType = "leave"
	DOCUMENT_TYPE_OVERTIME DocumentType = "overtime"
	DOCUMENT_TYPE_ADDTIME  DocumentType = "addTime"
)

func GetDocumentType(t DocumentType) DocumentType {
	switch t {
	case DOCUMENT_TYPE_LEAVE:
		return DOCUMENT_TYPE_LEAVE
	case DOCUMENT_TYPE_OVERTIME:
		return DOCUMENT_TYPE_OVERTIME
	case DOCUMENT_TYPE_ADDTIME:
		return DOCUMENT_TYPE_ADDTIME
	default:
		return DOCUMENT_TYPE_LEAVE
	}
}
