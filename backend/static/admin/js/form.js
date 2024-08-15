// Toolbar functions for the rich text editor
function wrapText(tag) {
  const editor = document.getElementById("editor");
  const start = editor.selectionStart;
  const end = editor.selectionEnd;
  const selectedText = editor.value.substring(start, end);
  const wrappedText = `<${tag}>${selectedText}</${tag}>`;
  editor.setRangeText(wrappedText, start, end, "end");

  //   Update preview
  updatePreview();
}
function addLink() {
  const url = prompt("Enter the URL", "https://example.com");
  if (url) {
    wrapTextWithCustomTag(`<a href="${url}">`, "</a>");
  }
}
function addCode() {
  wrapTextWithCustomTag("<pre><code>", "</code></pre>");
}
function addTable() {
  const rows = prompt("Enter the number of rows", "2");
  const cols = prompt("Enter the number of columns", "2");
  if (rows && cols) {
    let table = '<table border="1">\n';
    for (let i = 0; i < rows; i++) {
      table += "  <tr>\n";
      for (let j = 0; j < cols; j++) {
        table += "    <td>&nbsp;</td>\n";
      }
      table += "  </tr>\n";
    }
    table += "</table>\n";
    insertTextAtCursor(table);

    //  Update preview
    updatePreview();
  }
}
function addBulletList() {
  wrapTextWithCustomTag("<ul><li>", "</li></ul>");
}
function addNumberedList() {
  wrapTextWithCustomTag("<ol><li>", "</li></ol>");
}

// Keydown (To handle return key)
function handleKeyDown(event) {
  const editor = document.getElementById("editor");
  if (document.activeElement === editor && event.key === "Enter") {
    event.preventDefault();
    insertTextAtCursor("\n<br>\n");
  }
}

// Preview window
function updatePreview() {
  const editorContent = document.getElementById("editor").value;
  const preview = document.getElementById("preview");
  preview.innerHTML = editorContent;
}

// Helper functions
function wrapTextWithCustomTag(openTag, closeTag) {
  const editor = document.getElementById("editor");
  const start = editor.selectionStart;
  const end = editor.selectionEnd;
  const selectedText = editor.value.substring(start, end);
  const wrappedText = `${openTag}${selectedText}${closeTag}`;
  editor.setRangeText(wrappedText, start, end, "end");

  //   Update preview
  updatePreview();
}
function insertTextAtCursor(text) {
  const editor = document.getElementById("editor");
  const start = editor.selectionStart;
  editor.setRangeText(text, start, start, "end");

  //   Update preview
  updatePreview();
}
