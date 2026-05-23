// Command sshsite runs an SSH server that serves a styled terminal landing
// page. Anyone can connect with `ssh -p 2222 localhost` (no password) and see
// the site rendered right in their terminal.
package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
	zone "github.com/lrstanley/bubblezone"
	"github.com/muesli/termenv"
)

const (
	host = "0.0.0.0"
	port = "2222"
)

func main() {
	// The server runs without a TTY, so lipgloss's default renderer would
	// detect "no color" and strip every style. Force a color profile so the
	// escape codes are emitted and rendered by each client's terminal.
	lipgloss.SetColorProfile(termenv.ANSI256)

	zone.NewGlobal() // enable clickable mouse zones

	srv, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		wish.WithMiddleware(
			bubbletea.Middleware(teaHandler),
			activeterm.Middleware(), // only allow interactive terminals
			logging.Middleware(),
		),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, "could not create server:", err)
		os.Exit(1)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	fmt.Printf("Starting SSH server on %s:%s\n", host, port)
	fmt.Printf("Connect with:  ssh -p %s localhost\n", port)
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			fmt.Fprintln(os.Stderr, "server error:", err)
			done <- syscall.SIGTERM
		}
	}()

	<-done
	fmt.Println("\nStopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		fmt.Fprintln(os.Stderr, "shutdown error:", err)
		os.Exit(1)
	}
}

// teaHandler builds a Bubble Tea program for each SSH session.
func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	pty, _, active := s.Pty()
	if !active {
		wish.Fatalln(s, "no active terminal, skipping")
		return nil, nil
	}
	m := model{
		term:   pty.Term,
		width:  pty.Window.Width,
		height: pty.Window.Height,
		user:   s.User(),
		page:   pageHome, // start on the landing screen, not page 0 (Creations)
	}
	return m, []tea.ProgramOption{tea.WithAltScreen(), tea.WithMouseCellMotion()}
}
