<script setup>
import { ref, computed } from 'vue'

const emit = defineEmits(['update:searchQuery', 'filterChange'])

const searchQuery = ref('')
const filterContainerActive = ref(false)
const activeFilters = ref([])
const selectedCourse = ref('')
const selectedYear = ref('')

const blockOptions = computed(() => {
  if (!selectedCourse.value || !selectedYear.value) return []
  
  const year = selectedYear.value.charAt(0) 
  const coursePrefix = selectedCourse.value.replace('BS', '') 
  const sections = ['A', 'B', 'C']
  
  return sections.map(section => `${coursePrefix}${year}1${section}`)
})

const toggleFilterContainer = () => {
  filterContainerActive.value = !filterContainerActive.value
}

const handleFilterClick = (filter, category) => {
  if (category === 'course') {
    if (selectedCourse.value === filter) {
      selectedCourse.value = ''
      activeFilters.value = activeFilters.value.filter(f => f !== filter)
    } else {
      selectedCourse.value = filter
      activeFilters.value = activeFilters.value.filter(f => !['BSIT', 'BSCS', 'BSIS', 'BSHM'].includes(f))
      activeFilters.value.push(filter)
    }
    activeFilters.value = activeFilters.value.filter(f => !f.includes('1A') && !f.includes('1B') && !f.includes('1C'))
  } 

  // YEAR SELECTION
  else if (category === 'year') {
    if (selectedYear.value === filter) {
      selectedYear.value = ''
      activeFilters.value = activeFilters.value.filter(f => f !== filter)
    } else {
      selectedYear.value = filter
      activeFilters.value = activeFilters.value.filter(f => !['1st', '2nd', '3rd', '4th'].includes(f))
      activeFilters.value.push(filter)
    }
    activeFilters.value = activeFilters.value.filter(f => !f.includes('1A') && !f.includes('1B') && !f.includes('1C'))
  }
  
  // BLOCK SELECTION
  else if (category === 'block') {
    if (activeFilters.value.includes(filter)) {
      activeFilters.value = activeFilters.value.filter(f => f !== filter)
    } else {
      activeFilters.value = activeFilters.value.filter(f => !blockOptions.value.includes(f))
      activeFilters.value.push(filter)
    }
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
            class="filter-button" 
            :class="{ active: activeFilters.includes(year) }"
            @click="handleFilterClick(year, 'year')">
            {{ year }}
          </button>
        </div>
      </div>

      <div class="filter-category">
        <span class="category-label">Course</span>
        <div class="filter-buttons-row">
          <button v-for="course in ['BSIT', 'BSCS', 'BSIS', 'BSHM']" :key="course"
            class="filter-button" 
            :class="{ active: activeFilters.includes(course) }"
            @click="handleFilterClick(course, 'course')">
            {{ course }}
          </button>
        </div>
      </div>

      <div class="filter-category">
        <span class="category-label">Block</span>
        <div class="filter-buttons-row">
          <button v-for="block in blockOptions" :key="block"
            class="filter-button" 
            :class="{ 
              active: activeFilters.includes(block),
              disabled: !selectedCourse || !selectedYear 
            }"
            @click="handleFilterClick(block, 'block')"
            :disabled="!selectedCourse || !selectedYear">
            {{ block }}
          </button>
        </div>
        <div v-if="!selectedCourse || !selectedYear" class="block-message">
          Select a course and year level first
        </div>
      </div>
    </div>
  </div>
</template>