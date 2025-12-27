package godbf

import (
	"fmt"
	"testing"

	. "github.com/onsi/gomega"
)

const testEncoding = "UTF-8"
const realEncoding = "cp866"

// For reference: https://en.wikipedia.org/wiki/.dbf#File_format_of_Level_5_DOS_dBASE

func TestDbfTable_New(t *testing.T) {
	g := NewGomegaWithT(t)

	tableUnderTest := New(testEncoding)

	g.Expect(tableUnderTest.NumberOfRecords()).To(BeZero())
	g.Expect(len(tableUnderTest.Fields())).To(BeZero())
}

func TestDbfTable_AddBooleanField(t *testing.T) {
	g := NewGomegaWithT(t)

	tableUnderTest := New(testEncoding)
	expectedFieldName := "testBool"
	additionError := tableUnderTest.AddBooleanField(expectedFieldName)
	g.Expect(additionError).To(BeNil())

	g.Expect(tableUnderTest.NumberOfRecords()).To(BeZero())
	g.Expect(len(tableUnderTest.Fields())).To(BeNumerically("==", 1))

	addedField := tableUnderTest.Fields()[0]
	g.Expect(addedField.name).To(Equal(expectedFieldName))
	g.Expect(addedField.fieldType).To(Equal(Logical))
}

func TestDbfTable_AddBooleanField_TooLongGetsTruncated(t *testing.T) {
	g := NewGomegaWithT(t)

	tableUnderTest := New(testEncoding)
	expectedFieldName := "FieldName!"
	suppliedFieldName := expectedFieldName + "shouldBeTruncated"

	tableUnderTest.AddBooleanField(suppliedFieldName)

	addedField := tableUnderTest.Fields()[0]
	g.Expect(addedField.name).To(Equal(expectedFieldName))
}

func TestDbfTable_AddBooleanField_SecondAttemptFails(t *testing.T) {
	g := NewGomegaWithT(t)

	tableUnderTest := New(testEncoding)
	expectedFieldName := "FieldName!"

	additionError := tableUnderTest.AddBooleanField(expectedFieldName)
	g.Expect(additionError).To(BeNil())

	secondAdditionError := tableUnderTest.AddBooleanField(expectedFieldName)
	g.Expect(secondAdditionError).ToNot(BeNil())

	t.Log(secondAdditionError)
}

func TestDbfTable_AddBooleanField_ErrorAfterDataEntryStart(t *testing.T) {
	g := NewGomegaWithT(t)

	tableUnderTest := New(testEncoding)
	expectedFieldName := "goodField"

	additionError := tableUnderTest.AddBooleanField(expectedFieldName)
	g.Expect(additionError).To(BeNil())

	tableUnderTest.AddNewRecord()

	postDataEntryField := "badField"

	secondAdditionError := tableUnderTest.AddBooleanField(postDataEntryField)
	g.Expect(secondAdditionError).ToNot(BeNil())

	t.Log(secondAdditionError)
}

func TestDbfTable_AddDateField(t *testing.T) {
	g := NewGomegaWithT(t)

	tableUnderTest := New(testEncoding)
	expectedFieldName := "testDate"
	additionError := tableUnderTest.AddDateField(expectedFieldName)
	g.Expect(additionError).To(BeNil())

	g.Expect(tableUnderTest.NumberOfRecords()).To(BeZero())
	g.Expect(len(tableUnderTest.Fields())).To(BeNumerically("==", 1))

	addedField := tableUnderTest.Fields()[0]
	g.Expect(addedField.name).To(Equal(expectedFieldName))
	g.Expect(addedField.fieldType).To(Equal(Date))
}

func TestDbfTable_AddTextField(t *testing.T) {
	g := NewGomegaWithT(t)

	tableUnderTest := New(testEncoding)
	expectedFieldName := "testText"
	expectedFieldLength := byte(20)
	additionError := tableUnderTest.AddTextField(expectedFieldName, expectedFieldLength)
	g.Expect(additionError).To(BeNil())

	g.Expect(tableUnderTest.NumberOfRecords()).To(BeZero())
	g.Expect(len(tableUnderTest.Fields())).To(BeNumerically("==", 1))

	addedField := tableUnderTest.Fields()[0]
	g.Expect(addedField.name).To(Equal(expectedFieldName))
	g.Expect(addedField.fieldType).To(Equal(Character))
	g.Expect(addedField.length).To(Equal(expectedFieldLength))
}

