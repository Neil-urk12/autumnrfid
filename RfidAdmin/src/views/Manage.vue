<script setup>
import { ref, defineAsyncComponent, computed } from 'vue'
const Sidebar = defineAsyncComponent(() => import("@/components/Sidebar.vue"))
const Searchbar = defineAsyncComponent(() => import("@/components/Searchbar.vue"))

// MOCK
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


const activeCategory = ref('information')
const isModalOpen = ref(false)
const selectedStudentId = ref(null)
const searchQuery = ref('')
const activeFilters = ref([])

// SWITCHES THE ACTIVE CATEGORY TAB IN THE MODAL
const switchCategory = (category) => {
  activeCategory.value = category
}

// OPENS THE STUDENT DETAILS MODAL AND SETS THE SELECTED STUDENT ID
const openModal = (studentId) => {
  selectedStudentId.value = studentId
  isModalOpen.value = true
}

// CLOSES THE MODAL AND RESETS THE SELECTED STUDENT ID
const closeModal = () => {
  isModalOpen.value = false
  selectedStudentId.value = null
}

// UPDATES THE SEARCH QUERY VALUE
const handleSearch = (query) => {
  searchQuery.value = query
}

// UPDATES THE ACTIVE FILTERS ARRAY
const handleFilterChange = (filters) => {
  activeFilters.value = filters
}

// HANDLES THE FILTER BUTTON CLICKS AND UPDATES THE ACTIVE FILTERS
const handleFilterClick = (filter) => {
  if (activeFilters.value.includes(filter)) {
    activeFilters.value = []
  } else {
    activeFilters.value = activeFilters.value.filter(f => 
      !['All Students', 'Continuing', 'Withdrawn', 'Dropped', 'Probationary'].includes(f)
    )
    activeFilters.value.push(filter)
  }
}


// HANDLES SEARCH AND FILTER
const filteredStudents = computed(() => {
  return students.value.filter(student => {
    const matchesSearch = searchQuery.value === '' ||
      student.id.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      student.name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      student.course.toLowerCase().includes(searchQuery.value.toLowerCase())

    const statusFilter = activeFilters.value.find(filter => 
      ['All Students', 'Continuing', 'Withdrawn', 'Dropped', 'Probationary'].includes(filter)
    )
    
    const matchesStatus = !statusFilter || 
      statusFilter === 'All Students' || 
      student.status === statusFilter

    const otherFilters = activeFilters.value.filter(filter => 
      !['All Students', 'Continuing', 'Withdrawn', 'Dropped', 'Probationary'].includes(filter)
    )
    
    const matchesOtherFilters = otherFilters.length === 0 ||
      otherFilters.some(filter =>
        student.yearLevel.includes(filter) ||
        student.course.includes(filter)
      )

    return matchesSearch && matchesStatus && matchesOtherFilters
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

        <!-- SEARCH BAR SECTION -->
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
    
        <!-- STUDENT'S TABLE -->
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

      <!-- VIEW STUDENT DETAILS MODAL -->
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

          <!-- STUDENT DETAILS -->
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

          <!-- STUDENT GRADES -->
          <div class="category-content" :class="{ active: activeCategory === 'grades' }" id="grades-content">
            <div class="grades-summary">
              <table class="grades-table">
                <thead>
                  <tr>
                    <th class="subject-col">Subject</th>
                    <th>Prelim</th>
                    <th>Midterm</th>
                    <th>Prefinals</th>
                    <th>Finals</th>
                    <th>GWA</th>
                  </tr>
                </thead>
                <tbody>
                  <tr>
                    <td class="subject-col">Computer Fundamentals</td>
                    <td class="grade-col">1.00</td>
                    <td class="grade-col">1.25</td>
                    <td class="grade-col">1.00</td>
                    <td class="grade-col">1.00</td>
                    <td class="grade-col">1.00</td>
                  </tr>
                  <tr>
                    <td class="subject-col">Information Management</td>
                    <td class="grade-col">1.00</td>
                    <td class="grade-col">1.00</td>
                    <td class="grade-col">1.00</td>
                    <td class="grade-col">1.00</td>
                    <td class="grade-col">1.00</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>

          <!-- STUDENT BILLS -->
          <div class="category-content" :class="{ active: activeCategory === 'bills' }" id="bills-content">
            <div class="bills-summary">
              <div class="balance-grid">
                <div class="balance-item">
                  <h3>Total Tuition</h3>
                  <span class="amount">₱27,380.00</span>
                </div>
                <div class="balance-item">
                  <h3>Total Paid</h3>
                  <span class="amount paid">₱22,500.00</span>
                </div>
                <div class="balance-item">
                  <h3>Remaining Balance</h3>
                  <span class="amount remaining">₱4,880.00</span>
                </div>
              </div>
              <div class="bills-table-container">
                <table class="bills-table">
                  <thead>
                    <tr>
                      <th class="desc-col">Description</th>
                      <th class="date-col">Date</th>
                      <th class="amount-col">Amount</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr>
                      <td class="desc-col">Tuition Fee</td>
                      <td class="date-col">2024-01-15</td>
                      <td class="amount-col">₱22,500.00</td>
                    </tr>
                    <tr>
                      <td class="desc-col">Prelim Exam</td>
                      <td class="date-col">2024-01-15</td>
                      <td class="amount-col">₱4,880.00</td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  </main>
</template>
