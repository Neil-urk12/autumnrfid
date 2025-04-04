<script setup>
import { ref } from 'vue'
import { defineAsyncComponent } from "vue";
const Sidebar = defineAsyncComponent(() => import("@/components/Sidebar.vue"));
const ConfirmationModal = defineAsyncComponent(() => import("@/components/ConfirmationModal.vue"));

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
  <div class="sidebar">
    <Sidebar />
  </div>

  <section>
    <div style="width: 100%; max-width: 1200px; display: flex; flex-direction: column; gap: 20px;">
      <div class="welcome-header" style="padding-bottom: 0;">
        <h1>Course Management</h1>
        <p>View and manage course information</p>
      </div>

      <div class="students-controls" style="margin-top: 5px;">
        <div class="search-filters">
          <div class="search-bar">
            <i class="fa-solid fa-search"></i>
            <input type="text" placeholder="Search students...">
            <button class="filter-toggle-btn" @click="toggleFilterContainer">
              <i class="fa-solid fa-filter"></i>
              Filter
            </button>

            <div class="filter-tooltip" :class="{ active: filterContainerActive }">
              <div class="filter-category">
                <span class="category-label">Year Level</span>
                <div class="filter-buttons-row">
                  <button class="filter-btn">1st</button>
                  <button class="filter-btn">2nd</button>
                  <button class="filter-btn">3rd</button>
                  <button class="filter-btn">4th</button>
                </div>
              </div>

              <div class="filter-category">
                <span class="category-label">Semester</span>
                <div class="filter-buttons-row">
                  <button class="filter-btn">1st Sem</button>
                  <button class="filter-btn">2nd Sem</button>
                </div>
              </div>

              <div class="filter-category">
                <span class="category-label">Course</span>
                <div class="filter-buttons-row">
                  <button class="filter-btn">BSIT</button>
                  <button class="filter-btn">BSCS</button>
                  <button class="filter-btn">BSIS</button>
                </div>
              </div>

              <div class="filter-category">
                <span class="category-label">Block</span>
                <div class="filter-buttons-row">
                  <button class="filter-btn">A</button>
                  <button class="filter-btn">B</button>
                  <button class="filter-btn">C</button>
                </div>
              </div>
            </div>
          </div>
          <button class="add-student-btn" @click="openAddCourseModal">
            <i class="fa-solid fa-plus"></i> Add New Course
          </button>
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
            <tr v-for="course in courses" :key="course.ecode">
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

    <ConfirmationModal 
      :is-open="isConfirmationModalOpen"
      :title="confirmationData.title"
      :item-name="confirmationData.itemName"
      :item-info="confirmationData.itemInfo"
      @close="isConfirmationModalOpen = false"
      @confirm="handleConfirmDelete"
    />
  </section>
</template>
