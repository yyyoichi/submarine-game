package game

import (
	"fmt"
	"math"
	"sync"
	"time"

	apiv1 "github.com/yyyoichi/submarine-game/internal/gen/api/v1"
)

type Game struct {
	Id        string
	Users     [2]string
	Island    [2]int64
	createdAt time.Time

	mu        sync.RWMutex
	NextUser  string
	Winner    string
	histories []history
}

type history struct {
	user string
	camp uint32
	t    apiv1.ActionType
	at   time.Time
}

func (g *Game) GetHistory(me string) *apiv1.HistoryResponse {
	g.mu.RLock()
	defer g.mu.RUnlock()

	resp := &apiv1.HistoryResponse{
		Island:    g.Island[:],
		MyTurn:    me == g.NextUser,
		Winner:    g.Winner,
		Histories: make([]*apiv1.History, len(g.histories)),
		Timeout:   g.getTimeout().UnixMilli(),
	}
	switch {
	case g.Winner != "":
		if me == g.Winner {
			resp.Description = "🏆勝利！！"
		} else {
			resp.Description = "💣敗北.."
		}
	case resp.MyTurn:
		if len(g.histories) < 2 {
			resp.Description = "📍行動を開始する海域を決定しよう。"
			resp.EnableTypes = []apiv1.ActionType{apiv1.ActionType_ACTION_TYPE_BOMB}
		} else {
			resp.Description = "🪖行動か、魚雷か。"
			resp.EnableTypes = []apiv1.ActionType{apiv1.ActionType_ACTION_TYPE_BOMB, apiv1.ActionType_ACTION_TYPE_MOVE}
			resp.EnableCamps = getAdjacentPositions(g.histories[len(g.histories)-2].camp)
		}
	default:
		resp.Description = "👀敵の行動を待機中.."
	}

	for i, hist := range g.histories {
		var respHistory *apiv1.History
		switch hist.user {
		// 操作ユーザの履歴
		case me:
			respHistory = &apiv1.History{
				UserId: hist.user,
				Turn:   int64(i),
				Camp:   hist.camp,
				Type:   hist.t,
			}
			switch respHistory.Type {
			case apiv1.ActionType_ACTION_TYPE_MOVE:
				if i < 2 {
					respHistory.Description = fmt.Sprintf("📍作戦開始海域を'%d'に決定。", hist.camp)
				} else {
					respHistory.Description = fmt.Sprintf("🌊海域'%d'に移動。", hist.camp)
				}
			case apiv1.ActionType_ACTION_TYPE_LEAVE:
				respHistory.Description = "⚠️敗走した。"
			case apiv1.ActionType_ACTION_TYPE_BOMB:
				respHistory.Description = fmt.Sprintf("💣海域'%d'に魚雷発射！", hist.camp)
			}

		// 対戦相手の履歴
		default:
			mask := hist.mask()
			respHistory = &apiv1.History{
				UserId: mask.user,
				Turn:   int64(i),
				Camp:   mask.camp,
				Type:   mask.t,
			}
			switch respHistory.Type {
			case apiv1.ActionType_ACTION_TYPE_MOVE:
				if i < 2 {
					respHistory.Description = "📍作戦開始海域を'?'に決定。"
				} else {
					b := g.histories[i-2].camp
					a := hist.camp
					dir := calcDirection(b, a)
					var jp string
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
					respHistory.Description = fmt.Sprintf("🌊'%s'進。", jp)
				}
			case apiv1.ActionType_ACTION_TYPE_LEAVE:
				respHistory.Description = "✨敗走した。"
			case apiv1.ActionType_ACTION_TYPE_BOMB:
				respHistory.Description = fmt.Sprintf("💣海域'%d'に魚雷発射！", hist.camp)
			}

		}
		if i < 2 {
			continue
		}
		prev := resp.Histories[i-1]
		if prev.Type != apiv1.ActionType_ACTION_TYPE_BOMB {
			continue
		}
		switch bi := bombImpact(prev.Camp, hist.camp); bi {
		case meichu:
			respHistory.Impact = "🎯命中！！"
		case omokaji:
			respHistory.Impact = "💡面舵一杯！！"
		case yosoro:
			respHistory.Impact = "🧭ヨーソロー！！"
		}

		resp.Histories[i] = respHistory
	}

	return resp
}

