// Модуль тестирования
window.testing = {
    async testProtectedEndpoint() {
        const resultEl = document.getElementById('test-result');
        this.hideResult(resultEl);

        try {
            const response = await fetch('/api/test', {
                headers: { 'Authorization': 'Bearer ' + window.auth.token }
            });

            if (response.ok) {
                const text = await response.text();
                this.showResult(resultEl, '✅ ' + text, 'success');
            } else {
                this.showResult(resultEl, '❌ Доступ запрещен. Токен невалиден', 'error');
            }
        } catch (error) {
            this.showResult(resultEl, '❌ Ошибка сети при тестировании', 'error');
        }
    },

    async testEmailConnection() {
        const resultEl = document.getElementById('test-result');
        this.hideResult(resultEl);

        try {
            const response = await fetch('/api/test-email', {
                method: 'POST',
                headers: { 
                    'Authorization': 'Bearer ' + window.auth.token,
                    'Content-Type': 'application/json'
                }
            });

            const data = await response.json();
            if (response.ok) {
                this.showResult(resultEl, '✅ ' + data.message, 'success');
            } else {
                this.showResult(resultEl, '❌ ' + data.message, 'error');
            }
        } catch (error) {
            this.showResult(resultEl, '❌ Ошибка сети при тестировании почты', 'error');
        }
    },

    async sendTestDigest() {
        const resultEl = document.getElementById('test-result');
        this.hideResult(resultEl);

        try {
            const response = await fetch('/api/send-test-digest', {
                method: 'POST',
                headers: { 
                    'Authorization': 'Bearer ' + window.auth.token,
                    'Content-Type': 'application/json'
                }
            });

            const data = await response.json();
            if (response.ok) {
                this.showResult(resultEl, '✅ ' + data.message + ' Проверьте вашу почту!', 'success');
            } else {
                this.showResult(resultEl, '❌ ' + data.message, 'error');
            }
        } catch (error) {
            this.showResult(resultEl, '❌ Ошибка сети при отправке дайджеста', 'error');
        }
    },

    showResult(element, message, type) {
        element.textContent = message;
        element.className = type;
        element.classList.remove('hidden');
    },

    hideResult(element) {
        element.classList.add('hidden');
    }
};