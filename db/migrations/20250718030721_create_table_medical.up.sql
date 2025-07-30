CREATE TABLE plan_types (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) UNIQUE NOT NULL,
    description TEXT
);

-- transaction_types
CREATE TABLE transaction_types (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) UNIQUE NOT NULL
);

-- limitation_types
CREATE TABLE limitation_types (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) UNIQUE NOT NULL
);

-- departments
CREATE TABLE departments (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

-- employees (Sekarang sebagai "induk" bagi pasien, jadi buat dulu)
CREATE TABLE employees (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    -- patient_id dihapus dari sini, relasi ownership ke patient sekarang ada di tabel patients
    department_id INT NOT NULL,
    position VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    phone VARCHAR(255) NOT NULL,
    birth_date DATE NOT NULL,
    gender ENUM('male', 'female') NOT NULL,
    plan_type_id INT NOT NULL,
    dependence VARCHAR(255),
    bank_number VARCHAR(255) NOT NULL,
    join_date DATE NOT NULL,
    CONSTRAINT fk_employees_department
        FOREIGN KEY (department_id) REFERENCES departments(id)
        ON DELETE RESTRICT -- Department tidak dihapus jika ada employee
        ON UPDATE CASCADE,
    CONSTRAINT fk_employees_plan_type
        FOREIGN KEY (plan_type_id) REFERENCES plan_types(id)
        ON DELETE RESTRICT -- PlanType tidak dihapus jika ada employee
        ON UPDATE CASCADE
);

-- family_members (Sekarang sebagai "induk" bagi pasien, buat setelah employees karena ada FK ke employees)
CREATE TABLE family_members (
    id INT PRIMARY KEY AUTO_INCREMENT,
    -- patient_id dihapus dari sini, relasi ownership ke patient sekarang ada di tabel patients
    employee_id INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    plan_type_id INT NOT NULL,
    birth_date DATE NOT NULL,
    gender ENUM('male', 'female') NOT NULL,
    CONSTRAINT fk_family_members_employee
        FOREIGN KEY (employee_id) REFERENCES employees(id)
        ON DELETE CASCADE -- Jika employee dihapus, family_member-nya juga dihapus
        ON UPDATE CASCADE,
    CONSTRAINT fk_family_members_plan_type
        FOREIGN KEY (plan_type_id) REFERENCES plan_types(id)
        ON DELETE RESTRICT -- PlanType tidak dihapus jika ada family_member
        ON UPDATE CASCADE
);


-- patients (Sekarang menjadi "anak" dari employees dan family_members secara logis)
CREATE TABLE patients (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    birth_date DATE NOT NULL,
    gender ENUM('male', 'female') NOT NULL,
    
    -- Foreign Key ke Employee: NULLABLE karena patient mungkin bukan employee
    employee_id INT UNIQUE, -- UNIQUE karena 1 employee punya 1 patient (jika ada)
    -- Foreign Key ke FamilyMember: NULLABLE karena patient mungkin bukan family member
    family_member_id INT UNIQUE, -- UNIQUE karena 1 family member punya 1 patient (jika ada)

    -- Kunci: ON DELETE CASCADE ada di sini
    CONSTRAINT fk_patients_employee_cascade
        FOREIGN KEY (employee_id) REFERENCES employees(id)
        ON DELETE CASCADE -- Jika employee dihapus, patient juga dihapus
        ON UPDATE CASCADE,
    CONSTRAINT fk_patients_family_member_cascade
        FOREIGN KEY (family_member_id) REFERENCES family_members(id)
        ON DELETE CASCADE -- Jika family_member dihapus, patient juga dihapus
        ON UPDATE CASCADE
);

-- benefits
CREATE TABLE benefits (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    plan_type_id INT NOT NULL,
    detail TEXT,
    code VARCHAR(255) UNIQUE NOT NULL,
    limitation_type_id INT NOT NULL,
    plafond INT NOT NULL,
    yearly_max INT NOT NULL,
    CONSTRAINT fk_benefits_plan_type
        FOREIGN KEY (plan_type_id) REFERENCES plan_types(id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE,
    CONSTRAINT fk_benefits_limitation_type
        FOREIGN KEY (limitation_type_id) REFERENCES limitation_types(id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE
);

-- claims
CREATE TABLE claims (
    id INT PRIMARY KEY AUTO_INCREMENT,
    patient_id INT NOT NULL, -- FK ke tabel patients
    employee_id INT NOT NULL, -- FK ke tabel employees
    benefit_id INT NOT NULL,
    claim_amount DECIMAL(10, 2) NOT NULL,
    transaction_type_id INT NOT NULL,
    transaction_date DATE NOT NULL,
    submission_date DATE NOT NULL,
    SLA ENUM('meet', 'overdue') NOT NULL,
    approved_amount DECIMAL(10, 2) NOT NULL,
    claim_status ENUM('On Plafond', 'Over Plafond') NOT NULL,
    medical_facility_name VARCHAR(255),
    city VARCHAR(255),
    diagnosis VARCHAR(255),
    doc_link VARCHAR(255),
    transaction_status ENUM('Successful', 'Pending', 'Failed') NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    deleted_at DATETIME NULL,
    CONSTRAINT fk_claims_patient
        FOREIGN KEY (patient_id) REFERENCES patients(id)
        ON DELETE RESTRICT -- Claims tidak dihapus jika patient dihapus (biasanya data klaim ingin dipertahankan)
        ON UPDATE CASCADE,
    CONSTRAINT fk_claims_employee
        FOREIGN KEY (employee_id) REFERENCES employees(id)
        ON DELETE RESTRICT -- Claims tidak dihapus jika employee dihapus
        ON UPDATE CASCADE,
    CONSTRAINT fk_claims_benefit
        FOREIGN KEY (benefit_id) REFERENCES benefits(id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE,
    CONSTRAINT fk_claims_transaction_type
        FOREIGN KEY (transaction_type_id) REFERENCES transaction_types(id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE
);
