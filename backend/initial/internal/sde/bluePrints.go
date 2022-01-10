package sde

import (
	"encoding/csv"
	"evelp/model"
	"os"
	"sort"
	"strconv"

	log "github.com/sirupsen/logrus"
)

const (
	productActivity = 1
	inventActivity  = 8
	reactActivity   = 11
)

type BluePrintsInit struct {
	ProductFilePath  string
	MaterialFilePath string
	bluePrintsMap    map[int]*model.BluePrint
}

func (b *BluePrintsInit) Refresh() error {
	log.Infof("Start load bluePrints %s and %s", b.ProductFilePath, b.MaterialFilePath)
	b.bluePrintsMap = make(map[int]*model.BluePrint)

	bluePrints, err := b.load()
	if err != nil {
		return err
	}
	log.Info("Load bluePrints finished.")

	log.Info("Save bluePrints to DB.")
	if err := model.SaveBluePrints(bluePrints); err != nil {
		return err
	}
	log.Info("BluePrints have saved to DB.")

	return nil
}

func (b *BluePrintsInit) load() (*model.BluePrints, error) {
	products, err := read(b.ProductFilePath)
	if err != nil {
		return nil, err
	}
	if err := b.covertProducts(products); err != nil {
		return nil, err
	}

	materials, err := read(b.MaterialFilePath)
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

func (b *BluePrintsInit) covertProducts(products [][]string) error {
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
		productQuantity, err := strconv.Atoi(product[3])
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

func (b *BluePrintsInit) covertMaterials(materails [][]string) error {
	for _, materail := range materails {
		bluePrintId, err := strconv.Atoi(materail[0])
		if err != nil {
			return err
		}

		if _, ok := b.bluePrintsMap[bluePrintId]; !ok {
			log.Warnf("Blueprint %d has no product item.", bluePrintId)
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
		materailQuantity, err := strconv.Atoi(materail[3])
		if err != nil {
			return err
		}
		manufactMaterial.Quantity = materailQuantity

		b.bluePrintsMap[bluePrintId].Materials = append(b.bluePrintsMap[bluePrintId].Materials, *manufactMaterial)
	}
	return nil
}

func read(filePath string) ([][]string, error) {
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
