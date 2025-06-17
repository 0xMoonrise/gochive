const ACTIVE_COLOR = "#FFD700";
const INACTIVE_COLOR = "#bb86fc";
const PAGINATION_BUTTONS = 5;

function getCookie(name) {
    return document.cookie
        .split('; ')
        .find(row => row.startsWith(name + '='))
        ?.split('=')[1];
}

function updateFileView(files) {
    const cardList = document.querySelector(".card-list");
    cardList.querySelectorAll(".card-container").forEach(card => card.remove());
    files.forEach(file => createCardElement(file));
}

function updatePagination(totalPages, currentPage) {
    const paginationSection = document.querySelector(".pagination-section");
    paginationSection.innerHTML = "";
    document.cookie = `page=${currentPage}`;
    
    const createPageButton = (page, isActive = false) => {
        const button = document.createElement("button");
        button.textContent = page;
        button.className = "pagination-button";
        if (isActive) button.style.backgroundColor = "#e0e0e0";
        return button;
    };

    paginationSection.appendChild(createPageButton(1, currentPage === 1));

    const start = Math.max(2, currentPage - 1);
    const end = Math.min(totalPages - 1, currentPage + 1);
    
    if (start > 2) paginationSection.appendChild(createPageButton("..."));
    
    for (let i = start; i <= end; i++) {
        paginationSection.appendChild(createPageButton(i, i === currentPage));
    }
    
    if (end < totalPages - 1) paginationSection.appendChild(createPageButton("..."));
    
    if (totalPages > 1) {
        paginationSection.appendChild(
            createPageButton(totalPages, currentPage === totalPages)
        );
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
            
            if (!response.ok) {
                file.favorite = wasFavorite;
                button.style.color = wasFavorite ? ACTIVE_COLOR : INACTIVE_COLOR;
                throw new Error('Favorite update failed');
            }
        } catch (error) {
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
    
    const titleLabel = document.createElement("label");
    titleLabel.textContent = "Titulo:";
    const titleInput = document.createElement("input");
    titleInput.type = "text";
    titleInput.name = "filename";
    titleInput.value = file.filename;
    titleInput.required = true;
    
    const publisherLabel = document.createElement("label");
    publisherLabel.textContent = "Editorial:";
    const publisherInput = document.createElement("input");
    publisherInput.type = "text";
    publisherInput.name = "editorial";
    publisherInput.value = file.editorial || "";
    
    const submitButton = document.createElement("button");
    submitButton.type = "submit";
    submitButton.textContent = "Guardar";
    
    form.appendChild(titleLabel);
    form.appendChild(titleInput);
    form.appendChild(publisherLabel);
    form.appendChild(publisherInput);
    form.appendChild(submitButton);
    
    modalContent.appendChild(closeButton);
    modalContent.appendChild(form);
    modal.appendChild(modalContent);
    document.body.appendChild(modal);
    
    // Event listeners
    closeButton.addEventListener("click", () => modal.remove());
    window.addEventListener("click", (event) => {
        if (event.target === modal) modal.remove();
    });
    
    return { modal, form };
}

function createEditButton(file){
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
			const updatedData = new URLSearchParams();

			updatedData.append("filename", formData.get("filename") || '');
			updatedData.append("editorial", formData.get("editorial") || '');
            try {
                const response = await fetch(`/edit/${file.id}`, {
                    method: "POST",
					    headers: {
					        "Content-Type": "application/x-www-form-urlencoded"
					    },
					    body: updatedData
                });
                
                if (response.ok) {
                    const currentPage = parseInt(getCookie('page')) || 1;
                    const query = document.querySelector('.form-input [name="query"]').value;
                    loadFiles(currentPage, query);
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

function createCardElement(file) {
    const cardContainer = document.createElement("div");
	const buttonContainer = document.createElement("div");

    cardContainer.className = "card-container";
    buttonContainer.className = "button-container";
    
    const card = document.createElement("a");
    card.className = "card";
    card.href = file.filename.includes('.pdf') ? 
        `view?file=${file.id}` : `get_file/${file.id}`;
    
    const title = document.createElement("div");
    title.className = "card-title";
    title.textContent = file.filename;
    
    if (file.filename.includes('.pdf')) {
        const img = document.createElement("img");
        img.className = "card-thumbnail";
        img.src = `static/thumbnails/${file.id}`;
        card.appendChild(img);
    } 
    else if (file.filename.includes('.md')) {
        const iframe = document.createElement("iframe");
        iframe.className = "card-thumbnail";
        iframe.src = `get_file/${file.id}`;
        iframe.scrolling = "no";
        card.appendChild(iframe);
    }
    
    card.appendChild(title);

    buttonContainer.appendChild(createFavoriteButton(file));
	buttonContainer.appendChild(createEditButton(file));
	cardContainer.appendChild(buttonContainer);
    cardContainer.appendChild(card);

    document.querySelector(".card-list").appendChild(cardContainer);
}

async function loadFiles(page, searchQuery = null) {
    try {
        const url = searchQuery ? 
            `${window.location.origin}/search/${page}` : `/get_files/${page}`;
        
        const options = searchQuery ? {
            method: "POST",
            body: new URLSearchParams({ search: searchQuery })
        } : { method: "GET" };
        
        const response = await fetch(url, options);
        const data = await response.json();
        
        updateFileView(data.files);
        updatePagination(data.pages, page);
    } catch (error) {
        console.error("Error loading files:", error);
    }
}

function handleSearch(event) {
    event.preventDefault();
    const query = document.querySelector('.form-input [name="query"]').value;
    loadFiles(1, query);
}

function handlePagination(event) {
    if (!event.target.classList.contains('pagination-button')) return;
    
    const pageText = event.target.textContent;
    if (pageText === "...") return;
    
    const page = parseInt(pageText);
    const query = document.querySelector('.form-input [name="query"]').value;
    loadFiles(page, query);
}

document.addEventListener("DOMContentLoaded", () => {
    const initialPage = parseInt(getCookie('page')) || 1;
    loadFiles(initialPage);
    
    document.querySelector(".form-input").addEventListener("submit", handleSearch);
    document.querySelector(".pagination-section").addEventListener("click", handlePagination);
});

function uploadFile() {
    const fileInput = document.getElementById('file-input');
    const file = fileInput.files[0];
    
    if (!file) return;
    
    const formData = new FormData();
    formData.append('file', file);
    
    fetch('/upload', {
        method: 'POST',
        body: formData
    })
    .then(response => response.json())
    .then(data => data.success && window.location.reload())
    .catch(console.error);
}
