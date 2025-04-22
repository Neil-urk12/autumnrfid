<script setup lang="ts">
import { ref, computed, defineAsyncComponent } from 'vue'
import { Student, Payment, FeeStructure } from '@/typescript/models'
import mockData from '@/mock/models.json'

const Searchbar = defineAsyncComponent(() => import("@/components/Searchbar.vue"))
const UnsavedChangesModal = defineAsyncComponent(() => import("@/components/UnsavedChangesModal.vue"))

const students = ref<Student[]>(mockData.students)
const feeStructures = ref<Record<string, FeeStructure>>(mockData.feeStructures)

const currentStudent = ref<Student>({
    id: '',
    firstName: '',
    lastName: '',
    middleName: '',
    suffix: '',
    birthday: '',
    course: '',
    block: '',
    yearLevel: '',
    status: '',
    email: '',
    phone: '',
    grades: {
        prelim: {},
        midterm: {},
        prefinals: {},
        finals: {},
        gwa: '',
        remarks: ''
    },
    billing: {
        totalTuition: 0,
        totalPaid: 0,
        remainingBalance: 0,
        tuitionFee: 0,
        miscellaneousFee: 0,
        initialPayment: 0,
        discount: 0,
        examFees: {
            prelim: 0,
            midterm: 0,
            prefinal: 0,
            final: 0
        },
        payments: []
    }
})

const originalStudentData = ref<Student | null>(null)

const viewStudent = ref<Student>({
    id: '',
    firstName: '',
    lastName: '',
    middleName: '',
    suffix: '',
    birthday: '',
    course: '',
    block: '',
    yearLevel: '',
    status: '',
    email: '',
    phone: '',
    grades: {
        prelim: {},
        midterm: {},
        prefinals: {},
        finals: {},
        gwa: '',
        remarks: ''
    },
    billing: {
        totalTuition: 0,
        totalPaid: 0,
        remainingBalance: 0,
        tuitionFee: 0,
        miscellaneousFee: 0,
        initialPayment: 0,
        discount: 0,
        examFees: {
            prelim: 0,
            midterm: 0,
            prefinal: 0,
            final: 0
        },
        payments: []
    }
})

const editStudentData = ref({
    id: '',
    name: '',
    course: '',
    yearLevel: ''
})

const newPayment = ref<Payment>({
    description: '',
    date: '',
    amount: 0
})

const newStudent = ref<Student>({
    id: '',
    firstName: '',
    lastName: '',
    middleName: '',
    suffix: '',
    birthday: '',
    course: '',
    block: '',
    yearLevel: '1st Year',
    status: 'New',
    email: '',
    phone: '',
    feePreset: 'BSCS',
    discountType: 'none',
    grades: {
        prelim: {},
        midterm: {},
        prefinals: {},
        finals: {},
        gwa: '',
        remarks: ''
    },
    billing: {
        totalTuition: 0,
        totalPaid: 0,
        remainingBalance: 0,
        tuitionFee: 0,
        miscellaneousFee: 0,
        initialPayment: 0,
        discount: 0,
        examFees: {
            prelim: 0,
            midterm: 0,
            prefinal: 0,
            final: 0
        },
        payments: []
    }
})

// STATE
const searchQuery = ref<string>('')
const activeFilters = ref<string[]>([])
const showStudentModal = ref<boolean>(false)
const showAddStudentModal = ref<boolean>(false)
const showEditFeeStructureModal = ref<boolean>(false)
const activeTab = ref<string>('info')
const showPaymentForm = ref<boolean>(false)
const showEditModal = ref<boolean>(false)
const showViewModal = ref<boolean>(false)
const isUnsavedChangesModalOpen = ref<boolean>(false)
const modalToClose = ref<string | null>(null)
const editFeePreset = ref<string>('BSCS')
const studentName = ref<string>('')
const matchingStudents = ref<Student[]>([])
const showStudentDropdown = ref<boolean>(false)
const studentIdInput = ref<string>('')

// CHECKS WHETHER THERE ARE UNSAVED CHANGES IN THE MODAL
const hasUnsavedChanges = computed(() => {
    if (showAddStudentModal.value) {
        const hasEnteredStudentId = studentIdInput.value.trim() !== '';
        const hasEnteredName = studentName.value.trim() !== '';

        if (!hasEnteredStudentId && !hasEnteredName) {
            return false;
        }
        return hasEnteredStudentId ||
            hasEnteredName ||
            newStudent.value.billing.tuitionFee > 0 ||
            newStudent.value.billing.miscellaneousFee > 0 ||
            newStudent.value.billing.initialPayment > 0;
    } else if (showStudentModal.value) {
        // CHECKS ANY MODIFICAIION FROM THE CURRENT STUDENT (EDIT STUDENT MODAL)
        if (!originalStudentData.value) return false
        return Object.keys(currentStudent.value).some(key => {
            if (key === 'payments') return false
            const currentValue = currentStudent.value[key as keyof Student];
            const originalValue = originalStudentData.value ? originalStudentData.value[key as keyof Student] : null;
            return JSON.stringify(currentValue) !== JSON.stringify(originalValue);
        })
    }
    return false
})
//

// HANDLES THE CONFIRMATION OF UNSAVED CHANGES
const handleUnsavedChanges = (confirm: boolean) => {
    if (confirm) {
        if (modalToClose.value === 'add') {
            showAddStudentModal.value = false
            newStudent.value = {
                id: '',
                firstName: '',
                lastName: '',
                middleName: '',
                suffix: '',
                birthday: '',
                course: 'BSCS',
                block: '',
                yearLevel: '1st Year',
                status: 'New',
                email: '',
                phone: '',
                feePreset: 'BSCS',
                discountType: 'none',
                grades: {
                    prelim: {},
                    midterm: {},
                    prefinals: {},
                    finals: {},
                    gwa: '',
                    remarks: ''
                },
                billing: {
                    totalTuition: 0,
                    totalPaid: 0,
                    remainingBalance: 0,
                    tuitionFee: 0,
                    miscellaneousFee: 0,
                    initialPayment: 0,
                    discount: 0,
                    examFees: {
                        prelim: 0,
                        midterm: 0,
                        prefinal: 0,
                        final: 0
                    },
                    payments: []
                }
            }
            studentName.value = ''
            studentIdInput.value = ''
            matchingStudents.value = []
            showStudentDropdown.value = false
        } else if (modalToClose.value === 'edit') {
            showStudentModal.value = false
            currentStudent.value = {
                id: '',
                firstName: '',
                lastName: '',
                middleName: '',
                suffix: '',
                birthday: '',
                course: '',
                block: '',
                yearLevel: '',
                status: '',
                email: '',
                phone: '',
                grades: {
                    prelim: {},
                    midterm: {},
                    prefinals: {},
                    finals: {},
                    gwa: '',
                    remarks: ''
                },
                billing: {
                    totalTuition: 0,
                    totalPaid: 0,
                    remainingBalance: 0,
                    tuitionFee: 0,
                    miscellaneousFee: 0,
                    initialPayment: 0,
                    discount: 0,
                    examFees: {
                        prelim: 0,
                        midterm: 0,
                        prefinal: 0,
                        final: 0
                    },
                    payments: []
                }
            }
            originalStudentData.value = null
        }
    }
    isUnsavedChangesModalOpen.value = false
    modalToClose.value = null
}
//

