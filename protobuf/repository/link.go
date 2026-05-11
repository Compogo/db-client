package repository

import (
	"github.com/Compogo/db-client/proto.back"
	"github.com/Compogo/db-client/repository"
	"github.com/Compogo/types/linker"
)

var (
	SortDirectionToProtoSortDirection = linker.NewLinker[repository.SortDirection, proto_back.SortDirection](
		linker.Link(repository.ASC, proto_back.SortDirection_ASC),
		linker.Link(repository.DESC, proto_back.SortDirection_DESC),
	)

	ProtoSortDirectionToSortDirection = linker.NewLinker[proto_back.SortDirection, repository.SortDirection](
		linker.Link(proto_back.SortDirection_ASC, repository.ASC),
		linker.Link(proto_back.SortDirection_DESC, repository.DESC),
	)

	ComparableToProtoComparable = linker.NewLinker[repository.Comparable, proto_back.Comparable](
		linker.Link(repository.Eq, proto_back.Comparable_Eq),
		linker.Link(repository.Neq, proto_back.Comparable_Neq),
		linker.Link(repository.Gt, proto_back.Comparable_Gt),
		linker.Link(repository.Gte, proto_back.Comparable_Gte),
		linker.Link(repository.Lt, proto_back.Comparable_Lt),
		linker.Link(repository.Lte, proto_back.Comparable_Lte),
		linker.Link(repository.LIKE, proto_back.Comparable_LIKE),
		linker.Link(repository.IN, proto_back.Comparable_IN),
	)

	ProtoComparableToComparable = linker.NewLinker[proto_back.Comparable, repository.Comparable](
		linker.Link(proto_back.Comparable_Eq, repository.Eq),
		linker.Link(proto_back.Comparable_Neq, repository.Neq),
		linker.Link(proto_back.Comparable_Gt, repository.Gt),
		linker.Link(proto_back.Comparable_Gte, repository.Gte),
		linker.Link(proto_back.Comparable_Lt, repository.Lt),
		linker.Link(proto_back.Comparable_Lte, repository.Lte),
		linker.Link(proto_back.Comparable_LIKE, repository.LIKE),
		linker.Link(proto_back.Comparable_IN, repository.IN),
	)
)
