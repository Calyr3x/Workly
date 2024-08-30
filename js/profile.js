// profile.js

document.addEventListener('DOMContentLoaded', () => {
    const avatarImage = document.getElementById('userAvatar');
    const changeAvatarBtn = document.getElementById('changeAvatarBtn');
    const saveNameBtn = document.getElementById('saveNameBtn');
    const teamForm = document.getElementById('teamForm');

    // Функция для смены аватара
    changeAvatarBtn.addEventListener('click', () => {
        // Заменим на логику выбора аватара (например, случайный выбор из массива)
        const avatars = [
            '/imgs/profileIcons/1.png',
            '/imgs/profileIcons/2.png',
            '/imgs/profileIcons/3.png',
            '/imgs/profileIcons/4.png',
            '/imgs/profileIcons/5.png',
            '/imgs/profileIcons/6.png',
            '/imgs/profileIcons/7.png'
        ];
        const randomAvatar = avatars[Math.floor(Math.random() * avatars.length)];
        avatarImage.src = randomAvatar;
    });

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
    notificationList.innerHTML = '<li class="notification-item">Нет новых уведомлений.</li>';
});
