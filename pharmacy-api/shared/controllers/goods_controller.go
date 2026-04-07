package controllers

import (
	"net/http"
	shared_models "pharmacy-api/shared/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetGoods(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := c.Query("q")
		pageStr := c.Query("page")
		limitStr := c.Query("limit")
		limit, err := strconv.Atoi(limitStr)

		if err != nil || limit < 1 {
			limit = 10
		} else if limit > 40 {
			limit = 40
		}
		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
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

		c.JSON(http.StatusOK, gin.H{
			"Items":       response,
			"TotalPages":  totalPages,
			"CurrentPage": page,
		})
	}
}

func GetGoodsByID(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var good shared_models.Goods

		if err := db.Preload("Producer").Preload("Tags").Find(&good, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
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

		c.JSON(http.StatusOK, response)
	}
}
