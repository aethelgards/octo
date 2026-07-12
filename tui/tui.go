package tui

import (
	"context"
	"io"
	"strings"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/glamour/v2"
	"charm.land/lipgloss/v2"
	"github.com/aethelgards/octo/llm"
	"github.com/aethelgards/octo/structs"
	"github.com/cloudwego/eino/schema"
)

type keyMap struct {
	Up     key.Binding
	Down   key.Binding
	Help   key.Binding
	Quit   key.Binding
	Send   key.Binding
	Scroll key.Binding
	Clear  key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Send, k.Clear, k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Scroll},
		{k.Send, k.Clear, k.Help, k.Quit},
	}
}

var keys = keyMap{
	Up:     key.NewBinding(key.WithKeys("up"), key.WithHelp("↑", "上翻")),
	Down:   key.NewBinding(key.WithKeys("down"), key.WithHelp("↓", "下翻")),
	Scroll: key.NewBinding(key.WithKeys("pgup", "pgdown"), key.WithHelp("PgUp/PgDn", "翻页")),
	Send:   key.NewBinding(key.WithKeys("enter"), key.WithHelp("Enter", "发送消息")),
	Clear:  key.NewBinding(key.WithKeys("ctrl+l"), key.WithHelp("Ctrl+L", "清空对话")),
	Quit:   key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("Ctrl+C", "退出")),
}

type streamMsg struct {
	content          string
	reasoningContent string
	err              error
	isDone           bool
	stream           *schema.StreamReader[*schema.Message]
}

type OctoModel struct {
	ctx              context.Context
	config           *structs.OctoConfig
	input            textinput.Model
	history          []*schema.Message
	reasoningHistory []string
	response         string
	reasoning        string
	stream           *schema.StreamReader[*schema.Message]
	viewport         viewport.Model
	renderer         *glamour.TermRenderer
	spinner          spinner.Model
	help             help.Model
	ready            bool
	thinking         bool
}

func (m *OctoModel) Init() tea.Cmd {
	return m.input.Focus()
}

func (m *OctoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		r, _ := glamour.NewTermRenderer(
			glamour.WithWordWrap(0),
			glamour.WithStyles(customStyle()),
		)
		m.renderer = r

		inputWidth := max(msg.Width-4, 20)
		m.input.SetWidth(inputWidth)

		viewportHeight := min(max(msg.Height-headerHeight-footerHeight, minViewportHeight), maxViewportHeight)

		if !m.ready {
			m.viewport = viewport.New(
				viewport.WithWidth(msg.Width),
				viewport.WithHeight(viewportHeight),
			)
			m.viewport.SoftWrap = true
			m.ready = true
		} else {
			m.viewport.SetWidth(msg.Width)
			m.viewport.SetHeight(viewportHeight)
		}
		m.viewport.SetContent(m.buildContent())
		return m, nil

	case tea.KeyPressMsg:
		if msg.String() == "ctrl+c" {
			if m.stream != nil {
				m.stream.Close()
			}
			return m, tea.Quit
		} else if msg.String() == "ctrl+l" {
			if m.stream != nil {
				m.stream.Close()
				m.stream = nil
			}
			m.history = nil
			m.reasoningHistory = nil
			m.response = ""
			m.reasoning = ""
			m.thinking = false
			m.viewport.SetContent(m.buildContent())
			return m, nil
		} else if msg.String() == "enter" {
			value := m.input.Value()
			if value == "" {
				return m, nil
			}
			m.clear()
			m.history = append(m.history, schema.UserMessage(value))
			m.thinking = true
			cmds := []tea.Cmd{m.startStream(), m.spinner.Tick}
			return m, tea.Batch(cmds...)
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case streamMsg:
		if msg.stream != nil {
			m.stream = msg.stream
			return m, m.readStream(m.stream)
		}
		if msg.err != nil {
			m.response = "错误: " + msg.err.Error()
			m.thinking = false
			if m.stream != nil {
				m.stream.Close()
				m.stream = nil
			}
			m.viewport.SetContent(m.buildContent())
			m.viewport.GotoBottom()
			return m, nil
		}
		if msg.isDone {
			if m.response != "" {
				m.history = append(m.history, schema.AssistantMessage(m.response, nil))
				m.reasoningHistory = append(m.reasoningHistory, m.reasoning)
			}
			m.stream = nil
			m.thinking = false
			m.viewport.SetContent(m.buildContent())
			m.viewport.GotoBottom()
			return m, nil
		}
		if msg.reasoningContent != "" {
			m.reasoning += msg.reasoningContent
		}
		if msg.content != "" {
			m.thinking = false
			m.response += msg.content
		}
		m.viewport.SetContent(m.buildContent())
		m.viewport.GotoBottom()
		return m, m.readStream(m.stream)
	}

	var cmds []tea.Cmd
	var cmd tea.Cmd
	m.help, cmd = m.help.Update(msg)
	cmds = append(cmds, cmd)
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)
	m.input, cmd = m.input.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
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

		return streamMsg{
			content:          msg.Content,
			reasoningContent: msg.ReasoningContent,
		}
	}
}

