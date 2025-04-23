<script setup lang="ts">
import { ref, defineAsyncComponent, computed, onMounted } from 'vue'
import type { Student, StudentGrades, StudentBilling, StudentAssessmentSummary, PaginatedStudentAssessmentResponse, PaginationMetadata } from '@/typescript/models'
import studentsData from '@/mock/models.json'

const Searchbar = defineAsyncComponent(() => import("@/components/Searchbar.vue"))

const students = ref<Student[]>(studentsData.students as Student[])
const manageStudents = ref<StudentAssessmentSummary[]>([])

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

  if (statusFilters.includes(filter)) {
    if (filter === 'All Students') {
      activeFilters.value = ['All Students'];
    } else {
      if (activeFilters.value.includes(filter)) {
        activeFilters.value = ['All Students'];
      } else {
        activeFilters.value = activeFilters.value.filter(f => !statusFilters.includes(f));
        activeFilters.value.push(filter); 
      }
    }
  } else {
    if (activeFilters.value.includes(filter)) {
      activeFilters.value = activeFilters.value.filter(f => f !== filter);
    } else {
      activeFilters.value.push(filter);
    }
  }
}

// HANDLES SEARCH AND FILTER
const filteredStudents = computed(() => {
  const activeStatusFilter = activeFilters.value.find(filter =>
    ['All Students', 'Continuing', 'Withdrawn', 'Dropped', 'Probationary'].includes(filter)
  );

  return manageStudents.value.filter((student: StudentAssessmentSummary) => {
    const matchesSearch = searchQuery.value === '' ||
      student.student_id.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      (student.name?.toLowerCase() ?? '').includes(searchQuery.value.toLowerCase()) ||
      (student.course?.toLowerCase() ?? '').includes(searchQuery.value.toLowerCase());

    const matchesStatus = !activeStatusFilter ||
      activeStatusFilter === 'All Students' ||
      (student.status ?? '') === activeStatusFilter;

    return matchesSearch && matchesStatus;
  });
});

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

// Define refs for pagination state (optional, for future UI implementation)
const currentPage = ref(1);
const pageSize = ref(10);
const totalStudents = ref(0);
const totalPages = ref(0);

// Function to fetch students for a specific assessment term with pagination
const fetchStudentsForTerm = async (termId: number, page: number = 1, limit: number = 10) => {
  try {
    // Construct URL with query parameters
    const url = `http://localhost:8080/students/assessment-term/${termId}?page=${page}&limit=${limit}`;
    console.log(`Fetching from URL: ${url}`); 

    const response = await fetch(url);
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    const responseData: PaginatedStudentAssessmentResponse = await response.json();

    console.log('Fetched Paginated Response:', responseData);
    console.log('Fetched Students Data (responseData.data):', responseData.data);
    console.log('Fetched Pagination Metadata (responseData.pagination):', responseData.pagination);

    currentPage.value = responseData.pagination.currentPage;
    pageSize.value = responseData.pagination.pageSize;
    totalStudents.value = responseData.pagination.totalItems;
    totalPages.value = responseData.pagination.totalPages;

    console.log('Mock Students Data Structure (first student):', students.value.length > 0 ? students.value[0] : 'No mock data');
    console.log('Fetched Students Data Structure (first student):', responseData.data.length > 0 ? responseData.data[0] : 'No fetched data');
    manageStudents.value = responseData.data as StudentAssessmentSummary[];

  } catch (error) {
    console.error('Error fetching students for assessment term:', error);
  }
};

onMounted(() => {
  fetchStudentsForTerm(1, currentPage.value, pageSize.value);
});
</script>

<template>
  <main>
    
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
              <tr v-for="student in filteredStudents" :key="student.student_id">
                <td>{{ student.student_id }}</td>
                <td>{{ student.name }}</td>
                <td>{{ student.course }}</td>
                <td>{{ student.year_level }}</td>
                <td>
                  <span :class="['status-badge', `status-${student.status.toLowerCase()}`]">
                    {{ student.status }}
                  </span>
                </td>
                <td>
                  <button class="action-btn" @click="openModal(student.student_id)">
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
