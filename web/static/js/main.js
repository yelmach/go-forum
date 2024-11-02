const loadData = (posts) => {
    const main = document.querySelector('.main');
    const postsContainer = document.createElement('div');
    postsContainer.classList.add('posts-container');

    for (const post of posts) {
        const postDiv = createPostElement(post);
        postsContainer.appendChild(postDiv);
    }
    main.appendChild(postsContainer);
};

fetch("http://localhost:8080/api/posts")
    .then((response) => response.json()) // parse the response from JSON
    .then(loadData); // .then will call the `loadData` function with the JSON value.

const createPostElement = (post) => {
    const postDiv = document.createElement("div");
    postDiv.classList.add("post");
    

    // api avatar
    // https://ui-avatars.com/api/?name=${post.by}
    postDiv.innerHTML = `
    <div class="user-info">
        <img src="https://xsgames.co/randomusers/assets/avatars/male/${post.id}.jpg" alt="User avatar" class="avatar">
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
                <i class="ri-thumb-up-line"></i><span>${post.likes ? post.likes.length : 0}</span>
            </div>
            <div class="stat">
                <i class="ri-thumb-down-line"></i><span>${post.dislikes ? post.dislikes.length : 0}</span>
            </div>
            <div class="stat">
                <i class="ri-chat-4-line" on-click="openPost(${post.id})"></i><span>${post.comments ? post.comments.length : 0}</span>
            </div>
        </div>
    </div>
    `
    return postDiv;
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

const openPost = (postId) => {
    const post = document.createElement('div');
}

/*
<div class="post">
    <div class="user-info">
        <img src="/assets/img/usrs/Ava01.png" alt="User avatar" class="avatar">
        <div>
            <div class="username">Golanginya</div>
            <div class="timestamp">5 min ago</div>
        </div>
    </div>
    <div class="post-content">
        <h3>How to patch KDE on FreeBSD?</h3>
        <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Consequat aliquet maecenas ut sit nulla
        </p>
    </div>
    <div class="tags-stats">
        <div class="tags">
            <span class="tag">golang</span>
            <span class="tag">linux</span>
            <span class="tag">overflow</span>
        </div>
        <div class="post-stats">
            <div class="stat">
                <i class="ri-thumb-up-line"></i><span>125</span>
            </div>
            <div class="stat">
                <i class="ri-thumb-down-line"></i><span>90</span>
            </div>
            <div class="stat">
                <i class="ri-chat-4-line"></i><span>21</span>
            </div>
        </div>
    </div>
</div>
*/

    
console.log(`
╱╱╱╱╱╱╱╱╱╱╱╱╱╭━━━╮╭╮
╱╱╱╱╱╱╱╱╱╱╱╱╱┃╭━╮┣╯┃
╭━━━┳━━┳━╮╭━━┫┃┃┃┣╮┃
┣━━┃┃╭╮┃╭╮┫┃━┫┃┃┃┃┃┃
┃┃━━┫╰╯┃┃┃┃┃━┫╰━╯┣╯╰╮
╰// ━━━┻━━┻╯╰┻━━┻━━━┻━━╯
`)