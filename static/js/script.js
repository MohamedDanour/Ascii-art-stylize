var form = document.getElementById("art-form");
var input = document.getElementById("input");
var counter = document.getElementById("counter");
var errorBox = document.getElementById("error");
var drawBtn = document.getElementById("draw");
var downloadBtn = document.getElementById("download");
var out = document.getElementById("output"); // null until there is a result
var copyBtn = document.getElementById("copy"); // null until there is a result

// Shrink the art until it fits the box, then grow the box to show every row.
function fit() {
    if (!out) return;
    var size = 14;
    out.style.fontSize = size + "px";

    while (out.scrollWidth > out.clientWidth && size > 4) {
        size--;
        out.style.fontSize = size + "px";
    }

    out.style.height = "auto";
    out.style.height = out.scrollHeight + "px";
}

// Count characters and lines, warn about characters we can't draw,
// and turn the buttons off when there is nothing valid to send.
function update(clearServerError) {
    var text = input.value;
    var chars = text.replace(/\n/g, "").length;
    var lines = text === "" ? 0 : text.split("\n").length;
    counter.textContent =
        chars + (chars === 1 ? " character" : " characters") + " · " +
        lines + (lines === 1 ? " line" : " lines");

    var bad = text.match(/[^\x20-\x7E\n]/);
    if (bad) {
        errorBox.textContent = '"' + bad[0] + '" can\'t be drawn. Use letters, numbers, spaces and keyboard symbols only.';
        input.setAttribute("aria-invalid", "true");
    } else {
        input.removeAttribute("aria-invalid");
        if (clearServerError) errorBox.textContent = "";
    }

    var blocked = bad !== null || text === "";
    drawBtn.disabled = blocked;
    downloadBtn.disabled = blocked;
}

// Ctrl + Enter (or Cmd + Enter on Mac) draws without reaching for the mouse.
input.addEventListener("keydown", function (e) {
    if (e.key === "Enter" && (e.ctrlKey || e.metaKey) && !drawBtn.disabled) {
        e.preventDefault();
        form.requestSubmit(drawBtn);
    }
});

// Show that something is happening after clicking Draw.
form.addEventListener("submit", function (e) {
    if (e.submitter === drawBtn) {
        drawBtn.textContent = "Drawing…";
    }
});

// Coming back with the browser's Back button: reset the button text.
window.addEventListener("pageshow", function () {
    drawBtn.textContent = "Draw";
});

if (copyBtn) {
    copyBtn.addEventListener("click", function () {
        navigator.clipboard.writeText(out.value).then(function () {
            copyBtn.textContent = "Copied!";
        }, function () {
            copyBtn.textContent = "Copy failed";
        });
        setTimeout(function () { copyBtn.textContent = "Copy"; }, 1500);
    });
}

input.addEventListener("input", function () { update(true); });
window.addEventListener("resize", fit);
update(false);
fit();