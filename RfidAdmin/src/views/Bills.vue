<script setup>
import { ref, computed, defineAsyncComponent } from 'vue'
const Sidebar = defineAsyncComponent(() => import("@/components/Sidebar.vue"))
const Searchbar = defineAsyncComponent(() => import("@/components/Searchbar.vue"))

const studentBills = ref([
    {
        id: '2023-0001',
        name: 'John Cez',
        course: 'Computer Science',
        yearLevel: '3rd Year',
        basicTuition: 25000.00,
        laboratoryFee: 5000.00,
        studentActivities: 2000.00,
        discount: 2500.00,
        totalFees: 29500.00,
        examFees: {
            prelim: 3000.00,
            midterm: 3000.00,
            prefinal: 3000.00,
            final: 3000.00
        },
        payments: [],
        remainingBalance: 29500.00
    },
    {
        id: '2023-0002',
        name: 'Maria Santos',
        course: 'Business Administration',
        yearLevel: '2nd Year',
        basicTuition: 22000.00,
        laboratoryFee: 2000.00,
        studentActivities: 1500.00,
        discount: 0.00,
        totalFees: 25500.00,
        examFees: {
            prelim: 2000.00,
            midterm: 2000.00,
            prefinal: 2000.00,
            final: 2000.00
        },
        payments: [],
        remainingBalance: 25500.00
    },
    {
        id: '2023-0003',
        name: 'Carlos Reyes',
        course: 'Information Technology',
        yearLevel: '1st Year',
        basicTuition: 20000.00,
        laboratoryFee: 4000.00,
        studentActivities: 1800.00,
        discount: 1000.00,
        totalFees: 24800.00,
        examFees: {
            prelim: 2500.00,
            midterm: 2500.00,
            prefinal: 2500.00,
            final: 2500.00
        },
        payments: [],
        remainingBalance: 24800.00
    }
])

const searchQuery = ref('')
const activeFilters = ref([])
const showStudentModal = ref(false)
const showAddStudentModal = ref(false)
const showEditFeeStructureModal = ref(false)
const activeTab = ref('info')
const showPaymentForm = ref(false)
const showEditModal = ref(false)
const showViewModal = ref(false)

const currentStudent = ref({
    id: '',
    name: '',
    course: '',
    yearLevel: '',
    basicTuition: 0,
    laboratoryFee: 0,
    studentActivities: 0,
    discount: 0,
    totalFees: 0,
    examFees: {
        prelim: 0,
        midterm: 0,
        prefinal: 0,
        final: 0
    },
    payments: [],
    remainingBalance: 0
})

const viewStudent = ref({
    id: '',
    name: '',
    course: '',
    yearLevel: '',
    basicTuition: 0,
    laboratoryFee: 0,
    studentActivities: 0,
    discount: 0,
    totalFees: 0,
    examFees: {
        prelim: 0,
        midterm: 0,
        prefinal: 0,
        final: 0
    },
    payments: [],
    remainingBalance: 0
})

const editStudentData = ref({
    id: '',
    name: '',
    course: '',
    yearLevel: ''
})

const newPayment = ref({
    type: 'Initial Payment',
    amount: '',
    date: ''
})

const newStudent = ref({
    id: '',
    name: '',
    course: '',
    yearLevel: '',
    feePreset: 'default',
    basicTuition: 0,
    laboratoryFee: 0,
    studentActivities: 0,
    discount: 0,
    discountType: 'none'
})

const feeStructures = ref({
    default: {
        tuition: {
            registration: 5000,
            subject: 15000
        },
        misc: {
            laboratory: 3000,
            library: 1000,
            computer: 2000,
            athletic: 1000
        }
    },
    computer: {
        tuition: {
            registration: 6000,
            subject: 19000
        },
        misc: {
            laboratory: 5000,
            library: 1500,
            computer: 3500,
            athletic: 1000
        }
    },
    business: {
        tuition: {
            registration: 5500,
            subject: 16500
        },
        misc: {
            laboratory: 2000,
            library: 1200,
            computer: 1800,
            athletic: 1000
        }
    }
})

const editFeePreset = ref('default')

const handleSearch = (query) => {
  searchQuery.value = query || ''
}

const handleFilterChange = (filters) => {
  activeFilters.value = filters || []
}
const filteredStudents = computed(() => {
  return studentBills.value.filter(student => {
    const matchesSearch = !searchQuery.value || 
      student.id.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      student.name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      student.course.toLowerCase().includes(searchQuery.value.toLowerCase())

    const matchesFilters = !activeFilters.value.length ||
      activeFilters.value.some(filter => 
        student.yearLevel.includes(filter) ||
        student.course.includes(filter))

    return matchesSearch && matchesFilters
  })
})

const openStudentModal = (studentId) => {
    const student = studentBills.value.find(s => s.id === studentId)
    if (!student) return

    currentStudent.value = JSON.parse(JSON.stringify(student))
    showStudentModal.value = true
    activeTab.value = 'info'
    showPaymentForm.value = false
    showEditModal.value = false
}

