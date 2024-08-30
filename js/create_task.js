document.addEventListener('DOMContentLoaded', () => {
    const addTaskButton = document.getElementById('addTaskButton');
    const taskModal = document.getElementById('taskModal');
    const closeModalButton = document.querySelector('.modal .close');
    const taskList = document.getElementById('taskList');
    let editingTaskID = null;

    // Функция для открытия модального окна
    const openModal = (task = null) => {
        if (task) {
            // Заполняем форму данными задачи для редактирования
            document.getElementById('taskName').value = task.name;
            document.getElementById('taskDescription').value = task.description;
            document.getElementById('taskDeadline').value = new Date(task.deadline).toISOString().split('T')[0];
            document.getElementById('taskUsers').value = ''; // Присвойте email пользователей сюда, если есть
            editingTaskID = task.id;
            document.getElementById('modalTitle').innerText = 'Редактировать задачу';
        } else {
            // Очищаем форму для новой задачи
            document.getElementById('taskName').value = '';
            document.getElementById('taskDescription').value = '';
            document.getElementById('taskDeadline').value = '';
            document.getElementById('taskUsers').value = '';
            editingTaskID = null;
            document.getElementById('modalTitle').innerText = 'Создать задачу';
        }
        taskModal.style.display = 'block';
    };

    // Функция для закрытия модального окна
    const closeModal = () => {
        taskModal.style.display = 'none';
    };

    // Открытие модального окна при нажатии на кнопку
    addTaskButton.addEventListener('click', () => openModal());

    // Закрытие модального окна при нажатии на крестик
    closeModalButton.addEventListener('click', closeModal);

    // Закрытие модального окна при клике за его пределами
    window.addEventListener('click', (event) => {
        if (event.target === taskModal) {
            closeModal();
        }
    });

    // Отправка данных задачи на сервер и обновление списка задач
    const taskForm = document.getElementById('taskForm');
    taskForm.addEventListener('submit', async (event) => {
        event.preventDefault();

        const name = document.getElementById('taskName').value;
        const description = document.getElementById('taskDescription').value;
        const deadlineDate = new Date(document.getElementById('taskDeadline').value);
        const deadlineIso = deadlineDate.toISOString(); // Преобразуем дату в формат ISO
        const userEmails = document.getElementById('taskUsers').value.split(',').map(email => email.trim()).filter(email => email);

        try {
            const userId = document.cookie.replace(/(?:(?:^|.*;\s*)user_id\s*\=\s*([^;]*).*$)|^.*$/, "$1");

            let response;
            if (editingTaskID) {
                // Редактирование существующей задачи
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
                        user_ids: userEmails // Преобразуйте email в user_id на сервере
                    })
                });
            } else {
                // Создание новой задачи
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
                        user_ids: userEmails // Преобразуйте email в user_id на сервере
                    })
                });
            }

            if (response.ok) {
                const result = await response.json();
                fetchTasks(userId);
                closeModal();
            } else {
                alert('Ошибка сохранения задачи');
            }
        } catch (error) {
            console.error('Error:', error);
        }
    });

    async function deleteTask(taskID) {
        try {
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



    // Функция для получения задач с сервера
    async function fetchTasks(userId) {
        try {
            const response = await fetch(`http://localhost:8081/tasks?user_id=${userId}`);
            if (response.ok) {
                const tasks = await response.json();
                taskList.innerHTML = '';
                tasks.forEach(task => {
                    const li = document.createElement('li');
                    li.className = 'task-item';
                    li.innerHTML = `
                        <h3>${task.name}</h3>
                        <p>${task.description}</p>
                        <span>Дедлайн: ${new Date(task.deadline).toLocaleDateString()}</span>
                        <div class="task-actions">
                            <button class="edit-btn" data-id="${task.id}">Редактировать</button>
                            <button class="delete-btn" data-id="${task.id}">Удалить</button>
                        </div>
                    `;
                    taskList.appendChild(li);
                });

                document.querySelectorAll('.edit-btn').forEach(button => {
                    button.addEventListener('click', (event) => {
                        const taskID = event.target.dataset.id;
                        fetch(`http://localhost:8081/tasks/${taskID}`, {
                            method: 'GET',
                            headers: {
                                'Content-Type': 'application/json'
                            }
                        })
                            .then(response => response.json())
                            .then(task => openModal(task))
                            .catch(error => {
                                console.error('Error fetching task:', error);
                                alert('Ошибка при загрузке задачи для редактирования');
                            });
                    });
                });


                document.querySelectorAll('.delete-btn').forEach(button => {
                    button.addEventListener('click', (event) => {
                        const taskID = event.target.dataset.id;
                        if (confirm('Вы уверены, что хотите удалить эту задачу?')) {
                            deleteTask(taskID);
                        }
                    });
                });
            } else {
                alert('Ошибка получения задач');
            }
        } catch (error) {
            console.error('Error:', error);
        }
    }


    // Получаем задачи при загрузке страницы
    const userId = document.cookie.replace(/(?:(?:^|.*;\s*)user_id\s*\=\s*([^;]*).*$)|^.*$/, "$1");
    if (userId) {
        fetchTasks(userId);
    }
});
