// Глобальные переменные и конфигурация
window.APP_CONFIG = {
    API_BASE: '/api',
    TOKEN_KEY: 'token'
};

window.auth = {
    token: localStorage.getItem(window.APP_CONFIG.TOKEN_KEY)
};