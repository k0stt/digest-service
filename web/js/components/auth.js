// Модуль аутентификации
window.auth = {
    ...window.auth,
    
    async login() {
        const username = document.getElementById('login-username').value;
        const password = document.getElementById('login-password').value;
        const errorEl = document.getElementById('login-error');
        
        this.hideError(errorEl);
        
        try {
            const response = await fetch('/api/login', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ username, password })
            });

            if (response.ok) {
                const data = await response.json();
                this.token = data.token;
                localStorage.setItem('token', this.token);
                
                // Устанавливаем имя пользователя
                this.setUsername(username);
                window.app.showApp();
                window.settings.loadSettings();
            } else {
                this.showError(errorEl, '❌ Ошибка входа. Проверьте логин и пароль');
            }
        } catch (error) {
            this.showError(errorEl, '❌ Ошибка сети. Проверьте подключение');
        }
    },

    async register() {
        const username = document.getElementById('register-username').value;
        const password = document.getElementById('register-password').value;
        const email = document.getElementById('register-email').value;
        const errorEl = document.getElementById('register-error');
        
        this.hideError(errorEl);
        
        try {
            const response = await fetch('/api/register', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ username, password, email })
            });

            if (response.ok) {
                const data = await response.json();
                this.token = data.token;
                localStorage.setItem('token', this.token);
                
                // Устанавливаем имя пользователя
                this.setUsername(username);
                window.app.showApp();
            } else {
                this.showError(errorEl, '❌ Ошибка регистрации. Возможно, пользователь уже существует');
            }
        } catch (error) {
            this.showError(errorEl, '❌ Ошибка сети. Проверьте подключение');
        }
    },

    setUsername(username) {
        const usernameElement = document.getElementById('username');
        if (usernameElement) {
            usernameElement.textContent = username;
        }
    },

    logout() {
        this.token = null;
        localStorage.removeItem('token');
        document.getElementById('auth-container').classList.remove('hidden');
        document.getElementById('app-container').classList.add('hidden');
        
        // Очищаем поля форм
        this.clearForms();
    },

    clearForms() {
        document.getElementById('login-username').value = '';
        document.getElementById('login-password').value = '';
        document.getElementById('register-username').value = '';
        document.getElementById('register-password').value = '';
        document.getElementById('register-email').value = '';
    },

    showError(element, message) {
        element.textContent = message;
        element.classList.remove('hidden');
    },

    hideError(element) {
        element.classList.add('hidden');
    }
};