<div align="center">

# Kaizen v1.1.0
</div>
<div align="center">

<img src="./logo/漸_2_-removebg-preview.png">
</div>
<div align="center">

[![stars](https://img.shields.io/github/stars/serene-brew/kaizen?style=social)](https://github.com/serene-brew/kaizen/stargazers)
[![forks](https://img.shields.io/github/forks/serene-brew/kaizen?style=social)](https://github.com/serene-brew/kaizen/network/members)
[![issues](https://img.shields.io/github/issues/serene-brew/kaizen?style=social&logo=github)](https://github.com/serene-brew/kaizen/issues?q=is%3Aissue+is%3Aopen+)
[![version](https://img.shields.io/github/v/release/serene-brew/kaizen?display_name=release&style=social&logo=github)](https://github.com/serene-brew/kaizen/releases/latest)
[![go_reportcard](https://img.shields.io/badge/go%20report-A+-brightgreen.svg?style=social&logo=github)](https://goreportcard.com/report/github.com/serene-brew/Kaizen)
</div>

## Why Watch Anime Like a Normie?
Ever tried watching anime from the terminal? 🖥️✨ Yeah, me neither. But one day, while staring into the abyss of CLI tools and minimalism-induced existential dread, we thought: Why not? 🤷‍♂️

Kaizen is here to fulfill that dream (or nightmare, depending on how much you love the terminal). A native, no-nonsense, super-aesthetic TUI for streaming your favorite anime—right from your beloved command line. 🎥🎌

No more context-switching, no more bloated browser tabs screaming at your CPU like it owes them money. With Kaizen, it’s just you, your terminal, and an uninterrupted anime binge-fest.

So, grab your instant ramen 🍜, fire up Kaizen, and let’s go full weeb mode. 🌸✨ Terminal supremacy, but make it anime. 🎉


https://github.com/user-attachments/assets/ed43bfb2-1338-4d63-86a1-1abb435382c5

## Installation, Update & Uninstallation

1. Clone the Repository
```bash
git clone https://github.com/serene-brew/Kaizen
```

2. Install it
```bash
chmod +x ./install.sh

./install.sh
```

### View Version Info

```bash
kaizen -v
```

### Update and Uninstallation

To update

```bash
kaizen -update
```

To uninstall
```bash
kaizen -uninstall
```
## Configuration
Once Kaizen is installed, you can find `config.yaml` file in your `~/.config/kaizen` directory. If it is not present, then copy the `config.yaml` into `~/.config/kaizen` directory. That file contains some default colors for kaizen, however you can modify it according to your own needs. 

<h2 align="center"> Keybinds </h2>

<table align=center>

<tr><th> In app </th><th> In player </th></tr>
<tr><td>
  
| Operation       | Keybind       |
|-----------------|----------------|
| Toggle b/w tabs |     `tab`      |
| Focus on Input  | `shift` + `1`  |
| Focus on Results| `shift` + `2`  |
| Focus on SUB ep | `shift` + `3`  |
| Focus on DUB ep | `shift` + `4`  |
| Exit the app    |     `esc`      |

</td><td>

| Operation       | Keybind       |
|-----------------|----------------|
| Pause and play  |    `space`     |
| Forward & Backward  | `→` & `←`  |
| Volume Control  |    `9` & `0`   |
| Mute Audio      |           `m`  |
| Adjust Brightness | `3` & `4`  |
| Exit the player |     `q`        |

</td></tr> 
</table>

> [!NOTE]
> <h2> Known issues </h2>
> <p>This is a TUI which requires a very specific kind of display size, so either fire up Kaizen in full screen terminal or decrease the size of terminal like show below if you have a very small resolution monitor.</p>
![kaizen_issue_solved](https://github.com/user-attachments/assets/9bad59ea-fb5c-4b00-97f5-17df9e2cd142)


<br/>

> [!IMPORTANT]
> <h2> Info about the API and the Backend </h2>
> <p>The API is Home grown by Serene Brew and due to the current circumstances of unofficial APIs being taken down, we decided to keep the codebase confidential, however the API endpoints are still being used in Kaizen</p>




<h2 align="center"> Contributors </h2>
<div align="center">
  
<a href="https://github.com/mintRaven-05"><img src="https://github.com/user-attachments/assets/2a8d88d5-58a7-43ed-ab2e-d8896c272904" height=70 width=70)></a>
<a href="https://github.com/ImonChakraborty"><img src="https://github.com/user-attachments/assets/50c4e290-6ec7-4390-bc76-51c0e9ea1cae" height=70 width=70></a>
<a href="https://github.com/noornee"><img src="https://github.com/user-attachments/assets/d60edfc5-873e-4fdb-93bd-3b289ddd1693" height=70 width=70></a>

</div>