func (m *OctoModel) View() tea.View {
	var sb strings.Builder

	sb.WriteString(logo)
	sb.WriteString("\nModel ")
	sb.WriteString(m.config.LLMConfig.Model)

	if m.config.LLMConfig.Thinking.Enabled {
		sb.WriteString(" | Thinking: ")
		sb.WriteString(colorThinking)
		sb.WriteString("enabled")
		sb.WriteString(colorReset)
		if m.config.LLMConfig.Thinking.ReasoningEffort != "" {
			sb.WriteString(" (")
			sb.WriteString(m.config.LLMConfig.Thinking.ReasoningEffort)
			sb.WriteString(")")
		}
		if m.config.LLMConfig.Thinking.Show {
			sb.WriteString(" | Show: on")
		}
	} else {
		sb.WriteString(" | Thinking: disabled")
	}

	sb.WriteString("\n")

	if m.ready && (len(m.history) > 0 || m.stream != nil) {
		sb.WriteString(m.viewport.View())
		sb.WriteString("\n")
	}

	if m.thinking {
		sb.WriteString(m.spinner.View())
		sb.WriteString(" 思考中...")
		sb.WriteString("\n")
	}

	sb.WriteString(m.input.View())
	sb.WriteString("\n")
	sb.WriteString(m.help.View(keys))

	v := tea.NewView(sb.String())
	v.MouseMode = tea.MouseModeCellMotion
	return v
}

func (m *OctoModel) renderMarkdown(content string) string {
	if m.renderer == nil {
		return content
	}
	out, err := m.renderer.Render(content)
	if err != nil {
		return content
	}
	return strings.TrimRight(out, "\n")
}

func (m *OctoModel) buildContent() string {
	var sb strings.Builder

	if len(m.history) == 0 && m.stream == nil {
		sb.WriteString(colorThinking)
		sb.WriteString("输入消息开始对话，Ctrl+L 清空，Ctrl+C 退出")
		sb.WriteString(colorReset)
		sb.WriteString("\n")
		return sb.String()
	}

	assistantIdx := 0

	for _, msg := range m.history {
		switch msg.Role {
		case schema.User:
			sb.WriteString(colorUser)
			sb.WriteString("You")
			sb.WriteString(colorReset)
			sb.WriteString(": ")
			sb.WriteString(msg.Content)
			sb.WriteString("\n\n")
		case schema.Assistant:
			sb.WriteString(colorAssistant)
			sb.WriteString("Assistant")
			sb.WriteString(colorReset)
			sb.WriteString(":\n")
			if m.config.LLMConfig.Thinking.Show && assistantIdx < len(m.reasoningHistory) && m.reasoningHistory[assistantIdx] != "" {
				sb.WriteString(colorThinking)
				sb.WriteString(m.reasoningHistory[assistantIdx])
				sb.WriteString(colorReset)
				sb.WriteString("\n")
			}
			sb.WriteString(m.renderMarkdown(msg.Content))
			sb.WriteString("\n\n")
			assistantIdx++
		}
	}

	if m.stream != nil {
		sb.WriteString(colorAssistant)
		sb.WriteString("Assistant")
		sb.WriteString(colorReset)
		sb.WriteString(":\n")
		if m.config.LLMConfig.Thinking.Show && m.reasoning != "" {
			sb.WriteString(colorThinking)
			sb.WriteString(m.reasoning)
			sb.WriteString(colorReset)
			sb.WriteString("\n")
		}
		if m.response != "" {
			sb.WriteString(m.renderMarkdown(m.response))
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

func (m *OctoModel) clear() {
	m.input.SetValue("")
	m.response = ""
	m.reasoning = ""
}

func Init(ctx context.Context, config *structs.OctoConfig) error {
	input := textinput.New()
	input.Placeholder = ""
	input.Focus()
	input.CharLimit = 500

	h := help.New()
	h.ShowAll = false

	s := spinner.New()
	s.Spinner = spinner.Monkey
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))

	p := tea.NewProgram(&OctoModel{
		input:   input,
		ctx:     ctx,
		config:  config,
		help:    h,
		spinner: s,
	})
	_, err := p.Run()
	return err
}
