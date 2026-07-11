package tui

import (
	"context"
	"io"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"github.com/aethelgards/octo/llm"
	"github.com/aethelgards/octo/structs"
	"github.com/cloudwego/eino/schema"
)

const logo = "\033[38;2;102;126;234m  ██████╗   ██████╗ ████████╗  ██████╗\033[0m\n" +
	"\033[38;2;118;134;240m ██╔═══██╗ ██╔════╝ ╚══██╔══╝ ██╔═══██╗\033[0m\n" +
	"\033[38;2;134;142;246m ██║   ██║ ██║         ██║    ██║   ██║\033[0m\n" +
	"\033[38;2;150;150;252m ██║   ██║ ██║         ██║    ██║   ██║\033[0m\n" +
	"\033[38;2;180;120;234m ╚██████╔╝ ╚██████╗    ██║    ╚██████╔╝\033[0m\n" +
	"\033[38;2;210;90;220m  ╚═════╝   ╚═════╝    ╚═╝     ╚═════╝\033[0m\n"

type streamMsg struct {
	content string
	err     error
	isDone  bool
	stream  *schema.StreamReader[*schema.Message]
}

type OctoModel struct {
	ctx      context.Context
	config   *structs.OctoConfig
	input    textinput.Model
	history  []*schema.Message
	response string
	stream   *schema.StreamReader[*schema.Message]
}

func (m *OctoModel) Init() tea.Cmd {
	return m.input.Focus()
}

func (m *OctoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if msg.String() == "ctrl+c" {
			if m.stream != nil {
				m.stream.Close()
			}
			return m, tea.Quit
		} else if msg.String() == "enter" {
			value := m.input.Value()
			if value == "" {
				return m, nil
			}
			m.clear()
			m.history = append(m.history, schema.UserMessage(value))
			return m, m.startStream()
		}
	case streamMsg:
		if msg.stream != nil {
			m.stream = msg.stream
			return m, m.readStream(m.stream)
		}
		if msg.err != nil {
			m.response = "错误: " + msg.err.Error()
			if m.stream != nil {
				m.stream.Close()
				m.stream = nil
			}
			return m, nil
		}
		if msg.isDone {
			if m.response != "" {
				m.history = append(m.history, schema.AssistantMessage(m.response, nil))
			}
			m.stream = nil
			return m, nil
		}
		if msg.content != "" {
			m.response += msg.content
		}
		return m, m.readStream(m.stream)
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m *OctoModel) startStream() tea.Cmd {
	return func() tea.Msg {
		stream, err := llm.ChatStream(m.ctx, m.history)
		if err != nil {
			return streamMsg{err: err}
		}
		return streamMsg{stream: stream}
	}
}

func (m *OctoModel) readStream(stream *schema.StreamReader[*schema.Message]) tea.Cmd {
	return func() tea.Msg {
		if stream == nil {
			return streamMsg{isDone: true}
		}

		msg, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return streamMsg{isDone: true}
			}
			return streamMsg{err: err}
		}

		return streamMsg{content: msg.Content}
	}
}

func (m *OctoModel) View() tea.View {
	view := logo
	view += "\nModel " + m.config.LLMConfig.Model + "\n"
	if m.response != "" {
		view += m.response + "\n\n"
	}

	view += "> " + m.input.View() + "\n"
	view += "(Ctrl+C 退出)\n"

	return tea.NewView(view)
}

func (m *OctoModel) clear() {
	m.input.SetValue("")
	m.response = ""
}

func Init(ctx context.Context, config *structs.OctoConfig) error {
	input := textinput.New()
	input.Placeholder = ""
	input.Focus()
	input.CharLimit = 500
	input.SetWidth(60)

	p := tea.NewProgram(&OctoModel{
		input:  input,
		ctx:    ctx,
		config: config,
	})
	_, err := p.Run()
	return err
}
