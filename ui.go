package main

import (
	"context"
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	connections  = 5
	downloadTime = 10 * time.Second
	uploadTime   = 10 * time.Second
	sparkWidth   = 20
	tickInterval = time.Second / 10
)

var (
	accentColor = lipgloss.Color("#2EF8BB")
	dimColor    = lipgloss.Color("240")

	speedStyle = lipgloss.NewStyle().Bold(true)
	unitStyle  = lipgloss.NewStyle().Foreground(dimColor)
	sparkStyle = lipgloss.NewStyle().Foreground(accentColor)
	peakStyle  = lipgloss.NewStyle().Foreground(dimColor)
	baseStyle  = lipgloss.NewStyle().Padding(1, 2)

	phaseStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("99"))
	doneStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("10"))
)

type phase int

const (
	phaseDownload phase = iota
	phaseUpload
	phaseDone
)

func (p phase) String() string {
	switch p {
	case phaseDownload:
		return "Download"
	case phaseUpload:
		return "Upload"
	default:
		return ""
	}
}

type tickMsg time.Time

func tickCmd(t time.Time) tea.Msg {
	return tickMsg(t)
}

type Model struct {
	targets []string

	bytes     *atomic.Int64
	ctx       context.Context
	cancel    context.CancelFunc
	start     time.Time
	speed     float64
	speeds    []float64
	peak      float64

	phase      phase
	done       bool
	quitting   bool

	dlSpeed     float64
	dlPeak      float64
	dlSpeeds    []float64
	ulSpeed     float64
	ulPeak      float64
	ulSpeeds    []float64
}

func NewModel(targets []string) Model {
	ctx, cancel := context.WithTimeout(context.Background(), downloadTime)
	bytes := &atomic.Int64{}

	return Model{
		targets: targets,
		bytes:   bytes,
		ctx:     ctx,
		cancel:  cancel,
		start:   time.Now(),
		phase:   phaseDownload,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(tea.Tick(tickInterval, tickCmd), m.measure)
}

func (m Model) measure() tea.Msg {
	for _, url := range m.targets {
		go download(m.ctx, url, m.bytes)
	}
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.quitting = true
			m.cancel()
			return m, tea.Quit
		}

	case tickMsg:
		elapsed := time.Since(m.start)
		m.speed = mbps(m.bytes.Load(), elapsed)
		m.speeds = append(m.speeds, m.speed)
		if m.speed > m.peak {
			m.peak = m.speed
		}

		switch m.phase {
		case phaseDownload:
			if elapsed >= downloadTime {
				m.dlSpeed = m.speed
				m.dlPeak = m.peak
				m.dlSpeeds = make([]float64, len(m.speeds))
				copy(m.dlSpeeds, m.speeds)

				m.cancel()
				m.phase = phaseUpload
				m.bytes = &atomic.Int64{}
				m.speed = 0
				m.speeds = nil
				m.peak = 0
				m.start = time.Now()

				ctx, cancel := context.WithTimeout(context.Background(), uploadTime)
				m.ctx = ctx
				m.cancel = cancel

				return m, tea.Batch(tea.Tick(tickInterval, tickCmd), m.measureUpload)
			}
			return m, tea.Tick(tickInterval, tickCmd)

		case phaseUpload:
			if elapsed >= uploadTime {
				m.ulSpeed = m.speed
				m.ulPeak = m.peak
				m.ulSpeeds = make([]float64, len(m.speeds))
				copy(m.ulSpeeds, m.speeds)

				m.cancel()
				m.phase = phaseDone
				m.done = true
				return m, tea.Quit
			}
			return m, tea.Tick(tickInterval, tickCmd)
		}
	}

	return m, nil
}

func (m Model) measureUpload() tea.Msg {
	for _, url := range m.targets {
		go upload(m.ctx, url, m.bytes)
	}
	return nil
}

func (m Model) View() string {
	if m.quitting {
		return ""
	}

	var s strings.Builder

	switch m.phase {
	case phaseDownload:
		s.WriteString(phaseStyle.Render("  Download"))
		s.WriteString("\n\n")
		speed, unit := scale(m.speed)
		s.WriteString(speedStyle.Render(fmt.Sprintf("%5.1f", speed)))
		s.WriteString(unitStyle.Render(" "+unit))
		s.WriteString(" ")
		s.WriteString(sparkStyle.Render(sparkline(m.speeds, m.peak, sparkWidth)))
		if m.peak > 0 {
			peak, peakUnit := scale(m.peak)
			label := fmt.Sprintf("  peak %.0f", peak)
			if peakUnit != unit {
				label += " " + peakUnit
			}
			s.WriteString(peakStyle.Render(label))
		}

	case phaseUpload:
		s.WriteString(phaseStyle.Render("  Upload"))
		s.WriteString("\n\n")
		speed, unit := scale(m.speed)
		s.WriteString(speedStyle.Render(fmt.Sprintf("%5.1f", speed)))
		s.WriteString(unitStyle.Render(" "+unit))
		s.WriteString(" ")
		s.WriteString(sparkStyle.Render(sparkline(m.speeds, m.peak, sparkWidth)))
		if m.peak > 0 {
			peak, peakUnit := scale(m.peak)
			label := fmt.Sprintf("  peak %.0f", peak)
			if peakUnit != unit {
				label += " " + peakUnit
			}
			s.WriteString(peakStyle.Render(label))
		}

	case phaseDone:
		s.WriteString(doneStyle.Render("  Speed Test Complete"))
		s.WriteString("\n\n")

		dlSpeed, dlUnit := scale(m.dlSpeed)
		ulSpeed, ulUnit := scale(m.ulSpeed)

		s.WriteString("  Download  ")
		s.WriteString(speedStyle.Render(fmt.Sprintf("%6.1f %s", dlSpeed, dlUnit)))
		s.WriteString("\n")
		s.WriteString("  Upload    ")
		s.WriteString(speedStyle.Render(fmt.Sprintf("%6.1f %s", ulSpeed, ulUnit)))
		s.WriteString("\n")

		if m.dlPeak > 0 {
			peak, peakUnit := scale(m.dlPeak)
			s.WriteString("  Peak DL   ")
			s.WriteString(peakStyle.Render(fmt.Sprintf("%6.1f %s", peak, peakUnit)))
			s.WriteString("\n")
		}
		if m.ulPeak > 0 {
			peak, peakUnit := scale(m.ulPeak)
			s.WriteString("  Peak UL   ")
			s.WriteString(peakStyle.Render(fmt.Sprintf("%6.1f %s", peak, peakUnit)))
		}
	}

	return baseStyle.Render(s.String())
}
