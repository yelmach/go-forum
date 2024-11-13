const data = {
    allPosts: [],
    allCategories: []
}

const loadData = async () => {
    data.allPosts = await fetch("/api/posts")
        .then((response) => response.json())
    data.allCategories = await fetch("/api/categories")
        .then(response => response.json());
};

const init = async () => {
    await loadData();
    const categContainer = document.querySelector('.categories');

    for (const category of data.allCategories) {
        const categoryElem = document.createElement('li')
        categoryElem.id = category.name;
        categoryElem.onclick = () => filterByCategory(category.name);
        categoryElem.innerHTML = `
        <i class="ri-hashtag"></i>${category.name}
        `
        categContainer.append(categoryElem);
    }
    displayPosts(data.allPosts);
    document.getElementById('select_1').classList.add('active');
}

const createPostElement = (post) => {
    const userId = parseInt(getCookie("user_id"));
    const postDiv = document.createElement("div");
    const likeActive = post.likes.includes(userId) ? ' liked' : ''
    const dislikeActive = post.dislikes.includes(userId) ? ' disliked' : ''
    postDiv.dataset.id = post.id;
    postDiv.classList.add("post");
    let splitedContent = post.content.split("\n")
    postDiv.innerHTML = `
    <div class="user-info">
        <img src="https://ui-avatars.com/api/?name=${post.by}" alt="User avatar" class="avatar">
        <div>
            <div class="username">${post.by}</div>
            <div class="timestamp">${timeAgo(new Date(post.createdAt).getTime())}</div>
        </div>
    </div>
    <div class="post-content">
        <h3 onclick="openPost(${post.id})">${post.title}</h3>
        <p>${post.content.replace(/\n/g, "</br>")}</p>
    </div>
    <div class="tags-stats">
        <div class="tags">
            ${post.categories.map(tag => `<span class="tag">${tag}</span>`).join('\n')}
        </div>
        <div class="post-stats">
            <div class="stat">
                <i class="ri-thumb-up-line${likeActive}" onclick="likeAction(${post.id}, true)"></i><span class="${likeActive}">${post.likes.length}</span>
            </div>
            <div class="stat">
                <i class="ri-thumb-down-line dislike${dislikeActive}" onclick="dislikeAction(${post.id}, true)"></i><span class="${dislikeActive}">${post.dislikes.length}</span>
            </div>
            <div class="stat">
                <i class="ri-chat-4-line" onclick="openPost(${post.id})"></i><span>${post.comments.length}</span>
            </div>
        </div>
    </div>
    `
    return postDiv;
}

const createCommentElement = (comment) => {
    const userId = parseInt(getCookie("user_id"));
    const commentDiv = document.createElement('div');
    const likeActive = comment.likes.includes(userId) ? ' liked' : ''
    const dislikeActive = comment.dislikes.includes(userId) ? ' disliked' : ''

    commentDiv.dataset.id = comment.id;
    commentDiv.classList.add('comment');
    commentDiv.innerHTML = `
    <div class="user-info">
        <img src="https://ui-avatars.com/api/?name=${comment.by}" alt="User avatar" class="avatar">
        <div>
            <div class="username">${comment.by}</div>
            <div class="timestamp">${timeAgo(new Date(comment.createdAt).getTime())}</div>
        </div>
    </div>
    <div class="content">
        <p>${comment.content}</p>
    </div>
    <div class="tags-stats">
        <div class="post-stats">
            <div class="stat">
                <i class="ri-thumb-up-line${likeActive}" onclick="likeAction(${comment.id}, false)"></i><span class="${likeActive}">${comment.likes.length}</span>
            </div>
            <div class="stat">
                <i class="ri-thumb-down-line dislike${dislikeActive}" onclick="dislikeAction(${comment.id}, false)"></i><span class="${dislikeActive}">${comment.dislikes.length}</span>
            </div>
        </div>
    </div>
    `
    return commentDiv;
}

const getPostData = async (postId) => {
    try {
        const response = await fetch(`/api/posts/${postId}`);
        if (!response.ok) {
            throw new Error(`Response status: ${response.status}`);
        }

        return response.json();
    } catch (error) {
        console.error(error.message);
    }
}

