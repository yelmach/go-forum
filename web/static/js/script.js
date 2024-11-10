const tags = await fetch("/api/categories")
        .then(response => response.json());
let selectedTags = [];
console.log(tags)

const selectedTagsContainer = document.getElementById('selectedTags');
const dropdown = document.getElementById('dropdown');
const searchBox = document.getElementById('searchBox');
const optionsContainer = document.getElementById('options');
const selectAllBtn = document.getElementById('selectAll');
const submitBtn = document.getElementById('submitBtn');

// Initialize
function renderTags() {
    selectedTagsContainer.innerHTML = selectedTags.map(tag => `
        <span class="tag">
            ${tag}
            <button class="tag-remove" data-tag="${tag}">×</button>
        </span>
    `).join('');

    // Update options
    const filteredTags = tags.filter(tag => 
        tag.toLowerCase().includes(searchBox.value.toLowerCase())
    );
    
    optionsContainer.innerHTML = filteredTags.map(tag => `
        <div class="dropdown-item ${selectedTags.includes(tag) ? 'selected' : ''}" data-tag="${tag}">
            ${tag}
        </div>
    `).join('');
}

// Event Listeners
selectedTagsContainer.addEventListener('click', () => {
    dropdown.classList.add('show');
});

document.addEventListener('click', (e) => {
    if (!dropdown.contains(e.target) && !selectedTagsContainer.contains(e.target)) {
        dropdown.classList.remove('show');
    }
});

selectedTagsContainer.addEventListener('click', (e) => {
    if (e.target.classList.contains('tag-remove')) {
        const tag = e.target.dataset.tag;
        selectedTags = selectedTags.filter(t => t !== tag);
        renderTags();
    }
});

optionsContainer.addEventListener('click', (e) => {
    if (e.target.classList.contains('dropdown-item')) {
        const tag = e.target.dataset.tag;
        if (selectedTags.includes(tag)) {
            selectedTags = selectedTags.filter(t => t !== tag);
        } else {
            selectedTags.push(tag);
        }
        renderTags();
    }
});

searchBox.addEventListener('input', renderTags);

selectAllBtn.addEventListener('click', () => {
    if (selectedTags.length === tags.length) {
        selectedTags = [];
    } else {
        selectedTags = [...tags];
    }
    renderTags();
});

submitBtn.addEventListener('click', () => {
    console.log('Selected tags:', selectedTags);
    // Add your submit logic here
});

// Initial render
renderTags();
