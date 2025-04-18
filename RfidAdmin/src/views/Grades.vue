<script setup>
import { ref, computed, onMounted, defineAsyncComponent } from 'vue'
const ConfirmationModal = defineAsyncComponent(() => import('@/components/ConfirmationModal.vue'))
const UnsavedChangesModal = defineAsyncComponent(() => import('@/components/UnsavedChangesModal.vue'))
const Sidebar = defineAsyncComponent(() => import("@/components/Sidebar.vue"))
const Searchbar = defineAsyncComponent(() => import("@/components/Searchbar.vue"))
import Toast from '@/components/Toast.vue'

// MOCK DATA
const validGrades = [
    '1', '1.25', '1.50', '1.75',
    '2', '2.25', '2.50', '2.75',
    '3', '3.25', '3.50', '3.75',
    '4', '5', '0', 'INC'
]

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

const students = ref([
    {
        id: '2023-0001',
        name: 'John Cez',
        course: 'Computer Science',
        grades: {
            prelim: {
                'ENG': '',
                'MATH': '',
                'CHAR': '',
                'SCI': '',
                'COMP': '',
                'CIG': ''
            },
            midterm: {
                'ENG': '',
                'MATH': '',
                'CHAR': '',
                'SCI': '',
                'COMP': '',
                'CIG': ''
            },
            prefinals: {
                'ENG': '',
                'MATH': '',
                'CHAR': '',
                'SCI': '',
                'COMP': '',
                'CIG': ''
            },
            finals: {
                'ENG': '',
                'MATH': '',
                'CHAR': '',
                'SCI': '',
                'COMP': '',
                'CIG': ''
            }
        },
        gwa: '',
        remarks: ''
    },
    {
        id: '2023-0002',
        name: 'Jan Rosa',
        course: 'Computer Science',
        grades: {
            prelim: {
                'ENG': '',
                'MATH': '',
                'CHAR': '',
                'SCI': '',
                'COMP': '',
                'CIG': ''
            },
            midterm: {
                'ENG': '',
                'MATH': '',
                'CHAR': '',
                'SCI': '',
                'COMP': '',
                'CIG': ''
            },
            prefinals: {
                'ENG': '',
                'MATH': '',
                'CHAR': '',
                'SCI': '',
                'COMP': '',
                'CIG': ''
            },
            finals: {
                'ENG': '',
                'MATH': '',
                'CHAR': '',
                'SCI': '',
                'COMP': '',
                'CIG': ''
            }
        },
        gwa: '',
        remarks: ''
    },
    {
        id: '2023-0003',
        name: 'Lightning Mcqueen',
        course: 'Information Technology',
        grades: {
            prelim: {
                'ENG': '',
                'MATH': '',
                'CHAR': '',
                'SCI': '',
                'COMP': '',
                'CIG': ''
            },
            midterm: {
                'ENG': '',
                'MATH': '',
                'CHAR': '',
                'SCI': '',
                'COMP': '',
                'CIG': ''
            },
            prefinals: {
                'ENG': '',
                'MATH': '',
                'CHAR': '',
                'SCI': '',
                'COMP': '',
                'CIG': ''
            },
            finals: {
                'ENG': '',
                'MATH': '',
                'CHAR': '',
                'SCI': '',
                'COMP': '',
                'CIG': ''
            }
        },
        gwa: '',
        remarks: ''
    }
])
//

// STATE 
const isConfirmationModalOpen = ref(false)
const isAddModalOpen = ref(false)
const isEditModalOpen = ref(false)
const isDeleteModalOpen = ref(false)
const isUnsavedChangesModalOpen = ref(false)
const modalToClose = ref(null)
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

const searchQuery = ref('')
const activeFilters = ref([])
const activeFilter = ref('all')
const selectedGradingPeriod = ref('')
const showNewSubjectForm = ref(false)
const showToast = ref(false)
const toastMessage = ref('')

const newSubject = ref({
    code: '',
    name: '',
    grade: ''
})

