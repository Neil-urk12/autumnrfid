<script setup>
import { ref, computed } from 'vue'
import { defineAsyncComponent } from "vue";
const Sidebar = defineAsyncComponent(() => import("@/components/Sidebar.vue"));
const ConfirmationModal = defineAsyncComponent(() => import("@/components/ConfirmationModal.vue"));
const UnsavedChangesModal = defineAsyncComponent(() => import("@/components/UnsavedChangesModal.vue"));
const Searchbar = defineAsyncComponent(() => import("@/components/Searchbar.vue"));

// MOCK 
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

const courseForm = ref({
  ecode: '',
  subjectCode: '',
  courseName: '',
  units: '',
  originalEcode: '' 
})

// STATE 
const searchQuery = ref('')
const activeFilters = ref([])
const isConfirmationModalOpen = ref(false)
const isAddModalOpen = ref(false)
const isEditModalOpen = ref(false)
const filterContainerActive = ref(false)
const isUnsavedChangesModalOpen = ref(false)
const modalToClose = ref(null)
const ecodeError = ref('')
const addEcodeError = ref('')
const editEcodeError = ref('')

// CONFIRMATION MODAL DATA
const confirmationData = ref({
  title: '',
  itemName: '',
  itemInfo: null
})

// HANDLES THE SEARCH QUERY UPDATE
const handleSearch = (query) => {
  searchQuery.value = query
}

// HANDLES THE FILTER CHANGES
const handleFilterChange = (filters) => {
  activeFilters.value = filters
}

// FILTERS THE COURSES BASED ON SEARCH QUERY AND ACTIVE FILTERS
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

// SHOWS THE DELETE CONFIRMATION MODAL
const showDeleteConfirmation = (title, itemName, itemInfo = null) => {
  confirmationData.value = { title, itemName, itemInfo }
  isConfirmationModalOpen.value = true
}

// HANDLES THE COURSE DELETION AFTER CONFIRMATION
const handleConfirmDelete = (itemInfo) => {
  const index = courses.value.findIndex(course => course.ecode === itemInfo.ecode)
  if (index !== -1) {
    courses.value.splice(index, 1)
  }
  isConfirmationModalOpen.value = false
}

// TOGGLES THE FILTER CONTAINER VISIBILITY
const toggleFilterContainer = () => {
  filterContainerActive.value = !filterContainerActive.value
}

// OPENS THE ADD COURSE MODAL WITH EMPTY FORM
const openAddCourseModal = () => {
  courseForm.value = {
    ecode: '',
    subjectCode: '',
    courseName: '',
    units: ''
  }
  isAddModalOpen.value = true
}

// CLOSES THE ADD COURSE MODAL WITH UNSAVED CHANGES CHECK
const closeAddCourseModal = (fromSave = false) => {
  if (hasUnsavedChanges.value && !modalToClose.value && !fromSave) {
    modalToClose.value = 'add'
    isUnsavedChangesModalOpen.value = true
  } else {
    isAddModalOpen.value = false
    courseForm.value = {
      ecode: '',
      subjectCode: '',
      courseName: '',
      units: ''
    }
    addEcodeError.value = ''
    modalToClose.value = null
  }
}

// OPENS THE EDIT COURSE MODAL WITH EXISTING COURSE DATA
const openEditCourse = (ecode) => {
  const course = courses.value.find(c => c.ecode === ecode)
  if (course) {
    courseForm.value = { 
      ...course,
      originalEcode: course.ecode
    }
    isEditModalOpen.value = true
  }
}

// CLOSES THE EDIT COURSE MODAL WITH UNSAVED CHANGES CHECK
const closeEditCourseModal = () => {
  if (hasUnsavedChanges.value) {
    modalToClose.value = 'edit'
    isUnsavedChangesModalOpen.value = true
  } else {
    isEditModalOpen.value = false
    editEcodeError.value = ''
  }
}

// TRANSFORMS THE COURSE NAME TO CAPITALIZE THE FIRST LETTER OF EACH WORD, EXCEPT FOR CERTAIN EXCEPTIONS
const formatCourseName = (name) => {
  const exceptions = ['of', 'to', 'and']
  return name
    .trim()
    .toLowerCase()
    .split(' ')
    .map((word, index) => {
      if (exceptions.includes(word)) {
        return word
      }
      return word.charAt(0).toUpperCase() + word.slice(1).toLowerCase()
    })
    .join(' ')
}

// HANDLES THE ADD COURSE FORM SUBMISSION
const handleAddCourseSubmit = (e) => {
  e.preventDefault()
  
  const courseData = { ...courseForm.value }
  
  // CHANGE ECODE AND SUBJECT CODE TO UPPERCASE AND FORMAT COURSE NAME
  courseData.ecode = courseData.ecode.toUpperCase()
  courseData.subjectCode = courseData.subjectCode.toUpperCase()
  courseData.courseName = formatCourseName(courseData.courseName)
  
  if (!validateNewEcode(courseData.ecode)) {
    return
  }
  
  courses.value.push(courseData)
  closeAddCourseModal(true)
}

// HANDLES THE EDIT COURSE FORM SUBMISSION
const handleEditCourseSubmit = (e) => {
  e.preventDefault()
  
  const courseData = { ...courseForm.value }
  
  // CHANGE ECODE AND SUBJECT CODE TO UPPERCASE AND FORMAT COURSE NAME
  courseData.ecode = courseData.ecode.toUpperCase()
  courseData.subjectCode = courseData.subjectCode.toUpperCase()
  courseData.courseName = formatCourseName(courseData.courseName)
  
  const index = courses.value.findIndex(course => 
    course.ecode === courseForm.value.originalEcode
  )

  if (index !== -1) {
    if (!validateEditEcode(courseData.ecode, courseForm.value.originalEcode)) {
      return
    }

    courses.value[index] = courseData
    isEditModalOpen.value = false
    editEcodeError.value = ''
  }
}

