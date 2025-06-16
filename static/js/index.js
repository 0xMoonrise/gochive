//My code is crap but is my code

// window.onload

//https://medium.com/@ryan_forrester_/how-to-get-a-cookie-by-name-in-javascript-ff36761e5356
function getCookieByName(name) {
    const cookieString = document.cookie;
    const cookies = cookieString.split(';');
    for (let i = 0; i < cookies.length; i++) {
        const cookie = cookies[i].trim();
        if (cookie.startsWith(name + '=')) {
            return cookie.substring(name.length + 1);
        }
    }
    return null;
}

const form = document.querySelector(".form-input");

form.addEventListener("submit", (event) => {
    event.preventDefault();
    sendData();
});

function uploadFile()
{
    const fileInput = document.getElementById('file-input');
    const file = fileInput.files[0];

    if (file)
    {
        const formData = new FormData();
        formData.append('file', file);

        fetch('/upload', {
            method: 'POST',
            body: formData
        })
        .then(response => response.json())
        .then(data => {
            if (data.success)
            {
                window.location.reload();
            }
        })
        .catch(error => {
            console.error('Error:', error);
        });
    }
}

async function sendData()
{

    const search = new FormData(form);
	let page = parseInt(getCookieByName('page'));
	page = page ? page : 1;
    const params = new URLSearchParams();
    params.append('search', search.get('query'));

    try
    {
	    const response = await fetch(`${window.location.origin}/search/${page}`,
	    {
	        method: "POST",
	        body: params,
	    });

	    archive = await response.json();

		const card_list = document.querySelector(".card-list");
		const cards = card_list.querySelectorAll(".card-container");
		cards.forEach(card => card.remove());
		// Loaded from search
		archive.files.forEach(file => {
    		make_cardElement(file);
		});
		make_paginationSection(archive.pages, page);
    }
    catch (e)
    {
        console.error(e);
    }
}

function createFavoriteButton(file) {
    const button_favorite = document.createElement("button");
    button_favorite.innerText = "★";
    button_favorite.className = "card-button-favorite";

    const activeColor = "#FFD700";
    const inactiveColor = "#bb86fc";

    button_favorite.style.color = file.favorite ? activeColor : inactiveColor;

    button_favorite.addEventListener("mouseenter", () => {
        button_favorite.style.color = file.favorite ? inactiveColor : activeColor;
    });

    button_favorite.addEventListener("mouseleave", () => {
        button_favorite.style.color = file.favorite ? activeColor : inactiveColor;
    });

    button_favorite.addEventListener("click", (event) => {
        event.stopPropagation();
        const prevFavorite = file.favorite;
        file.favorite = !prevFavorite;

        button_favorite.style.color = file.favorite ? activeColor : inactiveColor;

        const formData = new FormData();
        formData.append('favorite', file.favorite);

        fetch(`/set_favorite/${file.id}`, {
            method: 'POST',
            body: formData,
        })
        .then(response => {
            if (!response.ok) throw new Error('Error setting up the value');
            return response.json();
        })
        .catch(error => {
            file.favorite = prevFavorite;
            button_favorite.style.color = file.favorite ? activeColor : inactiveColor;
        });
    });

    return button_favorite;
}

function make_cardElement(file)
{

	const cardList = document.querySelector(".card-list");
	const card_container = document.createElement("div");

	card_container.classList.add("card-container");
	const card = document.createElement("a");

	card.classList.add("card");
	const card_tittle = document.createElement("div");
	card_tittle.classList.add("card-title");
	card_tittle.innerText = file.filename;

	const button_favorite = createFavoriteButton(file);
	
	if (file.favorite) {
		button_favorite.style.color = "#FFD700";
	} else {
		button_favorite.style.color = "#bb86fc";
	}
	// #FFD700 -> true
	// #bb86fc -> false

	if(file.filename.includes('.pdf'))
	{
		const img = document.createElement("img");

		img.classList.add("card-thumbnail");
		img.src = `static/thumbnails/${file.filename.replace(".pdf", ".webp")}`;

		card.append(img);
		card.href = `view?file=${file.id}`;
	}

	if(file.filename.includes('.md'))
	{
		const iframe = document.createElement("iframe");

		iframe.classList.add("card-thumbnail")
		iframe.src = `get_file/${file.id}`;
		iframe.scrolling = "no";

		card.append(iframe);
		card.href = `get_file/${file.id}`;
	}
	
	card_container.append(button_favorite);
	card_container.append(card);
	card.append(card_tittle);
	cardList.append(card_container);
}