// HANDLES SEARCH AND FILTER
const filteredStudents = computed(() => {
    return students.value.filter(student => {
        const matchesSearch = !searchQuery.value ||
            student.id.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
            student.name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
            student.course.toLowerCase().includes(searchQuery.value.toLowerCase())

        const matchesStatus = activeFilter.value === 'all' ||
            (activeFilter.value === 'passed' && student.remarks.toLowerCase() === 'passed') ||
            (activeFilter.value === 'failed' && student.remarks.toLowerCase() === 'failed') ||
            (activeFilter.value === 'incomplete' && student.remarks.toLowerCase() === 'incomplete')

        return matchesSearch && matchesStatus
    })
})

// GETS THE CURRENT GRADING PERIOD FOR A STUDENT
const currentPeriod = computed(() => {
    const student = students.value.find(s => s.id === studentInfo.value.id)
    if (!student) return 'prelim'

    const periods = ['prelim', 'midterm', 'prefinals', 'finals']

    // FIND THE NEXT PERIOD WHERE NO GRADES ARE ENTERED
    const nextPeriod = periods.find(period => {
        return !Object.values(student.grades[period]).some(grade => grade)
    })

    return nextPeriod || ''
})

// GETS THE LIST OF EDITABLE PERIODS FOR A STUDENT
const editablePeriods = computed(() => {
    const student = students.value.find(s => s.id === studentInfo.value.id)
    if (!student) return []

    const periods = ['prelim', 'midterm', 'prefinals', 'finals']
    return periods.filter(period => hasGradesForPeriod(period))
})

// CHECKS IF THERE ARE UNSAVED CHANGES IN THE MODAL
const hasUnsavedChanges = computed(() => {
    if (isAddModalOpen.value) {
        return subjects.value.some(subject => {
            const cleanGrade = subject.grade.trim()
            return cleanGrade !== '' && validGrades.includes(cleanGrade)
        })
    } else if (isEditModalOpen.value) {
        const student = students.value.find(s => s.id === studentInfo.value.id)
        if (!student || !selectedGradingPeriod.value) return false

        return editSubjects.value.some(subject => {
            const currentGrade = subject.grade.trim()
            const originalGrade = student.grades[selectedGradingPeriod.value][subject.code] || ''

            return currentGrade !== originalGrade &&
                (currentGrade === '' || validGrades.includes(currentGrade))
        })
    }
    return false
})

// SHOWS THE DELETE CONFIRMATION MODAL
const showDeleteConfirmation = (title, itemName, itemInfo = null) => {
    confirmationData.value = {
        title,
        itemName,
        itemInfo
    }
    isConfirmationModalOpen.value = true
}

// HANDLES THE CONFIRMATION OF DELETION
const handleConfirmDelete = (itemInfo) => {
    if (itemInfo.name) {
        deleteSubjects.value = deleteSubjects.value.filter(
            subject => subject.name !== itemInfo.name
        )
    }
    isConfirmationModalOpen.value = false
}

// TOGGLES THE FILTER CONTAINER VISIBILITY
const toggleFilterContainer = () => {
    filterContainerActive.value = !filterContainerActive.value
}

// OPENS THE ADD GRADES MODAL
const openAddGradesModal = (studentId) => {
    const student = students.value.find(s => s.id === studentId)
    if (!student) return

    studentInfo.value = {
        id: studentId,
        name: student.name,
        course: student.course
    }

    // AUTOMATICALLY SET THE GRADING PERIOD
    selectedGradingPeriod.value = currentPeriod.value
    isAddModalOpen.value = true
}

// CLOSES THE ADD GRADES MODAL
const closeAddGradesModal = () => {
    if (hasUnsavedChanges.value) {
        modalToClose.value = 'add'
        isUnsavedChangesModalOpen.value = true
    } else {
        isAddModalOpen.value = false
        selectedGradingPeriod.value = ''
        subjects.value.forEach(subject => subject.grade = '')
    }
}

// CLOSES THE EDIT GRADES MODAL
const closeEditGradesModal = () => {
    if (hasUnsavedChanges.value) {
        modalToClose.value = 'edit'
        isUnsavedChangesModalOpen.value = true
    } else {
        isEditModalOpen.value = false
        selectedGradingPeriod.value = ''
        editSubjects.value.forEach(subject => subject.grade = '')
    }
}

