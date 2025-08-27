package seed

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/thoriqwildan/aino-medical-be/internal/entity"
	"gorm.io/gorm"
)

func ptrFloat64(v float64) *float64 { return &v }
func ptrString(s string) *string    { return &s }

func ptrDate(date time.Time) *time.Time { return &date }

// SeedBenefits populates the database with benefit data based on predefined plans.
func SeedBenefits(db *gorm.DB) {

	// 2) Ensure PlanTypes for grades A-D exist.
	grades := []string{"PLAN A", "PLAN B", "PLAN C", "PLAN D"}
	planTypeIDs := map[string]uint{}
	for _, g := range grades {
		var pt entity.PlanType
		if err := db.Where("name = ?", g).First(&pt).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				pt = entity.PlanType{Name: g}
				if err := db.Create(&pt).Error; err != nil {
					log.Printf("create plan type %s: %v", g, err)
					continue
				}
			} else {
				log.Printf("query plan type %s: %v", g, err)
				continue
			}
		}
		planTypeIDs[g] = pt.ID
	}

	yearlyClaimCodes := []*entity.YearlyBenefitClaim{
		{
			Code:        "RI-PLANB",
			YearlyClaim: 75000000,
		},
		{
			Code:        "RI-PLANC",
			YearlyClaim: 60000000,
		},
		{
			Code:        "RI-PLAND",
			YearlyClaim: 50000000,
		},
		{
			Code:        "RJ-PLANA",
			YearlyClaim: 12000000,
		},
		{
			Code:        "RJ-PLANB",
			YearlyClaim: 9000000,
		},
		{
			Code:        "RJ-PLANC",
			YearlyClaim: 8000000,
		},
		{
			Code:        "RJ-PLAND",
			YearlyClaim: 6000000,
		},
		{
			Code:        "RG-PLANA",
			YearlyClaim: 5000000,
		},
		{
			Code:        "RG-PLANB",
			YearlyClaim: 4000000,
		},
		{
			Code:        "RG-PLANC",
			YearlyClaim: 3000000,
		},
		{
			Code:        "RG-PLAND",
			YearlyClaim: 2500000,
		},
	}
	yearlyClaimCodeIDs := map[string]uint{}
	for _, yearlyClaim := range yearlyClaimCodes {
		var yc entity.YearlyBenefitClaim
		if err := db.Where("code = ?", yearlyClaim.Code).First(&yc).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				yc = *yearlyClaim
				if err := db.Create(&yc).Error; err != nil {
					log.Printf("create yearly benefit claim %s: %v", yc.Code, err)
					continue
				}
			} else {
				log.Printf("query yearly benefit claim %s: %v", yc.Code, err)
				continue
			}
		}
		yearlyClaimCodeIDs[yc.Code] = yc.ID
	}

	// Helper struct to define benefit data rows.
	type row struct {
		Code                string
		Name                string
		Detail              *string
		Limitation          entity.LimitationType
		YearlyBenefitsClaim map[string]*string
		Plafond             map[string]*float64 // by PlanType Name (GRADE X). nil => as charged
		YearlyMax           map[string]*float64 // by PlanType Name (GRADE X).
	}

	// Data sourced from "Tabel Benefit 2025.xlsx"
	rows := []row{
		{Code: "RI-1", Name: "Rawat Inap", Detail: ptrString("Biaya Kamar dan Makan (Maks. 365 hari per kasus)"), Limitation: entity.LimitationTypePerDay,
			YearlyMax: map[string]*float64{
				"PLAN A": nil,
				"PLAN B": nil,
				"PLAN C": nil,
				"PLAN D": nil,
			},
			YearlyBenefitsClaim: map[string]*string{
				"PLAN A": nil,
				"PLAN B": ptrString("RI-PLANB"),
				"PLAN C": ptrString("RI-PLANC"),
				"PLAN D": ptrString("RI-PLAND"),
			},
			Plafond: map[string]*float64{
				"PLAN A": ptrFloat64(1500000),
				"PLAN B": ptrFloat64(1000000),
				"PLAN C": ptrFloat64(650000),
				"PLAN D": ptrFloat64(500000),
			},
		},
		{Code: "RI-2", Name: "Rawat Inap", Detail: ptrString("Biaya Perawatan di Rumah Sakit"), Limitation: entity.LimitationTypePerIncident,
			YearlyMax: map[string]*float64{
				"PLAN A": nil,
				"PLAN B": nil,
				"PLAN C": nil,
				"PLAN D": nil,
			}, YearlyBenefitsClaim: map[string]*string{
				"PLAN A": nil,
				"PLAN B": ptrString("RI-PLANB"),
				"PLAN C": ptrString("RI-PLANC"),
				"PLAN D": ptrString("RI-PLAND"),
			},
			Plafond: map[string]*float64{
				"PLAN A": ptrFloat64(20000000),
				"PLAN B": ptrFloat64(13000000),
				"PLAN C": ptrFloat64(10000000),
				"PLAN D": ptrFloat64(8000000),
			},
		},
		{Code: "RI-3", Name: "Rawat Inap", Detail: ptrString("Biaya Kamar Semi ICU dan ICU (Maks. 365 hari per kasus)"), Limitation: entity.LimitationTypePerDay,
			YearlyMax: map[string]*float64{
				"PLAN A": nil,
				"PLAN B": nil,
				"PLAN C": nil,
				"PLAN D": nil,
			},
			YearlyBenefitsClaim: map[string]*string{
				"PLAN A": nil,
				"PLAN B": ptrString("RI-PLANB"),
				"PLAN C": ptrString("RI-PLANC"),
				"PLAN D": ptrString("RI-PLAND"),
			},
			Plafond: map[string]*float64{
				"PLAN A": ptrFloat64(2000000),
				"PLAN B": ptrFloat64(1500000),
				"PLAN C": ptrFloat64(1000000),
				"PLAN D": ptrFloat64(800000),
			},
		},
		{Code: "RI-4", Name: "Rawat Inap", Detail: ptrString("Operasi Kompleks: Biaya Operasi (Termasuk Dokter Bedah, Kamar Operasi dan Anestesi)"), Limitation: entity.LimitationTypePerIncident,
			YearlyMax: map[string]*float64{
				"PLAN A": nil,
				"PLAN B": nil,
				"PLAN C": nil,
				"PLAN D": nil,
			},
			YearlyBenefitsClaim: map[string]*string{
				"PLAN A": nil,
				"PLAN B": ptrString("RI-PLANB"),
				"PLAN C": ptrString("RI-PLANC"),
				"PLAN D": ptrString("RI-PLAND"),
			},
			Plafond: map[string]*float64{
				"PLAN A": ptrFloat64(70000000),
				"PLAN B": ptrFloat64(49000000),
				"PLAN C": ptrFloat64(28000000),
				"PLAN D": ptrFloat64(24500000),
			},
		},
		{Code: "RI-5", Name: "Rawat Inap", Detail: ptrString("Operasi Besar: Biaya Operasi (Termasuk Dokter Bedah, Kamar Operasi dan Anestesi)"), Limitation: entity.LimitationTypePerIncident,
			YearlyMax: map[string]*float64{
				"PLAN A": nil,
				"PLAN B": nil,
				"PLAN C": nil,
				"PLAN D": nil,
			},
			YearlyBenefitsClaim: map[string]*string{
				"PLAN A": nil,
				"PLAN B": ptrString("RI-PLANB"),
				"PLAN C": ptrString("RI-PLANC"),
				"PLAN D": ptrString("RI-PLAND"),
			},
			Plafond: map[string]*float64{
				"PLAN A": ptrFloat64(35000000),
				"PLAN B": ptrFloat64(24500000),
				"PLAN C": ptrFloat64(14000000),
				"PLAN D": ptrFloat64(12250000),
			},
		},
		{Code: "RI-6", Name: "Rawat Inap", Detail: ptrString("Operasi Sedang: Biaya Operasi (Termasuk Dokter Bedah, Kamar Operasi dan Anestesi)"), Limitation: entity.LimitationTypePerIncident,
			YearlyMax: map[string]*float64{
				"PLAN A": nil,
				"PLAN B": nil,
				"PLAN C": nil,
				"PLAN D": nil,
			},
			YearlyBenefitsClaim: map[string]*string{
				"PLAN A": nil,
				"PLAN B": ptrString("RI-PLANB"),
				"PLAN C": ptrString("RI-PLANC"),
				"PLAN D": ptrString("RI-PLAND"),
			},
			Plafond: map[string]*float64{
				"PLAN A": ptrFloat64(22400000),
				"PLAN B": ptrFloat64(15680000),
				"PLAN C": ptrFloat64(8960000),
				"PLAN D": ptrFloat64(7840000),
			},
		},
		{Code: "RI-7", Name: "Rawat Inap", Detail: ptrString("Operasi Kecil: Biaya Operasi (Termasuk Dokter Bedah, Kamar Operasi dan Anestesi)"), Limitation: entity.LimitationTypePerIncident,
			YearlyMax: map[string]*float64{
				"PLAN A": nil,
				"PLAN B": nil,
				"PLAN C": nil,
				"PLAN D": nil,
			},
			YearlyBenefitsClaim: map[string]*string{
				"PLAN A": nil,
				"PLAN B": ptrString("RI-PLANB"),
				"PLAN C": ptrString("RI-PLANC"),
				"PLAN D": ptrString("RI-PLAND"),
			},
			Plafond: map[string]*float64{
				"PLAN A": ptrFloat64(15400000),
				"PLAN B": ptrFloat64(10780000),
				"PLAN C": ptrFloat64(6160000),
				"PLAN D": ptrFloat64(5390000),
			},
		},
		{Code: "RI-8", Name: "Rawat Inap", Detail: ptrString("Biaya Kunjungan Dokter (Maks. 365 hari per kasus)"), Limitation: entity.LimitationTypePerDay,
			YearlyMax: map[string]*float64{
				"PLAN A": nil,
				"PLAN B": nil,
				"PLAN C": nil,
				"PLAN D": nil,
			},
			YearlyBenefitsClaim: map[string]*string{
				"PLAN A": nil,
				"PLAN B": ptrString("RI-PLANB"),
				"PLAN C": ptrString("RI-PLANC"),
				"PLAN D": ptrString("RI-PLAND"),
			},
			Plafond: map[string]*float64{
				"PLAN A": ptrFloat64(500000),
				"PLAN B": ptrFloat64(350000),
				"PLAN C": ptrFloat64(200000),
				"PLAN D": ptrFloat64(175000),
			},
		},
		{Code: "RI-9", Name: "Rawat Inap", Detail: ptrString("Biaya Konsultasi dengan Dokter Spesialis"), Limitation: entity.LimitationTypePerIncident,
			YearlyMax: map[string]*float64{
				"PLAN A": nil,
				"PLAN B": nil,
				"PLAN C": nil,
				"PLAN D": nil,
			},
			YearlyBenefitsClaim: map[string]*string{
				"PLAN A": nil,
				"PLAN B": ptrString("RI-PLANB"),
				"PLAN C": ptrString("RI-PLANC"),
				"PLAN D": ptrString("RI-PLAND"),
			},
			Plafond: map[string]*float64{
				"PLAN A": ptrFloat64(600000),
				"PLAN B": ptrFloat64(450000),
				"PLAN C": ptrFloat64(300000),
				"PLAN D": ptrFloat64(275000),
			},
		},
		{Code: "RI-10", Name: "Rawat Inap", Detail: ptrString("Biaya Ambulan"), Limitation: entity.LimitationTypePerIncident,
			YearlyMax: map[string]*float64{
				"PLAN A": nil,
				"PLAN B": nil,
				"PLAN C": nil,
				"PLAN D": nil,
			},
			YearlyBenefitsClaim: map[string]*string{
				"PLAN A": nil,
				"PLAN B": ptrString("RI-PLANB"),
				"PLAN C": ptrString("RI-PLANC"),
				"PLAN D": ptrString("RI-PLAND"),
			},
			Plafond: map[string]*float64{
				"PLAN A": ptrFloat64(500000),
				"PLAN B": ptrFloat64(400000),
				"PLAN C": ptrFloat64(300000),
				"PLAN D": ptrFloat64(250000),
			},
		},
		{Code: "RJ-1", Name: "Rawat Jalan", Detail: ptrString("Biaya Konsultasi Dokter Umum"), Limitation: entity.LimitationTypePerIncident,
			YearlyMax: map[string]*float64{
				"PLAN A": nil,
				"PLAN B": nil,
				"PLAN C": nil,
				"PLAN D": nil,
			},
			YearlyBenefitsClaim: map[string]*string{
				"PLAN A": ptrString("RJ-PLANA"),
				"PLAN B": ptrString("RJ-PLANB"),
				"PLAN C": ptrString("RJ-PLANC"),
				"PLAN D": ptrString("RJ-PLAND"),
			},
			Plafond: map[string]*float64{
				"PLAN A": ptrFloat64(140000),
				"PLAN B": ptrFloat64(130000),
				"PLAN C": ptrFloat64(120000),
				"PLAN D": ptrFloat64(100000),
			},
		},
		{Code: "RJ-2", Name: "Rawat Jalan", Detail: ptrString("Biaya Konsultasi Dokter Spesialis (Tanpa Surat Pengantar)"), Limitation: entity.LimitationTypePerIncident,
			YearlyMax: map[string]*float64{
				"PLAN A": nil,
				"PLAN B": nil,
				"PLAN C": nil,
				"PLAN D": nil,
			},
			YearlyBenefitsClaim: map[string]*string{
				"PLAN A": ptrString("RJ-PLANA"),
				"PLAN B": ptrString("RJ-PLANB"),
				"PLAN C": ptrString("RJ-PLANC"),
				"PLAN D": ptrString("RJ-PLAND"),
			},
			Plafond: map[string]*float64{
				"PLAN A": ptrFloat64(240000),
				"PLAN B": ptrFloat64(220000),
				"PLAN C": ptrFloat64(200000),
				"PLAN D": ptrFloat64(160000),
			},
		},
		{Code: "RJ-3", Name: "Rawat Jalan", Detail: ptrString("Biaya Konsultasi Dokter dan Obat-obatan"), Limitation: entity.LimitationTypePerIncident,
			YearlyMax: map[string]*float64{
				"PLAN A": nil,
				"PLAN B": nil,
				"PLAN C": nil,
				"PLAN D": nil,
			},
			YearlyBenefitsClaim: map[string]*string{
				"PLAN A": ptrString("RJ-PLANA"),
				"PLAN B": ptrString("RJ-PLANB"),
				"PLAN C": ptrString("RJ-PLANC"),
				"PLAN D": ptrString("RJ-PLAND"),
			},
			Plafond: map[string]*float64{
				"PLAN A": ptrFloat64(140000),
				"PLAN B": ptrFloat64(130000),
				"PLAN C": ptrFloat64(120000),
				"PLAN D": ptrFloat64(100000),
			},
		},
		{Code: "RJ-4", Name: "Rawat Jalan", Detail: ptrString("Biaya Pembelian Obat-obatan Sesuai dengan Resep Dokter"), Limitation: entity.LimitationTypePerYear,
			YearlyMax: map[string]*float64{
				"PLAN A": ptrFloat64(2540000),
				"PLAN B": ptrFloat64(2040000),
				"PLAN C": ptrFloat64(1540000),
				"PLAN D": ptrFloat64(1040000),
			},
			YearlyBenefitsClaim: map[string]*string{
				"PLAN A": ptrString("RJ-PLANA"),
				"PLAN B": ptrString("RJ-PLANB"),
				"PLAN C": ptrString("RJ-PLANC"),
				"PLAN D": ptrString("RJ-PLAND"),
			},
			Plafond: map[string]*float64{
				"PLAN A": nil,
				"PLAN B": nil,
				"PLAN C": nil,
				"PLAN D": nil,
			},
		},
		{Code: "RJ-5", Name: "Rawat Jalan", Detail: ptrString("Biaya Pemeriksaan Laboratorium"), Limitation: entity.LimitationTypePerYear,
			YearlyMax: map[string]*float64{
				"PLAN A": ptrFloat64(2540000),
				"PLAN B": ptrFloat64(2040000),
				"PLAN C": ptrFloat64(1790000),
				"PLAN D": ptrFloat64(1390000),
			},
			YearlyBenefitsClaim: map[string]*string{
				"PLAN A": ptrString("RJ-PLANA"),
				"PLAN B": ptrString("RJ-PLANB"),
				"PLAN C": ptrString("RJ-PLANC"),
				"PLAN D": ptrString("RJ-PLAND"),
			},
			Plafond: map[string]*float64{
				"PLAN A": nil,
				"PLAN B": nil,
				"PLAN C": nil,
				"PLAN D": nil,
			},
		},
		{Code: "RJ-6", Name: "Rawat Jalan", Detail: ptrString("Biaya Fisioterapi"), Limitation: entity.LimitationTypePerYear,
			YearlyMax: map[string]*float64{
				"PLAN A": ptrFloat64(1540000),
				"PLAN B": ptrFloat64(1340000),
				"PLAN C": ptrFloat64(1040000),
				"PLAN D": ptrFloat64(790000),
			},
			YearlyBenefitsClaim: map[string]*string{
				"PLAN A": ptrString("RJ-PLANA"),
				"PLAN B": ptrString("RJ-PLANB"),
				"PLAN C": ptrString("RJ-PLANC"),
				"PLAN D": ptrString("RJ-PLAND"),
			},
			Plafond: map[string]*float64{
				"PLAN A": nil,
				"PLAN B": nil,
				"PLAN C": nil,
				"PLAN D": nil,
			},
		},
		{Code: "RJ-7", Name: "Rawat Jalan", Detail: ptrString("Biaya Imunisasi Dasar untuk Anak s/d 5 Tahun"), Limitation: entity.LimitationTypePerYear,
			YearlyMax: map[string]*float64{
				"PLAN A": ptrFloat64(840000),
				"PLAN B": ptrFloat64(790000),
				"PLAN C": ptrFloat64(640000),
				"PLAN D": ptrFloat64(540000),
			},
			YearlyBenefitsClaim: map[string]*string{
				"PLAN A": ptrString("RJ-PLANA"),
				"PLAN B": ptrString("RJ-PLANB"),
				"PLAN C": ptrString("RJ-PLANC"),
				"PLAN D": ptrString("RJ-PLAND"),
			},
			Plafond: map[string]*float64{
				"PLAN A": nil,
				"PLAN B": nil,
				"PLAN C": nil,
				"PLAN D": nil,
			},
		},
		{Code: "RG-1", Name: "Rawat Gigi", Detail: ptrString("Biaya Perawatan Dasar"), Limitation: entity.LimitationTypePerYear,
			YearlyMax: map[string]*float64{
				"PLAN A": nil,
				"PLAN B": nil,
				"PLAN C": nil,
				"PLAN D": nil,
			},
			YearlyBenefitsClaim: map[string]*string{
				"PLAN A": ptrString("RG-PLANA"),
				"PLAN B": ptrString("RG-PLANB"),
				"PLAN C": ptrString("RG-PLANC"),
				"PLAN D": ptrString("RG-PLAND"),
			},
			Plafond: map[string]*float64{
				"PLAN A": nil,
				"PLAN B": nil,
				"PLAN C": nil,
				"PLAN D": nil,
			},
		},
		{Code: "RG-2", Name: "Rawat Gigi", Detail: ptrString("Biaya Perawatan Gusi"), Limitation: entity.LimitationTypePerYear,
			YearlyMax: map[string]*float64{
				"PLAN A": nil,
				"PLAN B": nil,
				"PLAN C": nil,
				"PLAN D": nil,
			},
			YearlyBenefitsClaim: map[string]*string{
				"PLAN A": ptrString("RG-PLANA"),
				"PLAN B": ptrString("RG-PLANB"),
				"PLAN C": ptrString("RG-PLANC"),
				"PLAN D": ptrString("RG-PLAND"),
			},
			Plafond: map[string]*float64{
				"PLAN A": nil,
				"PLAN B": nil,
				"PLAN C": nil,
				"PLAN D": nil,
			},
		},
		{Code: "RG-3", Name: "Rawat Gigi", Detail: ptrString("Biaya Perawatan Pencegahan"), Limitation: entity.LimitationTypePerYear,
			YearlyMax: map[string]*float64{
				"PLAN A": nil,
				"PLAN B": nil,
				"PLAN C": nil,
				"PLAN D": nil,
			},
			YearlyBenefitsClaim: map[string]*string{
				"PLAN A": ptrString("RG-PLANA"),
				"PLAN B": ptrString("RG-PLANB"),
				"PLAN C": ptrString("RG-PLANC"),
				"PLAN D": ptrString("RG-PLAND"),
			},
			Plafond: map[string]*float64{
				"PLAN A": nil,
				"PLAN B": nil,
				"PLAN C": nil,
				"PLAN D": nil,
			},
		},
		{Code: "RG-4", Name: "Rawat Gigi", Detail: ptrString("Biaya Perawatan Kompleks"), Limitation: entity.LimitationTypePerYear,
			YearlyMax: map[string]*float64{
				"PLAN A": nil,
				"PLAN B": nil,
				"PLAN C": nil,
				"PLAN D": nil,
			},
			YearlyBenefitsClaim: map[string]*string{
				"PLAN A": ptrString("RG-PLANA"),
				"PLAN B": ptrString("RG-PLANB"),
				"PLAN C": ptrString("RG-PLANC"),
				"PLAN D": ptrString("RG-PLAND"),
			},
			Plafond: map[string]*float64{
				"PLAN A": nil,
				"PLAN B": nil,
				"PLAN C": nil,
				"PLAN D": nil,
			},
		},
		{Code: "RG-5", Name: "Rawat Gigi", Detail: ptrString("Biaya Perawatan Perbaikan"), Limitation: entity.LimitationTypePerYear,
			YearlyMax: map[string]*float64{
				"PLAN A": nil,
				"PLAN B": nil,
				"PLAN C": nil,
				"PLAN D": nil,
			},
			YearlyBenefitsClaim: map[string]*string{
				"PLAN A": ptrString("RG-PLANA"),
				"PLAN B": ptrString("RG-PLANB"),
				"PLAN C": ptrString("RG-PLANC"),
				"PLAN D": ptrString("RG-PLAND"),
			},
			Plafond: map[string]*float64{
				"PLAN A": nil,
				"PLAN B": nil,
				"PLAN C": nil,
				"PLAN D": nil,
			},
		},
		{Code: "RG-6", Name: "Rawat Gigi", Detail: ptrString("Biaya Gigi Palsu"), Limitation: entity.LimitationTypePerYear,
			YearlyMax: map[string]*float64{
				"PLAN A": nil,
				"PLAN B": nil,
				"PLAN C": nil,
				"PLAN D": nil,
			},
			YearlyBenefitsClaim: map[string]*string{
				"PLAN A": ptrString("RG-PLANA"),
				"PLAN B": ptrString("RG-PLANB"),
				"PLAN C": ptrString("RG-PLANC"),
				"PLAN D": ptrString("RG-PLAND"),
			},
			Plafond: map[string]*float64{
				"PLAN A": nil,
				"PLAN B": nil,
				"PLAN C": nil,
				"PLAN D": nil,
			},
		},
		{Code: "RM-1", Name: "Melahirkan", Detail: ptrString("Biaya Melahirkan Normal"), Limitation: entity.LimitationTypePerPregnancy,
			YearlyMax: map[string]*float64{
				"PLAN A": nil,
				"PLAN B": nil,
				"PLAN C": nil,
				"PLAN D": nil,
			},
			Plafond: map[string]*float64{
				"PLAN A": ptrFloat64(10000000),
				"PLAN B": ptrFloat64(8500000),
				"PLAN C": ptrFloat64(7000000),
				"PLAN D": ptrFloat64(6750000),
			},
		},
		{Code: "RM-2", Name: "Melahirkan", Detail: ptrString("Biaya Melahirkan dengan Pembedahan (Caesar)"), Limitation: entity.LimitationTypePerPregnancy,
			YearlyMax: map[string]*float64{
				"PLAN A": nil,
				"PLAN B": nil,
				"PLAN C": nil,
				"PLAN D": nil,
			},
			Plafond: map[string]*float64{
				"PLAN A": ptrFloat64(15000000),
				"PLAN B": ptrFloat64(12000000),
				"PLAN C": ptrFloat64(9000000),
				"PLAN D": ptrFloat64(8500000),
			},
		},
		{Code: "RM-3", Name: "Melahirkan", Detail: ptrString("Biaya Pengguguran Kehamilan (Aborsi) Atas Pertimbangan Medis"), Limitation: entity.LimitationTypePerPregnancy,
			YearlyMax: map[string]*float64{
				"PLAN A": nil,
				"PLAN B": nil,
				"PLAN C": nil,
				"PLAN D": nil,
			},
			Plafond: map[string]*float64{
				"PLAN A": ptrFloat64(8500000),
				"PLAN B": ptrFloat64(7500000),
				"PLAN C": ptrFloat64(6500000),
				"PLAN D": ptrFloat64(6230000),
			},
		},
		{Code: "RM-4", Name: "Melahirkan", Detail: ptrString("Biaya Komplikasi Kehamilan dan Komplikasi Pasca Melahirkan"), Limitation: entity.LimitationTypePerYear,
			YearlyMax: map[string]*float64{
				"PLAN A": ptrFloat64(8500000),
				"PLAN B": ptrFloat64(7500000),
				"PLAN C": ptrFloat64(6500000),
				"PLAN D": ptrFloat64(6230000),
			},
			Plafond: map[string]*float64{
				"PLAN A": nil,
				"PLAN B": nil,
				"PLAN C": nil,
				"PLAN D": nil,
			},
		},
		{Code: "RM-5", Name: "Melahirkan", Detail: ptrString("Biaya Perawatan Sebelum dan 40 Hari Setelah Melahirkan"), Limitation: entity.LimitationTypePerYear,
			YearlyMax: map[string]*float64{
				"PLAN A": ptrFloat64(7500000),
				"PLAN B": ptrFloat64(6750000),
				"PLAN C": ptrFloat64(6000000),
				"PLAN D": ptrFloat64(5880000),
			},
			Plafond: map[string]*float64{
				"PLAN A": nil,
				"PLAN B": nil,
				"PLAN C": nil,
				"PLAN D": nil,
			},
		},
		{Code: "RK-1", Name: "Kacamata", Detail: ptrString("Biaya Bingkai Kacamata"), Limitation: entity.LimitationTypePerYear,
			YearlyMax: map[string]*float64{
				"PLAN A": ptrFloat64(1250000),
				"PLAN B": ptrFloat64(1000000),
				"PLAN C": ptrFloat64(750000),
				"PLAN D": ptrFloat64(500000),
			},
			Plafond: map[string]*float64{
				"PLAN A": nil,
				"PLAN B": nil,
				"PLAN C": nil,
				"PLAN D": nil,
			},
		},
		{Code: "RK-2", Name: "Kacamata", Detail: ptrString("Biaya Lensa Kacamata atau Lensa Kontak"), Limitation: entity.LimitationTypePerYear,
			YearlyMax: map[string]*float64{
				"PLAN A": ptrFloat64(1250000),
				"PLAN B": ptrFloat64(1000000),
				"PLAN C": ptrFloat64(750000),
				"PLAN D": ptrFloat64(500000),
			},
			Plafond: map[string]*float64{
				"PLAN A": nil,
				"PLAN B": nil,
				"PLAN C": nil,
				"PLAN D": nil,
			},
		},
		{Code: "RK-3", Name: "Kacamata", Detail: ptrString("Biaya Refraksi Mata / Konsultasi Dokter Mata"), Limitation: entity.LimitationTypePerYear,
			YearlyMax: map[string]*float64{
				"PLAN A": ptrFloat64(250000),
				"PLAN B": ptrFloat64(200000),
				"PLAN C": ptrFloat64(150000),
				"PLAN D": ptrFloat64(150000),
			},
			Plafond: map[string]*float64{
				"PLAN A": nil,
				"PLAN B": nil,
				"PLAN C": nil,
				"PLAN D": nil,
			},
		},
		{Code: "MCU", Name: "Medical Check Up", Detail: ptrString("Biaya Medical Check Up"), Limitation: entity.LimitationTypePerYear,
			YearlyMax: map[string]*float64{
				"PLAN A": ptrFloat64(1000000),
				"PLAN B": ptrFloat64(800000),
				"PLAN C": ptrFloat64(600000),
				"PLAN D": ptrFloat64(500000),
			},
			Plafond: map[string]*float64{
				"PLAN A": nil,
				"PLAN B": nil,
				"PLAN C": nil,
				"PLAN D": nil,
			},
		},
	}
	// 3) Upsert benefits per grade
	for _, rw := range rows {
		for planName, ptID := range planTypeIDs {
			plaf, ok := rw.Plafond[planName]
			if !ok {
				continue
			}

			var yearlyMax *float64
			if rw.YearlyMax != nil {
				if ym, ok := rw.YearlyMax[planName]; ok {
					yearlyMax = ym
				}
			}

			b := entity.Benefit{
				Name:           rw.Name,
				Detail:         rw.Detail,
				Code:           fmt.Sprintf("%s-%s", strings.ReplaceAll(planName, " ", ""), rw.Code), // e.g., GRADEA-RI-1
				PlanTypeID:     ptID,
				LimitationType: rw.Limitation,
				YearlyMax:      yearlyMax,
				Plafond:        plaf,
			}

			if rw.YearlyBenefitsClaim[planName] != nil {
				yearlyBenefitClaimId, ok := yearlyClaimCodeIDs[*rw.YearlyBenefitsClaim[planName]]
				if ok {
					b.YearlyBenefitClaimID = &yearlyBenefitClaimId
				}
			}

			var existing entity.Benefit
			err := db.Where("code = ?", b.Code).First(&existing).Error
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					if err := db.Create(&b).Error; err != nil {
						log.Printf("create benefit %s: %v", b.Code, err)
					}
				} else {
					log.Printf("query benefit %s: %v", b.Code, err)
				}
			} else {
				// Update existing record
				updateData := map[string]interface{}{
					"Name":           b.Name,
					"Detail":         b.Detail,
					"PlanTypeID":     b.PlanTypeID,
					"LimitationType": b.LimitationType,
					"Plafond":        b.Plafond,
					"YearlyMax":      b.YearlyMax,
				}
				if err := db.Model(&existing).Updates(updateData).Error; err != nil {
					log.Printf("update benefit %s: %v", existing.Code, err)
				}
			}
		}
	}
}
