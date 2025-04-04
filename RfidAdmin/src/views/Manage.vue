<script setup>
import { ref } from 'vue'
import Sidebar from "@/components/Sidebar.vue";

const activeCategory = ref('information')
const isModalOpen = ref(false)
const selectedStudentId = ref(null)
const filterContainerActive = ref(false)
const activeFilters = ref([])
const searchQuery = ref('')

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

const toggleFilterContainer = () => {
  filterContainerActive.value = !filterContainerActive.value
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
</script>

<template>
      <div class="sidebar">
        <Sidebar />
    </div>

  <section>
    <div style="width: 100%; max-width: 1200px; margin: 0 auto; display: flex; flex-direction: column; gap: 20px;">
      <div class="welcome-header">
        <h1>Management</h1>
        <p>Manage and monitor student information</p>
      </div>

      <div class="students-controls">
        <div class="search-filters">
          <div class="search-bar">
            <i class="fa-solid fa-search"></i>
            <input type="text" placeholder="Search students..." v-model="searchQuery">
            <button class="filter-toggle-btn" @click="toggleFilterContainer">
              <i class="fa-solid fa-filter"></i>
              Filter
            </button>

            <div class="filter-tooltip" :class="{ active: filterContainerActive }" id="filterContainer">
              <div class="filter-category">
                <span class="category-label">Year Level</span>
                <div class="filter-buttons-row">
                  <button v-for="year in ['1st', '2nd', '3rd', '4th']" :key="year" class="filter-btn"
                    :class="{ active: activeFilters.includes(year) }" @click="handleFilterClick(year)">
                    {{ year }}
                  </button>
                </div>
              </div>

              <div class="filter-category">
                <span class="category-label">Semester</span>
                <div class="filter-buttons-row">
                  <button v-for="semester in ['1st Sem', '2nd Sem']" :key="semester" class="filter-btn"
                    :class="{ active: activeFilters.includes(semester) }" @click="handleFilterClick(semester)">
                    {{ semester }}
                  </button>
                </div>
              </div>

              <div class="filter-category">
                <span class="category-label">Course</span>
                <div class="filter-buttons-row">
                  <button v-for="course in ['BSIT', 'BSCS', 'BSIS']" :key="course" class="filter-btn"
                    :class="{ active: activeFilters.includes(course) }" @click="handleFilterClick(course)">
                    {{ course }}
                  </button>
                </div>
              </div>

              <div class="filter-category">
                <span class="category-label">Block</span>
                <div class="filter-buttons-row">
                  <button v-for="block in ['22A', '1B', '2B']" :key="block" class="filter-btn"
                    :class="{ active: activeFilters.includes(block) }" @click="handleFilterClick(block)">
                    {{ block }}
                  </button>
                </div>
              </div>
            </div>
          </div>

          <div class="filter-buttons">
            <button v-for="status in ['All Students', 'Continuing', 'Withdrawn', 'Dropped', 'Probationary']"
              :key="status" class="filter-btn" :class="{ active: activeFilters.includes(status) }"
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
            <tr v-for="student in students" :key="student.id">
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
              <div v-for="(value, key) in students.find(s => s.id === selectedStudentId)" :key="key" class="info-group">
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
</template>