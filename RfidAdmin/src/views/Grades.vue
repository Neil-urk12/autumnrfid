<script setup lang="ts">
import { ref, computed, defineAsyncComponent } from 'vue'
import type { Subject, SubjectEditable, SubjectWithGrades, Student, ConfirmationData, StudentBasicInfo } from '@/typescript/models'
import mockData from '@/mock/models.json'
import Toast from '@/components/Toast.vue'

const ConfirmationModal = defineAsyncComponent(() => import('@/components/ConfirmationModal.vue'))
const UnsavedChangesModal = defineAsyncComponent(() => import('@/components/UnsavedChangesModal.vue'))
const Sidebar = defineAsyncComponent(() => import("@/components/Sidebar.vue"))
const Searchbar = defineAsyncComponent(() => import("@/components/Searchbar.vue"))


const students = ref<Student[]>(mockData.students)
const validGrades = ref<string[]>(mockData.validGrades)
const subjects = ref<SubjectEditable[]>(mockData.subjects.map(subject => ({ ...subject, grade: '' })))
const editSubjects = ref<SubjectEditable[]>(mockData.subjects.map(subject => ({ ...subject, grade: '' })))
const deleteSubjects = ref<SubjectWithGrades[]>([])
const availableSubjects = ref<Subject[]>(mockData.availableSubjects)

const isConfirmationModalOpen = ref<boolean>(false)
const isAddModalOpen = ref<boolean>(false)
const isEditModalOpen = ref<boolean>(false)
const isDeleteModalOpen = ref<boolean>(false)
const isUnsavedChangesModalOpen = ref<boolean>(false)
const modalToClose = ref<'add' | 'edit' | null>(null)

const confirmationData = ref<ConfirmationData>({
    title: '',
    itemName: '',
    itemInfo: null
})

const studentInfo = ref<StudentBasicInfo>({
    id: '',
    name: '',
    course: ''
})

const searchQuery = ref<string>('')
const activeFilters = ref<string[]>([])
const activeFilter = ref<string>('all')
const selectedGradingPeriod = ref<'prelim' | 'midterm' | 'prefinals' | 'finals' | ''>('')
const showNewSubjectForm = ref<boolean>(false)
const showToast = ref<boolean>(false)
const toastMessage = ref<string>('')
const selectedSubject = ref<string>('')

const newSubject = ref<SubjectEditable>({
    code: '',
    name: '',
    grade: ''
})

// HANDLES SEARCH AND FILTER
const filteredStudents = computed(() => {
    return students.value.filter(student => {
        const matchesSearch = !searchQuery.value ||
            student.id.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
            `${student.firstName} ${student.lastName}`.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
            student.course.toLowerCase().includes(searchQuery.value.toLowerCase())

        const matchesStatus = activeFilter.value === 'all' ||
            (activeFilter.value === 'passed' && student.grades.gwa && Number(student.grades.gwa) <= 3.00) ||
            (activeFilter.value === 'failed' && student.grades.gwa && Number(student.grades.gwa) > 3.00) ||
            (activeFilter.value === 'incomplete' && !student.grades.gwa)

        return matchesSearch && matchesStatus
    })
})

// GETS THE CURRENT GRADING PERIOD FOR A STUDENT AND PROCEEDS TO NEXT PERIOD
const currentPeriod = computed((): 'prelim' | 'midterm' | 'prefinals' | 'finals' | '' => {
    const student = students.value.find(s => s.id === studentInfo.value.id)
    if (!student) return ''

    const periods: ('prelim' | 'midterm' | 'prefinals' | 'finals')[] = ['prelim', 'midterm', 'prefinals', 'finals']

    const allPeriodsHaveGrades = periods.every(period =>
        Object.values(student.grades[period]).every(grade => grade !== '')
    )

    if (allPeriodsHaveGrades) {
        return ''
    }
    const nextPeriod = periods.find(period =>
        Object.values(student.grades[period]).every(grade => grade === '')
    )

    return nextPeriod || ''
})

// GETS THE LIST OF EDITABLE PERIODS FOR A STUDENT
const editablePeriods = computed(() => {
    const student = students.value.find(s => s.id === studentInfo.value.id)
    if (!student) return []

    const periods: ('prelim' | 'midterm' | 'prefinals' | 'finals')[] = ['prelim', 'midterm', 'prefinals', 'finals']
    return periods.filter(period => hasGradesForPeriod(period))
})

