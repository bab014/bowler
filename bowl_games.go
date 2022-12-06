package main

import (
	"encoding/json"
	"net/url"
	"os"
	"strings"
)

// bowlGame struct holds the game data
type bowlGame struct {
	Name       string `json:"name,omitempty"`
	Team1      string `json:"team1,omitempty"`
	Team2      string `json:"team2,omitempty"`
	Team1Score int    `json:"team1Score,omitempty"`
	Team2Score int    `json:"team2Score,omitempty"`
	Date       string `json:"date,omitempty"`
	Time       string `json:"time,omitempty"`
}

// BowGame holds the game data with the proper team type
type BowlGame struct {
	Name       string
	Team1      Team   `json:"team_1,omitempty"`
	Team2      Team   `json:"team_2,omitempty"`
	Team1Score int    `json:"team_1_score,omitempty"`
	Team2Score int    `json:"team_2_score,omitempty"`
	Date       string `json:"date,omitempty"`
	Time       string `json:"time,omitempty"`
}

type Team struct {
	ID           int      `json:"id"`
	School       string   `json:"school"`
	Mascot       string   `json:"mascot"`
	Abbreviation string   `json:"abbreviation"`
	AltName1     string   `json:"alt_name_1"`
	AltName2     string   `json:"alt_name_2"`
	AltName3     string   `json:"alt_name_3"`
	Conference   string   `json:"conference"`
	Division     string   `json:"division"`
	Color        string   `json:"color"`
	AltColor     string   `json:"alt_color"`
	Logos        []string `json:"logos"`
	Twitter      string   `json:"twitter"`
	Location     Location `json:"location"`
}

type Location struct {
	VenueID         int     `json:"venue_id"`
	Name            string  `json:"name"`
	City            string  `json:"city"`
	State           string  `json:"state"`
	Zip             string  `json:"zip"`
	CountryCode     string  `json:"country_code"`
	TimeZone        string  `json:"time_zone"`
	Latitude        float64 `json:"latitude"`
	Longitude       float64 `json:"longitude"`
	Elevation       string  `json:"elevation"`
	Capacity        int     `json:"capacity"`
	YearConstructed int     `json:"year_constructed"`
	Grass           bool    `json:"grass"`
	Dome            bool    `json:"dome"`
}

// Selection is the game and winner selection from a user
type Selection struct {
	GameName       string `json:"game_name,omitempty"`
	SelectedWinner string `json:"selected_winner,omitempty"`
}

// Selections is a completed selections from a user
type Selections []Selection

func NewSelctions(v url.Values) Selections {
	var s Selections
	for k, v := range v {
		if k == "submitter" {
			continue
		}
		s = append(s, Selection{
			GameName:       k,
			SelectedWinner: v[0],
		})
	}
	return s
}

func (s Selections) MakeFile(name string) error {
	var err error

	name = strings.ToLower(name)

	f, err := os.OpenFile(SELECTIONS_DIR+name+".json", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	if err = json.NewEncoder(f).Encode(s); err != nil {
		return err
	}

	return err
}

func SelectionsFromFile(name string) (Selections, error) {
	var s Selections
	name = strings.ToLower(name)

	// TODO: check if file exists first
	f, err := os.Open(SELECTIONS_DIR + name + ".json")
	if err != nil {
		return s, err
	}
	defer f.Close()

	if err = json.NewDecoder(f).Decode(&s); err != nil {
		return s, err
	}

	return s, err
}
