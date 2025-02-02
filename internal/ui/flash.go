package ui

import (
	"context"

	"github.com/derailed/k9s/internal/config"
	"github.com/derailed/k9s/internal/model"
	"github.com/derailed/tview"
	"github.com/gdamore/tcell"
	"github.com/rs/zerolog/log"
)

const (
	emoHappy = "😎"
	emoDoh   = "😗"
	emoRed   = "😡"
)

// Flash represents a flash message indicator.
type Flash struct {
	*tview.TextView

	app      *App
	testMode bool
}

// NewFlash returns a new flash view.
func NewFlash(app *App) *Flash {
	f := Flash{
		app:      app,
		TextView: tview.NewTextView(),
	}
	f.SetTextColor(tcell.ColorAqua)
	f.SetTextAlign(tview.AlignCenter)
	f.SetBorderPadding(0, 0, 1, 1)
	f.app.Styles.AddListener(&f)

	return &f
}

// SetTestMode for testing ONLY!
func (f *Flash) SetTestMode(b bool) {
	f.testMode = b
}

// StylesChanged notifies listener the skin changed.
func (f *Flash) StylesChanged(s *config.Styles) {
	f.SetBackgroundColor(s.BgColor())
	f.SetTextColor(s.FgColor())
}

// Watch watches for flash changes.
func (f *Flash) Watch(ctx context.Context, c model.FlashChan) {
	defer log.Debug().Msgf("Flash Canceled!")
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-c:
			f.SetMessage(msg)
		}
	}
}

// SetMessage sets flash message and level.
func (f *Flash) SetMessage(m model.LevelMessage) {
	fn := func() {
		if m.Text == "" {
			f.Clear()
			return
		}
		f.SetTextColor(flashColor(m.Level))
		f.SetText(flashEmoji(m.Level) + " " + m.Text)
	}

	if f.testMode {
		fn()
	} else {
		f.app.QueueUpdateDraw(fn)
	}
}

func flashEmoji(l model.FlashLevel) string {
	switch l {
	case model.FlashWarn:
		return emoDoh
	case model.FlashErr:
		return emoRed
	default:
		return emoHappy
	}
}

func flashColor(l model.FlashLevel) tcell.Color {
	switch l {
	case model.FlashWarn:
		return tcell.ColorOrange
	case model.FlashErr:
		return tcell.ColorOrangeRed
	default:
		return tcell.ColorNavajoWhite
	}
}
