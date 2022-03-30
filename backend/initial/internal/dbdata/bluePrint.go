package dbdata

import (
	"encoding/csv"
	"evelp/log"
	"evelp/model"
	"os"
	"sort"
	"strconv"
)

const (
	productActivity = 1
	inventActivity  = 8
	reactActivity   = 11
)

type bluePrintData struct {
	productFilePath  string
	materialFilePath string
	bluePrintsMap    map[int]*model.BluePrint
}

func (b *bluePrintData) Refresh() error {
	log.Infof("start to load bluePrints %s and %s", b.productFilePath, b.materialFilePath)
	b.bluePrintsMap = make(map[int]*model.BluePrint)

	bluePrints, err := b.loadBluePrints()
	if err != nil {
		return err
	}
	log.Info("loading bluePrints finished")

	log.Info("start to save bluePrints to DB")
	if err := model.SaveBluePrints(bluePrints); err != nil {
		return err
	}
	log.Info("bluePrints saved to DB")

	return nil
}

func (b *bluePrintData) loadBluePrints() (*model.BluePrints, error) {
	products, err := readCsv(b.productFilePath)
	if err != nil {
		return nil, err
	}
	if err := b.covertProducts(products); err != nil {
		return nil, err
	}

	materials, err := readCsv(b.materialFilePath)
	if err != nil {
		return nil, err
	}
	if err := b.covertMaterials((materials)); err != nil {
		return nil, err
	}

	var bluePrints model.BluePrints
	for _, v := range b.bluePrintsMap {
		bluePrints = append(bluePrints, *v)
	}
	sort.Sort(bluePrints)
	return &bluePrints, nil

}

func (b *bluePrintData) covertProducts(products [][]string) error {
	for _, product := range products {
		bluePrintId, err := strconv.Atoi(product[0])
		if err != nil {
			return err
		}

		activityId, err := strconv.Atoi(product[1])
		if err != nil {
			return err
		}
		if activityId != productActivity && activityId != reactActivity {
			continue
		}

		manufactProduct := new(model.ManufactProduct)
		productId, err := strconv.Atoi(product[2])
		if err != nil {
			return err
		}
		manufactProduct.ItemId = productId
		productQuantity, err := strconv.ParseInt(product[3], 10, 64)
		if err != nil {
			return err
		}
		manufactProduct.Quantity = productQuantity

		if _, ok := b.bluePrintsMap[bluePrintId]; !ok {
			bluePrint := new(model.BluePrint)
			bluePrint.BlueprintId = bluePrintId
			bluePrint.Products = append(bluePrint.Products, *manufactProduct)
			b.bluePrintsMap[bluePrintId] = bluePrint
		} else {
			b.bluePrintsMap[bluePrintId].Products = append(b.bluePrintsMap[bluePrintId].Products, *manufactProduct)
		}
	}

	return nil
}

func (b *bluePrintData) covertMaterials(materails [][]string) error {
	for _, materail := range materails {
		bluePrintId, err := strconv.Atoi(materail[0])
		if err != nil {
			return err
		}

		if _, ok := b.bluePrintsMap[bluePrintId]; !ok {
			log.Warnf("blueprint %d has no product item", bluePrintId)
			continue
		}

		activityId, err := strconv.Atoi(materail[1])
		if err != nil {
			return err
		}
		if activityId != productActivity && activityId != reactActivity {
			continue
		}

		manufactMaterial := new(model.ManufactMaterial)
		materailId, err := strconv.Atoi(materail[2])
		if err != nil {
			return err
		}
		manufactMaterial.ItemId = materailId
		materailQuantity, err := strconv.ParseInt(materail[3], 10, 64)
		if err != nil {
			return err
		}
		manufactMaterial.Quantity = materailQuantity

		b.bluePrintsMap[bluePrintId].Materials = append(b.bluePrintsMap[bluePrintId].Materials, *manufactMaterial)
	}
	return nil
}

func readCsv(filePath string) ([][]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	return records[1:], nil
}
