<template>
  <div class="class-selector">
    <label for="class-filter">Filter by Class:</label>
    <select id="class-filter" v-model="selectedClass" @change="onClassChange">
      <option v-for="classId in classIds" :key="classId" :value="classId">
        Class {{ classId }}
      </option>
    </select>
  </div>
</template>

<script>
export default {
  props: {
    classIds: {
      type: Array,
      required: true, // Expect an array of available class IDs
    },
  },
  data() {
    return {
      selectedClass: this.classIds[0], // Default to the first class in the list
    };
  },
  methods: {
    onClassChange() {
      this.$emit("class-filter", this.selectedClass); // Emit the selected class ID to the parent
    },
  },
  watch: {
    classIds(newClassIds) {
      if (!newClassIds.includes(this.selectedClass)) {
        // If the current selected class is no longer valid, reset to the first in the updated list
        this.selectedClass = newClassIds[0];
        this.onClassChange();
      }
    },
  },
  created() {
    // Emit the default selected class when the component is created
    if (this.classIds.length > 0) {
      this.onClassChange();
    }
  },
};
</script>

<style scoped>
.class-selector {
  margin-bottom: 1rem;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

label {
  font-weight: bold;
}

select {
  padding: 0.5rem;
  border: 1px solid #ccc;
  border-radius: 4px;
  font-size: 1rem;
}
</style>
