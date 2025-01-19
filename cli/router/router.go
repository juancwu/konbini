package router

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss"
)

type NavigationMsg struct {
	To     string
	Params map[string]interface{}
}

type PageBuilder func(params map[string]interface{}) tea.Model

type Router struct {
	head         *HistoryNode
	current      *HistoryNode
	pageBuilders map[string]PageBuilder
}

type HistoryNode struct {
	route  string
	params map[string]interface{}
	model  tea.Model
	next   *HistoryNode
	prev   *HistoryNode
}

func NewRouter() *Router {
	return &Router{
		head:         nil,
		current:      nil,
		pageBuilders: make(map[string]PageBuilder),
	}
}

func NewNavigationMsg(to string, params map[string]interface{}) tea.Cmd {
	return func() tea.Msg {
		return NavigationMsg{
			To:     to,
			Params: params,
		}
	}
}

func (r *Router) RegisterPage(route string, builder PageBuilder) error {
	_, ok := r.pageBuilders[route]
	if ok {
		return fmt.Errorf("Page with route '%s' already exists", route)
	}

	r.pageBuilders[route] = builder

	return nil
}

func (r *Router) SetInitialPage(route string, params map[string]interface{}) (tea.Cmd, error) {
	builder, ok := r.pageBuilders[route]
	if !ok {
		return nil, fmt.Errorf("No page builder registered for route: %s", route)
	}

	model, cmd := r.initalizePage(builder, params)

	r.head = &HistoryNode{
		route:  route,
		params: params,
		model:  model,
	}
	r.current = r.head

	return cmd, nil
}

func (r *Router) Navigate(to string, params map[string]interface{}) (tea.Cmd, error) {
	if r.current.route == to {
		return nil, nil
	}

	builder, ok := r.pageBuilders[to]
	if !ok {
		return nil, fmt.Errorf("No page builder registered for route: %s", to)
	}

	var cmds []tea.Cmd

	if hook, ok := r.current.model.(LifecycleHooks); ok {
		if !hook.BeforeNavigateAway() {
			return nil, fmt.Errorf("Navigation cancelled by current page: %s", r.current.route)
		}

		if cmd := hook.AfterNavigateAway(); cmd != nil {
			cmds = append(cmds, cmd)
		}

		if cmd := hook.OnExit(); cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	model, cmd := r.initalizePage(builder, params)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	node := &HistoryNode{
		route:  to,
		model:  model,
		params: params,
		next:   nil,
		prev:   r.current,
	}

	r.current.next = node
	r.current = node

	return tea.Batch(cmds...), nil
}

func (r *Router) Back() (tea.Model, tea.Cmd) {
	if r.current.prev == nil {
		return r.current.model, nil
	}

	var cmds []tea.Cmd

	if hook, ok := r.current.model.(LifecycleHooks); ok {
		if !hook.BeforeNavigateAway() {
			return r.current.model, nil
		}

		if cmd := hook.AfterNavigateAway(); cmd != nil {
			cmds = append(cmds, cmd)
		}

		if cmd := hook.OnExit(); cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	r.current = r.current.prev

	// Only call OnEnter when going back (Init already called during initialization of page)
	if hook, ok := r.current.model.(LifecycleHooks); ok {
		if cmd := hook.OnEnter(r.current.params); cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	return r.current.model, tea.Batch(cmds...)
}

func (r *Router) Forward() (tea.Model, tea.Cmd) {
	if r.current.next == nil {
		return r.current.model, nil
	}

	var cmds []tea.Cmd

	if hook, ok := r.current.model.(LifecycleHooks); ok {
		if !hook.BeforeNavigateAway() {
			return r.current.model, nil
		}

		if cmd := hook.AfterNavigateAway(); cmd != nil {
			cmds = append(cmds, cmd)
		}

		if cmd := hook.OnExit(); cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	r.current = r.current.next

	// Only call OnEnter when going forward (skip Init since page was already initialized)
	if hook, ok := r.current.model.(LifecycleHooks); ok {
		if cmd := hook.OnEnter(r.current.params); cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	return r.current.model, tea.Batch(cmds...)
}

func (r *Router) CurrentModel() tea.Model {
	return r.current.model
}

func (r *Router) UpdateCurrentModel(model tea.Model) {
	r.current.model = model
}

func (r *Router) HistoryString() string {
	var builder strings.Builder
	var node *HistoryNode = r.head
	var i int = 0
	for node != nil {
		if i > 0 {
			builder.WriteString(" -> ")
		}
		if r.current.route == node.route {
			builder.WriteString(lipgloss.NewStyle().Underline(true).Render(node.route))
		} else {
			builder.WriteString(node.route)
		}
		node = node.next
		i += 1
	}
	return builder.String()
}

func (r *Router) initalizePage(builder PageBuilder, params map[string]interface{}) (tea.Model, tea.Cmd) {
	model := builder(params)

	var cmds []tea.Cmd
	var cmd tea.Cmd

	model, cmd = model.Init()
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	if hook, ok := model.(LifecycleHooks); ok {
		if cmd := hook.OnEnter(params); cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	return model, tea.Batch(cmds...)
}
