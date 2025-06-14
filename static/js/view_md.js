import mermaid from 'https://cdn.jsdelivr.net/npm/mermaid@11/dist/mermaid.esm.min.mjs';

hljs.registerAliases("mermaid", { languageName: "plaintext" });
hljs.highlightAll();

mermaid.initialize({
  startOnLoad: false,
  theme: 'dark',
  securityLevel: 'loose'
});

mermaid.run({
  querySelector: '.language-mermaid',
});

if (window.MathJax && MathJax.typesetPromise) {
  MathJax.typesetPromise([document.querySelector('.markdown-body')])
    .catch((err) => console.error("MathJax render error:", err));
}
