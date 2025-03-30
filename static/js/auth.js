// auth.js - Handle authentication logic

const AUTH_KEY = 'task_management_auth';
const USERS_KEY = 'created_users';

const auth = {
    // Save authentication data to localStorage
    saveAuth: (userData, token, tokenExpiry) => {
        const authData = {
            user: userData,
            token: token,
            tokenExpiry: tokenExpiry
        };
        localStorage.setItem(AUTH_KEY, JSON.stringify(authData));
        return authData;
    },

    // Get current auth data from localStorage
    getAuth: () => {
        const authData = localStorage.getItem(AUTH_KEY);
        if (!authData) return null;

        try {
            return JSON.parse(authData);
        } catch (error) {
            console.error('Error parsing auth data:', error);
            return null;
        }
    },

    // Clear auth data from localStorage
    clearAuth: () => {
        localStorage.removeItem(AUTH_KEY);
    },

    // Check if user is authenticated
    isAuthenticated: () => {
        const authData = auth.getAuth();
        if (!authData) return false;

        // Check token expiry
        if (authData.tokenExpiry && authData.tokenExpiry < Date.now() / 1000) {
            auth.clearAuth();
            return false;
        }

        return true;
    },

    // Get current user data
    getCurrentUser: () => {
        const authData = auth.getAuth();
        return authData ? authData.user : null;
    },

    // Get auth token
    getToken: () => {
        const authData = auth.getAuth();
        return authData ? authData.token : null;
    },

    // Save created user data
    saveCreatedUser: (userData, password) => {
        let users = auth.getCreatedUsers() || [];

        // Add password to user data for demo purposes
        const userWithPassword = {
            ...userData,
            password: password
        };

        users.push(userWithPassword);
        localStorage.setItem(USERS_KEY, JSON.stringify(users));
        return users;
    },

    // Get all created users
    getCreatedUsers: () => {
        const usersData = localStorage.getItem(USERS_KEY);
        if (!usersData) return [];

        try {
            return JSON.parse(usersData);
        } catch (error) {
            console.error('Error parsing users data:', error);
            return [];
        }
    }
};
