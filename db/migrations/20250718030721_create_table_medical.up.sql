CREATE TABLE IF NOT EXISTS plan_types (
        id INT PRIMARY KEY AUTO_INCREMENT,
        name VARCHAR(255) UNIQUE NOT NULL,
        description TEXT
    );

CREATE TABLE IF NOT EXISTS transaction_types (
         id INT PRIMARY KEY AUTO_INCREMENT,
         name VARCHAR(255) UNIQUE NOT NULL
    );



CREATE TABLE IF NOT EXISTS departments (
        id INT PRIMARY KEY AUTO_INCREMENT,
        name VARCHAR(255) UNIQUE NOT NULL,
        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        deleted_at TIMESTAMP NULL
    );

CREATE TABLE IF NOT EXISTS employees (
        id INT PRIMARY KEY AUTO_INCREMENT,
        name VARCHAR(255) NOT NULL,
        department_id INT NOT NULL,
        position VARCHAR(255) NOT NULL,
        email VARCHAR(255) UNIQUE NOT NULL,
        phone VARCHAR(255) NOT NULL,
        birth_date DATE NOT NULL,
        gender ENUM('male', 'female') NOT NULL,
        plan_type_id INT NOT NULL,
        dependence VARCHAR(255),
        bank_number VARCHAR(255) NOT NULL,
        pro_rate DECIMAL(5,2) UNSIGNED NOT NULL,
        join_date DATE NOT NULL,
        CONSTRAINT fk_employees_department
        FOREIGN KEY (department_id) REFERENCES departments(id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE,
        CONSTRAINT fk_employees_plan_type
        FOREIGN KEY (plan_type_id) REFERENCES plan_types(id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE
    );

CREATE TABLE IF NOT EXISTS family_members (
        id INT PRIMARY KEY AUTO_INCREMENT,
        employee_id INT NOT NULL,
        name VARCHAR(255) NOT NULL,
        plan_type_id INT NOT NULL,
        birth_date DATE NOT NULL,
        relationship_type ENUM('wife', 'husband', 'father', 'mother', 'child') NOT NULL,
        gender ENUM('male', 'female') NOT NULL,
        CONSTRAINT fk_family_members_employee
        FOREIGN KEY (employee_id) REFERENCES employees(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
        CONSTRAINT fk_family_members_plan_type
        FOREIGN KEY (plan_type_id) REFERENCES plan_types(id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE
    );

CREATE TABLE IF NOT EXISTS patients (
        id INT PRIMARY KEY AUTO_INCREMENT,
        name VARCHAR(255) NOT NULL,
        birth_date DATE NOT NULL,
        gender ENUM('male', 'female') NOT NULL,
        employee_id INT UNIQUE,
        family_member_id INT UNIQUE,
        plan_type_id INT NOT NULL,
        CONSTRAINT fk_patients_employee_cascade
        FOREIGN KEY (employee_id) REFERENCES employees(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
        CONSTRAINT fk_patients_family_member_cascade
        FOREIGN KEY (family_member_id) REFERENCES family_members(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
        CONSTRAINT fk_patients_plan_type
        FOREIGN KEY (plan_type_id) REFERENCES plan_types(id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE
    );


CREATE TABLE IF NOT EXISTS yearly_benefit_claims (
    id INT PRIMARY KEY AUTO_INCREMENT,
    code VARCHAR(255) UNIQUE,
    yearly_claim DECIMAL(18, 2) NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
    );

CREATE TABLE IF NOT EXISTS benefits (
        id INT PRIMARY KEY AUTO_INCREMENT,
        name VARCHAR(255) NOT NULL,
        plan_type_id INT NOT NULL,
        yearly_benefit_claim_id INT NULL,
        detail TEXT,
        code VARCHAR(255) UNIQUE NOT NULL,
        limitation_type ENUM('Per Day', 'Per Month', 'Per Year', 'Per Pregnancy', 'Per Incident') NOT NULL,
        plafond DECIMAL(18, 2) NULL,
        yearly_max DECIMAL(18, 2) NULL,
        CONSTRAINT fk_benefits_plan_type
        FOREIGN KEY (plan_type_id) REFERENCES plan_types(id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE,
        CONSTRAINT fk_yearly_benefit_claims
        FOREIGN KEY (yearly_benefit_claim_id) REFERENCES yearly_benefit_claims(id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE
    );

CREATE TABLE IF NOT EXISTS patient_benefits (
        id INT PRIMARY KEY AUTO_INCREMENT,
        patient_id INT NOT NULL,
        benefit_id INT NOT NULL,
        yearly_max DECIMAL(18, 2) NULL,
        remaining_plafond DECIMAL(18, 2),
        initial_plafond DECIMAL(18, 2),
        start_date DATE NOT NULL,
        end_date DATE,
        status ENUM('active', 'exhausted', 'expired') DEFAULT 'active',
        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        FOREIGN KEY (patient_id) REFERENCES patients(id),
        FOREIGN KEY (benefit_id) REFERENCES benefits(id),
        UNIQUE (patient_id, benefit_id)
    );

CREATE TABLE IF NOT EXISTS claims (
        id INT PRIMARY KEY AUTO_INCREMENT,
        patient_benefit_id INT NOT NULL,
        patient_id INT NOT NULL,
        employee_id INT NOT NULL,
        claim_amount DECIMAL(10, 2) NOT NULL,
        transaction_type_id INT NULL,
        transaction_date DATE NULL,
        submission_date DATE NULL,
        SLA ENUM('meet', 'overdue') NULL,
        approved_amount DECIMAL(10, 2) NULL,
        claim_status ENUM('On Plafond', 'Over Plafond') NOT NULL,
        medical_facility_name VARCHAR(255),
        city VARCHAR(255),
        diagnosis VARCHAR(255),
        doc_link VARCHAR(255),
        transaction_status ENUM('Successful', 'Pending', 'Failed') NOT NULL,
        created_at DATETIME NOT NULL,
        updated_at DATETIME NULL,
        deleted_at DATETIME NULL,
        CONSTRAINT fk_claims_patient_benefit
        FOREIGN KEY (patient_benefit_id) REFERENCES patient_benefits(id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE,
        CONSTRAINT fk_claims_patient
        FOREIGN KEY (patient_id) REFERENCES patients(id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE,
        CONSTRAINT fk_claims_employee
        FOREIGN KEY (employee_id) REFERENCES employees(id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE,
        CONSTRAINT fk_claims_transaction_type
        FOREIGN KEY (transaction_type_id) REFERENCES transaction_types(id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE
    );
