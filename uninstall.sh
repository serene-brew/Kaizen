#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[0;33m'
NC='\033[0m'

# ASCII Art
echo -e "${RED}"
cat << "EOF"

          ⠠⣤⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣷⡄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⠻⣶⡄⠀⠀⣼⠟⠁⠀⠀⠀⠀⠀
⠀⣄⣀⣠⣤⣴⣶⣦⠀⢰⣿⠁⠀⠀⠀⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠠⣤⣤⣾⣷⣶⡾⠿⠾⠿⠗⠀⠀⠀⠀
⠀⠈⠛⠉⠁⠀⣿⠃⢀⣿⣷⣶⠶⣿⡿⠟⠃⠀⠀⠀⠀⠀⠀⠀⠀⢠⣬⣤⣤⣤⣿⣦⣴⣶⣶⡄⠀⠀⠀⠀
⠀⣤⣀⣀⣀⣼⣿⢀⣾⠃⠀⠀⢠⣿⠃⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣉⣉⣉⣉⣿⣥⣤⣤⣤⣄⠀⠀⠀⠀
⠀⢸⡟⠋⠉⠉⠁⠞⠁⠳⣄⢀⣾⠏⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⠛⣿⡋⠉⣿⡏⠉⣽⠟⠁⠀⠀⠀⠀
⠀⢸⡇⠀⠀⠀⣤⠀⠀⠀⢙⣿⣏⠀⠀⠀⠀⠀⠀⠀⠀⠀⢠⣄⣀⣠⣤⣼⣿⣤⣿⣷⡾⠿⠶⠾⠿⢿⣦⠀
⠀⢸⣇⠀⠀⢀⣿⡀⠀⢠⡾⠋⢻⣧⡀⠀⠀⠀⠀⠀⠀⠀⠀⠙⠛⠉⢩⣤⣤⣤⣤⣤⣶⣶⣶⣄⠀⠀⠀⠀
⠀⠈⠻⠿⠿⠿⠟⢃⣴⠟⠁⠀⠀⠹⣿⣦⣄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⣿⡉⠀⠀⠀⠀⢰⣿⠁⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠶⠋⠁⠀⠀⠀⠀⠀⠈⠛⠛⠛⠀⠀⠀⠀⠀⠀⠀⠀⠀⣿⣧⣤⣤⣤⣴⣾⡇⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠘⠉⠉
~ developed by serene-brew, 2024

EOF
echo -e "${NC}"
sleep .5

echo -e "${BLUE}[+] Starting uninstaller... ${NC}"

if [ -d "$HOME/.local/share/kaizen" ]; then
  echo -e "${GREEN}[+]removing  $HOME/.local/share/kaizen ${NC}"
  rm -rf ~/.local/share/kaizen 
else
  echo -e "${YELLOW}[!]could not locate $HOME/.local/share/kaizen ${NC}"
fi

if [ -f "/usr/bin/kaizen" ]; then
  echo -e "${GREEN}[+]removing /usr/bin/kaizen ${NC}"
  sudo rm /usr/bin/kaizen
else
  echo -e "${YELLOW}[!]could not locate /usr/bin/kaizen ${NC}"
fi

if [ -d "$HOME/.config/kaizen"  ]; then
  echo -e "${GREEN}[+]removing $HOME/.config/kaizen/config.yaml ${NC}"
  rm -rf ~/.config/kaizen
else
  echo -e "${YELLOW}[!]could not locate $HOME/.config/kaizen ${NC}"
fi

sleep 1

echo ""
echo "Thank you for trying out Kaizen."
echo "If you have any suggestions or found out any bugs, you can report it"
echo "in our github project repository."
echo "Visit & follow https://www.github.com/serene-brew for more projects."
echo "contact: serene.brew.git@gmail.com"

