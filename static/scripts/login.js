document.getElementById("loginForm").addEventListener('submit', function(event) {
    event.preventDefault();

    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;

    loginUser(email, password);
});

async function loginUser(email, password) {
    const data = { email, password };

    const response = await fetch('/v1/login', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(data),
    })

    if (response.ok) {
        const { id, status, access_token, refresh_token } = await response.json();

        setAccessToken(access_token);
        localStorage.setItem('refreshToken', refresh_token);

        console.log(`ID: ${id}, Status: ${status}`)
        window.location.href = '/home';
    } else {
        const errorData = await response.json();
        console.error('Login failed:', errorData);
    }
}

let currentAccessToken = '';

function setAccessToken(token) {
    currentAccessToken = token;
}