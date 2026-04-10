package main

import (
	"database/sql"
	shared_models "pharmacy-api/shared/models"

	"gorm.io/gorm"
)

func test_data(db *gorm.DB) {
	// Clear tables (in reverse order due to foreign keys)
	db.Exec("DELETE FROM goods_orders")
	db.Exec("DELETE FROM goods_tags")
	db.Exec("DELETE FROM orders")
	db.Exec("DELETE FROM goods")
	db.Exec("DELETE FROM tags")
	db.Exec("DELETE FROM producers")
	db.Exec("DELETE FROM work_times")

	// WorkTime data
	workTimes := []shared_models.WorkTime{
		{TimeStart: sql.NullString{String: "09:00", Valid: true}, TimeEnd: sql.NullString{String: "18:00", Valid: true}, Date: sql.NullString{String: "2024-01-01", Valid: true}},
		{TimeStart: sql.NullString{String: "09:00", Valid: true}, TimeEnd: sql.NullString{String: "18:00", Valid: true}, Date: sql.NullString{String: "2024-01-02", Valid: true}},
		{TimeStart: sql.NullString{String: "10:00", Valid: true}, TimeEnd: sql.NullString{String: "16:00", Valid: true}, Date: sql.NullString{String: "2024-01-03", Valid: true}},
	}
	db.Create(&workTimes)

	// Producers
	producers := []shared_models.Producer{
		{ProducerName: "Pharmakor"},
		{ProducerName: "Medpreparat"},
		{ProducerName: "Biofarm"},
		{ProducerName: "HealthPlus"},
		{ProducerName: "Pharmacy Factory"},
		{ProducerName: "Pharmacist"},
		{ProducerName: "Medisorb"},
		{ProducerName: "Vitalogica"},
	}
	db.Create(&producers)

	// Tags
	tags := []shared_models.Tag{
		{TagName: "Painkiller"},
		{TagName: "Antipyretic"},
		{TagName: "Antibiotic"},
		{TagName: "Vitamins"},
		{TagName: "Antiseptic"},
		{TagName: "Cough"},
		{TagName: "Stomach"},
		{TagName: "Heart"},
		{TagName: "Allergy"},
		{TagName: "Sedative"},
	}
	db.Create(&tags)

	// Goods (40+ items)
	goods := []shared_models.Goods{
		{
			Name:                 "Paracetamol 500mg",
			Image:                "/goods/pill1.jpg",
			ProducerID:           1,
			IsInStock:            true,
			Instruction:          "Take 1 tablet 3-4 times a day",
			Description:          "Effective antipyretic",
			IsPrescriptionNeeded: false,
			Price:                150,
		},
		{
			Name:                 "Ibuprofen 400mg",
			Image:                "/goods/pill2.jpg",
			ProducerID:           2,
			IsInStock:            true,
			Instruction:          "1 tablet every 6-8 hours",
			Description:          "Non-steroidal anti-inflammatory drug",
			IsPrescriptionNeeded: false,
			Price:                200,
		},
		{
			Name:                 "Amoxicillin 250mg",
			Image:                "/goods/pill3.jpg",
			ProducerID:           3,
			IsInStock:            true,
			Instruction:          "1 capsule 3 times a day",
			Description:          "Broad-spectrum antibiotic",
			IsPrescriptionNeeded: true,
			Price:                350,
		},
		{
			Name:                 "Aspirin 500mg",
			Image:                "/goods/pill4.jpg",
			ProducerID:           1,
			IsInStock:            true,
			Instruction:          "1 tablet as needed",
			Description:          "Painkiller and antipyretic",
			IsPrescriptionNeeded: false,
			Price:                120,
		},
		{
			Name:                 "Vitamin C 1000mg",
			Image:                "/goods/pill5.jpg",
			ProducerID:           4,
			IsInStock:            true,
			Instruction:          "1 tablet per day",
			Description:          "Vitamin supplement for immunity",
			IsPrescriptionNeeded: false,
			Price:                300,
		},
		{
			Name:                 "No-Spa 40mg",
			Image:                "/goods/pill1.jpg",
			ProducerID:           2,
			IsInStock:            true,
			Instruction:          "1-2 tablets 3 times a day",
			Description:          "Antispasmodic",
			IsPrescriptionNeeded: false,
			Price:                180,
		},
		{
			Name:                 "Loratadine 10mg",
			Image:                "/goods/pill2.jpg",
			ProducerID:           5,
			IsInStock:            true,
			Instruction:          "1 tablet per day",
			Description:          "Antiallergic",
			IsPrescriptionNeeded: false,
			Price:                220,
		},
		{
			Name:                 "Validol 60mg",
			Image:                "/goods/pill3.jpg",
			ProducerID:           6,
			IsInStock:            true,
			Instruction:          "1 tablet under the tongue",
			Description:          "Sedative for heart pain",
			IsPrescriptionNeeded: false,
			Price:                80,
		},
		{
			Name:                 "Activated Charcoal",
			Image:                "/goods/pill4.jpg",
			ProducerID:           7,
			IsInStock:            true,
			Instruction:          "4-6 tablets for poisoning",
			Description:          "Adsorbent",
			IsPrescriptionNeeded: false,
			Price:                50,
		},
		{
			Name:                 "Amoxiclav 625mg",
			Image:                "/goods/pill5.jpg",
			ProducerID:           3,
			IsInStock:            false,
			Instruction:          "1 tablet 3 times a day",
			Description:          "Combined antibiotic",
			IsPrescriptionNeeded: true,
			Price:                450,
		},
		{
			Name:                 "Citramon P",
			Image:                "/goods/pill1.jpg",
			ProducerID:           1,
			IsInStock:            true,
			Instruction:          "1 tablet 2-3 times a day",
			Description:          "Combined painkiller",
			IsPrescriptionNeeded: false,
			Price:                90,
		},
		{
			Name:                 "Mukaltin",
			Image:                "/goods/pill2.jpg",
			ProducerID:           4,
			IsInStock:            true,
			Instruction:          "1-2 tablets 3 times a day",
			Description:          "Expectorant",
			IsPrescriptionNeeded: false,
			Price:                60,
		},
		{
			Name:                 "Corvalol",
			Image:                "/goods/pill3.jpg",
			ProducerID:           6,
			IsInStock:            true,
			Instruction:          "15-30 drops per dose",
			Description:          "Sedative and hypnotic",
			IsPrescriptionNeeded: false,
			Price:                110,
		},
		{
			Name:                 "Smecta",
			Image:                "/goods/pill4.jpg",
			ProducerID:           7,
			IsInStock:            true,
			Instruction:          "1 sachet 3 times a day",
			Description:          "Antidiarrheal",
			IsPrescriptionNeeded: false,
			Price:                170,
		},
		{
			Name:                 "Nasivin Spray",
			Image:                "/goods/pill5.jpg",
			ProducerID:           5,
			IsInStock:            true,
			Instruction:          "1 spray 2-3 times a day",
			Description:          "Nasal decongestant",
			IsPrescriptionNeeded: false,
			Price:                190,
		},
		{
			Name:                 "Omeprazole 20mg",
			Image:                "/goods/pill1.jpg",
			ProducerID:           2,
			IsInStock:            true,
			Instruction:          "1 capsule in the morning",
			Description:          "Heartburn relief",
			IsPrescriptionNeeded: false,
			Price:                280,
		},
		{
			Name:                 "Vitamin D3 2000IU",
			Image:                "/goods/pill2.jpg",
			ProducerID:           4,
			IsInStock:            true,
			Instruction:          "1 capsule per day",
			Description:          "Vitamin for bones and immunity",
			IsPrescriptionNeeded: false,
			Price:                320,
		},
		{
			Name:                 "Glycine 100mg",
			Image:                "/goods/pill3.jpg",
			ProducerID:           8,
			IsInStock:            true,
			Instruction:          "1 tablet 2-3 times a day",
			Description:          "Improves cerebral circulation",
			IsPrescriptionNeeded: false,
			Price:                70,
		},
		{
			Name:                 "Fenistil Gel",
			Image:                "/goods/pill4.jpg",
			ProducerID:           5,
			IsInStock:            true,
			Instruction:          "Apply thin layer 2-4 times a day",
			Description:          "Antiallergic gel",
			IsPrescriptionNeeded: false,
			Price:                380,
		},
		{
			Name:                 "Levomycetin",
			Image:                "/goods/pill5.jpg",
			ProducerID:           3,
			IsInStock:            false,
			Instruction:          "1 tablet 3-4 times a day",
			Description:          "Broad-spectrum antibiotic",
			IsPrescriptionNeeded: true,
			Price:                95,
		},
		{
			Name:                 "Bepanthen Cream",
			Image:                "/goods/pill1.jpg",
			ProducerID:           7,
			IsInStock:            true,
			Instruction:          "Apply to damaged skin 1-2 times a day",
			Description:          "Healing cream",
			IsPrescriptionNeeded: false,
			Price:                420,
		},
		{
			Name:                 "Nurofen Express",
			Image:                "/goods/pill2.jpg",
			ProducerID:           2,
			IsInStock:            true,
			Instruction:          "1 capsule up to 3 times a day",
			Description:          "Fast-acting painkiller",
			IsPrescriptionNeeded: false,
			Price:                270,
		},
		{
			Name:                 "Enterosgel",
			Image:                "/goods/pill3.jpg",
			ProducerID:           1,
			IsInStock:            true,
			Instruction:          "1 tablespoon 3 times a day",
			Description:          "Enterosorbent for poisoning",
			IsPrescriptionNeeded: false,
			Price:                380,
		},
		{
			Name:                 "Valerian Drops",
			Image:                "/goods/pill4.jpg",
			ProducerID:           6,
			IsInStock:            true,
			Instruction:          "20-30 drops 3-4 times a day",
			Description:          "Herbal sedative",
			IsPrescriptionNeeded: false,
			Price:                65,
		},
		{
			Name:                 "Suprastin 25mg",
			Image:                "/goods/pill5.jpg",
			ProducerID:           5,
			IsInStock:            true,
			Instruction:          "1 tablet 2-3 times a day",
			Description:          "Antihistamine",
			IsPrescriptionNeeded: true,
			Price:                130,
		},
		{
			Name:                 "Mezim Forte",
			Image:                "/goods/pill1.jpg",
			ProducerID:           4,
			IsInStock:            true,
			Instruction:          "1-2 tablets with meals",
			Description:          "Enzyme supplement",
			IsPrescriptionNeeded: false,
			Price:                290,
		},
		{
			Name:                 "Azithromycin 500mg",
			Image:                "/goods/pill2.jpg",
			ProducerID:           3,
			IsInStock:            true,
			Instruction:          "1 tablet per day for 3 days",
			Description:          "Azalide antibiotic",
			IsPrescriptionNeeded: true,
			Price:                510,
		},
		{
			Name:                 "Nise Gel",
			Image:                "/goods/pill3.jpg",
			ProducerID:           2,
			IsInStock:            true,
			Instruction:          "Apply 3-4 times a day",
			Description:          "Pain relief gel",
			IsPrescriptionNeeded: false,
			Price:                340,
		},
		{
			Name:                 "Kagocel",
			Image:                "/goods/pill4.jpg",
			ProducerID:           1,
			IsInStock:            true,
			Instruction:          "According to scheme 2+2+2 tablets",
			Description:          "Antiviral",
			IsPrescriptionNeeded: false,
			Price:                480,
		},
		{
			Name:                 "Linex",
			Image:                "/goods/pill5.jpg",
			ProducerID:           7,
			IsInStock:            true,
			Instruction:          "2 capsules 3 times a day",
			Description:          "Probiotic for intestines",
			IsPrescriptionNeeded: false,
			Price:                520,
		},
		{
			Name:                 "Afobazol 10mg",
			Image:                "/goods/pill1.jpg",
			ProducerID:           8,
			IsInStock:            true,
			Instruction:          "1 tablet 3 times a day",
			Description:          "Anti-anxiety",
			IsPrescriptionNeeded: false,
			Price:                410,
		},
		{
			Name:                 "Furacilin",
			Image:                "/goods/pill2.jpg",
			ProducerID:           4,
			IsInStock:            true,
			Instruction:          "Dissolve 1 tablet in glass of water",
			Description:          "Antiseptic for gargling",
			IsPrescriptionNeeded: false,
			Price:                75,
		},
		{
			Name:                 "Teraflu",
			Image:                "/goods/pill3.jpg",
			ProducerID:           2,
			IsInStock:            true,
			Instruction:          "1 sachet per cup of hot water",
			Description:          "Cold remedy",
			IsPrescriptionNeeded: false,
			Price:                360,
		},
		{
			Name:                 "Hexoral Spray",
			Image:                "/goods/pill4.jpg",
			ProducerID:           5,
			IsInStock:            true,
			Instruction:          "2 sprays 2 times a day",
			Description:          "Throat antiseptic",
			IsPrescriptionNeeded: false,
			Price:                290,
		},
		{
			Name:                 "Lazolvan Syrup",
			Image:                "/goods/pill5.jpg",
			ProducerID:           6,
			IsInStock:            true,
			Instruction:          "5 ml 3 times a day",
			Description:          "Mucolytic",
			IsPrescriptionNeeded: false,
			Price:                330,
		},
		{
			Name:                 "Vitrum Vitamins",
			Image:                "/goods/pill1.jpg",
			ProducerID:           4,
			IsInStock:            true,
			Instruction:          "1 tablet per day",
			Description:          "Vitamin and mineral complex",
			IsPrescriptionNeeded: false,
			Price:                890,
		},
		{
			Name:                 "Ceftriaxone",
			Image:                "/goods/pill2.jpg",
			ProducerID:           3,
			IsInStock:            false,
			Instruction:          "For injections only",
			Description:          "Cephalosporin antibiotic",
			IsPrescriptionNeeded: true,
			Price:                120,
		},
		{
			Name:                 "Panthenol Spray",
			Image:                "/goods/pill3.jpg",
			ProducerID:           7,
			IsInStock:            true,
			Instruction:          "Spray on affected area",
			Description:          "Burn treatment",
			IsPrescriptionNeeded: false,
			Price:                280,
		},
		{
			Name:                 "Analgin 500mg",
			Image:                "/goods/pill4.jpg",
			ProducerID:           1,
			IsInStock:            true,
			Instruction:          "1 tablet 2-3 times a day",
			Description:          "Painkiller and antipyretic",
			IsPrescriptionNeeded: false,
			Price:                85,
		},
		{
			Name:                 "Essentiale Forte",
			Image:                "/goods/pill5.jpg",
			ProducerID:           8,
			IsInStock:            true,
			Instruction:          "2 capsules 3 times a day",
			Description:          "Hepatoprotector",
			IsPrescriptionNeeded: false,
			Price:                670,
		},
	}
	db.Create(&goods)

	// Goods-Tags relationships
	goodsTags := [][]int{
		{1, 1}, {1, 2}, // Paracetamol - painkiller, antipyretic
		{2, 1}, {2, 2}, // Ibuprofen - painkiller, antipyretic
		{3, 3},         // Amoxicillin - antibiotic
		{4, 1}, {4, 2}, // Aspirin - painkiller, antipyretic
		{5, 4},          // Vitamin C - vitamins
		{6, 7},          // No-Spa - stomach
		{7, 9},          // Loratadine - allergy
		{8, 8}, {8, 10}, // Validol - heart, sedative
		{9, 7},   // Charcoal - stomach
		{10, 3},  // Amoxiclav - antibiotic
		{11, 1},  // Citramon - painkiller
		{12, 6},  // Mukaltin - cough
		{13, 10}, // Corvalol - sedative
		{14, 7},  // Smecta - stomach
		{15, 5},  // Nasivin - antiseptic
		{16, 7},  // Omeprazole - stomach
		{17, 4},  // Vitamin D - vitamins
		{18, 10}, // Glycine - sedative
		{19, 9},  // Fenistil - allergy
		{20, 3},  // Levomycetin - antibiotic
	}

	for _, gt := range goodsTags {
		if len(gt) == 2 {
			db.Exec("INSERT INTO goods_tags (goods_id, tag_id) VALUES (?, ?)", gt[0], gt[1])
		}
	}

	// Orders
	orders := []shared_models.Order{
		{ClientFIO: "Ivanov Ivan Ivanovich", ClientEmail: "ivanov@mail.ru", ClientPhone: "+79161234567"},
		{ClientFIO: "Petrova Maria Sergeevna", ClientEmail: "petrova@gmail.com", ClientPhone: "+79037654321"},
		{ClientFIO: "Sidorov Alexey Petrovich", ClientEmail: "", ClientPhone: "+79219876543"},
	}
	db.Create(&orders)

	// GoodsOrders
	goodsOrders := []shared_models.GoodsOrders{
		{OrderID: 1, GoodsID: 1, Quantity: 2},
		{OrderID: 1, GoodsID: 5, Quantity: 1},
		{OrderID: 2, GoodsID: 2, Quantity: 1},
		{OrderID: 2, GoodsID: 7, Quantity: 2},
		{OrderID: 3, GoodsID: 3, Quantity: 1},
		{OrderID: 3, GoodsID: 8, Quantity: 1},
		{OrderID: 3, GoodsID: 12, Quantity: 3},
	}
	db.Create(&goodsOrders)

}
