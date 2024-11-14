package internal

import "time"

const (
	MotiveElderlyPerson  = "1. PERSONA MAYOR DE 60 AÑOS"
	MotiveChronicDisease = "2. PERSONA CON ENFERMEDAD CRÓNICA"
	MotiveDisabledPerson = "3. PERSONA CON DISCAPACIDAD"
	MotivePregnantPerson = "4. GESTANTE"
	MotivePQRS           = "5. USUARIO QUE INTERPUSO PQRS"
	MotiveOther          = "6. OTRO"
)

const (
	DocumentTypeCC  = "CC"
	DocumentTypeTI  = "TI"
	DocumentTypeRC  = "RC"
	DocumentTypeCE  = "CE"
	DocumentTypePEP = "PEP"
	DocumentTypeDNI = "DNI"
	DocumentTypeSCR = "SCR"
	DocumentTypePA  = "PA"
)

const (
	DistrictUsaquen       = "USAQUÉN"
	DistrictChapinero     = "CHAPINERO"
	DistrictSantaFe       = "SANTA FE"
	DistrictSanCristobal  = "SAN CRISTÓBAL"
	DistrictUsme          = "USME"
	DistrictKenedy        = "KENNEDY"
	DistrictFontibon      = "FONTIBÓN"
	DistrictEngativa      = "ENGATIVÁ"
	DistrictSuba          = "SUBA"
	DistrictBarriosUnidos = "BARRIOS UNIDOS"
	DistrictTeusaquillo   = "TEUSAQUILLO"
	DistrictLosMartires   = "LOS MARTIRES"
	DistrictAntonioNarino = "ANTONIO NARINO"
	DistrictPuenteAranda  = "PUENTE ARANDA"
	DistrictLaCandelaria  = "LA CANDELARIA"
	DistrictRafaelUribe   = "RAFAEL URIBE URIBE"
	DistrictCiudadBolivar = "CIUDAD BOLÍVAR"
	DistrictSumapaz       = "SUMAPAZ"
	DistrictBosa          = "BOSA"
)

const (
	DivisionNorth     = "NORTE"
	DivisionSouth     = "SUR"
	DivisionSouthWest = "SUR OCCIDENTE"
	DivisionSouthEast = "CENTRO ORIENTE"
)

const (
	HasPriorityTrue  = "SI"
	HasPriorityFalse = "NO"
)

type Entry struct {
	FirstName      string
	SecondName     string
	FirstLastName  string
	SecondLastName string
	BirthDate      time.Time
	Medicine       string
	MedicineGiven  time.Time
	Motive         string
	DocumentType   string
	DocumentNumber string
	Address        string
	District       string
	Division       string
	MedicineOrder  time.Time
	HasPriority    bool
}

type LinedError struct {
	Err  error
	Page string
	Line int
}
