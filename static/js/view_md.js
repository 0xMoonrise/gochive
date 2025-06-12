import mermaid from 'https://cdn.jsdelivr.net/npm/mermaid@11/dist/mermaid.esm.min.mjs';

hljs.registerAliases("mermaid", { languageName: "plaintext" })
hljs.highlightAll();

mermaid.initialize({
  startOnLoad: false,
  theme: 'dark',
  securityLevel: 'loose'
});

mermaid.run({
  querySelector: '.language-mermaid',
});

window.onload = function() {
    if (window.MathJax) {
        MathJax.Hub.Queue(["Typeset", MathJax.Hub]);
    } else {
        setTimeout(arguments.callee, 100);
    }
}
MathJax.Hub.Config({
    tex2jax: {
        inlineMath: [['$', '$']],
        displayMath: [['$$', '$$']]
    }
});
