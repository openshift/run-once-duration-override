package runoncedurationoverride

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertExternalConfig(t *testing.T) {
	external := &RunOnceDurationOverride{
		Spec: RunOnceDurationOverrideSpec{
			ActiveDeadlineSeconds: 10,
		},
	}

	configGot := ConvertExternalConfig(external)
	assert.NotNil(t, configGot)
	assert.Equal(t, int64(10), configGot.ActiveDeadlineSeconds)
}

func TestDecodeWithFile(t *testing.T) {
	tests := []struct {
		name   string
		file   string
		assert func(t *testing.T, objGot *RunOnceDurationOverride, errGot error)
	}{
		{
			name: "WithValidObject",
			file: "testdata/external.yaml",
			assert: func(t *testing.T, objGot *RunOnceDurationOverride, errGot error) {
				assert.NoError(t, errGot)
				assert.NotNil(t, objGot)

				assert.Equal(t, int64(10), objGot.Spec.ActiveDeadlineSeconds)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			objGot, errGot := DecodeWithFile(tt.file)

			tt.assert(t, objGot, errGot)
		})
	}
}
