document.addEventListener('DOMContentLoaded', () => {
    const teamsList = document.getElementById('teamsList');
    const userId = document.cookie.replace(/(?:(?:^|.*;\s*)user_id\s*\=\s*([^;]*).*$)|^.*$/, "$1");

    // Функция для получения команд пользователя
    async function getTeams(userId) {
        const response = await fetch(`http://localhost:8080/getTeams?user_id=${userId}`);
        if (response.ok) {
            const teams = await response.json();
            displayTeams(teams);
        } else {
            console.error('Ошибка при получении команд');
        }
    }

    // Функция для получения аватара пользователя
    async function getUserAvatar(username) {
        const response = await fetch(`http://localhost:8080/getUserAvatar?username=${username}`);
        if (response.ok) {
            const data = await response.json();
            return data.avatar;
        }
        return ''; // Возвращаем пустую строку, если аватар не найден
    }

    // Функция для отображения команд на странице
    async function displayTeams(teams) {
        teamsList.innerHTML = '';  // Очищаем список перед добавлением новых данных
        for (const team of teams) {
            const teamElement = document.createElement('li');
            teamElement.classList.add('team-item');
            const isOwner = team.owner_id === userId; // Проверка, является ли пользователь владельцем

            let teamHTML = `<h3>${team.name}</h3>`;
            teamHTML += '<ul class="team-members">';

            for (const member of team.members) {
                const avatar = await getUserAvatar(member); // Получаем аватар пользователя
                teamHTML += `
                        <li>
                            <img src="${avatar}" alt="${member}'s avatar" onerror="this.src='default-avatar.png';" />
                            ${member}
                            ${isOwner && member !== userId ? `<button class="remove-member" data-team-id="${team.id}" data-member="${member}">Удалить</button>` : ''}
                        </li>`;
            }

            teamHTML += '</ul>';
            if (isOwner) {
                teamHTML += `<button class="add-member" data-team-id="${team.id}">Добавить участника</button>`;
            }

            teamElement.innerHTML = teamHTML;
            teamsList.appendChild(teamElement);
        }

        // Добавляем обработчики для кнопок добавления/удаления участников
        document.querySelectorAll('.remove-member').forEach(button => {
            button.addEventListener('click', handleRemoveMember);
        });

        document.querySelectorAll('.add-member').forEach(button => {
            button.addEventListener('click', handleAddMember);
        });
    }

    // Удалить участника
    async function handleRemoveMember(event) {
        const teamId = event.target.dataset.teamId;
        const member = event.target.dataset.member;

        const response = await fetch('http://localhost:8080/removeMember', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                TeamID: teamId,
                Member: member,
            }),
        });

        if (response.ok) {
            event.target.closest('li').remove();
        } else {
            alert('Ошибка при удалении участника');
        }
    }

    // Добавить участника
    async function handleAddMember(event) {
        const teamId = event.target.dataset.teamId;
        const member = prompt('Введите юзернейм участника для добавления:');

        if (member) {
            const response = await fetch('http://localhost:8080/addMember', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    TeamID: teamId,
                    Member: member,
                }),
            });

            if (response.ok) {
                // Обновляем список команд после добавления участника
                getTeams(userId);
            } else {
                alert('Ошибка при добавлении участника');
            }
        }
    }

    // Загрузить команды при загрузке страницы
    getTeams(userId);
});
