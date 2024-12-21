package main
import (
	"fmt"
	"strconv"
	"strings"
	"os/exec"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

/*
generateRows is a method of Tab1Model that converts a two-dimensional slice of interface{} data into a slice of table.Row.
Each row in the table is constructed by extracting and formatting specific fields from the input data.
The rows are intended to be displayed in a tabular Bubble Tea model.
*/
func (m *Tab1Model) generateRows(data [][]interface{}) []table.Row {
	m.loading = false
	rows := []table.Row{}
	for i, item := range data {
		rows = append(rows, table.Row{
			strconv.Itoa(i+1),
			item[1].(string),
			strconv.Itoa(int(item[2].(float64))),
			strconv.Itoa(int(item[3].(float64))),
		})
	}
	return rows
}

/*
streamSubAnime is a method of Tab1Model that streams a selected subbed anime episode.
It extracts the selected episode number, determines the streaming link using getStreamLink,
and invokes the MPV media player to play the episode in full-screen mode.
*/
func (m *Tab1Model) streamSubAnime() {
	SubEpisodeString := fmt.Sprintf("%s", m.listOne.SelectedItem()) //nolint:govet // Ignore the warning for this line
	SubEpisodeString = SubEpisodeString[8:12]
	SubEpisodeString = strings.ReplaceAll(SubEpisodeString, " ", "")
	m.subSelectedNum = SubEpisodeString
	m.episodeType = "sub"
	link, _ := getStreamLink(m.animeID, m.episodeType, m.subSelectedNum)
	m.streamLink = link
	if m.streamLink != "" {
		stream := exec.Command("mpv", "-fs", m.streamLink)
		stream.Output()
	} else {
		fmt.Println("no link found")
	}
}

/*
streamDubAnime is a method of Tab1Model that streams a selected dubbed anime episode.
It operates similarly to streamSubAnime, but it sets the episode type to "dub"
and fetches the streaming link accordingly before playing the episode with MPV.
*/
func (m *Tab1Model) streamDubAnime() {
	DubEpisodeString := fmt.Sprintf("%s",m.listTwo.SelectedItem()) //nolint:govet // Ignore the warning for this line
	DubEpisodeString = DubEpisodeString[8:12]
	DubEpisodeString = strings.ReplaceAll(DubEpisodeString, " ", "")
	m.dubSelectedNum = DubEpisodeString
	m.episodeType = "dub"
	link, _ := getStreamLink(m.animeID, m.episodeType, m.dubSelectedNum)
	m.streamLink = link
	if m.streamLink != "" {
		stream := exec.Command("mpv", "-fs", m.streamLink)
		stream.Output()
	} else {
		fmt.Println("no link found")
	}
}

/*
generateSubEpisodes is a method of Tab1Model that generates a list of items representing subbed episodes.
It creates a list of episodes from 1 to the given number, formatted with default styles.
*/
func (m *Tab1Model) generateSubEpisodes(num int) []list.Item {
	items := []list.Item{}
	for i := 1; i <= num; i++ {
		items = append(items, item{title: "Episode " + strconv.Itoa(i) + "               ", style: "default"})
	}
	return items
}

/*
generateDubEpisodes is a method of Tab1Model that generates a list of items representing dubbed episodes.
Like generateSubEpisodes, it creates a list of episodes from 1 to the given number, formatted with default styles.
*/
func (m *Tab1Model) generateDubEpisodes(num int) []list.Item {
	items := []list.Item{}
	for i := 1; i <= num; i++ {
		items = append(items, item{title: "Episode " + strconv.Itoa(i) + "               ", style: "default"})
	}
	return items
}

/*
fetchAnimeData is a method of Tab1Model that retrieves anime data based on a given query.
It returns a Bubble Tea command (tea.Cmd) that fetches the data asynchronously.
If an error occurs during data retrieval, it is returned as the command's message.
*/
func (m *Tab1Model) fetchAnimeData(query string) tea.Cmd {
	return func() tea.Msg {
		data, err := extractInfo(query)
		if err != nil {
			return err
		}
		return data
	}
}