// OPENS THE EDIT GRADES MODAL
const openEditGrades = (studentId) => {
    const student = students.value.find(s => s.id === studentId)
    if (!student) return

    studentInfo.value = {
        id: student.id,
        name: student.name,
        course: student.course
    }

    selectedGradingPeriod.value = ''

    editSubjects.value = subjects.value.map(subject => ({
        ...subject,
        grade: ''
    }))

    isEditModalOpen.value = true
}

// OPENS THE DELETE GRADES MODAL
const openDeleteGrades = (studentId) => {
    const student = students.value.find(s => s.id === studentId)
    if (!student) return

    studentInfo.value = {
        id: student.id,
        name: student.name,
        course: student.course
    }

    loadGradesForDelete()

    selectedGradingPeriod.value = ''
    isDeleteModalOpen.value = true
}

// CALCULATES THE AVERAGE GRADE FOR A SUBJECT
const calculateSubjectAverage = (student, subjectCode) => {
    const periods = ['prelim', 'midterm', 'prefinals', 'finals']
    const grades = periods.map(period => student.grades[period][subjectCode])
        .filter(grade => grade !== '')
        .map(Number)

    if (grades.length === 0) return '--'

    const total = grades.reduce((sum, grade) => sum + grade, 0)
    return (total / grades.length).toFixed(2)
}

// LOADS GRADES FOR DELETION
const loadGradesForDelete = () => {
    const student = students.value.find(s => s.id === studentInfo.value.id)
    if (!student) return

    const periods = ['prelim', 'midterm', 'prefinals', 'finals']

    deleteSubjects.value = subjects.value.map(subject => {
        const existingGrades = periods
            .map(period => student.grades[period][subject.code])
            .filter(grade => grade !== '')
            .map(Number)

        // CALCULATE OVERALL GRADES
        let totalGrade = '--'
        if (existingGrades.length > 0) {
            const sum = existingGrades.reduce((acc, grade) => acc + grade, 0)
            totalGrade = (sum / existingGrades.length).toFixed(2)
        }

        return {
            ...subject,
            grade: student.grades[selectedGradingPeriod.value]?.[subject.code] || '--', // CURRENT PERIOD GRADE
            totalGrade, // OVERALL AVERAGE OF THE SUBJECT ON ALL PERIODS
            periodGrades: periods.reduce((acc, period) => {
                acc[period] = student.grades[period][subject.code] || ''
                return acc
            }, {})
        }
    })
}

// CLOSES THE DELETE GRADES MODAL
const closeDeleteGradesModal = () => {
    isDeleteModalOpen.value = false
}

// ADDS A NEW SUBJECT
const addNewSubject = (subject = null) => {
    if (subject) {
        subjects.value.push({ ...subject, grade: '' })
    } else {
        showNewSubjectForm.value = true
    }
}

// HANDLES ADDING A CUSTOM SUBJECT
const handleAddCustomSubject = () => {
    if (newSubject.value.name && newSubject.value.code) {
        subjects.value.push({ ...newSubject.value, grade: '' })
        newSubject.value = { code: '', name: '', grade: '' }
        showNewSubjectForm.value = false
    }
}

// CANCELLS ADDING A NEW SUBJECT
const cancelNewSubject = () => {
    showNewSubjectForm.value = false
    newSubject.value = { code: '', name: '', grade: '' }
}

// CONFIRMS DELETION OF SUBJECTS
const confirmDeleteSubjects = () => {
    const student = students.value.find(s => s.id === studentInfo.value.id)
    if (!student) return

    Object.keys(student.grades[selectedGradingPeriod.value]).forEach(code => {
        student.grades[selectedGradingPeriod.value][code] = ''
    })

    student[selectedGradingPeriod.value] = ''
    computeFinalGrade(student.id)

    closeDeleteGradesModal()
}

// HANDLES GRADE INPUT VALIDATION
const handleGradeInput = (subject, value) => {
    const cleanValue = value.trim()

    if (cleanValue === '') {
        subject.grade = ''
        showToast.value = false
        toastMessage.value = ''
        return
    }

    // Check if the input is a valid grade (either whole number or decimal)
    const isValidGrade = validGrades.some(validGrade => {
        // Convert both to numbers for comparison
        const inputNum = Number(cleanValue)
        const validNum = Number(validGrade)
        return inputNum === validNum
    })

    if (!isValidGrade) {
        subject.grade = cleanValue
        toastMessage.value = 'Please enter a valid grade (1, 1.25, 1.50, etc)'
        showToast.value = true
    } else {
        // Store the grade as-is during editing
        subject.grade = cleanValue
        showToast.value = false
        toastMessage.value = ''
    }
}

