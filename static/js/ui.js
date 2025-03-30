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

        employerTaskList: document.getElementById('employerTaskList'),
        employeeTaskList: document.getElementById('employeeTaskList'),
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
            ui.loadEmployeeData();
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
                ${user.email} (${user.role})
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
                ui.updateAssigneeDropdown(employeeSummary.data.employees);
                ui.updateAssigneeFilterDropdown(employeeSummary.data.employees);
            }
        }
    },

    // Update the assignee filter dropdown with employees
    updateAssigneeFilterDropdown: (employees) => {
        ui.elements.taskAssigneeFilter.innerHTML = '<option value="">All Assignees</option>';

        employees.forEach(employeeData => {
            const employee = employeeData.employee;
            if (employee.role === 'employee') {
                const option = document.createElement('option');
                option.value = employee.id;
                option.textContent = `${employee.name}`;
                ui.elements.taskAssigneeFilter.appendChild(option);
            }
        });
    },

    // Load data for employee dashboard
    loadEmployeeData: async () => {
        const taskResult = await tasks.getAssignedTasks();
        if (taskResult && taskResult.status === 'success') {
            ui.updateEmployeeTaskList(taskResult.data.tasks);
        }
    },

    // Update the assignee dropdown with employees
    updateAssigneeDropdown: (employees) => {
        ui.elements.taskAssignee.innerHTML = '<option value="">Select Employee</option>';

        employees.forEach(employeeData => {
            const employee = employeeData.employee;
            if (employee.role === 'employee') {
                const option = document.createElement('option');
                option.value = employee.id;
                option.textContent = `${employee.name} (${employee.email})`;
                ui.elements.taskAssignee.appendChild(option);
            }
        });
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

            row.innerHTML = `
                <td>${task.id}</td>
                <td>${task.title}</td>
                <td>${task.description || 'No description'}</td>
                <td>${tasks.formatStatus(task.status)}</td>
                <td>${dueDate}</td>
                <td>${assigneeName}</td>
            `;

            tableBody.appendChild(row);
        });
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

            taskCard.innerHTML = `
                <h4>${task.title}</h4>
                <p>${task.description || 'No description'}</p>
                <p><strong>Status:</strong> ${tasks.formatStatus(task.status)}</p>
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

    // Update the task list for employees
    updateEmployeeTaskList: (tasks) => {
        ui.elements.employeeTaskList.innerHTML = '';

        if (!tasks || tasks.length === 0) {
            ui.elements.employeeTaskList.innerHTML = '<p>No tasks assigned to you.</p>';
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

            taskCard.innerHTML = `
                <h4>${task.title}</h4>
                <p>${task.description || 'No description'}</p>
                <p><strong>Status:</strong> ${tasks.formatStatus(task.status)}</p>
                <p><strong>Due Date:</strong> ${dueDate}</p>
                <div class="task-actions">
                    <button class="task-status-btn" data-task-id="${task.id}" data-status="pending">Set Pending</button>
                    <button class="task-status-btn" data-task-id="${task.id}" data-status="in_progress">Set In Progress</button>
                    <button class="task-status-btn" data-task-id="${task.id}" data-status="completed">Set Completed</button>
                </div>
            `;

            ui.elements.employeeTaskList.appendChild(taskCard);
        });

        // Add event listeners to status buttons
        document.querySelectorAll('.task-status-btn').forEach(button => {
            button.addEventListener('click', async (e) => {
                const taskId = e.target.getAttribute('data-task-id');
                const status = e.target.getAttribute('data-status');

                const result = await tasks.updateTaskStatus(taskId, status);
                if (result && result.status === 'success') {
                    ui.loadEmployeeData(); // Reload the tasks
                }
            });
        });
    }
};
