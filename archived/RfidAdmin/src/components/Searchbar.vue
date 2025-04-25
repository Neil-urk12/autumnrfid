<script setup>
import { ref, computed } from 'vue'
import studentsData from '@/mock/models.json'

const emit = defineEmits(['update:searchQuery', 'filterChange'])

const searchQuery = ref('')
const filterContainerActive = ref(false)
const activeFilters = ref([])
const selectedCourse = ref('')
const selectedYear = ref('')

// Reference to all students
const students = studentsData.students

// Updated blockOptions to extract real blocks from student data
const blockOptions = computed(() => {
  if (!selectedCourse.value || !selectedYear.value) return []
  
  // Get year prefix: 1st Year -> 1, 2nd Year -> 2, etc.
  const yearPrefix = selectedYear.value.charAt(0)
  
  // Filter blocks that match the selected course and year
  const matchingBlocks = students
    .filter(student => 
      student.course === selectedCourse.value && 
      student.yearLevel === selectedYear.value && 
      student.block
    )
    .map(student => student.block)
  
  // Remove duplicates
  const uniqueBlocks = [...new Set(matchingBlocks)]
  
  // If no blocks found, generate some defaults
  if (uniqueBlocks.length === 0) {
    const coursePrefix = selectedCourse.value.replace('BS', '')
    const sections = ['A', 'B', 'C']
    return sections.map(section => `${coursePrefix}${yearPrefix}1${section}`)
  }
  
  return uniqueBlocks
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
      activeFilters.value = activeFilters.value.filter(f => !['BSIT', 'BSCS', 'BSA', 'BSHM'].includes(f))
      activeFilters.value.push(filter)
    }
    // Clear any selected blocks when course changes
    activeFilters.value = activeFilters.value.filter(f => !/^[A-Z]{2,}\d{1,2}[A-Z]$/.test(f))
  } 

  // YEAR SELECTION
  else if (category === 'year') {
    // Convert short year (e.g., '1st') to full year level (e.g., '1st Year')
    const fullYearLevel = filter.includes('Year') ? filter : `${filter} Year`;
    
    if (selectedYear.value === filter) {
      selectedYear.value = ''
      activeFilters.value = activeFilters.value.filter(f => f !== filter)
    } else {
      selectedYear.value = filter
      activeFilters.value = activeFilters.value.filter(f => !['1st Year', '2nd Year', '3rd Year', '4th Year', '1st', '2nd', '3rd', '4th'].includes(f))
      activeFilters.value.push(filter)
    }
    // Clear any selected blocks when year changes
    activeFilters.value = activeFilters.value.filter(f => !/^[A-Z]{2,}\d{1,2}[A-Z]$/.test(f))
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
      <span v-if="activeFilters.length > 0" class="filter-count">{{ activeFilters.length }}</span>
    </button>

    <div class="filter-tooltip" :class="{ active: filterContainerActive }">
      <div class="filter-category">
        <span class="category-label">Year Level</span>
        <div class="filter-buttons-row">
          <button v-for="year in ['1st', '2nd', '3rd', '4th']" :key="year"
            class="filter-button" 
            :class="{ active: activeFilters.includes(year) || activeFilters.includes(year.replace(' Year', '')) }"
            @click="handleFilterClick(year, 'year')">
            {{ year }}
          </button>
        </div>
      </div>

      <div class="filter-category">
        <span class="category-label">Course</span>
        <div class="filter-buttons-row">
          <button v-for="course in ['BSIT', 'BSCS', 'BSA', 'BSHM']" :key="course"
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
        <div v-else-if="blockOptions.length === 0" class="block-message">
          No blocks available for this selection
        </div>
      </div>
      
      <div v-if="activeFilters.length > 0" class="active-filters">
        <span class="category-label">Active Filters:</span>
        <div class="active-filters-row">
          <span v-for="filter in activeFilters" :key="filter" class="active-filter-tag">
            {{ filter }}
            <button class="remove-filter" @click="handleFilterClick(filter, filter.includes('Year') || /^[1-4](st|nd|rd|th)$/.test(filter) ? 'year' : 
                                                                   /^[A-Z]{2,}\d{1,2}[A-Z]$/.test(filter) ? 'block' : 'course')">
              &times;
            </button>
          </span>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.search-bar {
  position: relative;
  display: flex;
  align-items: center;
  max-width: 500px;
  width: 100%;
}

.search-bar input {
  flex: 1;
  padding: 10px 10px 10px 35px;
  border-radius: 4px;
  font-size: 14px;
}

.search-bar i {
  position: absolute;
  left: 10px;
  color: #999;
}

.filter-toggle-btn {
  display: flex;
  align-items: center;
  margin-left: 10px;
  padding: 8px 12px;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.3s;
}



.filter-toggle-btn i {
  position: static;
  margin-right: 5px;
}

.filter-tooltip {
  position: absolute;
  top: 100%;
  right: 0;
  width: 300px;
  border-radius: 4px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  padding: 15px;
  margin-top: 10px;
  z-index: 100;
  display: none;
}

.filter-tooltip.active {
  display: block;
}


.category-label {
  display: block;

  color: rgb(179, 176, 176);
}

.filter-buttons-row {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.filter-button {
  padding: 6px 12px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 13px;
  transition: all 0.2s;
}


.filter-button.active {
  background-color: var(--accent);
  color: white;
  border-color: var(--accent);
}

.filter-button.disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.block-message {
  font-size: 13px;
  color: #999;
  font-style: italic;
}

.filter-count {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 20px;
  height: 20px;
  background-color: var(--accent);
  color: white;
  border-radius: 50%;
  font-size: 12px;
  margin-left: 5px;
}

.active-filters {
  padding-top: 15px;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}

.active-filters-row {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 8px;
}

.active-filter-tag {
  display: inline-flex;
  align-items: center;
  padding: 4px 8px 4px 12px;
  border-radius: 16px;
  font-size: 12px;
  background: rgba(255, 255, 255, 0.05);
    border: 1px solid rgba(255, 255, 255, 0.1);
    font-weight: normal;
  color: white;
}

.remove-filter {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 16px;
  height: 16px;
  background-color: rgba(37, 37, 37, 0.2);
  border-radius: 50%;
  margin-left: 6px;
  font-size: 12px;
  border: none;
  cursor: pointer;
  padding: 0;
  color: white;
}

.remove-filter:hover {
  background-color: rgba(0, 0, 0, 0.2);
}
</style>