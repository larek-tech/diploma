package get_documents

import (
	"context"

	"github.com/google/uuid"
	"github.com/larek-tech/diploma/data/internal/data/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	documentsStore documentsStore
}

func New(documentsStore documentsStore) *Handler {
	return &Handler{
		documentsStore: documentsStore,
	}
}

func (h Handler) GetDocuments(ctx context.Context, in *pb.GetDocumentsIn) (*pb.GetDocumentsOut, error) {
	err := validate(in)
	if err != nil {
		return nil, err
	}
	total, documents, err := h.documentsStore.GetMany(ctx, in.SourceId, int(in.Page), int(in.Size))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get documents: %v", err)
	}
	if total == 0 {
		return nil, status.Error(codes.NotFound, "documents not found")
	}
	pbDocs := make([]*pb.Document, 0, len(documents))
	for _, v := range documents {
		pbDocs = append(pbDocs, &pb.Document{
			Id:       v.ID,
			SourceId: v.SourceID,
			Name:     v.Name,
			Content:  v.Content,
			Metadata: "",
		})
	}

	return &pb.GetDocumentsOut{
		Size:      in.Size,
		Page:      in.Page,
		Total:     uint32(total),
		Documents: []*pb.Document{},
	}, nil
}

func validate(in *pb.GetDocumentsIn) error {
	if in.SourceId == "" {
		return status.Errorf(codes.InvalidArgument, "empty source id")
	}
	_, err := uuid.Parse(in.SourceId)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "invalid source id: %v", err)
	}
	if in.Page < 0 {
		return status.Errorf(codes.InvalidArgument, "invalid page value must be positive int greater than 0")
	}
	if in.Size < 0 || in.Size > 50 {
		return status.Errorf(codes.InvalidArgument, "invalid size value must be positive int not greater than 50")
	}
	return nil
}
