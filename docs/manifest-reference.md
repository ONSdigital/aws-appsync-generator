# Manifest Reference

- [Manifest Reference](#manifest-reference)
  - [Blocks](#blocks)
    - [Sources Block](#sources-block)
    - [Enums Block](#enums-block)
    - [Objects Block](#objects-block)
    - [Queries Block](#queries-block)
    - [Mutations Block](#mutations-block)
  - [Sub-Blocks](#sub-blocks)
    - [Field Sub-Block](#field-sub-block)
    - [Resolver Sub-Block](#resolver-sub-block)

## Blocks

_Blocks_ declare parts of the configuration manifest

---

### Sources Block

The `sources` block defines configuration of appsync data sources

Exactly one of `dynamo` or `sql` subblocks _must_ be supplied

**sources** [Hash, required]: Specifies the data sources available to the appsync api

- **&lt;sourcekey&gt;** [String, required]: Name of the data source that may be referenced in resolvers
  - **name** [String, required]: Name of the data source. When the type is `dynamodb` this will be the table name
  - **dynamo** [Hash, optional]: Describes `dynamodb` specific configuration options.
    - **hash_key** [DynamoKey, required]: Specifies the field to be used as the table `hash key`
      - **name** [String, required]: Name of the field
      - **type** [String, optional]: The dynamodb type of the field (default `S` (string))
    - **sort_key** [DynamoKey, optional]: Specifies the field to be used as the table `sort key`
      - **name** [String, required]: Name of the field
      - **type** [String, optional]: The dynamodb type of the field (default `S` (string))
    - **backup** [Bool, optional]: Sets whether to enable incrementatal backup on the table. Default `false`
      - _be aware, enabling backup has a cost implication, so only use for tables that require it_
  - **sql** [Hash, optional]
    - _not yet implemented_

Example

```yml
sources:

  users:
    name: users
    dynamo:
      hash_key:
        name: email
      backup: true

  customers:
    name: customers
    dynamo:
      hash_key:
        name: login
      sort_key:
        name: priority
        type: N
```

---

### Enums Block

The `enums` block defines configuration of enumeration types for the graphql schema. A schema need not declare any enums

**enums** [Array, optional]: Declares a set of enums

- **name** [String, required]: The type name of the enumeration
- **values** [Array required]: The set of values that make up the enumeration

Example:

```yml
enums:
  - name: Status
    values: [UNREAD,READ,REPLIED,FORWARDED]
  - name: Priority
    values: [LOW,MEDIUM,HIGH]
```

---

### Objects Block

The `objects` block defines graphql object types

**objects** [Array, required]: Declare a set of objects

- **name** [String, required]: The name identifier for the object
- **fields** [Array, required]: Each sub-block specifies a field in the object

See [field](#field-sub-block) for more information

Example:

```yml
objects:
  - name: animal
    fields:
      - name: id
        type: ID!
      - name: name
      - name: isFluffy
        type: Boolean
      - name: favouriteToys
        type: [String]
      - name: license
        type: Licence
        inputType: ID

  - name: License
    fields:
      # ... omitted
```

---

### Queries Block

The `queries` block declares queries to be created in the schema. A schema need not declare any queries.

**queries** [Array, optional]: Declare a set of queries

- **name** [String, required]: Name of the query. No restriction, but by convention prefix with the action (e.g. `get`, `list`)
- **resolver** [Hash, required]: The resolver configuration to satisfy the query

See [resolver](#resolver-sub-block) for more information

Example

```yml
queries:

  - name: getAnimal
    resolver:
      action: get
      type: Animal
      keyFields:
        - name: id
          type: ID

  - name: listAnimals
    resolver:
      action: list
      type: Animal
```

---

### Mutations Block

The `mutations` block declares mutations to be created in the schema. A schema need not declare any mutations.

**mutations** [Array, optional]: Declare a set of mutations. Each mutation definition is the same structure as used for [queries](#queries-block)

Example

```yml
mutations:
  - name: createAnimal
    resolver:
      action: insert
      type: Animal
      keyFields:
        - name: id
          type: ID`
```

---

## Sub-Blocks

_Sub-Blocks_ declare smaller resuable chunks of configuration

### Field Sub-Block

The `field` sub-block is a field inside another type (`object`, `resolver` etc)

- **name**: [String, required]: The name of the field
- **type**: [String, optional]: The `type` of the field. May be a normal graphql `scalar type`, a `graphql object type` or `AWS scalar type`. If not specified, will default to `String`.
  - _If an exclaimation mark is provided, this denotes the field as `non nullable`_
  - _If surrounded with square brackets, this denotes the field as an `array` type_
- **inputType**: [String, optional]: Only applies to fields used in [object blocks](#objects-block). If present, will override the field type in any associated generated `input object`. This is useful if you want to return a nested type when reading an object, but only specify an ID to object when creating it

Example:

```yml
objects:
  - name: anobject
    fields:

      # An ID field with a non-nullable marker
      - name: id
        type: ID!

      # A field defaulting to a String type
      - name: email

      # A field with object type, overriden for input
      - name: appointment
        type: Apppointment
        inputType: ID

      # An array type
      - name: cc
        type: [String]
```

---

### Resolver Sub-Block

The `resolver` sub-block is a resolver for a [query](#queries-block) or a [mutation](#mutations-block)

**resolver** [Hash, required]

- **action** [String, required]: Defines the action the resolver should take. Must be one of `get`,`list`,`update`,`delete` or `insert`
- **type** [String, required]: The `type` returned by the resolver.
  - _If the resolver is of a kind that returns multiple values, this will automatically become an array. There is no need to mark up the type with square brackets_
- **source** [String, optional]: If present, must be `source key` as declared in the [sources](#sources-block) block. If omitted it will be set to the _default_ `source key` (if one has been declared)
- **keyFields** [Array, optional]: Used to denote which field (defined in the type being returned) to use as the look up key fields. This will become a mandatory field in the query/mutation definition
  - _Not applicable to `list` action types_
  - _Each field is [field](#field-sub-block) sub-block_

Additional resources will be created in the schema appropriate to the the action specified (e.g. input and filter object types)

Example:

```yml
  # In a get query
  resolver:
    action: get
    type: Animal
    keyFields:
      - name: id
        type: ID

  # In a list query with a specific source
  resolver:
    action: list
    type: Animal
    source: ZooAnimals

  # In a create mutation
  resolver:
    action: create
    type: Animal
    keyFields:
      - name: id
        type: ID
```