// CHECKS IF THERE ARE UNSAVED CHANGES IN THE MODAL
const hasUnsavedChanges = computed(() => {
    if (isAddModalOpen.value) {
        return subjects.value.some(subject => {
            const cleanGrade = subject.grade.trim()
            return cleanGrade !== '' && validGrades.value.includes(cleanGrade)
        })
    } else if (isEditModalOpen.value) {
        const student = students.value.find(s => s.id === studentInfo.value.id)
        if (!student || !selectedGradingPeriod.value) return false

        return editSubjects.value.some(subject => {
            const currentGrade = subject.grade.trim()
            const period = selectedGradingPeriod.value as keyof typeof student.grades
            const originalGrade = period ? (student.grades[period] as Record<string, string>)[subject.code] || '' : ''

            return currentGrade !== originalGrade &&
                (currentGrade === '' || validGrades.value.includes(currentGrade))
        })
    }
    return false
})

// SHOWS THE DELETE CONFIRMATION MODAL
const showDeleteConfirmation = (title: string, itemName: string, itemInfo: any = null) => {
    confirmationData.value = {
        title,
        itemName,
        itemInfo
    }
    isConfirmationModalOpen.value = true
}

// HANDLES THE CONFIRMATION OF DELETION
const handleConfirmDelete = (itemInfo: any) => {
    if (itemInfo?.name) {
        deleteSubjects.value = deleteSubjects.value.filter(
            subject => subject.name !== itemInfo.name
        )
    }
    isConfirmationModalOpen.value = false
}


