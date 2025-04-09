<script setup>
import { ref, computed } from 'vue'
import { defineAsyncComponent } from "vue";
const Sidebar = defineAsyncComponent(() => import("@/components/Sidebar.vue"));
const ConfirmationModal = defineAsyncComponent(() => import("@/components/ConfirmationModal.vue"));
const Searchbar = defineAsyncComponent(() => import("@/components/Searchbar.vue")); // Add this import

const searchQuery = ref('')
const activeFilters = ref([])

const handleSearch = (query) => {
  searchQuery.value = query
}

const handleFilterChange = (filters) => {
  activeFilters.value = filters
}

const filteredCourses = computed(() => {
  return courses.value.filter(course => {
    const matchesSearch = searchQuery.value === '' ||
      course.ecode.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      course.subjectCode.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      course.courseName.toLowerCase().includes(searchQuery.value.toLowerCase())

    const matchesFilters = activeFilters.value.length === 0 ||
      activeFilters.value.some(filter =>
        course.units.toString().includes(filter)
      )

    return matchesSearch && matchesFilters
  })
})

const isConfirmationModalOpen = ref(false)
const isAddModalOpen = ref(false)
const isEditModalOpen = ref(false)
const filterContainerActive = ref(false)

const confirmationData = ref({
  title: '',
  itemName: '',
  itemInfo: null
})

const courseForm = ref({
  ecode: '',
  subjectCode: '',
  courseName: '',
  units: ''
})

const courses = ref([
  {
    ecode: 'CS101',
    subjectCode: 'COMP1001',
    courseName: 'Introduction to Programming',
    units: 3
  },
  {
    ecode: 'CS102',
    subjectCode: 'COMP1002',
    courseName: 'Data Structures and Algorithms',
    units: 4
  }
])

const showDeleteConfirmation = (title, itemName, itemInfo = null) => {
  confirmationData.value = { title, itemName, itemInfo }
  isConfirmationModalOpen.value = true
}

const handleConfirmDelete = (itemInfo) => {
  console.log('Deleting:', itemInfo)
}

const toggleFilterContainer = () => {
  filterContainerActive.value = !filterContainerActive.value
}

const openAddCourseModal = () => {
  courseForm.value = {
    ecode: '',
    subjectCode: '',
    courseName: '',
    units: ''
  }
  isAddModalOpen.value = true
}

const closeAddCourseModal = () => {
  isAddModalOpen.value = false
}

const openEditCourse = (ecode) => {
  const course = courses.value.find(c => c.ecode === ecode)
  if (course) {
    courseForm.value = { ...course }
    isEditModalOpen.value = true
  }
}

const closeEditCourseModal = () => {
  isEditModalOpen.value = false
}

const handleAddCourseSubmit = (e) => {
  e.preventDefault()
  courses.value.push({ ...courseForm.value })
  closeAddCourseModal()
}

const handleEditCourseSubmit = (e) => {
  e.preventDefault()
  const index = courses.value.findIndex(c => c.ecode === courseForm.value.ecode)
  if (index !== -1) {
    courses.value[index] = { ...courseForm.value }
  }
  closeEditCourseModal()
}
</script>

<template>
  <main>
    <div class="sidebar">
      <Sidebar />
    </div>

    <section>
      <div class="container">
        <div class="welcome-header" style="padding-bottom: 0;">
          <h1>Course Management</h1>
          <p>View and manage course information</p>
        </div>

        <div class="students-controls">
          <div class="search-filters">
            <Searchbar 
              v-model="searchQuery"
              @update:search-query="handleSearch"
              @filter-change="handleFilterChange"
            />
            <div class="filter-buttons">
            <button class="add-student-btn" @click="openAddCourseModal">
              <i class="fa-solid fa-plus"></i> Add New Course
            </button>
          </div>
          </div>
        </div>

        <div class="students-table">
          <table>
            <thead>
              <tr>
                <th>Course ECODE</th>
                <th>Subject Code</th>
                <th>Course Name</th>
                <th>Units</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="course in filteredCourses" :key="course.ecode">
                <td>{{ course.ecode }}</td>
                <td>{{ course.subjectCode }}</td>
                <td>{{ course.courseName }}</td>
                <td>{{ course.units }}</td>
                <td class="action-buttons">
                  <button class="action-btn edit-btn" @click="openEditCourse(course.ecode)">
                    <i class="fa-solid fa-pen-to-square"></i>
                  </button>
                  <button class="action-btn delete-btn" @click="showDeleteConfirmation(
                    'Delete Course',
                    `${course.courseName}`,
                    course
                  )">
                    <i class="fa-solid fa-trash"></i>
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <div class="modal" :class="{ active: isAddModalOpen }" @click="closeAddCourseModal">
        <div class="modal-content" @click.stop>
          <div class="modal-header">
            <h2>Add New Course</h2>
            <button class="close-modal" @click="closeAddCourseModal">&times;</button>
          </div>
          <form @submit="handleAddCourseSubmit" class="student-form">
            <div class="course-info-grid">
              <div class="info-group">
                <label>Course ECODE</label>
                <input type="text" v-model="courseForm.ecode" required>
              </div>
              <div class="info-group">
                <label>Subject Code</label>
                <input type="text" v-model="courseForm.subjectCode" required>
              </div>
              <div class="info-group full-width">
                <label>Course Name</label>
                <input type="text" v-model="courseForm.courseName" required>
              </div>
              <div class="info-group">
                <label>Units</label>
                <input type="number" min="1" max="6" v-model="courseForm.units" required>
              </div>
            </div>
            <div class="form-actions">
              <button type="button" class="cancel-btn" @click="closeAddCourseModal">Cancel</button>
              <button type="submit" class="submit-btn">Add Course</button>
            </div>
          </form>
        </div>
      </div>

      <div class="modal" :class="{ active: isEditModalOpen }" @click="closeEditCourseModal">
        <div class="modal-content" @click.stop>
          <div class="modal-header">
            <h2>Edit Course</h2>
            <button class="close-modal" @click="closeEditCourseModal">&times;</button>
          </div>
          <form @submit="handleEditCourseSubmit" class="student-form">
            <div class="course-info-grid">
              <div class="info-group">
                <label>Course ECODE</label>
                <input type="text" v-model="courseForm.ecode" required>
              </div>
              <div class="info-group">
                <label>Subject Code</label>
                <input type="text" v-model="courseForm.subjectCode" required>
              </div>
              <div class="info-group full-width">
                <label>Course Name</label>
                <input type="text" v-model="courseForm.courseName" required>
              </div>
              <div class="info-group">
                <label>Units</label>
                <input type="number" min="1" max="6" v-model="courseForm.units" required>
              </div>
            </div>
            <div class="form-actions">
              <button type="button" class="cancel-btn" @click="closeEditCourseModal">Cancel</button>
              <button type="submit" class="submit-btn">Save Changes</button>
            </div>
          </form>
        </div>
      </div>

      <ConfirmationModal :is-open="isConfirmationModalOpen" :title="confirmationData.title"
        :item-name="confirmationData.itemName" :item-info="confirmationData.itemInfo"
        @close="isConfirmationModalOpen = false" @confirm="handleConfirmDelete" />
    </section>
  </main>
</template>
