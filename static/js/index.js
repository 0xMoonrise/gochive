const ACTIVE_COLOR = "#FFD700";
const INACTIVE_COLOR = "#bb86fc";

const dom = {
  cardList: () => document.querySelector(".card-list"),
  pagination: () => document.querySelector(".pagination-section"),
  searchInput: () => document.querySelector('.form-input [name="query"]'),
};

function getCookie(name) {
  return document.cookie
    .split('; ')
    .find(row => row.startsWith(name + '='))
    ?.split('=')[1];
}

function currentPage() {
  const page = parseInt(getCookie('page'), 10);
  return Number.isInteger(page) && page > 0 ? page : 1;
}

function currentQuery() {
  return dom.searchInput().value;
}

function createPageButton(label, isActive = false) {
  const button = document.createElement("button");
  button.textContent = label;
  button.className = "pagination-button";
  if (isActive) button.style.backgroundColor = "#e0e0e0";
  return button;
}

function updatePagination(totalPages, page) {
  const paginationSection = dom.pagination();
  paginationSection.innerHTML = "";
  document.cookie = `page=${page}`;

  paginationSection.appendChild(createPageButton(1, page === 1));

  const start = Math.max(2, page - 1);
  const end = Math.min(totalPages - 1, page + 1);

  if (start > 2) paginationSection.appendChild(createPageButton("..."));
  for (let i = start; i <= end; i++) {
    paginationSection.appendChild(createPageButton(i, i === page));
  }
  if (end < totalPages - 1) paginationSection.appendChild(createPageButton("..."));

  if (totalPages > 1) {
    paginationSection.appendChild(createPageButton(totalPages, page === totalPages));
  }
}

function createFavoriteButton(file) {
  const button = document.createElement("button");
  button.className = "card-button-favorite";
  button.textContent = "★";
  button.style.color = file.favorite ? ACTIVE_COLOR : INACTIVE_COLOR;

  button.addEventListener("click", async (event) => {
    event.stopPropagation();
    const wasFavorite = file.favorite;
    file.favorite = !wasFavorite;
    button.style.color = file.favorite ? ACTIVE_COLOR : INACTIVE_COLOR;

    try {
      const response = await fetch(`/set_favorite/${file.id}`, {
        method: 'POST',
        body: new URLSearchParams({ favorite: file.favorite }),
      });
      if (!response.ok) throw new Error('Favorite update failed');
    } catch (error) {
      file.favorite = wasFavorite;
      button.style.color = wasFavorite ? ACTIVE_COLOR : INACTIVE_COLOR;
      console.error('Error updating favorite:', error);
    }
  });

  return button;
}

function createEditModal(file) {
  const modal = document.createElement("div");
  modal.className = "modal";
  modal.style.display = "block";

  const modalContent = document.createElement("div");
  modalContent.className = "modal-content";

  const closeButton = document.createElement("span");
  closeButton.className = "close-button";
  closeButton.innerHTML = "&times;";

  const form = document.createElement("form");
  form.className = "edit-form";

  const addField = (labelText, name, value, required = false) => {
    const label = document.createElement("label");
    label.textContent = labelText;
    const input = document.createElement("input");
    input.type = "text";
    input.name = name;
    input.value = value || "";
    input.required = required;
    form.appendChild(label);
    form.appendChild(input);
  };

  addField("Titulo:", "filename", file.filename, true);
  addField("Editorial:", "editorial", file.editorial);

  const submitButton = document.createElement("button");
  submitButton.type = "submit";
  submitButton.textContent = "Guardar";
  form.appendChild(submitButton);

  modalContent.appendChild(closeButton);
  modalContent.appendChild(form);
  modal.appendChild(modalContent);
  document.body.appendChild(modal);

  closeButton.addEventListener("click", () => modal.remove());
  window.addEventListener("click", (event) => {
    if (event.target === modal) modal.remove();
  });

  return { modal, form };
}

function createEditButton(file) {
  const buttonEdit = document.createElement("button");
  buttonEdit.className = "card-button-edit";
  buttonEdit.title = 'Edit';
  buttonEdit.innerHTML = `
    <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" fill="currentColor" viewBox="0 0 24 24">
      <path d="M3 21v-3.75l11.06-11.06 3.75 3.75L6.75 21H3zm15.41-11.34l-3.75-3.75 1.41-1.41a1 1 0 011.42 0l2.33 2.34a1 1 0 010 1.41l-1.41 1.41z"/>
    </svg>
  `;

  buttonEdit.addEventListener("click", (event) => {
    event.stopPropagation();
    const { modal, form } = createEditModal(file);

    form.addEventListener("submit", async (e) => {
      e.preventDefault();
      const formData = new FormData(form);

      try {
        const response = await fetch(`/edit/${file.id}`, {
          method: "PATCH",
          headers: { "Content-Type": "application/x-www-form-urlencoded" },
          body: new URLSearchParams({
            filename: formData.get("filename") || '',
            editorial: formData.get("editorial") || '',
          }),
        });

        if (response.ok) {
          loadFiles(currentPage(), currentQuery());
          modal.remove();
        } else {
          console.error("Error al actualizar el archivo");
        }
      } catch (error) {
        console.error("Error:", error);
      }
    });
  });

  return buttonEdit;
}

