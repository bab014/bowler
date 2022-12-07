package main

import (
	"encoding/json"
	"io"
	"os"
)

// getTeamData function opens the team data file and returns the team data
func getTeamData() ([]Team, error) {
	// Open the team data file
	file, err := os.Open(TEAM_FILE)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Decode the team data
	var teams []Team
	err = json.NewDecoder(file).Decode(&teams)
	if err != nil {
		return nil, err
	}

	return teams, nil
}

// getTeam function returns the team data for the given team name (school)
func convGame(game bowlGame, teams []Team) BowlGame {
	teamsArr := make([]Team, 0, 2)
	var bowlGame BowlGame

	for _, team := range teams {
		if game.Team1 == team.School {
			teamsArr = append(teamsArr, team)
			bowlGame.Team1 = team
		}
		if game.Team2 == team.School {
			teamsArr = append(teamsArr, team)
			bowlGame.Team2 = team
		}
		if len(teamsArr) == 2 {
			bowlGame.Name = game.Name
			bowlGame.Date = game.Date
			bowlGame.Time = game.Time
			bowlGame.Team1Score = game.Team1Score
			bowlGame.Team2Score = game.Team2Score
			break
		}
	}
	return bowlGame
}

// createBowlGames function creates the bowl game data from provided Teams data and bowl game data
func createBowlGames(teams []Team, bowlGames []bowlGame) []BowlGame {
	var games []BowlGame
	for _, game := range bowlGames {
		games = append(games, convGame(game, teams))
	}
	return games
}

// getBowlsData function returns the bowl game data
func getBowlsData() ([]BowlGame, error) {
	// Open the bowl game file
	file, err := os.Open(BOWL_FILE)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Decode the bowl game data
	var bowlGames []bowlGame
	err = json.NewDecoder(file).Decode(&bowlGames)
	if err != nil {
		return nil, err
	}

	teams, err := getTeamData()
	if err != nil {
		return nil, err
	}

	// var bowlGames2 []BowlGame

	bGames := createBowlGames(teams, bowlGames)

	return bGames, nil
}

func SelectionsDirIsEmpty() (bool, error) {
	f, err := os.Open(SELECTIONS_DIR)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err
}
