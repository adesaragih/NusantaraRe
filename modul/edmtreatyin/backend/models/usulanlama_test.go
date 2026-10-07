package models

// Uji salinan SuggestList dokumen lama endorsemen ke riwayat produksi (pemetaan `petaUsulan` jalur biasa).

import "testing"

func TestUsulanDokumenLamaEDM(t *testing.T) {
	h, err := PecahDokumenEDM(barisUjiEDM(dokumenUjiEDMProp))
	if err != nil || len(h.Galat) > 0 {
		t.Fatalf("%v %+v", err, h.Galat)
	}
	if len(h.Usulan) != 2 {
		t.Fatalf("baris IsSave kosong saja yang disalin, dapat %+v", h.Usulan)
	}
	u := h.Usulan[0]
	if u.TypePolis != "EDMT" || u.TglInp != "2017-10-02 08:00:00" || u.PIC != "UJI-OPERATOR" || u.Approval != "Accept" ||
		u.Keterangan != "UJI catatan pertama" || u.Type != "T" || u.AksesLogin != "" || u.Posisi != PosisiUsulanProduksi {
		t.Errorf("baris pertama %+v", u)
	}
	if v := h.Usulan[1]; v.PIC != "" || v.Approval != "Reject" || v.TglInp != "2017-10-02 09:00:00" {
		t.Errorf("baris kedua %+v", v)
	}
	if nama, ok := pmAnggotaUsulanLama(DaftarUsulan + "().Suggest"); !ok || nama != "Suggest" {
		t.Error("anggota SuggestList tidak dikenali")
	}
	if _, ok := pmAnggotaUsulanLama(DaftarUsulan + "().UJILain"); ok {
		t.Error("anggota tak dikenal diterima")
	}
}
