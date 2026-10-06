package models

// Untuk apa berkas ini: TEMPAT BERPERAN - tiket 05 (spec §5.4; AC 12, 13, 81,
// 82, 91). Di Pega beberapa bagian layar muncul atau wajib hanya untuk orang
// tertentu (`OperatorID.pyUserIdentifier == '<ID-operator-N>'`). Pemetaan
// nama -> peran TIDAK ada di korpus; sistem baru memetakan TEMPAT -> PERAN ->
// ARAH, dan peran pengguna dibaca dari `inti.Pelaku.Peran` (workbasket akun).
//
// ⛔ KONSTANTA KODE, BUKAN TABEL. `[keputusan work owner]` K16 (03-10-2026,
// PROMPT-NB-TREATY-IN-PUTARAN-2 bab 2): tabel `M_NBTRIN_PERAN_TEMPAT` (migrasi
// 330 putaran 1) tidak ada di diagram grilling dan DIHAPUS; pemetaannya hidup
// di `PemetaanPeranTempat` di bawah - KOSONG sampai IAM menjawab (K12). Tiket
// 05 tetap needs-info. Perilaku saat kosong = perilaku tabel kosong putaran 1:
//
//	tempat tanpa baris  -> TERTUNDA: bagian itu tidak tampil, syaratnya tidak
//	                       dianggap terpenuhi (AC 81) - tidak ditebak
//	ARAH MUNCUL         -> tampil hanya bagi pemegang peran itu
//	ARAH KECUALI        -> tampil bagi semua KECUALI pemegang peran itu
//	kedua arah / arah asing di satu tempat -> tetap tertunda (AC 82)
//
// ⛔ Nol peran karangan (AC 91): uji memakai peran fiktif berawalan UJI-.
//
// KEDUA BELAS TEMPAT (`[terverifikasi]` grilling ronde 2 P28 / PERTANYAAN-untuk-
// IAM P28: "4 layar · 12 tempat" = `pyUserIdentifier` di
// `DetailDeptHeadTreatyIn_UW` 3× · `GeneralDeptHeadTreatyIn_UW` 3× · `ListSuggest`
// 4×; `pxInsName` di `DetailPoliciesNonProportional` 2×), dibaca ulang dari XML
// 2026-10-04 - lihat `DaftarTempat`. Identitas orang disamarkan
// `<ID-operator-N>` seperti INVENTARIS-XML bab 5 (nilainya tidak disalin ke
// artefak mana pun). Arah setiap tempat TIDAK ditebak dari bentuk `==`/`!=`
// syarat XML: itu jawaban IAM (AC 82).

// PeranTempat adalah satu baris pemetaan tempat -> peran -> arah.
type PeranTempat struct {
	KodeTempat string
	Peran      string
	// Arah - ArahMuncul atau ArahKecuali; TIDAK ditebak (AC 82).
	Arah string
}

// Kode kedua belas tempat (tiket 05).
const (
	// `Section/DetailDeptHeadTreatyIn_UW` - tiga tombol `.pyTemplateButton` "Submit".
	TempatDetailDHSubmitLetterNoDeptHead = "DETAILDEPTHEADTREATYIN_UW_SUBMIT_LETTERNO_TREATYINDEPTHEAD"
	TempatDetailDHSubmitLetterNoKosong   = "DETAILDEPTHEADTREATYIN_UW_SUBMIT_LETTERNO_KOSONG"
	TempatDetailDHSubmitOperator1        = "DETAILDEPTHEADTREATYIN_UW_SUBMIT_OPERATOR_1"
	// `Section/GeneralDeptHeadTreatyIn_UW` - salinan sertakan ketiga tombol yang
	// sama (berkas XML-nya memuat aliran `DetailDeptHeadTreatyIn_UW`).
	TempatGeneralDHSubmitLetterNoDeptHead = "GENERALDEPTHEADTREATYIN_UW_SUBMIT_LETTERNO_TREATYINDEPTHEAD"
	TempatGeneralDHSubmitLetterNoKosong   = "GENERALDEPTHEADTREATYIN_UW_SUBMIT_LETTERNO_KOSONG"
	TempatGeneralDHSubmitOperator1        = "GENERALDEPTHEADTREATYIN_UW_SUBMIT_OPERATOR_1"
	// `Section/ListSuggest` `.ProductionDate` - syarat tampil (`pyVisible OTHER`)
	// dan syarat wajib (`pyRequiredWhen`), masing-masing dua identitas.
	TempatProduksiTampilOperator3 = "LISTSUGGEST_PRODUCTIONDATE_TAMPIL_OPERATOR_3"
	TempatProduksiTampilOperator4 = "LISTSUGGEST_PRODUCTIONDATE_TAMPIL_OPERATOR_4"
	TempatProduksiWajibOperator3  = "LISTSUGGEST_PRODUCTIONDATE_WAJIB_OPERATOR_3"
	TempatProduksiWajibOperator4  = "LISTSUGGEST_PRODUCTIONDATE_WAJIB_OPERATOR_4"
	// `Section/DetailPoliciesNonProportional` - dua LABEL (`OperatorID.pxInsName`).
	TempatLabelNonEDM = "DETAILPOLICIESNONPROPORTIONAL_LABEL_NON_EDM"
	TempatLabelEDM    = "DETAILPOLICIESNONPROPORTIONAL_LABEL_EDM"
)

