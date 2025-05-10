package object_store

type ObjectStore struct {
	ID       string `db:"id"`        // ID uuid идентификатор источника
	SourceID string `db:"source_id"` // ID uuid идентификатора источника
	Config   Config `db:"config"`    // Config конфигурация для доступа к объектному хранилищу
}

type Config struct {
	Endpoint  string `json:"endpoint"`
	Url       string `json:"url"`
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
	Api       string `json:"api"`
	Path      string `json:"path"`
}

type Object struct {
	ID           string `db:"id"` // ID uuid идентификатор объекта
	Data         []byte // Data данные объекта
	Size         int    `db:"content_size"`  // Size размер объекта
	ContentType  string `db:"content_type"`  // ContentType тип содержимого объекта
	RawContentID string `db:"raw_object_id"` // raw_content_id идентификатор в dwh хранилище
	CreatedAt    string `db:"created_at"`    // CreatedAt дата создания объекта
	UpdatedAt    string `db:"updated_at"`    // UpdatedAt дата обновления объекта
}
