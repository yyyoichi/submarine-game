package game

import (
	"fmt"
	"math"
	"slices"
	"sync"
	"time"

	apiv1 "github.com/yyyoichi/submarine-game/internal/gen/api/v1"
)

type Game struct {
	Id     string
	Users  [2]string
	Island [2]uint32
	clock  clock

	mu        sync.RWMutex
	NextUser  *string
	Winner    string
	histories []history
	mines     map[string][]uint32 // 残機雷
}

type history struct {
	user   string
	camp   uint32
	t      apiv1.ActionType
	impact bombImpactType // t がbombのときの結果
	mines  []uint32       // 残機雷
	at     time.Time
}

func (g *Game) GetHistory(me string) *apiv1.HistoryResponse {
	g.mu.RLock()
	defer g.mu.RUnlock()

	var latestMyHistory = g.getLatestHistory(me)
	var latestPlaceHistory = g.getLatestPlaceHistory(me)
	var resp = &apiv1.HistoryResponse{
		Camps:     g.getCampStatus(latestPlaceHistory, latestMyHistory),
		MyTurn:    false,
		Winner:    g.Winner,
		Histories: make([]*apiv1.History, 0, len(g.histories)),
		Timeout:   g.clock.getActionTimeout().UnixMilli(),
	}
	if g.NextUser == nil || *g.NextUser == me {
		resp.MyTurn = true
	}
	if resp.Winner != "" {
		resp.MyTurn = false
		resp.Timeout = 0
	}

	// description
	switch {
	case g.Winner != "":
		if me == g.Winner {
			resp.Description = "🏆勝利！！"
		} else {
			resp.Description = "💣敗北.."
		}
	case resp.MyTurn:
		if latestMyHistory == nil {
			resp.Description = "📍行動を開始する海域を決定して、機雷を敷設しよう。"
		} else {
			resp.Description = "🪖行動か、魚雷か、機雷か。"
		}
	default:
		resp.Description = "👀敵の行動を待機中.."
	}

	// histories
	for trun, hist := range g.loopHistories() {
		// histユーザの前回の配置
		var prevPlaceHistory = g.getPrevPlaceHistory(hist.user, trun)

		// respにpushする構造体
		var respHistory apiv1.History

		// description
		switch hist.user {
		// 操作ユーザの履歴
		case me:
			respHistory = apiv1.History{
				UserId: hist.user,
				Turn:   trun,
				Camp:   hist.camp,
				Type:   hist.t,
			}
			switch respHistory.Type {
			case apiv1.ActionType_ACTION_TYPE_FIRST:
				respHistory.Description = fmt.Sprintf("📍作戦開始海域を'%d'に決定。", hist.camp)

			case apiv1.ActionType_ACTION_TYPE_MOVE:
				respHistory.Description = fmt.Sprintf("🌊海域'%d'に移動。", hist.camp)

			case apiv1.ActionType_ACTION_TYPE_LEAVE:
				respHistory.Description = "🏳️敗走した。"

			case apiv1.ActionType_ACTION_TYPE_BOMB:
				respHistory.Description = fmt.Sprintf("💣海域'%d'に魚雷発射！", hist.camp)

			case apiv1.ActionType_ACTION_TYPE_MINE:
				respHistory.Description = fmt.Sprintf("海域'%d'の機雷作動！", hist.camp)

			}

		// 対戦相手の履歴
		default:
			mask := hist.mask()
			respHistory = apiv1.History{
				UserId: mask.user,
				Turn:   trun,
				Camp:   mask.camp,
				Type:   mask.t,
			}

			switch respHistory.Type {
			case apiv1.ActionType_ACTION_TYPE_FIRST:
				respHistory.Description = "📍作戦開始海域を'?'に決定。"

			case apiv1.ActionType_ACTION_TYPE_MOVE:
				var jp string
				dir := calcDirection(prevPlaceHistory.camp, hist.camp)
				switch dir {
				case north:
					jp = "北"
				case south:
					jp = "南"
				case west:
					jp = "西"
				case east:
					jp = "東"
				}
				respHistory.Description = fmt.Sprintf("🌊%s進。", jp)

			case apiv1.ActionType_ACTION_TYPE_LEAVE:
				respHistory.Description = "🏳️敗走した。"

			case apiv1.ActionType_ACTION_TYPE_BOMB:
				respHistory.Description = fmt.Sprintf("💣海域'%d'に魚雷発射！", hist.camp)

			case apiv1.ActionType_ACTION_TYPE_MINE:
				respHistory.Description = fmt.Sprintf("💣海域'%d'の機雷作動！", hist.camp)
			}

		}

		// impact
		switch hist.impact {
		case meichu:
			respHistory.Impact = "🎯命中！！"
		case omokaji:
			respHistory.Impact = "💡面舵一杯！！"
		case yosoro:
			respHistory.Impact = "🧭ヨーソロー！！"
		}

		resp.Histories = append(resp.Histories, &respHistory)
	}

	return resp
}

