<template>
  <form class="auth-form" @submit.prevent="handleSubmit">
    <div class="form-group">
      <label>Username</label>
      <input 
        v-model="username" 
        class="form-input"
        :class="{ 'error-border': authStore.errors.username }"
      />
      <p class="error-message">{{ authStore.errors.username }}</p>
    </div>
    
    <div class="form-group">
      <label>Password</label>
      <input 
        v-model="password" 
        type="password" 
        class="form-input"
        :class="{ 'error-border': authStore.errors.password }"
      />
      <p class="error-message">{{ authStore.errors.password }}</p>
    </div>

    <div class="form-group">
      <label>2FA Code</label>
      <input 
        v-model="twoFACode" 
        class="form-input"
        :class="{ 'error-border': authStore.errors.twoFA }"
        placeholder="6-digit code"
      />
      <p class="error-message">{{ authStore.errors.twoFA }}</p>
    </div>
    
    <button 
      class="primary-button"
      :disabled="isSubmitting"
    >
      {{ isSubmitting ? 'Logging in...' : 'Login' }}
    </button>
  </form>
</template>

<script setup>
import { ref } from 'vue';
import { useAuthStore } from '@/stores/auth';
import { useRouter } from 'vue-router';

const router = useRouter();
const authStore = useAuthStore();
const username = ref('testuser');
const password = ref('testpassword');
const twoFACode = ref('123456'); // Default test code
const isSubmitting = ref(false);

const handleSubmit = async () => {
  isSubmitting.value = true;
  try {
    const success = await authStore.login({
      username: username.value,
      password: password.value,
      twoFA_code: twoFACode.value,
    });
    
    if (success) {
      router.push('/dashboard');
    }
  } finally {
    isSubmitting.value = false;
  }
};
</script>

<style scoped>
.auth-form {
  width: 100%;
  max-width: 400px; /* Matches the login box width */
  margin: 0 auto;
}

.form-group {
  margin-bottom: 1.5rem;
}

/* Input fields styled to match the login box width */
.form-input {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #d1d5db; /* Light gray border */
  border-radius: 0.375rem;
  font-size: 1rem;
  box-sizing: border-box; /* Ensure padding doesn't affect width */
}

.form-input:focus {
  border-color: #3b82f6; /* Blue border on focus */
  outline: none;
  box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.5); /* Subtle focus shadow */
}

.error-border {
  border-color: #ef4444; /* Red border for errors */
}

/* Error message styling */
.error-message {
  color: #ef4444;
  font-size: 0.875rem;
  margin-top: 0.25rem;
}

/* Button styling */
.primary-button {
  width: 100%;
  padding: 0.75rem;
  background-color: #3b82f6;
  color: white;
  border: none;
  border-radius: 0.375rem;
  cursor: pointer;
  transition: background-color 0.2s;
}

.primary-button:hover {
  background-color: #2563eb;
}

.primary-button:disabled {
  background-color: #94a3b8;
  cursor: not-allowed;
}
</style>
