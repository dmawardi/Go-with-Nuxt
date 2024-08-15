package adminpanel

import (
	"fmt"
	"strings"

	"github.com/dmawardi/Go-Template/internal/helpers/data"
	"github.com/dmawardi/Go-Template/internal/models"
)

// TABLE DEFINITIONS AND BUILDER
//
// Table data is used to store all data pertaining to a table
type TableData struct {
	AdminSchemaUrl string // eg. /users/
	TableHeaders   []TableHeader
	TableRows      []TableRow
	MetaData       models.ExtendedSchemaMetaData
}

// Used for table header information. Also holds information for sorting and pointer + data type
type TableHeader struct {
	Label string
	// Is the column sortable
	Sortable bool
	// label used in db to sort
	ColumnSortLabel string
	// Is data type a pointer
	Pointer bool
	// Used for pointer to string extraction
	DataType string
	// Used for foreign key extraction
	ForeignKeyRepKeyName string
}

// Data to complete a table row
type TableRow struct {
	Data []TableCell
	Edit EditInfo
}

// Data for a single cell
type TableCell struct {
	Label   string
	RowSpan int
	// Primarily used for the policy table
	EditLink string
}

// Edit info for the Edit column in the table
type EditInfo struct {
	EditAllowed bool
	EditUrl   string // eg. admin/users/1
	DeleteUrl string // eg. admin/users/delete/1
}

// Used for bulk delete form on find all pages
type BulkDeleteRequest struct {
	SelectedItems []string `json:"selected_items"`
}

// Function to build table data from slice of adminpanel schema objects, admin schema url (eg. /admin/users) and table headers
func BuildTableData(listOfSchemaObjects []models.AdminPanelSchema, metaData models.SchemaMetaData, adminSchemaBaseUrl string, tableHeaders []TableHeader, allowEdit bool) TableData {
	// Calculate currently showing records and total pages
	currentlyShowing := metaData.CalculateCurrentlyShowingRecords()
	// Init table data
	tableData := TableData{
		AdminSchemaUrl: adminSchemaBaseUrl,
		TableHeaders:   tableHeaders,
		TableRows:      []TableRow{},
		// Build extended metadata
		MetaData: models.NewExtendedSchemaMetaData(metaData, currentlyShowing),
	}

	// Loop through listOfSchemaObjects and build table rows
	for _, object := range listOfSchemaObjects {
		// Init table row
		row := TableRow{
			Data: []TableCell{},
			// Fill in edit info
			Edit: EditInfo{
				EditAllowed: allowEdit,
				EditUrl:   fmt.Sprintf("%s/%s", adminSchemaBaseUrl, object.GetID()),
				DeleteUrl: fmt.Sprintf("%s/delete/%s", adminSchemaBaseUrl, object.GetID()),
			},
		}

		// Iterate through tableheaders to extract table data
		for _, header := range tableHeaders {
			// If found to be a foreign key, extract the foreign key value
			if header.DataType == "foreign" {
				// Attempt to extract foreign key representative value
				fieldData := object.ObtainValue(header.Label)
				// Convert string value to map[string]string
				foreignKeyDataMap, err := data.StringToMap(fieldData)
				if err != nil {
					fmt.Printf("Error converting struct to map: %v\n", err)
				}
				// Extract foreign key id
				foreignKeyID := foreignKeyDataMap["ID"]
				// Extract foreign key representative value
				foreignKeyValue := foreignKeyDataMap[header.ForeignKeyRepKeyName]

				// Append to row data
				row.Data = append(row.Data, TableCell{Label: fmt.Sprintf("(%s) %s", foreignKeyID, foreignKeyValue)})

				// If not a foreign key, extract the value from the schema object
			} else {
				// Grab data from the schema object
				fieldData := object.ObtainValue(header.Label)

				// convert fieldData to string
				stringFieldData := fmt.Sprint(fieldData)

				// Use header string values to get values from schema object and append
				row.Data = append(row.Data, TableCell{Label: stringFieldData})
			}
		}

		// Append row to table data
		tableData.TableRows = append(tableData.TableRows, row)
	}

	return tableData
}



// POLICY TABLE
//
// Used in template file
type PolicyEditSelectors struct {
	RoleSelection   []FormFieldSelector
	ActionSelection []FormFieldSelector
}

// Formatted data for ideal edit page rendering
type PolicyEditDataRow struct {
	Resource string
	Role     string
	Actions  []PolicyActionCell
}

// Used for rendering the policy table row
type PolicyActionCell struct {
	Action  string
	Granted bool
}

