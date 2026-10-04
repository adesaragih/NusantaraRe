package models

// IDPEGA di T_GENERAL_POLIS adalah kolom DASAR nbfacin 182 (VARCHAR2(50),
// keputusan WO 04-10-2026), bukan VARCHAR2(128) rancangan lama. Pemuat
// dokumen lama menulis `pyWorkPage.pzInsKey` apa adanya (`<kelas> <pyID>`);
// kelas Treaty `ASM-FW-GISFW-WORK-NB` 20 karakter + spasi + pyID sampai 32
// = 53 > 50. Dokumen yang IDPEGA-nya tidak muat ditolak di pemecah (galat
// dokumen), bukan ORA-12899 di tengah transaksi. Pelebaran diminta ke
// pemilik nbfacin (PERMINTAAN-TIM-INTI C10).

import (
	"errors"
	"strings"
	"testing"
)

func TestIDPegaMelebihiKolomDasarDitolak(t *testing.T) {
	if PanjangIDPega != 50 {
		t.Fatalf("PanjangIDPega = %d, harap 50 (T_GENERAL_POLIS.IDPEGA nbfacin 182)", PanjangIDPega)
	}
	kelas := strings.ToUpper(KelasDeret) + " "
	muat := kelas + "NB-" + strings.Repeat("9", 50-len(kelas)-3) // tepat 50
	if id, err := IDKasusDariIDPega(muat); err != nil || id != strings.TrimPrefix(muat, kelas) {
		t.Fatalf("IDPEGA 50 karakter: %q %v", id, err)
	}
	lebih := muat + "9" // 51; pyID tetap <= 32
	if len(strings.TrimPrefix(lebih, kelas)) > 32 {
		t.Fatal("persiapan: pyID harus tetap <= 32 supaya yang diuji batas IDPEGA")
	}
	if _, err := IDKasusDariIDPega(lebih); !errors.Is(err, ErrIDPega) || !strings.Contains(err.Error(), "50") {
		t.Fatalf("IDPEGA 51 karakter harus ditolak menyebut batas 50: %v", err)
	}
}