// TempatBerperan - satu tempat identitas di rule terjangkau.
type TempatBerperan struct {
	Kode string
	// Section - berkas `Section\<nama>.xml`.
	Section string
	// Sel - sel dan atribut yang memuat syarat identitas.
	Sel string
	// Syarat - bentuk syarat XML (identitas disamarkan).
	Syarat string
	// Gerbang - apa yang digerbang tempat ini di sistem baru.
	Gerbang string
}

const (
	syaratSubmitDeptHead = `.IsApproved == 1 && OperatorID.pyUserIdentifier!='<ID-operator-1>' && pyWorkPage.LetterNo=='TREATYINDEPTHEAD'`
	syaratSubmitKosong   = `.IsApproved == 1 && OperatorID.pyUserIdentifier!='<ID-operator-1>' && pyWorkPage.LetterNo==''`
	syaratSubmitOperator = `.IsApproved == 1 && OperatorID.pyUserIdentifier=='<ID-operator-1>'`
	syaratProduksi       = `.IsApproved == 1 && (OperatorID.pyUserIdentifier=='<ID-operator-3>' || OperatorID.pyUserIdentifier=='<ID-operator-4>')`
	syaratLabel          = `OperatorID.pxInsName = '<ID-operator-2>'`

	// gerbangTombolAtasan - ketiga tombol Submit atasan DIGANTI posisi kasus
	// (`TombolUntuk`): Dept Head - jenjang terakhir tangga P13 - satu-satunya
	// yang menerbitkan nomor polis, Sec Head selalu menaikkan (AC 8, K2).
	// Tempat ini terdaftar tetapi TIDAK dibaca layanan (catatan tiket 05).
	gerbangTombolAtasan = "tidak dibaca - tombol ditentukan posisi kasus (TombolUntuk; AC 8, K2, P13)"
)

// DaftarTempat - kedua belas tempat tiket 05, urutan berkas XML.
var DaftarTempat = []TempatBerperan{
	{TempatDetailDHSubmitLetterNoDeptHead, "DetailDeptHeadTreatyIn_UW", "pxButton Submit (finishAssignment) pyVisible", syaratSubmitDeptHead, gerbangTombolAtasan},
	{TempatDetailDHSubmitLetterNoKosong, "DetailDeptHeadTreatyIn_UW", "pxButton Submit (GeneratePolicyNoTreaty_Act) pyVisible", syaratSubmitKosong, gerbangTombolAtasan},
	{TempatDetailDHSubmitOperator1, "DetailDeptHeadTreatyIn_UW", "pxButton Submit (GeneratePolicyNoTreaty_Act) pyVisible", syaratSubmitOperator, gerbangTombolAtasan},
	{TempatGeneralDHSubmitLetterNoDeptHead, "GeneralDeptHeadTreatyIn_UW", "pxButton Submit (finishAssignment) pyVisible", syaratSubmitDeptHead, gerbangTombolAtasan},
	{TempatGeneralDHSubmitLetterNoKosong, "GeneralDeptHeadTreatyIn_UW", "pxButton Submit (GeneratePolicyNoTreaty_Act) pyVisible", syaratSubmitKosong, gerbangTombolAtasan},
	{TempatGeneralDHSubmitOperator1, "GeneralDeptHeadTreatyIn_UW", "pxButton Submit (GeneratePolicyNoTreaty_Act) pyVisible", syaratSubmitOperator, gerbangTombolAtasan},
	{TempatProduksiTampilOperator3, "ListSuggest", ".ProductionDate pyVisible (identitas ke-1 dari 2)", syaratProduksi, "tampil + diterima dari layar (TanggalProduksiTampil)"},
	{TempatProduksiTampilOperator4, "ListSuggest", ".ProductionDate pyVisible (identitas ke-2 dari 2)", syaratProduksi, "tampil + diterima dari layar (TanggalProduksiTampil)"},
	{TempatProduksiWajibOperator3, "ListSuggest", ".ProductionDate pyRequiredWhen (identitas ke-1 dari 2)", syaratProduksi, "wajib (TanggalProduksiWajib)"},
	{TempatProduksiWajibOperator4, "ListSuggest", ".ProductionDate pyRequiredWhen (identitas ke-2 dari 2)", syaratProduksi, "wajib (TanggalProduksiWajib)"},
	{TempatLabelNonEDM, "DetailPoliciesNonProportional", `LABEL "NON EDM" pyVisible`, syaratLabel, "label layar NonProp (frontend/tempat.ts)"},
	{TempatLabelEDM, "DetailPoliciesNonProportional", `LABEL "EDM" pyVisible`, syaratLabel, "label layar NonProp (frontend/tempat.ts)"},
}