func (g *Game) FirstAction(me string, place uint32, mines []uint32) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.Winner != "" {
		return ErrGameIsOver
	}
	if g.clock.isActionTimeout() {
		// timeoutよりも500milsec大きいときにタイムアウト判定
		g.leave(me)
		return ErrTimeout
	}
	if place >= campSize {
		return ErrOutOfCampSize
	}
	if len(mines) != 2 {
		return fmt.Errorf("%w: mines must have 2 length", ErrInvalidAction)
	}
	if mines[0] >= campSize || mines[1] >= campSize {
		return ErrOutOfCampSize
	}
	var latestHistory = g.getLatestHistory(me)
	var latestPlaceHistory = g.getLatestPlaceHistory(me)

	// check enable action or not
	var enableCamps = g.getCampStatus(latestPlaceHistory, latestHistory)
	// valid place
	row, col := place/lineSize, place%lineSize
	enableStatus := enableCamps[row].Camps[col].Status
	if !slices.Contains(enableStatus, apiv1.CampStatus_CAMP_STATUS_PLACE) {
		return fmt.Errorf("%w: %s", ErrInvalidAction, apiv1.ActionType_ACTION_TYPE_FIRST)
	}
	// valid mines place
	for _, camp := range mines {
		row, col := camp/lineSize, camp%lineSize
		enableStatus := enableCamps[row].Camps[col].Status
		if !slices.Contains(enableStatus, apiv1.CampStatus_CAMP_STATUS_MINE) {
			return fmt.Errorf("%w: %s", ErrInvalidAction, apiv1.ActionType_ACTION_TYPE_FIRST)
		}
	}

	g.changeTurn(me)
	g.appendHistory(history{
		user:  me,
		camp:  place,
		t:     apiv1.ActionType_ACTION_TYPE_FIRST,
		mines: mines,
	})
	g.mines[me] = mines
	return nil
}

