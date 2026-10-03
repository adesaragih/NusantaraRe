package tiruan

import "nusantarare/modul/nbtreatyin/backend/services"

// Tiruan harus tetap memenuhi antarmuka layanan.
var _ services.Gudang = (*Gudang)(nil)
