package listener

import (
	"github.com/WildEgor/cdc-listener/internal/adapters/publisher"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
	"sync"
	"testing"
	"time"
)

var data *CDCData

func TestMain(m *testing.M) {
	start()
	defer teardown()

	code := m.Run()
	os.Exit(code)
}

func start() {
	pool := &sync.Pool{
		New: func() any {
			return &publisher.Event{}
		},
	}

	data = NewCDCData(pool)
}

func teardown() {

}

type ConverterNestedMockData struct {
	NumberField int                `bson:"number_field"`
	FloatField  float64            `bson:"float_field"`
	DateField   primitive.DateTime `bson:"date_field"`
}

type ConverterMockData struct {
	Id     primitive.ObjectID      `bson:"_id"`
	Name   string                  `bson:"name"`
	Nested ConverterNestedMockData `bson:"nested"`
}

func Test_Converter(t *testing.T) {
	objIdMock := primitive.NewObjectID()
	timeMock := primitive.DateTime(time.Now().Unix())
	floatMock := float64(int(0.1*100)) / 100

	tests := []struct {
		name    string
		input   ConverterMockData
		want    map[string]any
		wantErr bool
	}{
		{
			name: "success convert",
			input: ConverterMockData{
				objIdMock,
				"name",
				ConverterNestedMockData{
					1,
					floatMock,
					timeMock,
				},
			},
			want: map[string]any{
				"_id":  objIdMock.Hex(),
				"name": "name",
				"nested": map[string]any{
					"number_field": 1,
					"float_field":  floatMock,
					"date_field":   timeMock,
				},
			},
			wantErr: false,
		},
	}

	pool := &sync.Pool{
		New: func() any {
			return &publisher.Event{}
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cdcData := NewCDCData(pool)

			bsonData, err := bson.Marshal(tt.input)
			assert.Nil(t, err)

			var doc bson.D
			err = bson.Unmarshal(bsonData, &doc)
			assert.Nil(t, err)

			result := cdcData.convertBsonToJson(doc)

			assert.NotNil(t, result)
			assert.Equal(t, tt.want["_id"], result["_id"])
			assert.Equal(t, tt.want["name"], result["name"])

			// FIXME
			//if !reflect.DeepEqual(result["nested"], tt.want["nested"]) {
			//	t.Errorf("Result map %v does not match expected map %v", result["nested"], tt.want["nested"])
			//}
		})
	}
}