// VALIDATES THE ECODE FOR BOTH ADD AND EDIT OPERATIONS
const validateEcode = (ecode, originalEcode = null) => {
  ecodeError.value = ''
  
  // CHECKS IF THE ECODE ALREADY EXISTS IN THE COURSES ARRAY
  const exists = courses.value.some(course => 
    course.ecode === ecode && course.ecode !== originalEcode
  )
  
  if (exists) {
    ecodeError.value = 'Course already exists'
    return false
  }
  
  return true
}

// VALIDATES THE ECODE FOR NEW COURSE ADDITION
const validateNewEcode = (ecode) => {
  addEcodeError.value = ''
  
  const exists = courses.value.some(course => 
    course.ecode.toUpperCase() === ecode.toUpperCase()
  )
  
  if (exists) {
    addEcodeError.value = 'Course ECODE already exists'
    setTimeout(() => {
      addEcodeError.value = ''
    }, 3000)
    return false
  }
  
  return true
}

// VALIDATES THE ECODE FOR COURSE EDITING
const validateEditEcode = (ecode, originalEcode) => {
  editEcodeError.value = ''
  
  const exists = courses.value.some(course => 
    course.ecode.toUpperCase() === ecode.toUpperCase() && 
    course.ecode.toUpperCase() !== originalEcode.toUpperCase()
  )
  
  if (exists) {
    editEcodeError.value = 'Course already exists'
    setTimeout(() => {
      editEcodeError.value = ''
    }, 3000)
    return false
  }
  
  return true
}

// CHECKS FOR UNSAVED CHANGES IN THE FORM
const hasUnsavedChanges = computed(() => {
  if (isAddModalOpen.value) {
    return Object.values(courseForm.value).some(value => value !== '')
  } else if (isEditModalOpen.value) {
    const originalCourse = courses.value.find(c => c.ecode === courseForm.value.originalEcode)
    if (!originalCourse) return false
    
    return Object.keys(courseForm.value).some(key => 
      key !== 'originalEcode' && courseForm.value[key] !== originalCourse[key]
    )
  }
  return false
})

// HANDLES THE UNSAVED CHANGES MODAL ACTIONS
const handleUnsavedChanges = (confirm) => {
  if (confirm) {
    if (modalToClose.value === 'add') {
      isAddModalOpen.value = false
      courseForm.value = {
        ecode: '',
        subjectCode: '',
        courseName: '',
        units: ''
      }
    } else if (modalToClose.value === 'edit') {
      isEditModalOpen.value = false
      editEcodeError.value = ''
    }
  }
  isUnsavedChangesModalOpen.value = false
  modalToClose.value = null
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

        <!-- SEARCH BAR SECTION -->
        <div class="students-controls">
          <div class="search-filters">
            <Searchbar v-model="searchQuery" @update:search-query="handleSearch" @filter-change="handleFilterChange" />
            <div class="filter-buttons">
              <button class="add-student-btn" @click="openAddCourseModal">
                <i class="fa-solid fa-plus"></i> Add New Course
              </button>
            </div>
          </div>
        </div>
        
        <!-- COURSE TABLE -->
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
     
      <!-- ADD NEW COURSE MODAL -->
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
                <input 
                  type="text" 
                  v-model="courseForm.ecode" 
                  :class="{ 'error': addEcodeError }" placeholder="Enter Course ECODE"
                  required
                >
                <span class="error-message" v-if="addEcodeError">{{ addEcodeError }}</span>
              </div>
              <div class="info-group">
                <label>Subject Code</label>
                <input type="text" v-model="courseForm.subjectCode" placeholder="Enter Course Subject Code" required>
              </div>
              <div class="info-group full-width">
                <label>Course Name</label>
                <input type="text" v-model="courseForm.courseName" placeholder="Enter Course Name" required>
              </div>
              <div class="info-group">
                <label>Units</label>
                <input type="number" min="1" max="6" v-model="courseForm.units" placeholder="Enter Course Units" required>
              </div>
            </div>
            <div class="form-actions">
              <button type="button" class="cancel-btn" @click="closeAddCourseModal">Cancel</button>
              <button type="submit" class="submit-btn">Add Course</button>
            </div>
          </form>
        </div>
      </div>
      
      <!-- EDIT COURSE MODAL -->
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
                <input 
                  type="text" 
                  v-model="courseForm.ecode" 
                  :class="{ 'error': editEcodeError }"
                  required
                >
                <span class="error-message" v-if="editEcodeError">{{ editEcodeError }}</span>
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
      
      <!-- DELETE CONFIRMATION MODAL -->
      <ConfirmationModal :is-open="isConfirmationModalOpen" :title="confirmationData.title"
        :item-name="confirmationData.itemName" :item-info="confirmationData.itemInfo"
        placeholder-text="Enter Course ECODE to Confirm Deletion" @close="isConfirmationModalOpen = false"
        @confirm="handleConfirmDelete" />
      
      <!-- UNSAVED CHANGES MODAL -->
      <UnsavedChangesModal 
        :is-open="isUnsavedChangesModalOpen"
        @close="handleUnsavedChanges(false)"
        @confirm="handleUnsavedChanges(true)"
      />
    </section>
  </main>
</template>


