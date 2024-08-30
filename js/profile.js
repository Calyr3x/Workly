document.addEventListener('DOMContentLoaded', () => {
    const avatarModal = document.getElementById('avatarModal');
    const changeAvatarButton = document.getElementById('changeAvatarButton');
    const selectAvatarButton = document.getElementById('selectAvatarButton');
    const avatarGallery = document.getElementById('avatarGallery');
    const avatar = document.getElementById('avatar');
    const closeModal = document.querySelector('.close');
    const saveNameBtn = document.getElementById('saveNameBtn');
    const teamForm = document.getElementById('teamForm');

    // Открыть модальное окно для выбора аватара
    changeAvatarButton.addEventListener('click', () => {
        avatarModal.style.display = 'block';
        loadAvatars();
    });

    // Закрыть модальное окно
    closeModal.addEventListener('click', () => {
        avatarModal.style.display = 'none';
    });

    // Выбрать аватар
    selectAvatarButton.addEventListener('click', () => {
        const selectedAvatar = document.querySelector('.avatar-gallery img.selected');
        if (selectedAvatar) {
            avatar.src = selectedAvatar.src;
            avatarModal.style.display = 'none';
        }
    });

    function loadAvatars() {
        const avatars = [
            '1.png', // Предполагается, что файлы в папке imgs/profileIcons
            '2.png',
            '3.png',
            '4.png'
            // Добавьте все имена файлов, доступных для выбора
        ];

        avatarGallery.innerHTML = ''; // Очистить галерею перед добавлением новых изображений

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

    // Сохранение имени пользователя
    saveNameBtn.addEventListener('click', () => {
        const userNameInput = document.getElementById('userName');
        const userName = userNameInput.value.trim();
        if (userName) {
            alert(`Имя сохранено: ${userName}`);
            // Логика сохранения имени на сервере
        } else {
            alert('Имя не может быть пустым.');
        }
    });

    // Обработка формы создания команды
    teamForm.addEventListener('submit', (event) => {
        event.preventDefault();
        const teamName = document.getElementById('teamName').value.trim();
        const teamMembers = document.getElementById('teamMembers').value.trim();

        if (teamName && teamMembers) {
            alert(`Команда "${teamName}" создана с участниками: ${teamMembers}`);
            // Логика создания команды на сервере
        } else {
            alert('Пожалуйста, заполните все поля.');
        }
    });

    // Заглушка для уведомлений
    const notificationList = document.getElementById('notificationList');
    notificationList.innerHTML = '<li class="notification-item">Нет новых уведомлений.</li>'
});
