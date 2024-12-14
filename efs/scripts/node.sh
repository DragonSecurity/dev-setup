curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.1/install.sh | bash
omz reload
nvm install --lts
nvm use --lts --default
npm install -g semantic-release
npm install -g yarn