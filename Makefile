define PARSER_HEADER
// Copyright 2014 Frustra Sofware. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

// This file is generated from parser.y

endef
export PARSER_HEADER

yacc:
	go tool yacc -o parser.y.go parser.y
	@echo "$$PARSER_HEADER" | cat - parser.y.go > parser.go
	@rm parser.y.go
	go fmt parser.go
