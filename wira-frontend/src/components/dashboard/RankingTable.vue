<template>
  <div>
    <table>
      <thead>
        <tr>
          <th>Username</th>
          <th>Class ID</th>
          <th>Rank</th>
          <th>Reward Score</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="(rank, index) in rankings" :key="rank.username">
          <td>{{ rank.username }}</td>
          <td>{{ rank.class_id }}</td>
          <td>{{ (currentPage - 1) * limit + index + 1 }}</td>
          <td>{{ rank.reward_score }}</td>
        </tr>
      </tbody>
    </table>

    <!-- <div class="pagination">
      <button :disabled="currentPage === 1" @click="changePage(currentPage - 1)">
        Previous
      </button>
      <span>Page {{ currentPage }}</span>
      <button :disabled="rankings.length < limit" @click="changePage(currentPage + 1)">
        Next
      </button>
    </div> -->
  </div>
</template>

<script>
export default {
  props: {
    rankings: Array, // Array of ranking objects
    currentPage: Number,
    limit: Number,
  },
  methods: {
    changePage(newPage) {
      this.$emit("page-change", newPage);
    },
  },
};
</script>

<style scoped>
table {
  width: 100%;
  border-collapse: collapse;
  margin-bottom: 20px;
}

th, td {
  border: 1px solid #ddd;
  padding: 8px;
}

.pagination {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
