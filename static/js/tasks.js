// tasks.js - Manage task-related functionality

const tasks = {
    // Create a new task
    createTask: async (taskData) => {
        const token = auth.getToken();
        if (!token) return null;

        return await api.createTask(taskData, token);
    },

    // Generate a random task for demo purposes
    generateRandomTask: (employeeId = null) => {
        const taskTitles = [
            'Create user authentication',
            'Implement file upload',
            'Fix navigation bug',
            'Design new dashboard',
            'Update documentation',
            'Add dark mode',
            'Optimize database queries',
            'Create mobile responsive layout',
            'Add unit tests',
            'Refactor legacy code'
        ];

        const taskDescriptions = [
            'This task involves implementing the specified feature or fixing the issue.',
            'Please complete this task as soon as possible.',
            'This is a high priority task that needs immediate attention.',
            'Coordinate with the team when implementing this feature.',
            'Make sure to follow the coding standards when implementing this feature.',
            'This task requires thorough testing before deployment.'
        ];

        const statuses = ['pending', 'in_progress', 'completed'];

        // Generate random date in the next 7 days
        const dueDate = new Date();
        dueDate.setDate(dueDate.getDate() + Math.floor(Math.random() * 7) + 1);

        return {
            title: taskTitles[Math.floor(Math.random() * taskTitles.length)],
            description: taskDescriptions[Math.floor(Math.random() * taskDescriptions.length)],
            status: statuses[Math.floor(Math.random() * statuses.length)],
            due_date: dueDate.toISOString(),
            assignee_id: employeeId
        };
    },

    // Update a task's status
    updateTaskStatus: async (taskId, status) => {
        const token = auth.getToken();
        if (!token) return null;

        return await api.updateTaskStatus(taskId, status, token);
    },

    // Get tasks assigned to the current user (for employees)
    getAssignedTasks: async (params = {}) => {
        const token = auth.getToken();
        if (!token) return null;

        return await api.getAssignedTasks(token, params);
    },

    // Get all tasks created by the current user (for employers)
    getTasks: async (params = {}) => {
        const token = auth.getToken();
        if (!token) return null;

        return await api.getTasks(token, params);
    },

    // Get summary of all employees' tasks (for employers)
    getEmployeeSummary: async () => {
        const token = auth.getToken();
        if (!token) return null;

        return await api.getEmployeeSummary(token);
    },

    // Format a task status for display
    formatStatus: (status) => {
        switch (status) {
            case 'pending': return 'Pending';
            case 'in_progress': return 'In Progress';
            case 'completed': return 'Completed';
            default: return status;
        }
    }
};
