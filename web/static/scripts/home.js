document.getElementById('logoutButton').addEventListener('click', function() {
    logout();
});

document.getElementById('createRecipeButton').addEventListener('click', function() {
    fetchNewRecipePage();
})

async function fetchNewRecipePage() {
    try {
        const response = await fetchWithAuth('/web/recipes/new', { method: 'GET' });

        if (response.ok) {
            const html = await response.text();
            document.open();
            document.write(html);
            document.close();
        
            window.history.pushState({}, '', '/web/recipes/new');
        } else {
            console.error('Failed to fetch home page:', await response.text());
        }
    } catch (error) {
        console.error('Error fetching home page:', error);
    }
}