// SYNCHRONIZES GRADES BETWEEN SOURCE AND TARGET ARRAYS
const syncGrades = (sourceArray, targetArray) => {
    sourceArray.forEach(sourceSubject => {
        const targetSubject = targetArray.find(t => t.code === sourceSubject.code)
        if (targetSubject) {
            targetSubject.grade = sourceSubject.grade
        }
    })
}

// HANDLES SUBMISSION OF ADD GRADES FORM
const handleSubmit = (e) => {
    e.preventDefault()

    const periodGrades = subjects.value.filter(subject => subject.grade)

    if (periodGrades.length > 0) {
        const hasInvalidGrades = periodGrades.some(subject => {
            const cleanGrade = subject.grade.trim()
            return !validGrades.some(validGrade => Number(cleanGrade) === Number(validGrade))
        })

        if (hasInvalidGrades) {
            toastMessage.value = 'Please enter valid grades for all subjects'
            showToast.value = true
            return
        }

        // UPDATE EACH SUBJECT GRADE
        const student = students.value.find(s => s.id === studentInfo.value.id)
        if (student) {
            periodGrades.forEach(subject => {
                // Convert whole numbers to decimal format when saving
                const grade = subject.grade.includes('.') ? subject.grade : `${subject.grade}.00`
                student.grades[selectedGradingPeriod.value][subject.code] = grade
            })

            // PERIOD AVERAGE
            const periodGradesArray = Object.values(student.grades[selectedGradingPeriod.value])
            const validGrades = periodGradesArray.filter(grade => grade)
            if (validGrades.length > 0) {
                const average = validGrades.reduce((sum, grade) => sum + Number(grade), 0) / validGrades.length
                student[selectedGradingPeriod.value] = average.toFixed(2)
            }
        }

        computeFinalGrade(studentInfo.value.id)

        editSubjects.value = subjects.value.map(subject => ({
            ...subject,
            grade: periodGrades.find(p => p.code === subject.code)?.grade || ''
        }))

        subjects.value.forEach(subject => subject.grade = '')
        selectedGradingPeriod.value = '' // Reset to default
        closeAddGradesModal()
    }
}

// HANDLES SUBMISSION OF EDIT GRADES FORM
const handleEditSubmit = (e) => {
    e.preventDefault()

    // GET SUBJECT GRADES FOR SELECTED PERIOD
    const periodGrades = editSubjects.value.filter(subject => subject.grade)

    if (periodGrades.length > 0) {
        const hasInvalidGrades = periodGrades.some(subject => {
            const cleanGrade = subject.grade.trim()
            return !validGrades.some(validGrade => Number(cleanGrade) === Number(validGrade))
        })

        if (hasInvalidGrades) {
            toastMessage.value = 'Please enter valid grades for all subjects'
            showToast.value = true
            return
        }

        const student = students.value.find(s => s.id === studentInfo.value.id)
        if (student) {
            periodGrades.forEach(subject => {
                // Convert whole numbers to decimal format when saving
                const grade = subject.grade.includes('.') ? subject.grade : `${subject.grade}.00`
                student.grades[selectedGradingPeriod.value][subject.code] = grade
            })

            const periodGradesArray = Object.values(student.grades[selectedGradingPeriod.value])
            const validGrades = periodGradesArray.filter(grade => grade)
            if (validGrades.length > 0) {
                const average = validGrades.reduce((sum, grade) => sum + Number(grade), 0) / validGrades.length
                student[selectedGradingPeriod.value] = average.toFixed(2)
            }
        }
        computeFinalGrade(studentInfo.value.id)

        editSubjects.value.forEach(subject => subject.grade = '')
        selectedGradingPeriod.value = ''
        closeEditGradesModal()
    }
}