// GETS THE FEE STRUCTURE FOR A SPECIFIC PRESET
const getFeeStructure = (preset: string): FeeStructure => {
    return feeStructures.value[preset] || feeStructures.value['BSCS']
}
//

// LOADS THE FEE STRUCTURE FOR A SPECIFIC PRESET
const loadFeeStructure = (preset: string) => {
    editFeePreset.value = preset
    newStudent.value.feePreset = preset
    applyCoursePreset(preset)
}

// SAVES THE CURRENT FEE STRUCTURE AND UPDATES STUDENT FEES
const saveFeeStructure = () => {
    const baseFees = getFeeStructure(editFeePreset.value)
    const totalTuition = baseFees.tuition.basicTuition + baseFees.tuition.laboratory
    const totalMisc = baseFees.misc.development + baseFees.misc.library +
        baseFees.misc.computer + baseFees.misc.athletic

    newStudent.value.billing.tuitionFee = totalTuition
    newStudent.value.billing.miscellaneousFee = totalMisc
    newStudent.value.billing.initialPayment = 2750

    if (newStudent.value.yearLevel) {
        const year = newStudent.value.yearLevel.split(' ')[0].toLowerCase()
        applyYearPreset(year)
    }

    showEditFeeStructureModal.value = false
}
//

// HANDLES THE SEARCH QUERY UPDATE
const handleSearch = (query: string) => {
    searchQuery.value = query || ''
}

// HANDLES THE FILTER CHANGE
const handleFilterChange = (filters: string[]) => {
    activeFilters.value = filters || []
}
//

// FILTERS STUDENTS BASED ON SEARCH QUERY AND ACTIVE FILTERS
const filteredStudents = computed(() => {
    return students.value.filter(student => {
        if (!student.billing ||
            student.billing.totalTuition === 0 ||
            student.billing.totalTuition === null) {
            return false
        }

        const matchesSearch = !searchQuery.value ||
            student.id.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
            `${student.firstName} ${student.lastName}`.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
            student.course.toLowerCase().includes(searchQuery.value.toLowerCase())

        const matchesFilters = !activeFilters.value.length ||
            activeFilters.value.some(filter =>
                student.yearLevel.includes(filter) ||
                student.course.includes(filter))

        return matchesSearch && matchesFilters
    })
})
//

// OPENS THE ADD STUDENT MODAL WITH DEFAULT VALUES
const openAddStudentModal = () => {
    newStudent.value = {
        id: '',
        firstName: '',
        lastName: '',
        middleName: '',
        suffix: '',
        birthday: '',
        course: 'BSCS',
        block: '',
        yearLevel: '1st Year',
        status: 'New',
        email: '',
        phone: '',
        feePreset: 'BSCS',
        discountType: 'none',
        grades: {
            prelim: {},
            midterm: {},
            prefinals: {},
            finals: {},
            gwa: '',
            remarks: ''
        },
        billing: {
            totalTuition: 0,
            totalPaid: 0,
            remainingBalance: 0,
            tuitionFee: 0,
            miscellaneousFee: 0,
            initialPayment: 0,
            discount: 0,
            examFees: {
                prelim: 0,
                midterm: 0,
                prefinal: 0,
                final: 0
            },
            payments: []
        }
    }

    // RESET THE STUDENT NAME AND ID INPUT
    studentName.value = ''
    studentIdInput.value = ''
    matchingStudents.value = []
    showStudentDropdown.value = false

    applyCoursePreset('BSCS')
    applyYearPreset('1st')

    showAddStudentModal.value = true
}

// CLOSES THE ADD STUDENT MODAL WITH UNSAVED CHANGES CHECK
const closeAddStudentModal = () => {
    if (hasUnsavedChanges.value && !modalToClose.value) {
        modalToClose.value = 'add'
        isUnsavedChangesModalOpen.value = true
    } else {
        showAddStudentModal.value = false
        resetNewStudentForm();
    }
}

const resetNewStudentForm = () => {
    newStudent.value = {
        id: '',
        firstName: '',
        lastName: '',
        middleName: '',
        suffix: '',
        birthday: '',
        course: 'BSCS',
        block: '',
        yearLevel: '1st Year',
        status: 'New',
        email: '',
        phone: '',
        feePreset: 'BSCS',
        discountType: 'none',
        grades: {
            prelim: {},
            midterm: {},
            prefinals: {},
            finals: {},
            gwa: '',
            remarks: ''
        },
        billing: {
            totalTuition: 0,
            totalPaid: 0,
            remainingBalance: 0,
            tuitionFee: 0,
            miscellaneousFee: 0,
            initialPayment: 0,
            discount: 0,
            examFees: {
                prelim: 0,
                midterm: 0,
                prefinal: 0,
                final: 0
            },
            payments: []
        }
    }
    studentName.value = ''
    studentIdInput.value = ''
    matchingStudents.value = []
    showStudentDropdown.value = false
    modalToClose.value = null
}

// SAVES A NEW STUDENT TO THE STUDENT TABLE
const saveNewStudent = () => {
    if (!newStudent.value.id || !studentName.value || !newStudent.value.course || !newStudent.value.yearLevel) {
        alert('Please fill in all required fields')
        return
    }

    // CHECK IF STUDENT IS EXISTING (HAS BILLING)
    const existingStudentInTable = students.value.find(s =>
        s.id === newStudent.value.id &&
        s.billing &&
        (s.billing.totalTuition !== 0 &&
            s.billing.totalTuition !== null)
    )

    if (existingStudentInTable) {
        alert('Student ID already exists in the billing records')
        return
    }

    // MO AUTOMATIC NALANG UG FILL SA NAME INPUT IF MAKASELECT NAG STUDENT ID
    // const nameParts = studentName.value.trim().split(' ')
    // if (nameParts.length < 2) {
    //     alert('Please enter both first and last name')
    //     return
    // }

    // const lastName = nameParts.pop() || '';
    // const firstName = nameParts.join(' ');

    // newStudent.value.lastName = lastName;
    // newStudent.value.firstName = firstName;

    // COMPUTATION
    const totalTuition = newStudent.value.billing.tuitionFee + newStudent.value.billing.miscellaneousFee - newStudent.value.billing.discount
    const remainingBalance = totalTuition - newStudent.value.billing.initialPayment
    newStudent.value.billing.examFees = {
        prelim: newStudent.value.billing.tuitionFee * 0.2,
        midterm: newStudent.value.billing.tuitionFee * 0.2,
        prefinal: newStudent.value.billing.tuitionFee * 0.2,
        final: newStudent.value.billing.tuitionFee * 0.2
    }

    newStudent.value.billing.totalTuition = totalTuition
    newStudent.value.billing.totalPaid = newStudent.value.billing.initialPayment
    newStudent.value.billing.remainingBalance = remainingBalance

    // AUTOMATICALLY ADDS INITIAL PAYMENT AFTER SAVING THE STUDENT
    if (newStudent.value.billing.initialPayment > 0) {
        newStudent.value.billing.payments.push({
            description: 'Initial Payment',
            date: new Date().toISOString().split('T')[0],
            amount: newStudent.value.billing.initialPayment
        })
    }

    // CHECKS IF STUDENT EXISTED (WITHOUT BILL)
    const existingStudentIndex = students.value.findIndex(s => s.id === newStudent.value.id)

    if (existingStudentIndex !== -1) {
        students.value[existingStudentIndex].billing = JSON.parse(JSON.stringify(newStudent.value.billing))
    } else {
        students.value.push(JSON.parse(JSON.stringify(newStudent.value)))
    }

    showAddStudentModal.value = false
    resetNewStudentForm()
}
//

