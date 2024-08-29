document.addEventListener('DOMContentLoaded', () => {
    const form = document.querySelector('#loginForm');
    const loginButton = document.querySelector('#loginButton');
    const registerButton = document.querySelector('#registerButton');

    if (!form || !loginButton || !registerButton) {
        console.error('Form or buttons not found');
        return;
    }

    form.addEventListener('submit', async (event) => {
        event.preventDefault();

        const email = document.querySelector('#email').value;
        const password = document.querySelector('#password').value;

        if (event.submitter === loginButton) {
            // Handle login
            await handleLogin(email, password);
        } else if (event.submitter === registerButton) {
            // Handle registration
            await handleRegistration(email, password);
        }
    });

    async function handleLogin(email, password) {
        try {
            const response = await fetch('http://localhost:8080/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ email, password })
            });

            if (response.ok) {
                // Parse response to get user ID
                const data = await response.json();
                const userId = data.user_id;

                // Set user_id cookie
                document.cookie = `user_id=${userId}; path=/;`;

                // Успешный вход, перенаправляем на главную страницу
                window.location.href = 'main.html';
            } else {
                // Неудачный вход
                const form = document.querySelector('.card');
                const errorMessage = document.querySelector('.error-message');

                // Показываем текст ошибки
                errorMessage.textContent = 'Неверный email или пароль. Попробуйте снова.';
                errorMessage.style.display = 'block';

                // Добавляем класс тряски и изменяем цвет границы
                form.classList.add('shake-animation');
                form.style.borderColor = 'red';

                // Убираем анимацию и восстанавливаем цвет границы через 1 секунду
                setTimeout(() => {
                    form.classList.remove('shake-animation');
                    form.style.borderColor = ''; // Вернуть к исходному цвету
                }, 2000);
            }
        } catch (error) {
            console.error('Error:', error);
        }
    }

    async function handleRegistration(email, password) {
        try {
            const response = await fetch('http://localhost:8080/register', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ email, password })
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
});