// COMPUTES THE FINAL GRADE FOR A STUDENT
const computeFinalGrade = (studentId) => {
    const student = students.value.find(s => s.id === studentId)
    if (!student) return

    const periods = ['prelim', 'midterm', 'prefinals', 'finals']
    const weights = {
        prelim: 0.20,
        midterm: 0.20,
        prefinals: 0.20,
        finals: 0.40
    }

    const hasAllPeriods = periods.every(period => {
        const periodGrades = Object.values(student.grades[period])
        return periodGrades.some(grade => grade !== '')
    })

    if (hasAllPeriods) {
        let weightedSum = 0

        periods.forEach(period => {
            const periodGrades = Object.values(student.grades[period])
            const periodAverage = periodGrades
                .filter(grade => grade !== '')
                .reduce((sum, grade) => sum + Number(grade), 0) / periodGrades.length

            weightedSum += periodAverage * weights[period]

            student[period] = periodAverage.toFixed(2)
        })

        student.gwa = weightedSum.toFixed(2)

        // SETS REMARKS BASED ON GWA
        student.remarks = Number(student.gwa) === 0 ? 'Incomplete' :
            Number(student.gwa) <= 3.00 ? 'Passed' : 'Failed'
    } else {
        student.gwa = ''
        student.remarks = ''
    }
}

// HANDLES SEARCH QUERY CHANGES
const handleSearch = (query) => {
    searchQuery.value = query || ''
}

// HANDLES FILTER CHANGES
const handleFilterChange = (filters) => {
    activeFilters.value = filters || []
}

// SETS THE ACTIVE FILTER
const setFilter = (filter) => {
    activeFilter.value = filter === activeFilter.value ? 'all' : filter
}

// CHECKS IF A PERIOD HAS GRADES
const hasGradesForPeriod = (period) => {
    const student = students.value.find(s => s.id === studentInfo.value.id)
    if (!student) return false

    return Object.values(student.grades[period]).some(grade => grade)
}

// LOADS GRADES FOR A SELECTED PERIOD
const loadPeriodGrades = () => {
    const student = students.value.find(s => s.id === studentInfo.value.id)
    if (!student) return

    editSubjects.value = subjects.value.map(subject => ({
        ...subject,
        grade: student.grades[selectedGradingPeriod.value][subject.code] || ''
    }))
}

