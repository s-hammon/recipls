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

        localStorage.setItem('access_token', access_token);
        document.cookie = `refresh_token=${refresh_token}; path=/; SameSite=Strict`;

        fetchHomePage();
    } else {
        const errorData = await response.json();
        console.error('Login failed:', errorData);
    }
}

async function fetchHomePage() {
    const access_token = localStorage.getItem('access_token');

    try {
        const response = await fetch("/home", {
            method: "GET",
            headers: {
                'Authorization': `Bearer ${access_token}`,
                'Content-Type': 'application/json'
            },
            credentials: 'include'
        });

        if (response.ok) {
            const html = await response.text();
            document.open();
            document.write(html);
            document.close();
        } else {
            console.error('Failed to fetch home page:', await response.text());
        }
    } catch (error) {
        console.error('Error fetching home page:', error);
    }
}

function getCookie(name) {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) return parts.pop().split(';').shift();
}