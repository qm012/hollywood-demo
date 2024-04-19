package constant

const (
	GlobalFilmActorNum = 15 // 每局游戏的演员数量，每局游戏会有多个回合

	GlobalCaptureRoundNum   int = 5 // 每局游戏的拍摄回合数，拍摄结束后进入推广回合
	GlobalPromotionRoundNum int = 3 // 每局游戏的推广回合数

	// GlobalTotalRoundNum 一局游戏的总回合数
	GlobalTotalRoundNum = GlobalCaptureRoundNum + GlobalPromotionRoundNum

	// GlobalRoundActorNum  每回合生成演员的人数
	// 		Deprecated:
	GlobalRoundActorNum = 2
)
