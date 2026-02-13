package metadata

// MetadataUseCase is a noop implementation of the metadata usecase
type MetadataUseCase struct{}

// NewMetadataUseCase creates a new metadata usecase
func NewMetadataUseCase() *MetadataUseCase {
	return &MetadataUseCase{}
}
