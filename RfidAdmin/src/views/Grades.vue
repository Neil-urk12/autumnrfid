<script setup>
import { ref, computed, onMounted } from 'vue'
import ConfirmationModal from '@/components/ConfirmationModal.vue'
import Sidebar from "@/components/Sidebar.vue"

const isConfirmationModalOpen = ref(false)
const isAddModalOpen = ref(false)
const isEditModalOpen = ref(false)
const isDeleteModalOpen = ref(false)
const filterContainerActive = ref(false)
const isAddSubjectDisabled = computed(() => !selectedGradingPeriod.value)
const addGradingPeriod = ref('')
const editGradingPeriod = ref('')

const confirmationData = ref({
    title: '',
    itemName: '',
    itemInfo: null
})

const studentInfo = ref({
    id: '',
    name: '',
    course: ''
})

const selectedGradingPeriod = ref('')
const searchQuery = ref('')

const subjects = ref([
    { code: "ENG", name: "English" },
    { code: "MATH", name: "Mathematics" },
    { code: "CHAR", name: "Conduct/Char. Build. Act" },
    { code: "SCI", name: "Science & Tech." },
    { code: "COMP", name: "Computer" },
    { code: "CIG", name: "Current Issues/Geography" }
])

const editSubjects = ref([
    { code: "SCI", name: "Science & Tech.", grade: "" },
    { code: "COMP", name: "Computer", grade: "" }
])

const deleteSubjects = ref([
    { code: "COM", name: "Computer Fundamentals", grade: "1.2" },
    { code: "TECH", name: "Information", grade: "1" }
])

const availableSubjects = ref([
    { code: "PHY", name: "Physics" },
    { code: "CHEM", name: "Chemistry" },
    { code: "BIO", name: "Biology" },
])

const newSubject = ref({
    code: '',
    name: '',
    grade: ''
})

const showNewSubjectForm = ref(false)

const showDeleteConfirmation = (title, itemName, itemInfo = null) => {
    confirmationData.value = { title, itemName, itemInfo }
    isConfirmationModalOpen.value = true
}

const handleConfirmDelete = (itemInfo) => {
    console.log('Deleting:', itemInfo)
    if (itemInfo.code) {
        deleteSubjects.value = deleteSubjects.value.filter(subject => subject.code !== itemInfo.code)
    }
    isConfirmationModalOpen.value = false
}

const toggleFilterContainer = () => {
    filterContainerActive.value = !filterContainerActive.value
}

const openAddGradesModal = (studentId) => {
    studentInfo.value = {
        id: studentId,
        name: "John Cez",
        course: "Computer Science"
    }
    isAddModalOpen.value = true
}

const closeAddGradesModal = () => {
    isAddModalOpen.value = false
}

const openEditGrades = (studentId) => {
    studentInfo.value = {
        id: studentId,
        name: "John Cez",
        course: "Computer Science"
    }
    isEditModalOpen.value = true
}

const closeEditGradesModal = () => {
    isEditModalOpen.value = false
}

const openDeleteGrades = (studentId) => {
    studentInfo.value = {
        id: studentId,
        name: "John Cez"
    }
    isDeleteModalOpen.value = true
}


const closeDeleteGradesModal = () => {
    isDeleteModalOpen.value = false
}

const addNewSubject = (subject = null) => {
    if (subject) {
        subjects.value.push({ ...subject, grade: '' })
    } else {
        showNewSubjectForm.value = true
    }
}

const handleAddCustomSubject = () => {
    if (newSubject.value.name && newSubject.value.code) {
        subjects.value.push({ ...newSubject.value, grade: '' })
        newSubject.value = { code: '', name: '', grade: '' }
        showNewSubjectForm.value = false
    }
}
const cancelNewSubject = () => {
    showNewSubjectForm.value = false
    newSubject.value = { code: '', name: '', grade: '' }
}

const deleteSubject = (subject) => {
    showDeleteConfirmation(

    )
}

const confirmDeleteSubjects = () => {
    closeDeleteGradesModal()
}

const handleSubmit = (e) => {
    e.preventDefault()
}

onMounted(() => {
    const handleClickOutside = (e, modalRef, closeFunction) => {
        if (e.target === modalRef) {
            closeFunction()
        }
    }
})
</script>