function make_paginationSection(n_button, current)
{
    const pagination_section = document.querySelector(".pagination-section");
    pagination_section.innerHTML = "";
	document.cookie = `page=${current}`; 

    const currentColor = "#e0e0e0"
    const max_buttons = 5;

    let start, end;

    if(n_button <= max_buttons)
    {
        start = current;
        end = n_button;
    }
    else
    {
        if (current <= 3) {
            start = 1;
            end = max_buttons - 1;
        } else if (current >= n_button - 2) {
            start = n_button - (max_buttons - 2);
            end = n_button;
        } else {
            start = current - 1;
            end = current + 1;
        }
    }

    let button1 = document.createElement("button");
    button1.innerText = 1;
    button1.classList.add("pagination-button");

    if(current === 1){
        //button1.classList.add("active");
        button1.style.backgroundColor = currentColor;
    }

    button1.addEventListener("click", function () {
        make_paginationSection(n_button, 1);
        button_event(this);
    });

    pagination_section.append(button1);

    if(start > 2) {
        let dots = document.createElement("button");
        dots.classList.add("pagination-button");
        dots.innerText = "...";
        pagination_section.append(dots);
    }

    for(let i = start; i <= end; i++) {
        if (i === 1 || i === n_button) continue;
        let button = document.createElement("button");

        button.innerText = i;
        button.classList.add("pagination-button");

        if (i === current){
             button.classList.add("active");
             button.style.backgroundColor = currentColor;
        }

        button.addEventListener("click", function () {
            make_paginationSection(n_button, i);
            button_event(this);
        });

        pagination_section.append(button);
    }

    if(end < n_button - 1) {
		let dots = document.createElement("button");

		dots.classList.add("pagination-button");
        dots.innerText = "...";

        pagination_section.append(dots);
    }

    if(n_button > 1) {
        let buttonLast = document.createElement("button");

        buttonLast.innerText = n_button;
        buttonLast.classList.add("pagination-button");

        if (current === n_button){
            buttonLast.classList.add("active");
            buttonLast.style.backgroundColor = currentColor;
        }

        buttonLast.addEventListener("click", function () {
            make_paginationSection(n_button, n_button);
            button_event(this);
        });

        pagination_section.append(buttonLast);
    }
}

document.addEventListener("DOMContentLoaded", async () => {
    try
    {
    	let page = parseInt(getCookieByName('page'));
    	page = page ? page : 1; 
    	
        const response = await fetch(`/get_files/${page}`);
        const archive = await response.json();

		//Loaded from all
        archive.files.forEach(files => {
            make_cardElement(files);
        });
        
		make_paginationSection(archive.pages, page);
    }
    catch (error)
    {
        console.error(error);
    }
});

async function button_event(e) {

    const search = new FormData(form);
    const params = new URLSearchParams();
    params.append('search', search.get('query'));

    let response;
    try {
        if (search.get('query')) {
            response = await fetch(`${window.location.origin}/search/${e.innerText}`, {
                method: "POST",
                body: params,
            });
        } else {
            response = await fetch(`/get_files/${e.innerText}`);
        }

        if (!response.ok) {
            throw new Error(`Error: ${response.status} ${response.statusText}`);
        }

        const archive = await response.json();

        const card_list = document.querySelector(".card-list");
        const cards = card_list.querySelectorAll(".card-container");

        cards.forEach(card => card.remove());

        archive.files.forEach(filename => {
            make_cardElement(filename);
        });

    } catch (error) {
        console.error("Error:", error);
    }
}