func (g *Game) Action(me string, camp uint32, action apiv1.ActionType) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.Winner != "" {
		return ErrGameIsOver
	}
	if g.NextUser == nil || *g.NextUser != me {
		return ErrIsnotYourTurn
	}
	if g.clock.isActionTimeout() {
		// timeoutよりも500milsec大きいときにタイムアウト判定
		g.leave(me)
		return ErrTimeout
	}
	if camp >= campSize {
		return ErrOutOfCampSize
	}

	var latestHistory = g.getLatestHistory(me)
	var latestPlaceHistory = g.getLatestPlaceHistory(me)

	// check enable action or not
	var enableCamps = g.getCampStatus(latestPlaceHistory, latestHistory)
	row, col := camp/lineSize, camp%lineSize
	enableStatus := enableCamps[row].Camps[col].Status
	switch action {
	case apiv1.ActionType_ACTION_TYPE_MOVE:
		if !slices.Contains(enableStatus, apiv1.CampStatus_CAMP_STATUS_MOVE) {
			return fmt.Errorf("%w: %s", ErrInvalidAction, action)
		}

	case apiv1.ActionType_ACTION_TYPE_BOMB:
		if !slices.Contains(enableStatus, apiv1.CampStatus_CAMP_STATUS_BOMB) {
			return fmt.Errorf("%w: %s", ErrInvalidAction, action)
		}

	case apiv1.ActionType_ACTION_TYPE_MINE:
		if !slices.Contains(enableStatus, apiv1.CampStatus_CAMP_STATUS_MINE) {
			return fmt.Errorf("%w: %s", ErrInvalidAction, action)
		}

	case apiv1.ActionType_ACTION_TYPE_FIRST:
		return fmt.Errorf("%w: %s", ErrInvalidAction, action)

	case apiv1.ActionType_ACTION_TYPE_UNSPECIFIED:
		return fmt.Errorf("%w: %s", ErrInvalidAction, action)

	case apiv1.ActionType_ACTION_TYPE_LEAVE:
	}
	// end validation

	defer g.changeTurn(me)
	if action == apiv1.ActionType_ACTION_TYPE_MINE {
		// 機雷作動分を削除
		g.mines[me] = slices.DeleteFunc(g.mines[me], func(c uint32) bool {
			return c == camp
		})
	}

	hist := history{
		user:   me,
		camp:   camp,
		t:      action,
		impact: unknownImpact,
		mines:  make([]uint32, len(g.mines[me])),
	}
	_ = copy(hist.mines, g.mines[me])

	// impact
	switch action {
	// 魚雷と機雷のときは命中判定をする。
	case apiv1.ActionType_ACTION_TYPE_BOMB, apiv1.ActionType_ACTION_TYPE_MINE:
		enemy := g.getEnemy(me)
		ehist := g.getLatestPlaceHistory(enemy)
		if ehist != nil { // nilは実装上あり得ない
			hist.impact = bombImpact(camp, ehist.camp)
		}
	}
	if hist.impact == meichu {
		g.Winner = me
	}
	g.appendHistory(hist)
	return nil
}

func (g *Game) Leave(me string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.leave(me)
}

func (g *Game) leave(user string) {
	if g.Winner != "" {
		return
	}
	g.appendHistory(history{
		user: user,
		camp: 0,
		t:    apiv1.ActionType_ACTION_TYPE_LEAVE,
	})
	enemy := g.getEnemy(user)
	g.Winner = enemy
}

func (g *Game) appendHistory(h history) {
	now := time.Now()
	h.at = now
	g.clock.latestAt = &now
	g.histories = append(g.histories, h)
}

// 最新の配置場所[latestHistory]と最新の行動内容から、海域情報を返す。
func (g *Game) getCampStatus(latestPlace *history, latestHist *history) []*apiv1.HistoryResponse_Line {
	if latestPlace == nil {
		return g.getInitCampStatus()
	}

	var resp = make([]*apiv1.HistoryResponse_Line, lineSize)
	moveEnable := getUDLRCamps(latestPlace.camp)
	bombEnable := getAdjacentCamps(latestPlace.camp)

	for c := range uint32(campSize) {
		row := c / lineSize
		col := c % lineSize
		if col == 0 {
			resp[row] = &apiv1.HistoryResponse_Line{
				Camps: make([]*apiv1.HistoryResponse_Camp, lineSize),
			}
		}
		if slices.Contains(g.Island[:], c) {
			// 島
			resp[row].Camps[col] = &apiv1.HistoryResponse_Camp{
				Status: []apiv1.CampStatus{apiv1.CampStatus_CAMP_STATUS_ISLAND},
				Camp:   c,
			}
			continue
		}

		var status = make([]apiv1.CampStatus, 0, 3)

		// 機雷敷設場所
		if slices.Contains(latestHist.mines, c) {
			status = append(status, apiv1.CampStatus_CAMP_STATUS_MINE)
		}

		if latestPlace.camp == c {
			// 自身
			status = append(status, apiv1.CampStatus_CAMP_STATUS_SUBMARINE)
			resp[row].Camps[col] = &apiv1.HistoryResponse_Camp{
				Status: status,
				Camp:   c,
			}
			continue
		}
		if slices.Contains(moveEnable, c) {
			status = append(status, apiv1.CampStatus_CAMP_STATUS_MOVE)
		}
		if slices.Contains(bombEnable, c) {
			status = append(status, apiv1.CampStatus_CAMP_STATUS_BOMB)
		}
		resp[row].Camps[col] = &apiv1.HistoryResponse_Camp{
			Status: status,
			Camp:   c,
		}
	}
	return resp
}

