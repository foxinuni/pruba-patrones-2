package internal

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/xuri/excelize/v2"
)

type Line struct {
	Line int
	Page string
	row  []string
}

type Parser interface {
	ParseFile(path string) (<-chan Entry, <-chan LinedError, error)
}

type ExcelParser struct {
	workers int
	wg      sync.WaitGroup
}

func NewExcelParser(workers int) *ExcelParser {
	return &ExcelParser{
		workers: workers,
	}
}

func (p *ExcelParser) ParseFile(path string) (<-chan Entry, <-chan LinedError, error) {
	// the file is openned
	file, err := excelize.OpenFile(path)
	if err != nil {
		return nil, nil, err
	}

	incomming := make(chan Line, 100)
	errors := make(chan LinedError)
	outgoing := make(chan Entry)

	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go p.EntryWorker(incomming, errors, outgoing)
	}

	go func() {
		defer file.Close()

		for sheet := 0; sheet < file.SheetCount; sheet++ {
			// get all rows from file
			rows, err := file.GetRows(file.GetSheetName(sheet))
			if err != nil {
				log.Printf("error getting rows from sheet %s: %s", file.GetSheetName(sheet), err)
				continue
			}

			// parse rows
			for number, row := range rows {
				// add row to incomming channel
				if number != 0 {
					incomming <- Line{
						row:  row,
						Line: number,
						Page: file.GetSheetName(sheet),
					}
				}
			}
		}

		// close the incomming channel
		close(incomming)

		// wait for workers threads to finish
		p.wg.Wait()

		// close the outgoing & error channels
		close(outgoing)
		close(errors)
	}()

	return outgoing, errors, nil
}

func (p *ExcelParser) EntryWorker(incomming <-chan Line, errors chan<- LinedError, outgoing chan<- Entry) {
	defer p.wg.Done()

	for line := range incomming {
		func() {
			defer func() {
				if err := recover(); err != nil {
					errors <- LinedError{
						Err:  err.(error),
						Page: line.Page,
						Line: line.Line,
					}
				}
			}()

			entry, err := ParseEntry(line.row)
			if err != nil {
				errors <- LinedError{
					Err:  err,
					Page: line.Page,
					Line: line.Line,
				}
			} else {
				outgoing <- *entry
			}
		}()
	}
}

func ParseEntry(row []string) (*Entry, error) {
	// check if the amount of rows is correct
	if len(row) != 15 {
		return nil, fmt.Errorf("invalid amount of columns")
	}

	// birthday
	birthdate, err := time.Parse("2/1/2006", strings.TrimSpace(row[4]))
	if err != nil {
		return nil, fmt.Errorf("error parsing date %s: %v", row[4], err)
	}

	// medicine order
	medicineOrder, err := time.Parse("2/1/2006", strings.TrimSpace(row[13]))
	if err != nil {
		return nil, fmt.Errorf("error parsing date %s: %v", row[13], err)
	}

	// medicine given
	medicineGiven, err := time.Parse("2/1/2006", strings.TrimSpace(row[6]))
	if err != nil {
		return nil, fmt.Errorf("error parsing date %s: %v", row[6], err)
	}

	// create entry
	entry := &Entry{
		FirstName:      strings.TrimSpace(row[0]),
		SecondName:     strings.TrimSpace(row[1]),
		FirstLastName:  strings.TrimSpace(row[2]),
		SecondLastName: strings.TrimSpace(row[3]),
		BirthDate:      birthdate,
		Medicine:       strings.TrimSpace(row[5]),
		MedicineGiven:  medicineGiven,
		Motive:         strings.TrimSpace(row[7]),
		DocumentType:   strings.TrimSpace(row[8]),
		DocumentNumber: strings.TrimSpace(row[9]),
		Address:        strings.TrimSpace(row[10]),
		District:       strings.TrimSpace(row[11]),
		Division:       strings.TrimSpace(row[12]),
		MedicineOrder:  medicineOrder,
		HasPriority:    strings.TrimSpace(row[14]) == HasPriorityTrue,
	}

	// validate entry
	if err := ValidateEntry(entry); err != nil {
		return nil, err
	}

	return entry, nil
}

func ValidateEntry(entry *Entry) error {
	// FirstName & FirstLastName must not be empty
	if entry.FirstName == "" || entry.FirstLastName == "" {
		return fmt.Errorf("first name and first last name must not be empty")
	}

	// motive must be one of motives
	motives := []string{
		MotiveElderlyPerson,
		MotiveChronicDisease,
		MotiveDisabledPerson,
		MotivePregnantPerson,
		MotivePQRS,
		MotiveOther,
	}
	if !InSlice(entry.Motive, motives) {
		return fmt.Errorf("motive %q must be one of %v", entry.Motive, motives)
	}

	// document type must be one of document types
	documentTypes := []string{
		DocumentTypeCC,
		DocumentTypeTI,
		DocumentTypeRC,
		DocumentTypeCE,
		DocumentTypePEP,
		DocumentTypeDNI,
		DocumentTypeSCR,
		DocumentTypePA,
	}

	if !InSlice(entry.DocumentType, documentTypes) {
		return fmt.Errorf("document type %q must be one of %v", entry.DocumentType, documentTypes)
	}

	// district must be one of districts
	districts := []string{
		DistrictUsaquen,
		DistrictChapinero,
		DistrictSantaFe,
		DistrictSanCristobal,
		DistrictUsme,
		DistrictKenedy,
		DistrictFontibon,
		DistrictEngativa,
		DistrictSuba,
		DistrictBarriosUnidos,
		DistrictTeusaquillo,
		DistrictLosMartires,
		DistrictAntonioNarino,
		DistrictPuenteAranda,
		DistrictLaCandelaria,
		DistrictRafaelUribe,
		DistrictCiudadBolivar,
		DistrictSumapaz,
		DistrictBosa,
	}
	if !InSlice(entry.District, districts) {
		return fmt.Errorf("district %q must be one of %v", entry.District, districts)
	}

	// Check division
	division := []string{
		DivisionNorth,
		DivisionSouth,
		DivisionSouthWest,
		DivisionSouthEast,
	}
	if !InSlice(entry.Division, division) {
		return fmt.Errorf("division %q must be one of %v", entry.Division, division)
	}

	// validate document alphanumeric only (no spaces or special chars)
	re := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	if !re.MatchString(entry.DocumentNumber) {
		return fmt.Errorf("document number %q must be alphanumeric", entry.DocumentNumber)
	}

	// validate that priority corresponds with it's motive
	prioties := []string{
		MotiveElderlyPerson,
		MotiveChronicDisease,
		MotiveDisabledPerson,
		MotivePregnantPerson,
	}
	if InSlice(entry.Motive, prioties) && !entry.HasPriority {
		return fmt.Errorf("entry with motive %q must have priority", entry.Motive)
	} else if !InSlice(entry.Motive, prioties) && entry.HasPriority {
		return fmt.Errorf("entry with motive %q must not have priority", entry.Motive)
	}

	return nil

}

func InSlice(needle string, haystack []string) bool {
	for _, value := range haystack {
		if value == needle {
			return true
		}
	}

	return false
}
