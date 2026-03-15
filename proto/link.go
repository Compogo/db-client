package proto

import (
	"github.com/Compogo/db-client/repository"
	"github.com/Compogo/types/linker"
)

var (
	SortDirectionToProtoSortDirection = linker.NewLinker[repository.SortDirection, SortDirection](
		linker.NewLink(repository.ASC, SortDirection_ASC),
		linker.NewLink(repository.DESC, SortDirection_DESC),
	)

	ProtoSortDirectionToSortDirection = linker.NewLinker[SortDirection, repository.SortDirection](
		linker.NewLink(SortDirection_ASC, repository.ASC),
		linker.NewLink(SortDirection_DESC, repository.DESC),
	)

	ComparableToProtoComparable = linker.NewLinker[repository.Comparable, Comparable](
		linker.NewLink(repository.Eq, Comparable_Eq),
		linker.NewLink(repository.Neq, Comparable_Neq),
		linker.NewLink(repository.Gt, Comparable_Gt),
		linker.NewLink(repository.Gte, Comparable_Gte),
		linker.NewLink(repository.Lt, Comparable_Lt),
		linker.NewLink(repository.Lte, Comparable_Lte),
		linker.NewLink(repository.LIKE, Comparable_LIKE),
		linker.NewLink(repository.IN, Comparable_IN),
	)

	ProtoComparableToComparable = linker.NewLinker[Comparable, repository.Comparable](
		linker.NewLink(Comparable_Eq, repository.Eq),
		linker.NewLink(Comparable_Neq, repository.Neq),
		linker.NewLink(Comparable_Gt, repository.Gt),
		linker.NewLink(Comparable_Gte, repository.Gte),
		linker.NewLink(Comparable_Lt, repository.Lt),
		linker.NewLink(Comparable_Lte, repository.Lte),
		linker.NewLink(Comparable_LIKE, repository.LIKE),
		linker.NewLink(Comparable_IN, repository.IN),
	)
)
