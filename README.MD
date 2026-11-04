# ASCII Art Web

## Description

ASCII Art Web is a small web server written in Go that turns text into ASCII art. The person types some text, picks one of three banner styles (standard, shadow or thinkertoy), and presses Submit. The drawing appears in a text area on the same page. It is the web version of the earlier ascii-art command-line project.

## Authors

- Muhammed
- Soliman

## Usage

1. Make sure Go is installed.
2. From the project folder, run:

   ```
   go run .
   ```

3. Open `http://localhost:8080` in a browser.
4. Type your text, choose a banner, and press Submit. Type `\n` to start a new line.

The project folder must contain `main.go`, the three banner files (`standard.txt`, `shadow.txt`, `thinkertoy.txt`) and a `templates` folder with `index.html`.

## Implementation details

### Endpoints

- `GET /` sends the main page, with an empty result box. Any other path answers 404.
- `POST /ascii-art` receives the text and banner from the form, draws the art, and sends the same page back with the result in the text area. Any other method answers 400.

### Status codes

- 200 when everything works. Empty text also gives 200 with an empty result.
- 400 for a bad request: an unknown banner, a character outside the printable range (32 to 126), or a wrong method on `/ascii-art`.
- 404 when the template, a banner file, or an unknown page cannot be found.
- Template execution errors are logged on the server.

### Algorithm

1. The handler reads `theInput` and `banner` from the form and validates them.
2. The banner file is read and split into a list of lines.
3. Each banner file stores 95 characters (ASCII 32 to 126), and each character takes 9 lines: a blank separator line, then 8 lines of art. The first line of a character is found with `(character - 32) * 9 + 1`.
4. The text is split at every `\n` into separate words. Text made only of `\n` gives one blank line fewer than the number of pieces, and an empty piece gives one blank line.
5. For each word, the art is built row by row: for each of the 8 rows, the matching line of every character is added side by side, followed by a newline. Rows are drawn on the outside and characters on the inside, because the characters sit next to each other.
6. The finished string is passed to the template with `html/template`, which places it in the result text area.

A few lines of JavaScript in the page shrink the font when the art is wider than the box, so the drawing always fits.