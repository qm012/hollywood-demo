package constant

type playerPromotionOptions []*PlayerPromotionOption

// PlayerPromotionOption 玩家推广选项
type PlayerPromotionOption struct {
}

func (p *PlayerPromotionOption) PhaseName() string {
	return ""
}

// PlayerPromotionOptionOutcome 玩家推广选项结果（基于每个选项有不同的概率）
type PlayerPromotionOptionOutcome string

const (
	PlayerPromotionOptionOutcomeSuccess         PlayerPromotionOptionOutcome = "success"
	PlayerPromotionOptionOutcomeFailure         PlayerPromotionOptionOutcome = "failure"
	PlayerPromotionOptionOutcomeCriticalSuccess PlayerPromotionOptionOutcome = "criticalSuccess"
	PlayerPromotionOptionOutcomeCriticalFailure PlayerPromotionOptionOutcome = "criticalFailure"
)

func (p PlayerPromotionOptionOutcome) String() string {
	switch p {
	case PlayerPromotionOptionOutcomeSuccess:
		return "Movie quality goes UP!"
	case PlayerPromotionOptionOutcomeFailure:
		return "Movie quality goes DOWN!"
	case PlayerPromotionOptionOutcomeCriticalSuccess:
		return "Movie quality drastically goes UP!"
	case PlayerPromotionOptionOutcomeCriticalFailure:
		return "Movie quality drastically goes down! Movie gets more coverage on SNS."
	}
	return ""
}
