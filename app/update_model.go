package app

import (
	"errors"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/platform"
	"github.com/opensourceways/xihe-server/domain/repository"
)

type ModelUpdateCmd struct {
	Name     domain.ModelName
	Desc     domain.ResourceDesc
	RepoType domain.RepoType
}

func (cmd *ModelUpdateCmd) toModel(
	p *domain.ModelModifiableProperty, repo *platform.RepoOption,
) (b bool) {
	f := func() {
		if !b {
			b = true
		}
	}

	if cmd.Name != nil && p.Name.ModelName() != cmd.Name.ModelName() {
		p.Name = cmd.Name
		repo.Name = cmd.Name
		f()
	}

	if cmd.Desc != nil && p.Desc.ResourceDesc() != cmd.Desc.ResourceDesc() {
		p.Desc = cmd.Desc
		f()
	}

	if cmd.RepoType != nil && p.RepoType.RepoType() != cmd.RepoType.RepoType() {
		p.RepoType = cmd.RepoType
		repo.RepoType = cmd.RepoType
		f()
	}

	return
}

func (s modelService) Update(
	p *domain.Model, cmd *ModelUpdateCmd, pr platform.Repository,
) (dto ModelDTO, err error) {
	opt := new(platform.RepoOption)
	if !cmd.toModel(&p.ModelModifiableProperty, opt) {
		s.toModelDTO(p, &dto)

		return

	}

	v, err := s.repo.Save(p)
	if err != nil {
		return
	}

	if opt.IsNotEmpty() {
		if err = pr.Update(p.RepoId, opt); err != nil {
			return
		}
	}

	s.toModelDTO(&v, &dto)

	return
}

func (s modelService) AddLike(owner domain.Account, rid string) error {
	return s.repo.AddLike(owner, rid)
}

func (s modelService) RemoveLike(owner domain.Account, rid string) error {
	return s.repo.RemoveLike(owner, rid)
}

func (s modelService) AddRelatedDataset(
	m *domain.Model, index *domain.ResourceIndex,
) error {
	if m.RelatedDatasets.Has(index) {
		return nil
	}

	if m.RelatedDatasets.Count()+1 > m.MaxRelatedResourceNum() {
		return ErrorExceedMaxRelatedResourceNum{
			errors.New("exceed max related reousrce num"),
		}
	}

	info := repository.RelatedResourceInfo{
		Owner:           m.Owner,
		ResourceId:      m.Id,
		Version:         m.Version,
		RelatedResource: *index,
	}

	return s.repo.AddRelatedDataset(&info)
}

func (s modelService) RemoveRelatedDataset(
	m *domain.Model, index *domain.ResourceIndex,
) error {
	if !m.RelatedDatasets.Has(index) {
		return nil
	}

	info := repository.RelatedResourceInfo{
		Owner:           m.Owner,
		ResourceId:      m.Id,
		Version:         m.Version,
		RelatedResource: *index,
	}

	return s.repo.RemoveRelatedDataset(&info)
}