// OPENS THE VIEW MODAL FOR A SPECIFIC STUDENT
const openViewModal = (studentId: string) => {
    const student = students.value.find(s => s.id === studentId)
    if (!student) return

    viewStudent.value = JSON.parse(JSON.stringify(student))
    showViewModal.value = true
}

// CLOSES THE VIEW MODAL AND RESETS VIEW STUDENT DATA
const closeViewModal = () => {
    showViewModal.value = false
    viewStudent.value = {
        id: '',
        firstName: '',
        lastName: '',
        middleName: '',
        suffix: '',
        birthday: '',
        course: '',
        block: '',
        yearLevel: '',
        status: '',
        email: '',
        phone: '',
        grades: {
            prelim: {},
            midterm: {},
            prefinals: {},
            finals: {},
            gwa: '',
            remarks: ''
        },
        billing: {
            totalTuition: 0,
            totalPaid: 0,
            remainingBalance: 0,
            tuitionFee: 0,
            miscellaneousFee: 0,
            initialPayment: 0,
            discount: 0,
            examFees: {
                prelim: 0,
                midterm: 0,
                prefinal: 0,
                final: 0
            },
            payments: []
        }
    }
}
//

// OPENS THE STUDENT MODAL FOR EDITING
const openStudentModal = (studentId: string) => {
    const student = students.value.find(s => s.id === studentId)
    if (!student) return

    originalStudentData.value = JSON.parse(JSON.stringify(student))
    currentStudent.value = JSON.parse(JSON.stringify(student))

    if (!currentStudent.value.billing.examFees) {
        currentStudent.value.billing.examFees = {
            prelim: currentStudent.value.billing.tuitionFee * 0.2,
            midterm: currentStudent.value.billing.tuitionFee * 0.2,
            prefinal: currentStudent.value.billing.tuitionFee * 0.2,
            final: currentStudent.value.billing.tuitionFee * 0.2
        }
    }

    showStudentModal.value = true
    activeTab.value = 'info'
    showPaymentForm.value = false
    showEditModal.value = false
}

// CLOSES THE STUDENT MODAL WITH UNSAVED CHANGES CHECK
const closeStudentModal = () => {
    if (hasUnsavedChanges.value && !modalToClose.value) {
        modalToClose.value = 'edit'
        isUnsavedChangesModalOpen.value = true
    } else {
        showStudentModal.value = false
        currentStudent.value = {
            id: '',
            firstName: '',
            lastName: '',
            middleName: '',
            suffix: '',
            birthday: '',
            course: '',
            block: '',
            yearLevel: '',
            status: '',
            email: '',
            phone: '',
            grades: {
                prelim: {},
                midterm: {},
                prefinals: {},
                finals: {},
                gwa: '',
                remarks: ''
            },
            billing: {
                totalTuition: 0,
                totalPaid: 0,
                remainingBalance: 0,
                tuitionFee: 0,
                miscellaneousFee: 0,
                initialPayment: 0,
                discount: 0,
                examFees: {
                    prelim: 0,
                    midterm: 0,
                    prefinal: 0,
                    final: 0
                },
                payments: []
            }
        }
        originalStudentData.value = null
        modalToClose.value = null
    }
}

// INITIATES THE EDIT MODE FOR A STUDENT
const editStudent = (studentId: string) => {
    const student = students.value.find(s => s.id === studentId)
    if (!student) return

    openStudentModal(studentId)

    editStudentData.value = {
        id: student.id,
        name: `${student.firstName} ${student.lastName}`,
        course: student.course,
        yearLevel: student.yearLevel
    }

    showEditModal.value = true
}

// SAVES THE CHANGES MADE TO A STUDENT'S BILLING INFORMATION
const saveStudentChanges = () => {
    const index = students.value.findIndex(s => s.id === currentStudent.value.id)
    if (index === -1) return

    currentStudent.value.billing.totalTuition =
        currentStudent.value.billing.tuitionFee +
        currentStudent.value.billing.miscellaneousFee -
        currentStudent.value.billing.discount

    currentStudent.value.billing.totalPaid =
        currentStudent.value.billing.payments.reduce((sum, payment) => sum + payment.amount, 0)

    currentStudent.value.billing.remainingBalance =
        currentStudent.value.billing.totalTuition - currentStudent.value.billing.totalPaid

    currentStudent.value.billing.examFees = {
        prelim: currentStudent.value.billing.tuitionFee * 0.2,
        midterm: currentStudent.value.billing.tuitionFee * 0.2,
        prefinal: currentStudent.value.billing.tuitionFee * 0.2,
        final: currentStudent.value.billing.tuitionFee * 0.2
    }

    students.value[index] = JSON.parse(JSON.stringify(currentStudent.value))

    originalStudentData.value = null
    showStudentModal.value = false
}
//

// APPLIES THE COURSE PRESET AND UPDATES FEES
const applyCoursePreset = (preset: string) => {
    const baseFees = getFeeStructure(preset)

    // CALCULATE THE TOTAL TUITION AND MISC FEE
    const totalTuition = baseFees.tuition.basicTuition + baseFees.tuition.laboratory
    const totalMisc = baseFees.misc.development + baseFees.misc.library +
        baseFees.misc.computer + baseFees.misc.athletic

    newStudent.value.feePreset = preset
    newStudent.value.course = preset
    newStudent.value.billing.tuitionFee = totalTuition
    newStudent.value.billing.miscellaneousFee = totalMisc
    newStudent.value.billing.initialPayment = 2750

    if (newStudent.value.yearLevel) {
        const year = newStudent.value.yearLevel.split(' ')[0].toLowerCase()
        applyYearPreset(year)
    }

    newStudent.value.billing.discount = 0
    newStudent.value.discountType = 'none'
}

