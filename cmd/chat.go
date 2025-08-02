package cmd

import (
	"askai/lib/llm"
	"askai/lib/utils"
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
	"time"
)

const (
	ExitCommand = "exit"
)

var (
	chatClient  llm.ChatClient
	rateLimiter utils.RateLimiter
)

func initLLMClient() (err error) {
	client, err := llm.NewClient(
		viper.GetString("provider"),
		viper.GetString("model"),
		viper.GetString("api_key"),
	)
	if err != nil {
		return err
	}
	chatClient = client.GetChatClient()
	if chatClient == nil {
		err = errors.New("failed to create client")
		return err
	}
	return nil
}

var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Chat with ai",
	PreRun: func(cmd *cobra.Command, args []string) {
		if err := initLLMClient(); err != nil {
			cobra.CheckErr(err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		rate := viper.GetInt("rate_limit")
		if rate > 0 {
			rateLimiter = utils.NewTokenBucket(ctx, rate, rate, time.Minute)
		}

		session, err := chatClient.NewSession(ctx)
		if err != nil {
			cobra.CheckErr(err)
		}

		scanner := bufio.NewScanner(os.Stdin)
		for {
			if rateLimiter != nil && !rateLimiter.Allow() {
				fmt.Println("rate limit exceeded, please wait for a moment ...")
				allowed, err := rateLimiter.WaitWithTimeout(ctx, time.Minute)
				if err != nil {
					cobra.CheckErr(err)
				}
				if !allowed {
					continue
				}
			}

			fmt.Print("User > ")
			if !scanner.Scan() {
				cobra.CheckErr(scanner.Err())
			}
			input := strings.TrimSpace(scanner.Text())
			if strings.ToLower(input) == ExitCommand {
				cancel()
				break
			} else if input == "" {
				continue
			}

			ch, err := session.Send(input)
			if err != nil {
				cobra.CheckErr(scanner.Err())
			}
			finished := false
			for !finished {
				select {
				case <-ctx.Done():
					fmt.Println("context canceled")
					return
				case msg, ok := <-ch:
					if !ok {
						finished = true
						fmt.Println()
					} else {
						fmt.Print(msg)
						fmt.Print(" ")
					}
				}
			}
		}
	},
}
