// Kontrol untuk TestJabatanAntreanGagalKompilasi: pemakaian yang benar WAJIB
// terkompilasi, supaya kegagalan pasangannya terbukti karena tipe, bukan karena
// berkas yang rusak.
package main

import "nusantarare/modul/nbfacin/backend/services/acceptance"

func main() {
	t := acceptance.Transisi{JabatanTujuan: acceptance.Jabatan("KADIVTEKNIK"), Antrean: acceptance.Antrean("ReasFacInGroupLeader")}
	_ = t
}