// APPLIES THE YEAR LEVEL PRESET AND UPDATES FEES
const applyYearPreset = (year: string) => {
    const feePreset = newStudent.value.feePreset || 'BSCS'
    const baseFees = getFeeStructure(feePreset)
    const totalTuition = baseFees.tuition.basicTuition + baseFees.tuition.laboratory
    const totalMisc = baseFees.misc.development + baseFees.misc.library +
        baseFees.misc.computer + baseFees.misc.athletic

    switch (year) {
        case '1st':
            newStudent.value.billing.tuitionFee = Math.round(totalTuition)
            newStudent.value.billing.miscellaneousFee = Math.round(totalMisc)
            newStudent.value.billing.initialPayment = 2750
            newStudent.value.yearLevel = '1st Year'
            break
        case '2nd':
            newStudent.value.billing.tuitionFee = Math.round(totalTuition * 1.05)
            newStudent.value.billing.miscellaneousFee = Math.round(totalMisc * 1.1)
            newStudent.value.billing.initialPayment = 2750
            newStudent.value.yearLevel = '2nd Year'
            break
        case '3rd':
            newStudent.value.billing.tuitionFee = Math.round(totalTuition * 1.1)
            newStudent.value.billing.miscellaneousFee = Math.round(totalMisc * 1.2)
            newStudent.value.billing.initialPayment = 2750
            newStudent.value.yearLevel = '3rd Year'
            break
        case '4th':
            newStudent.value.billing.tuitionFee = Math.round(totalTuition * 1.15)
            newStudent.value.billing.miscellaneousFee = Math.round(totalMisc * 1.3)
            newStudent.value.billing.initialPayment = 2750
            newStudent.value.yearLevel = '4th Year'
            break
    }

    newStudent.value.discountType = 'none'
    newStudent.value.billing.discount = 0
}

// APPLIES THE SELECTED DISCOUNT TYPE AND UPDATES FEES
const applyDiscountType = () => {
    const totalBeforeDiscount = newStudent.value.billing.tuitionFee +
        newStudent.value.billing.miscellaneousFee +
        newStudent.value.billing.initialPayment

    switch (newStudent.value.discountType) {
        case 'honor':
            newStudent.value.billing.discount = Math.round(totalBeforeDiscount * 0.15)
            break
        case 'highHonor':
            newStudent.value.billing.discount = Math.round(totalBeforeDiscount * 0.30)
            break
        case 'highestHonor':
            newStudent.value.billing.discount = Math.round(totalBeforeDiscount * 0.50)
            break
        case 'freshman':
        case 'continuing':
            newStudent.value.billing.discount = Math.round(totalBeforeDiscount * 0.10)
            break
        default:
            newStudent.value.billing.discount = 0
    }
}

// UPDATES THE PAYMENT AMOUNT BASED ON PAYMENT TYPE
const updatePaymentAmount = () => {
    if (newPayment.value.description === "Full Payment") {
        newPayment.value.amount = currentStudent.value.billing.remainingBalance
    } else if (newPayment.value.description === "Full Payment (before exam)") {
        newPayment.value.amount = Math.max(0, currentStudent.value.billing.remainingBalance - 1200)
    } else {
        newPayment.value.amount = 0
    }
}

// SUBMITS A NEW PAYMENT FOR THE CURRENT STUDENT
const submitPayment = () => {
    if (!currentStudent.value.id) return

    const description = newPayment.value.description
    let amount = parseFloat(newPayment.value.amount.toString())
    const date = newPayment.value.date

    if (!amount || isNaN(amount) || !date || !description) {
        alert('Please fill in all payment details')
        return
    }

    if (description !== "Full Payment (before exam)" && amount > currentStudent.value.billing.remainingBalance) {
        alert('Payment amount cannot exceed remaining balance')
        return
    }

    const payment: Payment = {
        description,
        date,
        amount
    }

    const index = students.value.findIndex(s => s.id === currentStudent.value.id)
    if (index === -1) return

    students.value[index].billing.payments.push(payment)
    currentStudent.value.billing.payments.push(payment)

    if (description === "Full Payment (before exam)") {
        students.value[index].billing.totalTuition -= 1200
        currentStudent.value.billing.totalTuition -= 1200
        students.value[index].billing.remainingBalance = 0
        currentStudent.value.billing.remainingBalance = 0
    } else {
        students.value[index].billing.totalPaid += amount
        currentStudent.value.billing.totalPaid += amount
        students.value[index].billing.remainingBalance -= amount
        if (students.value[index].billing.remainingBalance < 0) {
            students.value[index].billing.remainingBalance = 0
        }
        currentStudent.value.billing.remainingBalance = students.value[index].billing.remainingBalance
    }

    newPayment.value = {
        description: 'Initial Payment',
        date: '',
        amount: 0
    }
    showPaymentForm.value = false
}

// INITIATES EDITING MODE FOR A PAYMENT
const startEditingPayment = (payment: Payment) => {
    currentStudent.value.billing.payments.forEach((p: Payment) => {
        p.isEditing = p.id === payment.id;
    });

    payment.isEditing = true;
    payment.editAmount = payment.amount;
}

// CANCELS THE PAYMENT EDITING MODE
const cancelEditingPayment = (payment: Payment) => {
    payment.isEditing = false;
    payment.editAmount = undefined;
}

// SAVES THE EDITED PAYMENT AMOUNT
const saveEditedPayment = (payment: Payment) => {
    if (!payment.editAmount || isNaN(payment.editAmount)) {
        alert('Please enter a valid amount');
        return;
    }

    const newAmount = payment.editAmount;
    if (newAmount <= 0) {
        alert('Amount must be greater than zero');
        return;
    }

    const amountDifference = payment.amount - newAmount;

    const studentIndex = students.value.findIndex(s => s.id === currentStudent.value.id);
    if (studentIndex === -1) return;

    const paymentIndex = students.value[studentIndex].billing.payments.findIndex((p: Payment) => p.id === payment.id);
    if (paymentIndex === -1) return;

    students.value[studentIndex].billing.payments[paymentIndex].amount = newAmount;

    students.value[studentIndex].billing.remainingBalance += amountDifference;
    if (students.value[studentIndex].billing.remainingBalance < 0) {
        students.value[studentIndex].billing.remainingBalance = 0;
    }

    currentStudent.value.billing.payments[paymentIndex].amount = newAmount;
    currentStudent.value.billing.remainingBalance = students.value[studentIndex].billing.remainingBalance;

    payment.isEditing = false;
    payment.editAmount = undefined;
}

const handlePaymentAction = (paymentIdentifier: string | number, action: string) => {
    if (action === 'delete') {
        const studentIndex = students.value.findIndex(s => s.id === currentStudent.value.id);
        if (studentIndex === -1) return;

        const paymentIndex = students.value[studentIndex].billing.payments.findIndex((p: Payment) => 
            (p.id && p.id === paymentIdentifier) || 
            (!p.id && p.description === paymentIdentifier)
        );

        if (paymentIndex === -1) return;

        const payment = students.value[studentIndex].billing.payments[paymentIndex];

        students.value[studentIndex].billing.remainingBalance += payment.amount;
        students.value[studentIndex].billing.totalPaid -= payment.amount;

        students.value[studentIndex].billing.payments.splice(paymentIndex, 1);

        currentStudent.value.billing.remainingBalance = students.value[studentIndex].billing.remainingBalance;
        currentStudent.value.billing.totalPaid = students.value[studentIndex].billing.totalPaid;
        currentStudent.value.billing.payments = [...students.value[studentIndex].billing.payments];
    }
}