const openPost = async (postId) => {
    const post = await getPostData(postId);
    const main = document.querySelector('.main');
    const widget = document.querySelector('.widget');
    const comments = document.createElement('div');
    comments.classList.add('comments');

    for (const comment of post.comments) {
        const commentDiv = createCommentElement(comment);
        comments.append(commentDiv);
    }
    main.innerHTML = createPostElement(post).outerHTML;
    main.innerHTML += `
    <div class="comment-box">
        <textarea class="comment-input" placeholder="Type here your wise suggestion"></textarea>
        <div class="button-group">
            <button class="btn btn-cancel">Cancel</button>
            <button class="btn btn-comment">
                <i class="ri-chat-new-line"></i>Comment
            </button>
        </div>
    </div>
    `
    widget.innerHTML = `
    <img src="https://ui-avatars.com/api/?name=${post.by}" alt="User avatar">
    <p class="username">@${post.by}</p>
    `
    main.append(comments);

    document.querySelector('.btn-cancel').onclick = () => {
        document.querySelector('.comment-input').value = '';
    };
    document.querySelector('.btn-comment').onclick = async () => {
        const commentArea = document.querySelector('.comment-input');
        const content = commentArea.value;
        commentArea.value = '';

        if (content.trim() == "") {
            commentArea.placeholder = 'Please type a valid comment ⚠'
            commentArea.style.setProperty('--placeholder-color', 'red');
            return
        }
        try {
            await fetch("/newcomment", {
                method: "POST",
                body: JSON.stringify({ postId, content })
            })
            openPost(postId);
        } catch (err) {
            console.error(err)
        }
    };
}

const displayPosts = (posts) => {
    const main = document.querySelector('.main');
    const postsContainer = document.createElement('div');
    postsContainer.classList.add('posts-container');
    disactive();
    main.innerHTML = ''
    if (!posts.length) {
        main.innerHTML += `
        <img id="no_data" src="/assets/img/no_data.svg" alt="no result"/>
        `
    } else {
        for (const post of posts) {
            // console.log(post.content)
            const postDiv = createPostElement(post);
            postsContainer.append(postDiv);
        }
        main.append(postsContainer)
    }
}

const likeAction = async (id, isPost) => {
    const reqData = isPost ? { postId: id, isLike: true } : { commentId: id, isLike: true }
    try {
        await fetch("/reaction", {
            method: "POST",
            body: JSON.stringify(reqData)
        })
        if (isPost) {
            const post = await getPostData(id);
            const postDiv = document.querySelector(`.post[data-id="${id}"]`)
            postDiv.innerHTML = createPostElement(post).innerHTML;
        } else {
            const postId = function () {
                return document.querySelector('.post').dataset.id
            }();
            openPost(postId);
        }
    } catch (err) {
        console.error(err)
    }
}

const dislikeAction = async (id, isPost) => {
    const reqData = isPost ? { postId: id, isDislike: true } : { commentId: id, isDislike: true }
    try {
        await fetch("/  ", {
            method: "POST",
            body: JSON.stringify(reqData)
        })
        if (isPost) {
            const post = await getPostData(id);
            const postDiv = document.querySelector(`.post[data-id="${id}"]`)
            postDiv.innerHTML = createPostElement(post).innerHTML;
        } else {
            const postId = function () {
                return document.querySelector('.post').dataset.id
            }();
            openPost(postId);
        }
    } catch (err) {
        console.error(err)
    }
}

const timeAgo = (time) => {
    const seconds = Math.floor((Date.now()) - time) / 1000;
    const intervals = {
        year: 31536000,
        month: 2592000,
        week: 604800,
        day: 86400,
        hour: 3600,
        minute: 60
    };
    for (const [unit, secondsInUnit] of Object.entries(intervals)) {
        const interval = Math.floor(seconds / secondsInUnit);
        if (interval >= 1) {
            return `${interval} ${unit}${interval === 1 ? '' : 's'} ago`;
        }
    }
    return 'just now';
}

function getCookie(name) {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) return parts.pop().split(';').shift();
}

const disactive = () => {
    const elems = [...document.querySelectorAll('.menu-select'), ...document.querySelectorAll('.categories>li')]
    for (const elem of elems) {
        elem.classList.remove('active');
        elem.classList.remove('activeCat');
    }
}

const widgetBack = () => {
    const widget = document.querySelector('.widget');
    widget.innerHTML = `
    <div class="section">
        <h2 class="section-title">
            <i class="ri-star-line"></i>
            Must-read posts
        </h2>
        <ul>
            <li><a href="#">Please read rules before you start working on a platform</a></li>
            <li><a href="#">Vision & Strategy of AIemhelp</a></li>
        </ul>
    </div>
    <div class="section">
        <h2 class="section-title">
            <i class="ri-links-line"></i>
            Featured links
        </h2>
        <ul>
            <li><a href="#">AIemhelp source-code on GitHub</a></li>
            <li><a href="#">Golang best-practices</a></li>
            <li><a href="#">AIem School dashboard</a></li>
        </ul>
    </div>
    `
}

