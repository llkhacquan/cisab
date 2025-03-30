// ui.js - Handle UI updates and interactions

const ui = {
    // UI elements
    elements: {
        userSidebar: document.getElementById('userSidebar'),
        userList: document.getElementById('userList'),
        createUserBtn: document.getElementById('createUserBtn'),

        authPanel: document.getElementById('authPanel'),
        loginForm: document.getElementById('loginForm'),
        loginEmail: document.getElementById('loginEmail'),
        loginPassword: document.getElementById('loginPassword'),
        loginBtn: document.getElementById('loginBtn'),

        userInfo: document.getElementById('userInfo'),
        currentUserDetails: document.getElementById('currentUserDetails'),
        logoutBtn: document.getElementById('logoutBtn'),

        taskPanel: document.getElementById('taskPanel'),
        employerPanel: document.getElementById('employerPanel'),
        employeePanel: document.getElementById('employeePanel'),

        // Modal elements
        taskModal: document.getElementById('taskModal'),
        closeModalBtn: document.querySelector('.close-modal'),
        showAddTaskBtn: document.getElementById('showAddTaskBtn'),

        // Task form elements
        taskTitle: document.getElementById('taskTitle'),
        taskDescription: document.getElementById('taskDescription'),
        taskDueDate: document.getElementById('taskDueDate'),
        taskAssignee: document.getElementById('taskAssignee'),
        createTaskBtn: document.getElementById('createTaskBtn'),
        createRandomTaskBtn: document.getElementById('createRandomTaskBtn'),

        // Task table and filters
        taskTable: document.getElementById('taskTable'),
        taskTableBody: document.getElementById('taskTableBody'),
        taskStatusFilter: document.getElementById('taskStatusFilter'),
        taskAssigneeFilter: document.getElementById('taskAssigneeFilter'),
        taskSortBy: document.getElementById('taskSortBy'),
        taskSortOrder: document.getElementById('taskSortOrder'),

        // Task assignment modal elements
        assignTaskModal: document.getElementById('assignTaskModal'),
        closeAssignModalBtn: document.querySelector('.close-modal-assign'),
        assignTaskTitle: document.getElementById('assignTaskTitle'),
        assignTaskId: document.getElementById('assignTaskId'),
        assignTaskAssignee: document.getElementById('assignTaskAssignee'),
        confirmAssignTaskBtn: document.getElementById('confirmAssignTaskBtn'),

        employerTaskList: document.getElementById('employerTaskList'),
        employeeTaskList: document.getElementById('employeeTaskList'),
        employeeTaskTable: document.getElementById('employeeTaskTable'),
        employeeTaskTableBody: document.getElementById('employeeTaskTableBody'),
        employeeSummary: document.getElementById('employeeSummary')
    },

    // Initialize UI based on auth state
    initUI: () => {
        ui.updateUserList();

        // Check if user is authenticated
        if (auth.isAuthenticated()) {
            const user = auth.getCurrentUser();
            ui.showAuthenticatedUI(user);
        } else {
            ui.showLoginForm();
        }
    },

    // Show login form
    showLoginForm: () => {
        ui.elements.userInfo.classList.add('hidden');
        ui.elements.loginForm.classList.remove('hidden');
        ui.elements.taskPanel.classList.add('hidden');
    },

    // Show authenticated UI based on user role
    showAuthenticatedUI: (user) => {
        // Update user info display with a compact format
        ui.elements.currentUserDetails.innerHTML = `${user.name} (${user.role}, ID:${user.id})`;

        ui.elements.loginForm.classList.add('hidden');
        ui.elements.userInfo.classList.remove('hidden');
        ui.elements.taskPanel.classList.remove('hidden');

        // Show appropriate panel based on role
        if (user.role === 'employer') {
            ui.elements.employerPanel.classList.remove('hidden');
            ui.elements.employeePanel.classList.add('hidden');
            ui.loadEmployerData();
        } else if (user.role === 'employee') {
            ui.elements.employerPanel.classList.add('hidden');
            ui.elements.employeePanel.classList.remove('hidden');
            ui.loadEmployeeAssignedTasks();
        }
    },

    // Update the user list in the sidebar
    updateUserList: () => {
        const users = auth.getCreatedUsers();
        ui.elements.userList.innerHTML = '';

        users.forEach(user => {
            const userEntry = document.createElement('div');
            userEntry.className = 'user-entry';
            userEntry.innerHTML = `
                <strong>${user.name}</strong><br>
                ${user.email} (${user.role})<br>
                <span class="user-id">ID: ${user.id || 'N/A'}</span>
            `;

            // Add click event to log in as this user
            userEntry.addEventListener('click', () => {
                ui.elements.loginEmail.value = user.email;
                ui.elements.loginPassword.value = user.password;
            });

            ui.elements.userList.appendChild(userEntry);
        });
    },

    // Generate a random user for demo purposes
    generateRandomUser: () => {
        const roles = ['employee', 'employer'];
        const names = ['Alice', 'Bob', 'Charlie', 'Diana', 'Edward', 'Fiona', 'George', 'Hannah'];
        const surnames = ['Smith', 'Johnson', 'Williams', 'Brown', 'Jones', 'Miller', 'Davis', 'Garcia', 'Wilson'];

        const name = `${names[Math.floor(Math.random() * names.length)]} ${surnames[Math.floor(Math.random() * surnames.length)]}`;
        const email = `${name.toLowerCase().replace(' ', '.')}${Math.floor(Math.random() * 1000)}@example.com`;
        const password = `password${Math.floor(Math.random() * 1000)}`;
        const role = roles[Math.floor(Math.random() * roles.length)];

        return {
            name,
            email,
            password,
            role
        };
    },

    // Load data for employer dashboard
    loadEmployerData: async () => {
        // Initialize the page with both the tasks table and employee data
        await ui.loadEmployeeData();
        await ui.loadTaskTable();
    },

    // Load the tasks table using GET /api/v1/tasks
    loadTaskTable: async () => {
        // Show loading indicator
        const tableBody = ui.elements.taskTableBody;
        tableBody.innerHTML = '<tr><td colspan="6" style="text-align: center;">Loading tasks...</td></tr>';

        // Get filter parameters from the UI controls
        const params = {
            status: ui.elements.taskStatusFilter.value,
            sort_by: ui.elements.taskSortBy.value,
            sort_order: ui.elements.taskSortOrder.value,
            limit: 50  // Default limit, could be made configurable
        };

        // Add assignee filter if selected
        if (ui.elements.taskAssigneeFilter.value) {
            params.assignee_id = parseInt(ui.elements.taskAssigneeFilter.value, 10);
        }

        // Remove empty values
        Object.keys(params).forEach(key => {
            if (params[key] === null || params[key] === '') {
                delete params[key];
            }
        });

        console.log('Fetching tasks with params:', params);

        // Load tasks created by this employer with filters
        const taskResult = await tasks.getTasks(params);
        if (taskResult && taskResult.status === 'success') {
            console.log('Received tasks:', taskResult.data.tasks);
            ui.updateTaskTable(taskResult.data.tasks);
        } else {
            console.error('Failed to load tasks:', taskResult);
            tableBody.innerHTML = '<tr><td colspan="6" style="text-align: center;">Error loading tasks</td></tr>';
        }
    },

    // Load employee summary data using GET /api/v1/employee-summary
    loadEmployeeData: async () => {
        // Load employees for the assignee dropdown and employee summary
        if (!ui.elements.taskAssignee.options.length || !ui.elements.taskAssigneeFilter.options.length) {
            // Only load employee data if dropdowns aren't populated yet
            const employeeSummary = await tasks.getEmployeeSummary();
            if (employeeSummary && employeeSummary.status === 'success') {
                ui.updateEmployeeSummary(employeeSummary.data.employees);

                // Use the getAllUsers API to get a complete list of users
                const allUsers = await tasks.getAllUsers();
                if (allUsers && allUsers.status === 'success' && allUsers.data.users) {
                    // Use the complete user list for the dropdowns
                    ui.updateAssigneeDropdown(allUsers.data.users);
                    ui.updateAssigneeFilterDropdown(allUsers.data.users);
                } else {
                    // Fallback to employee summary if the getAllUsers API fails
                    ui.updateAssigneeDropdown(employeeSummary.data.employees);
                    ui.updateAssigneeFilterDropdown(employeeSummary.data.employees);
                }
            }
        }
    },

    // Update the assignee filter dropdown with employees
    updateAssigneeFilterDropdown: (employees) => {
        ui.elements.taskAssigneeFilter.innerHTML = '<option value="">All Assignees</option>';

        // Check if we're using employee summary format or users/all format
        if (employees[0] && employees[0].employee) {
            // Using employee summary format
            employees.forEach(employeeData => {
                const employee = employeeData.employee;
                if (employee.role === 'employee') {
                    const option = document.createElement('option');
                    option.value = employee.id;
                    option.textContent = `${employee.name}`;
                    ui.elements.taskAssigneeFilter.appendChild(option);
                }
            });
        } else {
            // Using users/all format
            employees.forEach(user => {
                if (user.role === 'employee') {
                    const option = document.createElement('option');
                    option.value = user.id;
                    option.textContent = `${user.name}`;
                    ui.elements.taskAssigneeFilter.appendChild(option);
                }
            });
        }
    },

    // Load tasks assigned to the current employee (for employee dashboard)
    loadEmployeeAssignedTasks: async () => {
        console.log('Loading assigned tasks for employee');
        const taskResult = await tasks.getAssignedTasks();
        console.log('Assigned tasks result:', taskResult);
        if (taskResult && taskResult.status === 'success') {
            console.log('Tasks assigned to employee:', taskResult.data.tasks);
            ui.updateEmployeeTaskList(taskResult.data.tasks);
        } else {
            console.error('Failed to load assigned tasks:', taskResult);
            ui.elements.employeeTaskList.innerHTML = '<p>Error loading tasks. Please try again.</p>';
        }
    },

    // Update the assignee dropdown with employees
    updateAssigneeDropdown: (employees) => {
        ui.elements.taskAssignee.innerHTML = '<option value="">Select Employee</option>';
        // Also update the assignment modal dropdown
        ui.elements.assignTaskAssignee.innerHTML = '<option value="">Select Employee</option>';

        // Check if we're using employee summary format or users/all format
        if (employees[0] && employees[0].employee) {
            // Using employee summary format
            employees.forEach(employeeData => {
                const employee = employeeData.employee;
                if (employee.role === 'employee') {
                    // Add to task creation dropdown
                    const option = document.createElement('option');
                    option.value = employee.id;
                    option.textContent = `${employee.name} (${employee.email})`;
                    ui.elements.taskAssignee.appendChild(option);

                    // Add to task assignment dropdown
                    const assignOption = document.createElement('option');
                    assignOption.value = employee.id;
                    assignOption.textContent = `${employee.name} (${employee.email})`;
                    ui.elements.assignTaskAssignee.appendChild(assignOption);
                }
            });
        } else {
            // Using users/all format
            employees.forEach(user => {
                if (user.role === 'employee') {
                    // Add to task creation dropdown
                    const option = document.createElement('option');
                    option.value = user.id;
                    option.textContent = `${user.name} (${user.email})`;
                    ui.elements.taskAssignee.appendChild(option);

                    // Add to task assignment dropdown
                    const assignOption = document.createElement('option');
                    assignOption.value = user.id;
                    assignOption.textContent = `${user.name} (${user.email})`;
                    ui.elements.assignTaskAssignee.appendChild(assignOption);
                }
            });
        }
    },

    // Update employee summary for employer dashboard
    updateEmployeeSummary: (employees) => {
        ui.elements.employeeSummary.innerHTML = '';

        employees.forEach(employeeData => {
            const employee = employeeData.employee;
            const stats = employeeData.statistics;

            const employeeCard = document.createElement('div');
            employeeCard.className = 'card';
            employeeCard.innerHTML = `
                <h4>${employee.name}</h4>
                <p>Email: ${employee.email}</p>
                <p>Total Tasks: ${stats.total_tasks}</p>
                <p>Pending: ${stats.pending}</p>
                <p>In Progress: ${stats.in_progress}</p>
                <p>Completed: ${stats.completed}</p>
            `;

            ui.elements.employeeSummary.appendChild(employeeCard);
        });
    },

    // Update the task table with tasks data
    updateTaskTable: (tasksData) => {
        const tableBody = ui.elements.taskTableBody;
        tableBody.innerHTML = '';

        if (!tasksData || tasksData.length === 0) {
            const emptyRow = document.createElement('tr');
            emptyRow.innerHTML = '<td colspan="6" style="text-align: center;">No tasks found</td>';
            tableBody.appendChild(emptyRow);
            return;
        }

        // Get employee map for assignee names
        const employeeMap = {};
        document.querySelectorAll('#taskAssigneeFilter option').forEach(option => {
            if (option.value) {
                employeeMap[option.value] = option.textContent;
            }
        });

        tasksData.forEach(task => {
            const row = document.createElement('tr');

            // Format the due date if available
            let dueDate = 'Not set';
            if (task.due_date) {
                dueDate = new Date(task.due_date).toLocaleString();
            }

            // Get assignee name if available
            let assigneeName = 'Unassigned';
            if (task.assignee_id && employeeMap[task.assignee_id]) {
                assigneeName = employeeMap[task.assignee_id];
            } else if (task.assignee_id) {
                assigneeName = `Employee ID: ${task.assignee_id}`;
            }

            // Format status
            let formattedStatus = task.status;
            switch (task.status) {
                case 'pending': formattedStatus = 'Pending'; break;
                case 'in_progress': formattedStatus = 'In Progress'; break;
                case 'completed': formattedStatus = 'Completed'; break;
            }

            row.innerHTML = `
                <td>${task.id}</td>
                <td>${task.title}</td>
                <td>${task.description || 'No description'}</td>
                <td>${formattedStatus}</td>
                <td>${dueDate}</td>
                <td>${assigneeName}</td>
                <td>
                    <button class="assign-task-btn" data-task-id="${task.id}" data-task-title="${task.title}">Assign</button>
                </td>
            `;

            tableBody.appendChild(row);
        });

        // Add event listeners to assign buttons
        document.querySelectorAll('.assign-task-btn').forEach(button => {
            button.addEventListener('click', async (e) => {
                const taskId = e.target.getAttribute('data-task-id');
                const taskTitle = e.target.getAttribute('data-task-title');
                await ui.showAssignTaskModal(taskId, taskTitle);
            });
        });
    },

    // Show the assign task modal
    showAssignTaskModal: async (taskId, taskTitle) => {
        // Set the task ID and title in the modal
        ui.elements.assignTaskId.value = taskId;
        ui.elements.assignTaskTitle.textContent = taskTitle;

        // Show the modal
        ui.elements.assignTaskModal.classList.remove('hidden');
        ui.elements.assignTaskModal.classList.add('active');

        // If the dropdown is empty, load users directly from API
        if (ui.elements.assignTaskAssignee.options.length <= 1) {
            const usersResult = await tasks.getAllUsers();
            if (usersResult && usersResult.status === 'success') {
                ui.updateAssigneeDropdown(usersResult.data.users);
            }
        }
    },

    // Close the assign task modal
    closeAssignTaskModal: () => {
        ui.elements.assignTaskModal.classList.remove('active');
        ui.elements.assignTaskModal.classList.add('hidden');
    },

    // Update the task list for employers (legacy method, replaced by table)
    updateEmployerTaskList: (tasks) => {
        if (!ui.elements.employerTaskList) return; // Element might be removed

        ui.elements.employerTaskList.innerHTML = '';

        if (!tasks || tasks.length === 0) {
            ui.elements.employerTaskList.innerHTML = '<p>No tasks created yet.</p>';
            return;
        }

        tasks.forEach(task => {
            const taskCard = document.createElement('div');
            taskCard.className = 'card';

            // Format the due date if available
            let dueDate = 'Not set';
            if (task.due_date) {
                dueDate = new Date(task.due_date).toLocaleString();
            }

            // Format status directly
            let formattedStatus = task.status;
            if (task.status === 'pending') formattedStatus = 'Pending';
            else if (task.status === 'in_progress') formattedStatus = 'In Progress';
            else if (task.status === 'completed') formattedStatus = 'Completed';
            
            taskCard.innerHTML = `
                <h4>${task.title}</h4>
                <p>${task.description || 'No description'}</p>
                <p><strong>Status:</strong> ${formattedStatus}</p>
                <p><strong>Due Date:</strong> ${dueDate}</p>
                <p><strong>Assignee ID:</strong> ${task.assignee_id || 'Unassigned'}</p>
                <div class="task-actions">
                    <button class="task-status-btn" data-task-id="${task.id}" data-status="pending">Set Pending</button>
                    <button class="task-status-btn" data-task-id="${task.id}" data-status="in_progress">Set In Progress</button>
                    <button class="task-status-btn" data-task-id="${task.id}" data-status="completed">Set Completed</button>
                </div>
            `;

            ui.elements.employerTaskList.appendChild(taskCard);
        });

        // Add event listeners to status buttons
        document.querySelectorAll('.task-status-btn').forEach(button => {
            button.addEventListener('click', async (e) => {
                const taskId = e.target.getAttribute('data-task-id');
                const status = e.target.getAttribute('data-status');

                const result = await tasks.updateTaskStatus(taskId, status);
                if (result && result.status === 'success') {
                    ui.loadEmployerData(); // Reload the tasks
                }
            });
        });
    },

    // Update the task list for employees using a table layout
    updateEmployeeTaskList: (tasksList) => {
        console.log('Updating employee task list with tasks:', tasksList);
        const tableBody = ui.elements.employeeTaskTableBody;
        tableBody.innerHTML = '';

        if (!tasksList || !Array.isArray(tasksList) || tasksList.length === 0) {
            console.log('No tasks assigned to employee');
            tableBody.innerHTML = '<tr><td colspan="6" style="text-align: center;">No tasks assigned to you.</td></tr>';
            return;
        }

        tasksList.forEach(task => {
            const row = document.createElement('tr');

            // Format the due date if available
            let dueDate = 'Not set';
            if (task.due_date) {
                dueDate = new Date(task.due_date).toLocaleString();
            }

            // Format the status properly without relying on external function
            let formattedStatus = task.status;
            if (task.status === 'pending') formattedStatus = 'Pending';
            else if (task.status === 'in_progress') formattedStatus = 'In Progress';
            else if (task.status === 'completed') formattedStatus = 'Completed';

            row.innerHTML = `
                <td>${task.id}</td>
                <td>${task.title}</td>
                <td>${task.description || 'No description'}</td>
                <td>${formattedStatus}</td>
                <td>${dueDate}</td>
                <td>
                    <div class="task-actions">
                        <button class="task-status-btn" data-task-id="${task.id}" data-status="pending">Set Pending</button>
                        <button class="task-status-btn" data-task-id="${task.id}" data-status="in_progress">Set In Progress</button>
                        <button class="task-status-btn" data-task-id="${task.id}" data-status="completed">Set Completed</button>
                    </div>
                </td>
            `;

            tableBody.appendChild(row);
        });

        // Add event listeners to status buttons
        document.querySelectorAll('.task-status-btn').forEach(button => {
            button.addEventListener('click', async (e) => {
                const taskId = e.target.getAttribute('data-task-id');
                const status = e.target.getAttribute('data-status');

                const result = await tasks.updateTaskStatus(taskId, status);
                if (result && result.status === 'success') {
                    ui.loadEmployeeAssignedTasks(); // Reload the assigned tasks
                }
            });
        });
    }
};
