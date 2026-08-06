//go:build (!amd64 && !arm64) || go1.27 || !go1.17 || (arm64 && !go1.20)
// +build !amd64,!arm64 go1.27 !go1.17 arm64,!go1.20

/*
 * Copyright 2022 ByteDance Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package ast

import (
	"github.com/bytedance/sonic/internal/native/types"
	"github.com/bytedance/sonic/unquote"
)

func (self *Parser) searchKey(match string) types.ParsingError {
	ns := len(self.s)
	if err := self.object(); err != 0 {
		return err
	}

	/* check for EOF */
	if self.p = self.lspace(self.p); self.p >= ns {
		return types.ERR_EOF
	}

	/* check for empty object */
	if self.s[self.p] == '}' {
		self.p++
		return _ERR_NOT_FOUND
	}

	var njs types.JsonState
	var err types.ParsingError
	/* decode each pair */
	for {

		/* decode the key */
		if njs = self.decodeValue(); njs.Vt != types.V_STRING {
			return types.ERR_INVALID_CHAR
		}

		/* extract the key */
		idx := self.p - 1
		key := self.s[njs.Iv:idx]

		/* check for escape sequence */
		if njs.Ep != -1 {
			if key, err = unquote.String(key); err != 0 {
				return err
			}
		}

		/* expect a ':' delimiter */
		if err = self.delim(); err != 0 {
			return err
		}

		/* skip value */
		if key != match {
			if _, err = self.skipFast(); err != 0 {
				return err
			}
		} else {
			return 0
		}

		/* check for EOF */
		self.p = self.lspace(self.p)
		if self.p >= ns {
			return types.ERR_EOF
		}

		/* check for the next character */
		switch self.s[self.p] {
		case ',':
			self.p++
		case '}':
			self.p++
			return _ERR_NOT_FOUND
		default:
			return types.ERR_INVALID_CHAR
		}
	}
}

func (self *Parser) searchIndex(idx int) types.ParsingError {
	ns := len(self.s)
	if err := self.array(); err != 0 {
		return err
	}

	/* check for EOF */
	if self.p = self.lspace(self.p); self.p >= ns {
		return types.ERR_EOF
	}

	/* check for empty array */
	if self.s[self.p] == ']' {
		self.p++
		return _ERR_NOT_FOUND
	}

	var err types.ParsingError
	/* allocate array space and parse every element */
	for i := 0; i < idx; i++ {

		/* decode the value */
		if _, err = self.skipFast(); err != 0 {
			return err
		}

		/* check for EOF */
		self.p = self.lspace(self.p)
		if self.p >= ns {
			return types.ERR_EOF
		}

		/* check for the next character */
		switch self.s[self.p] {
		case ',':
			self.p++
		case ']':
			self.p++
			return _ERR_NOT_FOUND
		default:
			return types.ERR_INVALID_CHAR
		}
	}

	return 0
}
