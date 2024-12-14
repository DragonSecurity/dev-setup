#!/usr/bin/env bash
set -ex

function runGit() {
  sudo -v

  # Keep-alive: update existing `sudo` time stamp until the script has finished
  while true; do sudo -n true; sleep 60; kill -0 "$$" || exit; done 2>/dev/null &

  OS=$(uname -s | tr A-Z a-z)

  case $OS in
    linux)
      source /etc/os-release
      case $ID in
        debian|ubuntu|mint|linuxmint)
          sudo apt update
          sudo apt install git -y
          ;;
        fedora|rhel|centos)
          sudo yum update
          sudo yum install git -y
          ;;
        *)
          echo -n $ID
          echo -n "unsupported linux distro"
          ;;
        esac
      ;;

    darwin)
      brew update
      brew install git
      ;;

    *)
      echo -n "unsupported linux distro"
      ;;
    esac
}

echo -n "This script may overwrite existing files in your home directory. Are you sure? (y/n) " -n 1;
runGit