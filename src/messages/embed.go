/*
	embed.go
	Purpose: Embed message packs as fs.FS.

	@author Evan Chen
	@version 1.0 2023/02/22
*/

package messages

import "embed"

//go:embed *.json
var FS embed.FS
