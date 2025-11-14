// Модуль настроек
window.settings = {
    async loadSettings() {
        try {
            const response = await fetch('/api/settings', {
                headers: { 'Authorization': 'Bearer ' + window.auth.token }
            });

            if (response.ok) {
                const settings = await response.json();
                this.updateForm(settings);
            }
        } catch (error) {
            console.error('Error loading settings:', error);
        }
    },

    updateForm(settings) {
        const elements = {
            'imap-server': settings.imap_server || 'imap.gmail.com:993',
            'email': settings.email || '',
            'schedule': settings.schedule || '09:00'
        };

        Object.entries(elements).forEach(([id, value]) => {
            const element = document.getElementById(id);
            if (element) element.value = value;
        });
    },

    async saveSettings() {
        const settingsData = {
            imap_server: document.getElementById('imap-server').value,
            email: document.getElementById('email').value,
            app_password: document.getElementById('app-password').value,
            schedule: document.getElementById('schedule').value
        };

        const messageEl = document.getElementById('settings-message');
        this.hideMessage(messageEl);

        try {
            const response = await fetch('/api/settings', {
                method: 'POST',
                headers: { 
                    'Content-Type': 'application/json',
                    'Authorization': 'Bearer ' + window.auth.token
                },
                body: JSON.stringify(settingsData)
            });

            if (response.ok) {
                this.showMessage(messageEl, '✅ Настройки успешно сохранены!', 'success');
            } else {
                this.showMessage(messageEl, '❌ Ошибка сохранения настроек', 'error');
            }
        } catch (error) {
            this.showMessage(messageEl, '❌ Ошибка сети при сохранении', 'error');
        }
    },

    showMessage(element, message, type) {
        element.textContent = message;
        element.className = type;
        element.classList.remove('hidden');
    },

    hideMessage(element) {
        element.classList.add('hidden');
    }
};