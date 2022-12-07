package main

import (
	"encoding/json"
	"errors"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
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

// teamsFromFile reads the data/teams2022.json file and returns a slice of Team
func teamsFromFile() ([]Team, error) {
	var teams []Team
	f, err := os.Open(TEAM_FILE)
	if err != nil {
		return teams, err
	}
	defer f.Close()

	if err = json.NewDecoder(f).Decode(&teams); err != nil {
		return teams, err
	}

	return teams, err
}

// Selection is the game and winner selection from a user
type Selection struct {
	GameName       string   `json:"game_name,omitempty"`
	SelectedWinner Team     `json:"selected_winner,omitempty"`
	GameInfo       BowlGame `json:"game_info,omitempty"`
}

// Selections is a completed selections from a user
type Selections []Selection

func NewSelections(v url.Values) Selections {
	tf, err := teamsFromFile()
	if err != nil {
		return nil
	}
	bgames, err := getBowlsData()
	if err != nil {
		return nil
	}
	var s Selections
	for k, v := range v {
		if k == "submitter" {
			continue
		}
		for _, bg := range bgames {
			if bg.Name == k {
				for _, t := range tf {
					if t.School == v[0] {
						s = append(s, Selection{
							GameName:       k,
							SelectedWinner: t,
							GameInfo:       bg,
						})
					}
				}
			}
		}
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
	f, err := os.Open(SELECTIONS_DIR + name)
	if err != nil {
		return s, err
	}
	defer f.Close()

	if err = json.NewDecoder(f).Decode(&s); err != nil {
		return s, err
	}

	return s, err
}

// UserSelections is a map of user selections
type UserSelections map[string]Selections

func NewUserSelections() (UserSelections, error) {
	us := make(UserSelections)

	// check to make sure selections directory is not empty
	t, err := SelectionsDirIsEmpty()
	if err != nil {
		return us, errors.New("no selections")
	}
	if t {
		return us, errors.New("selections directory is empty")
	}

	// walk the data/selections directory and create the map of user selections
	err = filepath.WalkDir(SELECTIONS_DIR, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && d.Name() == "selections" {
			return nil
		}

		// get the selections from the file
		s, err := SelectionsFromFile(d.Name())
		if err != nil {
			return err
		}

		bgames, err := getBowlsData()
		if err != nil {
			return err
		}
		for idx, game := range bgames {
			for idx2, sel := range s {
				if sel.GameName == game.Name {
					s[idx], s[idx2] = sel, s[idx]
				}
			}

		}

		// use the filename as the key but capitalize the first letter
		iName := d.Name()
		il := strings.ToTitle(string(iName[0]))
		bn := iName[1:]
		name := il + strings.TrimSuffix(bn, filepath.Ext(bn))

		// add the selections to the map with the name key
		us[name] = s

		return nil
	})
	if err != nil {
		return us, err
	}

	return us, nil
}
