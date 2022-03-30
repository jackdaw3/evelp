package dto

type CorporationDTO struct {
	CorporationId   int
	CorporationName string
}

type CorporationDTOs []*CorporationDTO

type FactionDTO struct {
	FactionId    int
	FactionName  string
	Corporations CorporationDTOs
}

type FactionDTOs []*FactionDTO

func (c CorporationDTOs) Len() int { return len(c) }

func (c CorporationDTOs) Less(i, j int) bool { return c[i].CorporationId < c[j].CorporationId }

func (c CorporationDTOs) Swap(i, j int) { c[i], c[j] = c[j], c[i] }

func (f FactionDTOs) Len() int { return len(f) }

func (f FactionDTOs) Less(i, j int) bool { return f[i].FactionId < f[j].FactionId }

func (f FactionDTOs) Swap(i, j int) { f[i], f[j] = f[j], f[i] }
