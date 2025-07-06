# NvFile

A File explorer for coding needs written in go with the bubbletea framework with a modern and intuitive interface.

NvFile is a TUI file explorer that works with your terminal text/code editors like- Nvim, vim, Nano etc. This works as an alternative to the default file explorer window in nvim though custom configuration for integration with text editors other than Nvim is in the works. NvFile has most of the common shortcuts for new file/ new subdirectory creation, deletion, viewing hidden files. 


## Demo

Deleting file or subdirectories- 
![ezgif com-video-to-gif-converter](https://github.com/user-attachments/assets/27f2faa0-d109-473a-bb25-a07347d89b4e)


## Useful Keyboard Shortcuts

#### Going back in the file tree

```http
  ctrl+b
```

#### Create new file

```http
  ctrl+n
```

#### Create new subdirectory

```http
  ctrl+d
```

#### Delete subdirectory or File

```http
  delete
```

#### Toggle between showing and not showing hidden files

```http
  ctrl+h
```

#### Opening a file in your editor of choice

```http
  enter
```
#### Entering a subdirectory 

```http
  enter
```
#### Quit creating a new file or subdirectory

```http
  ctrl+b
```
#### Quit NvFile

```http
  ctrl+q, ctrl+z
```
## Development Roadmap

This projct has reached v1.0. All the necessary features have been implementated. 
Some more useful optimizations that can be made are still in the works. Though  in the current state the programme is usable in any linux system with neovim. 

Some smaller features on the development roadmap- 
1. A seperate preview window for files and directories.

## Be sure to Checkout 

 - [Lipgloss](https://pkg.go.dev/github.com/charmbracelet/lipgloss)
 - [BubbleTea](https://pkg.go.dev/github.com/charmbracelet/bubbletea)
 - [Bubbles](https://pkg.go.dev/github.com/charmbracelet/bubbles)