const isPaymentButtonDisabled = computed(() => {
    return !newPayment.value.description ||
        !newPayment.value.date ||
        !newPayment.value.amount ||
        isNaN(parseFloat(newPayment.value.amount.toString()));
})

// APPLIES THE DISCOUNT TO THE CURRENT STUDENT
const applyDiscount = () => {
    const totalBeforeDiscount = currentStudent.value.billing.tuitionFee +
        currentStudent.value.billing.miscellaneousFee;

    currentStudent.value.billing.totalTuition = totalBeforeDiscount - currentStudent.value.billing.discount;

    const totalPaid = currentStudent.value.billing.totalPaid;
    currentStudent.value.billing.remainingBalance = currentStudent.value.billing.totalTuition - totalPaid;

    if (currentStudent.value.billing.remainingBalance < 0) {
        currentStudent.value.billing.remainingBalance = 0;
    }
}

const hasBillingInfo = (student: Student): boolean => {
    return student.billing &&
        ((typeof student.billing.totalTuition === 'number' && student.billing.totalTuition > 0) ||
            (typeof student.billing.totalTuition === 'string' && student.billing.totalTuition !== ''));
}

const safeLocaleString = (value: string | number | null | undefined): string => {
    if (value === null || value === undefined || value === '') {
        return '0.00';
    }

    if (typeof value === 'string') {
        const parsedValue = parseFloat(value);
        if (!isNaN(parsedValue)) {
            return parsedValue.toLocaleString();
        }
        return '0.00';
    }

    if (typeof value === 'number') {
        return value.toLocaleString();
    }

    return '0.00';
}


// SELECTS A STUDENT ID FROM THE DROPDOWN (ADD STUDENT BILLING MODAL CAN ONLY SELECT STUDENT WITHOUT BILLING)
const selectStudent = (student: Student) => {
    const hasBilling = hasBillingInfo(student);

    newStudent.value.id = student.id;
    studentIdInput.value = student.id;
    studentName.value = `${student.firstName} ${student.lastName}`;

    showStudentDropdown.value = false;
    if (hasBilling) {
        return;
    }

    newStudent.value.firstName = student.firstName;
    newStudent.value.lastName = student.lastName;
    newStudent.value.middleName = student.middleName || '';
    newStudent.value.suffix = student.suffix || '';
    newStudent.value.course = student.course;
    newStudent.value.yearLevel = student.yearLevel;
    newStudent.value.block = student.block || '';
    newStudent.value.birthday = student.birthday || '';
    newStudent.value.email = student.email || '';
    newStudent.value.phone = student.phone || '';
    newStudent.value.status = student.status || 'New';

    applyCoursePreset(student.course);
    const year = student.yearLevel.split(' ')[0].toLowerCase();
    applyYearPreset(year);
}

const handleStudentIdChange = (event: Event) => {
    const target = event.target as HTMLInputElement;
    const id = target.value;
    studentIdInput.value = id;
    searchStudentById(id);
}

    // ATUOMATICALLY MULTIPLIES THE FEES BASED ON YEAR LEVEL
const getYearAdjustedFees = (course: string, yearLevel: string) => {
    const courseToUse = course || 'BSCS'
    const baseFees = getFeeStructure(courseToUse);
    const totalBaseTuition = baseFees.tuition.basicTuition;
    const totalBaseLaboratory = baseFees.tuition.laboratory;
    const totalBaseDevelopment = baseFees.misc.development;
    const totalBaseLibrary = baseFees.misc.library;
    const totalBaseComputer = baseFees.misc.computer;
    const totalBaseAthletic = baseFees.misc.athletic;

    let tuitionMultiplier = 1;
    let miscMultiplier = 1;

    if (yearLevel.includes('2nd')) {
        tuitionMultiplier = 1.05;
        miscMultiplier = 1.1;
    } else if (yearLevel.includes('3rd')) {
        tuitionMultiplier = 1.1;
        miscMultiplier = 1.2;
    } else if (yearLevel.includes('4th')) {
        tuitionMultiplier = 1.15;
        miscMultiplier = 1.3;
    }

    return {
        tuition: {
            basicTuition: Math.round(totalBaseTuition * tuitionMultiplier),
            laboratory: Math.round(totalBaseLaboratory * tuitionMultiplier)
        },
        misc: {
            development: Math.round(totalBaseDevelopment * miscMultiplier),
            library: Math.round(totalBaseLibrary * miscMultiplier),
            computer: Math.round(totalBaseComputer * miscMultiplier),
            athletic: Math.round(totalBaseAthletic * miscMultiplier)
        }
    };
}

// SEARCH STUDENT BY ID
const searchStudentById = (id: string) => {
    if (!id.trim()) {
        matchingStudents.value = [];
        showStudentDropdown.value = false;
        return;
    }

    const matches = students.value.filter(student =>
        student.id.toLowerCase().includes(id.toLowerCase())
    );

    matchingStudents.value = matches;
    showStudentDropdown.value = true; 
}

// HANDLES THE COURSE CHANGE AND UPDATES FEE STRUCTURE
const handleCourseChange = () => {
    loadFeeStructure(newStudent.value.course)
}

// HANDLES THE YEAR LEVEL CHANGE AND UPDATES FEES
const handleYearChange = () => {
    const year = newStudent.value.yearLevel.split(' ')[0].toLowerCase()
    applyYearPreset(year)
}
</script>


