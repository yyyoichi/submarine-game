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

	stateMachine *stateless.StateMachine
	istate       *istate
	setIstate    state.SetStateFunc[istate]

	ebiten.Game
}

func New() *Screen {
	var s Screen
	s.istate, s.setIstate = state.UseState[istate](state.WithInitialValue(istate{}))
	s.stateMachine = stateless.NewStateMachine(pstate_sector)
	// 開始海域選定 <-> 機雷設置 -> 待機

	// 開始海域選定 -> 機雷設置
	s.stateMachine.Configure(pstate_sector).Permit(trigger_mines, pstate_mines, func(context.Context, ...any) bool {
		// istateに開始海域が設定されていたら許可
		// return s.istate.selectedSector != nil
		return true
	})
	//  開始海域選定 <- 機雷設置
	s.stateMachine.Configure(pstate_mines).Permit(trigger_sector, pstate_sector, func(context.Context, ...any) bool {
		// 機雷設置画面から行動開始位置設定は無条件で許可
		return true
	})

	// 機雷設置 -> 待機
	s.stateMachine.Configure(pstate_mines).Permit(trigger_pending, pstate_pending, func(context.Context, ...any) bool {
		// istateに機雷が設定されていたら許可
		// return len(s.istate.selectedMines) == int(s.istate.config.MineCount)
		return true
	})

	// TODO APIから取得
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
				s.Game = newMinesScreen()
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

// func (s *Screen) Draw(screen *ebiten.Image) {
// 	// 1. 初めの一秒間は黒背景に「作戦開始海域の選定」
// 	// 2. その後、海域の選定画面
// 	// 3. 完了したらpenndingに遷移
// 	s.game.Draw(screen)
// }

// func (s *Screen) Layout(outsideWidth, outsideHeight int) (int, int) {
// 	return config.ScreenWidth, config.ScreenHeight
// }
