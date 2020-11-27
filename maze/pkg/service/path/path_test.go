package path

import (
	"context"
	"errors"
	"github.com/avanticaTest/maze/pkg/db"
	"github.com/go-kit/kit/log"
	"github.com/stretchr/testify/mock"
	"gotest.tools/v3/assert"
	"testing"
)

func TestDeletePath(t *testing.T) {

	tests := []struct {
		name         string
		request      string
		response     int
		expectedResp int
		expectedErr  error
		success      bool
		mongoOK      bool
	}{
		{
			name:         "OK",
			request:      "5fbecff95f80a305742abb10",
			response:     123456,
			expectedResp: 123456,
			expectedErr:  nil,
			success:      true,
			mongoOK:      true,
		},
		{
			name:         "Not OK",
			request:      "5fbecff95f80a305742abb10",
			response:     0,
			expectedResp: 0,
			expectedErr:  errors.New("mongo error"),
			success:      false,
			mongoOK:      false,
		},
	}

	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			logger := log.NewNopLogger()
			db := db.Mock{}
			if tt.mongoOK {
				db.On("DeleteOne", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(123456, nil)
			} else {
				db.On("DeleteOne", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(0, errors.New("mongo error"))
			}
			p := New(logger, db)

			resp, err := p.DeletePath(ctx, tt.request)

			assert.DeepEqual(t, tt.expectedResp, resp)
			if tt.success {
				assert.NilError(t, err)
			} else {
				assert.Error(t, err, errors.New("mongo error").Error())
			}

		})
	}

}
