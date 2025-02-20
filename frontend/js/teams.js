document.addEventListener('DOMContentLoaded', () => {
    // Элементы модального окна создания команды
    const createTeamModal = document.getElementById('createTeamModal');
    const createTeamForm = document.getElementById('createTeamForm');
    const teamNameInput = document.getElementById('teamName');
    const teamMemberInput = document.getElementById('teamMemberInput');
    const addMemberButton = document.getElementById('addMemberButton');
    const teamMembersList = document.getElementById('teamMembersList');
    const errorMessage = document.getElementById('errorMessage');
    const successMessage = document.getElementById('successMessage');
    const successMembersList = document.getElementById('successMembersList');
    const userId = document.cookie.replace(/(?:(?:^|.*;\s*)user_id\s*\=\s*([^;]*).*$)|^.*$/, "$1");

    const teamMembers = [];

    // Модальное окно успешного создания команды
    const successModal = document.getElementById('successModal');
    const closeSuccessButton = document.getElementById('closeSuccessButton');

    // Открыть модальное окно создания команды
    const createTeamButton = document.getElementById('createTeamButton');
    createTeamButton.addEventListener('click', () => {
        createTeamModal.style.display = 'block';
    });

    // Закрыть модальное окно создания команды
    const closeCreateTeamModal = createTeamModal.querySelector('.close');
    closeCreateTeamModal.addEventListener('click', () => {
        createTeamModal.style.display = 'none';
        resetCreateTeamForm();
    });

    // Закрыть модальное окно успешного создания команды
    closeSuccessButton.addEventListener('click', () => {
        successModal.style.display = 'none';
    });

    // Добавление участника в список
    addMemberButton.addEventListener('click', async () => {
        const username = teamMemberInput.value.trim();
        errorMessage.style.display = 'none';

        if (!username) return;

        const response = await fetch(`https://workly-production-8296.up.railway.app/getUserAvatar?username=${username}`);
        if (response.ok) {
            const data = await response.json();
            if (teamMembers.includes(username)) {
                errorMessage.textContent = 'Этот участник уже добавлен.';
                errorMessage.style.display = 'block';
                return;
            }

            teamMembers.push(username);

            const listItem = document.createElement('li');
            const avatarImg = document.createElement('img');
            avatarImg.src = data.avatar;
            avatarImg.alt = 'Avatar';
            avatarImg.classList.add('avatar-thumbnail');
            listItem.appendChild(avatarImg);
            listItem.appendChild(document.createTextNode(username));
            teamMembersList.appendChild(listItem);

            teamMemberInput.value = '';
            errorMessage.style.display = 'none';
        } else {
            errorMessage.textContent = 'Пользователь не найден.';
            errorMessage.style.display = 'block';
        }
    });

    // Отправка формы создания команды
    createTeamForm.addEventListener('submit', async (event) => {
        event.preventDefault();

        const teamName = teamNameInput.value.trim();

        if (teamName && teamMembers.length > 0) {
            const response = await fetch(`https://workly-production-8296.up.railway.app/createTeam?user_id=${userId}`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    name: teamName,
                    members: teamMembers,
                }),
            });

            if (response.ok) {
                console.log('Команда успешно создана.');
                createTeamModal.style.display = 'none';
                displaySuccessModal(teamName, teamMembers);
            } else {
                console.error('Ошибка при создании команды.');
            }
        }
        resetCreateTeamForm()
    });

    // Функция отображения модального окна успешного создания команды
    async function displaySuccessModal(teamName, members) {
        successMessage.textContent = `Команда "${teamName}" создана!`;
        successMembersList.innerHTML = '';

        // Создаем массив промисов для загрузки аватаров всех участников
        const avatarPromises = members.map(async (username) => {
            const avatarUrl = await getUserAvatar(username);
            return { username, avatarUrl };
        });

        // Ждем, пока все аватары загрузятся
        const memberAvatars = await Promise.all(avatarPromises);

        // Добавляем каждого участника в список
        memberAvatars.forEach(({ username, avatarUrl }) => {
            const listItem = document.createElement('li');
            const avatarImg = document.createElement('img');
            avatarImg.src = avatarUrl;
            avatarImg.alt = 'Avatar';
            avatarImg.classList.add('avatar-thumbnail');
            listItem.appendChild(avatarImg);
            listItem.appendChild(document.createTextNode(username));
            successMembersList.appendChild(listItem);
        });

        successModal.style.display = 'block';
    }


    async function getUserAvatar(username) {
        const response = await fetch(`https://workly-production-8296.up.railway.app/getUserAvatar?username=${username}`);
        if (response.ok) {
            const data = await response.json();
            return data.avatar;
        }
        return ''; // Возвращаем пустую строку, если аватар не найден
    }

    // Функция сброса формы создания команды
    function resetCreateTeamForm() {
        teamNameInput.value = '';
        teamMemberInput.value = '';
        teamMembersList.innerHTML = '';
        teamMembers.length = 0;
        errorMessage.style.display = 'none';
    }
});
