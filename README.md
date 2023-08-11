## Credhub plugin

A CF client plugin that allows a user to easily interact with the secrets maintenance api from a compatible credhub service broker (https://github.com/rabobank/credhub-service-broker >v1.0.6)

## Intro

The secrets maintenance API provides 5 functionalities to interact with the credentials of an existing credhub service instance:

* **List Secrets**

  This will list all keys present in the credentials stored in an existing credhub service instance. If the stored credentials has encapsulated objects,
keys will reflect the encapsulation by concatenating the child keys to the parent keys by means of a "." character (i.e. the listed key for a credential
like { "a" : { "b" : "value" } } will be "a.b")

* **Add Secret**
  
  Allows to add one or more secrets to an existing service instance. Values with conflicting keys will overwrite the existing values even if they
are of different types (so, even if the new value is a string, for example, and the existing value was a sub-map, the value will be overwritten as
a simple value)

* **Delete Secret**

  Deletes individual secrets from the credentials referencing them by the keys (as presented when listing the secrets).

* **List Versions**

  This will list up to 20 latest versions of an existing credhub service instance. The list will show the date of the creation/update and
an ID which can be used to reinstate a previous version (rollback)

* **Reinstate Version**

  By providing an ID as listed by the "List Versions" functionality, a value can be reinstated or an update can be rolled back by reinstating
the previous version. Reinstating a version will not work by providing a version id belonging to another service. Only ids from versions of
the same service.

## List Secret

```
NAME:
   list-credhub-secrets - List all secret keys in the credhub service instance

USAGE:
   
cf list-credhub-secrets <SERVICE_INSTANCE>

  SERVICE_INSTANCE - Credhub service instance name.


ALIAS:
   lcs
```
### example

```
$ cf lcs test-broker-api

output:
   a
   b.a
   b.c
   doe-iets-leuks
```

***NOTE*** List of keys is not sorted
## Add Secret

```
NAME:
   add-credhub-secrets - Add secrets to credhub service

USAGE:
   
cf add-credhub-secrets <SERVICE_INSTANCE> <JSON_OBJECT>
cf add-credhub-secrets <SERVICE_INSTANCE> <KEY> <VALUE>

  SERVICE_INSTANCE - Credhub service instance name the keys are being added to.

  JSON_OBJECT      - A well formed json object map. Key values will either replace existing keys or added to the existing credentials if not present
                     This will only be interpreted as a json object it the KEY/VALUE parameters are not provided.
  KEY              - When a VALUE is provided, instead of JSON_OBJECT, the first parameter will be interpreted as the secret key.
                     If updating/setting encapsulated values, dots may be used to reference the inner-keys (i.e. a.b to reference {"a":{"b":"value"}})
  VALUE            - Secret value.


ALIAS:
   acs
```
### example

```
$ cf acs test-broker-api c something
$ cf acs test-broker-api d.a something
$ cf acs test-broker-api '{ "b": {"b": "value"}, "f" : "value"}'

$ cf lcs test-broker-api 

output:
   c
   d.a
   doe-iets-leuks
   f
   a
   b.a
   b.b
   b.c
```

## Delete Secret

```
NAME:
   delete-credhub-secrets - Delete a key from the credhub service instance

USAGE:
   
cf delete-credhub-secrets <SERVICE_INSTANCE> <KEYS>...

  SERVICE_INSTANCE - Credhub service instance name the keys are being deleted from.

  KEYS             - Secret keys to delete. Multiple keys can be provided separated by spaces.


ALIAS:
   dcs
```

### example

```
$ cf dcs test-broker-api d b.b a

$ cf lcs test-broker-api 

output:
   doe-iets-leuks
   f
   b.a
   b.c
   c
```

## List Versions

```
NAME:
   list-credhub-secrets-versions - List up to 20 latest versions for a credhub service instance credentials

USAGE:
   
cf list-credhub-secrets-versions <SERVICE_INSTANCE>

  SERVICE_INSTANCE - Credhub service instance name.


ALIAS:
   lcv
```

### example

```
$ cf lcv test-broker-api

output:

ID: 03142a68-254b-4739-a7a7-abdae54ec40d
Created: 2023-08-11 15:41:33 +0000 UTC

   doe-iets-leuks
   f
   b.a
   b.c
   c

ID: 5942167f-8d02-4b03-a1b8-1116c41069f3
Created: 2023-08-11 15:37:08 +0000 UTC

   a
   b.a
   b.b
   b.c
   c
   d.a
   doe-iets-leuks
   f

ID: 836bfc7b-1d75-4a0d-98c8-1772ff318709
Created: 2023-08-11 15:35:39 +0000 UTC

   d.a
   doe-iets-leuks
   a
   b.c
   b.a
   c

ID: a8520eae-39da-419f-b268-88edd5d9cbfe
Created: 2023-08-11 15:35:31 +0000 UTC

   a
   b.a
   b.c
   c
   doe-iets-leuks

ID: 5a31c5a0-e41a-4cea-ae5c-d2a873fbcaf7
Created: 2023-08-11 14:30:37 +0000 UTC

   a
   b.a
   b.c
   doe-iets-leuks
```

## Reinstate Version

```
NAME:
   reinstate-credhub-secrets-version - Reinstate a previous version of the credhub service instance credentials

USAGE:
   
cf reinstate-credhub-secrets-version <SERVICE_INSTANCE> <VERSION_ID>

  SERVICE_INSTANCE - Credhub service instance name.
  VERSION_ID       - The credentials version id to reinstate. Can be obtained from the list-credhub-secrets-versions command.


ALIAS:
   rcv
```

### example

***NOTE*** in the example an ID from the listed version in the previous section is used.

```
$ cf rcv test-broker-api 5942167f-8d02-4b03-a1b8-1116c41069f3
$ cf lcv test-broker-api

output:

ID: c121044f-9d2a-4333-af80-a332c571fd55
Created: 2023-08-11 15:49:16 +0000 UTC

   c
   d.a
   doe-iets-leuks
   f
   a
   b.a
   b.b
   b.c

ID: 03142a68-254b-4739-a7a7-abdae54ec40d
Created: 2023-08-11 15:41:33 +0000 UTC

   b.c
   b.a
   c
   doe-iets-leuks
   f

ID: 5942167f-8d02-4b03-a1b8-1116c41069f3
Created: 2023-08-11 15:37:08 +0000 UTC

   f
   a
   b.a
   b.b
   b.c
   c
   d.a
   doe-iets-leuks

ID: 836bfc7b-1d75-4a0d-98c8-1772ff318709
Created: 2023-08-11 15:35:39 +0000 UTC

   a
   b.a
   b.c
   c
   d.a
   doe-iets-leuks

ID: a8520eae-39da-419f-b268-88edd5d9cbfe
Created: 2023-08-11 15:35:31 +0000 UTC

   a
   b.a
   b.c
   c
   doe-iets-leuks

ID: 5a31c5a0-e41a-4cea-ae5c-d2a873fbcaf7
Created: 2023-08-11 14:30:37 +0000 UTC

   a
   b.a
   b.c
   doe-iets-leuks
```

***NOTE*** The value of the reinstated version is reset, but effectively a new version is created, with a new ID.