// OPENS THE ADD GRADES MODAL
const openAddGradesModal = (studentId: string) => {
    const student = students.value.find(s => s.id === studentId)
    if (!student) return

    studentInfo.value = {
        id: studentId,
        name: `${student.firstName} ${student.lastName}`,
        course: student.course
    }

    subjects.value = mockData.subjects.map(subject => ({ ...subject, grade: '' }))
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
const openEditGrades = (studentId: string) => {
    const student = students.value.find(s => s.id === studentId)
    if (!student) return

    studentInfo.value = {
        id: student.id,
        name: `${student.firstName} ${student.lastName}`,
        course: student.course
    }

    selectedGradingPeriod.value = ''

    editSubjects.value = mockData.subjects.map(subject => ({
        ...subject,
        grade: ''
    }))

    isEditModalOpen.value = true
}

// OPENS THE DELETE GRADES MODAL
const openDeleteGrades = (studentId: string) => {
    const student = students.value.find(s => s.id === studentId)
    if (!student) return

    studentInfo.value = {
        id: student.id,
        name: `${student.firstName} ${student.lastName}`,
        course: student.course
    }

    deleteSubjects.value = mockData.subjects.map(subject => {
        const periodGrades = {
            prelim: student.grades.prelim[subject.code] || '',
            midterm: student.grades.midterm[subject.code] || '',
            prefinals: student.grades.prefinals[subject.code] || '',
            finals: student.grades.finals[subject.code] || ''
        }

        const currentGrade = calculateSubjectGrade(periodGrades)

        return {
            ...subject,
            grade: currentGrade,
            periodGrades
        }
    })

    isDeleteModalOpen.value = true
}

// CALCULATES THE AVERAGE GRADE FOR A SUBJECT ACROSS PERIODS
const calculateSubjectGrade = (periodGrades: Record<string, string>) => {
    const grades = Object.values(periodGrades)
        .filter(grade => grade !== '')
        .map(Number)

    if (grades.length === 0) return '--'

    const total = grades.reduce((sum, grade) => sum + grade, 0)
    return (total / grades.length).toFixed(2)
}


// CLOSES THE DELETE GRADES MODAL
const closeDeleteGradesModal = () => {
    isDeleteModalOpen.value = false
}

// ADDS A NEW SUBJECT
const addNewSubject = (subject: Subject | null = null) => {
    if (subject) {
        subjects.value.push({ ...subject, grade: '' })
        selectedSubject.value = ''
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

    const periods: ('prelim' | 'midterm' | 'prefinals' | 'finals')[] = ['prelim', 'midterm', 'prefinals', 'finals']
    periods.forEach(period => {
        Object.keys(student.grades[period]).forEach(code => {
            student.grades[period][code] = ''
        })
    })

    student.grades.gwa = ''
    student.grades.remarks = ''

    closeDeleteGradesModal()
}


// HANDLES GRADE INPUT VALIDATION
const validateGradeInput = (subject: SubjectEditable) => {
    const cleanValue = subject.grade.trim();

    if (cleanValue === '') {
        subject.grade = ''
        showToast.value = false;
        toastMessage.value = '';
        return;
    }

    const isValidGrade = validGrades.value.some(validGrade => {
        const inputNum = Number(cleanValue);
        const validNum = Number(validGrade);
        return inputNum === validNum;
    });

    if (!isValidGrade) {
        subject.grade = cleanValue
        toastMessage.value = 'Please enter a valid grade (1, 1.25, 1.50, etc)';
        showToast.value = true;
    } else {
        showToast.value = false;
        toastMessage.value = '';
    }
}

// HANDLES SUBMISSION OF ADD GRADES FORM
const handleSubmit = (e: Event) => {
    e.preventDefault()

    const student = students.value.find(s => s.id === studentInfo.value.id)
    if (!student || !selectedGradingPeriod.value) return

    const periodGrades = subjects.value.filter(subject => subject.grade.trim() !== '')

    if (periodGrades.length === 0) {
        toastMessage.value = 'Please enter at least one grade'
        showToast.value = true
        return
    }

    const hasInvalidGrades = periodGrades.some(subject => {
        const cleanGrade = subject.grade.trim()
        return !validGrades.value.some(validGrade => Number(cleanGrade) === Number(validGrade))
    })

    if (hasInvalidGrades) {
        toastMessage.value = 'Please enter valid grades for all subjects'
        showToast.value = true
        return
    }

    // UPDATE EACH SUBJECT GRADE
    periodGrades.forEach(subject => {
        // Convert whole numbers to decimal format when saving
        const grade = subject.grade.includes('.') ? subject.grade : `${subject.grade}.00`
        const period = selectedGradingPeriod.value as keyof typeof student.grades
        if (period) {
            (student.grades[period] as Record<string, string>)[subject.code] = grade
        }
    })

    // Compute final grade after updating grades
    computeFinalGrade(studentInfo.value.id)

    // Reset form
    subjects.value = mockData.subjects.map(subject => ({ ...subject, grade: '' }))
    selectedGradingPeriod.value = ''
    isAddModalOpen.value = false
    showToast.value = false
}

// HANDLES SUBMISSION OF EDIT GRADES FORM
const handleEditSubmit = (e: Event) => {
    e.preventDefault()

    // GET SUBJECT GRADES FOR SELECTED PERIOD
    const periodGrades = editSubjects.value.filter(subject => subject.grade)

    if (periodGrades.length > 0 && selectedGradingPeriod.value) {
        const hasInvalidGrades = periodGrades.some(subject => {
            const cleanGrade = subject.grade.trim()
            return !validGrades.value.some(validGrade => Number(cleanGrade) === Number(validGrade))
        })

        if (hasInvalidGrades) {
            toastMessage.value = 'Please enter valid grades for all subjects'
            showToast.value = true
            return
        }

        const student = students.value.find(s => s.id === studentInfo.value.id)
        if (student) {
            periodGrades.forEach(subject => {
                const grade = subject.grade.includes('.') ? subject.grade : `${subject.grade}.00`
                const period = selectedGradingPeriod.value as keyof typeof student.grades
                if (period) {
                    (student.grades[period] as Record<string, string>)[subject.code] = grade
                }
            })

            computeFinalGrade(studentInfo.value.id)
        }

        editSubjects.value.forEach(subject => subject.grade = '')
        selectedGradingPeriod.value = ''
        closeEditGradesModal()
    }
}

// LOADS GRADES FOR A SELECTED PERIOD
const loadPeriodGrades = () => {
    const student = students.value.find(s => s.id === studentInfo.value.id)
    if (!student || !selectedGradingPeriod.value) return

    editSubjects.value = mockData.subjects.map(subject => {
        const period = selectedGradingPeriod.value as keyof typeof student.grades
        const grade = period ? (student.grades[period] as Record<string, string>)[subject.code] || '' : ''
        return {
            ...subject,
            grade
        }
    })
}

// COMPUTES THE FINAL GRADE FOR A STUDENT
const computeFinalGrade = (studentId: string) => {
    const student = students.value.find(s => s.id === studentId)
    if (!student) return

    const periods: ('prelim' | 'midterm' | 'prefinals' | 'finals')[] = ['prelim', 'midterm', 'prefinals', 'finals']
    const weights = {
        prelim: 0.20,
        midterm: 0.20,
        prefinals: 0.20,
        finals: 0.40
    }

    const hasAllPeriods = periods.every(period =>
        Object.values(student.grades[period]).some(grade => grade !== '')
    )

    if (hasAllPeriods) {
        let weightedSum = 0
        let totalWeight = 0

        periods.forEach(period => {
            const periodGrades = Object.values(student.grades[period])
                .filter(grade => grade !== '')
                .map(Number)

            if (periodGrades.length > 0) {
                const periodAverage = periodGrades.reduce((sum, grade) => sum + grade, 0) / periodGrades.length
                weightedSum += periodAverage * weights[period]
                totalWeight += weights[period]
            }
        })

        if (totalWeight > 0) {
            const gwa = (weightedSum / totalWeight).toFixed(2)
            const remarks = Number(gwa) === 0 ? 'Incomplete' :
                Number(gwa) <= 3.00 ? 'Passed' : 'Failed'

            student.grades.gwa = gwa
            student.grades.remarks = remarks
        }
    }
}

// HANDLES SEARCH QUERY CHANGES
const handleSearch = (query: string) => {
    searchQuery.value = query || ''
}

// HANDLES FILTER CHANGES
const handleFilterChange = (filters: string[]) => {
    activeFilters.value = filters || []
}

// SETS THE ACTIVE FILTER
const setFilter = (filter: string) => {
    activeFilter.value = filter === activeFilter.value ? 'all' : filter
}

// CHECKS IF A PERIOD HAS GRADES
const hasGradesForPeriod = (period: 'prelim' | 'midterm' | 'prefinals' | 'finals') => {
    const student = students.value.find(s => s.id === studentInfo.value.id)
    if (!student) return false

    return Object.values(student.grades[period]).some(grade => grade !== '')
}

// HANDLES UNSAVED CHANGES MODAL ACTIONS
const handleUnsavedChanges = (confirm: boolean) => {
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

// CALCULATES THE AVERAGE GRADE FOR A PERIOD
const calculatePeriodAverage = (student: Student, period: 'prelim' | 'midterm' | 'prefinals' | 'finals') => {
    const grades = Object.values(student.grades[period])
        .filter(grade => grade !== '')
        .map(Number)

    if (grades.length === 0) return '-'

    const total = grades.reduce((sum, grade) => sum + grade, 0)
    return (total / grades.length).toFixed(2)
}

// HANDLES SUBJECT SELECTION FROM DROPDOWN
const handleSubjectSelect = () => {
    if (selectedSubject.value) {
        try {
            const subject = JSON.parse(selectedSubject.value);
            addNewSubject(subject);
        } catch (error) {
            console.error('Error parsing subject JSON:', error);
        }
    }
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
                                <td>{{ student.firstName }} {{ student.lastName }}</td>
                                <td>{{ student.course }}</td>
                                <td>{{ calculatePeriodAverage(student, 'prelim') }}</td>
                                <td>{{ calculatePeriodAverage(student, 'midterm') }}</td>
                                <td>{{ calculatePeriodAverage(student, 'prefinals') }}</td>
                                <td>{{ calculatePeriodAverage(student, 'finals') }}</td>
                                <td>{{ student.grades.gwa || '-' }}</td>
                                <td>
                                    <span v-if="student.grades.gwa" :class="[
                                        'status-badge',
                                        `status-${student.grades.remarks.toLowerCase()}`
                                    ]">
                                        {{ student.grades.remarks }}
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
                                <input type="text" :value="selectedGradingPeriod
                                    ? selectedGradingPeriod.charAt(0).toUpperCase() + selectedGradingPeriod.slice(1)
                                    : 'All periods have grades'" readonly required class="readonly-input"
                                    :class="{ 'no-periods': !selectedGradingPeriod }">
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
                                        <select @change="handleSubjectSelect"
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
                                                    @input="() => validateGradeInput(subject)"
                                                    :class="{ 'error': subject.grade && !validGrades.includes(subject.grade) }">
                                            </div>
                                        </div>
                                    </div>
                                </template>
                                <p v-else style="color: var(--text-secondary); text-align: center;">
                                    All grading periods have been completed for this student
                                </p>
                            </div>
                        </div>

                        <div class="form-actions">
                            <button type="button" class="cancel-btn" @click="closeAddGradesModal">Cancel</button>
                            <button type="submit" class="submit-btn" :disabled="!selectedGradingPeriod">Add
                                Grades</button>
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
                                <div class="subject-actions" v-if="selectedGradingPeriod === 'prelim'">
                                    <button type="button" class="add-subject-btn" @click="addNewSubject()">
                                        <i class="fa-solid fa-plus"></i> Add Custom Subject
                                    </button>
                                    <div class="subject-select">
                                        <select @change="handleSubjectSelect"
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
                                                    @input="() => validateGradeInput(subject)"
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
                                                Grade: {{ subject.grade }}
                                            </span>
                                            <div class="period-grades" v-if="subject.periodGrades">
                                                <span v-if="subject.periodGrades.prelim" class="period-grade">
                                                    Prelim: {{ subject.periodGrades.prelim }}
                                                </span>
                                                <span v-if="subject.periodGrades.midterm" class="period-grade">
                                                    Midterm: {{ subject.periodGrades.midterm }}
                                                </span>
                                                <span v-if="subject.periodGrades.prefinals" class="period-grade">
                                                    Prefinals: {{ subject.periodGrades.prefinals }}
                                                </span>
                                                <span v-if="subject.periodGrades.finals" class="period-grade">
                                                    Finals: {{ subject.periodGrades.finals }}
                                                </span>
                                            </div>
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

            <!-- UNSAVED CHANGES WARNING -->
            <UnsavedChangesModal :is-open="isUnsavedChangesModalOpen" @close="handleUnsavedChanges(false)"
                @confirm="handleUnsavedChanges(true)" />


            <Toast v-if="showToast" :message="toastMessage" type="error" :duration="null" class="modal-toast" />
        </section>
    </main>
</template>