// Function build permission policy table data for Permissions
func BuildPolicyTableData(policySlice []models.PolicyRuleCombinedActions, adminSchemaBaseUrl string, tableHeaders []TableHeader) TableData {
	var tableRows []TableRow

	// Loop through policy slice to build table rows
	for _, policy := range policySlice {
		var rowData []TableCell

		// Iterate through tableheaders
		for _, header := range tableHeaders {
			var value string
			var found bool

			// Check the header label and assign the corresponding field from policy
			switch header.Label {
			case "role":
				value = policy.Role
				found = true
			case "resource":
				value = policy.Resource
				found = true
			case "action":
				value = strings.Join(policy.Action, ", ")
				found = true
			}

			// If the key is found, append the value to the row data
			if found {
				// Append with edit link if it's the first column (resource)
				if header.Label == "resource" {
					// Create the edit link from the label value (slugify)
					editLink := strings.ReplaceAll(value, "/", "-")
					// Append to row data with edit link
					rowData = append(rowData, TableCell{Label: fmt.Sprintf("%v", value), EditLink: editLink})

					// else if other column
				} else {
					rowData = append(rowData, TableCell{Label: fmt.Sprintf("%v", value), EditLink: ""})
				}

				// If the key is not found, append an empty string
			} else {
				rowData = append(rowData, TableCell{Label: ""}) // Add an empty string if the key is not found
			}
		}

		// Append to table rows
		tableRows = append(tableRows, TableRow{Data: rowData})
	}
	return TableData{
		AdminSchemaUrl: adminSchemaBaseUrl, // You can set this value as needed
		TableHeaders:   tableHeaders,
		TableRows:      tableRows,
	}
}

// Build auth role table data
func BuildRoleTableData(roleSlice []string, adminSchemaBaseUrl string, tableHeaders []TableHeader) TableData {
	var tableRows []TableRow

	// Loop through policy slice to build table rows
	for _, role := range roleSlice {
		var rowData []TableCell

		// Grab data from the schema object
		value := role

		rowData = append(rowData, TableCell{Label: fmt.Sprintf("%v", value), EditLink: ""})

		// Append to table rows
		tableRows = append(tableRows, TableRow{Data: rowData})
	}
	return TableData{
		AdminSchemaUrl: adminSchemaBaseUrl, // You can set this value as needed
		TableHeaders:   tableHeaders,
		TableRows:      tableRows,
	}
}

// Build auth role inheritance table data
func BuildRoleInheritanceTableData(policySlice []models.GRecord, adminSchemaBaseUrl string, tableHeaders []TableHeader) TableData {
	var tableRows []TableRow

	// Loop through policy slice to build table rows
	for _, policy := range policySlice {
		var rowData []TableCell
		policySlug := fmt.Sprintf("%s,%s", policy.Role, policy.InheritsFrom)

		// Iterate through tableheaders
		for _, header := range tableHeaders {
			var value string
			var found bool

			// Map header.Label to struct fields
			switch header.Label {
			case "role":
				value = policy.Role
				found = true
			case "inherits_from":
				value = policy.InheritsFrom
				found = true
			}

			// If the key is found, append the value to the row data
			if found {
				rowData = append(rowData, TableCell{Label: fmt.Sprintf("%v", value), EditLink: ""})
			}
		}

		// Append to table rows
		tableRows = append(tableRows, TableRow{Data: rowData, Edit: EditInfo{DeleteUrl: fmt.Sprintf("%s/delete-inheritance/%s", adminSchemaBaseUrl, policySlug)}})
	}
	return TableData{
		AdminSchemaUrl: adminSchemaBaseUrl, // You can set this value as needed
		TableHeaders:   tableHeaders,
		TableRows:      tableRows,
	}
}

// Edit a slice of table rows to add row span to first column and remove <td> tags from
// subsequent rows (Used on policy table)
func editTableDataRowSpan(tableRows []TableRow) {
	var lastRecordedStart struct {
		resource string
		index    int
	}
	for i, row := range tableRows {
		// Extract row variables
		rowData := row.Data
		currentResource := rowData[0].Label

		// If current resource is different from last recorded resource, then must edit row span
		if currentResource != lastRecordedStart.resource {
			// Calculate row span
			rowSpan := i - lastRecordedStart.index

			// If the difference between the current index and the last recorded index is greater than 1, then must edit row span
			if i-lastRecordedStart.index > 1 {
				// Add row span to first cell of last recorded start row
				tableRows[lastRecordedStart.index].Data[0].RowSpan = rowSpan

				// Remove <td> tags from subsequent rows: count from last recorded index + 1 till before current index
				for j := lastRecordedStart.index + 1; j < i; j++ {
					// Chop first element off of data
					tableRows[j].Data = tableRows[j].Data[1:]
				}
			}
			// Reassign lastRecordedResource details
			lastRecordedStart.resource = currentResource
			lastRecordedStart.index = i

		}
	}
}
