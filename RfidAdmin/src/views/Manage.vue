<script setup lang="ts">
import { ref, defineAsyncComponent, computed } from 'vue'
import type { Student, StudentGrades, StudentBilling } from '@/typescript/models'
import studentsData from '@/mock/models.json'

const Sidebar = defineAsyncComponent(() => import("@/components/Sidebar.vue"))
const Searchbar = defineAsyncComponent(() => import("@/components/Searchbar.vue"))

const students = ref<Student[]>(studentsData.students as Student[])

const activeCategory = ref<'information' | 'grades' | 'bills'>('information')
const isModalOpen = ref<boolean>(false)
const selectedStudentId = ref<string | null>(null)
const searchQuery = ref<string>('')
const activeFilters = ref<string[]>([])

// SWITCHES THE ACTIVE CATEGORY TAB IN THE MODAL
const switchCategory = (category: 'information' | 'grades' | 'bills') => {
  activeCategory.value = category
}

// OPENS THE STUDENT DETAILS MODAL AND SETS THE SELECTED STUDENT ID
const openModal = (studentId: string) => {
  selectedStudentId.value = studentId
  isModalOpen.value = true
}

// CLOSES THE MODAL AND RESETS THE SELECTED STUDENT ID
const closeModal = () => {
  isModalOpen.value = false
  selectedStudentId.value = null
}

// UPDATES THE SEARCH QUERY VALUE
const handleSearch = (query: string) => {
  searchQuery.value = query
}

// UPDATES THE ACTIVE FILTERS ARRAY
const handleFilterChange = (filters: string[]) => {
  activeFilters.value = filters
}

// HANDLES THE FILTER BUTTON CLICKS AND UPDATES THE ACTIVE FILTERS
const handleFilterClick = (filter: string) => {
  const statusFilters = ['All Students', 'Continuing', 'Withdrawn', 'Dropped', 'Probationary'];
  
  if (activeFilters.value.includes(filter)) {
    activeFilters.value = []
  } else {
    if (statusFilters.includes(filter)) {
      activeFilters.value = activeFilters.value.filter(f => !statusFilters.includes(f))
    }
    activeFilters.value.push(filter)
  }
}

// HANDLES SEARCH AND FILTER
const filteredStudents = computed(() => {
  return students.value.filter((student: Student) => {
    const matchesSearch = searchQuery.value === '' ||
      student.id.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      `${student.firstName} ${student.lastName}`.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
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
    
    const blockFilters = otherFilters.filter(filter => 
      /^[A-Z]{2,}\d{1,2}[A-Z]$/.test(filter)
    )
    
    const courseYearFilters = otherFilters.filter(filter => 
      !blockFilters.includes(filter)
    )
    
    const matchesBlock = blockFilters.length === 0 || 
      blockFilters.some(filter => student.block === filter)
    
    const matchesCourseYear = courseYearFilters.length === 0 ||
      courseYearFilters.some(filter => {
        if (/^[1-4](st|nd|rd|th)(\sYear)?$/.test(filter)) {
          const yearFilter = filter.includes('Year') ? filter : `${filter} Year`;
          return student.yearLevel === yearFilter;
        }
        return student.course.includes(filter);
      })

    return matchesSearch && matchesStatus && matchesBlock && matchesCourseYear
  })
})
//

const getStudentInfo = (studentId: string | null) => {
  if (!studentId) return null
  const student = students.value.find((s: Student) => s.id === studentId)
  if (!student) return null
  
  return {
    id: student.id,
    firstName: student.firstName,
    lastName: student.lastName,
    middleName: student.middleName,
    suffix: student.suffix,
    birthday: student.birthday,
    course: student.course,
    block: student.block,
    yearLevel: student.yearLevel,
    status: student.status,
    email: student.email,
    phone: student.phone
  }
}

const getStudentGrades = (studentId: string | null): StudentGrades | null => {
  if (!studentId) return null
  const student = students.value.find((s: Student) => s.id === studentId)
  return student?.grades || null
}

const subjects = ref(studentsData.subjects)
const getStudentBilling = (studentId: string | null): StudentBilling | null => {
  if (!studentId) return null
  const student = students.value.find((s: Student) => s.id === studentId)
  return student?.billing || null
}