func TestDbfTable_AddNumericField(t *testing.T) {
	g := NewGomegaWithT(t)

	tableUnderTest := New(testEncoding)
	expectedFieldName := "testNumber"
	expectedFieldLength := byte(20)
	expectedFDecimalPlaces := byte(2)
	additionError := tableUnderTest.AddNumberField(expectedFieldName, expectedFieldLength, expectedFDecimalPlaces)
	g.Expect(additionError).To(BeNil())

	g.Expect(tableUnderTest.NumberOfRecords()).To(BeZero())
	g.Expect(len(tableUnderTest.Fields())).To(BeNumerically("==", 1))

	addedField := tableUnderTest.Fields()[0]
	g.Expect(addedField.name).To(Equal(expectedFieldName))
	g.Expect(addedField.fieldType).To(Equal(Numeric))
	g.Expect(addedField.length).To(Equal(expectedFieldLength))
	g.Expect(addedField.decimalPlaces).To(Equal(expectedFDecimalPlaces))
}

func TestDbfTable_AddFloatField(t *testing.T) {
	g := NewGomegaWithT(t)

	tableUnderTest := New(testEncoding)
	expectedFieldName := "testFloat"
	expectedFieldLength := byte(20)
	expectedFDecimalPlaces := byte(2)
	additionError := tableUnderTest.AddFloatField(expectedFieldName, expectedFieldLength, expectedFDecimalPlaces)
	g.Expect(additionError).To(BeNil())

	g.Expect(tableUnderTest.NumberOfRecords()).To(BeZero())
	g.Expect(len(tableUnderTest.Fields())).To(BeNumerically("==", 1))

	addedField := tableUnderTest.Fields()[0]
	g.Expect(addedField.name).To(Equal(expectedFieldName))
	g.Expect(addedField.fieldType).To(Equal(Float))
	g.Expect(addedField.length).To(Equal(expectedFieldLength))
	g.Expect(addedField.decimalPlaces).To(Equal(expectedFDecimalPlaces))
}

func TestDbfTable_FieldNames(t *testing.T) {
	g := NewGomegaWithT(t)

	tableUnderTest := New(testEncoding)

	expectedFieldNames := []string{"first", "second"}

	for _, name := range expectedFieldNames {
		additionError := tableUnderTest.AddBooleanField(name)
		g.Expect(additionError).To(BeNil())
	}

	fieldNamesUnderTest := tableUnderTest.FieldNames()
	g.Expect(fieldNamesUnderTest).To(Equal(expectedFieldNames))
}

func TestDbfTable_DecimalPlacesInField_ValidField(t *testing.T) {
	g := NewGomegaWithT(t)

	tableUnderTest := New(testEncoding)

	numberFieldName := "numField"
	expectedNumberDecimalPlaces := uint8(0)
	tableUnderTest.AddNumberField(numberFieldName, 5, expectedNumberDecimalPlaces)
	actualNumberDecimalPlaces, numberError := tableUnderTest.DecimalPlacesInField(numberFieldName)

	g.Expect(numberError).To(BeNil())
	g.Expect(actualNumberDecimalPlaces).To(BeNumerically("==", expectedNumberDecimalPlaces))

	floatFieldName := "floatField"
	expectedFloatDecimalPlaces := uint8(2)
	tableUnderTest.AddFloatField(floatFieldName, 10, expectedFloatDecimalPlaces)
	actualFloatDecimalPlaces, floatError := tableUnderTest.DecimalPlacesInField(floatFieldName)

	g.Expect(floatError).To(BeNil())
	g.Expect(actualFloatDecimalPlaces).To(BeNumerically("==", expectedFloatDecimalPlaces))
}

func TestDbfTable_DecimalPlacesInField_NonExistentField(t *testing.T) {
	g := NewGomegaWithT(t)

	tableUnderTest := New(testEncoding)

	_, numberError := tableUnderTest.DecimalPlacesInField("missingField")

	g.Expect(numberError).ToNot(BeNil())
	t.Log(numberError)
}

func TestDbfTable_DecimalPlacesInField_InvalidField(t *testing.T) {
	g := NewGomegaWithT(t)

	tableUnderTest := New(testEncoding)

	textFieldName := "textField"
	tableUnderTest.AddTextField(textFieldName, 5)
	_, numberError := tableUnderTest.DecimalPlacesInField(textFieldName)

	g.Expect(numberError).ToNot(BeNil())
	t.Log(numberError)
}

