package seed

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/thoriqwildan/aino-medical-be/internal/entity"
	"gorm.io/gorm"
)

func ptrFloat64(v float64) *float64 { return &v }
func ptrString(s string) *string    { return &s }

// SeedBenefits populates the database with benefit data based on predefined plans.
func SeedBenefits(db *gorm.DB) {
	// 1) Build limitation type map from DB (use records seeded by SeedLimitationTypes)
	var lts []entity.LimitationType
	if err := db.Find(&lts).Error; err != nil {
		log.Printf("failed to load limitation types: %v", err)
		return
	}
	limID := map[string]uint{}
	for _, lt := range lts {
		limID[strings.ToLower(lt.Name)] = lt.ID
	}
	lookupLim := func(name string) uint {
		if id, ok := limID[strings.ToLower(strings.TrimSpace(name))]; ok {
			return id
		}
		// fallback: create if missing
		lt := entity.LimitationType{Name: name}
		if err := db.Create(&lt).Error; err != nil {
			log.Printf("create limitation type %s: %v", name, err)
			return 0
		}
		limID[strings.ToLower(name)] = lt.ID
		return lt.ID
	}

	// 2) Ensure PlanTypes for grades A-D exist.
	grades := []string{"GRADE A", "GRADE B", "GRADE C", "GRADE D"}
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

	// Helper struct to define benefit data rows.
	type row struct {
		Code       string
		Name       string
		Detail     *string
		Limitation string
		Plafond    map[string]*float64 // by PlanType Name (GRADE X). nil => as charged
		YearlyMax  map[string]*float64 // by PlanType Name (GRADE X).
	}

	// Data sourced from "Tabel Benefit 2025.xlsx"
	rows := []row{
		{Code: "RI-1", Name: "Rawat Inap", Detail: ptrString("Biaya Kamar dan Makan (Maks. 365 hari per kasus)"), Limitation: "Per Day",
			YearlyMax: map[string]*float64{
				"GRADE A": nil,
				"GRADE B": ptrFloat64(75000000),
				"GRADE C": ptrFloat64(60000000),
				"GRADE D": ptrFloat64(50000000),
			},
			Plafond: map[string]*float64{
				"GRADE A": ptrFloat64(1500000),
				"GRADE B": ptrFloat64(1000000),
				"GRADE C": ptrFloat64(650000),
				"GRADE D": ptrFloat64(500000),
			},
		},
		{Code: "RI-2", Name: "Rawat Inap", Detail: ptrString("Biaya Perawatan di Rumah Sakit"), Limitation: "Per Incident",
			YearlyMax: map[string]*float64{
				"GRADE A": nil,
				"GRADE B": ptrFloat64(75000000),
				"GRADE C": ptrFloat64(60000000),
				"GRADE D": ptrFloat64(50000000),
			},
			Plafond: map[string]*float64{
				"GRADE A": ptrFloat64(20000000),
				"GRADE B": ptrFloat64(13000000),
				"GRADE C": ptrFloat64(10000000),
				"GRADE D": ptrFloat64(8000000),
			},
		},
		{Code: "RI-3", Name: "Rawat Inap", Detail: ptrString("Biaya Kamar Semi ICU dan ICU (Maks. 365 hari per kasus)"), Limitation: "Per Day",
			YearlyMax: map[string]*float64{
				"GRADE A": nil,
				"GRADE B": ptrFloat64(75000000),
				"GRADE C": ptrFloat64(60000000),
				"GRADE D": ptrFloat64(50000000),
			},
			Plafond: map[string]*float64{
				"GRADE A": ptrFloat64(2000000),
				"GRADE B": ptrFloat64(1500000),
				"GRADE C": ptrFloat64(1000000),
				"GRADE D": ptrFloat64(800000),
			},
		},
		{Code: "RI-4", Name: "Rawat Inap", Detail: ptrString("Operasi Kompleks: Biaya Operasi (Termasuk Dokter Bedah, Kamar Operasi dan Anestesi)"), Limitation: "Per Incident",
			YearlyMax: map[string]*float64{
				"GRADE A": nil,
				"GRADE B": ptrFloat64(75000000),
				"GRADE C": ptrFloat64(60000000),
				"GRADE D": ptrFloat64(50000000),
			},
			Plafond: map[string]*float64{
				"GRADE A": ptrFloat64(70000000),
				"GRADE B": ptrFloat64(49000000),
				"GRADE C": ptrFloat64(28000000),
				"GRADE D": ptrFloat64(24500000),
			},
		},
		{Code: "RI-5", Name: "Rawat Inap", Detail: ptrString("Operasi Besar: Biaya Operasi (Termasuk Dokter Bedah, Kamar Operasi dan Anestesi)"), Limitation: "Per Incident",
			YearlyMax: map[string]*float64{
				"GRADE A": nil,
				"GRADE B": ptrFloat64(75000000),
				"GRADE C": ptrFloat64(60000000),
				"GRADE D": ptrFloat64(50000000),
			},
			Plafond: map[string]*float64{
				"GRADE A": ptrFloat64(35000000),
				"GRADE B": ptrFloat64(24500000),
				"GRADE C": ptrFloat64(14000000),
				"GRADE D": ptrFloat64(12250000),
			},
		},
		{Code: "RI-6", Name: "Rawat Inap", Detail: ptrString("Operasi Sedang: Biaya Operasi (Termasuk Dokter Bedah, Kamar Operasi dan Anestesi)"), Limitation: "Per Incident",
			YearlyMax: map[string]*float64{
				"GRADE A": nil,
				"GRADE B": ptrFloat64(75000000),
				"GRADE C": ptrFloat64(60000000),
				"GRADE D": ptrFloat64(50000000),
			},
			Plafond: map[string]*float64{
				"GRADE A": ptrFloat64(22400000),
				"GRADE B": ptrFloat64(15680000),
				"GRADE C": ptrFloat64(8960000),
				"GRADE D": ptrFloat64(7840000),
			},
		},
		{Code: "RI-7", Name: "Rawat Inap", Detail: ptrString("Operasi Kecil: Biaya Operasi (Termasuk Dokter Bedah, Kamar Operasi dan Anestesi)"), Limitation: "Per Incident",
			YearlyMax: map[string]*float64{
				"GRADE A": nil,
				"GRADE B": ptrFloat64(75000000),
				"GRADE C": ptrFloat64(60000000),
				"GRADE D": ptrFloat64(50000000),
			},
			Plafond: map[string]*float64{
				"GRADE A": ptrFloat64(15400000),
				"GRADE B": ptrFloat64(10780000),
				"GRADE C": ptrFloat64(6160000),
				"GRADE D": ptrFloat64(5390000),
			},
		},
		{Code: "RI-8", Name: "Rawat Inap", Detail: ptrString("Biaya Kunjungan Dokter (Maks. 365 hari per kasus)"), Limitation: "Per Day",
			YearlyMax: map[string]*float64{
				"GRADE A": nil,
				"GRADE B": ptrFloat64(75000000),
				"GRADE C": ptrFloat64(60000000),
				"GRADE D": ptrFloat64(50000000),
			},
			Plafond: map[string]*float64{
				"GRADE A": ptrFloat64(500000),
				"GRADE B": ptrFloat64(350000),
				"GRADE C": ptrFloat64(200000),
				"GRADE D": ptrFloat64(175000),
			},
		},
		{Code: "RI-9", Name: "Rawat Inap", Detail: ptrString("Biaya Konsultasi dengan Dokter Spesialis"), Limitation: "Per Incident",
			YearlyMax: map[string]*float64{
				"GRADE A": nil,
				"GRADE B": ptrFloat64(75000000),
				"GRADE C": ptrFloat64(60000000),
				"GRADE D": ptrFloat64(50000000),
			},
			Plafond: map[string]*float64{
				"GRADE A": ptrFloat64(600000),
				"GRADE B": ptrFloat64(450000),
				"GRADE C": ptrFloat64(300000),
				"GRADE D": ptrFloat64(275000),
			},
		},
		{Code: "RI-10", Name: "Rawat Inap", Detail: ptrString("Biaya Ambulan"), Limitation: "Per Incident",
			YearlyMax: map[string]*float64{
				"GRADE A": nil,
				"GRADE B": ptrFloat64(75000000),
				"GRADE C": ptrFloat64(60000000),
				"GRADE D": ptrFloat64(50000000),
			},
			Plafond: map[string]*float64{
				"GRADE A": ptrFloat64(500000),
				"GRADE B": ptrFloat64(400000),
				"GRADE C": ptrFloat64(300000),
				"GRADE D": ptrFloat64(250000),
			},
		},
		{Code: "RJ-1", Name: "Rawat Jalan", Detail: ptrString("Biaya Konsultasi Dokter Umum"), Limitation: "Per Incident",
			YearlyMax: map[string]*float64{
				"GRADE A": ptrFloat64(12000000),
				"GRADE B": ptrFloat64(75000000),
				"GRADE C": ptrFloat64(60000000),
				"GRADE D": ptrFloat64(50000000),
			},
			Plafond: map[string]*float64{
				"GRADE A": ptrFloat64(140000),
				"GRADE B": ptrFloat64(130000),
				"GRADE C": ptrFloat64(120000),
				"GRADE D": ptrFloat64(100000),
			},
		},
		{Code: "RJ-2", Name: "Rawat Jalan", Detail: ptrString("Biaya Konsultasi Dokter Spesialis (Tanpa Surat Pengantar)"), Limitation: "Per Incident",
			YearlyMax: map[string]*float64{
				"GRADE A": ptrFloat64(12000000),
				"GRADE B": ptrFloat64(75000000),
				"GRADE C": ptrFloat64(60000000),
				"GRADE D": ptrFloat64(50000000),
			},
			Plafond: map[string]*float64{
				"GRADE A": ptrFloat64(240000),
				"GRADE B": ptrFloat64(220000),
				"GRADE C": ptrFloat64(200000),
				"GRADE D": ptrFloat64(160000),
			},
		},
		{Code: "RJ-3", Name: "Rawat Jalan", Detail: ptrString("Biaya Konsultasi Dokter dan Obat-obatan"), Limitation: "Per Incident",
			YearlyMax: map[string]*float64{
				"GRADE A": ptrFloat64(12000000),
				"GRADE B": ptrFloat64(75000000),
				"GRADE C": ptrFloat64(60000000),
				"GRADE D": ptrFloat64(50000000),
			},
			Plafond: map[string]*float64{
				"GRADE A": ptrFloat64(140000),
				"GRADE B": ptrFloat64(130000),
				"GRADE C": ptrFloat64(120000),
				"GRADE D": ptrFloat64(100000),
			},
		},
		{Code: "RJ-4", Name: "Rawat Jalan", Detail: ptrString("Biaya Pembelian Obat-obatan Sesuai dengan Resep Dokter"), Limitation: "Per Year",
			YearlyMax: map[string]*float64{
				"GRADE A": ptrFloat64(12000000),
				"GRADE B": ptrFloat64(75000000),
				"GRADE C": ptrFloat64(60000000),
				"GRADE D": ptrFloat64(50000000),
			},
			Plafond: map[string]*float64{
				"GRADE A": ptrFloat64(2540000),
				"GRADE B": ptrFloat64(2040000),
				"GRADE C": ptrFloat64(1540000),
				"GRADE D": ptrFloat64(1040000),
			},
		},
		{Code: "RJ-5", Name: "Rawat Jalan", Detail: ptrString("Biaya Pemeriksaan Laboratorium"), Limitation: "Per Year",
			YearlyMax: map[string]*float64{
				"GRADE A": ptrFloat64(12000000),
				"GRADE B": ptrFloat64(75000000),
				"GRADE C": ptrFloat64(60000000),
				"GRADE D": ptrFloat64(50000000),
			},
			Plafond: map[string]*float64{
				"GRADE A": ptrFloat64(2540000),
				"GRADE B": ptrFloat64(2040000),
				"GRADE C": ptrFloat64(1790000),
				"GRADE D": ptrFloat64(1390000),
			},
		},
		{Code: "RJ-6", Name: "Rawat Jalan", Detail: ptrString("Biaya Fisioterapi"), Limitation: "Per Year",
			YearlyMax: map[string]*float64{
				"GRADE A": ptrFloat64(12000000),
				"GRADE B": ptrFloat64(75000000),
				"GRADE C": ptrFloat64(60000000),
				"GRADE D": ptrFloat64(50000000),
			},
			Plafond: map[string]*float64{
				"GRADE A": ptrFloat64(1540000),
				"GRADE B": ptrFloat64(1340000),
				"GRADE C": ptrFloat64(1040000),
				"GRADE D": ptrFloat64(790000),
			},
		},
		{Code: "RJ-7", Name: "Rawat Jalan", Detail: ptrString("Biaya Imunisasi Dasar untuk Anak s/d 5 Tahun"), Limitation: "Per Year",
			YearlyMax: map[string]*float64{
				"GRADE A": ptrFloat64(12000000),
				"GRADE B": ptrFloat64(75000000),
				"GRADE C": ptrFloat64(60000000),
				"GRADE D": ptrFloat64(50000000),
			},
			Plafond: map[string]*float64{
				"GRADE A": ptrFloat64(840000),
				"GRADE B": ptrFloat64(790000),
				"GRADE C": ptrFloat64(640000),
				"GRADE D": ptrFloat64(540000),
			},
		},
		{Code: "RG-1", Name: "Rawat Gigi", Detail: ptrString("Biaya Perawatan Dasar"), Limitation: "Per Year",
			YearlyMax: map[string]*float64{
				"GRADE A": ptrFloat64(12000000),
				"GRADE B": ptrFloat64(9000000),
				"GRADE C": ptrFloat64(8000000),
				"GRADE D": ptrFloat64(6000000),
			},
			Plafond: map[string]*float64{
				"GRADE A": nil,
				"GRADE B": nil,
				"GRADE C": nil,
				"GRADE D": nil,
			},
		},
		{Code: "RG-2", Name: "Rawat Gigi", Detail: ptrString("Biaya Perawatan Gusi"), Limitation: "Per Year",
			YearlyMax: map[string]*float64{
				"GRADE A": ptrFloat64(12000000),
				"GRADE B": ptrFloat64(9000000),
				"GRADE C": ptrFloat64(8000000),
				"GRADE D": ptrFloat64(6000000),
			},
			Plafond: map[string]*float64{
				"GRADE A": nil,
				"GRADE B": nil,
				"GRADE C": nil,
				"GRADE D": nil,
			},
		},
		{Code: "RG-3", Name: "Rawat Gigi", Detail: ptrString("Biaya Perawatan Pencegahan"), Limitation: "Per Year",
			YearlyMax: map[string]*float64{
				"GRADE A": ptrFloat64(12000000),
				"GRADE B": ptrFloat64(9000000),
				"GRADE C": ptrFloat64(8000000),
				"GRADE D": ptrFloat64(6000000),
			},
			Plafond: map[string]*float64{
				"GRADE A": nil,
				"GRADE B": nil,
				"GRADE C": nil,
				"GRADE D": nil,
			},
		},
		{Code: "RG-4", Name: "Rawat Gigi", Detail: ptrString("Biaya Perawatan Kompleks"), Limitation: "Per Year",
			YearlyMax: map[string]*float64{
				"GRADE A": ptrFloat64(12000000),
				"GRADE B": ptrFloat64(9000000),
				"GRADE C": ptrFloat64(8000000),
				"GRADE D": ptrFloat64(6000000),
			},
			Plafond: map[string]*float64{
				"GRADE A": nil,
				"GRADE B": nil,
				"GRADE C": nil,
				"GRADE D": nil,
			},
		},
		{Code: "RG-5", Name: "Rawat Gigi", Detail: ptrString("Biaya Perawatan Perbaikan"), Limitation: "Per Year",
			YearlyMax: map[string]*float64{
				"GRADE A": ptrFloat64(12000000),
				"GRADE B": ptrFloat64(9000000),
				"GRADE C": ptrFloat64(8000000),
				"GRADE D": ptrFloat64(6000000),
			},
			Plafond: map[string]*float64{
				"GRADE A": nil,
				"GRADE B": nil,
				"GRADE C": nil,
				"GRADE D": nil,
			},
		},
		{Code: "RG-6", Name: "Rawat Gigi", Detail: ptrString("Biaya Gigi Palsu"), Limitation: "Per Year",
			YearlyMax: map[string]*float64{
				"GRADE A": ptrFloat64(12000000),
				"GRADE B": ptrFloat64(9000000),
				"GRADE C": ptrFloat64(8000000),
				"GRADE D": ptrFloat64(6000000),
			},
			Plafond: map[string]*float64{
				"GRADE A": nil,
				"GRADE B": nil,
				"GRADE C": nil,
				"GRADE D": nil,
			},
		},
		{Code: "RM-1", Name: "Melahirkan", Detail: ptrString("Biaya Melahirkan Normal"), Limitation: "Per Pregnancy",
			YearlyMax: map[string]*float64{
				"GRADE A": nil,
				"GRADE B": nil,
				"GRADE C": nil,
				"GRADE D": nil,
			},
			Plafond: map[string]*float64{
				"GRADE A": ptrFloat64(10000000),
				"GRADE B": ptrFloat64(8500000),
				"GRADE C": ptrFloat64(7000000),
				"GRADE D": ptrFloat64(6750000),
			},
		},
		{Code: "RM-2", Name: "Melahirkan", Detail: ptrString("Biaya Melahirkan dengan Pembedahan (Caesar)"), Limitation: "Per Pregnancy",
			YearlyMax: map[string]*float64{
				"GRADE A": nil,
				"GRADE B": nil,
				"GRADE C": nil,
				"GRADE D": nil,
			},
			Plafond: map[string]*float64{
				"GRADE A": ptrFloat64(15000000),
				"GRADE B": ptrFloat64(12000000),
				"GRADE C": ptrFloat64(9000000),
				"GRADE D": ptrFloat64(8500000),
			},
		},
		{Code: "RM-3", Name: "Melahirkan", Detail: ptrString("Biaya Pengguguran Kehamilan (Aborsi) Atas Pertimbangan Medis"), Limitation: "Per Pregnancy",
			YearlyMax: map[string]*float64{
				"GRADE A": nil,
				"GRADE B": nil,
				"GRADE C": nil,
				"GRADE D": nil,
			},
			Plafond: map[string]*float64{
				"GRADE A": ptrFloat64(8500000),
				"GRADE B": ptrFloat64(7500000),
				"GRADE C": ptrFloat64(6500000),
				"GRADE D": ptrFloat64(6230000),
			},
		},
		{Code: "RM-4", Name: "Melahirkan", Detail: ptrString("Biaya Komplikasi Kehamilan dan Komplikasi Pasca Melahirkan"), Limitation: "Per Year",
			YearlyMax: map[string]*float64{
				"GRADE A": nil,
				"GRADE B": nil,
				"GRADE C": nil,
				"GRADE D": nil,
			},
			Plafond: map[string]*float64{
				"GRADE A": ptrFloat64(8500000),
				"GRADE B": ptrFloat64(7500000),
				"GRADE C": ptrFloat64(6500000),
				"GRADE D": ptrFloat64(6230000),
			},
		},
		{Code: "RM-5", Name: "Melahirkan", Detail: ptrString("Biaya Perawatan Sebelum dan 40 Hari Setelah Melahirkan"), Limitation: "Per Year",
			YearlyMax: map[string]*float64{
				"GRADE A": nil,
				"GRADE B": nil,
				"GRADE C": nil,
				"GRADE D": nil,
			},
			Plafond: map[string]*float64{
				"GRADE A": ptrFloat64(7500000),
				"GRADE B": ptrFloat64(6750000),
				"GRADE C": ptrFloat64(6000000),
				"GRADE D": ptrFloat64(5880000),
			},
		},
		{Code: "RK-1", Name: "Kacamata", Detail: ptrString("Biaya Bingkai Kacamata"), Limitation: "Per Year",
			YearlyMax: map[string]*float64{
				"GRADE A": nil,
				"GRADE B": nil,
				"GRADE C": nil,
				"GRADE D": nil,
			},
			Plafond: map[string]*float64{
				"GRADE A": ptrFloat64(1250000),
				"GRADE B": ptrFloat64(1000000),
				"GRADE C": ptrFloat64(750000),
				"GRADE D": ptrFloat64(500000),
			},
		},
		{Code: "RK-2", Name: "Kacamata", Detail: ptrString("Biaya Lensa Kacamata atau Lensa Kontak"), Limitation: "Per Year",
			YearlyMax: map[string]*float64{
				"GRADE A": nil,
				"GRADE B": nil,
				"GRADE C": nil,
				"GRADE D": nil,
			},
			Plafond: map[string]*float64{
				"GRADE A": ptrFloat64(1250000),
				"GRADE B": ptrFloat64(1000000),
				"GRADE C": ptrFloat64(750000),
				"GRADE D": ptrFloat64(500000),
			},
		},
		{Code: "RK-3", Name: "Kacamata", Detail: ptrString("Biaya Refraksi Mata / Konsultasi Dokter Mata"), Limitation: "Per Year",
			YearlyMax: map[string]*float64{
				"GRADE A": nil,
				"GRADE B": nil,
				"GRADE C": nil,
				"GRADE D": nil,
			},
			Plafond: map[string]*float64{
				"GRADE A": ptrFloat64(250000),
				"GRADE B": ptrFloat64(200000),
				"GRADE C": ptrFloat64(150000),
				"GRADE D": ptrFloat64(150000),
			},
		},
		{Code: "MCU", Name: "Medical Check Up", Detail: ptrString("Biaya Medical Check Up"), Limitation: "Per Year",
			YearlyMax: map[string]*float64{
				"GRADE A": ptrFloat64(3000000),
				"GRADE B": ptrFloat64(3200000),
				"GRADE C": ptrFloat64(5400000),
				"GRADE D": ptrFloat64(12500000),
			},
			Plafond: map[string]*float64{
				"GRADE A": ptrFloat64(1000000),
				"GRADE B": ptrFloat64(800000),
				"GRADE C": ptrFloat64(600000),
				"GRADE D": ptrFloat64(500000),
			},
		},
	}
	// 3) Upsert benefits per grade
	for _, rw := range rows {
		ltID := lookupLim(rw.Limitation)
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
				Name:             rw.Name,
				Detail:           rw.Detail,
				Code:             fmt.Sprintf("%s-%s", rw.Code, strings.ReplaceAll(planName, " ", "")), // e.g., RI-1-GRADEA
				PlanTypeID:       ptID,
				LimitationTypeID: ltID,
				YearlyMax:        yearlyMax,
				Plafond:          plaf,
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
					"Name":             b.Name,
					"Detail":           b.Detail,
					"PlanTypeID":       b.PlanTypeID,
					"LimitationTypeID": b.LimitationTypeID,
					"Plafond":          b.Plafond,
					"YearlyMax":        b.YearlyMax,
				}
				if err := db.Model(&existing).Updates(updateData).Error; err != nil {
					log.Printf("update benefit %s: %v", existing.Code, err)
				}
			}
		}
	}
}
