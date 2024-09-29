window.addEventListener('popstate', function (event) {
    redirectToHomePage();
})

document.getElementById("cancelButton").addEventListener("click", redirectToHomePage);

document.getElementById("recipeForm").addEventListener("submit", function(event) {
    event.preventDefault();

    const title = document.getElementById('title').value;
    const description = document.getElementById('description').value;
    const difficulty = document.getElementById('difficulty').value;
    const ingredients = document.getElementById('ingredients').value;
    const instructions = document.getElementById('instructions').value;
    const category = document.getElementById('category').value;

    uploadRecipe({ title, description, difficulty, ingredients, instructions, category });
});

async function uploadRecipe(recipe) {
    const response = await fetchWithAuth('/v1/recipes', {
        method: 'POST',
        body: JSON.stringify(recipe),
    });

    if (response.ok) {
        window.location.href = "/home"
    } else {
        const errorData = await response.json();
        console.error('Upload failed:', errorData);
    }
}

function redirectToHomePage() {
    window.location.href = '/home';
    window.history.pushState({}, '', '/home');
}