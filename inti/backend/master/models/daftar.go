package models

// Daftar delapan master (rencana `modul/masterprovince/docs/issues/01-rencana-master-data.md`, cakupan "Yang dipakai
// NB FacIn dulu"; sejak 04-10-2026 tiap master satu modul menu - `03-pemecahan-delapan-modul.md`).
//
// Kolom dan lebar `[terverifikasi]` dari DDL `D:\migrasi\RNM\DDL\<NAMA>.txt`: NATION (TABEL), OBJECTITEMTYPE (TABEL),
// dan tabel flat migrasi 880 (PROVINCE, ACCUMULATEDTYPE, CZONE, ACCUMULATION, CITYINPUT, DISTRICTINPUT - kolom view
// asal, VARCHAR2(4000) A180). Rujukan turunan = definisi view asal: PROVINCE.NATIONNAME = Note M_Nation ber-ID
// NationID; CZONE.GROUPOFNAME = Description CZONE ber-Code GroupOf; ACCUMULATION.ACCUMULATIONNAME = AccumulationType
// ACCUMULATEDTYPE ber-ID Accumulation. ACCUMULATION.PROVINCE dari PROVINCE.NOTE ber-PROVINCEID `[dugaan]` (view asal
// membaca JSON ProvinceName apa adanya - MD-3).

const lebarFlat = 4000

func kf(nama, json string) Kolom { return Kolom{Nama: nama, JSON: json, Lebar: lebarFlat} }