func TestDbfTable_GetRowAsSlice_InitiallyEmptyStrings(t *testing.T) {
	g := NewGomegaWithT(t)

	tableUnderTest := New(testEncoding)

	booldFieldName := "boolField"
	tableUnderTest.AddBooleanField(booldFieldName)

	textFieldName := "textField"
	tableUnderTest.AddBooleanField(textFieldName)

	dateFieldName := "dateField"
	tableUnderTest.AddBooleanField(dateFieldName)

	numFieldName := "numField"
	tableUnderTest.AddBooleanField(numFieldName)

	floatFieldName := "floatField"
	tableUnderTest.AddBooleanField(floatFieldName)

	recordIndex, _ := tableUnderTest.AddNewRecord()

	fieldValues := tableUnderTest.GetRowAsSlice(recordIndex)

	for _, value := range fieldValues {
		g.Expect(value).To(Equal(""))
	}
}

func TestDbfTable_GetRowAsSlice(t *testing.T) {
	g := NewGomegaWithT(t)

	tableUnderTest := New(testEncoding)

	boolFieldName := "boolField"
	expectedBoolFieldValue := "T"
	tableUnderTest.AddBooleanField(boolFieldName)

	textFieldName := "textField"
	expectedTextFieldValue := "some text"
	tableUnderTest.AddTextField(textFieldName, 10)

	dateFieldName := "dateField"
	expectedDateFieldValue := "20181201"
	tableUnderTest.AddDateField(dateFieldName)

	numFieldName := "numField"
	expectedNumFieldValue := "640"
	tableUnderTest.AddNumberField(numFieldName, 3, 0)

	floatFieldName := "floatField"
	expectedFloatFieldValue := "640.42"
	tableUnderTest.AddFloatField(floatFieldName, 6, 2)

	recordIndex, _ := tableUnderTest.AddNewRecord()

	tableUnderTest.SetFieldValueByName(recordIndex, boolFieldName, expectedBoolFieldValue)
	tableUnderTest.SetFieldValueByName(recordIndex, textFieldName, expectedTextFieldValue)
	tableUnderTest.SetFieldValueByName(recordIndex, dateFieldName, expectedDateFieldValue)
	tableUnderTest.SetFieldValueByName(recordIndex, numFieldName, expectedNumFieldValue)
	tableUnderTest.SetFieldValueByName(recordIndex, floatFieldName, expectedFloatFieldValue)

	fieldValues := tableUnderTest.GetRowAsSlice(recordIndex)

	g.Expect(fieldValues[0]).To(Equal(expectedBoolFieldValue))
	g.Expect(fieldValues[1]).To(Equal(expectedTextFieldValue))
	g.Expect(fieldValues[2]).To(Equal(expectedDateFieldValue))
	g.Expect(fieldValues[3]).To(Equal(expectedNumFieldValue))
	g.Expect(fieldValues[4]).To(Equal(expectedFloatFieldValue))
}

func TestDbfTable_FieldValueByName(t *testing.T) {
	g := NewGomegaWithT(t)

	tableUnderTest := New(testEncoding)

	boolFieldName := "boolField"
	expectedBoolFieldValue := "T"
	tableUnderTest.AddBooleanField(boolFieldName)

	textFieldName := "textField"
	expectedTextFieldValue := "some text"
	tableUnderTest.AddTextField(textFieldName, 10)

	dateFieldName := "dateField"
	expectedDateFieldValue := "20181201"
	tableUnderTest.AddDateField(dateFieldName)

	numFieldName := "numField"
	expectedNumFieldValue := "640"
	tableUnderTest.AddNumberField(numFieldName, 3, 0)

	floatFieldName := "floatField"
	expectedFloatFieldValue := "640.42"
	tableUnderTest.AddFloatField(floatFieldName, 6, 2)

	recordIndex, _ := tableUnderTest.AddNewRecord()

	tableUnderTest.SetFieldValueByName(recordIndex, boolFieldName, expectedBoolFieldValue)
	tableUnderTest.SetFieldValueByName(recordIndex, textFieldName, expectedTextFieldValue)
	tableUnderTest.SetFieldValueByName(recordIndex, dateFieldName, expectedDateFieldValue)
	tableUnderTest.SetFieldValueByName(recordIndex, numFieldName, expectedNumFieldValue)
	tableUnderTest.SetFieldValueByName(recordIndex, floatFieldName, expectedFloatFieldValue)

	g.Expect(tableUnderTest.FieldValueByName(recordIndex, boolFieldName)).To(Equal(expectedBoolFieldValue))
	g.Expect(tableUnderTest.FieldValueByName(recordIndex, textFieldName)).To(Equal(expectedTextFieldValue))
	g.Expect(tableUnderTest.FieldValueByName(recordIndex, dateFieldName)).To(Equal(expectedDateFieldValue))
	g.Expect(tableUnderTest.FieldValueByName(recordIndex, numFieldName)).To(Equal(expectedNumFieldValue))
	g.Expect(tableUnderTest.FieldValueByName(recordIndex, floatFieldName)).To(Equal(expectedFloatFieldValue))
}

