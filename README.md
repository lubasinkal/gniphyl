# Gniphyl

Gniphyl is a cross-platform command-line interface (CLI) tool written in **Go** that allows users to organize files within specified directories. It enables users to add, delete, list, and organize directories, ensuring that files are neatly sorted based on their type (e.g., `images`, `videos`, `documents`). 

**Built with Go's standard library only** - no external dependencies required! The tool is fast, lightweight, and produces standalone executables for **Windows**, **macOS**, and **Linux**.

---
Running the tool
![start](assets/image.png "starting from terminal")
___
Output after running the command:  **gniphyl run**

![output](assets/output.png "Organized Downloads folder")
___
## Features

- Add directory paths to a configuration file.
- Delete paths by their index.
- List all configured paths.
- Organize files into categorized folders based on their extensions (e.g., `images`, `videos`, `documents`).

---

## Installation

### Option 1: Download Precompiled Binaries

Download the precompiled executables from the [Releases](https://github.com/lubasinkal/gniphyl/releases) section:

- **For Windows**: Download `gniphyl.exe`
- **For macOS**: Download `gniphyl`
- **For Linux**: Download `gniphyl-linux`

### Option 2: Build from Source

**Prerequisites:**
- Go 1.18 or higher ([Download Go](https://go.dev/dl/))

**Build steps:**

```bash
# Clone the repository
git clone https://github.com/lubasinkal/gniphyl.git
cd gniphyl

# Build for your platform
go build -o gniphyl

# Or build for specific platforms
# Windows:
go build -o gniphyl.exe

# macOS/Linux:
go build -o gniphyl
```

### Adding to PATH

To run `gniphyl` from anywhere, add it to your system's PATH:

   To make it easier to run the `gniphyl` commands from any terminal session, move the downloaded executable to a directory that is part of your system's `PATH`.

   #### On **Windows**:
   - Move `gniphyl.exe` to a directory like `C:\Windows\System32` or another directory already in your system’s `PATH`.
   - Alternatively, add the directory where `gniphyl.exe` is located to your system’s `PATH` environment variable.

     To add a directory to the `PATH`:
     1. Right-click on **This PC** or **Computer** and select **Properties**.
     2. Select **Advanced system settings**.
     3. Click the **Environment Variables** button.
     4. In the **System variables** section, scroll to find the `Path` variable and click **Edit**.
     5. Add the path to the directory containing `gniphyl.exe` and click **OK**.

   #### On **macOS** and **Linux**:
   - Move the executable to `/usr/local/bin` (or another directory in your system’s `PATH`).
     ```bash
     sudo mv gniphyl /usr/local/bin/
     ```
     For **Linux**:
     ```bash
     sudo mv gniphyl-linux /usr/local/bin/
     ```
   - Ensure the executable is accessible by running:
     ```bash
     sudo chmod +x /usr/local/bin/gniphyl
     ```

     For **Linux**:
     ```bash
     sudo chmod +x /usr/local/bin/gniphyl-linux
     ```

---

## Usage

### Commands

```bash
# Show help
gniphyl --help

# Add a directory path to organize
gniphyl add <directory_path>

# Remove a path from configuration
gniphyl rm <directory_path>

# List all configured paths
gniphyl list

# Organize files in all configured paths
gniphyl run
```

---

## Example Workflow

```bash
# Add your Downloads folder
gniphyl add /users/name/downloads

# List configured paths
gniphyl list

# Output:
# Configured Paths:
# --------------------------------------------------
# 1. /users/name/downloads
# --------------------------------------------------

# Organize all files
gniphyl run

# Files will be sorted into folders:
# - images/     (jpg, png, gif, etc.)
# - videos/     (mp4, mkv, avi, etc.)
# - documents/  (pdf, docx, txt, etc.)
# - compressed/ (zip, rar, tar, etc.)
# - audio/      (mp3, wav, flac, etc.)
# - code/       (html, css, js, py, etc.)
# - executables/ (exe, msi)
# - others/     (unrecognized extensions)
```

## File Categories

The tool organizes files based on their extensions into the following categories:

- **images**: jpg, jpeg, png, gif, bmp, webp
- **videos**: mp4, mkv, webm, flv, avi, mov
- **documents**: pdf, doc, docx, xls, xlsx, ppt, pptx, txt, csv
- **compressed**: zip, rar, tar, gz, 7z
- **executables**: exe, msi
- **audio**: mp3, wav, flac, m4a, aac
- **code**: html, css, js, py, java, c, cpp, h, hpp, php, sql
- **others**: any other file type

You can customize these categories by modifying `config.json` in the source and rebuilding.

## Technical Details

- **Language**: Go (using only standard library)
- **Dependencies**: None (stdlib only)
- **Config Storage**: 
  - Windows: `%LOCALAPPDATA%\fileOrg\config.toml`
  - macOS/Linux: `~/.config/fileOrg/config.toml`
- **Config Format**: JSON
- **Binary Size**: ~3MB (standalone, no runtime required)

---

## Contributing

Contributions are welcome! Feel free to fork this repository, make improvements, and submit a pull request.

---

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

## Contact

For any questions or issues, please contact:
- **Name:** Lubasi Nkalolang
- **Email:** lubasinkal@outlook.com
- **GitHub:** [lubasinkal](https://github.com/lubasinkal)