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
  backspace, ctrl+b
```
#### Quit NvFile

```http
  ctrl+q, ctrl+z
```
## Development Roadmap

This project is being actively developed. Currently features on the roadmap -

- More optimization in bubble tea tui interface for dynamic window sizing. 
- A json configuration file to define keyboard shortcut and theme colors
- Integration with other terminal vased text/code editor such as nano and gedit through the aforementioned config file.

## Be sure to Checkout 

 - [Lipgloss](https://pkg.go.dev/github.com/charmbracelet/lipgloss)
 - [BubbleTea](https://pkg.go.dev/github.com/charmbracelet/bubbletea)
 - [Bubbles](https://pkg.go.dev/github.com/charmbracelet/bubbles)

