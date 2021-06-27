package graveler

import (
	"context"
	"errors"
	"regexp"

	"github.com/treeverse/lakefs/pkg/ident"
)

type Store interface {
	GetBranch(ctx context.Context, repositoryID RepositoryID, branchID BranchID) (*Branch, error)
	GetTag(ctx context.Context, repositoryID RepositoryID, tagID TagID) (*CommitID, error)
	GetCommitByPrefix(ctx context.Context, repositoryID RepositoryID, prefix CommitID) (*Commit, error)
	GetCommit(ctx context.Context, repositoryID RepositoryID, prefix CommitID) (*Commit, error)
}

type revResolverFunc func(context.Context, Store, ident.AddressProvider, RepositoryID, string) (*ResolvedRef, error)

var hashRegexp = regexp.MustCompile("^[a-fA-F0-9]{1,64}$")

func isAHash(part string) bool {
	return hashRegexp.MatchString(part)
}

// revResolve return the first resolve of 'rev' - by hash, branch or tag
func revResolve(ctx context.Context, store Store, addressProvider ident.AddressProvider, repositoryID RepositoryID, rev string) (*ResolvedRef, error) {
	resolvers := []revResolverFunc{revResolveAHash, revResolveBranch, revResolveTag}
	for _, resolveHelper := range resolvers {
		r, err := resolveHelper(ctx, store, addressProvider, repositoryID, rev)
		if err != nil {
			return nil, err
		}
		if r != nil {
			return r, nil
		}
	}
	return nil, ErrNotFound
}

func ResolveRef(ctx context.Context, store Store, addressProvider ident.AddressProvider, repositoryID RepositoryID, rev ParsedRev) (*ResolvedRef, error) {
	rr, err := revResolve(ctx, store, addressProvider, repositoryID, rev.BaseRev)
	if err != nil {
		return nil, err
	}
	// return the matched reference, when no modifiers on ref or use the commit id as base
	if len(rev.Modifiers) == 0 {
		return rr, nil
	}
	baseCommit := rr.CommitID

	for _, mod := range rev.Modifiers {
		// lastly, apply modifier
		switch mod.Type {
		case RevModTypeTilde:
			// skip mod.ValueNumeric iterations
			for i := 0; i < mod.Value; i++ {
				commit, err := store.GetCommit(ctx, repositoryID, baseCommit)
				if err != nil {
					return nil, err
				}
				if len(commit.Parents) == 0 {
					return nil, ErrNotFound
				}
				baseCommit = commit.Parents[0]
			}
		case RevModTypeCaret:
			switch mod.Value {
			case 0:
				continue // ^0 = the commit itself
			default:
				// get the commit and extract parents
				c, err := store.GetCommitByPrefix(ctx, repositoryID, baseCommit)
				if err != nil {
					return nil, err
				}
				if mod.Value > len(c.Parents) {
					return nil, ErrInvalidRef
				}
				baseCommit = c.Parents[mod.Value-1]
			}

		default:
			return nil, ErrInvalidRef
		}
	}

	return &ResolvedRef{
		Type:     ReferenceTypeCommit,
		CommitID: baseCommit,
	}, nil
}

func revResolveAHash(ctx context.Context, store Store, addressProvider ident.AddressProvider, repositoryID RepositoryID, rev string) (*ResolvedRef, error) {
	if !isAHash(rev) {
		return nil, nil
	}
	commit, err := store.GetCommitByPrefix(ctx, repositoryID, CommitID(rev))
	if errors.Is(err, ErrNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	commitID := CommitID(addressProvider.ContentAddress(commit))
	return &ResolvedRef{
		Type:     ReferenceTypeCommit,
		CommitID: commitID,
	}, nil
}

func revResolveBranch(ctx context.Context, store Store, _ ident.AddressProvider, repositoryID RepositoryID, rev string) (*ResolvedRef, error) {
	branchID := BranchID(rev)
	branch, err := store.GetBranch(ctx, repositoryID, branchID)
	if errors.Is(err, ErrNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &ResolvedRef{
		Type:         ReferenceTypeBranch,
		BranchID:     branchID,
		StagingToken: branch.StagingToken,
		CommitID:     branch.CommitID,
	}, nil
}

func revResolveTag(ctx context.Context, store Store, _ ident.AddressProvider, repositoryID RepositoryID, rev string) (*ResolvedRef, error) {
	commitID, err := store.GetTag(ctx, repositoryID, TagID(rev))
	if errors.Is(err, ErrNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &ResolvedRef{
		Type:     ReferenceTypeTag,
		CommitID: *commitID,
	}, nil
}