const closeStudentModal = () => {
    showStudentModal.value = false
    currentStudent.value = {
        id: '',
        name: '',
        course: '',
        yearLevel: '',
        basicTuition: 0,
        laboratoryFee: 0,
        studentActivities: 0,
        discount: 0,
        totalFees: 0,
        examFees: {
            prelim: 0,
            midterm: 0,
            prefinal: 0,
            final: 0
        },
        payments: [],
        remainingBalance: 0
    }
}

const openAddStudentModal = () => {
    newStudent.value = {
        id: '',
        name: '',
        course: '',
        yearLevel: '',
        feePreset: 'default',
        basicTuition: 0,
        laboratoryFee: 0,
        studentActivities: 0,
        discount: 0,
        discountType: 'none'
    }
    applyCoursePreset('default')
    showAddStudentModal.value = true
}

const closeAddStudentModal = () => {
    showAddStudentModal.value = false
}

const saveNewStudent = () => {
    if (!newStudent.value.id || !newStudent.value.name ||
        !newStudent.value.course || !newStudent.value.yearLevel) {
        alert('Please fill in all required fields')
        return
    }

    if (studentBills.value.some(s => s.id === newStudent.value.id)) {
        alert('Student ID already exists')
        return
    }

    const totalFees = newStudent.value.basicTuition +
        newStudent.value.laboratoryFee +
        newStudent.value.studentActivities -
        newStudent.value.discount

    const examFee = newStudent.value.basicTuition * 0.2

    const student = {
        id: newStudent.value.id,
        name: newStudent.value.name,
        course: newStudent.value.course,
        yearLevel: newStudent.value.yearLevel,
        basicTuition: newStudent.value.basicTuition,
        laboratoryFee: newStudent.value.laboratoryFee,
        studentActivities: newStudent.value.studentActivities,
        discount: newStudent.value.discount,
        totalFees,
        examFees: {
            prelim: examFee,
            midterm: examFee,
            prefinal: examFee,
            final: examFee
        },
        payments: [],
        remainingBalance: totalFees
    }

    studentBills.value.push(student)
    closeAddStudentModal()
}

const editStudent = (studentId) => {
    const student = studentBills.value.find(s => s.id === studentId)
    if (!student) return

    openStudentModal(studentId)

    editStudentData.value = {
        id: student.id,
        name: student.name,
        course: student.course,
        yearLevel: student.yearLevel
    }

    showEditModal.value = true
}

const closeEditModal = () => {
    showEditModal.value = false
    currentStudent.value = {
        id: '',
        name: '',
        course: '',
        yearLevel: '',
        basicTuition: 0,
        laboratoryFee: 0,
        studentActivities: 0,
        discount: 0,
        totalFees: 0,
        examFees: {
            prelim: 0,
            midterm: 0,
            prefinal: 0,
            final: 0
        },
        payments: [],
        remainingBalance: 0
    }
}

const applyCoursePreset = (preset) => {
    let tuition = 0
    let miscFee = 0
    let activities = 0

    switch (preset) {
        case 'computer':
            tuition = 25000
            miscFee = 5000
            activities = 2000
            break
        case 'business':
            tuition = 22000
            miscFee = 2000
            activities = 1500
            break
        default:
            tuition = 20000
            miscFee = 3000
            activities = 1800
    }

    newStudent.value.feePreset = preset
    newStudent.value.basicTuition = tuition
    newStudent.value.laboratoryFee = miscFee
    newStudent.value.studentActivities = activities
    newStudent.value.discount = 0
}

const applyYearPreset = (year) => {
    let tuition = newStudent.value.basicTuition || 0
    let miscFee = newStudent.value.laboratoryFee || 0

    switch (year) {
        case '1st':
            tuition = tuition * 1.0
            miscFee = miscFee * 1.0
            newStudent.value.yearLevel = '1st Year'
            break
        case '2nd':
            tuition = tuition * 1.05
            miscFee = miscFee * 1.1
            newStudent.value.yearLevel = '2nd Year'
            break
        case '3rd':
            tuition = tuition * 1.1
            miscFee = miscFee * 1.2
            newStudent.value.yearLevel = '3rd Year'
            break
        case '4th':
            tuition = tuition * 1.15
            miscFee = miscFee * 1.3
            newStudent.value.yearLevel = '4th Year'
            break
    }

    newStudent.value.basicTuition = Math.round(tuition)
    newStudent.value.laboratoryFee = Math.round(miscFee)
    newStudent.value.discountType = 'none'
    newStudent.value.discount = 0
}

const applyDiscountType = () => {
    const totalBeforeDiscount = newStudent.value.basicTuition +
        newStudent.value.laboratoryFee +
        newStudent.value.studentActivities

    switch (newStudent.value.discountType) {
        case 'honor':
            newStudent.value.discount = Math.round(totalBeforeDiscount * 0.15)
            break
        case 'highHonor':
            newStudent.value.discount = Math.round(totalBeforeDiscount * 0.30)
            break
        case 'highestHonor':
            newStudent.value.discount = Math.round(totalBeforeDiscount * 0.50)
            break
        case 'freshman':
        case 'continuing':
            newStudent.value.discount = Math.round(totalBeforeDiscount * 0.10)
            break
        default:
            newStudent.value.discount = 0
    }
}