// 初期値の移動可能海域情報を返す。
func (g *Game) getInitCampStatus() []*apiv1.HistoryResponse_Line {
	var resp = make([]*apiv1.HistoryResponse_Line, lineSize)
	for c := range uint32(campSize) {
		row := c / lineSize
		col := c % lineSize
		if col == 0 {
			resp[row] = &apiv1.HistoryResponse_Line{
				Camps: make([]*apiv1.HistoryResponse_Camp, lineSize),
			}
		}
		var status []apiv1.CampStatus
		if slices.Contains(g.Island[:], c) {
			status = []apiv1.CampStatus{
				apiv1.CampStatus_CAMP_STATUS_ISLAND,
			}
		} else {
			status = []apiv1.CampStatus{
				apiv1.CampStatus_CAMP_STATUS_PLACE, apiv1.CampStatus_CAMP_STATUS_MINE,
			}
		}
		resp[row].Camps[col] = &apiv1.HistoryResponse_Camp{
			Status: status,
			Camp:   c,
		}
	}
	return resp
}

// ユーザの最新の行動履歴を返す。
func (g *Game) getLatestHistory(user string) *history {
	for i := len(g.histories) - 1; 0 <= i; i-- {
		if g.histories[i].user == user {
			return &g.histories[i]
		}
	}
	return nil
}

// ユーザの最新の配置履歴を返す。
func (g *Game) getLatestPlaceHistory(user string) *history {
	for i := len(g.histories) - 1; 0 <= i; i-- {
		if g.histories[i].user != user {
			continue
		}
		if g.histories[i].t == apiv1.ActionType_ACTION_TYPE_MOVE || g.histories[i].t == apiv1.ActionType_ACTION_TYPE_FIRST {
			return &g.histories[i]
		}
	}
	return nil
}

// ユーザの[trun]未満の最新の行動を返す。
func (g *Game) getPrevPlaceHistory(user string, trun int32) *history {
	for i := int(trun) - 2; 0 <= i; i-- {
		if g.histories[i].user != user {
			continue
		}
		if g.histories[i].t == apiv1.ActionType_ACTION_TYPE_MOVE || g.histories[i].t == apiv1.ActionType_ACTION_TYPE_FIRST {
			return &g.histories[i]
		}
	}
	return nil
}

// trunNumberと行動履歴を返す。
func (g *Game) loopHistories() map[int32]history {
	var resp = make(map[int32]history, len(g.histories))
	for i, h := range g.histories {
		resp[int32(i)+1] = h
	}
	return resp
}

// 相手のidを返す。
func (g *Game) getEnemy(me string) string {
	if g.Users[0] == me {
		return g.Users[1]
	}
	return g.Users[0]
}

func (g *Game) changeTurn(user string) {
	if g.Users[0] == user {
		g.NextUser = &g.Users[1]
	} else {
		g.NextUser = &g.Users[0]
	}
}

func (h history) mask() history {
	h.user = ""
	if h.t == apiv1.ActionType_ACTION_TYPE_MOVE {
		// 移動の場合場所はマスクする
		h.camp = 0
	}
	return h
}

type direction int

const (
	unknownDirection direction = iota
	north
	south
	west
	east
)