func (g *Game) Action(user string, camp uint32, action apiv1.ActionType) error {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if g.Winner != "" {
		return ErrGameIsOver
	}
	if g.NextUser != user {
		return ErrIsnotYourTurn
	}
	if g.isTimeout() {
		// timeoutよりも500milsec大きいときにタイムアウト判定
		g.leave(user)
		return ErrTimeout
	}
	if camp > campSize {
		return ErrOutOfCampSize
	}

	if len(g.histories) < 2 {
		// place
		g.histories = append(g.histories, history{
			user: user,
			camp: camp,
			t:    apiv1.ActionType_ACTION_TYPE_PLACE,
			at:   time.Now(),
		})
		g.changeTurn(user)
		return nil
	}

	switch action {
	case apiv1.ActionType_ACTION_TYPE_MOVE, apiv1.ActionType_ACTION_TYPE_BOMB:
		if d := calcDirection(g.histories[len(g.histories)-2].camp, camp); d == unknownDirection {
			return fmt.Errorf("%w: '%d'", ErrInvalidCamp, camp)
		}
	default:
		return fmt.Errorf("%w: %s", ErrInvalidAction, action)
	}

	g.histories = append(g.histories, history{
		user: user,
		camp: camp,
		t:    action,
		at:   time.Now(),
	})
	g.changeTurn(user)
	return nil
}

func (g *Game) Leave(user string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.leave(user)
}

func (g *Game) leave(user string) {
	g.histories = append(g.histories, history{
		user: user,
		camp: 0,
		t:    apiv1.ActionType_ACTION_TYPE_LEAVE,
		at:   time.Now(),
	})
	if g.Users[0] == user {
		g.Winner = g.Users[1]
	} else {
		g.Winner = g.Users[0]
	}
}

func (g *Game) changeTurn(user string) {
	if g.Users[0] == user {
		g.NextUser = g.Users[1]
	} else {
		g.NextUser = g.Users[0]
	}
}

func (g *Game) isTimeout() bool {
	return g.getTimeout().Add(time.Duration(500) * time.Millisecond).Before(time.Now())
}

func (g *Game) getTimeout() time.Time {
	if len(g.histories) == 0 {
		return g.createdAt.Add(time.Duration(30) * time.Second)
	}
	return g.histories[len(g.histories)-1].at.Add(time.Duration(30) * time.Second)
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

func calcDirection(b, a uint32) direction {
	switch {
	case b-a == lineSize:
		return north
	case a-b == lineSize:
		return south
	case b-a == 1:
		return west
	case a-b == 1:
		return east
	}
	return unknownDirection
}

type bombImpactType int

const (
	meichu = iota
	omokaji
	yosoro
)

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

// Check if the positions are adjacent and not the same position
func isAdjacent(pos1, pos2 int) bool {
	if pos1 == pos2 {
		return false
	}
	// Calculate row and column for both positions
	row1, col1 := pos1/lineSize, pos1%6
	row2, col2 := pos2/lineSize, pos2%6

	// Calculate the difference between the rows and columns
	rowDiff := math.Abs(float64(row1 - row2))
	colDiff := math.Abs(float64(col1 - col2))

	return (rowDiff <= 1 && colDiff <= 1)
}

// 行動可能な位置リスト
func getAdjacentPositions(pos uint32) []uint32 {
	var adjPositions = make([]uint32, 0, 8)

	// Calculate row and column for the given position
	row, col := int(pos)/lineSize, int(pos)%lineSize

	// Loop through all possible adjacent positions (including diagonally)
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			// Skip the same position
			if i == 0 && j == 0 {
				continue
			}

			// Calculate the new row and column
			newRow, newCol := row+i, col+j

			// Check if the new position is within the bounds of the grid
			if newRow >= 0 && newRow < lineSize && newCol >= 0 && newCol < lineSize {
				// Calculate the position number and add it to the list
				adjPos := newRow*lineSize + newCol
				adjPositions = append(adjPositions, uint32(adjPos))
			}
		}
	}

	return adjPositions
}
