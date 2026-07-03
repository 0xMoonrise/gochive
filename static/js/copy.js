function copyText(text) {
  if (navigator.clipboard && window.isSecureContext) {
    return navigator.clipboard.writeText(text).catch(() => fallbackCopy(text));
  }
  fallbackCopy(text);
}

function fallbackCopy(text) {
  const tempTextArea = document.createElement("textarea");
  tempTextArea.value = text;
  tempTextArea.style.position = "fixed";
  tempTextArea.style.opacity = "0";
  document.body.appendChild(tempTextArea);
  tempTextArea.select();
  document.execCommand("copy");
  document.body.removeChild(tempTextArea);
}

export function addCopyButtons(root = document) {
  root.querySelectorAll("pre").forEach((preBlock) => {
    if (preBlock.querySelector(".copy-btn")) return;

    const code = preBlock.querySelector("code:not(.language-mermaid)");
    if (!code) return;

    const button = document.createElement("button");
    button.innerText = "Copy";
    button.classList.add("copy-btn");

    preBlock.style.position = "relative";
    preBlock.appendChild(button);

    button.addEventListener("click", () => {
      copyText(code.innerText);
      button.innerText = "Copied!";
      setTimeout(() => (button.innerText = "Copy"), 2000);
    });
  });
}
