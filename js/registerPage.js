document.addEventListener('DOMContentLoaded', () => {
    const form = document.querySelector('#loginForm');
    const registerButton = document.querySelector('#registerButton');

    if (!form || !registerButton) {
        console.error('Form or buttons not found');
        return;
    }

    form.addEventListener('submit', async (event) => {
        event.preventDefault();

        const email = document.querySelector('#email').value;
        const password = document.querySelector('#password').value;
        const username = document.querySelector('#text').value;

        if (event.submitter === registerButton) {
            // Handle registration
            await handleRegistration(email, password, username);
        }
    });

    async function handleRegistration(email, password, username) {
        try {
            const response = await fetch('http://localhost:8080/register', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ email, password, username })
            });

            if (response.ok) {
                alert('Registration successful');
            } else {
                alert('Registration failed');
            }
        } catch (error) {
            console.error('Error:', error);
        }
    }
    //Включение\выключение видимости пароля
    const togglePassword = document.getElementById('toggle-password');
    const passwordField = document.getElementById('password');

    togglePassword.addEventListener('click', function () {
        const type = passwordField.getAttribute('type') === 'password' ? 'text' : 'password';
        passwordField.setAttribute('type', type);
        this.textContent = type === 'password' ? 'visibility' : 'visibility_off';
    });

});
