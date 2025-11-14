// Основной модуль приложения
window.app = {
    init() {
        if (window.auth.token) {
            this.showApp();
            window.settings.loadSettings();
        }
        
        // Добавляем обработчики для кнопок
        this.bindEvents();
    },

    bindEvents() {
        // Обработчики для Enter в формах
        document.getElementById('login-password')?.addEventListener('keypress', (e) => {
            if (e.key === 'Enter') window.auth.login();
        });
        
        document.getElementById('register-password')?.addEventListener('keypress', (e) => {
            if (e.key === 'Enter') window.auth.register();
        });
    },

    showApp() {
        document.getElementById('auth-container').classList.add('hidden');
        document.getElementById('app-container').classList.remove('hidden');
        this.showTab('settings');
    },

    showTab(tabName) {
        // Скрываем все табы
        ['settings', 'test'].forEach(tab => {
            const element = document.getElementById(tab + '-tab');
            if (element) element.classList.add('hidden');
        });
        
        // Убираем активный класс со всех кнопок
        document.querySelectorAll('.tab-btn').forEach(btn => {
            btn.classList.remove('active');
        });
        
        // Показываем выбранный таб и активируем кнопку
        const targetTab = document.getElementById(tabName + '-tab');
        const targetBtn = event?.target;
        
        if (targetTab) targetTab.classList.remove('hidden');
        if (targetBtn) targetBtn.classList.add('active');
    }
};