package messages

import (
	bigmodeldomain "github.com/opensourceways/xihe-server/bigmodel/domain"
	common "github.com/opensourceways/xihe-server/common/domain/message"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/message"
	"github.com/opensourceways/xihe-server/utils"
)

var _ message.Sender = (*sender)(nil)

func NewMessageSender(topic *Topics, p common.Publisher) *sender {
	return &sender{topics: *topic, publisher: p}
}

type sender struct {
	topics    Topics
	publisher common.Publisher
}

// Like
func (s *sender) AddLike(msg *domain.ResourceObject) error {
	return s.sendLike(msg, actionAdd)
}

func (s *sender) RemoveLike(msg *domain.ResourceObject) error {
	return s.sendLike(msg, actionRemove)
}

func (s *sender) sendLike(msg *domain.ResourceObject, action string) error {
	v := msgLike{Action: action}

	toMsgResourceObject(msg, &v.Resource)

	return s.send(s.topics.Like, &v)
}

// Finetune
func (s *sender) CreateFinetune(info *domain.FinetuneIndex) error {
	v := msgFinetune{
		User: info.Owner.Account(),
		Id:   info.Id,
	}

	return s.send(s.topics.Finetune, &v)
}

// Inference
func (s *sender) CreateInference(info *domain.InferenceInfo) error {
	v := s.toInferenceMsg(&info.InferenceIndex)
	v.Action = actionCreate
	v.ProjectName = info.ProjectName.ResourceName()
	v.ResourceLevel = info.ResourceLevel

	return s.send(s.topics.Inference, &v)

}

func (s *sender) ExtendInferenceSurvivalTime(info *message.InferenceExtendInfo) error {
	v := s.toInferenceMsg(&info.InferenceIndex)
	v.Action = actionExtend
	v.Expiry = info.Expiry
	v.ProjectName = info.ProjectName.ResourceName()
	v.ResourceLevel = info.ResourceLevel

	return s.send(s.topics.Inference, &v)
}

func (s *sender) toInferenceMsg(index *domain.InferenceIndex) msgInference {
	return msgInference{
		ProjectId:    index.Project.Id,
		LastCommit:   index.LastCommit,
		InferenceId:  index.Id,
		ProjectOwner: index.Project.Owner.Account(),
	}
}

// Evaluate
func (s *sender) CreateEvaluate(info *message.EvaluateInfo) error {
	v := msgEvaluate{
		Type:         info.Type,
		OBSPath:      info.OBSPath,
		ProjectId:    info.Project.Id,
		TrainingId:   info.TrainingId,
		EvaluateId:   info.Id,
		ProjectOwner: info.Project.Owner.Account(),
	}

	return s.send(s.topics.Evaluate, &v)
}

// Competition
func (s *sender) CalcScore(info *message.SubmissionInfo) error {
	v := msgSubmission{
		CId:   info.Index.Id,
		Phase: info.Index.Phase.CompetitionPhase(),
		SId:   info.Id,
		Path:  info.OBSPath,
	}

	return s.send(s.topics.Submission, &v)
}

// operate log
func (s *sender) AddOperateLogForAccessBigModel(u domain.Account, t bigmodeldomain.BigmodelType) error {
	return s.sendOperateLog(u, "bigmodel", map[string]string{
		"bigmodel": string(t),
	})
}

func (s *sender) sendOperateLog(u domain.Account, t string, info map[string]string) error {
	return s.send(s.topics.OperateLog, operateLog(u, t, info))
}

// send
func (s *sender) send(topic string, v interface{}) error {
	return s.publisher.Publish(topic, v, nil)
}

func operateLog(u domain.Account, t string, info map[string]string) common.MsgOperateLog {
	a := ""
	if u != nil {
		a = u.Account()
	}

	return common.MsgOperateLog{
		When: utils.Now(),
		User: a,
		Type: t,
		Info: info,
	}
}