const recentPosts = async () => {
    await loadData();
    widgetBack();
    displayPosts(data.allPosts);
    document.getElementById('select_1').classList.add('active');
}

const createdPosts = async () => {
    const username = getCookie("username");
    const posts = data.allPosts.filter(post => post.by == username)
    await loadData();
    widgetBack();
    displayPosts(posts);
    document.getElementById('select_2').classList.add('active');
}

const likedPosts = async () => {
    await loadData();
    const userId = parseInt(getCookie("user_id"));
    const posts = data.allPosts.filter(post => post.likes.includes(userId))
    widgetBack();
    displayPosts(posts);
    document.getElementById('select_3').classList.add('active');
}

const filterByCategory = async (category) => {
    const posts = data.allPosts.filter(post => post.categories.includes(category))
    await loadData();
    widgetBack();
    displayPosts(posts);
    document.getElementById(category).classList.add('activeCat');
}

const logout = () => {
    try {
        fetch('/auth/logout', {
            method: 'POST'
        });
        location.href = "/";
    } catch (err) {
        console.error(err);
    }
}

const newPost = () => {
    document.querySelector('.main').innerHTML = `
    <form id="newPostForm">
        <div class="multi-select">
            <div class="selected-tags" id="selectedTags" data-placeholder='Select categories (optional)'></div>
            <div class="dropdown" id="dropdown">
                <input type="text" class="search-box" placeholder="Search categories..." id="searchBox">
                <div class="select-all" id="selectAll">select all</div>
                <div class="options" id="options"></div>
            </div>
        </div>
        <input type="text" name="title" placeholder="Type catching attention title" required>
        <textarea name="content" placeholder="Type some content" required></textarea>
        <div class="button-container">
            <button class="btn btn-add-image">
                <i class="ri-image-add-line"></i>Add Image
            </button>
            <button class="btn btn-publish">
            <i class="fa-regular fa-paper-plane"></i>Publish
            </button>
        </div>
    </form>
    `
    const tags = data.allCategories;
    let selectedTags = [];

    const selectedTagsContainer = document.getElementById('selectedTags');
    const dropdown = document.getElementById('dropdown');
    const searchBox = document.getElementById('searchBox');
    const optionsContainer = document.getElementById('options');
    const selectAllBtn = document.getElementById('selectAll');
    const newPostForm = document.getElementById('newPostForm');

    function renderTags() {
        selectedTagsContainer.innerHTML = selectedTags.map(tag => `
            <span class="tag">
                ${tag}
                <button class="tag-remove" data-tag="${tag}">×</button>
            </span>
        `).join('');

        const filteredTags = tags.filter(tag =>
            tag.name.toLowerCase().includes(searchBox.value.toLowerCase())
        );

        optionsContainer.innerHTML = filteredTags.map(tag => `
            <div class="dropdown-item ${selectedTags.includes(tag) ? 'selected' : ''}" data-tag="${tag.name}">
                ${tag.name}
            </div>
        `).join('');
    }

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
            selectedTags = tags.map(tag => tag.name);
        }
        
        renderTags();
    });

    newPostForm.addEventListener("submit", async () => {
        const title = document.querySelector('input[name="title"]').value;
        const content = document.querySelector('textarea[name="content"]').value;
        try {
            await fetch('/newpost', {
                method: 'POST',
                body: JSON.stringify({
                    Title: title,
                    Content: content,
                    Categories: selectedTags
                })
            })
            // location.href = "/";
        } catch (err) {
            console.error(err);
        }
    });

    renderTags();
}

init();
console.log(`
╱╱╱╱╱╱╱╱╱╱╱╱╱╭━━━╮╭╮
╱╱╱╱╱╱╱╱╱╱╱╱╱┃╭━╮┣╯┃
╭━━━┳━━┳━╮╭━━┫┃┃┃┣╮┃
┣━━┃┃╭╮┃╭╮┫┃━┫┃┃┃┃┃┃
┃┃━━┫╰╯┃┃┃┃┃━┫╰━╯┣╯╰╮
╰━━━┻━━┻╯╰┻━━┻━━━┻━━╯
`)