const data = {
    allPosts: [],
    allCategories: [],
    currentPage: 1,
    hasMore: true,
    currentView: 'recent',
    currentCategory: null,
};

var currentPage = 1
var hasMoreComments = true

// load post with its comments (pagination implemented)
const initCommentInfiniteScroll = (postId) => {
    const commentsContainer = document.querySelector('.comments');
    if (!commentsContainer) return;

    const observer = new IntersectionObserver((elem) => {
        if (elem[0].isIntersecting && hasMoreComments) {
            loadMoreComments(postId);
        }
    }, {
        root: document.querySelector('.posts-container'),
        threshold: 0.1
    });

    const trigger = document.createElement('div');
    trigger.className = 'scroll-trigger';
    commentsContainer.appendChild(trigger);
    observer.observe(trigger);
};

const loadMoreComments = async (postId) => {
    const nextPage = currentPage + 1;
    const post = await getPostData(postId, nextPage);
    hasMoreComments = post.hasMoreComments

    const commentsContainer = document.querySelector('.comments');
    post.comments.forEach(comment => {
        const commentDiv = createCommentElement(comment);
        commentsContainer.insertBefore(commentDiv, commentsContainer.lastElementChild);
    });

    currentPage = nextPage;
};

const getPostData = async (postId, page = 1) => {
    console.log("getPostData");

    try {
        const response = await fetch(`/api/posts/${postId}?page=${page}`);
        const res = await response.json();

        if (response.ok) {
            return res
        } else {
            console.error(res.msg)
        }
    } catch (err) {
        console.error(err);
    }
}

const openPost = async (postId) => {
    currentPage = 1;
    const post = await getPostData(postId);
    hasMoreComments = post.hasMoreComments
    const postsContainer = document.querySelector('.posts-container');
    const widget = document.querySelector('.widget');
    const comments = document.createElement('div');
    const userId = parseInt(getCookie("user_id"));
    const likeActive = post.likes.includes(userId) ? ' liked' : ''
    const dislikeActive = post.dislikes.includes(userId) ? ' disliked' : ''

    comments.classList.add('comments');
    postsContainer.innerHTML = `
    <div class="post" data-id="${postId}">
        <div class="user-info">
            <img src="https://ui-avatars.com/api/?name=${post.by}" alt="User avatar" class="avatar">
            <div>
                <div class="username">${post.by}</div>
                <div class="timestamp">${timeAgo(new Date(post.createdAt).getTime())}</div>
            </div>
        </div>
        <div class="post-content">
            <h2>${filterContent(post.title)}</h2>
            <p>${filterContent(post.content)}</p>
        </div>
        <div class="tags-stats">
            <div class="tags">
                ${post.categories.map(tag => `<span class="tag">${tag}</span>`).join('')}
            </div>
            <div class="post-stats">
                <div class="stat${likeActive}">
                    <i class="like-icon" onclick="likeAction(${post.id}, true)"></i><span>${post.likes.length}</span>
                </div>
                <div class="stat${dislikeActive}">
                    <i class="dislike-icon dislike" onclick="dislikeAction(${post.id}, true)"></i><span>${post.dislikes.length}</span>
                </div>
                <div class="stat">
                    <i class="comment-icon" onclick="openPost(${post.id})"></i><span>${post.totalComments}</span>
                </div>
            </div>
        </div>
    </div>

    <form id="commentForm" class="comment-box">
        <textarea class="comment-input" placeholder="Type here your wise suggestion" required></textarea>
        <div class="button-group">
            <button class="btn btn-cancel">Cancel</button>
            <button class="btn btn-comment">
                <i class="comment-icon"></i>Comment
            </button>
        </div>
    </div>
    `

    widget.innerHTML = `
    <img src="https://ui-avatars.com/api/?name=${post.by}" alt="User avatar">
    <p class="username">@${post.by}</p>
    `

    for (const comment of post.comments) {
        const commentDiv = createCommentElement(comment);
        comments.append(commentDiv);
    }
    postsContainer.append(comments);

    document.querySelector('.btn-cancel').onclick = () => {
        document.querySelector('.comment-input').value = '';
    };
    document.getElementById('commentForm').addEventListener('submit', async (e) => {
        e.preventDefault();

        const commentArea = document.querySelector('.comment-input');
        const content = commentArea.value;
        commentArea.value = '';

        if (content.trim() == "") {
            commentArea.placeholder = 'Please type a valid comment ⚠'
            commentArea.style.setProperty('--placeholder-color', 'red');
            return
        }

        try {
            const response = await fetch("/newcomment", {
                method: "POST",
                body: JSON.stringify({ postId, content })
            })

            if (response.ok) {
                openPost(postId);
            } else {
                const res = await response.json();
                console.error(res.msg);
                document.getElementById("loginPopup").style.display = "block";
                if (res.msg != "unauthorized user") {
                    document.querySelector(".popup-content").innerHTML = `
                    <h2>Nice try!</h2>
                    <ul>
                        <li>${res.msg}</li>
                    </ul>
                    `
                }

            }
        } catch (err) {
            console.error(err)
        }
    })

    initCommentInfiniteScroll(postId);
}

