package messages

var topics Topics

type Topics struct {
	Like            string `json:"like"             required:"true"`
	Fork            string `json:"fork"             required:"true"`
	Training        string `json:"training"         required:"true"`
	Finetune        string `json:"finetune"         required:"true"`
	Following       string `json:"following"        required:"true"`
	Inference       string `json:"inference"        required:"true"`
	Evaluate        string `json:"evaluate"         required:"true"`
	Submission      string `json:"submission"       required:"true"`
	OperateLog      string `json:"operate_log"      required:"true"`
	RelatedResource string `json:"related_resource" required:"true"`
}