const categories = ['information', 'grades', 'bills'] as const
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
                <td>{{ student.firstName }} {{ student.lastName }}</td>
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
            <button v-for="category in categories" 
              :key="category" 
              class="category-btn"
              :class="{ active: activeCategory === category }" 
              @click="switchCategory(category)">
              {{ category.charAt(0).toUpperCase() + category.slice(1) }}
            </button>
          </div>

          <!-- STUDENT DETAILS -->
          <div class="category-content" :class="{ active: activeCategory === 'information' }" id="information-content">
            <div class="scrollable-content">
              <div class="student-info-grid">
                <template v-if="selectedStudentId">
                  <div v-for="(value, key) in getStudentInfo(selectedStudentId)" :key="key"
                    class="info-group">
                    <label>{{ String(key).charAt(0).toUpperCase() + String(key).slice(1) }}</label>
                    <span>{{ value }}</span>
                  </div>
                </template>
              </div>
            </div>
          </div>

          <!-- STUDENT GRADES -->
          <div class="category-content" :class="{ active: activeCategory === 'grades' }" id="grades-content">
            <div class="scrollable-content">
              <div class="grades-summary">
                <template v-if="getStudentGrades(selectedStudentId)">
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
                      <tr v-for="subject in subjects" :key="subject.code">
                        <td class="subject-col">{{ subject.name }}</td>
                        <td class="grade-col">{{ getStudentGrades(selectedStudentId)?.prelim[subject.code] || '-' }}</td>
                        <td class="grade-col">{{ getStudentGrades(selectedStudentId)?.midterm[subject.code] || '-' }}</td>
                        <td class="grade-col">{{ getStudentGrades(selectedStudentId)?.prefinals[subject.code] || '-' }}</td>
                        <td class="grade-col">{{ getStudentGrades(selectedStudentId)?.finals[subject.code] || '-' }}</td>
                        <td class="grade-col">{{ getStudentGrades(selectedStudentId)?.gwa || '-' }}</td>
                      </tr>
                    </tbody>
                  </table>
                  <div class="remarks-container" v-if="getStudentGrades(selectedStudentId)?.remarks">
                    <span :class="['status-badge', `status-${getStudentGrades(selectedStudentId)?.remarks.toLowerCase()}`]">
                      Remarks: {{ getStudentGrades(selectedStudentId)?.remarks }}
                    </span>
                  </div>
                </template>
                <p v-else>No grades available for this student.</p>
              </div>
            </div>
          </div>

          <!-- STUDENT BILLS -->
          <div class="category-content" :class="{ active: activeCategory === 'bills' }" id="bills-content">
            <div class="scrollable-content">
              <div class="bills-summary">
                <template v-if="getStudentBilling(selectedStudentId)">
                  <div class="balance-grid">
                    <div class="balance-item">
                      <h3>Total Tuition</h3>
                      <span class="amount">₱{{ getStudentBilling(selectedStudentId)?.totalTuition.toFixed(2) }}</span>
                    </div>
                    <div class="balance-item">
                      <h3>Total Paid</h3>
                      <span class="amount paid">₱{{ getStudentBilling(selectedStudentId)?.totalPaid.toFixed(2) }}</span>
                    </div>
                    <div class="balance-item">
                      <h3>Remaining Balance</h3>
                      <span class="amount remaining">₱{{ getStudentBilling(selectedStudentId)?.remainingBalance.toFixed(2) }}</span>
                    </div>
                  </div>

                  <h3 class="section-title">Payment History</h3>
                  <div class="bills-table-container">
                    <table class="bills-table">
                      <thead>
                        <tr>
                          <th class="desctit">Description</th>
                          <th class="desctit">Date</th>
                          <th class="desctit">Amount</th>
                        </tr>
                      </thead>
                      <tbody>
                        <tr class="paydat" v-for="(payment, index) in getStudentBilling(selectedStudentId)?.payments" :key="index">
                          <td class="desc-col">{{ payment.description }}</td>
                          <td class="desc-col">{{ payment.date }}</td>
                          <td class="desc-col">₱{{ payment.amount.toFixed(2) }}</td>
                        </tr>
                        <tr v-if="getStudentBilling(selectedStudentId)?.payments.length === 0">
                          <td colspan="3" class="no-data">No payment records found</td>
                        </tr>
                      </tbody>
                    </table>
                  </div>
                </template>
                <p v-else>No billing information available for this student.</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  </main>
</template>