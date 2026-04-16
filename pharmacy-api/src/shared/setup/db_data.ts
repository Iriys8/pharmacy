import { DataSource } from 'typeorm';

export async function testData(dataSource: DataSource): Promise<void> {
  const queryRunner = dataSource.createQueryRunner();
  
  try {
    await queryRunner.connect();
    await queryRunner.startTransaction();
    
    // 1. Добавление расписаний
    const schedules = [
      { time_start: "09:00", time_end: "18:00", date: "2026-04-10", is_opened: true },
      { time_start: "09:00", time_end: "18:00", date: "2026-04-11", is_opened: true },
      { time_start: "10:00", time_end: "16:00", date: "2026-04-12", is_opened: true },
    ];
    
    for (const schedule of schedules) {
      await queryRunner.query(
        "INSERT INTO schedules (timeStart, timeEnd, date, isOpened) VALUES (?, ?, ?, ?)",
        [schedule.time_start, schedule.time_end, schedule.date, schedule.is_opened]
      );
    }
    
    // 2. Добавление производителей
    const producers = [
      { producer_name: "Pharmakor" },
      { producer_name: "Medpreparat" },
      { producer_name: "Biofarm" },
      { producer_name: "HealthPlus" },
      { producer_name: "Pharmacy Factory" },
      { producer_name: "Pharmacist" },
      { producer_name: "Medisorb" },
      { producer_name: "Vitalogica" },
    ];
    
    const producerIds: number[] = [];
    for (const producer of producers) {
      const result = await queryRunner.query(
        "INSERT INTO producers (producerName) VALUES (?)",
        [producer.producer_name]
      );
      producerIds.push(result.insertId);
    }
    
    // 3. Добавление тегов
    const tags = [
      { tag_name: "Painkiller" },
      { tag_name: "Antipyretic" },
      { tag_name: "Antibiotic" },
      { tag_name: "Vitamins" },
      { tag_name: "Antiseptic" },
      { tag_name: "Cough" },
      { tag_name: "Stomach" },
      { tag_name: "Heart" },
      { tag_name: "Allergy" },
      { tag_name: "Sedative" },
      { tag_name: "Anti-inflammatory" },
      { tag_name: "Probiotic" },
      { tag_name: "Antiviral" },
      { tag_name: "Enzyme" },
      { tag_name: "Hepatoprotector" },
    ];
    
    const tagIds: number[] = [];
    for (const tag of tags) {
      const result = await queryRunner.query(
        "INSERT INTO tags (tagName) VALUES (?)",
        [tag.tag_name]
      );
      tagIds.push(result.insertId);
    }
    
    // 4. Добавление товаров
    const goods = [
      { name: "Paracetamol 500mg", image: "/goods/pill1.jpg", producer_id: producerIds[0], is_in_stock: true, instruction: "Take 1 tablet 3-4 times a day", description: "Effective antipyretic", is_prescription_needed: false, price: 150 },
      { name: "Ibuprofen 400mg", image: "/goods/pill2.jpg", producer_id: producerIds[1], is_in_stock: true, instruction: "1 tablet every 6-8 hours", description: "Non-steroidal anti-inflammatory drug", is_prescription_needed: false, price: 200 },
      { name: "Amoxicillin 250mg", image: "/goods/pill3.jpg", producer_id: producerIds[2], is_in_stock: true, instruction: "1 capsule 3 times a day", description: "Broad-spectrum antibiotic", is_prescription_needed: true, price: 350 },
      { name: "Aspirin 500mg", image: "/goods/pill4.jpg", producer_id: producerIds[0], is_in_stock: true, instruction: "1 tablet as needed", description: "Painkiller and antipyretic", is_prescription_needed: false, price: 120 },
      { name: "Vitamin C 1000mg", image: "/goods/pill5.jpg", producer_id: producerIds[3], is_in_stock: true, instruction: "1 tablet per day", description: "Vitamin supplement for immunity", is_prescription_needed: false, price: 300 },
      { name: "No-Spa 40mg", image: "/goods/pill1.jpg", producer_id: producerIds[1], is_in_stock: true, instruction: "1-2 tablets 3 times a day", description: "Antispasmodic", is_prescription_needed: false, price: 180 },
      { name: "Loratadine 10mg", image: "/goods/pill2.jpg", producer_id: producerIds[4], is_in_stock: true, instruction: "1 tablet per day", description: "Antiallergic", is_prescription_needed: false, price: 220 },
      { name: "Validol 60mg", image: "/goods/pill3.jpg", producer_id: producerIds[5], is_in_stock: true, instruction: "1 tablet under the tongue", description: "Sedative for heart pain", is_prescription_needed: false, price: 80 },
      { name: "Activated Charcoal", image: "/goods/pill4.jpg", producer_id: producerIds[6], is_in_stock: true, instruction: "4-6 tablets for poisoning", description: "Adsorbent", is_prescription_needed: false, price: 50 },
      { name: "Amoxiclav 625mg", image: "/goods/pill5.jpg", producer_id: producerIds[2], is_in_stock: false, instruction: "1 tablet 3 times a day", description: "Combined antibiotic", is_prescription_needed: true, price: 450 },
      { name: "Citramon P", image: "/goods/pill1.jpg", producer_id: producerIds[0], is_in_stock: true, instruction: "1 tablet 2-3 times a day", description: "Combined painkiller", is_prescription_needed: false, price: 90 },
      { name: "Mukaltin", image: "/goods/pill2.jpg", producer_id: producerIds[3], is_in_stock: true, instruction: "1-2 tablets 3 times a day", description: "Expectorant", is_prescription_needed: false, price: 60 },
      { name: "Corvalol", image: "/goods/pill3.jpg", producer_id: producerIds[5], is_in_stock: true, instruction: "15-30 drops per dose", description: "Sedative and hypnotic", is_prescription_needed: false, price: 110 },
      { name: "Smecta", image: "/goods/pill4.jpg", producer_id: producerIds[6], is_in_stock: true, instruction: "1 sachet 3 times a day", description: "Antidiarrheal", is_prescription_needed: false, price: 170 },
      { name: "Nasivin Spray", image: "/goods/pill5.jpg", producer_id: producerIds[4], is_in_stock: true, instruction: "1 spray 2-3 times a day", description: "Nasal decongestant", is_prescription_needed: false, price: 190 },
      { name: "Omeprazole 20mg", image: "/goods/pill1.jpg", producer_id: producerIds[1], is_in_stock: true, instruction: "1 capsule in the morning", description: "Heartburn relief", is_prescription_needed: false, price: 280 },
      { name: "Vitamin D3 2000IU", image: "/goods/pill2.jpg", producer_id: producerIds[3], is_in_stock: true, instruction: "1 capsule per day", description: "Vitamin for bones and immunity", is_prescription_needed: false, price: 320 },
      { name: "Glycine 100mg", image: "/goods/pill3.jpg", producer_id: producerIds[7], is_in_stock: true, instruction: "1 tablet 2-3 times a day", description: "Improves cerebral circulation", is_prescription_needed: false, price: 70 },
      { name: "Fenistil Gel", image: "/goods/pill4.jpg", producer_id: producerIds[4], is_in_stock: true, instruction: "Apply thin layer 2-4 times a day", description: "Antiallergic gel", is_prescription_needed: false, price: 380 },
      { name: "Levomycetin", image: "/goods/pill5.jpg", producer_id: producerIds[2], is_in_stock: false, instruction: "1 tablet 3-4 times a day", description: "Broad-spectrum antibiotic", is_prescription_needed: true, price: 95 },
      { name: "Bepanthen Cream", image: "/goods/pill1.jpg", producer_id: producerIds[6], is_in_stock: true, instruction: "Apply to damaged skin 1-2 times a day", description: "Healing cream", is_prescription_needed: false, price: 420 },
      { name: "Nurofen Express", image: "/goods/pill2.jpg", producer_id: producerIds[1], is_in_stock: true, instruction: "1 capsule up to 3 times a day", description: "Fast-acting painkiller", is_prescription_needed: false, price: 270 },
      { name: "Enterosgel", image: "/goods/pill3.jpg", producer_id: producerIds[0], is_in_stock: true, instruction: "1 tablespoon 3 times a day", description: "Enterosorbent for poisoning", is_prescription_needed: false, price: 380 },
      { name: "Valerian Drops", image: "/goods/pill4.jpg", producer_id: producerIds[5], is_in_stock: true, instruction: "20-30 drops 3-4 times a day", description: "Herbal sedative", is_prescription_needed: false, price: 65 },
      { name: "Suprastin 25mg", image: "/goods/pill5.jpg", producer_id: producerIds[4], is_in_stock: true, instruction: "1 tablet 2-3 times a day", description: "Antihistamine", is_prescription_needed: true, price: 130 },
      { name: "Mezim Forte", image: "/goods/pill1.jpg", producer_id: producerIds[3], is_in_stock: true, instruction: "1-2 tablets with meals", description: "Enzyme supplement", is_prescription_needed: false, price: 290 },
      { name: "Azithromycin 500mg", image: "/goods/pill2.jpg", producer_id: producerIds[2], is_in_stock: true, instruction: "1 tablet per day for 3 days", description: "Azalide antibiotic", is_prescription_needed: true, price: 510 },
      { name: "Nise Gel", image: "/goods/pill3.jpg", producer_id: producerIds[1], is_in_stock: true, instruction: "Apply 3-4 times a day", description: "Pain relief gel", is_prescription_needed: false, price: 340 },
      { name: "Kagocel", image: "/goods/pill4.jpg", producer_id: producerIds[0], is_in_stock: true, instruction: "According to scheme 2+2+2 tablets", description: "Antiviral", is_prescription_needed: false, price: 480 },
      { name: "Linex", image: "/goods/pill5.jpg", producer_id: producerIds[6], is_in_stock: true, instruction: "2 capsules 3 times a day", description: "Probiotic for intestines", is_prescription_needed: false, price: 520 },
      { name: "Afobazol 10mg", image: "/goods/pill1.jpg", producer_id: producerIds[7], is_in_stock: true, instruction: "1 tablet 3 times a day", description: "Anti-anxiety", is_prescription_needed: false, price: 410 },
      { name: "Furacilin", image: "/goods/pill2.jpg", producer_id: producerIds[3], is_in_stock: true, instruction: "Dissolve 1 tablet in glass of water", description: "Antiseptic for gargling", is_prescription_needed: false, price: 75 },
      { name: "Teraflu", image: "/goods/pill3.jpg", producer_id: producerIds[1], is_in_stock: true, instruction: "1 sachet per cup of hot water", description: "Cold remedy", is_prescription_needed: false, price: 360 },
      { name: "Hexoral Spray", image: "/goods/pill4.jpg", producer_id: producerIds[4], is_in_stock: true, instruction: "2 sprays 2 times a day", description: "Throat antiseptic", is_prescription_needed: false, price: 290 },
      { name: "Lazolvan Syrup", image: "/goods/pill5.jpg", producer_id: producerIds[5], is_in_stock: true, instruction: "5 ml 3 times a day", description: "Mucolytic", is_prescription_needed: false, price: 330 },
      { name: "Vitrum Vitamins", image: "/goods/pill1.jpg", producer_id: producerIds[3], is_in_stock: true, instruction: "1 tablet per day", description: "Vitamin and mineral complex", is_prescription_needed: false, price: 890 },
      { name: "Ceftriaxone", image: "/goods/pill2.jpg", producer_id: producerIds[2], is_in_stock: false, instruction: "For injections only", description: "Cephalosporin antibiotic", is_prescription_needed: true, price: 120 },
      { name: "Panthenol Spray", image: "/goods/pill3.jpg", producer_id: producerIds[6], is_in_stock: true, instruction: "Spray on affected area", description: "Burn treatment", is_prescription_needed: false, price: 280 },
      { name: "Analgin 500mg", image: "/goods/pill4.jpg", producer_id: producerIds[0], is_in_stock: true, instruction: "1 tablet 2-3 times a day", description: "Painkiller and antipyretic", is_prescription_needed: false, price: 85 },
      { name: "Essentiale Forte", image: "/goods/pill5.jpg", producer_id: producerIds[7], is_in_stock: true, instruction: "2 capsules 3 times a day", description: "Hepatoprotector", is_prescription_needed: false, price: 670 },
    ];
    
    const goodsIds: number[] = [];
    for (const item of goods) {
      const result = await queryRunner.query(
        "INSERT INTO goods (name, image, producerId, isInStock, instruction, description, isPrescriptionNeeded, price) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
        [item.name, item.image, item.producer_id, item.is_in_stock, item.instruction, item.description, item.is_prescription_needed, item.price]
      );
      goodsIds.push(result.insertId);
    }
    
    // 5. Добавление связей товаров с тегами (расширенная версия)
    const goodsTags = [
      // Paracetamol
      { goods_id: goodsIds[0], tag_id: tagIds[0] }, // Painkiller
      { goods_id: goodsIds[0], tag_id: tagIds[1] }, // Antipyretic
      // Ibuprofen
      { goods_id: goodsIds[1], tag_id: tagIds[0] }, // Painkiller
      { goods_id: goodsIds[1], tag_id: tagIds[1] }, // Antipyretic
      { goods_id: goodsIds[1], tag_id: tagIds[10] }, // Anti-inflammatory
      // Amoxicillin
      { goods_id: goodsIds[2], tag_id: tagIds[2] }, // Antibiotic
      // Aspirin
      { goods_id: goodsIds[3], tag_id: tagIds[0] }, // Painkiller
      { goods_id: goodsIds[3], tag_id: tagIds[1] }, // Antipyretic
      // Vitamin C
      { goods_id: goodsIds[4], tag_id: tagIds[3] }, // Vitamins
      // No-Spa
      { goods_id: goodsIds[5], tag_id: tagIds[6] }, // Stomach
      // Loratadine
      { goods_id: goodsIds[6], tag_id: tagIds[8] }, // Allergy
      // Validol
      { goods_id: goodsIds[7], tag_id: tagIds[7] }, // Heart
      { goods_id: goodsIds[7], tag_id: tagIds[9] }, // Sedative
      // Activated Charcoal
      { goods_id: goodsIds[8], tag_id: tagIds[6] }, // Stomach
      // Amoxiclav
      { goods_id: goodsIds[9], tag_id: tagIds[2] }, // Antibiotic
      // Citramon
      { goods_id: goodsIds[10], tag_id: tagIds[0] }, // Painkiller
      // Mukaltin
      { goods_id: goodsIds[11], tag_id: tagIds[5] }, // Cough
      // Corvalol
      { goods_id: goodsIds[12], tag_id: tagIds[9] }, // Sedative
      // Smecta
      { goods_id: goodsIds[13], tag_id: tagIds[6] }, // Stomach
      // Nasivin
      { goods_id: goodsIds[14], tag_id: tagIds[4] }, // Antiseptic
      // Omeprazole
      { goods_id: goodsIds[15], tag_id: tagIds[6] }, // Stomach
      // Vitamin D3
      { goods_id: goodsIds[16], tag_id: tagIds[3] }, // Vitamins
      // Glycine
      { goods_id: goodsIds[17], tag_id: tagIds[9] }, // Sedative
      // Fenistil
      { goods_id: goodsIds[18], tag_id: tagIds[8] }, // Allergy
      // Levomycetin
      { goods_id: goodsIds[19], tag_id: tagIds[2] }, // Antibiotic
      // Bepanthen
      { goods_id: goodsIds[20], tag_id: tagIds[4] }, // Antiseptic
      // Nurofen
      { goods_id: goodsIds[21], tag_id: tagIds[0] }, // Painkiller
      { goods_id: goodsIds[21], tag_id: tagIds[10] }, // Anti-inflammatory
      // Enterosgel
      { goods_id: goodsIds[22], tag_id: tagIds[6] }, // Stomach
      // Valerian
      { goods_id: goodsIds[23], tag_id: tagIds[9] }, // Sedative
      // Suprastin
      { goods_id: goodsIds[24], tag_id: tagIds[8] }, // Allergy
      // Mezim
      { goods_id: goodsIds[25], tag_id: tagIds[13] }, // Enzyme
      // Azithromycin
      { goods_id: goodsIds[26], tag_id: tagIds[2] }, // Antibiotic
      // Nise
      { goods_id: goodsIds[27], tag_id: tagIds[0] }, // Painkiller
      { goods_id: goodsIds[27], tag_id: tagIds[10] }, // Anti-inflammatory
      // Kagocel
      { goods_id: goodsIds[28], tag_id: tagIds[12] }, // Antiviral
      // Linex
      { goods_id: goodsIds[29], tag_id: tagIds[11] }, // Probiotic
      // Afobazol
      { goods_id: goodsIds[30], tag_id: tagIds[9] }, // Sedative
      // Furacilin
      { goods_id: goodsIds[31], tag_id: tagIds[4] }, // Antiseptic
      // Teraflu
      { goods_id: goodsIds[32], tag_id: tagIds[0] }, // Painkiller
      { goods_id: goodsIds[32], tag_id: tagIds[1] }, // Antipyretic
      // Hexoral
      { goods_id: goodsIds[33], tag_id: tagIds[4] }, // Antiseptic
      // Lazolvan
      { goods_id: goodsIds[34], tag_id: tagIds[5] }, // Cough
      // Vitrum
      { goods_id: goodsIds[35], tag_id: tagIds[3] }, // Vitamins
      // Ceftriaxone
      { goods_id: goodsIds[36], tag_id: tagIds[2] }, // Antibiotic
      // Panthenol
      { goods_id: goodsIds[37], tag_id: tagIds[4] }, // Antiseptic
      // Analgin
      { goods_id: goodsIds[38], tag_id: tagIds[0] }, // Painkiller
      { goods_id: goodsIds[38], tag_id: tagIds[1] }, // Antipyretic
      // Essentiale
      { goods_id: goodsIds[39], tag_id: tagIds[14] }, // Hepatoprotector
    ];
    
    for (const goodsTag of goodsTags) {
      await queryRunner.query(
        "INSERT INTO goods_tags (goodsId, tagId) VALUES (?, ?)",
        [goodsTag.goods_id, goodsTag.tag_id]
      );
    }
    
    // 6. Добавление заказов
    const orders = [
      { client_fio: "Ivanov Ivan Ivanovich", client_email: "ivanov@mail.ru", client_phone: "+79161234567" },
      { client_fio: "Petrova Maria Sergeevna", client_email: "petrova@gmail.com", client_phone: "+79037654321" },
      { client_fio: "Sidorov Alexey Petrovich", client_email: "", client_phone: "+79219876543" },
      { client_fio: "Kuznetsova Elena Vladimirovna", client_email: "kuznetsova@yandex.ru", client_phone: "+79185556677" },
      { client_fio: "Popov Dmitry Nikolaevich", client_email: "popov@gmail.com", client_phone: "+79203334455" },
    ];
    
    const orderIds: number[] = [];
    for (const order of orders) {
      const result = await queryRunner.query(
        "INSERT INTO orders (clientFIO, clientEmail, clientPhone) VALUES (?, ?, ?)",
        [order.client_fio, order.client_email, order.client_phone]
      );
      orderIds.push(result.insertId);
    }
    
    // 7. Добавление связей заказов с товарами
    const goodsOrders = [
      { order_id: orderIds[0], goods_id: goodsIds[0], quantity: 2 },
      { order_id: orderIds[0], goods_id: goodsIds[4], quantity: 1 },
      { order_id: orderIds[0], goods_id: goodsIds[15], quantity: 1 },
      { order_id: orderIds[1], goods_id: goodsIds[1], quantity: 1 },
      { order_id: orderIds[1], goods_id: goodsIds[6], quantity: 2 },
      { order_id: orderIds[1], goods_id: goodsIds[24], quantity: 1 },
      { order_id: orderIds[2], goods_id: goodsIds[2], quantity: 1 },
      { order_id: orderIds[2], goods_id: goodsIds[7], quantity: 1 },
      { order_id: orderIds[2], goods_id: goodsIds[11], quantity: 3 },
      { order_id: orderIds[3], goods_id: goodsIds[10], quantity: 2 },
      { order_id: orderIds[3], goods_id: goodsIds[21], quantity: 1 },
      { order_id: orderIds[4], goods_id: goodsIds[29], quantity: 1 },
      { order_id: orderIds[4], goods_id: goodsIds[35], quantity: 2 },
      { order_id: orderIds[4], goods_id: goodsIds[38], quantity: 1 },
    ];
    
    for (const goodsOrder of goodsOrders) {
      await queryRunner.query(
        "INSERT INTO goods_orders (orderId, goodsId, quantity) VALUES (?, ?, ?)",
        [goodsOrder.order_id, goodsOrder.goods_id, goodsOrder.quantity]
      );
    }
    
    // 8. Добавление объявлений
    const announcements = [
      { date_time: new Date('2026-04-10 10:00:00'), from: "Admin", announce: "New batch of antibiotics arrived!" },
      { date_time: new Date('2026-04-11 15:30:00'), from: "Pharmacy", announce: "Discount on vitamins up to 20%" },
      { date_time: new Date('2026-04-12 09:15:00'), from: "Admin", announce: "Free delivery on orders over 1000 rubles" },
      { date_time: new Date('2026-04-13 12:00:00'), from: "Pharmacy", announce: "Pre-order for new medications" },
    ];
    
    for (const announcement of announcements) {
      await queryRunner.query(
        "INSERT INTO announcements (dateTime, `from`, announce) VALUES (?, ?, ?)",
        [announcement.date_time, announcement.from, announcement.announce]
      );
    }
    
    await queryRunner.commitTransaction();
    console.log("Test data inserted successfully!");
    
  } catch (error) {
    await queryRunner.rollbackTransaction();
    console.error("Failed to insert test data:", error);
    throw error;
  } finally {
    await queryRunner.release();
  }
}