<template>
    <main>

        <section>
            <div class="container">
                <div class="welcome-header">
                    <h1>Student Bills</h1>
                    <p>Manage and monitor student payments and balances</p>
                </div>

                <!-- SEARCH BAR SECTION -->
                <div class="students-controls">
                    <div class="search-filters">
                        <Searchbar v-model="searchQuery" @update:search-query="handleSearch"
                            @filter-change="handleFilterChange" />
                        <div class="filter-buttons">
                            <button class="add-student-btn" @click="openAddStudentModal">
                                <i class="fa-solid fa-plus"></i> Add New Student Bill
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
                                <th>Total Tuition</th>
                                <th>Paid Amount</th>
                                <th>Remaining Balance</th>
                                <th>Action</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="student in filteredStudents" :key="student.id">
                                <td>{{ student.id }}</td>
                                <td>{{ `${student.firstName} ${student.lastName}` }}</td>
                                <td>{{ student.course }}</td>
                                <td>₱{{ safeLocaleString(student.billing?.totalTuition) }}</td>
                                <td>₱{{ safeLocaleString(student.billing?.totalPaid) }}</td>
                                <td>₱{{ safeLocaleString(student.billing?.remainingBalance) }}</td>
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
            <div class="modal" :class="{ active: showAddStudentModal }" @click="closeAddStudentModal">
                <div class="modal-content" @click.stop>
                    <div class="modal-header">
                        <h2>Add New Student Bill</h2>
                        <button class="close-modal" @click="closeAddStudentModal">&times;</button>
                    </div>

                    <form class="student-form">
                        <div class="student-info-grid">
                            <div class="info-group">
                                <label>Student ID</label>
                                <div class="search-input-container">
                                    <input type="text" v-model="studentIdInput" @input="handleStudentIdChange($event)"
                                        placeholder="Enter Student ID">
                                    <div v-if="showStudentDropdown" class="student-dropdown">
                                        <div v-if="matchingStudents.length === 0" class="no-students">
                                            No matching students found
                                        </div>
                                        <div v-else v-for="student in matchingStudents" :key="student.id"
                                            @click="!hasBillingInfo(student) && selectStudent(student)" :class="[
                                                'student-item',
                                                { 'disabled-item': hasBillingInfo(student) }
                                            ]">
                                            <div class="student-info">
                                                <span class="student-id">{{ student.id }}</span>
                                                <span class="student-name">{{ student.firstName }} {{ student.lastName
                                                    }}</span>
                                            </div>
                                            <div class="status-tag" :class="{
                                                'has-billing': hasBillingInfo(student)
                                            }">
                                                {{ hasBillingInfo(student)
                                                    ? 'Already has billing'
                                                    : 'No billing yet' }}
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </div>
                            <div class="info-group">
                                <label>Name</label>
                                <input type="text" v-model="studentName"
                                    placeholder="Student Name">
                            </div>
                            <div class="info-group">
                                <label>Course</label>
                                <select v-model="newStudent.course" @change="handleCourseChange">
                                    <option value="">Select course</option>
                                    <option value="BSCS">BSCS</option>
                                    <option value="BSIT">BSIT</option>
                                    <option value="BSBA">BSBA</option>
                                    <option value="BSA">BSA</option>
                                </select>
                            </div>
                            <div class="info-group">
                                <label>Year Level</label>
                                <select v-model="newStudent.yearLevel" @change="handleYearChange">
                                    <option value="">Select year level</option>
                                    <option value="1st Year">1st Year</option>
                                    <option value="2nd Year">2nd Year</option>
                                    <option value="3rd Year">3rd Year</option>
                                    <option value="4th Year">4th Year</option>
                                </select>
                            </div>
                        </div>

                        <div class="course-selector">
                            <div class="course-option" :class="{ active: editFeePreset === 'BSCS' }"
                                @click="loadFeeStructure('BSCS')">
                                BSCS
                            </div>
                            <div class="course-option" :class="{ active: editFeePreset === 'BSIT' }"
                                @click="loadFeeStructure('BSIT')">
                                BSIT
                            </div>
                            <div class="course-option" :class="{ active: editFeePreset === 'BSBA' }"
                                @click="loadFeeStructure('BSBA')">
                                BSBA
                            </div>
                            <div class="course-option" :class="{ active: editFeePreset === 'BSA' }"
                                @click="loadFeeStructure('BSA')">
                                BSA
                            </div>
                        </div>

                        <!-- MODAL FOR CHANGING THE FEE STRUCTURE FOR EVERY COURSE -->
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
                                    <input type="number" v-model="newStudent.billing.tuitionFee">
                                </div>
                                <div class="info-group">
                                    <label>Miscellaneous Fee</label>
                                    <input type="number" v-model="newStudent.billing.miscellaneousFee">
                                </div>
                                <div class="info-group">
                                    <label>Initial Payment</label>
                                    <input type="number" v-model="newStudent.billing.initialPayment">
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
                                <button type="button" class="preset-btn"
                                    :class="{ active: newStudent.yearLevel === '1st Year' }"
                                    @click.prevent="applyYearPreset('1st')">
                                    1st Year
                                </button>
                                <button type="button" class="preset-btn"
                                    :class="{ active: newStudent.yearLevel === '2nd Year' }"
                                    @click.prevent="applyYearPreset('2nd')">
                                    2nd Year
                                </button>
                                <button type="button" class="preset-btn"
                                    :class="{ active: newStudent.yearLevel === '3rd Year' }"
                                    @click.prevent="applyYearPreset('3rd')">
                                    3rd Year
                                </button>
                                <button type="button" class="preset-btn"
                                    :class="{ active: newStudent.yearLevel === '4th Year' }"
                                    @click.prevent="applyYearPreset('4th')">
                                    4th Year
                                </button>
                            </div>
                        </div>

                        <div class="exam-fees">
                            <h3 class="labels">Exam Fees</h3>
                            <div class="exam-fees-grid">
                                <div class="exam-fee-item">
                                    <span class="exam-fee-label">Prelim Exam</span>
                                    <span class="exam-fee-value">₱{{ safeLocaleString(newStudent.billing.tuitionFee *
                                        0.2) }}</span>
                                </div>
                                <div class="exam-fee-item">
                                    <span class="exam-fee-label">Midterm Exam</span>
                                    <span class="exam-fee-value">₱{{ safeLocaleString(newStudent.billing.tuitionFee *
                                        0.2) }}</span>
                                </div>
                                <div class="exam-fee-item">
                                    <span class="exam-fee-label">Pre-final Exam</span>
                                    <span class="exam-fee-value">₱{{ safeLocaleString(newStudent.billing.tuitionFee *
                                        0.2) }}</span>
                                </div>
                                <div class="exam-fee-item">
                                    <span class="exam-fee-label">Final Exam</span>
                                    <span class="exam-fee-value">₱{{ safeLocaleString(newStudent.billing.tuitionFee *
                                        0.2) }}</span>
                                </div>
                            </div>
                        </div>

                        <div class="summary-cards">
                            <div class="summary-card">
                                <h4>Total Tuition</h4>
                                <p>₱{{ safeLocaleString(newStudent.billing.tuitionFee +
                                    newStudent.billing.miscellaneousFee -
                                    newStudent.billing.discount) }}</p>
                            </div>
                            <div class="summary-card">
                                <h4>Total Paid</h4>
                                <p>₱{{ safeLocaleString(newStudent.billing.initialPayment) }}</p>
                            </div>
                            <div class="summary-card highlight">
                                <h4>Remaining Balance</h4>
                                <p>₱{{ safeLocaleString(newStudent.billing.tuitionFee +
                                    newStudent.billing.miscellaneousFee -
                                    newStudent.billing.discount - newStudent.billing.initialPayment) }}</p>
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
                            <div class="course-option" :class="{ active: editFeePreset === 'BSCS' }"
                                @click="loadFeeStructure('BSCS')">
                                BSCS
                            </div>
                            <div class="course-option" :class="{ active: editFeePreset === 'BSIT' }"
                                @click="loadFeeStructure('BSIT')">
                                BSIT
                            </div>
                            <div class="course-option" :class="{ active: editFeePreset === 'BSBA' }"
                                @click="loadFeeStructure('BSBA')">
                                BSBA
                            </div>
                            <div class="course-option" :class="{ active: editFeePreset === 'BSA' }"
                                @click="loadFeeStructure('BSA')">
                                BSA
                            </div>
                        </div>


                        <form class="student-form">
                            <h3 class="labels">Tuition Fee</h3>
                            <div class="student-info-grid">
                                <div class="info-group">
                                    <label>Basic Tuition Fee</label>
                                    <input type="number" v-model="feeStructures[editFeePreset].tuition.basicTuition">
                                </div>
                                <div class="info-group">
                                    <label>Laboratory Fee</label>
                                    <input type="number" v-model="feeStructures[editFeePreset].tuition.laboratory">
                                </div>
                            </div>

                            <h3 class="labels">Miscellaneous Fee</h3>
                            <div class="student-info-grid">
                                <div class="info-group">
                                    <label>Development Fee</label>
                                    <input type="number" v-model="feeStructures[editFeePreset].misc.development">
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
            <div class="modal" :class="{ active: showStudentModal }" @click="closeStudentModal">
                <div class="modal-content" @click.stop>
                    <div class="modal-header">
                        <h2>Edit Billing Details</h2>
                        <button class="close-modal" @click="closeStudentModal">&times;</button>
                    </div>

                    <form class="student-form">
                        <div class="student-info-grid">
                            <div class="info-group">
                                <label>Student ID</label>
                                <input type="text" v-model="currentStudent.id" readonly disabled>
                            </div>
                            <div class="info-group">
                                <label>Name</label>
                                <input type="text" :value="`${currentStudent.firstName} ${currentStudent.lastName}`"
                                    readonly disabled>
                            </div>
                            <div class="info-group">
                                <label>Tuition Fee</label>
                                <input type="number" v-model="currentStudent.billing.tuitionFee">
                            </div>
                            <div class="info-group">
                                <label>Miscellaneous Fee</label>
                                <input type="number" v-model="currentStudent.billing.miscellaneousFee">
                            </div>
                            <div class="info-group">
                                <label>Initial Payment</label>
                                <input type="number" v-model="currentStudent.billing.initialPayment">
                            </div>
                            <div class="info-group">
                                <label>Discount</label>
                                <select v-model="currentStudent.billing.discount" @change="applyDiscount">
                                    <template v-if="currentStudent.yearLevel === '1st Year'">
                                        <option :value="0">No Discount</option>
                                        <option :value="currentStudent.billing.tuitionFee * 0.15">Honor Student (15%)
                                        </option>
                                        <option :value="currentStudent.billing.tuitionFee * 0.30">High Honor (30%)
                                        </option>
                                        <option :value="currentStudent.billing.tuitionFee * 0.50">Highest Honor (50%)
                                        </option>
                                        <option :value="currentStudent.billing.tuitionFee * 0.10">Freshman (10%)
                                        </option>
                                    </template>
                                    <template v-else>
                                        <option :value="0">No Discount</option>
                                        <option :value="currentStudent.billing.tuitionFee * 0.10">Continuing Student
                                            (10%)
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
                                        currentStudent.billing.examFees?.prelim.toLocaleString()
                                        }}</span>
                                </div>
                                <div class="exam-fee-item">
                                    <span class="exam-fee-label">Midterm Exam</span>
                                    <span class="exam-fee-value">₱{{
                                        currentStudent.billing.examFees?.midterm.toLocaleString()
                                        }}</span>
                                </div>
                                <div class="exam-fee-item">
                                    <span class="exam-fee-label">Pre-final Exam</span>
                                    <span class="exam-fee-value">₱{{
                                        currentStudent.billing.examFees?.prefinal.toLocaleString()
                                        }}</span>
                                </div>
                                <div class="exam-fee-item">
                                    <span class="exam-fee-label">Final Exam</span>
                                    <span class="exam-fee-value">₱{{
                                        currentStudent.billing.examFees?.final.toLocaleString()
                                        }}</span>
                                </div>
                            </div>
                        </div>

                        <div class="payment-actions">
                            <h3 class="labels">Payments</h3>
                            <div>
                                <button class="action-button" @click="showPaymentForm = true" v-if="!showPaymentForm"
                                    :disabled="currentStudent.billing.remainingBalance <= 0"
                                    :class="{ 'disabled': currentStudent.billing.remainingBalance <= 0 }">
                                    <i class="fas fa-plus"></i>
                                    {{ currentStudent.billing.remainingBalance <= 0 ? 'No Balance Remaining'
                                        : 'Add Payment' }} </button>
                            </div>
                        </div>

                        <!-- ADD STUDENT'S PAYMENT OR BILLS -->
                        <div class="payment-form" :class="{ active: showPaymentForm }">
                            <div class="payment-form-grid">
                                <div class="info-group">
                                    <label>Payment Type</label>
                                    <select v-model="newPayment.description" @change="updatePaymentAmount">
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

                        <!-- THIS TABLE WILL APPEAR IF THERE'S A PAYMENT TRANSACTION THAT HAS BEEN DONE -->
                        <table class="payment-table"
                            v-if="currentStudent.billing.payments && currentStudent.billing.payments.length > 0">
                            <thead>
                                <tr>
                                    <th>Payment Type</th>
                                    <th>Date</th>
                                    <th>Amount</th>
                                    <th>Actions</th>
                                </tr>
                            </thead>
                            <tbody>
                                <tr v-for="payment in currentStudent.billing.payments"
                                    :key="payment.id || payment.description">
                                    <td>{{ payment.description }}</td>
                                    <td>{{ payment.date }}</td>
                                    <td>₱{{ safeLocaleString(payment.amount) }}</td>
                                    <td>
                                        <div class="payment-actions-row">
                                            <button class="payment-action-btn payment-edit" v-if="!payment.isEditing"
                                                @click="startEditingPayment(payment)">
                                                <i class="fa-solid fa-pen-to-square"></i>
                                            </button>
                                            <div v-else class="payment-edit-form">
                                                <input type="number" v-model.number="payment.editAmount"
                                                    :max="currentStudent.billing.remainingBalance + payment.amount"
                                                    min="1">
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
                                                @click="handlePaymentAction(payment.id || payment.description, 'delete')">
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
                                <p>₱{{ safeLocaleString(currentStudent.billing.tuitionFee +
                                    currentStudent.billing.miscellaneousFee
                                    - currentStudent.billing.discount) }}</p>
                            </div>
                            <div class="summary-card">
                                <h4>Total Paid</h4>
                                <p>₱{{ safeLocaleString(currentStudent.billing.totalPaid) }}</p>
                            </div>
                            <div class="summary-card highlight">
                                <h4>Remaining Balance</h4>
                                <p>₱{{ safeLocaleString(currentStudent.billing.tuitionFee +
                                    currentStudent.billing.miscellaneousFee
                                    - currentStudent.billing.discount - currentStudent.billing.totalPaid) }}</p>
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
            <div class="modal" :class="{ active: showViewModal }" @click="closeViewModal">
                <div class="modal-content" @click.stop>
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
                                <span>{{ `${viewStudent.firstName} ${viewStudent.lastName}` }}</span>
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
                                    <label>Tuition Fee</label>
                                    <span>₱{{ safeLocaleString(viewStudent.billing?.tuitionFee) }}</span>
                                </div>
                                <div class="info-group">
                                    <label>Miscellaneous Fee</label>
                                    <span>₱{{ safeLocaleString(viewStudent.billing?.miscellaneousFee) }}</span>
                                </div>
                                <div class="info-group">
                                    <label>Initial Payment</label>
                                    <span>₱{{ safeLocaleString(viewStudent.billing?.initialPayment) }}</span>
                                </div>
                                <div class="info-group">
                                    <label>Discount</label>
                                    <span>₱{{ safeLocaleString(viewStudent.billing?.discount) }}</span>
                                </div>
                            </div>
                        </div>

                        <div class="fee-structure-details">
                            <div class="fee-structure-section">
                                <h3 class="labels">Tuition Fee</h3>
                                <div class="fee-breakdown-grid">
                                    <div class="fee-item">
                                        <span class="fee-label">Basic Tuition Fee</span>
                                        <span class="fee-value">₱{{ getYearAdjustedFees(viewStudent.course, viewStudent.yearLevel).tuition.basicTuition.toLocaleString() }}</span>
                                    </div>
                                    <div class="fee-item">
                                        <span class="fee-label">Laboratory Fee</span>
                                        <span class="fee-value">₱{{ getYearAdjustedFees(viewStudent.course, viewStudent.yearLevel).tuition.laboratory.toLocaleString() }}</span>
                                    </div>
                                </div>
                            </div>

                            <div class="fee-structure-section">
                                <h3 class="labels">Miscellaneous Fee</h3>
                                <div class="fee-breakdown-grid">
                                    <div class="fee-item">
                                        <span class="fee-label">Development Fee</span>
                                        <span class="fee-value">₱{{ getYearAdjustedFees(viewStudent.course, viewStudent.yearLevel).misc.development.toLocaleString() }}</span>
                                    </div>
                                    <div class="fee-item">
                                        <span class="fee-label">Library Fee</span>
                                        <span class="fee-value">₱{{ getYearAdjustedFees(viewStudent.course, viewStudent.yearLevel).misc.library.toLocaleString() }}</span>
                                    </div>
                                    <div class="fee-item">
                                        <span class="fee-label">Computer Laboratory Fee</span>
                                        <span class="fee-value">₱{{ getYearAdjustedFees(viewStudent.course, viewStudent.yearLevel).misc.computer.toLocaleString() }}</span>
                                    </div>
                                    <div class="fee-item">
                                        <span class="fee-label">Athletic Fee</span>
                                        <span class="fee-value">₱{{ getYearAdjustedFees(viewStudent.course, viewStudent.yearLevel).misc.athletic.toLocaleString() }}</span>
                                    </div>
                                </div>
                            </div>
                        </div>


                        <div class="exam-fees">
                            <h3 class="labels">Exam Fees</h3>
                            <div class="exam-fees-grid">
                                <div class="exam-fee-item">
                                    <span class="exam-fee-label">Prelim Exam</span>
                                    <span class="exam-fee-value">₱{{
                                        safeLocaleString(viewStudent.billing?.examFees?.prelim)
                                        }}</span>
                                </div>
                                <div class="exam-fee-item">
                                    <span class="exam-fee-label">Midterm Exam</span>
                                    <span class="exam-fee-value">₱{{
                                        safeLocaleString(viewStudent.billing?.examFees?.midterm)
                                        }}</span>
                                </div>
                                <div class="exam-fee-item">
                                    <span class="exam-fee-label">Pre-final Exam</span>
                                    <span class="exam-fee-value">₱{{
                                        safeLocaleString(viewStudent.billing?.examFees?.prefinal)
                                        }}</span>
                                </div>
                                <div class="exam-fee-item">
                                    <span class="exam-fee-label">Final Exam</span>
                                    <span class="exam-fee-value">₱{{
                                        safeLocaleString(viewStudent.billing?.examFees?.final)
                                        }}</span>
                                </div>
                            </div>
                        </div>

                              <!-- THIS PAYMENT HISTORY TABLE WILL APPEAR IF THERE IS AN EXISTING PAYMENT TRANSACTION -->
                              <div v-if="viewStudent.billing?.payments && viewStudent.billing.payments.length > 0">
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
                                    <tr v-for="payment in viewStudent.billing.payments"
                                        :key="payment.id || payment.description">
                                        <td>{{ payment.description }}</td>
                                        <td>{{ payment.date }}</td>
                                        <td>₱{{ safeLocaleString(payment.amount) }}</td>
                                    </tr>
                                </tbody>
                            </table>
                        </div>

                        <div class="summary-cards">
                            <div class="summary-card">
                                <h4>Total Tuition</h4>
                                <p>₱{{ safeLocaleString(viewStudent.billing?.totalTuition) }}</p>
                            </div>
                            <div class="summary-card">
                                <h4>Total Paid</h4>
                                <p>₱{{ safeLocaleString(viewStudent.billing?.totalPaid) }}</p>
                            </div>
                            <div class="summary-card highlight">
                                <h4>Remaining Balance</h4>
                                <p>₱{{ safeLocaleString(viewStudent.billing?.remainingBalance) }}</p>
                            </div>
                        </div>

                  
                    </div>

                    <div class="form-actions">
                        <button class="btn cancel-btn" @click="closeViewModal">Close</button>
                    </div>
                </div>
            </div>

            <!-- UNSAVED CHANGES WARNING -->
            <UnsavedChangesModal :is-open="isUnsavedChangesModalOpen" @close="handleUnsavedChanges(false)"
                @confirm="handleUnsavedChanges(true)" />
        </section>

    </main>