func TestDbfTable_FieldValueByName_NonExistentField(t *testing.T) {
	g := NewGomegaWithT(t)

	tableUnderTest := New(testEncoding)

	textFieldName := "textField"
	tableUnderTest.AddTextField(textFieldName, 10)

	_, valueError := tableUnderTest.FieldValueByName(0, "missingField")

	g.Expect(valueError).ToNot(BeNil())
	t.Log(valueError)
}

func TestDbfTable_SetFieldValueByName_NonExistentField(t *testing.T) {
	g := NewGomegaWithT(t)

	tableUnderTest := New(testEncoding)

	setError := tableUnderTest.SetFieldValueByName(0, "missingField", "someText")

	g.Expect(setError).ToNot(BeNil())
	t.Log(setError)
}

func TestDbfTable_AddRecordWithNoFieldsDefined_Errors(t *testing.T) {
	g := NewGomegaWithT(t)

	tableUnderTest := New(testEncoding)

	recordIndex, addErr := tableUnderTest.AddNewRecord()
	g.Expect(addErr).ToNot(BeNil())
	g.Expect(recordIndex).To(BeEquivalentTo(-1))
}

func TestDbfTable_Int64FieldValueByName(t *testing.T) {
	g := NewGomegaWithT(t)

	tableUnderTest := New(testEncoding)

	intFieldName := "intField"
	expectedIntValue := 640
	expectedIntFieldValue := fmt.Sprintf("%d", expectedIntValue)
	tableUnderTest.AddNumberField(intFieldName, 6, 2)

	recordIndex, addErr := tableUnderTest.AddNewRecord()
	g.Expect(addErr).To(BeNil())

	tableUnderTest.SetFieldValueByName(recordIndex, intFieldName, expectedIntFieldValue)

	actualIntFieldValue, valueError := tableUnderTest.Int64FieldValueByName(recordIndex, intFieldName)

	g.Expect(valueError).To(BeNil())
	g.Expect(actualIntFieldValue).To(BeNumerically("==", expectedIntValue))
}

func TestDbfTable_Float64FieldValueByName(t *testing.T) {
	g := NewGomegaWithT(t)

	tableUnderTest := New(testEncoding)

	floatFieldName := "floatField"
	expectedFloatValue := 640.42
	expectedFloatFieldValue := fmt.Sprintf("%.2f", expectedFloatValue)
	tableUnderTest.AddFloatField(floatFieldName, 10, 2)

	recordIndex, addErr := tableUnderTest.AddNewRecord()
	g.Expect(addErr).To(BeNil())

	tableUnderTest.SetFieldValueByName(recordIndex, floatFieldName, expectedFloatFieldValue)

	actualFloatFieldValue, valueError := tableUnderTest.Float64FieldValueByName(recordIndex, floatFieldName)

	g.Expect(valueError).To(BeNil())
	g.Expect(actualFloatFieldValue).To(BeNumerically("==", expectedFloatValue))
}

func TestDbfTable_FieldDescriptor(t *testing.T) {
	g := NewGomegaWithT(t)

	tableUnderTest := New(testEncoding)

	const fieldName = "floatField"
	const fieldLength = uint8(10)
	const decimalPlaces = uint8(2)

	floatFieldName := fieldName
	tableUnderTest.AddFloatField(floatFieldName, fieldLength, decimalPlaces)

	fieldUnderTest := tableUnderTest.Fields()[0]

	g.Expect(fieldUnderTest.Name()).To(Equal(fieldName))
	g.Expect(fieldUnderTest.FieldType()).To(Equal(Float))
	g.Expect(fieldUnderTest.FieldType()).To(Equal(Float))
	g.Expect(fieldUnderTest.Length()).To(Equal(fieldLength))
	g.Expect(fieldUnderTest.DecimalPlaces()).To(Equal(decimalPlaces))
}

