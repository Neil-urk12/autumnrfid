export interface Student {
    id: string;
    firstName: string;
    lastName: string;
    middleName: string;
    suffix: string;
    birthday: string;
    course: string;
    block: string;
    yearLevel: string;
    status: string;
    email: string;
    phone: string;
    grades: StudentGrades;
    billing: StudentBilling;
    feePreset?: string;
    discountType?: string;
}

export interface BaseStudent {
    id: string;
    firstName: string;
    lastName: string;
    middleName: string;
    suffix?: string;
    birthday: string;
    course: string;
    block: string;
    yearLevel: string;
    status: string;
    email: string;
    phone: string;
}

// export interface Student {
//     id: string; // Corresponds to student_ID (string)
//     firstName: string; // Corresponds to first_Name (*string)
//     lastName: string; // Corresponds to last_Name (*string)
//     middleName: string; // Corresponds to middle_Name (*string)
//     suffix?: string; // Added back as optional, accessed in frontend
//     birthday: string; // Corresponds to birthday (*string)
//     course: string; // Corresponds to program (*string)
//     block: string;
//     yearLevel: number | null | string;
//     status: string;
//     email: string;
//     phone: string;
//     grades?: StudentGrades;
//     billing?: StudentBilling;
//     feePreset?: string;
//     discountType?: string;
// }

// export interface Student {
//     id: string; // Corresponds to student_ID (string)
//     firstName: string; // Corresponds to first_Name (*string)
//     lastName: string; // Corresponds to last_Name (*string)
//     middleName: string; // Corresponds to middle_Name (*string)
//     suffix?: string; // Added back as optional, accessed in frontend
//     birthday: string; // Corresponds to birthday (*string)
//     course: string; // Corresponds to program (*string)
//     block: string;
//     yearLevel: number | null | string;
//     status: string;
//     email: string;
//     phone: string;
//     grades?: StudentGrades;
//     billing?: StudentBilling;
//     feePreset?: string;
//     discountType?: string;
// }

export interface StudentGrades {
    prelim: Record<string, string>;
    midterm: Record<string, string>;
    prefinals: Record<string, string>;
    finals: Record<string, string>;
    gwa: string;
    remarks: string;
}

export interface StudentBilling {
    totalTuition: number;
    totalPaid: number;
    remainingBalance: number;
    tuitionFee: number;
    miscellaneousFee: number;
    initialPayment: number;
    discount: number;
    examFees: ExamFees;
    payments: Payment[];
}

export interface StudentBill {
    description: string;
    date: string;
    amount: number;
}

export interface Course {
    ecode: string;
    subjectCode: string;
    courseName: string;
    units: number;
}

export interface CourseForm {
    ecode: string;
    subjectCode: string;
    courseName: string;
    units: string;
    originalEcode: string;
}

export interface Subject {
    code: string;
    name: string;
}

export interface SubjectEditable extends Subject {
    grade: string;
}

export interface SubjectWithGrades extends Subject {
    grade: string;
    totalGrade?: string;
    periodGrades?: Record<string, string>;
}

export interface StudentBasicInfo {
    id: string;
    name: string;
    course: string;
}

export interface ConfirmationData {
    title: string;
    itemName: string;
    itemInfo: any;
}

export interface FeeStructure {
    tuition: {
        basicTuition: number;
        laboratory: number;
    };
    misc: {
        development: number;
        library: number;
        computer: number;
        athletic: number;
    };
}

export interface ExamFees {
    prelim: number;
    midterm: number;
    prefinal: number;
    final: number;
}

export interface Payment {
    description: string;
    date: string;
    amount: number;
    id?: string | number;
    isEditing?: boolean;
    editAmount?: number;
}

// Added type based on Go backend model
export interface StudentAssessmentSummary {
    student_id: string;
    name: string;
    course: string;
    year_level: string;
    status: string;
}

// Added types to match Go backend pagination response
export interface PaginationMetadata {
    currentPage: number;
    pageSize: number;
    totalItems: number;
    totalPages: number;
}

export interface PaginatedStudentAssessmentResponse {
    data: StudentAssessmentSummary[];
    pagination: PaginationMetadata;
}
