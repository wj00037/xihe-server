package repositories

import (
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

func (impl project) IncreaseFork(index *domain.ResourceIndex) error {
	err := impl.mapper.IncreaseFork(
		index.ResourceOwner.Account(),
		index.ResourceId,
	)
	if err != nil {
		err = convertError(err)
	}

	return err
}

func (impl project) AddLike(owner domain.Account, pid string) error {
	err := impl.mapper.AddLike(owner.Account(), pid)
	if err != nil {
		err = convertError(err)
	}

	return err
}

func (impl project) RemoveLike(owner domain.Account, pid string) error {
	err := impl.mapper.RemoveLike(owner.Account(), pid)
	if err != nil {
		err = convertError(err)
	}

	return err
}

func (impl project) AddRelatedModel(info *repository.RelatedResourceInfo) error {
	do := toRelatedResourceDO(info)

	if err := impl.mapper.AddRelatedModel(&do); err != nil {
		return convertError(err)
	}

	return nil
}

func (impl project) RemoveRelatedModel(info *repository.RelatedResourceInfo) error {
	do := toRelatedResourceDO(info)

	if err := impl.mapper.RemoveRelatedModel(&do); err != nil {
		return convertError(err)
	}

	return nil
}

func (impl project) AddRelatedDataset(info *repository.RelatedResourceInfo) error {
	do := toRelatedResourceDO(info)

	if err := impl.mapper.AddRelatedDataset(&do); err != nil {
		return convertError(err)
	}

	return nil
}

func (impl project) RemoveRelatedDataset(info *repository.RelatedResourceInfo) error {
	do := toRelatedResourceDO(info)

	if err := impl.mapper.RemoveRelatedDataset(&do); err != nil {
		return convertError(err)
	}

	return nil
}

func (impl project) UpdateProperty(info *repository.ProjectPropertyUpdateInfo) error {
	p := &info.Property

	do := ProjectPropertyDO{
		ResourceToUpdateDO: toResourceToUpdateDO(&info.ResourceToUpdate),

		Name:     p.Name.ProjName(),
		Desc:     p.Desc.ResourceDesc(),
		CoverId:  p.CoverId.CoverId(),
		RepoType: p.RepoType.RepoType(),
		Tags:     p.Tags,
	}

	if err := impl.mapper.UpdateProperty(&do); err != nil {
		return convertError(err)
	}

	return nil
}

type ProjectPropertyDO struct {
	ResourceToUpdateDO

	Name     string
	Desc     string
	CoverId  string
	RepoType string
	Tags     []string
}

func toRelatedResourceDO(info *repository.RelatedResourceInfo) RelatedResourceDO {
	return RelatedResourceDO{
		ResourceToUpdateDO: toResourceToUpdateDO(&info.ResourceToUpdate),
		ResourceOwner:      info.RelatedResource.ResourceOwner.Account(),
		ResourceId:         info.RelatedResource.ResourceId,
	}
}

type RelatedResourceDO struct {
	ResourceToUpdateDO

	ResourceOwner string
	ResourceId    string
}

type ResourceToUpdateDO struct {
	Id        string
	Owner     string
	Version   int
	UpdatedAt int64
}

func toResourceToUpdateDO(info *repository.ResourceToUpdate) ResourceToUpdateDO {
	return ResourceToUpdateDO{
		Id:        info.Id,
		Owner:     info.Owner.Account(),
		Version:   info.Version,
		UpdatedAt: info.UpdatedAt,
	}
}
