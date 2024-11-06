let POSTS = {};
const loadData = async (posts) => {
    POSTS = [...posts];
    const categContainer = document.querySelector('.categories');
    const categories = await fetch("http://localhost:8080/api/categories")
        .then(response => response.json())

    for (const category of categories) {
        const categoryElem = document.createElement('li')
        categoryElem.id = category;
        categoryElem.onclick = () => filterByCategory(category);
        categoryElem.innerHTML = `
        <i class="ri-hashtag"></i>${category}
        `
        categContainer.append(categoryElem);
    }
    recentPosts();
};

fetch("http://localhost:8080/api/posts")
    .then((response) => response.json())
    .then(loadData);

// const generateAvatar = async () => {
//     const [color, background] = await fetch("https://random-flat-colors.vercel.app/api/random?count=2")
//         .then(response => response.json())
//         .then(json => json.colors)

//     return [color, background]
// }
// const [color, background] = await generateAvatar();

const createPostElement = (post) => {
    const userId = parseInt(getCookie("user_id"));
    const postDiv = document.createElement("div");
    postDiv.classList.add("post");
    const likeActive = post.likes.includes(userId) ? ' liked' : ''
    const dislikeActive = post.dislikes.includes(userId) ? ' disliked' : ''

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
                <i class="ri-thumb-down-line dislike${dislikeActive}"></i><span class="${dislikeActive}">${post.dislikes.length}</span>
            </div>
            <div class="stat">
                <i class="ri-chat-4-line" onclick="openPost(${post.id})"></i><span>${post.comments ? post.comments.length : 0}</span>
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
                <i class="ri-thumb-up-line${likeActive}"></i><span class="${likeActive}">${comment.likes.length}</span>
            </div>
            <div class="stat">
                <i class="ri-thumb-down-line dislike${dislikeActive}"></i><span class="${dislikeActive}">${comment.dislikes.length}</span>
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
        <form action="/" method="POST">
            <textarea class="comment-input" name="content" placeholder="Type here your wise suggestion"></textarea>
            <input type="hidden" name="postId" value="${postId}">
            <div class="button-group">
                <button class="btn btn-cancel">Cancel</button>
                <button type="submit" class="btn btn-comment">
                    <i class="ri-chat-4-line"></i>Comment
                </button>
            </div>
        </form>
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
        fetch('/', {
            method: 'POST',
            body: new URLSearchParams({
                content: content,
                postId: postId,
            }),
        }).catch(error => {
            console.error('Error posting comment:', error);
        });

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

const recentPosts = () => {
    widgetBack();
    displayPosts(POSTS);
    document.getElementById('select_1').classList.add('active');
}

const createdPosts = () => {
    const username = getCookie("username");
    const posts = POSTS.filter(post => post.by == username)
    widgetBack();
    displayPosts(posts);
    document.getElementById('select_2').classList.add('active');
}

const likedPosts = () => {
    const userId = parseInt(getCookie("user_id"));
    const posts = POSTS.filter(post => post.likes.includes(userId))
    widgetBack();
    displayPosts(posts);
    document.getElementById('select_3').classList.add('active');
}

const filterByCategory = (category) => {
    const posts = POSTS.filter(post => post.categories.includes(category))
    widgetBack();
    displayPosts(posts);
    document.getElementById(category).classList.add('activeCat');
}

const displayPosts = (posts) => {
    const main = document.querySelector('.main');
    const postsContainer = document.createElement('div');
    postsContainer.classList.add('posts-container');
    disactive();
    main.innerHTML = ''
    if (!posts.length) {
        // tmp figure for empty meaning
        main.innerHTML += `
        <img style="display: block; width:300px; margin: 3rem auto;" src="/assets/img/no_data.svg" alt="no result"/>
        <h2 style="text-align:center">Noting Found</h2>
        `
    } else {
        for (const post of posts) {
            const postDiv = createPostElement(post);
            postsContainer.append(postDiv);
        }
        main.append(postsContainer)
    }
}

console.log(`
╱╱╱╱╱╱╱╱╱╱╱╱╱╭━━━╮╭╮
╱╱╱╱╱╱╱╱╱╱╱╱╱┃╭━╮┣╯┃
╭━━━┳━━┳━╮╭━━┫┃┃┃┣╮┃
┣━━┃┃╭╮┃╭╮┫┃━┫┃┃┃┃┃┃
┃┃━━┫╰╯┃┃┃┃┃━┫╰━╯┣╯╰╮
╰━━━┻━━┻╯╰┻━━┻━━━┻━━╯
`)