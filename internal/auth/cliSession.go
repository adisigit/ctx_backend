package auth

import (
	"errors"
	"sync"
	"time"
)

type CLISessionStatus string

const (
	CLISessionPending   CLISessionStatus = "pending"
	CLISessionCompleted CLISessionStatus = "completed"
)

type CLISession struct {
	PublicKeyB64   string
	Status         CLISessionStatus
	EncryptedToken string
	ExpiresAt      time.Time
}

type CLISessionService struct {
	mu       sync.Mutex
	sessions map[string]*CLISession
}

func NewCLISessionService() *CLISessionService {
	s := &CLISessionService{sessions: make(map[string]*CLISession)}
	go s.cleanupLoop()
	return s
}

func (s *CLISessionService) Create(publicKeyB64 string) string {
	sessionID := GenerateRandomString(16)
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[sessionID] = &CLISession{
		PublicKeyB64: publicKeyB64,
		Status:       CLISessionPending,
		ExpiresAt:    time.Now().Add(10 * time.Minute),
	}
	return sessionID
}

func (s *CLISessionService) Complete(sessionID, accessToken string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	sess, ok := s.sessions[sessionID]
	if !ok || time.Now().After(sess.ExpiresAt) {
		return errors.New("session not found or expired")
	}
	encrypted, err := SealTokenForCLI(sess.PublicKeyB64, accessToken)
	if err != nil {
		return err
	}
	sess.EncryptedToken = encrypted
	sess.Status = CLISessionCompleted
	return nil
}

func (s *CLISessionService) Poll(sessionID string) (*CLISession, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	sess, ok := s.sessions[sessionID]
	if !ok || time.Now().After(sess.ExpiresAt) {
		return nil, errors.New("session not found or expired")
	}
	if sess.Status == CLISessionCompleted {
		delete(s.sessions, sessionID)
	}
	return sess, nil
}

func (s *CLISessionService) cleanupLoop() {
	ticker := time.NewTicker(time.Minute)
	for {
		select {
		case <-ticker.C:
			s.mu.Lock()
			for id, session := range s.sessions {
				if time.Now().After(session.ExpiresAt) {
					delete(s.sessions, id)
				}
			}
			s.mu.Unlock()
		}
	}
}