// HANDLES UNSAVED CHANGES MODAL ACTIONS
const handleUnsavedChanges = (confirm) => {
    if (confirm) {
        if (modalToClose.value === 'add') {
            closeAddGradesModal()
        } else if (modalToClose.value === 'edit') {
            closeEditGradesModal()
        }
    }
    isUnsavedChangesModalOpen.value = false
    modalToClose.value = null
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
    <main>
        <div class="sidebar">
            <Sidebar />
        </div>

        <section>
            <div class="container">
                <div class="welcome-header" style="padding-bottom: 0;">
                    <h1>Student Grades</h1>
                    <p>View and manage student grades</p>
                </div>

                <!-- SEARCH BAR SECTION -->
                <div class="students-controls">
                    <div class="search-filters">
                        <Searchbar v-model="searchQuery" @update:search-query="handleSearch"
                            @filter-change="handleFilterChange" />
                        <div class="filter-buttons">
                            <button v-for="status in ['All', 'Passed', 'Failed', 'Incomplete']" :key="status"
                                class="filter-btn" :class="{ active: activeFilter === status.toLowerCase() }"
                                @click="setFilter(status.toLowerCase())">
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
                                <th>Prelim</th>
                                <th>Midterm</th>
                                <th>Prefinals</th>
                                <th>Finals</th>
                                <th>GWA</th>
                                <th>Remarks</th>
                                <th>Actions</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="student in filteredStudents" :key="student.id">
                                <td>{{ student.id }}</td>
                                <td>{{ student.name }}</td>
                                <td>{{ student.course }}</td>
                                <td>{{ student.prelim || '-' }}</td>
                                <td>{{ student.midterm || '-' }}</td>
                                <td>{{ student.prefinals || '-' }}</td>
                                <td>{{ student.finals || '-' }}</td>
                                <td>{{ student.gwa || '-' }}</td>
                                <td>
                                    <span v-if="student.remarks" :class="[
                                        'status-badge',
                                        `status-${student.remarks.toLowerCase()}`
                                    ]">
                                        {{ student.remarks }}
                                    </span>
                                    <span v-else>-</span>
                                </td>
                                <td class="action-buttons">
                                    <button class="action-btn" @click="openAddGradesModal(student.id)">
                                        <i class="fa-solid fa-plus"></i>
                                    </button>
                                    <button class="action-btn edit-btn" @click="openEditGrades(student.id)">
                                        <i class="fa-solid fa-pen-to-square"></i>
                                    </button>
                                    <button class="action-btn delete-btn" @click="openDeleteGrades(student.id)">
                                        <i class="fa-solid fa-trash"></i>
                                    </button>
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>

            <!-- ADD STUDENT GRADES MODAL -->
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
                                <input type="text"
                                    :value="selectedGradingPeriod.charAt(0).toUpperCase() + selectedGradingPeriod.slice(1)"
                                    readonly required class="readonly-input">
                            </div>
                        </div>

                        <div class="subjects-section">
                            <div class="subjects-header">
                                <h3>Subject Grades</h3>
                                <div class="subject-actions" v-if="selectedGradingPeriod === 'prelim'">
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
                                        <div v-for="(subject, index) in subjects" :key="subject.code"
                                            class="subject-row" :class="{ 'right-column': index % 2 !== 0 }">
                                            <div class="grade-input">
                                                <label>{{ subject.name }}</label>
                                                <input type="text" required
                                                    placeholder="Enter grade (1.00, 1.25, 1.50, etc)"
                                                    v-model="subject.grade"
                                                    @input="(e) => handleGradeInput(subject, e.target.value)"
                                                    :class="{ 'error': subject.grade && !validGrades.includes(subject.grade) }">
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

            <!-- EDIT STUDENT GRADES MODAL -->
            <div class="modal" :class="{ active: isEditModalOpen }" @click="closeEditGradesModal">
                <div class="modal-content" @click.stop>
                    <div class="modal-header">
                        <h2>Edit Student Grades</h2>
                        <button class="close-modal" @click="closeEditGradesModal">&times;</button>
                    </div>
                    <form @submit="handleEditSubmit" class="student-form">
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
                                <select v-model="selectedGradingPeriod" required @change="loadPeriodGrades">
                                    <option value="" disabled selected>
                                        {{ editablePeriods.length > 0
                                            ? 'Select Period'
                                            : 'No grades added yet'
                                        }}
                                    </option>
                                    <option v-for="period in editablePeriods" :key="period" :value="period">
                                        {{ period.charAt(0).toUpperCase() + period.slice(1) }}
                                    </option>
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
                                                <input type="text" required
                                                    placeholder="Enter grade (1.00, 1.25, 1.50, etc)"
                                                    v-model="subject.grade"
                                                    @input="(e) => handleGradeInput(subject, e.target.value)"
                                                    :class="{ 'error': subject.grade && !validGrades.includes(subject.grade) }">
                                            </div>
                                        </div>
                                    </div>
                                </template>
                                <p v-else style="color: var(--text-secondary); text-align: center;">
                                    {{ editablePeriods.length > 0
                                        ? 'Please select a grading period'
                                        : 'No grades added yet'
                                    }}
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

            <!-- DELETE SUBJECT MODAL -->
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
                                        <div class="grade-details">
                                            <span class="total-grade">
                                                Grade: {{ subject.totalGrade }}
                                            </span>
                                        </div>
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

            <!-- DELETE CONFIRMATION MODAL -->
            <ConfirmationModal :is-open="isConfirmationModalOpen" :title="confirmationData.title"
                :item-name="confirmationData.itemName" :item-info="confirmationData.itemInfo"
                placeholder-text="Enter Subject Name to Confirm Deletion" @close="isConfirmationModalOpen = false"
                @confirm="handleConfirmDelete" />

            <UnsavedChangesModal :is-open="isUnsavedChangesModalOpen" @close="handleUnsavedChanges(false)"
                @confirm="handleUnsavedChanges(true)" />


            <Toast v-if="showToast" :message="toastMessage" type="error" :duration="null" class="modal-toast" />
        </section>
    </main>
</template>
