package helper

import (
	protobuf "github.com/Compogo/db-client/protobuf/repository"
	"github.com/Compogo/db-client/repository"
	"github.com/Compogo/types/linker"
)

var (
	SortDirectionLinkProtoSortDirection = linker.NewLinker[repository.SortDirection, protobuf.SortDirection](
		linker.Link(repository.ASC, protobuf.SortDirection_ASC),
		linker.Link(repository.DESC, protobuf.SortDirection_DESC),
	)

	ProtoSortDirectionLinkSortDirection = linker.NewLinker[protobuf.SortDirection, repository.SortDirection](
		linker.Link(protobuf.SortDirection_ASC, repository.ASC),
		linker.Link(protobuf.SortDirection_DESC, repository.DESC),
	)

	ComparableLinkProtoComparable = linker.NewLinker[repository.Comparable, protobuf.Comparable](
		linker.Link(repository.Eq, protobuf.Comparable_Eq),
		linker.Link(repository.Neq, protobuf.Comparable_Neq),
		linker.Link(repository.Gt, protobuf.Comparable_Gt),
		linker.Link(repository.Gte, protobuf.Comparable_Gte),
		linker.Link(repository.Lt, protobuf.Comparable_Lt),
		linker.Link(repository.Lte, protobuf.Comparable_Lte),
		linker.Link(repository.LIKE, protobuf.Comparable_LIKE),
		linker.Link(repository.IN, protobuf.Comparable_IN),
	)

	ProtoComparableLinkComparable = linker.NewLinker[protobuf.Comparable, repository.Comparable](
		linker.Link(protobuf.Comparable_Eq, repository.Eq),
		linker.Link(protobuf.Comparable_Neq, repository.Neq),
		linker.Link(protobuf.Comparable_Gt, repository.Gt),
		linker.Link(protobuf.Comparable_Gte, repository.Gte),
		linker.Link(protobuf.Comparable_Lt, repository.Lt),
		linker.Link(protobuf.Comparable_Lte, repository.Lte),
		linker.Link(protobuf.Comparable_LIKE, repository.LIKE),
		linker.Link(protobuf.Comparable_IN, repository.IN),
	)
)
