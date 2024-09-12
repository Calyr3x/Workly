document.addEventListener('DOMContentLoaded', () => {
    const addTaskButton = document.getElementById('addTaskButton');
    const taskModal = document.getElementById('taskModal');
    const taskViewModal = document.getElementById('taskViewModal');
    const closeModalButtons = document.querySelectorAll('.modal .close');
    const taskList = document.getElementById('taskList');
    const loaderContainer = document.querySelector('.loader-container');
    let editingTaskID = null;

    // Функция для открытия модального окна создания/редактирования задачи
    const openModal = (task = null) => {
        if (task) {
            document.getElementById('taskName').value = task.name;
            document.getElementById('taskDescription').value = task.description;
            document.getElementById('taskDeadline').value = new Date(task.deadline).toISOString().split('T')[0];
            document.getElementById('taskUsers').value = '';
            editingTaskID = task.id;
            document.getElementById('modalTitle').innerText = 'Редактировать задачу';
        } else {
            document.getElementById('taskName').value = '';
            document.getElementById('taskDescription').value = '';
            document.getElementById('taskDeadline').value = '';
            document.getElementById('taskUsers').value = '';
            editingTaskID = null;
            document.getElementById('modalTitle').innerText = 'Создать задачу';
        }
        taskModal.style.display = 'block';
    };

    // Функция для открытия окна просмотра задачи
    const openTaskViewModal = (task) => {
        const deadline = new Date(task.deadline);
        const createdAt = new Date(task.created_at);
        const now = new Date();

        document.getElementById('viewTaskTitle').textContent = task.name;
        document.getElementById('viewTaskDescription').textContent = task.description;
        document.getElementById('viewTaskDeadline').querySelector('span').textContent = deadline.toLocaleDateString();

        // Обработчик кнопки "Редактировать"
        document.getElementById('editTaskButton').onclick = () => {
            closeModal();  // Закрываем окно просмотра перед открытием окна редактирования
            openModal(task);  // Открываем окно редактирования
        };

        // Обработчик кнопки "Удалить"
        document.getElementById('deleteTaskButton').onclick = () => {
            if (confirm('Вы уверены, что хотите удалить эту задачу?')) {
                deleteTask(task.id);
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

        const daysRemaining = Math.floor(timeRemaining / (1000 * 60 * 60 * 24));
        const hoursRemaining = Math.floor((timeRemaining % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60));
        const minutesRemaining = Math.floor((timeRemaining % (1000 * 60 * 60)) / (1000 * 60));

        // Создаем массив частей времени, исключая нулевые значения
        const timeParts = [];
        if (daysRemaining > 0) timeParts.push(`${daysRemaining}д`);
        if (hoursRemaining > 0) timeParts.push(`${hoursRemaining}ч`);
        if (minutesRemaining > 0) timeParts.push(`${minutesRemaining}м`);

        // Соединяем части в одну строку
        const timeText = timeParts.join(' ');

        // Обновление текста
        document.getElementById('timeRemainingValue').textContent = timeText;

        // Обновление прогресс-бара
        const progressBar = document.getElementById('progressBar');
        const progressPercentage = Math.max(0, Math.min(100, (timePassed / totalDuration) * 100));
        progressBar.style.width = `${progressPercentage}%`;

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

    const taskForm = document.getElementById('taskForm');
    taskForm.addEventListener('submit', async (event) => {
        event.preventDefault();

        const name = document.getElementById('taskName').value;
        const description = document.getElementById('taskDescription').value;
        const deadlineDate = new Date(document.getElementById('taskDeadline').value);
        const deadlineIso = deadlineDate.toISOString();
        const userEmails = document.getElementById('taskUsers').value.split(',').map(email => email.trim()).filter(email => email);

        try {
            const userId = document.cookie.replace(/(?:(?:^|.*;\s*)user_id\s*\=\s*([^;]*).*$)|^.*$/, "$1");

            let response;
            if (editingTaskID) {
                loaderContainer.style.display = 'flex';
                response = await fetch(`http://localhost:8081/tasks/update`, {
                    method: 'PUT',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({
                        id: editingTaskID,
                        name,
                        description,
                        deadline: deadlineIso,
                        user_ids: userEmails
                    })
                });
            } else {
                loaderContainer.style.display = 'flex';
                response = await fetch('http://localhost:8081/tasks/create', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({
                        name,
                        description,
                        deadline: deadlineIso,
                        creator_id: userId,
                        user_ids: userEmails
                    })
                });
            }

            if (response.ok) {
                loaderContainer.style.display = 'none';
                const result = await response.json();
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
            const response = await fetch(`http://localhost:8081/tasks/delete?id=${taskID}`, {
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
            const response = await fetch(`http://localhost:8081/tasks?user_id=${userId}`);
            if (response.ok) {
                loaderContainer.style.display = 'none';
                const tasks = await response.json();
                taskList.innerHTML = '';
                tasks.forEach(task => {
                    const li = document.createElement('li');
                    li.className = 'task-item';
                    li.innerHTML = `
                        <h3>${task.name}</h3>
                        <p>${task.description}</p>
                        <span>Дедлайн: ${new Date(task.deadline).toLocaleString()}</span>
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
