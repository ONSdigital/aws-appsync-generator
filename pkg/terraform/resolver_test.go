package terraform_test

// func TestArgsSource(t *testing.T) {
// 	type tc struct {
// 		scenario string
// 		resolver terraform.Resolver
// 		expected string
// 	}

// 	cases := []tc{
// 		{
// 			scenario: "Query should be 'args'",
// 			resolver: terraform.Resolver{ParentType: "Query"},
// 			expected: "args",
// 		},
// 		{
// 			scenario: "Mutation should be 'args'",
// 			resolver: terraform.Resolver{ParentType: "Mutation"},
// 			expected: "args",
// 		},
// 		{
// 			scenario: "Object should be 'source'",
// 			resolver: terraform.Resolver{ParentType: "Animal"},
// 			expected: "source",
// 		},
// 	}

// 	for _, c := range cases {
// 		assert.Equal(t, c.expected, c.resolver.ArgsSource(), c.scenario)
// 	}
// }

// func TestKeyFieldJSONMapAndList(t *testing.T) {

// 	// Using the same scenarios to test both
// 	// - KeyFieldJSONMap()
// 	// - KeyFieldJSONList()

// 	type tc struct {
// 		scenario     string
// 		resolver     terraform.Resolver
// 		expectedMap  string
// 		expectedList string
// 	}

// 	cases := []tc{
// 		{
// 			scenario:     "No fields",
// 			resolver:     terraform.Resolver{},
// 			expectedMap:  "{}",
// 			expectedList: "[]",
// 		},
// 		{
// 			scenario: "One mandatory field",
// 			resolver: terraform.Resolver{
// 				KeyFields: []manifest.Field{
// 					{Name: "id:ID!"},
// 				},
// 			},
// 			expectedMap:  `{"id":"ID"}`,
// 			expectedList: `["id"]`,
// 		},
// 		{
// 			scenario: "Several fields (with defaults)",
// 			resolver: terraform.Resolver{
// 				KeyFields: []manifest.Field{
// 					{Name: "id:ID!"},
// 					{Name: "color"},
// 					{Name: "size:Int"},
// 				},
// 			},
// 			expectedMap:  `{"id":"ID","color":"String","size":"Int"}`,
// 			expectedList: `["id","color","size"]`,
// 		},
// 	}

// 	for _, c := range cases {
// 		assert.Equal(t, c.expectedMap, c.resolver.KeyFieldJSONMap(), c.scenario+" (map)")
// 		assert.Equal(t, c.expectedList, c.resolver.KeyFieldJSONList(), c.scenario+" (list)")
// 	}
// }
