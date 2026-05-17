package arx

// Person er vår normaliserte modell for en ARX-person.
//
// Denne brukes av API/UI, ikke direkte XML.
type Person struct {
	ID                   string           `json:"id"`
	FirstName            string           `json:"firstName"`
	LastName             string           `json:"lastName"`
	Description          string           `json:"description"`
	PinCode              string           `json:"pinCode"`
	PictureBlob          string           `json:"pictureBlob"`
	PictureURL           string           `json:"pictureUrl"`
	Username             string           `json:"username"`
	Password             string           `json:"password"`
	ForcePinChange       string           `json:"forcePinChange"`
	Deleted              bool             `json:"deleted"`
	OfflineUnlockFeature string           `json:"offlineUnlockFeature"`
	ResourceAccesses     []string         `json:"resourceAccesses"`
	ExtraFields          []ExtraField     `json:"extraFields"`
	AccessCategories     []AccessCategory `json:"accessCategories"`
	AccessUpdates        AccessUpdates    `json:"accessUpdates"`
	Cards                []Card           `json:"cards"`
}

// Card er vår normaliserte modell for adgangskort.
type Card struct {
	Number      string `json:"number"`
	FormatName  string `json:"formatName"`
	Description string `json:"description"`
	PersonID    string `json:"personId"`
	OwnerName   string `json:"ownerName"`
	Inhibited   bool   `json:"inhibited"`
	Deleted     bool   `json:"deleted"`
}

// ExtraField er et ekstra personfelt fra ARX.
type ExtraField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// AccessCategory er adgangsgruppe/adgangskategori knyttet til person.
type AccessCategory struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}

// AccessUpdates brukes hvis exporten inneholder add/remove/replace for adgangsgrupper.
type AccessUpdates struct {
	Add     []AccessCategory `json:"add"`
	Remove  []AccessCategory `json:"remove"`
	Replace []AccessCategory `json:"replace"`
}

// PersonsExport er samlet normalisert resultat fra /arx/export.
type PersonsExport struct {
	Timestamp string   `json:"timestamp"`
	Persons   []Person `json:"persons"`
	Cards     []Card   `json:"cards"`
}