const loadFeeStructure = (preset) => {
    editFeePreset.value = preset
}

const saveFeeStructure = () => {
    if (newStudent.value.feePreset === editFeePreset.value) {
        const tuition = feeStructures.value[editFeePreset.value].tuition.registration +
            feeStructures.value[editFeePreset.value].tuition.subject

        const miscFee = feeStructures.value[editFeePreset.value].misc.laboratory +
            feeStructures.value[editFeePreset.value].misc.library +
            feeStructures.value[editFeePreset.value].misc.computer +
            feeStructures.value[editFeePreset.value].misc.athletic

        newStudent.value.basicTuition = tuition
        newStudent.value.laboratoryFee = miscFee
    }

    showEditFeeStructureModal.value = false
}

const startEditingPayment = (payment) => {
    currentStudent.value.payments.forEach(p => {
        if (p.id !== payment.id) {
            p.isEditing = false;
        }
    });

    payment.isEditing = true;
    payment.editAmount = payment.amount;
}

const cancelEditingPayment = (payment) => {
    payment.isEditing = false;
    delete payment.editAmount;
}

const saveEditedPayment = (payment) => {
    if (!payment.editAmount || isNaN(payment.editAmount)) {
        alert('Please enter a valid amount');
        return;
    }

    const newAmount = parseFloat(payment.editAmount);
    if (newAmount <= 0) {
        alert('Amount must be greater than zero');
        return;
    }

    const amountDifference = payment.amount - newAmount;

    const studentIndex = studentBills.value.findIndex(s => s.id === currentStudent.value.id);
    if (studentIndex === -1) return;

    const paymentIndex = studentBills.value[studentIndex].payments.findIndex(p => p.id === payment.id);
    if (paymentIndex === -1) return;

    studentBills.value[studentIndex].payments[paymentIndex].amount = newAmount;

    studentBills.value[studentIndex].remainingBalance += amountDifference;
    if (studentBills.value[studentIndex].remainingBalance < 0) {
        studentBills.value[studentIndex].remainingBalance = 0;
    }

    currentStudent.value.payments[paymentIndex].amount = newAmount;
    currentStudent.value.remainingBalance = studentBills.value[studentIndex].remainingBalance;

    payment.isEditing = false;
    delete payment.editAmount;
}

const saveStudentChanges = () => {
    const index = studentBills.value.findIndex(s => s.id === currentStudent.value.id);
    if (index === -1) return;

    studentBills.value[index] = JSON.parse(JSON.stringify(currentStudent.value));
    closeEditModal();
}

const updatePaymentAmount = () => {
    if (newPayment.value.type === "Full Payment") {
        newPayment.value.amount = currentStudent.value.remainingBalance
    } else if (newPayment.value.type === "Full Payment (before exam)") {
        newPayment.value.amount = Math.max(0, currentStudent.value.remainingBalance - 1200)
    } else {
        newPayment.value.amount = ''
    }
}

const submitPayment = () => {
    if (!currentStudent.value.id) return

    const type = newPayment.value.type
    let amount = parseFloat(newPayment.value.amount)
    const date = newPayment.value.date

    if (!amount || isNaN(amount) || !date || !type) {
        alert('Please fill in all payment details')
        return
    }

    if (type !== "Full Payment (before exam)" && amount > currentStudent.value.remainingBalance) {
        alert('Payment amount cannot exceed remaining balance')
        return
    }

    const paymentId = 'pay' + Date.now()
    const payment = {
        id: paymentId,
        type,
        amount,
        date,
        status: 'Pending'
    }

    const index = studentBills.value.findIndex(s => s.id === currentStudent.value.id)
    if (index === -1) return

    studentBills.value[index].payments = [...studentBills.value[index].payments, payment]
    currentStudent.value.payments = [...currentStudent.value.payments, payment]

    if (type === "Full Payment (before exam)") {
        studentBills.value[index].totalFees -= 1200
        currentStudent.value.totalFees -= 1200
        studentBills.value[index].remainingBalance = 0
        currentStudent.value.remainingBalance = 0
    } else {
        studentBills.value[index].remainingBalance -= amount
        if (studentBills.value[index].remainingBalance < 0) {
            studentBills.value[index].remainingBalance = 0
        }
        currentStudent.value.remainingBalance = studentBills.value[index].remainingBalance
    }

    newPayment.value = {
        type: 'Initial Payment',
        amount: '',
        date: ''
    }
    showPaymentForm.value = false
}

const openViewModal = (studentId) => {
    const student = studentBills.value.find(s => s.id === studentId)
    if (!student) return

    let feePreset = 'default'
    if (student.course && (student.course.includes('Computer') || student.course.includes('Information'))) {
        feePreset = 'computer'
    } else if (student.course && student.course.includes('Business')) {
        feePreset = 'business'
    }

    viewStudent.value = JSON.parse(JSON.stringify(student))
    viewStudent.value.feePreset = feePreset
    showViewModal.value = true
}

