<script setup>
import { ref } from 'vue';
import { api } from '../services/api';
import { notification } from '../services/notification';
import BaseButton from '../components/ui/BaseButton.vue';

const emit = defineEmits(['login-success']);

const username = ref('');
const password = ref('');
const isLoading = ref(false);
const error = ref('');

const handleLogin = async () => {
  if (!username.value || !password.value) {
    error.value = 'Please enter both username and password.';
    return;
  }

  isLoading.value = true;
  error.value = '';

  try {
    const response = await api.login(username.value, password.value);
    console.log('[Login] Received response:', response);
    
    // Store token
    if (response.token) {
      localStorage.setItem('auth_token', response.token);
    }
    
    if (!response.user) {
        throw new Error('User data missing from response');
    }

    emit('login-success', {
      user: response.user,
      role: response.user.role || 'manager'
    });
    notification.success(`Welcome back, ${response.user.username}`);
  } catch (err) {
    error.value = err.message || 'Login failed';
    notification.error(error.value);
  } finally {
    isLoading.value = false;
  }
};
</script>

<template>
  <div class="login-page">
    <div class="background-elements">
      <div class="blob blob-1"></div>
      <div class="blob blob-2"></div>
      <div class="blob blob-3"></div>
    </div>

    <div class="login-container">
      <div class="login-card">
        <div class="login-header">
          <div class="logo">
            <i class="ri-shield-flash-line"></i>
            <span>CPD Nexus</span>
          </div>
          <h1 class="login-title">Welcome Back</h1>
          <p class="login-subtitle">Secure access to construction management</p>
        </div>

        <form @submit.prevent="handleLogin" class="login-form">
          <div v-if="error" class="error-banner">
            <i class="ri-error-warning-line"></i>
            {{ error }}
          </div>

          <div class="form-group">
            <label for="username">Username</label>
            <div class="input-wrapper">
              <i class="ri-user-line"></i>
              <input 
                id="username" 
                v-model="username" 
                type="text" 
                placeholder="Enter your username" 
                required
                :disabled="isLoading"
              />
            </div>
          </div>

          <div class="form-group">
            <label for="password">Password</label>
            <div class="input-wrapper">
              <i class="ri-lock-2-line"></i>
              <input 
                id="password" 
                v-model="password" 
                type="password" 
                placeholder="••••••••" 
                required
                :disabled="isLoading"
              />
            </div>
          </div>

          <div class="form-options">
            <label class="checkbox-label">
              <input type="checkbox" />
              <span>Remember me</span>
            </label>
            <a href="#" class="forgot-link">Forgot password?</a>
          </div>

          <BaseButton 
            type="submit" 
            class="submit-btn" 
            :loading="isLoading"
            block
          >
            Sign In
          </BaseButton>
        </form>

        <div class="login-footer">
          <p>Don't have an account? <a href="#">Contact Support</a></p>
        </div>
      </div>
      
      <div class="login-info">
        <div class="info-item">
          <div class="info-icon admin"><i class="ri-admin-line"></i></div>
          <div class="info-text">
            <h4>Manager Access</h4>
            <p>Use any email with 'manager' to test Manager views.</p>
          </div>
        </div>
        <div class="info-item">
          <div class="info-icon client"><i class="ri-user-follow-line"></i></div>
          <div class="info-text">
            <h4>Client Access</h4>
            <p>Use any other email to test Client/Project Manager views.</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--color-bg);
  position: relative;
  overflow: hidden;
  font-family: 'Inter', sans-serif;
}

/* Background Blobs */
.background-elements {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 1;
}

.blob {
  position: absolute;
  filter: blur(80px);
  border-radius: 50%;
  opacity: 0.15;
  animation: float 20s infinite alternate;
}

.blob-1 {
  width: 500px;
  height: 500px;
  background: var(--color-accent);
  top: -100px;
  left: -100px;
}

.blob-2 {
  width: 400px;
  height: 400px;
  background: #10b981;
  bottom: -50px;
  right: -50px;
  animation-duration: 25s;
}

.blob-3 {
  width: 300px;
  height: 300px;
  background: #f59e0b;
  bottom: 40%;
  left: 30%;
  animation-duration: 30s;
}

@keyframes float {
  0% { transform: translate(0, 0) scale(1); }
  100% { transform: translate(50px, 50px) scale(1.1); }
}

.login-container {
  position: relative;
  z-index: 2;
  width: 100%;
  max-width: 420px;
  padding: 20px;
}

.login-card {
  background: rgba(255, 255, 255, 0.03);
  backdrop-filter: blur(12px);
  border: 1px solid rgba(255, 255, 255, 0.05);
  border-radius: var(--radius-lg);
  padding: 40px;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.5);
}

.login-header {
  text-align: center;
  margin-bottom: 32px;
}

.logo {
  display: inline-flex;
  align-items: center;
  gap: 12px;
  font-size: 24px;
  font-weight: 700;
  color: var(--color-text-primary);
  margin-bottom: 24px;
}

.logo i {
  font-size: 32px;
  color: var(--color-accent);
}

.login-title {
  font-size: 28px;
  font-weight: 700;
  color: var(--color-text-primary);
  margin-bottom: 8px;
  letter-spacing: -0.5px;
}

.login-subtitle {
  font-size: 14px;
  color: var(--color-text-secondary);
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.error-banner {
  background: rgba(239, 68, 68, 0.1);
  border: 1px solid rgba(239, 68, 68, 0.2);
  color: #ef4444;
  padding: 12px;
  border-radius: var(--radius-sm);
  font-size: 13px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-group label {
  font-size: 13px;
  font-weight: 500;
  color: var(--color-text-secondary);
}

.input-wrapper {
  position: relative;
  display: flex;
  align-items: center;
}

.input-wrapper i {
  position: absolute;
  left: 14px;
  color: var(--color-text-muted);
  font-size: 18px;
}

.input-wrapper input {
  width: 100%;
  padding: 12px 14px 12px 42px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: var(--radius-md);
  color: var(--color-text-primary);
  font-size: 14px;
  transition: all var(--transition-fast);
}

.input-wrapper input:focus {
  outline: none;
  border-color: var(--color-accent);
  background: rgba(255, 255, 255, 0.08);
  box-shadow: 0 0 0 4px rgba(59, 130, 246, 0.1);
}

.form-options {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-size: 13px;
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--color-text-secondary);
  cursor: pointer;
}

.forgot-link {
  color: var(--color-accent);
  text-decoration: none;
  font-weight: 500;
}

.forgot-link:hover {
  text-decoration: underline;
}

.login-footer {
  margin-top: 32px;
  text-align: center;
  font-size: 14px;
  color: var(--color-text-secondary);
}

.login-footer a {
  color: var(--color-accent);
  text-decoration: none;
  font-weight: 600;
}

.login-info {
  margin-top: 24px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.info-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.05);
  border-radius: var(--radius-md);
}

.info-icon {
  width: 36px;
  height: 36px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
}

.info-icon.admin { background: rgba(59, 130, 246, 0.1); color: var(--color-accent); }
.info-icon.client { background: rgba(16, 185, 129, 0.1); color: #10b981; }

.info-text h4 {
  font-size: 13px;
  font-weight: 600;
  color: var(--color-text-primary);
  margin-bottom: 2px;
}

.info-text p {
  font-size: 12px;
  color: var(--color-text-secondary);
}
</style>
