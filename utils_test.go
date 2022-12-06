package main

import (
	"testing"
)

func Test_getBowlsData(t *testing.T) {
	got, err := getBowlsData()
	if err != nil {
		t.Errorf("Error getting bowl data: %v", err)
	}

	if len(got) == 0 {
		t.Errorf("No bowl data found")
	}
}

func Test_getTeamData(t *testing.T) {
	got, err := getTeamData()
	if err != nil {
		t.Errorf("Error getting team data: %v", err)
	}

	if len(got) == 0 {
		t.Errorf("No team data found")
	}
}

func Test_no_teams(t *testing.T) {
	noTeam := 0
	noTeamsArr := make([]string, noTeam)

	bg, err := getBowlsData()
	if err != nil {
		t.Fatal(err)
	}

	for _, bowlGame := range bg {
		if bowlGame.Team1.ID == noTeam {
			noTeamsArr = append(noTeamsArr, bowlGame.Name)
		} else if bowlGame.Team2.ID == noTeam {
			noTeamsArr = append(noTeamsArr, bowlGame.Name)
		}
	}

	if len(noTeamsArr) != 0 {
		t.Log("The following bowl games have no teams:")
		for _, name := range noTeamsArr {
			t.Log(name)
		}
		t.Fail()
	}
}
