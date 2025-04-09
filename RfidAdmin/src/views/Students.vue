<script setup>
import { ref, defineAsyncComponent, computed } from 'vue'
const ConfirmationModal = defineAsyncComponent(() => import('@/components/ConfirmationModal.vue'))
const Sidebar = defineAsyncComponent(() => import("@/components/Sidebar.vue"))
const Searchbar = defineAsyncComponent(() => import("@/components/Searchbar.vue"))

const isAddModalOpen = ref(false)
const isEditModalOpen = ref(false)
const selectedStudentId = ref(null)
const filterContainerActive = ref(false)
const searchQuery = ref('')
const isConfirmationModalOpen = ref(false)
const activeFilters = ref([])

const students = ref([
  {
    id: '2023-0001',
    name: 'John Cez',
    course: 'Computer Science',
    yearLevel: '3rd Year',
    status: 'Active',
    email: 'john.cez@example.com',
    phone: '+63 123456789',
    address: 'Compostella'
  }
])

const formData = ref({
  id: '',
  name: '',
  course: '',
  yearLevel: '',
  status: '',
  email: '',
  phone: '',
  address: ''
})

const courses = ['BSIT', 'BSCS', 'BSIS', 'BSHM']
const yearLevels = ['1st Year', '2nd Year', '3rd Year', '4th Year']
const statusOptions = ['Active', 'Continuing', 'Dropped', 'Probationary']

const toggleFilterContainer = () => {
  filterContainerActive.value = !filterContainerActive.value
}

const openAddModal = () => {
  formData.value = {
    id: '',
    name: '',
    course: '',
    yearLevel: '',
    status: '',
    email: '',
    phone: '',
    address: ''
  }
  isAddModalOpen.value = true
}

const closeAddModal = () => {
  isAddModalOpen.value = false
}

const openEditModal = (studentId) => {
  const student = students.value.find(s => s.id === studentId)
  if (student) {
    formData.value = { ...student }
    selectedStudentId.value = studentId
    isEditModalOpen.value = true
  }
}

const closeEditModal = () => {
  isEditModalOpen.value = false
  selectedStudentId.value = null
}

const handleSubmitAdd = (e) => {
  e.preventDefault()
  students.value.push({ ...formData.value })
  closeAddModal()
}

const handleSubmitEdit = (e) => {
  e.preventDefault()
  const index = students.value.findIndex(s => s.id === selectedStudentId.value)
  if (index !== -1) {
    students.value[index] = { ...formData.value }
  }
  closeEditModal()
}

const confirmationData = ref({
  title: '',
  itemName: '',
  itemInfo: null
})

const showDeleteConfirmation = (title, itemName, itemInfo = null) => {
  confirmationData.value = { title, itemName, itemInfo }
  isConfirmationModalOpen.value = true
}

const handleConfirmDelete = (itemInfo) => {
  console.log('Deleting:', itemInfo)
}

const handleSearch = (query) => {
  searchQuery.value = query || ''
}

const handleFilterChange = (filters) => {
  activeFilters.value = filters || []
}

