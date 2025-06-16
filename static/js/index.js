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

function createCardElement(file) {
    const cardContainer = document.createElement("div");
    cardContainer.className = "card-container";
    
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
        img.src = `static/thumbnails/${file.filename.replace(".pdf", ".webp")}`;
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
    cardContainer.appendChild(createFavoriteButton(file));
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