const closeViewModal = () => {
    showViewModal.value = false
    viewStudent.value = {
        id: '',
        name: '',
        course: '',
        yearLevel: '',
        basicTuition: 0,
        laboratoryFee: 0,
        studentActivities: 0,
        discount: 0,
        totalFees: 0,
        examFees: {
            prelim: 0,
            midterm: 0,
            prefinal: 0,
            final: 0
        },
        payments: [],
        remainingBalance: 0
    }
}

const getFeeStructure = (preset) => {
    return feeStructures.value[preset] || feeStructures.value['default']
}
</script>


<template>
    <main>
        <div class="sidebar">
            <Sidebar />
        </div>
        <section>
            <div class="container">
                <div class="welcome-header">
                    <h1>Student Bills</h1>
                    <p>Manage and monitor student payments and balances</p>
                </div>

                <div class="students-controls">
                    <div class="search-filters">
                        <Searchbar 
                            v-model="searchQuery"
                            @update:search-query="handleSearch"
                            @filter-change="handleFilterChange"
                        />
                        <div class="filter-buttons">
                            <button class="add-student-btn" @click="openAddStudentModal">
                                <i class="fa-solid fa-plus"></i> Add New Student Bill
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
                                <th>Total Tuition</th>
                                <th>Paid Amount</th>
                                <th>Remaining Balance</th>
                                <th>Action</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="student in filteredStudents" :key="student.id">
                                <td>{{ student.id }}</td>
                                <td>{{ student.name }}</td>
                                <td>{{ student.course }}</td>
                                <td>₱{{ student.totalFees.toLocaleString() }}</td>
                                <td>₱{{ (student.totalFees - student.remainingBalance).toLocaleString() }}</td>
                                <td>₱{{ student.remainingBalance.toLocaleString() }}</td>
                                <td class="action-buttons">
                                    <button class="action-btn view" @click="openViewModal(student.id)">
                                        <i class="fas fa-eye"></i>
                                    </button>
                                    <button class="action-btn edit-btn" @click="editStudent(student.id)">
                                        <i class="fa-solid fa-pen-to-square"></i>
                                    </button>

                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>

            <!-- ADD NEW STUDENT BILL MODAL -->
            <div class="modal" :class="{ active: showAddStudentModal }">
                <div class="modal-content">
                    <div class="modal-header">
                        <h2>Add New Student Bill</h2>
                        <button class="close-modal" @click="closeAddStudentModal">&times;</button>
                    </div>

                    <form class="student-form">
                        <div class="student-info-grid">
                            <div class="info-group">
                                <label>Student ID</label>
                                <input type="text" v-model="newStudent.id" placeholder="e.g. 2023-0001">
                            </div>
                            <div class="info-group">
                                <label>Name</label>
                                <input type="text" v-model="newStudent.name" placeholder="Full name">
                            </div>
                            <div class="info-group">
                                <label>Course</label>
                                <select v-model="newStudent.course">
                                    <option value="">Select course</option>
                                    <option value="Computer Science">Computer Science</option>
                                    <option value="Information Technology">Information Technology</option>
                                    <option value="Business Administration">Business Administration</option>
                                    <option value="Accountancy">Accountancy</option>
                                </select>
                            </div>
                            <div class="info-group">
                                <label>Year Level</label>
                                <select v-model="newStudent.yearLevel">
                                    <option value="">Select year level</option>
                                    <option value="1st Year">1st Year</option>
                                    <option value="2nd Year">2nd Year</option>
                                    <option value="3rd Year">3rd Year</option>
                                    <option value="4th Year">4th Year</option>
                                </select>
                            </div>
                        </div>

                        <div class="course-selector">
                            <div class="course-option" :class="{ active: newStudent.feePreset === 'default' }"
                                @click="applyCoursePreset('default')">
                                Default Fees
                            </div>
                            <div class="course-option" :class="{ active: newStudent.feePreset === 'computer' }"
                                @click="applyCoursePreset('computer')">
                                Computer Courses
                            </div>
                            <div class="course-option" :class="{ active: newStudent.feePreset === 'business' }"
                                @click="applyCoursePreset('business')">
                                Business Courses
                            </div>
                        </div>

                        <div class="fee-presets">
                            <div class="fee-presets-header">
                                <h4>Fee Details</h4>
                                <button class="action-button" @click.prevent="showEditFeeStructureModal = true">
                                    <i class="fas fa-cog"></i> Edit Fee Details
                                </button>
                            </div>

                            <div class="student-info-grid">
                                <div class="info-group">
                                    <label>Tuition Fee</label>
                                    <input type="number" v-model="newStudent.basicTuition">
                                </div>
                                <div class="info-group">
                                    <label>Miscellaneous Fee</label>
                                    <input type="number" v-model="newStudent.laboratoryFee">
                                </div>
                                <div class="info-group">
                                    <label>Initial Payment</label>
                                    <input type="number" v-model="newStudent.studentActivities">
                                </div>
                                <div class="info-group">
                                    <label>Discount</label>
                                    <select v-model="newStudent.discountType" @change="applyDiscountType">
                                        <option value="none">No Discount</option>
                                        <template v-if="newStudent.yearLevel === '1st Year'">
                                            <option value="honor">Honor Student (15%)</option>
                                            <option value="highHonor">High Honor (30%)</option>
                                            <option value="highestHonor">Highest Honor (50%)</option>
                                            <option value="freshman">Freshman (10%)</option>
                                        </template>
                                        <template v-else-if="newStudent.yearLevel">
                                            <option value="continuing">Continuing Student (10%)</option>
                                        </template>
                                    </select>
                                </div>
                            </div>

                            <div class="preset-buttons">
                                <button type="button" class="preset-btn" @click.prevent="applyYearPreset('1st')">
                                    1st Year
                                </button>
                                <button type="button" class="preset-btn" @click.prevent="applyYearPreset('2nd')">
                                    2nd Year
                                </button>
                                <button type="button" class="preset-btn" @click.prevent="applyYearPreset('3rd')">
                                    3rd Year
                                </button>
                                <button type="button" class="preset-btn" @click.prevent="applyYearPreset('4th')">
                                    4th Year
                                </button>
                            </div>
                        </div>

                        <div class="exam-fees">
                            <h3 class="labels">Exam Fees</h3>
                            <div class="exam-fees-grid">
                                <div class="exam-fee-item">
                                    <span class="exam-fee-label">Prelim Exam</span>
                                    <span class="exam-fee-value">₱{{ (newStudent.basicTuition * 0.2).toLocaleString()
                                        }}</span>
                                </div>
                                <div class="exam-fee-item">
                                    <span class="exam-fee-label">Midterm Exam</span>
                                    <span class="exam-fee-value">₱{{ (newStudent.basicTuition * 0.2).toLocaleString()
                                        }}</span>
                                </div>
                                <div class="exam-fee-item">
                                    <span class="exam-fee-label">Pre-final Exam</span>
                                    <span class="exam-fee-value">₱{{ (newStudent.basicTuition * 0.2).toLocaleString()
                                        }}</span>
                                </div>
                                <div class="exam-fee-item">
                                    <span class="exam-fee-label">Final Exam</span>
                                    <span class="exam-fee-value">₱{{ (newStudent.basicTuition * 0.2).toLocaleString()
                                        }}</span>
                                </div>
                            </div>
                        </div>

                        <div class="summary-cards">
                            <div class="summary-card">
                                <h4>Total Tuition</h4>
                                <p>₱{{ (newStudent.basicTuition + newStudent.laboratoryFee +
                                    newStudent.studentActivities - newStudent.discount).toLocaleString() }}</p>
                            </div>
                            <div class="summary-card">
                                <h4>Total Paid</h4>
                                <p>₱0.00</p>
                            </div>
                            <div class="summary-card highlight">
                                <h4>Remaining Balance</h4>
                                <p>₱{{ (newStudent.basicTuition + newStudent.laboratoryFee +
                                    newStudent.studentActivities - newStudent.discount).toLocaleString() }}</p>
                            </div>
                        </div>
                    </form>

                    <div class="form-actions">
                        <button class="submit-btn" @click="saveNewStudent">Add Student</button>
                        <button class="cancel-btn" @click="closeAddStudentModal">Cancel</button>
                    </div>
                </div>
            </div>

            <!-- EDIT FEE DETAILS MODAL -->
            <div class="modal" :class="{ active: showEditFeeStructureModal }">
                <div class="modal-content">
                    <div class="modal-header">
                        <h2>Edit Fee Structure</h2>
                        <button class="close-modal" @click="showEditFeeStructureModal = false">&times;</button>
                    </div>
                    <div class="modal-body">
                        <div class="course-selector">
                            <div class="course-option" :class="{ active: editFeePreset === 'default' }"
                                @click="loadFeeStructure('default')">
                                Default Course
                            </div>
                            <div class="course-option" :class="{ active: editFeePreset === 'computer' }"
                                @click="loadFeeStructure('computer')">
                                Computer Courses
                            </div>
                            <div class="course-option" :class="{ active: editFeePreset === 'business' }"
                                @click="loadFeeStructure('business')">
                                Business Courses
                            </div>
                        </div>


                        <form class="student-form">
                            <h3 class="labels">Tuition Fee</h3>
                            <div class="student-info-grid">
                                <div class="info-group">
                                    <label>Basic Tuition Fee</label>
                                    <input type="number" v-model="feeStructures[editFeePreset].tuition.registration">
                                </div>
                                <div class="info-group">
                                    <label>Laboratory Fee</label>
                                    <input type="number" v-model="feeStructures[editFeePreset].tuition.subject">
                                </div>
                            </div>

                            <h3 class="labels">Miscellaneous Fee</h3>
                            <div class="student-info-grid">
                                <div class="info-group">
                                    <label>Development Fee</label>
                                    <input type="number" v-model="feeStructures[editFeePreset].misc.laboratory">
                                </div>
                                <div class="info-group">
                                    <label>Library Fee</label>
                                    <input type="number" v-model="feeStructures[editFeePreset].misc.library">
                                </div>
                                <div class="info-group">
                                    <label>Computer Laboratory Fee</label>
                                    <input type="number" v-model="feeStructures[editFeePreset].misc.computer">
                                </div>
                                <div class="info-group">
                                    <label>Athletic Fee</label>
                                    <input type="number" v-model="feeStructures[editFeePreset].misc.athletic">
                                </div>
                            </div>
                        </form>
                    </div>
                    <div class="form-actions">
                        <button class="submit-btn" @click="saveFeeStructure">Save Changes</button>
                        <button class="cancel-btn" @click="showEditFeeStructureModal = false">Cancel</button>
                    </div>
                </div>
            </div>

            <!-- EDIT STUDENT'S BILL MODAL -->
            <div class="modal" :class="{ active: showStudentModal }">
                <div class="modal-content">
                    <div class="modal-header">
                        <h2>Edit Billing Details</h2>
                        <button class="close-modal" @click="closeStudentModal">&times;</button>
                    </div>


                    <form class="student-form">
                        <div class="student-info-grid">
                            <div class="info-group">
                                <label>Name</label>
                                <input type="text" v-model="currentStudent.name" readonly disabled>
                            </div>
                            <div class="info-group">
                                <label>Year Level</label>
                                <input type="text" v-model="currentStudent.yearLevel" readonly disabled>
                            </div>
                            <div class="info-group">
                                <label>Basic Tuition</label>
                                <input type="number" v-model="currentStudent.basicTuition">
                            </div>
                            <div class="info-group">
                                <label>Laboratory Fee</label>
                                <input type="number" v-model="currentStudent.laboratoryFee">
                            </div>
                            <div class="info-group">
                                <label>Initial Payment</label>
                                <input type="number" v-model="currentStudent.studentActivities">
                            </div>
                            <div class="info-group">
                                <label>Discount</label>
                                <select v-model="currentStudent.discount" @change="applyDiscount">
                                    <template v-if="currentStudent.yearLevel === '1st Year'">
                                        <option :value="0">No Discount</option>
                                        <option :value="currentStudent.basicTuition * 0.15">Honor Student (15%)</option>
                                        <option :value="currentStudent.basicTuition * 0.30">High Honor (30%)</option>
                                        <option :value="currentStudent.basicTuition * 0.50">Highest Honor (50%)</option>
                                        <option :value="currentStudent.basicTuition * 0.10">Freshman (10%)</option>
                                    </template>
                                    <template v-else>
                                        <option :value="0">No Discount</option>
                                        <option :value="currentStudent.basicTuition * 0.10">Continuing Student (10%)
                                        </option>
                                    </template>
                                </select>
                            </div>
                        </div>

                        <div class="exam-fees">
                            <h3 class="labels">Exam Fees</h3>
                            <div class="exam-fees-grid">
                                <div class="exam-fee-item">
                                    <span class="exam-fee-label">Prelim Exam</span>
                                    <span class="exam-fee-value">₱{{
                                        currentStudent.examFees.prelim.toLocaleString() }}</span>
                                </div>
                                <div class="exam-fee-item">
                                    <span class="exam-fee-label">Midterm Exam</span>
                                    <span class="exam-fee-value">₱{{
                                        currentStudent.examFees.midterm.toLocaleString() }}</span>
                                </div>
                                <div class="exam-fee-item">
                                    <span class="exam-fee-label">Pre-final Exam</span>
                                    <span class="exam-fee-value">₱{{
                                        currentStudent.examFees.prefinal.toLocaleString() }}</span>
                                </div>
                                <div class="exam-fee-item">
                                    <span class="exam-fee-label">Final Exam</span>
                                    <span class="exam-fee-value">₱{{
                                        currentStudent.examFees.final.toLocaleString() }}</span>
                                </div>
                            </div>
                        </div>

                        <div class="payment-actions">
                            <h3 class="labels">Payments</h3>
                            <div>
                                <button class="action-button" @click="showPaymentForm = true" v-if="!showPaymentForm">
                                    <i class="fas fa-plus"></i> Add Payment
                                </button>
                            </div>
                        </div>

                        <div class="payment-form" :class="{ active: showPaymentForm }">
                            <div class="payment-form-grid">
                                <div class="info-group">
                                    <label>Payment Type</label>
                                    <select v-model="newPayment.type" @change="updatePaymentAmount">
                                        <option value="Initial Payment">Initial Payment</option>
                                        <option value="Prelim">Prelim</option>
                                        <option value="Midterm">Midterm</option>
                                        <option value="Pre-final">Pre-final</option>
                                        <option value="Final">Final</option>
                                        <option value="Full Payment">Full Payment</option>
                                        <option value="Full Payment (before exam)">Full Payment (before exam)</option>
                                        <option value="Other">Other</option>
                                    </select>
                                </div>
                                <div class="info-group">
                                    <label>Amount</label>
                                    <input type="number" v-model="newPayment.amount" placeholder="Enter amount">
                                </div>
                                <div class="info-group">
                                    <label>Date</label>
                                    <input type="date" v-model="newPayment.date">
                                </div>
                            </div>
                            <div class="payment-form-buttons">
                                <button type="button" class="payment-form-btn payment-form-cancel"
                                    @click="showPaymentForm = false">
                                    Cancel
                                </button>
                                <button type="button" class="payment-form-btn payment-form-submit"
                                    @click="submitPayment" :disabled="isPaymentButtonDisabled"
                                    :class="{ 'disabled': isPaymentButtonDisabled }">
                                    Add Payment
                                </button>
                            </div>
                        </div>
 
                        <table class="payment-table"
                            v-if="currentStudent.payments && currentStudent.payments.length > 0">
                            <thead>
                                <tr>
                                    <th>Payment Type</th>
                                    <th>Date</th>
                                    <th>Amount</th>
                                    <th>Actions</th>
                                </tr>
                            </thead>
                            <tbody>
                                <tr v-for="payment in currentStudent.payments" :key="payment.id">
                                    <td>{{ payment.type }}</td>
                                    <td>{{ payment.date }}</td>
                                    <td>₱{{ payment.amount.toLocaleString() }}</td>
                                    <td>
                                        <div class="payment-actions-row">
                                            <button class="payment-action-btn payment-edit" v-if="!payment.isEditing"
                                                @click="startEditingPayment(payment)">
                                                <i class="fa-solid fa-pen-to-square"></i>
                                            </button>
                                            <div v-else class="payment-edit-form">
                                                <input type="number" v-model="payment.editAmount"
                                                    :max="currentStudent.remainingBalance + payment.amount" min="1">
                                                <button class="payment-action-btn payment-save"
                                                    @click="saveEditedPayment(payment)">
                                                    <i class="fas fa-check"></i>
                                                </button>
                                                <button class="payment-action-btn payment-cancel"
                                                    @click="cancelEditingPayment(payment)">
                                                    <i class="fas fa-times"></i>
                                                </button>
                                            </div>
                                            <button class="payment-action-btn payment-delete"
                                                @click="handlePaymentAction(payment.id, 'delete')">
                                                <i class="fas fa-trash"></i>
                                            </button>
                                        </div>
                                    </td>
                                </tr>
                            </tbody>
                        </table>

                        <div class="summary-cards">
                            <div class="summary-card">
                                <h4>Total Tuition</h4>
                                <p>₱{{ newPayment.type === "Full Payment (before exam)"
                                    ? (currentStudent.totalFees - 1200).toLocaleString()
                                    : currentStudent.totalFees.toLocaleString() }}</p>
                            </div>
                            <div class="summary-card">
                                <h4>Total Paid</h4>
                                <p>₱{{ (currentStudent.totalFees -
                                    currentStudent.remainingBalance).toLocaleString() }}</p>
                            </div>
                            <div class="summary-card highlight">
                                <h4>Remaining Balance</h4>
                                <p>₱{{ newPayment.type === "Full Payment (before exam)"
                                    ? (currentStudent.remainingBalance - 1200).toLocaleString()
                                    : currentStudent.remainingBalance.toLocaleString() }}</p>
                            </div>
                        </div>
                    </form>



                    <div class="form-actions">
                        <button class="submit-btn" @click="saveStudentChanges">Save Changes</button>
                        <button class="cancel-btn" @click="closeStudentModal">Close</button>
                    </div>
                </div>
            </div>

            <!-- VIEW STUDENT'S BILL MODAL -->
            <div class="modal" :class="{ active: showViewModal }">
                <div class="modal-content">
                    <div class="modal-header">
                        <h2>Student Billing Information</h2>
                        <button class="close-modal" @click="closeViewModal">&times;</button>
                    </div>


                    <div class="student-form">
                        <div class="student-info-grid">
                            <div class="info-group">
                                <label>Student ID</label>
                                <span>{{ viewStudent.id }}</span>
                            </div>
                            <div class="info-group">
                                <label>Name</label>
                                <span>{{ viewStudent.name }}</span>
                            </div>
                            <div class="info-group">
                                <label>Course</label>
                                <span>{{ viewStudent.course }}</span>
                            </div>
                            <div class="info-group">
                                <label>Year Level</label>
                                <span>{{ viewStudent.yearLevel }}</span>
                            </div>
                        </div>

                        <div class="fee-presets">
                            <div class="fee-presets-header">
                                <h4>Fee Details</h4>
                            </div>
                            <div class="student-info-grid">
                                <div class="info-group">
                                    <label>Basic Tuition</label>
                                    <span>₱{{ viewStudent.basicTuition?.toLocaleString() }}</span>
                                </div>
                                <div class="info-group">
                                    <label>Laboratory Fee</label>
                                    <span>₱{{ viewStudent.laboratoryFee?.toLocaleString() }}</span>
                                </div>
                                <div class="info-group">
                                    <label>Student Activities</label>
                                    <span>₱{{ viewStudent.studentActivities?.toLocaleString() }}</span>
                                </div>
                                <div class="info-group">
                                    <label>Discount</label>
                                    <span>₱{{ viewStudent.discount?.toLocaleString() }}</span>
                                </div>
                            </div>
                        </div>

                        <div class="fee-structure-details">
                            <div class="fee-structure-section">
                                <h3 class="labels">Tuition Fee</h3>
                                <div class="fee-breakdown-grid">
                                    <div class="fee-item">
                                        <span class="fee-label">Basic Tuition Fee</span>
                                        <span class="fee-value">₱{{ getFeeStructure(viewStudent.feePreset ||
                                            'default').tuition.registration.toLocaleString() }}</span>
                                    </div>
                                    <div class="fee-item">
                                        <span class="fee-label">Laboratory Fee</span>
                                        <span class="fee-value">₱{{ getFeeStructure(viewStudent.feePreset ||
                                            'default').tuition.subject.toLocaleString() }}</span>
                                    </div>
                                </div>
                            </div>

                            <div class="fee-structure-section">
                                <h3 class="labels">Miscellaneous Fee</h3>
                                <div class="fee-breakdown-grid">
                                    <div class="fee-item">
                                        <span class="fee-label">Development Fee</span>
                                        <span class="fee-value">₱{{ getFeeStructure(viewStudent.feePreset ||
                                            'default').misc.laboratory.toLocaleString() }}</span>
                                    </div>
                                    <div class="fee-item">
                                        <span class="fee-label">Library Fee</span>
                                        <span class="fee-value">₱{{ getFeeStructure(viewStudent.feePreset ||
                                            'default').misc.library.toLocaleString() }}</span>
                                    </div>
                                    <div class="fee-item">
                                        <span class="fee-label">Computer Laboratory Fee</span>
                                        <span class="fee-value">₱{{ getFeeStructure(viewStudent.feePreset ||
                                            'default').misc.computer.toLocaleString() }}</span>
                                    </div>
                                    <div class="fee-item">
                                        <span class="fee-label">Athletic Fee</span>
                                        <span class="fee-value">₱{{ getFeeStructure(viewStudent.feePreset ||
                                            'default').misc.athletic.toLocaleString() }}</span>
                                    </div>
                                </div>
                            </div>
                        </div>

                        <div class="exam-fees">
                            <h3 class="labels">Exam Fees</h3>
                            <div class="exam-fees-grid">
                                <div class="exam-fee-item">
                                    <span class="exam-fee-label">Prelim Exam</span>
                                    <span class="exam-fee-value">₱{{ viewStudent.examFees?.prelim.toLocaleString()
                                        }}</span>
                                </div>
                                <div class="exam-fee-item">
                                    <span class="exam-fee-label">Midterm Exam</span>
                                    <span class="exam-fee-value">₱{{ viewStudent.examFees?.midterm.toLocaleString()
                                        }}</span>
                                </div>
                                <div class="exam-fee-item">
                                    <span class="exam-fee-label">Pre-final Exam</span>
                                    <span class="exam-fee-value">₱{{ viewStudent.examFees?.prefinal.toLocaleString()
                                        }}</span>
                                </div>
                                <div class="exam-fee-item">
                                    <span class="exam-fee-label">Final Exam</span>
                                    <span class="exam-fee-value">₱{{ viewStudent.examFees?.final.toLocaleString()
                                        }}</span>
                                </div>
                            </div>
                        </div>

                        <div class="summary-cards">
                            <div class="summary-card">
                                <h4>Total Tuition</h4>
                                <p>₱{{ viewStudent.totalFees?.toLocaleString() }}</p>
                            </div>
                            <div class="summary-card">
                                <h4>Total Paid</h4>
                                <p>₱{{ (viewStudent.totalFees - viewStudent.remainingBalance)?.toLocaleString() }}</p>
                            </div>
                            <div class="summary-card highlight">
                                <h4>Remaining Balance</h4>
                                <p>₱{{ viewStudent.remainingBalance?.toLocaleString() }}</p>
                            </div>
                        </div>

                        <div v-if="viewStudent.payments && viewStudent.payments.length > 0">
                            <h3 class="labels history">Payment History</h3>
                            <table class="payment-table">
                                <thead>
                                    <tr>
                                        <th>Payment Type</th>
                                        <th>Date</th>
                                        <th>Amount</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    <tr v-for="payment in viewStudent.payments" :key="payment.id">
                                        <td>{{ payment.type }}</td>
                                        <td>{{ payment.date }}</td>
                                        <td>₱{{ payment.amount.toLocaleString() }}</td>
                                    </tr>
                                </tbody>
                            </table>
                        </div>
                    </div>

                    <div class="form-actions">
                        <button class="btn cancel-btn" @click="closeViewModal">Close</button>
                    </div>
                </div>
            </div>
        </section>
    </main>
</template>