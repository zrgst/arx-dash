package arx

import "encoding/xml"

// rawARXData matcher XML-roten:
//
// <arxdata timestamp="...">
//
//	<persons>...</persons>
//	<cards>...</cards>
//
// </arxdata>
type rawARXData struct {
	XMLName   xml.Name    `xml:"arxdata"`
	Timestamp string      `xml:"timestamp,attr"`
	Persons   []rawPerson `xml:"persons>person"`
	Cards     []rawCard   `xml:"cards>card"`
}

type rawPerson struct {
	ID                   string                  `xml:"id"`
	FirstName            string                  `xml:"first_name"`
	LastName             string                  `xml:"last_name"`
	Description          string                  `xml:"description"`
	PinCode              string                  `xml:"pin_code"`
	PictureBlob          string                  `xml:"picture_blob"`
	PictureURL           string                  `xml:"picture_url"`
	Username             string                  `xml:"username"`
	Password             string                  `xml:"password"`
	ForcePinChange       string                  `xml:"force_pin_change"`
	Deleted              *struct{}               `xml:"deleted"`
	OfflineUnlockFeature string                  `xml:"offline_unlock_feature"`
	ResourceAccesses     []string                `xml:"resource_accesses>access_name"`
	ExtraFields          []rawExtraField         `xml:"extra_fields>extra_field"`
	AccessCategories     []rawAccessCategory     `xml:"access_categories>access_category"`
	AccessUpdates        rawAccessCategoryUpdate `xml:"access_categories_update"`
}

type rawExtraField struct {
	Name  string `xml:"name"`
	Value string `xml:"value"`
}

type rawAccessCategory struct {
	ID        string `xml:"id"`
	Name      string `xml:"name"`
	StartDate string `xml:"start_date"`
	EndDate   string `xml:"end_date"`
}

type rawAccessCategoryUpdate struct {
	Add     []rawAccessCategory `xml:"access_category_add"`
	Remove  []rawAccessCategory `xml:"access_category_remove"`
	Replace []rawAccessCategory `xml:"access_category_replace"`
}

type rawCard struct {
	Number      string    `xml:"number"`
	FormatName  string    `xml:"format_name"`
	Description string    `xml:"description"`
	PersonID    string    `xml:"person_id"`
	Inhibited   *struct{} `xml:"inhibited"`
	Deleted     *struct{} `xml:"deleted"`
}

// ParsePersonsExport parser rå XML fra /arx/export og returnerer vår egen modell.
func ParsePersonsExport(data []byte) (PersonsExport, error) {
	var raw rawARXData

	if err := xml.Unmarshal(data, &raw); err != nil {
		return PersonsExport{}, err
	}

	return normalizePersonsExport(raw), nil
}

func normalizePersonsExport(raw rawARXData) PersonsExport {
	cards := normalizeCards(raw.Cards)

	// Brukes for å koble kort til eiernavn.
	personNameByID := make(map[string]string)

	for _, person := range raw.Persons {
		personNameByID[person.ID] = fullName(person.FirstName, person.LastName)
	}

	for i := range cards {
		cards[i].OwnerName = personNameByID[cards[i].PersonID]
	}

	persons := make([]Person, 0, len(raw.Persons))

	for _, rawPerson := range raw.Persons {
		personCards := cardsForPerson(cards, rawPerson.ID)

		persons = append(persons, Person{
			ID:                   rawPerson.ID,
			FirstName:            rawPerson.FirstName,
			LastName:             rawPerson.LastName,
			Description:          rawPerson.Description,
			PinCode:              rawPerson.PinCode,
			PictureBlob:          rawPerson.PictureBlob,
			PictureURL:           rawPerson.PictureURL,
			Username:             rawPerson.Username,
			Password:             rawPerson.Password,
			ForcePinChange:       rawPerson.ForcePinChange,
			Deleted:              rawPerson.Deleted != nil,
			OfflineUnlockFeature: rawPerson.OfflineUnlockFeature,
			ResourceAccesses:     rawPerson.ResourceAccesses,
			ExtraFields:          normalizeExtraFields(rawPerson.ExtraFields),
			AccessCategories:     normalizeAccessCategories(rawPerson.AccessCategories),
			AccessUpdates:        normalizeAccessUpdates(rawPerson.AccessUpdates),
			Cards:                personCards,
		})
	}

	return PersonsExport{
		Timestamp: raw.Timestamp,
		Persons:   persons,
		Cards:     cards,
	}
}

func normalizeCards(rawCards []rawCard) []Card {
	cards := make([]Card, 0, len(rawCards))

	for _, rawCard := range rawCards {
		cards = append(cards, Card{
			Number:      rawCard.Number,
			FormatName:  rawCard.FormatName,
			Description: rawCard.Description,
			PersonID:    rawCard.PersonID,
			Inhibited:   rawCard.Inhibited != nil,
			Deleted:     rawCard.Deleted != nil,
		})
	}

	return cards
}

func normalizeExtraFields(rawFields []rawExtraField) []ExtraField {
	fields := make([]ExtraField, 0, len(rawFields))

	for _, rawField := range rawFields {
		fields = append(fields, ExtraField{
			Name:  rawField.Name,
			Value: rawField.Value,
		})
	}

	return fields
}

func normalizeAccessCategories(rawCategories []rawAccessCategory) []AccessCategory {
	categories := make([]AccessCategory, 0, len(rawCategories))

	for _, rawCategory := range rawCategories {
		categories = append(categories, AccessCategory{
			ID:        rawCategory.ID,
			Name:      rawCategory.Name,
			StartDate: rawCategory.StartDate,
			EndDate:   rawCategory.EndDate,
		})
	}

	return categories
}

func normalizeAccessUpdates(rawUpdates rawAccessCategoryUpdate) AccessUpdates {
	return AccessUpdates{
		Add:     normalizeAccessCategories(rawUpdates.Add),
		Remove:  normalizeAccessCategories(rawUpdates.Remove),
		Replace: normalizeAccessCategories(rawUpdates.Replace),
	}
}

func cardsForPerson(cards []Card, personID string) []Card {
	result := []Card{}

	for _, card := range cards {
		if card.PersonID == personID {
			result = append(result, card)
		}
	}

	return result
}

func fullName(firstName string, lastName string) string {
	switch {
	case firstName == "" && lastName == "":
		return ""
	case firstName == "":
		return lastName
	case lastName == "":
		return firstName
	default:
		return firstName + " " + lastName
	}
}
