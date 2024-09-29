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
        setCookie('recipls_token', refresh_token);

        console.log(status);
        fetchHomePage();
    } else {
        const errorData = await response.json();
        console.error('Login failed:', errorData);
    }
}

async function fetchHomePage() {
    try {
        const response = await fetchWithAuth("/web/home", { method: 'GET' });

        if (response.ok) {
            const html = await response.text();
            document.open();
            document.write(html);
            document.close();

            window.history.pushState({}, '', '/web/home');
        } else {
            console.error('Failed to fetch home page:', await response.text());
        }
    } catch (error) {
        console.error('Error fetching home page:', error);
    }
}