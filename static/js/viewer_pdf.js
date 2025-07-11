import * as pdfjsLib from '/static/js/pdfjs/src/pdf.js';
import '/static/js/pdfjs/src/pdf.worker.js';

pdfjsLib.GlobalWorkerOptions.workerSrc = '/static/js/pdfjs/src/pdf.worker.js';

const container = document.getElementById('pdf-container');
const PDF_URL = '/get_file/1';

async function renderPDF(url) {
  const loadingTask = pdfjsLib.getDocument(url);
  const pdf = await loadingTask.promise;

  for (let i = 1; i <= pdf.numPages; i++) {
    const page = await pdf.getPage(i);
    const scale = 1.5;
    const viewport = page.getViewport({ scale });

    const canvas = document.createElement('canvas');
    canvas.width = viewport.width;
    canvas.height = viewport.height;
    const context = canvas.getContext('2d');

    await page.render({
      canvasContext: context,
      viewport,
    }).promise;

    container.appendChild(canvas);
  }
}

renderPDF(PDF_URL);

window.scrollToPage = (pageNumber) => {
  const canvas = container.querySelectorAll('canvas')[pageNumber - 1];
  if (canvas) canvas.scrollIntoView({ behavior: 'smooth' });
};
