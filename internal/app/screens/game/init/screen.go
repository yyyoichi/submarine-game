package gameinit

import (
	"context"
	"math/rand"
	"slices"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/qmuntal/stateless"
	"github.com/yyyoichi/submarine-game/internal/app/actions"
	"github.com/yyyoichi/submarine-game/internal/app/state"
	"github.com/yyyoichi/submarine-game/internal/core"
)

type Screen struct {
	cGame core.Game
	me    string

	stateMachine *stateless.StateMachine
	istate       *istate
	setIstate    state.SetStateFunc[istate]

	ebiten.Game
}

func New() *Screen {
	var s Screen
	// TODO APIから取得
	s.me = "me"
	s.cGame = core.Game{
		GameId: "gameId",
		PlayerIds: [2]string{
			"me",
			"enemy",
		},
		OceanMap: core.DefaultOceanMap,
		Submarines: map[string]core.Submarine{
			"me":    core.DefaultSubmarine,
			"enemy": core.DefaultSubmarine,
		},
		Islands:   []core.Sector{},
		Timestamp: time.Now(),
	}
	getRangeIsland := func() core.Sector {
		for {
			if i := rand.Intn(int(s.cGame.OceanMap.H) * int(s.cGame.OceanMap.W)); !slices.Contains(s.cGame.Islands, core.Sector(i)) {
				return core.Sector(i)
			}
		}
	}
	for range s.cGame.OceanMap.IslandCount {
		s.cGame.Islands = append(s.cGame.Islands, getRangeIsland())
	}

	s.istate, s.setIstate = state.UseState[istate](state.WithInitialValue(istate{}))
	s.stateMachine = stateless.NewStateMachine(pstate_sector)
	// 開始海域選定 <-> 機雷設置 <-> 待機

	// 開始海域選定 -> 機雷設置
	s.stateMachine.Configure(pstate_sector).Permit(trigger_mines, pstate_mines, func(context.Context, ...any) bool {
		// istateに開始海域が設定されていたら許可
		return s.istate.selectedSector != nil
	})
	//  開始海域選定 <- 機雷設置
	s.stateMachine.Configure(pstate_mines).Permit(trigger_sector, pstate_sector, func(context.Context, ...any) bool {
		// 機雷設置画面から行動開始位置設定は設定なければ許可
		return len(s.istate.selectedMines) == 0
	})

	// 機雷設置 -> 待機
	s.stateMachine.Configure(pstate_mines).Permit(trigger_pending, pstate_pending, func(context.Context, ...any) bool {
		// istateに機雷が設定されていたら許可
		return len(s.istate.selectedMines) == int(s.cGame.Submarines[s.me].MineCount)
	})

	// 機雷設置 <- 待機
	s.stateMachine.Configure(pstate_pending).Permit(trigger_mines, pstate_mines, func(context.Context, ...any) bool {
		// 無条件で許可
		return true
	})

	//　初期画面
	s.Game = s.newSectorScreen()
	return &s
}

func (s *Screen) Update() error {

	s.Game.Update()

	ctx := context.Background()
	if actions.IsKeyAJustPressed() {
		switch st, _ := s.stateMachine.State(ctx); st {
		case pstate_sector:
			if err := s.stateMachine.Fire(trigger_mines); err == nil {
				s.Game = s.newMinesScreen()
			}

		case pstate_mines:
			if err := s.stateMachine.Fire(trigger_pending); err == nil {
				s.Game = newPendingScreen()
			}
		}
	}

	if actions.IsKeyBJustPressed() {
		switch st, _ := s.stateMachine.State(ctx); st {
		case pstate_mines:
			if err := s.stateMachine.Fire(trigger_sector); err == nil {
				s.Game = s.newSectorScreen()
			}

		case pstate_pending:
			if err := s.stateMachine.Fire(trigger_mines); err == nil {
				s.Game = s.newMinesScreen()
			}
		}
	}

	return nil
}

func (s *Screen) newSectorScreen() *sectorScreen {
	config := initSectorConfig{
		w:             int(s.cGame.OceanMap.W),
		h:             int(s.cGame.OceanMap.H),
		islandSectors: make([]int, len(s.cGame.Islands)),
		enableCursor: func(i int, d core.Direction) (int, bool) {
			s, ok := s.cGame.Move(core.Sector(i), d)
			return int(s), ok
		},
		setSelectedSector: func(i *int) {
			var news = istate{
				selectedSector: i,
				selectedMines:  s.istate.selectedMines,
			}
			_ = s.setIstate(news)
		},
		selectedSector: s.istate.selectedSector,
	}
	for i, v := range s.cGame.Islands {
		config.islandSectors[i] = int(v)
	}
	return newSectorScreen(config)
}

func (s *Screen) newMinesScreen() *minesScreen {
	config := initMinesConfig{
		w:             int(s.cGame.OceanMap.W),
		h:             int(s.cGame.OceanMap.H),
		islandSectors: make([]int, len(s.cGame.Islands)),
		enableCursor: func(i int, d core.Direction) (int, bool) {
			s, ok := s.cGame.Move(core.Sector(i), d)
			return int(s), ok
		},
		mineCount:       int(s.cGame.Submarines[s.me].MineCount),
		selectedSectors: &s.istate.selectedMines,
	}
	for i, v := range s.cGame.Islands {
		config.islandSectors[i] = int(v)
	}
	// 同じ位置には置けない仕様
	config.pushSelectedSector = func(i int) int {
		if slices.Contains(s.istate.selectedMines, i) {
			return len(s.istate.selectedMines)
		}
		nms := append(s.istate.selectedMines, i)
		if len(nms) > config.mineCount {
			nms = nms[1:]
		}
		var news = istate{
			selectedSector: s.istate.selectedSector,
			selectedMines:  nms,
		}
		_ = s.setIstate(news)
		return len(nms)
	}
	config.popSelectedSector = func() *int {
		if len(s.istate.selectedMines) == 0 {
			return nil
		}
		nms := s.istate.selectedMines
		i := s.istate.selectedMines[len(s.istate.selectedMines)-1]
		var news = istate{
			selectedSector: s.istate.selectedSector,
			selectedMines:  nms[:len(nms)-1],
		}
		_ = s.setIstate(news)
		return &i
	}
	return newMinesScreen(config)
}
