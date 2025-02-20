document.addEventListener('DOMContentLoaded', () => {
    const avatarModal = document.getElementById('avatarModal');
    const changeAvatarButton = document.getElementById('changeAvatarButton');
    const selectAvatarButton = document.getElementById('selectAvatarButton');
    const avatarGallery = document.getElementById('avatarGallery');
    const avatar = document.getElementById('avatar');
    const usernameDisplay = document.getElementById('usernameDisplay');
    const emailDisplay = document.getElementById('emailDisplay');
    const usernameInput = document.getElementById('usernameInput');
    const editUsernameButton = document.getElementById('editUsernameButton');
    const profileForm = document.getElementById('profileForm');
    const userId = document.cookie.replace(/(?:(?:^|.*;\s*)user_id\s*\=\s*([^;]*).*$)|^.*$/, "$1");
    const loaderContainer = document.querySelector('.loader-container');

    let selectedAvatarSrc = '';

    // Открыть модальное окно для выбора аватара
    changeAvatarButton.addEventListener('click', () => {
        avatarModal.style.display = 'block';
        loadAvatars();
    });

    // Закрыть модальное окно
    const closeModal = document.querySelector('.close');
    closeModal.addEventListener('click', () => {
        avatarModal.style.display = 'none';
    });

    // Выбрать аватар
    selectAvatarButton.addEventListener('click', async () => {
        const selectedAvatar = document.querySelector('.avatar-gallery img.selected');
        if (selectedAvatar) {
            selectedAvatarSrc = selectedAvatar.src;
            avatar.src = selectedAvatarSrc;
            avatarModal.style.display = 'none';

            // Отправить выбранный аватар на сервер
            loaderContainer.style.display = 'flex';
            await saveAvatar(selectedAvatarSrc, userId);
        }
    });

    async function saveAvatar(avatarSrc, userId) {
        const response = await fetch(`https://workly-production-8296.up.railway.app/updateAvatar?user_id=${userId}`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                avatar: avatarSrc,
            }),
        });

        if (response.ok) {
            console.log('Аватар обновлен успешно.');
        } else {
            console.log('Ошибка при обновлении аватара.');
        }
        loaderContainer.style.display = 'none';
    }

    function loadAvatars() {
        const avatars = [
            '1.png', '2.png', '3.png', '4.png', '5.png', '6.png', '7.png', '8.png', '9.png', '10.png',
            '11.png', '12.png', '13.png', '14.png', '15.png', '16.png', '17.png', '18.png', '19.png', '20.png'
        ];

        avatarGallery.innerHTML = '';

        avatars.forEach(fileName => {
            const img = document.createElement('img');
            img.src = `/imgs/profileIcons/${fileName}`;
            img.alt = 'Avatar';
            img.classList.add('avatar-thumbnail');
            img.addEventListener('click', () => {
                document.querySelectorAll('.avatar-gallery img').forEach(img => img.classList.remove('selected'));
                img.classList.add('selected');
            });
            avatarGallery.appendChild(img);
        });
    }

    async function getUserData(userId) {
        const response = await fetch(`https://workly-production-8296.up.railway.app/getUserData?user_id=${userId}`);
        if (response.ok) {
            const data = await response.json();
            if (data.Username) {
                usernameDisplay.textContent = data.Username;
                usernameInput.value = data.Username;
            }
            if (data.Email) {
                emailDisplay.textContent = data.Email;
            }
            if (data.Avatar) {
                avatar.src = data.Avatar;
            }
        }
        loaderContainer.style.display = 'none';
    }

    // Редактировать имя пользователя
    editUsernameButton.addEventListener('click', () => {
        usernameDisplay.classList.add('hidden');
        editUsernameButton.classList.add('hidden');
        usernameInput.classList.remove('hidden');
        usernameInput.focus();
    });

    // Сохранение изменений профиля
    profileForm.addEventListener('submit', async (event) => {
        event.preventDefault();

        const newUsername = usernameInput.value;

        const response = await fetch(`https://workly-production-8296.up.railway.app/updateUsername?user_id=${userId}`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                username: newUsername,
            }),
        });

        if (response.ok) {
            console.log('Имя пользователя обновлено успешно.');
            usernameDisplay.textContent = newUsername;
            usernameDisplay.classList.remove('hidden');
            editUsernameButton.classList.remove('hidden');
            usernameInput.classList.add('hidden');
        } else {
            console.log('Ошибка при обновлении имени пользователя.');
        }
        loaderContainer.style.display = 'none';
    });
    loaderContainer.style.display = 'flex';
    // Загрузить текущий аватар и данные пользователя при загрузке страницы
    getUserData(userId);

    // Заглушка для уведомлений
    const notificationList = document.getElementById('notificationList');
    notificationList.innerHTML = '<li class="notification-item">Нет новых уведомлений.</li>';
});