function createConfirmModal(message) {
  const modal = document.createElement("div");
  modal.className = "modal";
  modal.style.display = "block";

  const modalContent = document.createElement("div");
  modalContent.className = "modal-content";

  const text = document.createElement("p");
  text.className = "confirm-text";
  text.textContent = message;

  const actions = document.createElement("div");
  actions.className = "confirm-actions";

  const yesButton = document.createElement("button");
  yesButton.className = "confirm-button confirm-yes";
  yesButton.textContent = "YES";

  const noButton = document.createElement("button");
  noButton.className = "confirm-button confirm-no";
  noButton.textContent = "NO";

  actions.appendChild(yesButton);
  actions.appendChild(noButton);
  modalContent.appendChild(text);
  modalContent.appendChild(actions);
  modal.appendChild(modalContent);
  document.body.appendChild(modal);

  noButton.addEventListener("click", () => modal.remove());
  window.addEventListener("click", (event) => {
    if (event.target === modal) modal.remove();
  });

  return { modal, yesButton };
}

function createDeleteButton(file) {
  const buttonDelete = document.createElement("button");
  buttonDelete.className = "card-button-delete";
  buttonDelete.title = 'Delete';
  buttonDelete.innerHTML = `
    <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" fill="currentColor" viewBox="0 0 24 24">
      <path d="M19 6.41L17.59 5 12 10.59 6.41 5 5 6.41 10.59 12 5 17.59 6.41 19 12 13.41 17.59 19 19 17.59 13.41 12z"/>
    </svg>
  `;

  buttonDelete.addEventListener("click", async (event) => {
    event.stopPropagation();

    const { modal, yesButton } = createConfirmModal(`Are you sure to delete: ${file.filename} ?`);

    yesButton.addEventListener("click", async () => {
      try {
        const response = await fetch(`/file/${file.id}`, { method: "DELETE" });
        if (response.ok) {
          loadFiles(currentPage(), currentQuery());
        } else {
          console.error("Error al eliminar el archivo");
        }
      } catch (error) {
        console.error("Error:", error);
      } finally {
        modal.remove();
      }
    });
  });

  return buttonDelete;
}

function createCardElement(file) {
  const cardContainer = document.createElement("div");
  const buttonContainer = document.createElement("div");
  const buttonContainerRight = document.createElement("div");
  cardContainer.className = "card-container";
  buttonContainer.className = "button-container";
  buttonContainerRight.className = "button-container-right";

  const card = document.createElement("a");
  card.className = "card";
  card.href = `view/${file.id}`;

  const title = document.createElement("div");
  title.className = "card-title";
  title.textContent = file.filename;

  if (file.filename.includes('.pdf')) {
    const img = document.createElement("img");
    img.className = "card-thumbnail";
    img.src = `images/${file.id}`;
    card.appendChild(img);
  } else if (file.filename.includes('.md')) {
    const iframe = document.createElement("iframe");
    iframe.className = "card-thumbnail";
    iframe.src = `/view/${file.id}`;
    iframe.scrolling = "no";
    card.appendChild(iframe);
  }

  card.appendChild(title);
  buttonContainer.appendChild(createFavoriteButton(file));
  buttonContainer.appendChild(createEditButton(file));
  buttonContainerRight.appendChild(createDeleteButton(file));
  cardContainer.appendChild(buttonContainer);
  cardContainer.appendChild(buttonContainerRight);
  cardContainer.appendChild(card);

  dom.cardList().appendChild(cardContainer);
}

function updateFileView(files) {
  const cardList = dom.cardList();
  cardList.querySelectorAll(".card-container").forEach(card => card.remove());
  files.forEach(file => createCardElement(file));
}

async function loadFiles(page, searchQuery = null) {
  try {
    const url = searchQuery ? `${window.location.origin}/search/${page}` : `/get_files/${page}`;
    const options = searchQuery ?
      { method: "POST", body: new URLSearchParams({ search: searchQuery }) } :
      { method: "GET" };

    const response = await fetch(url, options);

    if (!response.ok) {
      if (page !== 1) {
        return loadFiles(1, searchQuery);
      }
      throw new Error(`Failed to load files: ${response.status}`);
    }

    const data = await response.json();
    updateFileView(data.files);
    updatePagination(data.pages, page);
  } catch (error) {
    console.error("Error loading files:", error);
  }
}

function handleSearch(event) {
  event.preventDefault();
  loadFiles(1, currentQuery());
}

function handlePagination(event) {
  if (!event.target.classList.contains('pagination-button')) return;
  const pageText = event.target.textContent;
  if (pageText === "...") return;
  loadFiles(parseInt(pageText), currentQuery());
}

async function uploadFile() {
  const fileInput = document.getElementById('file-input');
  const file = fileInput.files[0];
  if (!file) return;

  const formData = new FormData();
  formData.append('file', file);

  try {
    const response = await fetch('/upload', { method: 'POST', body: formData });
    const data = await response.json();
    if (!data.success) return;

    fileInput.value = "";
    loadFiles(currentPage(), currentQuery());
  } catch (error) {
    console.error("Error uploading file:", error);
  }
}

document.addEventListener("DOMContentLoaded", () => {
  loadFiles(currentPage());
  document.querySelector(".form-input").addEventListener("submit", handleSearch);
  dom.pagination().addEventListener("click", handlePagination);
});
