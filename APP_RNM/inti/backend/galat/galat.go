package galat

import (
	"errors"
)

// ErrPermintaanTidakSah menandai permintaan pendaftaran yang tidak lengkap.
var ErrPermintaanTidakSah = errors.New("services: permintaan pendaftaran tidak sah")
