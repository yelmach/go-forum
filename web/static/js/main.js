let POSTS = {};
const loadData = (posts) => {
    POSTS = [...posts]
    const main = document.querySelector('.main');
    const postsContainer = document.createElement('div');
    postsContainer.classList.add('posts-container');

    // posts = posts.slice(0, 10);
    for (const post of posts) {
        const postDiv = createPostElement(post);
        postsContainer.append(postDiv);
    }
    main.append(postsContainer);
};

fetch("http://localhost:8080/api/posts")
    .then((response) => response.json()) // parse the response from JSON
    .then(loadData); // .then will call the `loadData` function with the JSON value.

// const generateAvatar = async () => {
//     const [color, background] = await fetch("https://random-flat-colors.vercel.app/api/random?count=2")
//         .then(response => response.json())
//         .then(json => json.colors)

//     return [color, background]
// }

const createPostElement = (post) => {
    const userId = parseInt(getCookie("user_id"));
    const postDiv = document.createElement("div");
    postDiv.classList.add("post");
    postDiv.onclick = () => openPost(post.id);
    // const [color, background] = await generateAvatar();
    const likeActive = post.likes.includes(userId) ? ' activeThumb' : ''
    const dislikeActive = post.dislikes.includes(userId) ? ' activeThumb' : ''

    postDiv.innerHTML = `
    <div class="user-info">
        <img src="https://ui-avatars.com/api/?name=${post.by}" alt="User avatar" class="avatar">
        <div>
            <div class="username">${post.by}</div>
            <div class="timestamp">${timeAgo(new Date(post.createdAt).getTime())}</div>
        </div>
    </div>
    <div class="post-content">
        <h3>${post.title}</h3>
        <p>${post.content}</p>
    </div>
    <div class="tags-stats">
        <div class="tags">
            ${post.categories.map(tag => `<span class="tag">${tag}</span>`).join('\n')}
        </div>
        <div class="post-stats">
            <div class="stat">
                <i class="ri-thumb-up-line${likeActive}"></i><span class="${likeActive}">${post.likes.length}</span>
            </div>
            <div class="stat">
                <i class="ri-thumb-down-line${dislikeActive}"></i><span class="${dislikeActive}">${post.dislikes.length}</span>
            </div>
            <div class="stat">
                <i class="ri-chat-4-line"></i><span>${post.comments ? post.comments.length : 0}</span>
            </div>
        </div>
    </div>
    `
    return postDiv;
}

const createCommentElement = (comment) => {
    const commentDiv = document.createElement('div');
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
                <i class="ri-thumb-up-line"></i><span>${comment.likes.length}</span>
            </div>
            <div class="stat">
                <i class="ri-thumb-down-line"></i><span>${comment.dislikes.length}</span>
            </div>
        </div>
    </div>
    `
    return commentDiv;
}

getPostData = async (postId) => {
    try {
        const response = await fetch(`http://localhost:8080/api/posts/${postId}`);
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
                <svg class="comment-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z" />
                </svg>
                Comment
            </button>
        </div>
    </div>
    `

    widget.innerHTML = `
    <img src="https://ui-avatars.com/api/?name=${post.by}" alt="User avatar">
    <p class="username">@${post.by}</p>
    `
    
    document.querySelector('.btn-cancel').addEventListener('click', function () {
        document.querySelector('.comment-input').value = '';
    });

    document.querySelector('.btn-comment').addEventListener('click', function () {
        const comment = document.querySelector('.comment-input').value;
        console.log('Comment submitted:', comment);
        document.querySelector('.comment-input').value = '';
    });
    main.append(comments);
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

const disactive = (elems) => {
    for (const elem of elems) {
        elem.classList.remove('active');
    }
}

const recentPosts = () => {
    location.reload();
}

const createdPosts = () => {
    const username = getCookie("username");
    const posts = POSTS.filter(post => post.by == username)
    const main = document.querySelector('.main');
    const postsContainer = document.createElement('div');
    postsContainer.classList.add('posts-container');
    disactive(document.querySelectorAll('.menu-select'))
    document.getElementById('select_2').classList.add('active')

    for (const post of posts) {
        const postDiv = createPostElement(post);
        postsContainer.append(postDiv);
    }
    main.innerHTML = ''
    main.append(postsContainer)
}

const likedPosts = () => {
    const userId = parseInt(getCookie("user_id"));
    const posts = POSTS.filter(post => post.likes.includes(userId))
    const main = document.querySelector('.main');
    const postsContainer = document.createElement('div');
    postsContainer.classList.add('posts-container');
    disactive(document.querySelectorAll('.menu-select'))
    document.getElementById('select_3').classList.add('active')

    for (const post of posts) {
        const postDiv = createPostElement(post);
        postsContainer.append(postDiv);
    }
    main.innerHTML = ''
    main.append(postsContainer)
}

console.log(`
╱╱╱╱╱╱╱╱╱╱╱╱╱╭━━━╮╭╮
╱╱╱╱╱╱╱╱╱╱╱╱╱┃╭━╮┣╯┃
╭━━━┳━━┳━╮╭━━┫┃┃┃┣╮┃
┣━━┃┃╭╮┃╭╮┫┃━┫┃┃┃┃┃┃
┃┃━━┫╰╯┃┃┃┃┃━┫╰━╯┣╯╰╮
╰━━━┻━━┻╯╰┻━━┻━━━┻━━╯
`)