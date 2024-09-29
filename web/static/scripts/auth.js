function getAccessToken() {
    return localStorage.getItem('access_token');
}

function isAuthenticated() {
    const token = getAccessToken();
    return token !== null;
}

async function fetchWithAuth(url, options= {}) {
    const access_token = localStorage.getItem('access_token');

    if (!access_token) {
        console.error("no access token found. Redirecting to login...");
        window.location.href = '/web/login';
        return;
    }

    options.headers = {
        ...options.headers,
        'Authorization': `Bearer ${access_token}`,
        'Content-Type': 'application/json',
    };

    console.log('Making request to:', url);
    console.log('Request headers:', options.headers);

    return fetch(url, options)
}

async function logout() {
    localStorage.removeItem('access_token');
    const refresh_token = getCookie('refresh_token');

    const response = await fetch('/v1/revoke', {
        method: 'POST',
        headers: {
            'Authorization': `Bearer ${refresh_token}`,
            'Content-Type': 'application/json',
        }
    });
    
    if (response.ok) {
        console.log("refresh token revoked by server");
        deleteCookie('recipls_token')
    } else {
        const errorData = await response.json();
        console.error('Login failed:', errorData);
    }
    
    console.log(`attempting to remove token: ${refresh_token}`);
    window.location.href = '/web/login';
}

function getCookie(name) {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) return parts.pop().split(';').shift();
}

function setCookie(name, value) {
    document.cookie = name +'='+ value +'; Path=/;';
}

function deleteCookie(name) {
    document.cookie = name +'=; Path=/; Expires=Thu, 01 Jan 1970 00:00:01 GMT;';
}