// load all posts (pagination implemented)
const initInfiniteScroll = () => {
    const postsContainer = document.querySelector('.posts-container');
    if (!postsContainer) return;
    const observer = new IntersectionObserver((elem) => {
        if (elem[0].isIntersecting && data.hasMore) {
            loadMorePosts();
        }
    }, {
        root: postsContainer,
        threshold: 0.1
    });

    const trigger = document.createElement('div');
    trigger.className = 'scroll-trigger';
    postsContainer.appendChild(trigger);
    observer.observe(trigger);
};

const loadMorePosts = async () => {
    const nextPage = data.currentPage + 1;
    await loadData(nextPage, false);

    const postsContainer = document.querySelector('.posts-container');

    data.allPosts.slice(-100).forEach(post => {
        const postDiv = createPostElement(post);
        postsContainer.insertBefore(postDiv, postsContainer.lastElementChild);
    });
};

const loadData = async (page = 1, resetData = false) => {
    if (resetData) {
        data.allPosts = [];
        data.currentPage = 1;
        data.hasMore = true;
    }

    if (!data.hasMore) return;

    try {
        let url = `/api/posts?page=${page}`;

        switch (data.currentView) {
            case 'created':
                url += `&filter=created&userId=${getCookie('user_id')}`;
            case 'liked':
                url += `&filter=liked&userId=${getCookie('user_id')}`;
            case 'category':
                url += `&filter=category&category=${data.currentCategory}`;
        }

        const response = await fetch(url);
        const result = await response.json();

        if (resetData) {
            data.allPosts = result.posts;
        } else {
            data.allPosts = [...data.allPosts, ...result.posts];
        }

        data.hasMore = result.hasMore;
        data.currentPage = result.page;

        if (data.allCategories.length === 0 || page === 1) {
            data.allCategories = await fetch("/api/categories")
                .then(response => response.json());
        }
    } catch (error) {
        console.error('Error loading posts:', error);
    }
};

const init = async () => {
    const categContainer = document.querySelector('.categories');

    await loadData(1, true);
    for (const category of data.allCategories) {
        const categoryElem = document.createElement('li')
        categoryElem.id = category;
        categoryElem.onclick = () => filterByCategory(category);
        categoryElem.innerHTML = `
        <i class="hash-icon"></i>${category}
        `;
        categContainer.append(categoryElem);
    }
    displayPosts(data.allPosts);
    document.getElementById('select_1').classList.add('active');
};

const displayPosts = (posts) => {
    const main = document.querySelector('.main');
    const postsContainer = document.createElement('div');
    postsContainer.classList.add('posts-container');
    disactive();
    main.innerHTML = '';
    if (!posts.length) {
        main.innerHTML += `
        <img id="no_data" src="/assets/img/no_data.svg" alt="no result"/>
        `;
    } else {
        for (const post of posts) {
            const postDiv = createPostElement(post);
            postsContainer.append(postDiv);
        }
        main.append(postsContainer);
    }

    initInfiniteScroll()
};

const recentPosts = async () => {
    data.currentView = 'recent';
    await loadData(1, true);
    widgetBack();
    displayPosts(data.allPosts);
    document.getElementById('select_1').classList.add('active');
};

const createdPosts = async () => {
    data.currentView = 'created';
    await loadData(1, true);
    widgetBack();
    displayPosts(data.allPosts);
    document.getElementById('select_2').classList.add('active');
};

const likedPosts = async () => {
    data.currentView = 'liked';
    await loadData(1, true);
    widgetBack();
    displayPosts(data.allPosts);
    document.getElementById('select_3').classList.add('active');
};

