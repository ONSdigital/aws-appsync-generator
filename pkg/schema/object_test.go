package schema_test

// func TestQuery(t *testing.T) {
// 	s := schema.MustParse(definition)
// 	o := s.Objects[0]

// 	// Test missing query
// 	qBad, err := o.Query("aardvark")
// 	assert.Error(t, err, "Call should error")
// 	assert.Nil(t, qBad, "Query should be returned nil")

// 	// Test existing query
// 	expected := &schema.Query{
// 		Name: "correspondence",
// 		Type: "get",
// 		Params: []schema.Field{
// 			{
// 				Name:      "reference",
// 				Type:      "ID",
// 				Mandatory: true,
// 			},
// 		},
// 		ParamString: "reference: ID!",
// 		Returns:     "Correspondence",
// 	}

// 	q, err := o.Query("correspondence")
// 	assert.NoError(t, err, "Call should not error")
// 	assert.NotNil(t, q, "Query should not be nil")
// 	assert.Equal(t, expected, q, "Query should be as expected")
// }

// func TestField(t *testing.T) {
// 	s := schema.MustParse(definition)
// 	o := s.Objects[0]

// 	// Test valid fields
// 	cases := []*schema.Field{
// 		{
// 			// Mandatory ID type
// 			Name:      "reference",
// 			Type:      "ID",
// 			Mandatory: true,
// 		},
// 		{
// 			// Standard string - default scalar type
// 			Name:      "subject",
// 			Type:      "String",
// 			Mandatory: false,
// 		},
// 	}

// 	for _, expected := range cases {
// 		f, err := o.Field(expected.Name)
// 		assert.NoError(t, err, "Call should not error")
// 		assert.NotNil(t, f, "Field should not be nil")
// 		assert.Equal(t, expected, f, "Field should be as expected")
// 	}

// 	// Test missing fields
// 	f, err := o.Field("notARealField")
// 	assert.Error(t, err, "Call should error")
// 	assert.Nil(t, f, "Field should be nil")
// }
