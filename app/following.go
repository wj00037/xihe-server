package app

import (
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

type FollowDTO struct {
	Account  string `json:"account"`
	AvatarId string `json:"avatar_id"`
	Bio      string `json:"bio"`
}

func (s userService) AddFollowing(user, follower domain.Account) error {
	f := domain.FollowerInfo{
		User:     user,
		Follower: follower,
	}
	err := s.repo.AddFollowing(&f)
	if err != nil {
		return err
	}

	// TODO: activity

	// send event
	_ = s.sender.AddFollowing(f)

	return nil
}

func (s userService) RemoveFollowing(user, follower domain.Account) error {
	f := domain.FollowerInfo{
		User:     user,
		Follower: follower,
	}
	err := s.repo.RemoveFollowing(&f)
	if err != nil {
		return err
	}

	// send event
	_ = s.sender.RemoveFollowing(f)

	return nil
}

func (s userService) ListFollowing(owner domain.Account) (
	dtos []FollowDTO, err error,
) {
	v, err := s.repo.FindFollowing(owner, repository.FollowFindOption{})
	if err != nil || len(v) == 0 {
		return
	}

	dtos = make([]FollowDTO, len(v))
	for i := range v {
		s.toFollowDTO(&v[i], &dtos[i])
	}

	return
}

func (s userService) toFollowDTO(f *domain.FollowUserInfo, dto *FollowDTO) {
	*dto = FollowDTO{
		Account:  f.Account.Account(),
		AvatarId: f.AvatarId.AvatarId(),
	}

	if f.Bio != nil {
		dto.Bio = f.Bio.Bio()
	}
}