// DaftarMaster - urutan tampil.
var DaftarMaster = []TabelMaster{
	{
		// NATION warisan tanpa kolom status: statusnya di T_MASTER_STATUS (migrasi 881, MD-2).
		Kunci: "nation", Judul: "Nation", Nama: "NATION", JejakTerpisah: true,
		Kolom: []Kolom{
			{Nama: "ID", JSON: "id", Lebar: 10, Wajib: true}, {Nama: "OLDID", JSON: "oldId", Lebar: 6},
			{Nama: "NOTE", JSON: "note", Lebar: 100, Wajib: true}, {Nama: "NATIONINITIAL", JSON: "nationInitial", Lebar: 20},
		},
		KolomCari: []string{"ID", "NOTE"}, KolomNama: "NOTE",
	},
	{
		Kunci: "province", Judul: "Province", Nama: "PROVINCE", KolomStatus: "STS_AKTIF",
		Kolom: []Kolom{
			{Nama: "ID", JSON: "id", Lebar: lebarFlat, Wajib: true}, kf("NATIONID", "nationId"),
			{Nama: "NOTE", JSON: "note", Lebar: lebarFlat, Wajib: true}, {Nama: "NATIONNAME", JSON: "nationName", Lebar: lebarFlat, Turunan: true},
		},
		KolomCari: []string{"ID", "NOTE"},
		KolomNama: "NOTE",
		Rujukan:   []Rujukan{{Sumber: "NATIONID", Tabel: "NATION", KolomKunci: "ID", KolomNilai: "NOTE", Turunan: "NATIONNAME", Master: "nation"}},
	},
	{
		// Keputusan work owner 04-10-2026 (diteruskan sesi 9d): "di menu city kolom email dan mo id di hapus saja" -
		// EMAIL / MOID tidak tampil dan tidak diisi menu, kolom DB-nya TETAP (tambah: NULL; ubah: tidak disentuh, karena
		// SQL hanya menyebut kolom di daftar ini). PROVINCENAME (migrasi 883, isi awal dari RW) kolom biasa, bukan turunan.
		Kunci: "city", Judul: "City", Nama: "CITYINPUT", KolomStatus: "STS_AKTIF",
		Kolom: []Kolom{
			{Nama: "ID", JSON: "id", Lebar: lebarFlat, Wajib: true}, kf("PROVINCEID", "provinceId"), kf("PROVINCENAME", "provinceName"),
			{Nama: "NOTE", JSON: "note", Lebar: lebarFlat, Wajib: true}, kf("BRANCHID", "branchId"),
			kf("JABODETABEKSTATUS", "jabodetabekStatus"),
		},
		KolomCari: []string{"ID", "NOTE"},
		KolomNama: "NOTE",
		Rujukan: []Rujukan{{Sumber: "PROVINCEID", Tabel: "PROVINCE", KolomKunci: "ID", Master: "province"},
			{Sumber: "BRANCHID", Tabel: "BRANCH", KolomKunci: "ID"}},
	},
	{
		Kunci: "district", Judul: "District", Nama: "DISTRICTINPUT", KolomStatus: "STS_AKTIF",
		Kolom: []Kolom{
			{Nama: "ID", JSON: "id", Lebar: lebarFlat, Wajib: true}, kf("CITYID", "cityId"),
			{Nama: "DISTRICTNAME", JSON: "districtName", Lebar: lebarFlat, Wajib: true},
		},
		KolomCari: []string{"ID", "DISTRICTNAME"}, KolomNama: "DISTRICTNAME",
		Rujukan: []Rujukan{{Sumber: "CITYID", Tabel: "CITYINPUT", KolomKunci: "ID", Master: "city"}},
	},
	{
		Kunci: "czone", Judul: "CZone", Nama: "CZONE", KolomStatus: "STS_AKTIF",
		Kolom: []Kolom{
			{Nama: "ID", JSON: "id", Lebar: lebarFlat, Wajib: true}, {Nama: "CODE", JSON: "code", Lebar: lebarFlat, Wajib: true},
			kf("DESCRIPTION", "description"), kf("GROUPOF", "groupOf"),
			{Nama: "GROUPOFNAME", JSON: "groupOfName", Lebar: lebarFlat, Turunan: true},
		},
		KolomCari: []string{"ID", "CODE", "DESCRIPTION"}, KolomNama: "DESCRIPTION",
		Rujukan: []Rujukan{{Sumber: "GROUPOF", Tabel: "CZONE", KolomKunci: "CODE", KolomNilai: "DESCRIPTION", Turunan: "GROUPOFNAME",
			Master: "czone"}},
	},
	{
		Kunci: "accumulatedtype", Judul: "Accumulated Type", Nama: "ACCUMULATEDTYPE", KolomStatus: "STS_AKTIF",
		Kolom: []Kolom{
			{Nama: "ID", JSON: "id", Lebar: lebarFlat, Wajib: true},
			{Nama: "ACCUMULATIONTYPE", JSON: "accumulationType", Lebar: lebarFlat, Wajib: true},
			kf("KEYWORD", "keyword"), kf("NOTE", "note"), kf("TYPE", "type"),
		},
		KolomCari: []string{"ID", "ACCUMULATIONTYPE"}, KolomNama: "ACCUMULATIONTYPE",
	},
	{
		Kunci: "accumulation", Judul: "Accumulation", Nama: "ACCUMULATION", KolomStatus: "STS_AKTIF", IDOtomatis: true,
		Kolom: []Kolom{
			{Nama: "ID", JSON: "id", Lebar: lebarFlat}, kf("ACCUMULATION", "accumulation"),
			{Nama: "ACCUMULATIONNAME", JSON: "accumulationName", Lebar: lebarFlat, Turunan: true},
			{Nama: "NOTE", JSON: "note", Lebar: lebarFlat, Wajib: true, HurufBesar: true}, kf("KEYWORD", "keyword"),
			kf("SCOPEAREA", "scopeArea"), kf("CZONE", "cZone"), kf("CZONEID", "cZoneId"),
			{Nama: "PROVINCE", JSON: "province", Lebar: lebarFlat, Turunan: true}, kf("PROVINCEID", "provinceId"),
			{Nama: "ZIPCODE", JSON: "zipCode", Lebar: lebarFlat, Wajib: true}, kf("ACCUMULATIONTYPE", "accumulationType"),
			kf("SYARIAHSTATUS", "syariahStatus"),
		},
		KolomCari: []string{"ID", "NOTE"}, KolomNama: "NOTE",
		Rujukan: []Rujukan{
			{Sumber: "ACCUMULATION", Tabel: "ACCUMULATEDTYPE", KolomKunci: "ID", KolomNilai: "ACCUMULATIONTYPE", Turunan: "ACCUMULATIONNAME",
				Master: "accumulatedtype"},
			{Sumber: "PROVINCEID", Tabel: "PROVINCE", KolomKunci: "ID", KolomNilai: "NOTE", Turunan: "PROVINCE", Master: "province"},
		},
	},
	{
		Kunci: "objectitemtype", Judul: "Object Item Type", Nama: "OBJECTITEMTYPE", KolomStatus: "ISACTIVE", JejakTerpisah: true,
		Kolom: []Kolom{
			{Nama: "ID", JSON: "id", Lebar: lebarFlat, Wajib: true},
			{Nama: "OBJECTITEMTYPE", JSON: "objectItemType", Lebar: lebarFlat, Wajib: true}, kf("NOTE", "note"),
			kf("PCTADJUSTABLE1", "pctAdjustable1"), kf("PCTADJUSTABLE2", "pctAdjustable2"),
			kf("OBJECTITEMTYPEINA", "objectItemTypeIna"), kf(`"GROUP"`, "group"), kf("TYPE", "type"),
		},
		KolomCari: []string{"ID", "OBJECTITEMTYPE"}, KolomNama: "OBJECTITEMTYPE",
	},
}

// CariMaster - master ber-kunci rute `kunci`.
func CariMaster(kunci string) (TabelMaster, bool) {
	for _, t := range DaftarMaster {
		if t.Kunci == kunci {
			return t, true
		}
	}
	return TabelMaster{}, false
}

// KolomAudit - jejak ubah setiap master (MD-7, keputusan work owner 04-10-2026 "Kolom diubah oleh/tanggal"; migrasi
// 882): pembuat, tanggal buat, pengubah terakhir, tanggal ubah terakhir - TURUNAN (diisi backend dari pelaku, baca-saja
// di layar). Nama kolom pola T_WORK_* (MD-10). Tabel JejakTerpisah: kolom yang sama di T_MASTER_STATUS.
var KolomAudit = []Kolom{
	{Nama: "CREATE_OP", JSON: "createOp", Lebar: 64, Turunan: true},
	{Nama: "TGL_CREATE", JSON: "tglCreate", Lebar: 19, Turunan: true, Tanggal: true},
	{Nama: "UPDATE_OP", JSON: "updateOp", Lebar: 64, Turunan: true},
	{Nama: "TGL_UPDATE", JSON: "tglUpdate", Lebar: 19, Turunan: true, Tanggal: true},
}