<template>
    <div class="sidebar">
        <Sidebar />
    </div>

    <section>
        <div style="width: 100%; max-width: 1200px; display: flex; flex-direction: column; gap: 20px;">
            <div class="welcome-header" style="padding-bottom: 0;">
                <h1>Student Grades</h1>
                <p>View and manage student grades</p>
            </div>

            <div class="students-controls" style="margin-top: 5px;">
                <div class="search-filters">
                    <div class="search-bar">
                        <i class="fa-solid fa-search"></i>
                        <input type="text" placeholder="Search students..." v-model="searchQuery">
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
                    <div class="filter-buttons">
                        <button class="filter-btn active">All</button>
                        <button class="filter-btn">Passed</button>
                        <button class="filter-btn">Failed</button>
                        <button class="filter-btn">Incomplete</button>
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
                            <th>Prelim</th>
                            <th>Midterm</th>
                            <th>Finals</th>
                            <th>Final Grade</th>
                            <th>Remarks</th>
                            <th>Actions</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr>
                            <td>2023-0001</td>
                            <td>John Cez</td>
                            <td>Computer Science</td>
                            <td>92</td>
                            <td>88</td>
                            <td>90</td>
                            <td>90</td>
                            <td><span class="status-badge status-active">Passed</span></td>
                            <td class="action-buttons">
                                <button class="action-btn" @click="openAddGradesModal('2023-0001')">
                                    <i class="fa-solid fa-plus"></i>
                                </button>
                                <button class="action-btn edit-btn" @click="openEditGrades('2023-0001')">
                                    <i class="fa-solid fa-pen-to-square"></i>
                                </button>
                                <button class="action-btn delete-btn" @click="openDeleteGrades('2023-0001')">
                                    <i class="fa-solid fa-trash"></i>
                                </button>
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>

        <div class="modal" :class="{ active: isAddModalOpen }" @click="closeAddGradesModal">
            <div class="modal-content" @click.stop>
                <div class="modal-header">
                    <h2>Add Student Grades</h2>
                    <button class="close-modal" @click="closeAddGradesModal">&times;</button>
                </div>
                <form @submit="handleSubmit" class="student-form">
                    <div class="student-info-grid">
                        <div class="info-group">
                            <label>Student ID</label>
                            <input type="text" v-model="studentInfo.id" required readonly>
                        </div>
                        <div class="info-group">
                            <label>Name</label>
                            <input type="text" v-model="studentInfo.name" required readonly>
                        </div>
                        <div class="info-group">
                            <label>Course</label>
                            <input type="text" v-model="studentInfo.course" required readonly>
                        </div>
                        <div class="info-group">
                            <label>Grading Period</label>
                            <select v-model="selectedGradingPeriod" required>
                                <option value="">Select Period</option>
                                <option value="prelim">Prelim</option>
                                <option value="midterm">Midterm</option>
                                <option value="finals">Finals</option>
                            </select>
                        </div>
                    </div>

                    <div class="subjects-section">
                        <div class="subjects-header">
                            <h3>Subject Grades</h3>
                            <div class="subject-actions" v-if="selectedGradingPeriod">
                                <button type="button" class="add-subject-btn" @click="addNewSubject()">
                                    <i class="fa-solid fa-plus"></i> Add Custom Subject
                                </button>
                                <div class="subject-select">
                                    <select @change="(e) => addNewSubject(JSON.parse(e.target.value))"
                                        v-model="selectedSubject">
                                        <option value="">Select Existing Subject</option>
                                        <option v-for="subject in availableSubjects" :key="subject.code"
                                            :value="JSON.stringify(subject)">
                                            {{ subject.name }}
                                        </option>
                                    </select>
                                </div>
                            </div>
                        </div>

                        <div v-if="showNewSubjectForm" class="new-subject-form">
                            <div class="form-row">
                                <div class="form-group">
                                    <label>Subject Code</label>
                                    <input type="text" v-model="newSubject.code" required>
                                </div>
                                <div class="form-group">
                                    <label>Subject Name</label>
                                    <input type="text" v-model="newSubject.name" required>
                                </div>
                            </div>
                            <div class="form-actions">
                                <button class="buttonss" type="button" @click="cancelNewSubject">Cancel</button>
                                <button class="buttonss" type="button" @click="handleAddCustomSubject">Add</button>
                            </div>
                        </div>

                        <div class="subjects-list">
                            <template v-if="selectedGradingPeriod">
                                <div class="subjects-grid">
                                    <div v-for="(subject, index) in subjects" :key="subject.code" class="subject-row"
                                        :class="{ 'right-column': index % 2 !== 0 }">
                                        <div class="grade-input">
                                            <label>{{ subject.name }}</label>
                                            <input type="number" min="0" max="100" required placeholder="Grade"
                                                v-model="subject.grade">
                                        </div>
                                    </div>
                                </div>
                            </template>
                            <p v-else style="color: var(--text-secondary); text-align: center;">
                                Please select a grading period
                            </p>
                        </div>
                    </div>

                    <div class="form-actions">
                        <button type="button" class="cancel-btn" @click="closeAddGradesModal">Cancel</button>
                        <button type="submit" class="submit-btn">Add Grades</button>
                    </div>
                </form>
            </div>
        </div>

        <div class="modal" :class="{ active: isEditModalOpen }" @click="closeEditGradesModal">
            <div class="modal-content" @click.stop>
                <div class="modal-header">
                    <h2>Edit Student Grades</h2>
                    <button class="close-modal" @click="closeEditGradesModal">&times;</button>
                </div>
                <form @submit="handleSubmit" class="student-form">
                    <div class="student-info-grid">
                        <div class="info-group">
                            <label>Student ID</label>
                            <input type="text" v-model="studentInfo.id" required readonly>
                        </div>
                        <div class="info-group">
                            <label>Name</label>
                            <input type="text" v-model="studentInfo.name" required readonly>
                        </div>
                        <div class="info-group">
                            <label>Course</label>
                            <input type="text" v-model="studentInfo.course" required readonly>
                        </div>
                        <div class="info-group">
                            <label>Grading Period</label>
                            <select v-model="selectedGradingPeriod" required>
                                <option value="">Select Period</option>
                                <option value="prelim">Prelim</option>
                                <option value="midterm">Midterm</option>
                                <option value="finals">Finals</option>
                            </select>
                        </div>
                    </div>

                    <div class="subjects-section">
                        <div class="subjects-header">
                            <h3>Subject Grades</h3>
                        </div>
                        <div class="subjects-list">
                            <template v-if="selectedGradingPeriod">
                                <div class="subjects-grid">
                                    <div v-for="(subject, index) in editSubjects" :key="subject.code"
                                        class="subject-row" :class="{ 'right-column': index % 2 !== 0 }">
                                        <div class="grade-input">
                                            <label>{{ subject.name }}</label>
                                            <input type="number" min="0" max="100" required placeholder="Grade"
                                                v-model="subject.grade">
                                        </div>
                                    </div>
                                </div>
                            </template>
                            <p v-else style="color: var(--text-secondary); text-align: center;">
                                Please select a grading period
                            </p>
                        </div>
                    </div>

                    <div class="form-actions">
                        <button type="button" class="cancel-btn" @click="closeEditGradesModal">Cancel</button>
                        <button type="submit" class="submit-btn">Save Changes</button>
                    </div>
                </form>
            </div>
        </div>

        <div class="modal" :class="{ active: isDeleteModalOpen }" @click="closeDeleteGradesModal">
            <div class="modal-content" @click.stop>
                <div class="modal-header">
                    <h2>Delete Subject Grades</h2>
                    <button class="close-modal" @click="closeDeleteGradesModal">&times;</button>
                </div>
                <div class="student-info-grid">
                    <div class="info-group">
                        <label>Student ID</label>
                        <span>{{ studentInfo.id }}</span>
                    </div>
                    <div class="info-group">
                        <label>Name</label>
                        <span>{{ studentInfo.name }}</span>
                    </div>
                </div>
                <div class="subjects-section">
                    <div class="subjects-header">
                        <h3>Select Subjects to Delete</h3>
                    </div>
                    <div class="delete-subjects-list">
                        <div class="subjects-grid">
                            <div v-for="(subject, index) in deleteSubjects" :key="subject.code"
                                class="delete-subject-item" :class="{ 'right-column': index % 2 !== 0 }"
                                :data-subject="subject.code">
                                <div class="subject-info">
                                    <span class="subject-name">{{ subject.name }}</span>
                                    <span class="subject-grade">Grade: {{ subject.grade }}</span>
                                </div>
                                <button type="button" class="delete-subject-btn" @click="showDeleteConfirmation(
                                    'Delete Subject Grade',
                                    `${subject.name}`,
                                    subject
                                )">
                                    <i class="fa-solid fa-times"></i>
                                </button>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="form-actions">
                    <button type="button" class="cancel-btn" @click="closeDeleteGradesModal">Cancel</button>
                    <button type="button" class="delete-confirm-btn" @click="confirmDeleteSubjects">Delete
                        Selected</button>
                </div>
            </div>
        </div>

        <ConfirmationModal :is-open="isConfirmationModalOpen" :title="confirmationData.title"
            :item-name="confirmationData.itemName" :item-info="confirmationData.itemInfo"
            @close="isConfirmationModalOpen = false" @confirm="handleConfirmDelete" />
    </section>
</template>
