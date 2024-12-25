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

EOF
echo -e "${NC}"
sleep 2

echo -e "${BLUE}[+] Starting Kaizen installation...${NC}"
sleep 1

# Check for MPV
echo -e "${BLUE}[?] Checking for dependencies...${NC}"
if ! command -v mpv &> /dev/null; then
    echo -e "${YELLOW}[!] Error: MPV is not installed. Please install MPV first.${NC}"
    sleep 2
    exit 1
fi
echo -e "${GREEN}[✓] MPV found${NC}"
sleep 1

# Create necessary directories
echo -e "${BLUE}[+] Creating directories...${NC}"
mkdir -p ~/.local/share/kaizen
mkdir -p ~/.config/kaizen
echo -e "${GREEN}[✓] Directories created${NC}"
sleep 1

# Copy and set permissions for maintenance scripts
echo -e "${BLUE}[+] Setting up maintenance scripts...${NC}"
cp update.sh uninstall.sh ~/.local/share/kaizen/
chmod +x ~/.local/share/kaizen/update.sh
chmod +x ~/.local/share/kaizen/uninstall.sh
echo -e "${GREEN}[✓] Maintenance scripts configured${NC}"
sleep 1

# Build using make
echo -e "${BLUE}[+] Building Kaizen...${NC}"
if ! make; then
    echo -e "${YELLOW}[!] Error: Build failed${NC}"
    exit 1
fi
echo -e "${GREEN}[✓] Build successful${NC}"
sleep 1

# Copy files
echo -e "${BLUE}[+] Installing Kaizen...${NC}"
sudo cp build/kaizen /usr/bin/
cp config.yaml ~/.config/kaizen/
echo -e "${GREEN}[✓] Kaizen installed${NC}"
sleep 1

# Copy version file
cp VERSION ~/.local/share/kaizen/
sleep 1

# Clean up build directory
echo -e "${BLUE}[+] Cleaning up...${NC}"
make clean
echo -e "${GREEN}[✓] Build directory cleaned${NC}"
sleep 1

echo -e "${GREEN}[+] Installation complete! You can now run Kaizen by typing 'kaizen' in your terminal.${NC}"
