document.addEventListener('DOMContentLoaded', () => {
    const addTaskButton = document.getElementById('addTaskButton');
    const taskModal = document.getElementById('taskModal');
    const taskViewModal = document.getElementById('taskViewModal');
    const closeModalButtons = document.querySelectorAll('.modal .close');
    const taskList = document.getElementById('taskList');
    const loaderContainer = document.querySelector('.loader-container');
    const isTeamTaskCheckbox = document.getElementById('isTeamTask');
    const teamSelection = document.getElementById('teamSelection');
    const memberSelection = document.getElementById('memberSelection');
    const teamSelect = document.getElementById('teamSelect');
    const memberSelect = document.getElementById('memberSelect');
    let editingTaskID = null;

    // Функция для открытия модального окна создания/редактирования задачи
    const openModal = (task = null) => {
        if (task) {
            const utcDate = new Date(task.Deadline);
            const localDate = new Date(utcDate.getTime() - utcDate.getTimezoneOffset() * 60000);

            document.getElementById('taskDeadline').value = localDate.toISOString().slice(0, 16);
            document.getElementById('taskName').value = task.Name;
            document.getElementById('taskDescription').value = task.Description;
            document.getElementById('taskStatus').value = task.Status;
            document.getElementById('statusContainer').style.display = 'block'; // Показываем статус задачи при редактировании
            editingTaskID = task.ID;
            document.getElementById('modalTitle').innerText = 'Редактировать задачу';
        } else {
            document.getElementById('taskName').value = '';
            document.getElementById('taskDescription').value = '';
            document.getElementById('taskDeadline').value = '';
            document.getElementById('statusContainer').style.display = 'none'; // Скрываем статус задачи при создании
            editingTaskID = null;
            document.getElementById('modalTitle').innerText = 'Создать задачу';
            document.getElementById('isTeamTask').value = false;
        }
        taskModal.style.display = 'block';
    };

    // Функция для открытия окна просмотра задачи
    const openTaskViewModal = (task) => {
        const deadline = new Date(task.Deadline);
        const createdAt = new Date(task.CreatedAt);
        const now = new Date();

        document.getElementById('viewTaskTitle').textContent = task.Name;
        document.getElementById('viewTaskDescription').textContent = task.Description;
        document.getElementById('viewTaskDeadline').querySelector('span').textContent = deadline.toLocaleString(undefined, {
            year: 'numeric',
            month: 'long',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit',
            hour12: false
        })

        // Обработчик кнопки "Редактировать"
        document.getElementById('editTaskButton').onclick = () => {
            closeModal();  // Закрываем окно просмотра перед открытием окна редактирования
            openModal(task);  // Открываем окно редактирования
        };

        // Обработчик кнопки "Удалить"
        document.getElementById('deleteTaskButton').onclick = () => {
            if (confirm('Вы уверены, что хотите удалить эту задачу?')) {
                deleteTask(task.ID);
                closeModal();
            }
        };

        // Логика для расчета времени до дедлайна и обновления прогресс-бара
        updateTimeRemaining(deadline, createdAt, now);

        taskViewModal.style.display = 'block';
    };

    // Функция для расчета и отображения оставшегося времени до дедлайна
    function updateTimeRemaining(deadline, createdAt, now) {
        const totalDuration = deadline - createdAt; // Общее время от создания до дедлайна в миллисекундах
        const timePassed = now - createdAt; // Прошедшее время от создания до текущего момента
        const timeRemaining = Math.max(0, deadline - now); // Оставшееся время до дедлайна в миллисекундах

        // Рассчитываем прогресс-бар
        const progressBar = document.getElementById('progressBar');
        let progressPercentage;

        if (now > deadline) {
            // Если задача просрочена
            progressPercentage = 100;
            progressBar.style.backgroundColor = 'red';
        } else {
            // Если задача не просрочена
            progressPercentage = Math.min(100, (timePassed / totalDuration) * 100);
            progressBar.style.backgroundColor = ''; // Убираем цвет, если ранее был установлен красный
        }

        progressBar.style.width = `${progressPercentage}%`;

        const daysRemaining = Math.floor(timeRemaining / (1000 * 60 * 60 * 24));
        const hoursRemaining = Math.floor((timeRemaining % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60));
        const minutesRemaining = Math.floor((timeRemaining % (1000 * 60 * 60)) / (1000 * 60));

        const timeParts = [];
        if (daysRemaining > 0) timeParts.push(`${daysRemaining}д`);
        if (hoursRemaining > 0) timeParts.push(`${hoursRemaining}ч`);
        if (minutesRemaining > 0) timeParts.push(`${minutesRemaining}м`);
        if (timeRemaining <= 0) timeParts.push('Задача просрочена');

        // Обновление текста
        document.getElementById('timeRemainingValue').textContent = timeParts.join(' ');

        // Обновление оставшегося времени каждую минуту
        setTimeout(() => updateTimeRemaining(deadline, createdAt, new Date()), 60000);
    }

    // Функция для закрытия модальных окон
    const closeModal = () => {
        taskModal.style.display = 'none';
        taskViewModal.style.display = 'none';
    };

    // Закрытие окон при нажатии на крестик
    closeModalButtons.forEach(button => {
        button.addEventListener('click', closeModal);
    });

    // Закрытие окон при клике вне модального окна
    window.addEventListener('click', (event) => {
        if (event.target === taskModal || event.target === taskViewModal) {
            closeModal();
        }
    });

    addTaskButton.addEventListener('click', () => openModal());

    // Отображаем выбор команды при активации чекбокса командного задания
    isTeamTaskCheckbox.addEventListener('change', () => {
        if (isTeamTaskCheckbox.checked) {
            teamSelection.style.display = 'block';
            loadTeams();
        } else {
            teamSelection.style.display = 'none';
            memberSelection.style.display = 'none';
        }
    });

    let teamsMap = {}; // Список участников команды

    async function loadTeams() {
        try {
            const response = await fetch(`http://localhost:8080/getTeams?user_id=${userId}`);
            if (response.ok) {
                const teams = await response.json();

                if (teams.length === 0) {
                    alert('У вас нет команд');
                    return;
                }

                teamSelect.innerHTML = '';

                // Добавляем команды в селектор
                teams.forEach(team => {
                    const option = document.createElement('option');
                    option.value = team.ID;
                    option.textContent = team.Name;
                    teamSelect.appendChild(option);
                });

                // Хранение участников каждой команды в объекте teamsMap
                teams.forEach(team => {
                    teamsMap[team.ID] = team.Members;
                });

                teamSelect.value = teams[0].ID;  // Устанавливаем значение в селекторе
                loadMembers(teamsMap[teams[0].ID]);  // Загружаем участников команды

                // Добавляем обработчик изменения команды
                teamSelect.addEventListener('change', () => {
                    const selectedTeamId = teamSelect.value;
                    if (selectedTeamId) {
                        loadMembers(teamsMap[selectedTeamId]);
                    } else {
                        memberSelection.style.display = 'none';
                    }
                });
            } else {
                alert('Ошибка загрузки команд');
            }
        } catch (error) {
            console.error('Ошибка загрузки команд:', error);
        }
    }

    // Функция для отображения участников выбранной команды
    function loadMembers(members) {
        memberSelect.innerHTML = '';  // Очищаем текущий список участников

        // Добавляем опцию "Назначить на всю команду"
        const allMembersOption = document.createElement('option');
        allMembersOption.value = 'all';
        allMembersOption.textContent = 'Назначить на всю команду';
        memberSelect.appendChild(allMembersOption);

        // Добавляем участников команды
        members.forEach(member => {
            const option = document.createElement('option');
            option.value = member;
            option.textContent = member;
            memberSelect.appendChild(option);
        });
        memberSelection.style.display = 'block';  // Отображаем блок выбора участников
    }


    const taskForm = document.getElementById('taskForm');
    taskForm.addEventListener('submit', async (event) => {
        event.preventDefault();

        const name = document.getElementById('taskName').value;
        const description = document.getElementById('taskDescription').value;
        const deadlineDate = new Date(document.getElementById('taskDeadline').value);
        const deadlineIso = deadlineDate.toISOString();
        const taskStatus = document.getElementById('taskStatus').value;

        try {
            const userId = document.cookie.replace(/(?:(?:^|.*;\s*)user_id\s*\=\s*([^;]*).*$)|^.*$/, "$1");

            let response;
            let accessUserIds = [];

            if (memberSelect.value === 'all') {
                const selectedTeamId = teamSelect.value;

                // Получаем ID участников команды по их юзернеймам
                const members = teamsMap[selectedTeamId]; // Массив юзернеймов
                const userIdsResponse = await fetch(`http://localhost:8080/getUserIds`, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({ usernames: members })
                });

                if (userIdsResponse.ok) {
                    const userIds = await userIdsResponse.json();
                    accessUserIds = userIds.map(user => user.ID); // Предполагаем, что ответ содержит id
                } else {
                    alert('Ошибка получения ID участников');
                    return;
                }
            } else {
                accessUserIds.push(memberSelect.value); // Добавляем только выбранного участника
            }

            // Отправка задачи на бэк
            if (editingTaskID) {
                loaderContainer.style.display = 'flex';
                response = await fetch(`http://localhost:8080/tasks/update`, {
                    method: 'PUT',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({
                        id: editingTaskID,
                        name,
                        description,
                        deadline: deadlineIso,
                        taskStatus
                    })
                });
            } else {
                loaderContainer.style.display = 'flex';
                response = await fetch('http://localhost:8080/tasks/create', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({
                        name,
                        description,
                        deadline: deadlineIso,
                        creator_id: userId,
                    })
                });
            }

            if (response.ok) {
                const result = await response.json();
                const taskId = result.task_id;

                // Вставка доступа к задаче
                await Promise.all(accessUserIds.map(async (accessUserId) => {
                    await fetch('http://localhost:8080/task_access', {
                        method: 'POST',
                        headers: {
                            'Content-Type': 'application/json'
                        },
                        body: JSON.stringify({
                            task_id: taskId,
                            user_id: accessUserId
                        })
                    });
                }));

                loaderContainer.style.display = 'none';
                fetchTasks(userId);
                closeModal();
            } else {
                alert('Ошибка сохранения задачи');
            }
        } catch (error) {
            console.error('Error:', error);
        }
        loaderContainer.style.display = 'none';
    });

    async function deleteTask(taskID) {
        try {
            loaderContainer.style.display = 'flex';
            const response = await fetch(`http://localhost:8080/tasks/delete?id=${taskID}`, {
                method: 'DELETE',
                headers: {
                    'Content-Type': 'application/json'
                }
            });

            if (response.ok) {
                const userId = document.cookie.replace(/(?:(?:^|.*;\s*)user_id\s*\=\s*([^;]*).*$)|^.*$/, "$1");
                fetchTasks(userId);
            } else {
                alert('Ошибка удаления задачи');
            }
        } catch (error) {
            console.error('Error:', error);
        }
    }

    async function fetchTasks(userId) {
        try {
            const response = await fetch(`http://localhost:8080/tasks?user_id=${userId}`);
            if (response.ok) {
                loaderContainer.style.display = 'none';
                const tasks = await response.json();
                taskList.innerHTML = '';
                tasks.forEach(task => {
                    const li = document.createElement('li');
                    li.className = 'task-item';
                    const deadline = new Date(task.Deadline);

                    li.innerHTML = `
                    <div class="task-header">
                        <h3>${task.Name}</h3>
                        <span class="status-badge ${task.Status.toLowerCase()}">${getStatusLabel(task.Status)}</span>
                    </div>
                    <p class="task-deadline">Дедлайн: ${deadline.toLocaleString(undefined, {
                        year: 'numeric',
                        month: 'long',
                        day: 'numeric',
                        hour: '2-digit',
                        minute: '2-digit',
                        hour12: false
                    })}</p> 
                    <p class="task-description">${task.Description}</p>
                `;
                    li.onclick = () => openTaskViewModal(task);  // Устанавливаем обработчик клика для каждой задачи
                    taskList.appendChild(li);
                    loaderContainer.style.display = 'none';
                });
            } else {
                loaderContainer.style.display = 'none';
                alert('Ошибка получения задач');
            }
        } catch (error) {
            console.error('Error:', error);
        }
        loaderContainer.style.display = 'none';
    }

    const userId = document.cookie.replace(/(?:(?:^|.*;\s*)user_id\s*\=\s*([^;]*).*$)|^.*$/, "$1");
    if (userId) {
        loaderContainer.style.display = 'flex';
        fetchTasks(userId);
    }
});

function getStatusLabel(status) {
    switch (status.toLowerCase()) {
        case 'new':
            return 'Новая';
        case 'in_progress':
            return 'В работе';
        case 'completed':
            return 'Выполнена';
        case 'postponed':
            return 'Отложена';
        default:
            return 'Неизвестно';
    }
}
