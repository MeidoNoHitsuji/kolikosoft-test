package shutdowner

import (
	"context"
	"os"
	"os/signal"
	"sort"
	"sync"
	"syscall"

	"github.com/rs/zerolog"
)

type Shutdowner struct {
	ctx         context.Context
	logger      *zerolog.Logger
	shutdowners []shutdownOption
	cancelApp   context.CancelFunc
}

func NewShutdowner(ctx context.Context, logger *zerolog.Logger) *Shutdowner {
	return &Shutdowner{ctx: ctx,
		logger: logger,
	}
}

func (s *Shutdowner) AddShutdownOption(finalFunc func() error, priority int) {
	option := shutdownOption{finalFunc: finalFunc, priority: priority}
	s.shutdowners = append(s.shutdowners, option)
}

func (s *Shutdowner) AddCancelApp(cancelApp context.CancelFunc) {
	s.cancelApp = cancelApp
}

func (s *Shutdowner) WaitShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(quit)

	select {
	case <-quit:
	case <-s.ctx.Done():
	}

	s.shutdown()
}

func (s *Shutdowner) shutdown() {
	if len(s.shutdowners) > 0 {
		var (
			finalFuncList [][]finalFunc
			finalF        []finalFunc
		)
		sort.SliceStable(s.shutdowners, func(i, j int) bool {
			return s.shutdowners[i].priority < s.shutdowners[j].priority
		})

		tmpPriority := s.shutdowners[0].priority
		for _, shutdownOpt := range s.shutdowners {
			if shutdownOpt.priority != tmpPriority {
				tmpPriority = shutdownOpt.priority
				finalFuncList = append(finalFuncList, finalF)
				finalF = nil
			}
			finalF = append(finalF, shutdownOpt.finalFunc)
		}
		finalFuncList = append(finalFuncList, finalF)
		// Выполняем отложенные функции завершения приложения
		for _, f := range finalFuncList {
			s.callShutdown(f)
		}
	}
	if s.cancelApp != nil {
		s.cancelApp()
	}
}

func (s *Shutdowner) callShutdown(finalFuncList []finalFunc) {
	var wg sync.WaitGroup
	for _, f := range finalFuncList {
		wg.Add(1)
		go func() {
			defer func() {
				wg.Done()
			}()
			defer func() {
				if r := recover(); r != nil {
					s.logger.Warn().Any("Recovered", r).Msg("Recovered from panic")
				}
			}()
			err := f()
			if err != nil {
				s.logger.Error().Err(err).Msg("ошибка завершения работы")
			}
		}()
	}
	wg.Wait()
}
