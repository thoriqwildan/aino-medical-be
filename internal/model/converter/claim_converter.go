package converter

import (
	"time" // Pastikan time diimport untuk time.Time{}

	"github.com/thoriqwildan/aino-medical-be/internal/entity"
	"github.com/thoriqwildan/aino-medical-be/internal/helper"
	"github.com/thoriqwildan/aino-medical-be/internal/model"
)

func ClaimToResponse(claim *entity.Claim) *model.ClaimResponse {
	// Inisialisasi variabel untuk menampung nilai yang sudah di-cek nil
	var transactionDate helper.CustomDate
	var submissionDate helper.CustomDate
	var slaStatus string
	var approvedAmount float64
	var medicalFacility string
	var city string
	var diagnosis string
	var docLink string
	var updatedAt *time.Time // UpdatedAt di entity adalah *time.Time, jadi di response bisa juga pointer atau nil

	// Lakukan pengecekan nil untuk setiap pointer
	if claim.TransactionDate != nil {
		transactionDate = helper.CustomDate(*claim.TransactionDate)
	} else {
		// Jika nil, berikan nilai default atau kosong
		transactionDate = helper.CustomDate(time.Time{}) // Atau sesuaikan default yang Anda inginkan
	}

	if claim.SubmissionDate != nil {
		submissionDate = helper.CustomDate(*claim.SubmissionDate)
	} else {
		submissionDate = helper.CustomDate(time.Time{})
	}

	if claim.SLA != nil {
		slaStatus = string(*claim.SLA)
	} else {
		slaStatus = "" // Atau nilai default lainnya untuk SLAStatus
	}

	if claim.ApprovedAmount != nil {
		approvedAmount = *claim.ApprovedAmount
	} else {
		approvedAmount = 0.0 // Atau nilai default lainnya untuk ApprovedAmount
	}

	if claim.MedicalFacilityName != nil {
		medicalFacility = *claim.MedicalFacilityName
	} else {
		medicalFacility = "" // String kosong jika nil
	}

	if claim.City != nil {
		city = *claim.City
	} else {
		city = ""
	}

	if claim.Diagnosis != nil {
		diagnosis = *claim.Diagnosis
	} else {
		diagnosis = ""
	}

	if claim.DocLink != nil {
		docLink = *claim.DocLink
	} else {
		docLink = ""
	}

	if claim.UpdatedAt != nil {
		updatedAt = claim.UpdatedAt // Cukup assign langsung karena sudah *time.Time
	} else {
		updatedAt = nil // Tetap nil jika dari DB nil
	}


	result := &model.ClaimResponse{
		ID:                claim.ID,
		ClaimAmount:       claim.ClaimAmount,
		TransactionDate:   transactionDate,
		SubmissionDate:    submissionDate,
		SLAStatus:         slaStatus,
		ApprovedAmount:    approvedAmount,
		ClaimStatus:       string(claim.ClaimStatus),       // ClaimStatus bukan pointer, aman
		MedicalFacility:   medicalFacility,
		City:              city,
		Diagnosis:         diagnosis,
		DocLink:           docLink,
		TransactionStatus: string(claim.TransactionStatus), // TransactionStatus bukan pointer, aman
		CreatedAt:         claim.CreatedAt,                 // CreatedAt bukan pointer, aman
		UpdatedAt:         updatedAt,                       // Menggunakan variabel yang sudah di-cek nil
	}

	// Bagian ini sudah cukup aman karena ada nil check-nya
	if claim.TransactionType != nil {
		result.TransactionType = *TransactionTypeToResponse(claim.TransactionType)
	}

	if claim.Patient.ID != 0 { // Patient bukan pointer, cek ID 0 adalah cara aman untuk struct
		result.Patient = model.PatientResponse{
			ID:        claim.Patient.ID,
			Name:      claim.Patient.Name,
			BirthDate: helper.CustomDate(claim.Patient.BirthDate),
			Gender:    string(claim.Patient.Gender),
		}
	}

	if claim.PatientBenefit.BenefitID != 0 { // PatientBenefit bukan pointer, cek BenefitID 0 adalah cara aman
		result.Benefit = *BenefitToResponse(&claim.PatientBenefit.Benefit)
	}

	return result
}

func PatientToResponse(patient *entity.Patient) *model.PatientResponse {
	return &model.PatientResponse{
		ID:        patient.ID,
		Name:      patient.Name,
		BirthDate: helper.CustomDate(patient.BirthDate),
		Gender: string(patient.Gender),
		PlanType: *PlanTypeToResponse(&patient.PlanType),
	}
}