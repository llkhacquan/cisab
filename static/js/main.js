// main.js - Main application logic

// Initialize the application when the DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    // Initialize UI based on authentication state
    ui.initUI();

    // Set up event listeners
    setupEventListeners();
});

// Set up all event listeners
function setupEventListeners() {
    // Create random user button
    ui.elements.createUserBtn.addEventListener('click', async () => {
        const randomUser = ui.generateRandomUser();

        // Register the user via API
        const result = await api.registerUser(randomUser);

        if (result && result.status === 'success') {
            // Save the created user with password for demo purposes
            auth.saveCreatedUser(result.data.user, randomUser.password);
            ui.updateUserList();

            // Auto-fill the login form with the new user's credentials
            ui.elements.loginEmail.value = randomUser.email;
            ui.elements.loginPassword.value = randomUser.password;
        } else {
            alert('Failed to create user: ' + (result?.error?.message || 'Unknown error'));
        }
    });

    // Login button
    ui.elements.loginBtn.addEventListener('click', async () => {
        const email = ui.elements.loginEmail.value;
        const password = ui.elements.loginPassword.value;

        if (!email || !password) {
            alert('Please enter email and password');
            return;
        }

        const result = await api.loginUser({ email, password });

        if (result && result.status === 'success') {
            // Save authentication data
            auth.saveAuth(result.data.user, result.data.token, result.data.token_expiry);

            // Update UI for authenticated user
            ui.showAuthenticatedUI(result.data.user);
        } else {
            alert('Login failed: ' + (result?.error?.message || 'Invalid credentials'));
        }
    });

    // Logout button
    ui.elements.logoutBtn.addEventListener('click', () => {
        auth.clearAuth();
        ui.showLoginForm();
    });

    // Show task modal
    ui.elements.showAddTaskBtn.addEventListener('click', () => {
        ui.elements.taskModal.classList.remove('hidden');
        ui.elements.taskModal.classList.add('active');
    });

    // Close task modal
    ui.elements.closeModalBtn.addEventListener('click', () => {
        ui.elements.taskModal.classList.remove('active');
        ui.elements.taskModal.classList.add('hidden');
    });

    // Close modal when clicking outside content
    window.addEventListener('click', (event) => {
        if (event.target === ui.elements.taskModal) {
            ui.elements.taskModal.classList.remove('active');
            ui.elements.taskModal.classList.add('hidden');
        }
    });

    // Create task button (for employers)
    ui.elements.createTaskBtn.addEventListener('click', async () => {
        const title = ui.elements.taskTitle.value;
        const description = ui.elements.taskDescription.value;
        const dueDate = ui.elements.taskDueDate.value;
        const assigneeId = ui.elements.taskAssignee.value;

        if (!title) {
            alert('Please enter a task title');
            return;
        }

        const taskData = {
            title,
            description,
            status: 'pending'
        };

        // Add due date if provided
        if (dueDate) {
            taskData.due_date = new Date(dueDate).toISOString();
        }

        // Add assignee if selected
        if (assigneeId) {
            taskData.assignee_id = parseInt(assigneeId, 10);
        }

        const result = await tasks.createTask(taskData);

        if (result && result.status === 'success') {
            // Clear form fields
            ui.elements.taskTitle.value = '';
            ui.elements.taskDescription.value = '';
            ui.elements.taskDueDate.value = '';
            ui.elements.taskAssignee.value = '';

            // Close the modal
            ui.elements.taskModal.classList.remove('active');
            ui.elements.taskModal.classList.add('hidden');

            // Reload both task table and employee data after creating a task
            ui.loadEmployerData();
        } else {
            alert('Failed to create task: ' + (result?.error?.message || 'Unknown error'));
        }
    });

    // Create random task button (for employers)
    ui.elements.createRandomTaskBtn.addEventListener('click', async () => {
        // Get a random employee if possible
        let randomEmployeeId = null;

        // Try to get a random employee from the assignee dropdown
        const assigneeOptions = ui.elements.taskAssignee.querySelectorAll('option');
        if (assigneeOptions.length > 1) { // Skip the first "Select Employee" option
            const randomIndex = Math.floor(Math.random() * (assigneeOptions.length - 1)) + 1;
            randomEmployeeId = parseInt(assigneeOptions[randomIndex].value, 10);
        }

        const randomTask = tasks.generateRandomTask(randomEmployeeId);
        const result = await tasks.createTask(randomTask);

        if (result && result.status === 'success') {
            // Close the modal
            ui.elements.taskModal.classList.remove('active');
            ui.elements.taskModal.classList.add('hidden');

            // Reload both task table and employee data after creating a random task
            ui.loadEmployerData();
        } else {
            alert('Failed to create random task: ' + (result?.error?.message || 'Unknown error'));
        }
    });

    // Add change event listeners to filter elements for automatic filtering
    // When filter changes, only update the tasks table, not the entire dashboard
    ui.elements.taskStatusFilter.addEventListener('change', () => ui.loadTaskTable());
    ui.elements.taskAssigneeFilter.addEventListener('change', () => ui.loadTaskTable());
    ui.elements.taskSortBy.addEventListener('change', () => ui.loadTaskTable());
    ui.elements.taskSortOrder.addEventListener('change', () => ui.loadTaskTable());

    // Employee task filter controls
    ui.elements.employeeTaskStatusFilter.addEventListener('change', () => ui.loadEmployeeAssignedTasks());
    ui.elements.employeeTaskSortBy.addEventListener('change', () => ui.loadEmployeeAssignedTasks());
    ui.elements.employeeTaskSortOrder.addEventListener('change', () => ui.loadEmployeeAssignedTasks());

    // Close assign task modal button
    ui.elements.closeAssignModalBtn.addEventListener('click', () => {
        ui.closeAssignTaskModal();
    });

    // Close modal when clicking outside
    ui.elements.assignTaskModal.addEventListener('click', (event) => {
        if (event.target === ui.elements.assignTaskModal) {
            ui.closeAssignTaskModal();
        }
    });

    // Confirm assign task button
    ui.elements.confirmAssignTaskBtn.addEventListener('click', async () => {
        const taskId = ui.elements.assignTaskId.value;
        const assigneeId = ui.elements.assignTaskAssignee.value;

        if (!assigneeId) {
            alert('Please select an employee to assign the task to.');
            return;
        }

        const result = await tasks.assignTask(taskId, parseInt(assigneeId));

        if (result && result.status === 'success') {
            // Close the modal
            ui.closeAssignTaskModal();

            // Reload the task list to see the updated assignment
            ui.loadTaskTable();
        } else {
            alert('Failed to assign task: ' + (result?.error?.message || 'Unknown error'));
        }
    });
}
