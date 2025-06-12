document.addEventListener("DOMContentLoaded", function () {
    document.querySelectorAll("pre").forEach((preBlock) => {
        let code = preBlock.querySelector("code:not(.language-mermaid)");
        if (!code) return;

        let button = document.createElement("button");
        button.innerText = "Copy";
        button.classList.add("copy-btn");

        preBlock.style.position = "relative";
        preBlock.appendChild(button);

		button.addEventListener("click", () => {
		    let text = code.innerText;
		    let tempTextArea = document.createElement("textarea");
		    tempTextArea.value = text;
		    document.body.appendChild(tempTextArea);
		    tempTextArea.select();
		    document.execCommand("copy");
		    document.body.removeChild(tempTextArea);
		    button.innerText = "Copied!";
		    setTimeout(() => (button.innerText = "Copy"), 2000);
		});
    });
});
