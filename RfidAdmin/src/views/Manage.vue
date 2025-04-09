<script setup>
import { ref, defineAsyncComponent, computed } from 'vue'
const Sidebar = defineAsyncComponent(() => import("@/components/Sidebar.vue"))
const Searchbar = defineAsyncComponent(() => import("@/components/Searchbar.vue"))

const activeCategory = ref('information')
const isModalOpen = ref(false)
const selectedStudentId = ref(null)

const searchQuery = ref('')
const activeFilters = ref([])

const switchCategory = (category) => {
  activeCategory.value = category
}

const openModal = (studentId) => {
  selectedStudentId.value = studentId
  isModalOpen.value = true
}

const closeModal = () => {
  isModalOpen.value = false
  selectedStudentId.value = null
}

const handleSearch = (query) => {
  searchQuery.value = query
}

const handleFilterChange = (filters) => {
  activeFilters.value = filters
}

const handleFilterClick = (filter) => {
  if (activeFilters.value.includes(filter)) {
    activeFilters.value = activeFilters.value.filter(f => f !== filter)
  } else {
    activeFilters.value.push(filter)
  }
}

const students = ref([
  {
    id: '2023-0001',
    name: 'John Cez',
    course: 'Computer Science',
    yearLevel: '3rd Year',
    status: 'Continuing',
    email: 'john.cez@example.com',
    phone: '+63 29483928',
    address: 'Compostella'
  },
  {
    id: '2023-0002',
    name: 'Emily Casupana',
    course: 'Information Technology',
    yearLevel: '2nd Year',
    status: 'Probationary',
    email: 'emily.casupana@example.com',
    phone: '+63 12345678',
    address: 'Davao City'
  },
  {
    id: '2023-0003',
    name: 'Jan Rosa',
    course: 'Computer Engineering',
    yearLevel: '1st Year',
    status: 'Dropped',
    email: 'jan.rosa@example.com',
    phone: '+63 87654321',
    address: 'Tagum City'
  }
])

const filteredStudents = computed(() => {
  return students.value.filter(student => {
    const matchesSearch = searchQuery.value === '' ||
      student.id.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      student.name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      student.course.toLowerCase().includes(searchQuery.value.toLowerCase())

    const matchesFilters = activeFilters.value.length === 0 ||
      activeFilters.value.some(filter =>
        student.yearLevel.includes(filter) ||
        student.course.includes(filter)
      )

    return matchesSearch && matchesFilters
  })
})
</script>

<template>
  <main>
    <div class="sidebar">
      <Sidebar />
    </div>

    <section>
      <div class="container">
        <div class="welcome-header">
          <h1>Management</h1>
          <p>Manage and monitor student information</p>
        </div>

        <div class="students-controls">
          <div class="search-filters">
     
              <Searchbar 
                v-model="searchQuery"
                @update:search-query="handleSearch"
                @filter-change="handleFilterChange"
              />
          

            <div class="filter-buttons">
              <button v-for="status in ['All Students', 'Continuing', 'Withdrawn', 'Dropped', 'Probationary']"
                :key="status" 
                class="filter-btn" 
                :class="{ active: activeFilters.includes(status) }"
                @click="handleFilterClick(status)">
                {{ status }}
              </button>
            </div>
          </div>
        </div>

        <div class="students-table">
          <table>
            <thead>
              <tr>
                <th>Student ID</th>
                <th>Name</th>
                <th>Course</th>
                <th>Year Level</th>
                <th>Status</th>
                <th>Action</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="student in filteredStudents" :key="student.id">
                <td>{{ student.id }}</td>
                <td>{{ student.name }}</td>
                <td>{{ student.course }}</td>
                <td>{{ student.yearLevel }}</td>
                <td>
                  <span :class="['status-badge', `status-${student.status.toLowerCase()}`]">
                    {{ student.status }}
                  </span>
                </td>
                <td>
                  <button class="action-btn" @click="openModal(student.id)">
                    View
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <div class="modal" :class="{ active: isModalOpen }" @click="closeModal">
        <div class="modal-content" @click.stop>
          <div class="modal-header">
            <h2>Student Information</h2>
            <button class="close-modal" @click="closeModal">&times;</button>
          </div>

          <div class="category-buttons">
            <button v-for="category in ['information', 'grades', 'bills']" :key="category" class="category-btn"
              :class="{ active: activeCategory === category }" @click="switchCategory(category)">
              {{ category.charAt(0).toUpperCase() + category.slice(1) }}
            </button>
          </div>

          <div class="category-content" :class="{ active: activeCategory === 'information' }" id="information-content">
            <div class="student-info-grid">
              <template v-if="selectedStudentId">
                <div v-for="(value, key) in students.find(s => s.id === selectedStudentId)" :key="key"
                  class="info-group">
                  <label>{{ key.charAt(0).toUpperCase() + key.slice(1) }}</label>
                  <span>{{ value }}</span>
                </div>
              </template>
            </div>
          </div>

          <div class="category-content" :class="{ active: activeCategory === 'grades' }" id="grades-content">
            <div class="grades-summary">
              <table class="grades-table">
                <thead>
                  <tr>
                    <th>Subject</th>
                    <th>Prelim</th>
                    <th>Midterm</th>
                    <th>Prefinals</th>
                    <th>Finals</th>
                    <th>GWA</th>
                  </tr>
                </thead>
                <tbody>
                  <tr>
                    <td>Computer Fundamentals</td>
                    <td>1</td>
                    <td>1.25</td>
                    <td>1</td>
                    <td>1</td>
                    <td>1</td>
                  </tr>
                  <tr>
                    <td>Information Management</td>
                    <td>1</td>
                    <td>1</td>
                    <td>1</td>
                    <td>1</td>
                    <td>1</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>

          <div class="category-content" :class="{ active: activeCategory === 'bills' }" id="bills-content">
            <div class="bills-summary">
              <div class="total-balance">
                <h3>Total Balance</h3>
                <span class="amount">₱25,000.00</span>
              </div>
              <table class="bills-table">
                <thead>
                  <tr>
                    <th>Description</th>
                    <th>Date</th>
                    <th>Amount</th>
                  </tr>
                </thead>
                <tbody>
                  <tr>
                    <td>Tuition Fee</td>
                    <td>2024-01-15</td>
                    <td>₱22,500.00</td>
                  </tr>
                  <tr>
                    <td>Prelim Exam</td>
                    <td>2024-01-15</td>
                    <td>₱4,880.00</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </div>
    </section>
  </main>
</template>