const filterByCategory = async (category) => {
    data.currentView = 'category';
    data.currentCategory = category;
    await loadData(1, true);
    widgetBack();
    displayPosts(data.allPosts);
    document.getElementById(category).classList.add('activeCat');
};

const createPostElement = (post) => {
    const userId = parseInt(getCookie("user_id"));
    const postDiv = document.createElement("div");
    const likeActive = post.likes.includes(userId) ? ' liked' : ''
    const dislikeActive = post.dislikes.includes(userId) ? ' disliked' : ''
    postDiv.dataset.id = post.id;
    postDiv.classList.add("post");
    postDiv.innerHTML = `
    <div class="user-info">
        <img src="https://ui-avatars.com/api/?name=${post.by}" alt="User avatar" class="avatar">
        <div>
            <div class="username">${post.by}</div>
            <div class="timestamp">${timeAgo(new Date(post.createdAt).getTime())}</div>
        </div>
    </div>
    <div class="post-content">
        <h3 onclick="openPost(${post.id})">${filterContent(post.title)}</h3>
        <p>${displayContent(post.content, post.id)}</p>
    </div>
    <div class="tags-stats">
        <div class="tags">
            ${post.categories.map(tag => `<span class="tag">${tag}</span>`).join('')}
        </div>
        <div class="post-stats">
            <div class="stat${likeActive}">
                <i class="like-icon" onclick="likeAction(${post.id}, true)"></i><span>${post.likes.length}</span>
            </div>
            <div class="stat${dislikeActive}">
                <i class="dislike-icon dislike" onclick="dislikeAction(${post.id}, true)"></i><span>${post.dislikes.length}</span>
            </div>
            <div class="stat">
                <i class="comment-icon" onclick="openPost(${post.id})"></i><span>${post.totalComments}</span>
            </div>
        </div>
    </div>
    `
    return postDiv;
}

const displayContent = (content, id) => {
    if (content.match(/\n/g) && content.match(/\n/g).length < 3 && content.length < 200) {
        return filterContent(content)
    } else if (content.match(/\n/g) && content.match(/\n/g).length < 2 && content.length > 200) {
        return filterContent(content).slice(0, 200) + `... <span class="read-more" onclick="openPost(${id})">Read-More</span>`
    } else if (content.match(/\n/g) && content.match(/\n/g).length > 1) {
        return filterContent(content.split('\n')[0]).slice(0, 200) + `... <span class="read-more" onclick="openPost(${id})">Read-More</span>`
    } else {
        return filterContent(content)
    }
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
        <p>${filterContent(comment.content)}</p>
    </div>
    <div class="tags-stats">
        <div class="post-stats">
            <div class="stat${likeActive}">
                <i class="like-icon" onclick="likeAction(${comment.id}, false)"></i><span>${comment.likes.length}</span>
            </div>
            <div class="stat${dislikeActive}">
                <i class="dislike-icon dislike" onclick="dislikeAction(${comment.id}, false)"></i><span>${comment.dislikes.length}</span>
            </div>
        </div>
    </div>
    `
    return commentDiv;
}

const getPostId = () => {
    return document.querySelector('.post').dataset.id
}

const closedPost = () => {
    return document.querySelector('.comment-box') == undefined
}

const likeAction = async (id, isPost) => {
    const reqData = isPost ? { postId: id, isLike: true } : { commentId: id, isLike: true }
    try {
        const response = await fetch("/reaction", {
            method: "POST",
            body: JSON.stringify(reqData)
        })
        if (response.ok) {
            if (isPost && closedPost()) {
                const post = await getPostData(id);
                const postDiv = document.querySelector(`.post[data-id="${id}"]`)
                postDiv.innerHTML = createPostElement(post).innerHTML;
            } else {
                const postId = getPostId();
                openPost(postId);
            }
        } else {
            const res = await response.json();
            console.error(res);
            document.getElementById("loginPopup").style.display = "block";
        }
    } catch (err) {
        console.error(err)
    }
}

const dislikeAction = async (id, isPost) => {
    const reqData = isPost ? { postId: id, isDislike: true } : { commentId: id, isDislike: true }
    try {
        const response = await fetch("/reaction", {
            method: "POST",
            body: JSON.stringify(reqData)
        })
        if (response.ok) {
            if (isPost && closedPost()) {
                const post = await getPostData(id);
                const postDiv = document.querySelector(`.post[data-id="${id}"]`)
                postDiv.innerHTML = createPostElement(post).innerHTML;
            } else {
                const postId = getPostId();
                openPost(postId);
            }
        } else {
            const res = await response.json();
            console.error(res);
            document.getElementById("loginPopup").style.display = "block";
        }

    } catch (err) {
        console.error(err)
    }
}

const timeAgo = (time) => {
    const seconds = Math.floor(Date.now() - time) / 1000;
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

const escapeHtml = (unsafe) => {
    return unsafe
        .replace(/&/g, "&amp;")
        .replace(/</g, "&lt;")
        .replace(/>/g, "&gt;")
        .replace(/"/g, "&quot;")
        .replace(/'/g, "&#039;");
}

const filterContent = (content) => {
    return escapeHtml(content)
        .replace(/&lt;pre&gt;/g, '<pre>')
        .replace(/&lt;\/pre&gt;/g, '</pre>')
        .replace(/\n/g, '<br>')
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
            <i class="star-icon"></i>
            Must-read posts
        </h2>
        <ul>
            <li><a href="#" target="_blank">Please read rules before you start using our platform</a></li>
            <li><a href="https://www.paypal.com/paypalme/outiskteanas" target="_blank">Donate for 01Forum</a></li>
        </ul>
    </div>
    <div class="section">
        <h2 class="section-title">
            <i class="link-icon"></i>
            Featured links
        </h2>
        <ul>
            <li><a href="https://github.com/ANAS-OU/go_forum" target="_blank">01Forum source-code on GitHub</a></li>
            <li><a href="https://medium.com/@golangda/golang-quick-reference-top-20-best-coding-practices-c0cea6a43f20" target="_blank">Golang best-practices</a></li>
            <li><a href="https://zone01oujda.ma/" target="_blank">Zone01-Oujda Company</a></li>
        </ul>
    </div>
    `
}