func TestDbfTable_RowIsDeleted(t *testing.T) {
	g := NewGomegaWithT(t)

	tableUnderTest := New(testEncoding)

	const textFieldName = "textField"
	tableUnderTest.AddTextField(textFieldName, 10)

	const floatFieldName = "floatField"
	const fieldLength = uint8(10)
	const decimalPlaces = uint8(2)
	tableUnderTest.AddFloatField(floatFieldName, fieldLength, decimalPlaces)

	recordIndex, _ := tableUnderTest.AddNewRecord()

	tableUnderTest.SetFieldValueByName(recordIndex, textFieldName, "some text")
	tableUnderTest.SetFieldValueByName(recordIndex, textFieldName, "0")

	g.Expect(tableUnderTest.RowIsDeleted(recordIndex)).To(BeFalse())

	deletionError := tableUnderTest.SetRowIsDeleted(recordIndex)

	g.Expect(deletionError).To(BeNil())
	g.Expect(tableUnderTest.RowIsDeleted(recordIndex)).To(BeTrue())
}

func TestDbfTable_RowIsDeleted_InvalidRow_Errors(t *testing.T) {
	g := NewGomegaWithT(t)

	tableUnderTest := New(testEncoding)

	isDeleted, deletionError := tableUnderTest.RowIsDeleted(1)
	g.Expect(isDeleted).To(BeFalse())
	g.Expect(deletionError).ToNot(BeNil())
}

func TestDbfTable_MiddleRowIsDeleted_(t *testing.T) {
	g := NewGomegaWithT(t)

	tableUnderTest := New(testEncoding)

	const textFieldName = "textField"
	tableUnderTest.AddTextField(textFieldName, 10)

	const floatFieldName = "floatField"
	const fieldLength = uint8(10)
	const decimalPlaces = uint8(2)
	tableUnderTest.AddFloatField(floatFieldName, fieldLength, decimalPlaces)

	recordIndex, _ := tableUnderTest.AddNewRecord()
	tableUnderTest.SetFieldValueByName(recordIndex, textFieldName, "some text 1")
	tableUnderTest.SetFieldValueByName(recordIndex, textFieldName, "1")

	recordIndex, _ = tableUnderTest.AddNewRecord()
	tableUnderTest.SetFieldValueByName(recordIndex, textFieldName, "some text 2")
	tableUnderTest.SetFieldValueByName(recordIndex, textFieldName, "2")

	recordIndex, _ = tableUnderTest.AddNewRecord()
	tableUnderTest.SetFieldValueByName(recordIndex, textFieldName, "some text 3")
	tableUnderTest.SetFieldValueByName(recordIndex, textFieldName, "3")

	deletionError := tableUnderTest.SetRowIsDeleted(1)

	g.Expect(deletionError).To(BeNil())
	g.Expect(tableUnderTest.RowIsDeleted(0)).To(BeFalse())
	g.Expect(tableUnderTest.RowIsDeleted(1)).To(BeTrue())
	g.Expect(tableUnderTest.RowIsDeleted(2)).To(BeFalse())

	deletionError = tableUnderTest.SetRowIsDeleted(2)

	g.Expect(deletionError).To(BeNil())
	g.Expect(tableUnderTest.RowIsDeleted(0)).To(BeFalse())
	g.Expect(tableUnderTest.RowIsDeleted(1)).To(BeTrue())
	g.Expect(tableUnderTest.RowIsDeleted(2)).To(BeTrue())

	deletionError = tableUnderTest.SetRowIsDeleted(0)

	g.Expect(deletionError).To(BeNil())
	g.Expect(tableUnderTest.RowIsDeleted(0)).To(BeTrue())
	g.Expect(tableUnderTest.RowIsDeleted(1)).To(BeTrue())
	g.Expect(tableUnderTest.RowIsDeleted(2)).To(BeTrue())
}

// dBASE III Compliance Tests

