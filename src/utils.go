package src

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

/*
generateRows is a method of Tab1Model that converts a two-dimensional slice of interface{} data into a slice of table.Row.
Each row in the table is constructed by extracting and formatting specific fields from the input data.
The rows are intended to be displayed in a tabular Bubble Tea model.
*/
func (m *Tab1Model) generateRows(data [][]any) []table.Row {
	m.loading = false
	rows := []table.Row{}
	for i, item := range data {
		centerText := func(text string, width int) string {
			if len(text) >= width {
				return text
			}
			paddingTotal := width - len(text)
			leftPadding := paddingTotal / 2
			rightPadding := paddingTotal - leftPadding
			return strings.Repeat(" ", leftPadding) + text + strings.Repeat(" ", rightPadding)
		}

		id := strconv.Itoa(i + 1)
		// Safely extract fields with type tolerance since the API may return
		// numbers as float64 or strings.
		title, _ := item[1].(string)

		// helper to coerce to int from float64 or string
		toInt := func(v any) int {
			switch t := v.(type) {
			case float64:
				return int(t)
			case int:
				return t
			case string:
				if t == "" {
					return 0
				}
				if n, err := strconv.Atoi(t); err == nil {
					return n
				}
				if f, err := strconv.ParseFloat(t, 64); err == nil {
					return int(f)
				}
				return 0
			default:
				return 0
			}
		}

		subEpisodes := strconv.Itoa(toInt(item[2]))
		dubEpisodes := strconv.Itoa(toInt(item[3]))

		// rating and status may be strings or other types; safely stringify
		rating := "-:-"
		if v := item[11]; v != nil {
			if s, ok := v.(string); ok && s != "" {
				rating = s
			}
		}

		status := "-:-"
		if v := item[9]; v != nil {
			if s, ok := v.(string); ok && s != "" {
				status = s
			}
		}

		scoreText := "N/A"
		switch s := item[12].(type) {
		case float64:
			scoreText = fmt.Sprintf("%.1f", s)
		case string:
			if s != "" {
				if f, err := strconv.ParseFloat(s, 64); err == nil {
					scoreText = fmt.Sprintf("%.1f", f)
				}
			}
		}

		rows = append(rows, table.Row{
			centerText(id, 10),
			title,
			centerText(subEpisodes, 20),
			centerText(dubEpisodes, 20),
			centerText(scoreText, 15),
			centerText(rating, 15),
			centerText(status, 20),
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
		streamTitle := fmt.Sprintf("--force-media-title=%s Episode %s (SUB)", m.animeName, m.subSelectedNum)
		stream := exec.Command("mpv", "-fs", m.streamLink, streamTitle)
		stream.Output() //nolint:errcheck
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
	DubEpisodeString := fmt.Sprintf("%s", m.listTwo.SelectedItem()) //nolint:govet // Ignore the warning for this line
	DubEpisodeString = DubEpisodeString[8:12]
	DubEpisodeString = strings.ReplaceAll(DubEpisodeString, " ", "")
	m.dubSelectedNum = DubEpisodeString
	m.episodeType = "dub"
	link, _ := getStreamLink(m.animeID, m.episodeType, m.dubSelectedNum)
	m.streamLink = link
	if m.streamLink != "" {
		streamTitle := fmt.Sprintf("--force-media-title=%s Episode %s (DUB)", m.animeName, m.dubSelectedNum)
		stream := exec.Command("mpv", "-fs", m.streamLink, streamTitle)
		stream.Output() //nolint:errcheck
	} else {
		fmt.Println("no link found")
	}
}

/*
generateSubEpisodes is a method of Tab1Model that generates a list of items representing subbed episodes.
It creates a list of episodes from 1 to the given number, formatted with default styles.
*/
func (m *Tab1Model) generateSubEpisodes() []list.Item {
	availableSubEpisodes := m.availableSubEpisodes
	items := []list.Item{}
	for _, episode := range availableSubEpisodes {
		items = append(items, item{title: "Episode " + episode + "               ", style: "default"})
	}
	return items
}

/*
generateDubEpisodes is a method of Tab1Model that generates a list of items representing dubbed episodes.
Like generateSubEpisodes, it creates a list of episodes from 1 to the given number, formatted with default styles.
*/
func (m *Tab1Model) generateDubEpisodes() []list.Item {
	availableDubEpisodes := m.availableDubEpisodes
	items := []list.Item{}
	for _, episode := range availableDubEpisodes {
		items = append(items, item{title: "Episode " + episode + "               ", style: "default"})
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
