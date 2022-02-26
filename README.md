# [ Learning ] Keyword File Scanner

This application scan through all the files in directory to search for a specific keyword.

## How to use

Create a `.folderignore` file to include the file or folder you would like to exclude

Run the following command to perform the scan,
`go run main.go --keyword="TODO" --absolute`

## Argument

| Argument | Description                                        |          | Example                   |
| -------- | -------------------------------------------------- | -------- | ------------------------- |
| keyword  | Keyword to search through entire project structure | Require  | `--keyword="console.log"` |
| absolute | Show absolute path of the system                   | Optional | `--absolute`       |
