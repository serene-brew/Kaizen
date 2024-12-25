#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[0;33m'
NC='\033[0m'


if [[ "$(grep "VERSION" ~/.local/share/kaizen/VERSION | cut -c 10-)" == "$(curl -s https://api.github.com/repos/serene-brew/Kaizen/releases/latest | grep tag_name | cut -c 16-21)" ]]; then
  echo -e "${GREEN}[+] no updates released !${NC}"
else
    echo -e "${YELLOW}[!]A new version of kaizen is released${NC}"
    echo -e "${GREEN}[+]Downloading latest update...${NC}"
        
    if [ -d "$HOME/Kaizen" ]; then
      rm -rf ~/Kaizen
      git clone https://github.com/serene-brew/Kaizen.git ~/Kaizen

      sleep .3
      echo -e "${GREEN}[+]latest release downloaded successfully !${NC}"
      echo -e "${BLUE}[+]removing current versioni${NC}"
      sleep .3
            
      kaizen -uninstall
    
      sleep 2
      clear
      sleep .3
            
      cd ~/Kaizen || exit

      chmod +x ./install.sh
      ./install.sh
    else
      git clone https://github.com/serene-brew/Kaizen.git ~/Kaizen
      sleep .3
            
      echo -e "${GREEN}[+]latest release downloaded successfully !${NC}"
      echo -e "${BLUE}[+]removing current versioni${NC}"
            
      sleep .3

      ~/.local/share/kaizen/uninstall.sh

      sleep 2
      clear
      sleep .3


      cd ~/Kaizen || exit
      chmod +x ./install.sh
      ./install.sh
    fi

fi

