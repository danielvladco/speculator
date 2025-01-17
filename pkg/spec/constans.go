// Copyright © 2021 Cisco Systems, Inc. and its affiliates.
// All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package spec

const (
	parametersInBody   = "body"
	parametersInHeader = "header"
	parametersInQuery  = "query"
	parametersInForm   = "formData"
	parametersInPath   = "path"
)

const (
	formatUUID = "uuid"
)

const (
	schemaTypeObject  = "object"
	schemaTypeArray   = "array"
	schemaTypeBoolean = "boolean"
	schemaTypeInteger = "integer"
	schemaTypeNumber  = "number"
	schemaTypeString  = "string"
)

const inBodyParameterName = "body"

const (
	collectionFormatComma = "csv"
	collectionFormatSpace = "ssv"
	collectionFormatTab   = "tsv"
	collectionFormatPipe  = "pipes"
	collectionFormatMulti = "multi"
)

const (
	contentTypeHeaderName       = "content-type"
	acceptTypeHeaderName        = "accept"
	authorizationTypeHeaderName = "authorization"
)

const (
	mediaTypeApplicationJSON    = "application/json"
	mediaTypeApplicationHalJSON = "application/hal+json"
	mediaTypeApplicationForm    = "application/x-www-form-urlencoded"
	mediaTypeMultipartFormData  = "multipart/form-data"
)
