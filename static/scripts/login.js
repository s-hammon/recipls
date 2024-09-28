document.getElementById("loginForm").addEventListener('submit', function(event) {
    event.preventDefault();

    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;

    const data = {
        email: email,
        password: password
    };

    fetch('/v1/login', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(data),
    })
    .then(response => {
        if (response.ok) {
            return response.json();
        }
        throw new Error('login failed');
    })
    .then(data => {
        console.log('Success:', data);
        window.location.href = '/home';
    })
    .catch((error) => {
        console.error('Error:', error);
    })
});