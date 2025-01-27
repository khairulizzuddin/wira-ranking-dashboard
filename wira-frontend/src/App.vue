<template>
  <router-view v-slot="{ Component }">
    <template v-if="isAuthenticated">
      <div class="app-container">
        <main class="main-content">
          <component :is="Component" />
        </main>
      </div>
    </template>
    
    <component v-else :is="Component" />
  </router-view>
</template>

<script setup>
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
const router = useRouter()

const isAuthenticated = computed(() => {
  return localStorage.getItem('wira_token') !== null
})

const isDashboard = computed(() => {
  return router.currentRoute.value.path === '/dashboard'
})

const logout = () => {
  localStorage.removeItem('wira_token')
  authStore.logout()
  router.push('/login')
}
</script>

<style scoped>
.app-container {
  display: flex;
  min-height: 100vh;
}

/* .sidebar {
  width: 250px;
  background: white;
  box-shadow: 2px 0 5px rgba(0,0,0,0.1);
  padding: 1rem;
} */

/* .sidebar-content {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
} */

.main-content {
  flex: 1;
  padding: 2rem;
  background: #f8fafc;
}

.nav-link {
  display: block;
  padding: 0.75rem;
  border-radius: 4px;
  color: #334155;
  text-decoration: none;
  transition: background-color 0.2s;
}

.nav-link:hover {
  background-color: #f1f5f9;
}

.nav-link.router-link-exact-active {
  background-color: #e2e8f0;
  font-weight: 500;
}

.logout-button {
  width: 100%;
  padding: 0.75rem;
  text-align: left;
  border: none;
  background: none;
  color: #ef4444;
  cursor: pointer;
  border-radius: 4px;
  transition: background-color 0.2s;
}

.logout-button:hover {
  background-color: #fee2e2;
}
</style>