const logout = async () => {
    try {
        const response = await fetch('/auth/logout', {
            method: 'POST',
        });

        if (response.redirected) location.href = "/";
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
                <i class="image-icon"></i>Add Image
            </button>
            <button class="btn btn-publish">
            <i class="send-icon"></i>Publish
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

    const renderTags = () => {
        selectedTagsContainer.innerHTML = selectedTags.map(tag => `
            <span class="tag">
                ${tag}
                <button class="tag-remove" data-tag="${tag}">×</button>
            </span>
        `).join('');

        const filteredTags = tags.filter(tag =>
            tag.toLowerCase().includes(searchBox.value.toLowerCase())
        );

        optionsContainer.innerHTML = filteredTags.map(tag => `
            <div class="dropdown-item ${selectedTags.includes(tag) ? 'selected' : ''}" data-tag="${tag}">
                ${tag}
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
            selectedTags = [...tags];
        }
        renderTags();
    });

    newPostForm.addEventListener("submit", async (e) => {
        e.preventDefault();

        const title = document.querySelector('input[name="title"]').value;
        const content = document.querySelector('textarea[name="content"]').value;
        console.log(title, content)
        try {
            const response = await fetch('/newpost', {
                method: 'POST',
                body: JSON.stringify({
                    Title: title,
                    Content: content,
                    Categories: selectedTags
                })
            })
            if (response.ok) {
                location.href = "/";
            } else {
                const res = await response.json();
                console.error(res);
                document.getElementById("loginPopup").style.display = "block";
                document.querySelector(".popup-content").innerHTML = `
                <h2>Nice try!</h2>
                <ul>
                    <li>${res.msg}</li>
                </ul>
                `
            }
        } catch (err) {
            console.error(err);
        }
    });

    renderTags();
}

const loginPopup = document.getElementById("loginPopup");
window.onclick = function (event) {
    if (event.target == loginPopup) {
        loginPopup.style.display = "none";
    }
};

init();

console.log(`
╱╱╱╱╱╱╱╱╱╱╱╱╱╭━━━╮╭╮
╱╱╱╱╱╱╱╱╱╱╱╱╱┃╭━╮┣╯┃
╭━━━┳━━┳━╮╭━━┫┃┃┃┣╮┃
┣━━┃┃╭╮┃╭╮┫┃━┫┃┃┃┃┃┃
┃┃━━┫╰╯┃┃┃┃┃━┫╰━╯┣╯╰╮
╰━━━┻━━┻╯╰┻━━┻━━━┻━━╯
`)