document.addEventListener('DOMContentLoaded', () => {
    const teamsList = document.getElementById('teamsList');
    const userId = document.cookie.replace(/(?:(?:^|.*;\s*)user_id\s*\=\s*([^;]*).*$)|^.*$/, "$1");
    const loaderContainer = document.querySelector('.loader-container');

    // Функция для получения команд пользователя
    async function getTeams(userId) {
        loaderContainer.style.display = 'flex';
        const response = await fetch(`https://workly-production-8296.up.railway.app/getTeams?user_id=${userId}`);
        if (response.ok) {
            const teams = await response.json();
            if (teams === null) {
                return [];
            }
            await displayTeams(teams);
        } else {
            console.error('Ошибка при получении команд');
        }
    }

    // Функция для получения аватара пользователя
    async function getUserAvatar(username) {
        const response = await fetch(`https://workly-production-8296.up.railway.app/getUserAvatar?username=${username}`);
        if (response.ok) {
            const data = await response.json();
            return data.avatar;
        }
        return ''; // Возвращаем пустую строку, если аватар не найден
    }

    async function getUserData(userId) {
        const response = await fetch(`https://workly-production-8296.up.railway.app/getUserData?user_id=${userId}`);
        if (response.ok) {
            const data = await response.json();
            if (data.Username) {
                return data.Username;
            }
        }
    }

    // Функция для отображения команд на странице
    async function displayTeams(teams) {
        teamsList.innerHTML = '';  // Очищаем список перед добавлением новых данных
        const username = await getUserData(userId);

        for (const team of teams) {
            const teamElement = document.createElement('li');
            teamElement.classList.add('team-item');
            const isOwner = team.OwnerID === userId; // Проверка, является ли пользователь владельцем

            // Получаем имя создателя команды для последующего отображения
            const ownerUsername = await getUserData(team.OwnerID);

            // Проверяем, если создатель не включён в список участников, добавляем его
            if (!team.Members.includes(ownerUsername)) {
                team.Members.unshift(ownerUsername);
            }

            // Создаём HTML для команды
            let teamHTML = `<h3>${team.Name}</h3>`;
            teamHTML += '<ul class="team-members">';

            // Рендерим каждого участника
            for (const member of team.Members) {
                const avatar = await getUserAvatar(member);
                teamHTML += `
                <li>
                    <img src="${avatar}" alt="${member}'s avatar" onerror="this.src='http://localhost:63342/frontend/imgs/profileIcons/default-avatar.png';" />
                    ${member} ${member === ownerUsername ? `(Создатель)` : ''}
                    ${isOwner && member !== username ? `<button class="remove-member" data-team-id="${team.ID}" data-member="${member}">Удалить</button>` : ''}
                </li>`;
            }

            teamHTML += '</ul>';

            // Если пользователь — создатель, добавляем кнопку для добавления участников
            if (isOwner) {
                teamHTML += `<button class="add-member" data-team-id="${team.ID}">Добавить участника</button>`;
            }

            teamElement.innerHTML = teamHTML;
            teamsList.appendChild(teamElement);
        }

        // Добавляем обработчики для кнопок удаления и добавления участников
        document.querySelectorAll('.remove-member').forEach(button => {
            button.addEventListener('click', handleRemoveMember);
        });

        document.querySelectorAll('.add-member').forEach(button => {
            button.addEventListener('click', handleAddMember);
        });

        loaderContainer.style.display = 'none';
    }


    // Удалить участника
    async function handleRemoveMember(event) {
        loaderContainer.style.display = 'flex';
        const teamId = Number(event.target.dataset.teamId);
        const member = event.target.dataset.member;

        const response = await fetch('https://workly-production-8296.up.railway.app/removeMember', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                team_id: teamId,
                Member: member,
            }),
        });

        if (response.ok) {
            getTeams(userId);
        } else {
            alert('Ошибка при удалении участника');
        }
    }

    // Открытие модального окна добавления участников
    async function handleAddMember(event) {
        const teamId = Number(event.target.dataset.teamId);
        document.getElementById('addMemberModal').style.display = 'block';
        const member = [];

        // Добавление участника в список
        document.getElementById('addMemberButton2').addEventListener('click', async () => {
            const username = document.getElementById('memberInput').value.trim();
            const errorMessage = document.getElementById('addMemberErrorMessage');
            errorMessage.style.display = 'none';

            if (!username) return;

            const response = await fetch(`https://workly-production-8296.up.railway.app/getUserAvatar?username=${username}`);
            if (response.ok) {
                const data = await response.json();

                if (member.includes(username)) {
                    errorMessage.textContent = 'Этот участник уже добавлен.';
                    errorMessage.style.display = 'block';
                    return;
                }

                member.push(username);

                const listItem = document.createElement('li');
                const avatarImg = document.createElement('img');
                avatarImg.src = data.avatar;
                avatarImg.alt = 'Avatar';
                avatarImg.classList.add('avatar-thumbnail');
                listItem.appendChild(avatarImg);
                listItem.appendChild(document.createTextNode(username));
                document.getElementById('addedMembersList').appendChild(listItem);

                document.getElementById('memberInput').value = '';
            } else {
                errorMessage.textContent = 'Пользователь не найден.';
                errorMessage.style.display = 'block';
            }
        });

        // Сохранение участников
        document.getElementById('saveMembersButton').addEventListener('click', async () => {
            if (member.length === 0) {
                alert('Добавьте хотя бы одного участника.');
                return;
            }

            if (member) {
                const response = await fetch('https://workly-production-8296.up.railway.app/addMember', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({
                        team_id: teamId,
                        Member: member,
                    }),
                });

                if (response.ok) {
                    // Обновляем список команд после добавления участника
                    document.getElementById('addMemberModal').style.display = 'none';  // Закрыть модальное окно
                } else {
                    alert('Ошибка при добавлении участника');
                }
            }
            // Загрузить команды при загрузке страницы
            getTeams(userId);
        });

        // Закрытие модального окна
        document.getElementById('closeAddMemberModal').onclick = () => {
            document.getElementById('addMemberModal').style.display = 'none';
        };
    }

    // Событие для открытия модального окна при клике на кнопку
    document.querySelectorAll('.add-member').forEach(button => {
        button.addEventListener('click', handleAddMember);
    });
    // Загрузить команды при загрузке страницы
    getTeams(userId);
});
