package metadata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSourceSuccesfullyForSingleSourceContent(t *testing.T) {
	testContent := Content{
		UUID: "7560aca3-986c-487b-8f9f-6b865872096f",
		Identifiers: []Identifier{
			{
				Authority: "http://api.ft.com/system/FTCOM-METHODE",
			},
		},
	}

	actual, err := testContent.getSource()
	assert.NoError(t, err, "Error while getting content source")
	assert.Equal(t, "METHODE", actual, "Actual source is different from expected source")
}

func TestGetSourceSuccesfullyForMultipleSourcesContent(t *testing.T) {
	testContent := Content{
		UUID: "9cc74217-7690-35be-a0d6-683d118561d4",
		Identifiers: []Identifier{
			{
				Authority: "http://api.ft.com/system/FT-LABS-WP-1-335",
			},
			{
				Authority: "http://api.ft.com/system/FT-LABS-WP-1-335",
			},
		},
	}

	actual, err := testContent.getSource()
	assert.NoError(t, err, "Error while getting content source")
	assert.Equal(t, "BLOGS", actual, "Actual source is different from expected source")
}

func TestGetSourceSuccesfullyForMultipleDifferentSourcesContent(t *testing.T) {
	testContent := Content{
		UUID: "9cc74217-7690-35be-a0d6-683d118561d4",
		Identifiers: []Identifier{
			{
				Authority: "http://api.ft.com/system/FT-LABS-WP-1-335",
			},
			{
				Authority: "http://api.ft.com/system/FTCOM-METHODE",
			},
		},
	}

	_, err := testContent.getSource()
	assert.Error(t, err, "Expecting error but no error was found")
}
