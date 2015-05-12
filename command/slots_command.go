package command

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type slotsCommand struct {
	cp           *CommandPool
	uses         map[string]time.Time
	currencyName string
}

func (T slotsCommand) Init() {

}

func (T slotsCommand) ID() string {
	return "moneygames"
}

func (T slotsCommand) Response(username, message string) string {
	if strings.HasPrefix(strings.ToLower(message), "!slots") {
		lastUse, ok := T.uses[username]
		if ok && time.Since(lastUse) < 15*time.Second {
			return ""
		}

		user, in := T.cp.channel.InChannel(username)
		if !in {
			return ""
		}

		cost := 1

		space := strings.Index(strings.TrimSpace(message), " ")
		if space > 0 {
			log.Printf("Args: %s", message[space+1:])
			cost, _ = strconv.Atoi(message[space+1:])
			if cost <= 0 {
				cost = 1
			}
		}

		if user.Money < cost {
			return fmt.Sprintf("@%s You don't have enough money.", username)
		}

		T.uses[username] = time.Now()

		x := rand.Float64()
		winnings := 0

		if x < 0.01 {
			winnings = cost * 100
		} else if x < 0.05 {
			winnings = cost * 20
		} else if x < 0.1 {
			winnings = cost * 10
		} else if x < 0.3 {
			winnings = cost * 5
		} else if x < 0.43 {
			winnings = cost * 2
		}

		user.Money = user.Money - cost + winnings
		if winnings > cost {
			return fmt.Sprintf("@%s You won %d %ss", username, winnings-cost, T.cp.currencyName)
		} else {
			return fmt.Sprintf("@%s You lost %d %ss", username, cost-winnings, T.cp.currencyName)
		}
	}

	return ""
}

func (T slotsCommand) String() string {
	return ""
}
