package metadata

// MetadataRepository is a noop implementation of the metadata repository
type MetadataRepository struct{}

// NewMetadataRepository creates a new metadata repository
func NewMetadataRepository() *MetadataRepository {
	return &MetadataRepository{}
}
