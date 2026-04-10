package controller

import (
	"log"
	"math/rand"
	shared_models "pharmacy-api/shared/models"
	"strconv"
	"time"

	"gorm.io/gorm"
)

func GetGoods(query string, pageStr string, limitStr string, db *gorm.DB) (result map[string]any, err error) {
	limit, strErr := strconv.Atoi(limitStr)

	if strErr != nil || limit < 1 {
		limit = 10
	} else if limit > 40 {
		limit = 40
	}
	page, strErr := strconv.Atoi(pageStr)
	if strErr != nil || page < 1 {
		page = 1
	}
	offset := (page - 1) * limit
	var goodsByTags, goodsByProducers, goodsByName []shared_models.Goods
	var totalCount int64
	var goods []shared_models.Goods
	var seenIDs = make(map[uint]bool)
	if query != "" {
		db.Preload("Producer").
			Joins("JOIN goods_tags ON goods.id = goods_tags.goods_id").
			Joins("JOIN tags ON goods_tags.tag_id = tags.id").
			Where("tags.tag_name LIKE ?", "%"+query+"%").
			Order("name ASC").Find(&goodsByTags)

		db.Preload("Producer").
			Joins("JOIN producers ON goods.producer_id = producers.id").
			Where("producers.producer_name LIKE ?", "%"+query+"%").
			Order("name ASC").Find(&goodsByProducers)

		db.Preload("Producer").
			Where("name LIKE ?", "%"+query+"%").
			Order("name ASC").Find(&goodsByName)
		for _, item := range goodsByTags {
			if !seenIDs[item.ID] {
				goods = append(goods, item)
				seenIDs[item.ID] = true
			}
		}
		for _, item := range goodsByProducers {
			if !seenIDs[item.ID] {
				goods = append(goods, item)
				seenIDs[item.ID] = true
			}
		}
		for _, item := range goodsByName {
			if !seenIDs[item.ID] {
				goods = append(goods, item)
				seenIDs[item.ID] = true
			}
		}
		totalCount = int64(len(goods))
		goods = goods[offset:min(len(goods), offset+limit)]
	} else {
		db.Model(&goods).Count(&totalCount)
		db.Order("name ASC").Limit(limit).Offset(offset).Find(&goods)
	}
	totalPages := (totalCount + int64(limit) - 1) / int64(limit)

	var response []shared_models.GoodsResponse
	for _, goodsitem := range goods {
		response = append(response, shared_models.GoodsResponse{
			ID:          goodsitem.ID,
			Name:        goodsitem.Name,
			Image:       goodsitem.Image,
			IsInStock:   goodsitem.IsInStock,
			Description: goodsitem.Description,
			Price:       goodsitem.Price,
		})
	}

	result = make(map[string]any)
	result["Items"] = response
	result["TotalPages"] = totalPages
	result["CurrentPage"] = page

	return
}

func GetGoodsByID(id int, db *gorm.DB) (result map[string]any, err error) {

	var good shared_models.Goods

	if err = db.Preload("Producer").Preload("Tags").Find(&good, id).Error; err != nil {
		return
	}

	var tagNames []string
	for _, tag := range good.Tags {
		tagNames = append(tagNames, tag.TagName)
	}

	response := shared_models.GoodsResponse{
		ID:                   good.ID,
		Name:                 good.Name,
		Image:                good.Image,
		Producer:             good.Producer.ProducerName,
		IsInStock:            good.IsInStock,
		Tags:                 tagNames,
		Instruction:          good.Instruction,
		Description:          good.Description,
		IsPrescriptionNeeded: good.IsPrescriptionNeeded,
		Price:                good.Price,
	}

	result = make(map[string]any)
	result["response"] = response

	return
}

func GetPromoItems(db *gorm.DB) (result map[string]any, err error) {
	var goods []shared_models.Goods
	if err = db.Where("is_in_stock LIKE ?", "1").Find(&goods).Error; err != nil {
		return
	}
	totalGoods := len(goods)
	if totalGoods == 0 {
		return
	}
	rand.NewSource(time.Now().UnixNano())
	indices := rand.Perm(totalGoods)[:min(5, totalGoods)]
	var promoItems []shared_models.PromoItem
	for _, index := range indices {
		promoItems = append(promoItems, shared_models.PromoItem{
			ID:          goods[index].ID,
			Name:        goods[index].Name,
			Description: goods[index].Description,
			Price:       goods[index].Price,
			Image:       goods[index].Image,
		})
	}

	result = make(map[string]any)
	result["responce"] = promoItems

	return
}

// если накроется - добавить * к shared_models.Claims; всё равно дорабатывать
func UpdateGoods(id int, updateData shared_models.GoodsUpdateRequest, claims shared_models.Claims, db *gorm.DB) (result map[string]any, err error) {

	// user, _ := c.Get("user")
	// claims := user.(*shared_models.Claims)

	var existingGood shared_models.Goods
	if err = db.First(&existingGood, id).Error; err != nil {
		log.Println("Good PATH error [" + claims.Username + "]" + err.Error())
		return
	}

	existingGood.Name = updateData.Name
	existingGood.Instruction = updateData.Instruction
	existingGood.Description = updateData.Description
	existingGood.IsPrescriptionNeeded = updateData.IsPrescriptionNeeded
	existingGood.IsInStock = updateData.IsInStock
	existingGood.Price = updateData.Price

	if err = db.Save(&existingGood).Error; err != nil {
		log.Println("Good PATH error [" + claims.Username + "]" + err.Error())
		return
	}

	log.Println("Good PATH [" + claims.Username + "]")

	result = make(map[string]any)
	result["message"] = "updated"
	result["data"] = existingGood

	return
}