// [b]eforeから[a]fterへの移動方向
func calcDirection(b, a uint32) direction {
	if a == b {
		return unknownDirection
	}
	i1, i2 := int(a), int(b)

	row1, col1 := i1/lineSize, i1%lineSize
	row2, col2 := i2/lineSize, i2%lineSize

	rowDiff := row1 - row2
	colDiff := col1 - col2
	if rowDiff == 0 {
		if colDiff == -1 {
			return west
		}
		if colDiff == 1 {
			return east
		}
	}
	if colDiff == 0 {
		if rowDiff == -1 {
			return north
		}
		if rowDiff == 1 {
			return south
		}
	}
	return unknownDirection
}

type bombImpactType int

const (
	unknownImpact bombImpactType = iota
	meichu
	omokaji
	yosoro
)

// [camp]海域に対する魚雷[bomb]の結果を返す。
func bombImpact(bomb uint32, camp uint32) bombImpactType {
	if bomb == camp {
		return meichu
	}
	a := isAdjacent(int(bomb), int(camp))
	if a {
		return omokaji
	}
	return yosoro
}

// 上下左右の関係にあるか
func isUDLR(camp1, camp2 uint32) bool {
	if camp1 == camp2 {
		return false
	}
	i1, i2 := int(camp1), int(camp2)

	row1, col1 := i1/lineSize, i1%lineSize
	row2, col2 := i2/lineSize, i2%lineSize

	rowDiff := int(math.Abs(float64(row1 - row2)))
	colDiff := int(math.Abs(float64(col1 - col2)))

	return (rowDiff == 1 && colDiff == 0) || (rowDiff == 0 && colDiff == 1)
}

// 上下左右海域リスト
func getUDLRCamps(camp uint32) []uint32 {
	var adjCamps = make([]uint32, 0, 4)

	// Calculate row and column for the given camp
	row, col := int(camp)/lineSize, int(camp)%lineSize

	// Loop through all possible adjacent camps (including diagonally)
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			// Skip the same camp
			if i == 0 && j == 0 {
				continue
			}
			// veri or hori
			if i != 0 && j != 0 {
				continue
			}
			// Calculate the new row and column
			newRow, newCol := row+i, col+j
			// Check if the new camp is within the bounds of the grid
			if newRow >= 0 && newRow < lineSize && newCol >= 0 && newCol < lineSize {
				// Calculate the camp number and add it to the list
				adjCamp := newRow*lineSize + newCol
				adjCamps = append(adjCamps, uint32(adjCamp))
			}
		}
	}

	return adjCamps
}

// 上下左右斜めの関係にあるか
func isAdjacent(camp1, camp2 int) bool {
	if camp1 == camp2 {
		return false
	}
	i1, i2 := int(camp1), int(camp2)

	row1, col1 := i1/lineSize, i1%lineSize
	row2, col2 := i2/lineSize, i2%lineSize

	rowDiff := int(math.Abs(float64(row1 - row2)))
	colDiff := int(math.Abs(float64(col1 - col2)))

	return (rowDiff <= 1 && colDiff <= 1)
}

// 上下左右斜め海域リスト
func getAdjacentCamps(camp uint32) []uint32 {
	var adjCamps = make([]uint32, 0, 8)

	// Calculate row and column for the given camp
	row, col := int(camp)/lineSize, int(camp)%lineSize

	// Loop through all possible adjacent camps (including diagonally)
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			// Skip the same camp
			if i == 0 && j == 0 {
				continue
			}
			// Calculate the new row and column
			newRow, newCol := row+i, col+j
			// Check if the new camp is within the bounds of the grid
			if newRow >= 0 && newRow < lineSize && newCol >= 0 && newCol < lineSize {
				// Calculate the camp number and add it to the list
				adjCamp := newRow*lineSize + newCol
				adjCamps = append(adjCamps, uint32(adjCamp))
			}
		}
	}

	return adjCamps
}
