// api.js - Handle API requests

const API_BASE_URL = 'http://localhost:8080/api/v1';

const api = {
    // User related API calls
    registerUser: async (userData) => {
        try {
            const response = await fetch(`${API_BASE_URL}/users`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(userData)
            });
            return await response.json();
        } catch (error) {
            console.error('Error registering user:', error);
            return { status: 'error', error: { message: error.message } };
        }
    },

    loginUser: async (credentials) => {
        try {
            const response = await fetch(`${API_BASE_URL}/login`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(credentials)
            });
            return await response.json();
        } catch (error) {
            console.error('Error logging in:', error);
            return { status: 'error', error: { message: error.message } };
        }
    },

    getCurrentUser: async (token) => {
        try {
            const response = await fetch(`${API_BASE_URL}/users/me`, {
                headers: {
                    'Authorization': `Bearer ${token}`
                }
            });
            return await response.json();
        } catch (error) {
            console.error('Error getting current user:', error);
            return { status: 'error', error: { message: error.message } };
        }
    },

    getUserById: async (userId, token) => {
        try {
            const response = await fetch(`${API_BASE_URL}/users/${userId}`, {
                headers: {
                    'Authorization': `Bearer ${token}`
                }
            });
            return await response.json();
        } catch (error) {
            console.error('Error getting user:', error);
            return { status: 'error', error: { message: error.message } };
        }
    },

    // Task related API calls
    createTask: async (taskData, token) => {
        try {
            const response = await fetch(`${API_BASE_URL}/tasks`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify(taskData)
            });
            return await response.json();
        } catch (error) {
            console.error('Error creating task:', error);
            return { status: 'error', error: { message: error.message } };
        }
    },

    updateTaskStatus: async (taskId, status, token) => {
        try {
            const response = await fetch(`${API_BASE_URL}/tasks/${taskId}/status`, {
                method: 'PATCH',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify({ status })
            });
            return await response.json();
        } catch (error) {
            console.error('Error updating task:', error);
            return { status: 'error', error: { message: error.message } };
        }
    },

    getAssignedTasks: async (token, params = {}) => {
        try {
            const queryParams = new URLSearchParams(params).toString();
            const url = `${API_BASE_URL}/tasks/assigned${queryParams ? `?${queryParams}` : ''}`;
            const response = await fetch(url, {
                headers: {
                    'Authorization': `Bearer ${token}`
                }
            });
            return await response.json();
        } catch (error) {
            console.error('Error getting assigned tasks:', error);
            return { status: 'error', error: { message: error.message } };
        }
    },

    getTasks: async (token, params = {}) => {
        try {
            const queryParams = new URLSearchParams(params).toString();
            const url = `${API_BASE_URL}/tasks${queryParams ? `?${queryParams}` : ''}`;
            const response = await fetch(url, {
                headers: {
                    'Authorization': `Bearer ${token}`
                }
            });
            return await response.json();
        } catch (error) {
            console.error('Error getting tasks:', error);
            return { status: 'error', error: { message: error.message } };
        }
    },

    getEmployeeSummary: async (token) => {
        try {
            const response = await fetch(`${API_BASE_URL}/employee-summary`, {
                headers: {
                    'Authorization': `Bearer ${token}`
                }
            });
            return await response.json();
        } catch (error) {
            console.error('Error getting employee summary:', error);
            return { status: 'error', error: { message: error.message } };
        }
    },

    // Assign task to an employee
    assignTask: async (taskId, assigneeId, token) => {
        try {
            const response = await fetch(`${API_BASE_URL}/tasks/${taskId}/assign`, {
                method: 'PATCH',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify({ assignee_id: assigneeId })
            });
            return await response.json();
        } catch (error) {
            console.error('Error assigning task:', error);
            return { status: 'error', error: { message: error.message } };
        }
    },

    // Get all users (for task assignment dropdown)
    getAllUsers: async (token) => {
        try {
            const response = await fetch(`${API_BASE_URL}/users/all`, {
                headers: {
                    'Authorization': `Bearer ${token}`
                }
            });
            return await response.json();
        } catch (error) {
            console.error('Error getting all users:', error);
            return { status: 'error', error: { message: error.message } };
        }
    }
};
