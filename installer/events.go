package installer

import (
	"encoding/base64"
	"errors"

	"github.com/flynn/flynn/pkg/random"
)

type httpPrompt struct {
	ID       string `json:"id"`
	Type     string `json:"type,omitempty"`
	Message  string `json:"message,omitempty"`
	Yes      bool   `json:"yes,omitempty"`
	Input    string `json:"input,omitempty"`
	Resolved bool   `json:"resolved,omitempty"`
	resChan  chan *httpPrompt
	api      *httpAPI
	stack    *Stack
}

func (prompt *httpPrompt) Resolve(res *httpPrompt) {
	prompt.stack.promptsMutex.Lock()
	defer prompt.stack.promptsMutex.Unlock()
	prompt.Resolved = true
	prompt.resChan <- res
}

type httpEvent struct {
	Type        string      `json:"type"`
	Description string      `json:"description,omitempty"`
	Prompt      *httpPrompt `json:"prompt,omitempty"`
}

type Subscription struct {
	EventIndex int
	EventChan  chan *httpEvent
	DoneChan   chan struct{}
	done       bool
}

func (sub *Subscription) sendEvents(s *Stack) {
	if sub.done {
		return
	}
	for index, event := range s.getEvents(sub.EventIndex) {
		sub.EventIndex = index
		sub.EventChan <- event
	}
}

func (sub *Subscription) handleDone() {
	if sub.done {
		return
	}
	sub.done = true
	close(sub.DoneChan)
}

func (s *Stack) Subscribe(eventChan chan *httpEvent) <-chan struct{} {
	s.subscribeMtx.Lock()
	defer s.subscribeMtx.Unlock()

	subscription := &Subscription{
		EventIndex: -1,
		EventChan:  eventChan,
		DoneChan:   make(chan struct{}),
	}

	go func() {
		subscription.sendEvents(s)
		if s.done {
			subscription.handleDone()
		}
	}()

	s.subscriptions = append(s.subscriptions, subscription)

	return subscription.DoneChan
}

func (s *Stack) getEvents(sinceIndex int) []*httpEvent {
	events := make([]*httpEvent, 0, len(s.events))
	for index, event := range s.events {
		if index <= sinceIndex {
			continue
		}
		events = append(events, event)
	}
	return events
}

func (s *Stack) findPrompt(id string) (*httpPrompt, error) {
	s.promptsMutex.Lock()
	defer s.promptsMutex.Unlock()
	for _, p := range s.Prompts {
		if p.ID == id {
			return p, nil
		}
	}
	return nil, errors.New("Prompt not found")
}

func (s *Stack) addPrompt(prompt *httpPrompt) {
	s.promptsMutex.Lock()
	defer s.promptsMutex.Unlock()
	s.Prompts = append(s.Prompts, prompt)
}

func (s *Stack) YesNoPrompt(msg string) bool {
	prompt := &httpPrompt{
		ID:      random.Hex(16),
		Type:    "yes_no",
		Message: msg,
		resChan: make(chan *httpPrompt),
		api:     s.api,
		stack:   s,
	}
	s.addPrompt(prompt)

	s.sendEvent(&httpEvent{
		Type:   "prompt",
		Prompt: prompt,
	})

	res := <-prompt.resChan

	s.sendEvent(&httpEvent{
		Type:   "prompt",
		Prompt: prompt,
	})

	return res.Yes
}

func (s *Stack) PromptInput(msg string) string {
	prompt := &httpPrompt{
		ID:      random.Hex(16),
		Type:    "input",
		Message: msg,
		resChan: make(chan *httpPrompt),
		api:     s.api,
		stack:   s,
	}
	s.addPrompt(prompt)

	s.sendEvent(&httpEvent{
		Type:   "prompt",
		Prompt: prompt,
	})

	res := <-prompt.resChan

	s.sendEvent(&httpEvent{
		Type:   "prompt",
		Prompt: prompt,
	})

	return res.Input
}

func (s *Stack) sendEvent(event *httpEvent) {
	s.eventsMtx.Lock()
	s.events = append(s.events, event)
	s.eventsMtx.Unlock()

	for _, sub := range s.subscriptions {
		go sub.sendEvents(s)
	}
}

func (s *Stack) handleError(err error) {
	s.sendEvent(&httpEvent{
		Type:        "error",
		Description: err.Error(),
	})
}

func (s *Stack) handleDone() {
	if s.Domain != nil {
		s.sendEvent(&httpEvent{
			Type:        "domain",
			Description: s.Domain.Name,
		})
	}
	if s.DashboardLoginToken != "" {
		s.sendEvent(&httpEvent{
			Type:        "dashboard_login_token",
			Description: s.DashboardLoginToken,
		})
	}
	if s.CACert != "" {
		s.sendEvent(&httpEvent{
			Type:        "ca_cert",
			Description: base64.URLEncoding.EncodeToString([]byte(s.CACert)),
		})
	}
	s.sendEvent(&httpEvent{
		Type: "done",
	})

	for _, sub := range s.subscriptions {
		go sub.handleDone()
	}
}

func (s *Stack) handleEvents() {
	for {
		select {
		case event := <-s.EventChan:
			s.api.logger.Info(event.Description)
			s.sendEvent(&httpEvent{
				Type:        "status",
				Description: event.Description,
			})
		case err := <-s.ErrChan:
			s.api.logger.Info(err.Error())
			s.handleError(err)
		case <-s.Done:
			s.handleDone()
			s.api.logger.Info(s.DashboardLoginMsg())
			return
		}
	}
}
