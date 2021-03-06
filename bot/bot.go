package bot

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Noah-Huppert/kube-bot/chat"
	"github.com/Noah-Huppert/kube-bot/config"
	"github.com/Noah-Huppert/kube-bot/defs"
	"github.com/nlopes/slack"
	"github.com/prometheus/client_golang/prometheus"
)

// Bot acts as a chat bot based interface to the Kubernetes API. Leveraging the
// Slack API to send and receive messages.
type Bot struct {
	// Config holds the setting values provided by the user. Used to tweak the
	// bot's behavior.
	Config config.Config

	// ctx is used to stop the Bot's execution. The Bot will only run while the
	// context is not expired.
	ctx context.Context

	// ctxCancelFn is the cancel function for ctx
	ctxCancelFn context.CancelFunc

	// logger is used to record debug messages
	logger *log.Logger

	// slackLogger is the logger used by the Slack API client
	slackLogger *log.Logger

	// slackAPI is the authenticated Slack API client used to interact with
	// Slack
	slackAPI *slack.Client

	// slackRTM is the Slack API real time messaging client used to receive and
	// respond to Slack API events
	slackRTM *slack.RTM

	// registry holds all commands the bot can respond to
	registry chat.Registry

	// allParser is used to run a suite of Parsers on received messages
	allParser *chat.AllParser

	// ims holds the information about personal conversations (aka ims) that
	// the bot has with users.
	//
	// Keys are im channel IDs. Values are
	// TODO: Bot.ims

	// chatEventCounter
	chatEventCounter prometheus.Counter
}

// NewBot creates a new Bot instance from the parameters specified in the
// Config object. An error is return if one occurs, nil on success.
func NewBot(ctx context.Context, cfg config.Config) (*Bot, error) {
	var bot Bot

	// Config
	bot.Config = cfg

	// Loggers
	bot.logger = log.New(os.Stdout, "bot: ", 0)
	bot.slackLogger = log.New(os.Stdout, "bot: slack api: ", 0)

	// Context
	ctx, ctxCancelFn := context.WithCancel(ctx)
	bot.ctx = ctx
	bot.ctxCancelFn = ctxCancelFn

	// Registry
	bot.registry = chat.NewDefaultRegistry()
	allLdr := defs.NewAllLoader()
	if err := allLdr.Load(bot.registry); err != nil {
		return nil, fmt.Errorf("error loading registry items: %s", err.Error())
	}

	bot.chatEventCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "chat_event_count",
		Help: "Counts the number of chat events received",
	})
	if err := prometheus.Register(bot.chatEventCounter); err != nil {
		return nil, fmt.Errorf("error registering chat event counter metric: %s", err.Error())
	}

	// Make
	return &bot, nil
}

// Run begins the process of receiving and responding to user messages. This
// process will continue until Stop() is called.
//
// Returns the error that stopped execution. If Run() was stopped by a context
// either the context.Canceled or context.DeadlineExceeded error will be
// returned.
func (b Bot) Run() error {
	b.logger.Println("running bot")

	// Init Slack Lib
	b.slackAPI = slack.New(b.Config.Slack.Token)
	slack.SetLogger(b.slackLogger)

	// Connect
	b.slackRTM = b.slackAPI.NewRTM()
	go b.slackRTM.ManageConnection()

	// allParser
	b.allParser = chat.NewAllParser(b.registry, b.slackAPI, b.slackRTM)

	// Receive
	go b.handleEvents(b.slackRTM.IncomingEvents)

	select {
	case <-b.ctx.Done():
		b.logger.Printf("received shutdown request: %s", b.ctx.Err().Error())
		return b.ctx.Err()
	}

	return nil
}

// handleEvents receives Slack events via the provided channel and processes
// them accordingly. Returns the error that stopped execution.
func (b Bot) handleEvents(in <-chan slack.RTMEvent) error {
	b.logger.Println("starting to receive Slack events")

	for {
		select {
		case <-b.ctx.Done():
			// Ctx has expired
			return b.ctx.Err()
		case msg := <-in:
			// Received Slack API event
			b.chatEventCounter.Inc()
			switch event := msg.Data.(type) {
			case *slack.MessageEvent:
				if err := b.handleMessage(event); err != nil {
					b.logger.Printf("error handling message: %s\n", err.Error())
				}
			case *slack.InvalidAuthEvent:
				b.logger.Println("invalid credentials")
			case *slack.ConnectionErrorEvent:
				b.logger.Printf("connection error: %s", event.Error())
			case *slack.HelloEvent, *slack.ConnectingEvent, *slack.ConnectedEvent:
				continue
			default:
				// If logging unhandled events
				if b.Config.Slack.LogUnhandledEvents {
					b.logger.Printf("received unhandled event: %s %#s", msg.Type, event)
				}
			}
		}
	}

	return nil
}

// handleMessage performs the appropriate actions for the provided message event.
// Returns an error on failure, nil on success.
func (b Bot) handleMessage(event *slack.MessageEvent) error {
	msg := event.Msg

	// Log
	b.logger.Printf("received message: %s\n", msg.Text)

	// Test augments
	if cmdReq, err := b.allParser.Parse(msg); err == nil {
		// Format message
		str := "I'm still learning, here are your arguments:"

		for key, val := range cmdReq.Augments {
			str += fmt.Sprintf("\n- %s=%s", key, val)
		}

		// Send
		b.SendTxt(str, msg.Channel)
	} else {
		b.SendTxt(fmt.Sprintf("Whoops I had a brain fart: %s", err.Error()), msg.Channel)
	}

	return nil
}

// SendTxt uses the slackRTM client to send a text message to the specified
// channel.
func (b Bot) SendTxt(txt string, channel string) {
	b.slackRTM.SendMessage(b.slackRTM.NewOutgoingMessage(txt, channel))
}

// Stop ends the process of receiving and responding to user messages. This
// will cause the Run() method to exit and return a context.Canceled error.
func (b Bot) Stop() {
	b.ctxCancelFn()
}