// TestDbfTable_DbaseIII_RecordInitializedWithSpaces verifies that new records
// are initialized with spaces (0x20) per dBASE III spec, not zeros.
func TestDbfTable_DbaseIII_RecordInitializedWithSpaces(t *testing.T) {
	g := NewGomegaWithT(t)

	tableUnderTest := New(testEncoding)
	tableUnderTest.AddTextField("textField", 10)
	tableUnderTest.AddNumberField("numField", 5, 0)
	tableUnderTest.AddBooleanField("boolField")
	tableUnderTest.AddDateField("dateField")

	recordIndex, _ := tableUnderTest.AddNewRecord()

	// Get raw bytes for the record
	recordOffset := int(tableUnderTest.numberOfBytesInHeader) + recordIndex*int(tableUnderTest.lengthOfEachRecord)
	recordBytes := tableUnderTest.dataStore[recordOffset : recordOffset+int(tableUnderTest.lengthOfEachRecord)]

	// First byte should be 0x20 (active record marker)
	g.Expect(recordBytes[0]).To(Equal(byte(0x20)), "Record deletion flag should be 0x20 (active)")

	// All field bytes should be spaces (0x20), not zeros
	for i := 1; i < len(recordBytes); i++ {
		g.Expect(recordBytes[i]).To(Equal(byte(0x20)), fmt.Sprintf("Byte at offset %d should be 0x20 (space), got 0x%02X", i, recordBytes[i]))
	}
}

// TestDbfTable_DbaseIII_LogicalFieldInitialization verifies that logical fields
// are initialized to space (0x20), not 'F' or 0x00.
func TestDbfTable_DbaseIII_LogicalFieldInitialization(t *testing.T) {
	g := NewGomegaWithT(t)

	tableUnderTest := New(testEncoding)
	tableUnderTest.AddBooleanField("boolField")

	recordIndex, _ := tableUnderTest.AddNewRecord()

	// Get raw byte for the logical field
	recordOffset := int(tableUnderTest.numberOfBytesInHeader) + recordIndex*int(tableUnderTest.lengthOfEachRecord)
	logicalFieldByte := tableUnderTest.dataStore[recordOffset+1] // +1 to skip deletion flag

	// Per dBASE III spec, uninitialized logical should be space (0x20)
	g.Expect(logicalFieldByte).To(Equal(byte(0x20)), "Logical field should be initialized to space (0x20)")
}

// TestDbfTable_DbaseIII_DateFieldEmpty verifies that empty date fields
// contain 8 spaces, not "00000000".
func TestDbfTable_DbaseIII_DateFieldEmpty(t *testing.T) {
	g := NewGomegaWithT(t)

	tableUnderTest := New(testEncoding)
	tableUnderTest.AddDateField("dateField")

	recordIndex, _ := tableUnderTest.AddNewRecord()

	// Get raw bytes for the date field (8 bytes)
	recordOffset := int(tableUnderTest.numberOfBytesInHeader) + recordIndex*int(tableUnderTest.lengthOfEachRecord)
	dateFieldBytes := tableUnderTest.dataStore[recordOffset+1 : recordOffset+1+8] // +1 to skip deletion flag

	// Per dBASE III spec, empty date should be 8 spaces
	expectedDateBytes := []byte("        ") // 8 spaces
	g.Expect(dateFieldBytes).To(Equal(expectedDateBytes), "Empty date field should be 8 spaces")
}

// TestDbfTable_DbaseIII_NumericFieldEmpty verifies that empty numeric fields
// contain spaces, not zeros or "0".
func TestDbfTable_DbaseIII_NumericFieldEmpty(t *testing.T) {
	g := NewGomegaWithT(t)

	tableUnderTest := New(testEncoding)
	tableUnderTest.AddNumberField("numField", 10, 2)

	recordIndex, _ := tableUnderTest.AddNewRecord()

	// Get raw bytes for the numeric field (10 bytes)
	recordOffset := int(tableUnderTest.numberOfBytesInHeader) + recordIndex*int(tableUnderTest.lengthOfEachRecord)
	numFieldBytes := tableUnderTest.dataStore[recordOffset+1 : recordOffset+1+10]

	// Per dBASE III spec, empty numeric should be all spaces
	for i, b := range numFieldBytes {
		g.Expect(b).To(Equal(byte(0x20)), fmt.Sprintf("Numeric field byte %d should be space (0x20), got 0x%02X", i, b))
	}
}

