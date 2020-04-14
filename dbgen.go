package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"git.sr.ht/~emersion/soju"
	"golang.org/x/crypto/bcrypt"
)

type user struct {
	name     string
	password string
	networks []network
}

func (u user) Validate() bool {
	for _, i := range []string{u.name, u.password} {
		if i == "" {
			return false
		}
	}

	if len(u.networks) == 0 {
		return false
	}

	return true
}

func userFromEnv(id int) user {
	return user{
		name:     os.Getenv(fmt.Sprintf("SOJU_USER_%d_NAME", id)),
		password: os.Getenv(fmt.Sprintf("SOJU_USER_%d_PASSWORD", id)),
	}
}

type network struct {
	name     string
	server   string
	nick     string
	password string
	channels []string
}

func (n network) Validate() bool {
	for _, i := range []string{n.name, n.server, n.nick} {
		if i == "" {
			return false
		}
	}
	return true
}

func networkFromEnv(userID, networkID int) network {
	return network{
		name:     os.Getenv(fmt.Sprintf("SOJU_USER_%d_NETWORK_%d_NAME", userID, networkID)),
		server:   os.Getenv(fmt.Sprintf("SOJU_USER_%d_NETWORK_%d_SERVER", userID, networkID)),
		nick:     os.Getenv(fmt.Sprintf("SOJU_USER_%d_NETWORK_%d_NICK", userID, networkID)),
		password: os.Getenv(fmt.Sprintf("SOJU_USER_%d_NETWORK_%d_PASSWORD", userID, networkID)),
		channels: strings.Split(os.Getenv(fmt.Sprintf("SOJU_USER_%d_NETWORK_%d_CHANNELS", userID, networkID)), ","),
	}
}

func main() {
	users := []user{}
	for i := 1; true; i++ {
		networks := []network{}
		for j := 1; true; j++ {
			n := networkFromEnv(i, j)
			if !n.Validate() {
				break
			}
			networks = append(networks, n)
		}

		u := userFromEnv(i)
		u.networks = networks
		if !u.Validate() {
			break
		}
		users = append(users, u)
	}

	db, err := soju.OpenSQLDB("sqlite3", "soju.db")
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	for _, u := range users {
		hashed, err := bcrypt.GenerateFromPassword([]byte(u.password), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("failed to hash password: %v", err)
		}

		user := soju.User{
			Username: u.name,
			Password: string(hashed),
		}
		if err := db.CreateUser(&user); err != nil {
			log.Fatalf("failed to create user: %v", err)
		}

		for _, n := range u.networks {
			network := &soju.Network{
				Name: n.name,
				Addr: n.server,
				Nick: n.nick,
				Pass: n.password,
			}

			if err := db.StoreNetwork(u.name, network); err != nil {
				log.Fatalf("failed to create network: %v", err)
			}

			for _, c := range n.channels {
				channel := &soju.Channel{
					Name: c,
				}

				if err := db.StoreChannel(network.ID, channel); err != nil {
					log.Fatalf("failed to create channel: %v", err)
				}
			}
		}
	}
}
