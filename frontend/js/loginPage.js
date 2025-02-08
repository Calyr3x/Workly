document.addEventListener('DOMContentLoaded', () => {
    const form = document.querySelector('#loginForm');
    const loginButton = document.querySelector('#loginButton');
    const registerButton = document.querySelector('#registerButton');
    const loaderContainer = document.querySelector('.loader-container');

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
            loaderContainer.style.display = 'flex';
            await handleLogin(email, password);
        }
    });

    async function handleLogin(email, password) {
        try {
            const response = await fetch('http://workly-production-8296.up.railway.app:8080/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ email, password })
            });

            if (response.ok) {
                // Получаем user_id из ответа сервера
                const data = await response.json();
                const userId = data.user_id;

                // Усттанавливаем куки
                document.cookie = `user_id=${userId}; path=/;`;
                loaderContainer.style.display = 'none';
                // Успешный вход, перенаправляем на главную страницу
                window.location.href = 'main.html';
            } else {
                loaderContainer.style.display = 'none';
                // Неудачный вход
                const form = document.querySelector('.card');
                const errorMessage = document.querySelector('.error-message');

                // Показываем текст ошибки
                errorMessage.textContent = 'Неверный email или пароль. Попробуйте снова.';
                errorMessage.style.display = 'block';

                // Добавляем класс тряски и изменяем цвет границы
                form.classList.add('shake-animation');
                form.style.borderColor = 'red';

                // Убираем анимацию и восстанавливаем цвет границы через 2 секунды
                setTimeout(() => {
                    form.classList.remove('shake-animation');
                    form.style.borderColor = ''; // Вернуть к исходному цвету
                }, 2000);
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

    //Редирект на страницу регистрации
    registerButton.addEventListener('click', function () {
        window.location.href = 'register.html';
    });
});