</template>

<style>
.search-input-container {
    position: relative;
    width: 100%;
}

.student-dropdown {
    position: absolute;
    top: 100%;
    left: 0;
    width: 100%;
    max-height: 200px;
    overflow-y: auto;
    color: white;
    background-color: white;
    border: 1px solid #ddd;
    border-radius: 4px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
    z-index: 1000;
    margin-top: 5px;
}

.student-item {
    padding: 10px 12px;
    border-bottom: 1px solid #eee;
    cursor: pointer;
    transition: background-color 0.2s;
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.student-item:hover:not(.disabled-item) {
    background-color: #f5f5f5;
}

.disabled-item {
    cursor: not-allowed;
    opacity: 0.7;
    background-color: #f0f0f0;
}

.student-item:last-child {
    border-bottom: none;
}

.student-info {
    display: flex;
    flex-direction: column;
}

.student-id {
    font-weight: bold;
    color: #444;
    font-size: 1em;
}

.student-name {
    color: #666;
    font-size: 0.9em;
    margin-top: 2px;
}

.status-tag {
    display: inline-block;
    padding: 4px 8px;
    font-size: 0.75em;
    border-radius: 4px;
    color: white;
    background-color: #4caf50;
    white-space: nowrap;
}

.status-tag.has-billing {
    background-color: #f44336;
}

.no-students {
    padding: 12px;
    text-align: center;
    color: #888;
    font-size: 0.9em;
}


</style>