// TestDbfTable_DbaseIII_CharacterFieldEmpty verifies that empty character fields
// contain spaces.
func TestDbfTable_DbaseIII_CharacterFieldEmpty(t *testing.T) {
	g := NewGomegaWithT(t)

	tableUnderTest := New(testEncoding)
	tableUnderTest.AddTextField("textField", 20)

	recordIndex, _ := tableUnderTest.AddNewRecord()

	// Get raw bytes for the character field
	recordOffset := int(tableUnderTest.numberOfBytesInHeader) + recordIndex*int(tableUnderTest.lengthOfEachRecord)
	textFieldBytes := tableUnderTest.dataStore[recordOffset+1 : recordOffset+1+20]

	// Per dBASE III spec, empty character field should be all spaces
	for i, b := range textFieldBytes {
		g.Expect(b).To(Equal(byte(0x20)), fmt.Sprintf("Character field byte %d should be space (0x20), got 0x%02X", i, b))
	}
}

// TestDbfTable_DbaseIII_EOFMarkerAtEnd verifies that EOF marker (0x1A) is always
// at the end of the file, not before records.
func TestDbfTable_DbaseIII_EOFMarkerAtEnd(t *testing.T) {
	g := NewGomegaWithT(t)

	tableUnderTest := New(testEncoding)
	tableUnderTest.AddTextField("textField", 10)

	// Add multiple records
	for i := 0; i < 3; i++ {
		tableUnderTest.AddNewRecord()
	}

	// EOF marker should be the last byte
	lastByte := tableUnderTest.dataStore[len(tableUnderTest.dataStore)-1]
	g.Expect(lastByte).To(Equal(byte(0x1A)), "EOF marker (0x1A) should be at the end of dataStore")

	// Verify there's only one EOF marker (not one before each record)
	eofCount := 0
	for _, b := range tableUnderTest.dataStore {
		if b == 0x1A {
			eofCount++
		}
	}
	g.Expect(eofCount).To(Equal(1), "There should be exactly one EOF marker")
}

// TestDbfTable_DbaseIII_FileSignature verifies the file signature is 0x03 for dBASE III.
func TestDbfTable_DbaseIII_FileSignature(t *testing.T) {
	g := NewGomegaWithT(t)

	tableUnderTest := New(testEncoding)
	tableUnderTest.AddTextField("textField", 10)

	// File signature should be 0x03 (dBASE III without memo)
	g.Expect(tableUnderTest.dataStore[0]).To(Equal(byte(0x03)), "File signature should be 0x03 for dBASE III")
}

// TestDbfTable_DbaseIII_NumericRightAligned verifies numeric values are right-aligned
// with spaces on the left.
func TestDbfTable_DbaseIII_NumericRightAligned(t *testing.T) {
	g := NewGomegaWithT(t)

	tableUnderTest := New(testEncoding)
	tableUnderTest.AddNumberField("numField", 10, 0)

	recordIndex, _ := tableUnderTest.AddNewRecord()
	tableUnderTest.SetFieldValueByName(recordIndex, "numField", "42")

	// Get raw bytes for the numeric field
	recordOffset := int(tableUnderTest.numberOfBytesInHeader) + recordIndex*int(tableUnderTest.lengthOfEachRecord)
	numFieldBytes := tableUnderTest.dataStore[recordOffset+1 : recordOffset+1+10]

	// "42" should be right-aligned: "        42"
	expectedBytes := []byte("        42")
	g.Expect(numFieldBytes).To(Equal(expectedBytes), "Numeric '42' should be right-aligned with spaces on left")
}

// TestDbfTable_DbaseIII_HeaderReservedBytesZero verifies that reserved header bytes
// (offset 12-31, except 28-29) are zeros.
func TestDbfTable_DbaseIII_HeaderReservedBytesZero(t *testing.T) {
	g := NewGomegaWithT(t)

	tableUnderTest := New(testEncoding)
	tableUnderTest.AddTextField("textField", 10)

	// Bytes 12-27 should be zeros (reserved)
	for i := 12; i <= 27; i++ {
		g.Expect(tableUnderTest.dataStore[i]).To(Equal(byte(0x00)),
			fmt.Sprintf("Reserved header byte at offset %d should be 0x00", i))
	}

	// Byte 28 is table flags, should be 0x00 for no memo
	g.Expect(tableUnderTest.dataStore[28]).To(Equal(byte(0x00)), "Table flags (byte 28) should be 0x00")

	// Bytes 30-31 should be zeros (reserved)
	g.Expect(tableUnderTest.dataStore[30]).To(Equal(byte(0x00)), "Reserved byte 30 should be 0x00")
	g.Expect(tableUnderTest.dataStore[31]).To(Equal(byte(0x00)), "Reserved byte 31 should be 0x00")
}
