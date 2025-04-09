<script setup>
import { ref } from 'vue'

const emit = defineEmits(['update:searchQuery', 'filterChange'])

const searchQuery = ref('')
const filterContainerActive = ref(false)
const activeFilters = ref([])

const toggleFilterContainer = () => {
  filterContainerActive.value = !filterContainerActive.value
}

const handleFilterClick = (filter) => {
  if (activeFilters.value.includes(filter)) {
    activeFilters.value = activeFilters.value.filter(f => f !== filter)
  } else {
    activeFilters.value.push(filter)
  }
  emit('filterChange', activeFilters.value)
}

const updateSearch = (e) => {
  searchQuery.value = e.target.value
  emit('update:searchQuery', searchQuery.value)
}
</script>

<template>
  <div class="search-bar">
    <i class="fa-solid fa-search"></i>
    <input 
      type="text" 
      placeholder="Search students..." 
      :value="searchQuery"
      @input="updateSearch"
    >
    <button class="filter-toggle-btn" @click="toggleFilterContainer">
      <i class="fa-solid fa-filter"></i>
      Filter
    </button>

    <div class="filter-tooltip" :class="{ active: filterContainerActive }">
      <div class="filter-category">
        <span class="category-label">Year Level</span>
        <div class="filter-buttons-row">
          <button v-for="year in ['1st', '2nd', '3rd', '4th']" :key="year"
            class="filter-button" :class="{ active: activeFilters.includes(year) }"
            @click="handleFilterClick(year)">
            {{ year }}
          </button>
        </div>
      </div>

      <div class="filter-category">
        <span class="category-label">Semester</span>
        <div class="filter-buttons-row">
          <button v-for="semester in ['1st Sem', '2nd Sem']" :key="semester"
            class="filter-button" :class="{ active: activeFilters.includes(semester) }"
            @click="handleFilterClick(semester)">
            {{ semester }}
          </button>
        </div>
      </div>

      <div class="filter-category">
        <span class="category-label">Course</span>
        <div class="filter-buttons-row">
          <button v-for="course in ['BSIT', 'BSCS', 'BSIS']" :key="course"
            class="filter-button" :class="{ active: activeFilters.includes(course) }"
            @click="handleFilterClick(course)">
            {{ course }}
          </button>
        </div>
      </div>

      <div class="filter-category">
        <span class="category-label">Block</span>
        <div class="filter-buttons-row">
          <button v-for="block in ['22A', '1B', '2B']" :key="block" 
            class="filter-button" :class="{ active: activeFilters.includes(block) }"
            @click="handleFilterClick(block)">
            {{ block }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>