
CREATE TABLE plan_types (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) UNIQUE NOT NULL,
    description TEXT
);

CREATE TABLE transaction_types (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) UNIQUE NOT NULL
);

CREATE TABLE limitation_types (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) UNIQUE NOT NULL
);

CREATE TABLE departments (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE TABLE patients (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    birth_date DATE NOT NULL,
    gender ENUM('male', 'female') NOT NULL
);

CREATE TABLE employees (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    patient_id INT UNIQUE NOT NULL,
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
    FOREIGN KEY (patient_id) REFERENCES patients(id),
    FOREIGN KEY (department_id) REFERENCES departments(id),
    FOREIGN KEY (plan_type_id) REFERENCES plan_types(id)
);

CREATE TABLE family_members (
    id INT PRIMARY KEY AUTO_INCREMENT,
    patient_id INT UNIQUE NOT NULL,
    employee_id INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    plan_type_id INT NOT NULL,
    birth_date DATE NOT NULL,
    gender ENUM('male', 'female') NOT NULL,
    FOREIGN KEY (patient_id) REFERENCES patients(id),
    FOREIGN KEY (employee_id) REFERENCES employees(id),
    FOREIGN KEY (plan_type_id) REFERENCES plan_types(id)
);

CREATE TABLE benefits (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    plan_type_id INT NOT NULL,
    detail TEXT,
    code VARCHAR(255) UNIQUE NOT NULL,
    limitation_type_id INT NOT NULL,
    plafond INT NOT NULL,
    yearly_max INT NOT NULL,
    FOREIGN KEY (plan_type_id) REFERENCES plan_types(id),
    FOREIGN KEY (limitation_type_id) REFERENCES limitation_types(id)
);

CREATE TABLE claims (
    id INT PRIMARY KEY AUTO_INCREMENT,
    patient_id INT NOT NULL,
    employee_id INT NOT NULL,
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
    FOREIGN KEY (patient_id) REFERENCES patients(id),
    FOREIGN KEY (employee_id) REFERENCES employees(id),
    FOREIGN KEY (benefit_id) REFERENCES benefits(id),
    FOREIGN KEY (transaction_type_id) REFERENCES transaction_types(id)
);