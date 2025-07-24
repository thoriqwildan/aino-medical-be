CREATE TABLE patient_benefits (
    id INT PRIMARY KEY AUTO_INCREMENT,
    patient_id INT NOT NULL,
    benefit_id INT NOT NULL,
    remaining_plafond DECIMAL(10, 2) NOT NULL,
    initial_plafond DECIMAL(10, 2) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE,
    status ENUM('active', 'exhausted', 'expired') DEFAULT 'active',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (patient_id) REFERENCES patients(id),
    FOREIGN KEY (benefit_id) REFERENCES benefits(id),
    UNIQUE (patient_id, benefit_id)
);

ALTER TABLE claims
ADD COLUMN patient_benefit_id INT NOT NULL AFTER benefit_id,
ADD CONSTRAINT fk_claims_patient_benefit
FOREIGN KEY (patient_benefit_id) REFERENCES patient_benefits(id);