const filteredStudents = computed(() => {
  return students.value.filter(student => {
    const matchesSearch = !searchQuery.value || 
      student.id.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      student.name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      student.course.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      student.yearLevel.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      student.status.toLowerCase().includes(searchQuery.value.toLowerCase())

    const matchesFilters = !activeFilters.value.length ||
      activeFilters.value.some(filter => 
        student.yearLevel.includes(filter) ||
        student.course.includes(filter) ||
        student.status.includes(filter))

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
          <h1>Manage Students</h1>
          <p>Add, edit, or remove student records</p>
        </div>


        <div class="students-controls">
          <div class="search-filters">
              <Searchbar 
                :v-model="searchQuery"
                @update:search-query="handleSearch"
                @filter-change="handleFilterChange"
              />
            <div class="filter-buttons">
              <button class="add-student-btn" @click="openAddModal">
                <i class="fa-solid fa-plus"></i> Add New Student
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
                <th>Actions</th>
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
                <td class="action-buttons">
                  <button class="action-btn edit-btn" @click="openEditModal(student.id)">
                    <i class="fa-solid fa-pen-to-square"></i>
                  </button>
                  <button class="action-btn delete-btn" @click="showDeleteConfirmation(
                    'Delete Student',
                    `student with ID ${student.id}`,
                    student
                  )">
                    <i class="fa-solid fa-trash"></i>
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>


      <!-- ADD STUDENT MODAL -->
      <div class="modal" :class="{ active: isAddModalOpen }" @click="closeAddModal">
        <div class="modal-content" @click.stop>
          <div class="modal-header">
            <h2>Add New Student</h2>
            <button class="close-modal" @click="closeAddModal">&times;</button>
          </div>
          <form @submit="handleSubmitAdd" class="student-form">
            <div class="student-info-grid">
              <div class="info-group">
                <label>Student ID</label>
                <input type="text" v-model="formData.id" required>
              </div>
              <div class="info-group">
                <label>Name</label>
                <input type="text" v-model="formData.name" required>
              </div>
              <div class="info-group">
                <label>Course</label>
                <select v-model="formData.course" required>
                  <option v-for="course in courses" :key="course" :value="course">
                    {{ course }}
                  </option>
                </select>
              </div>
              <div class="info-group">
                <label>Year Level</label>
                <select v-model="formData.yearLevel" required>
                  <option v-for="year in yearLevels" :key="year" :value="year">
                    {{ year }}
                  </option>
                </select>
              </div>
              <div class="info-group">
                <label>Status</label>
                <select v-model="formData.status" required>
                  <option v-for="status in statusOptions" :key="status" :value="status">
                    {{ status }}
                  </option>
                </select>
              </div>
              <div class="info-group">
                <label>Email</label>
                <input type="email" v-model="formData.email" required>
              </div>
              <div class="info-group">
                <label>Phone</label>
                <input type="tel" v-model="formData.phone" required>
              </div>
              <div class="info-group">
                <label>Address</label>
                <input type="text" v-model="formData.address" required>
              </div>
            </div>
            <div class="form-actions">
              <button type="button" class="cancel-btn" @click="closeAddModal">Cancel</button>
              <button type="submit" class="submit-btn">Add Student</button>
            </div>
          </form>
        </div>
      </div>


      <!-- EDIT STUDENT MODAL -->
      <div class="modal" :class="{ active: isEditModalOpen }" @click="closeEditModal">
        <div class="modal-content" @click.stop>
          <div class="modal-header">
            <h2>Edit Student Information</h2>
            <button class="close-modal" @click="closeEditModal">&times;</button>
          </div>
          <form @submit="handleSubmitEdit" class="student-form">
            <div class="student-info-grid">
              <div class="info-group">
                <label>Student ID</label>
                <input type="text" v-model="formData.id" required>
              </div>
              <div class="info-group">
                <label>Name</label>
                <input type="text" v-model="formData.name" required>
              </div>
              <div class="info-group">
                <label>Course</label>
                <select v-model="formData.course" required>
                  <option v-for="course in courses" :key="course" :value="course">
                    {{ course }}
                  </option>
                </select>
              </div>
              <div class="info-group">
                <label>Year Level</label>
                <select v-model="formData.yearLevel" required>
                  <option v-for="year in yearLevels" :key="year" :value="year">
                    {{ year }}
                  </option>
                </select>
              </div>
              <div class="info-group">
                <label>Status</label>
                <select v-model="formData.status" required>
                  <option v-for="status in statusOptions" :key="status" :value="status">
                    {{ status }}
                  </option>
                </select>
              </div>
              <div class="info-group">
                <label>Email</label>
                <input type="email" v-model="formData.email" required>
              </div>
              <div class="info-group">
                <label>Phone</label>
                <input type="tel" v-model="formData.phone" required>
              </div>
              <div class="info-group">
                <label>Address</label>
                <input type="text" v-model="formData.address" required>
              </div>
            </div>
            <div class="form-actions">
              <button type="button" class="cancel-btn" @click="closeEditModal">Cancel</button>
              <button type="submit" class="submit-btn">Save Changes</button>
            </div>
          </form>
        </div>
      </div>



      <!-- DELETE STUDENT MODAL -->
      <ConfirmationModal :is-open="isConfirmationModalOpen" :title="confirmationData.title"
        :item-name="confirmationData.itemName" :item-info="confirmationData.itemInfo"
        @close="isConfirmationModalOpen = false" @confirm="handleConfirmDelete" />
    </section>
  </main>
</template>
