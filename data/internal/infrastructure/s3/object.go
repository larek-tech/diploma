package s3

import "fmt"

type ContentType string

const (
	ContentTypeJSON ContentType = "application/json"
	ContentTypeXML  ContentType = "application/xml"
	ContentTypeText ContentType = "text/plain"
	ContentTypeHTML ContentType = "text/html"
	ContentTypeCSV  ContentType = "text/csv"
	ContentTypePDF  ContentType = "application/pdf"
	ContentTypeDocx ContentType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	ContentTypeXlsx ContentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
)

var ErrObjectNotFound = fmt.Errorf("object not found")
var ErrObjectAlreadyExists = fmt.Errorf("object already exists")

type Object struct {
	bucketName  string
	Key         string
	data        []byte
	contentType ContentType
	metadata    map[string]string
	versionID   string
}

func (o *Object) GetBucketName() string {
	return o.bucketName
}
func (o *Object) GetKey() string {
	return o.Key
}
func (o *Object) GetData() []byte {
	return o.data
}
func (o *Object) GetContentType() ContentType {
	return o.contentType
}
func (o *Object) GetMetadata() map[string]string {
	return o.metadata
}
func (o *Object) GetVersionID() string {
	return o.versionID
}

func NewObject(bucketName, key string, data []byte, contentType ContentType, metadata ...map[string]string) *Object {
	obj := &Object{
		bucketName:  bucketName,
		Key:         key,
		data:        data,
		contentType: contentType,
		metadata:    make(map[string]string),
	}

	if len(metadata) > 0 && metadata[0] != nil {
		obj.metadata = metadata[0]
	}

	return obj
}
