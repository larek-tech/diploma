package file

import (
	"time"

	"github.com/google/uuid"
)

// Archive формат данных при получении файлов от пользователей подразумевает преобразование в плоскую директорию
// и хранение в базе данных
type Archive struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}

type File struct {
	ID        string    `db:"id"`
	SourceID  string    `db:"source_id"`
	Filename  string    `db:"filename"`
	Extension string    `db:"extension"`
	Raw       []byte    `db:"-" json:"-"`
	Size      int64     `db:"size"`
	ObjectURL string    `db:"object_key"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func NewFile(sourceID string, extension string) *File {
	return &File{
		ID:        uuid.NewString(),
		SourceID:  sourceID,
		Extension: extension,
		Raw:       []byte{},
		Size:      0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