// SemuaTempat - kode kedua belas tempat (`DaftarTempat`).
var SemuaTempat = func() []string {
	out := make([]string, len(DaftarTempat))
	for i, t := range DaftarTempat {
		out[i] = t.Kode
	}
	return out
}()

// Arah pemeriksaan.
const (
	ArahMuncul  = "MUNCUL"
	ArahKecuali = "KECUALI"
)

// PemetaanPeranTempat - pemetaan yang berlaku. ⛔ KOSONG sampai IAM bersama
// work owner menyerahkan peran dan arah ke-12 tempat (tiket 05, K12). Mengisi
// baris di sini adalah keputusan wewenang - bukan pekerjaan agen.
var PemetaanPeranTempat = []PeranTempat{}

// TempatTampil menghitung tampil/tidaknya setiap tempat `SemuaTempat` bagi
// pemegang peran `punyaPeran` (lazimnya `inti.Pelaku.PunyaPeran`).
func TempatTampil(pemetaan []PeranTempat, punyaPeran func(string) bool) map[string]bool {
	out := map[string]bool{}
	for _, t := range SemuaTempat {
		out[t] = false
	}
	type keadaan struct{ muncul, kecuali, asing, punyaMuncul, punyaKecuali bool }
	per := map[string]*keadaan{}
	for _, b := range pemetaan {
		k := per[b.KodeTempat]
		if k == nil {
			k = &keadaan{}
			per[b.KodeTempat] = k
		}
		switch b.Arah {
		case ArahMuncul:
			k.muncul = true
			k.punyaMuncul = k.punyaMuncul || punyaPeran(b.Peran)
		case ArahKecuali:
			k.kecuali = true
			k.punyaKecuali = k.punyaKecuali || punyaPeran(b.Peran)
		default:
			k.asing = true
		}
	}
	for kode, k := range per {
		if _, dikenal := out[kode]; !dikenal || k.asing {
			continue
		}
		switch {
		case k.muncul && !k.kecuali:
			out[kode] = k.punyaMuncul
		case k.kecuali && !k.muncul:
			out[kode] = !k.punyaKecuali
		}
		// kedua arah di satu tempat = bertentangan: tetap tertunda.
	}
	return out
}

// TanggalProduksiTampil = syarat tampil `Section/ListSuggest` `.ProductionDate`:
// `.IsApproved == 1 && (<identitas-3> || <identitas-4>)` - bagian identitas
// diganti dua tempat berperannya. Sel tak tampil tidak terkirim di Pega, maka
// nilainya hanya diterima bila tampil (`GabungMasukanLayar`).
func TanggalProduksiTampil(h *Halaman, tempat map[string]bool) bool {
	return h.Ambil(HalamanPolis+".IsApproved") == "1" &&
		(tempat[TempatProduksiTampilOperator3] || tempat[TempatProduksiTampilOperator4])
}

// TanggalProduksiWajib = `pyRequiredWhen` sel yang sama (syarat berbunyi sama,
// tempat sendiri). Sel yang tidak ter-render tidak menegakkan wajibnya, maka
// juga harus tampil.
func TanggalProduksiWajib(h *Halaman, tempat map[string]bool) bool {
	return TanggalProduksiTampil(h, tempat) &&
		(tempat[TempatProduksiWajibOperator3] || tempat[TempatProduksiWajibOperator4])
}
