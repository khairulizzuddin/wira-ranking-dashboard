<template>
  <div>
    <h1 class="dashboard-header">WIRA's Ranking Dashboard</h1>

    <!-- Filter by Class -->
    <ClassSelector :classIds="classIds" @class-filter="filterByClass" />

    <!-- Search Box -->
    <input
      v-model="searchQuery"
      @input="fetchRanks"
      placeholder="Search by username"
      class="search-box"
    />

    <!-- Rankings Table -->
    <RankingTable
      :rankings="rankings"
      :currentPage="currentPage"
      :limit="limit"
      @page-change="fetchRanks"
    />

    <!-- Pagination -->
    <div class="pagination">
      <button @click="previousPage" :disabled="currentPage <= 1">Previous</button>
      <span>Page {{ currentPage }}</span>
      <button @click="nextPage" :disabled="rankings.length < limit">Next</button>
    </div>

    <button class="logout-btn" @click="logout">Logout</button>
  </div>
</template>

<script>
import ClassSelector from '@/components/dashboard/ClassSelector.vue';
import RankingTable from '@/components/dashboard/RankingTable.vue';
import { useAuthStore } from '@/stores/auth';
import axios from 'axios';

export default {
  components: { ClassSelector, RankingTable },
  data() {
    return {
      classIds: [1,2,3,4,5,6,7,8], // Fetch class IDs dynamically if needed
      rankings: [],
      currentPage: 1,
      limit: 10,
      selectedClass: 1,
      searchQuery: "",
    };
  },
  computed: {
    authToken() {
      return useAuthStore().authToken;
    },
  },
  methods: {
    async fetchRanks() {
      this.loading = true;

      try {
        const token = localStorage.getItem('wira_token');
        const response = await axios.get('http://localhost:8080/dashboard', {
          params: {
            class_id: this.selectedClass,
            page: this.currentPage,
            limit: this.limit,
            search: this.searchQuery,
          },
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });
        this.rankings = response.data;
      } catch (error) {
        console.error('Failed to fetch rankings:', error);
        alert('Failed to fetch rankings. Please try again.');
      } finally {
        this.loading = false;
      }
    },
    filterByClass(classId) {
      this.selectedClass = classId;
      this.currentPage = 1; // Reset to page 1 when filtering
      this.fetchRanks();
    },
    nextPage() {
      this.currentPage++;
      this.fetchRanks();
    },
    previousPage() {
      if (this.currentPage > 1) {
        this.currentPage--;
        this.fetchRanks();
      }
    },
    logout() {
      useAuthStore().logout(); // Clear token
      this.$router.push("/login");
    },
  },
  created() {
    this.fetchRanks();
  },
};
</script>

<style scoped>
.dashboard-header {
  font-size: 2rem;
  font-weight: bold;
  text-align: center;
  margin: 5px 0;
}

.search-box {
  margin-bottom: 1rem;
  padding: 0.5rem;
  width: 100%;
  max-width: 300px;
}

.pagination {
  margin-top: 1rem;
  display: flex;
  justify-content: center;
  gap: 1rem;
}

.pagination button:disabled {
  background-color: #ccc;
}

.logout-btn {
  margin-top: 20px;
  padding: 10px 20px;
  background-color: #ff4d4f;
  color: white;
  border: none;
  border-radius: 5px;
  cursor: pointer;
}
</style>
