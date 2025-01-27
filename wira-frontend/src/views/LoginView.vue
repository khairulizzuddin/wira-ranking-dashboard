<template>
  <div class="login-container">
    <div class="login-box">
      <h1 class="login-title">WIRA Login</h1>
      <LoginForm />
      
      <!-- Error display for general errors -->
      <div v-if="authStore.errors.general" class="error-message">
        {{ authStore.errors.general }}
      </div>
    </div>
  </div>
</template>

<script setup>
import { useAuthStore } from '@/stores/auth';
import LoginForm from '@/components/auth/LoginForm.vue';

const authStore = useAuthStore();
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh; 
  background-color: #f3f4f6;
}

.login-box {
  width: 100%;
  max-width: 400px;
  padding: 2rem;
  background-color: #ffffff; 
  border-radius: 12px; 
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1); 
}

.login-title {
  font-size: 2rem;
  font-weight: bold;
  text-align: center;
  margin-bottom: 1.5rem;
  color: #1f2937; 
}

.error-message {
  color: #ef4444;
  text-align: center;
  margin-top: 1rem;
  font-weight: 500;
}
</style>
