package services

// Proteksi objek Aneka (tiket 73 P1): port `NB FacIn\Activity\ProtectionObjectAneka_Act.xml` - post-processing flow action
// `InputDtlObjectAneka_FacIn` (FromWhere "0", pesan tampil) atas SATU baris AnekaList; di sini dijalankan untuk setiap objek
// saat Save tab. Dispatcher 12 langkah `Call ProtectionObject<Lini>_Act` (gerbang per langkah), tiap sub-activity berpasangan
// `Property-Set local.cekX = "<teks> \n"` + `Property-Set-Messages Field=<prop>` bila properti kosong (`@String.equals(.X,"")`).
// Bukti per baris di riset form tiket 73 bagian A (subagen sesi c3; diperiksa ulang: HE L344, GrowingTree L1103 / L1253 /
// L1310, ElectronicEq L1002 / L2056, Burglary L811, dispatcher L372 / L495 / L1602; 63 literal pesan = grep dua cara).
//
// K73-4 (tiru apa adanya): 10 dari 12 gerbang membandingkan BusinessCode dengan kode 2 karakter, padahal BusinessCode FacIn
// 5 digit - gerbang itu tidak pernah terbuka; yang hidup praktis hanya HE (IsHe) dan Electronic Equipment
// (isElectronicEquipment, pembanding bertepi spasi " ElectronicEquipment " - belum terverifikasi apakah dipangkas).
// Langkah ber-label `//` (HE Quantity / UnitQuantity, GrowingTree TSI) TIDAK diport. Anomali dibawa apa adanya: Golf / Fidelity
// menempel pesan Section ke `.Quantity`, ElectronicEq pesan Quantity ke `.NoOfTree`, AllRisk menguji halaman kerja
// `pyWorkPage.Aneka.TerritorialLimits` (tidak ada di model - selalu kosong), Golf / Fidelity menguji `.LimitOfLiability` (kolom
// grid okupasi).

import (
	"fmt"
	"log"
	"strings"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/services/rules"
)

// PelanggaranAneka - satu pesan proteksi objek: posisi (0-based; Okupasi nil = objek Golf), medan kabel tempat pesan ditempel
// (Field Pega, apa adanya), pesan.
type PelanggaranAneka struct {
	Lokasi  int
	Okupasi *int
	Objek   int
	Medan   string
	Pesan   string
}

// GalatProteksiAneka - 422: objek tidak lolos ProtectionObjectAneka_Act.
type GalatProteksiAneka struct{ Pelanggaran []PelanggaranAneka }

func (g GalatProteksiAneka) Error() string {
	p := make([]string, len(g.Pelanggaran))
	for i, x := range g.Pelanggaran {
		p[i] = x.Pesan
	}
	return "services: proteksi objek Aneka: " + strings.Join(p, "; ")
}

// akhiranPesanPega - setiap literal pesan berakhir spasi + DUA karakter `\n` (backslash, n) di XML.
const akhiranPesanPega = ` \n`

// pesanProteksi - teks layar: akhiran ` \n` literal dipotong `[keputusan agent]` (diusulkan ke sesi 9d).
func pesanProteksi(verbatim string) string { return strings.TrimSuffix(verbatim, akhiranPesanPega) }

// cekAneka - satu pasangan set + pesan: `kosong` = prakondisi langkah pesan; `teks` = pesan verbatim XML.
type cekAneka struct {
	kosong func(models.ObjekAneka) bool
	medan  string
	teks   string
}

func teksKosong(f func(models.ObjekAneka) string) func(models.ObjekAneka) bool {
	return func(o models.ObjekAneka) bool { return f(o) == "" }
}

func desKosong(f func(models.ObjekAneka) *apd.Decimal) func(models.ObjekAneka) bool {
	return func(o models.ObjekAneka) bool { return f(o) == nil }
}

var (
	cekTahun      = cekAneka{teksKosong(func(o models.ObjekAneka) string { return o.Year }), "year", `Tahun harus diisi \n`}
	cekJumlah     = cekAneka{teksKosong(func(o models.ObjekAneka) string { return o.Quantity }), "quantity", `Jumlah harus diisi \n`}
	cekSatuan     = cekAneka{teksKosong(func(o models.ObjekAneka) string { return o.UnitQuantity }), "unitQuantity", `Satuan harus diisi \n`}
	cekNamaObyek  = cekAneka{teksKosong(func(o models.ObjekAneka) string { return o.ObjectName }), "objectName", `Nama Obyek harus diisi \n`}
	cekSerial     = cekAneka{teksKosong(func(o models.ObjekAneka) string { return o.SerialNo }), "serialNo", `Serial No harus diisi \n`}
	cekKetObjek   = cekAneka{teksKosong(func(o models.ObjekAneka) string { return o.ObjectName }), "objectName", `Keterangan Object harus di isi \n`}
	cekMesinLima  = []cekAneka{cekTahun, cekJumlah, cekSatuan, cekNamaObyek, cekSerial} // LandRig / ContractorPlant / Boiler / MB
	cekSeksiGolf  = cekAneka{teksKosong(func(o models.ObjekAneka) string { return o.Section }), "quantity", `Seksi / Section harus dipilih \n`}
	cekBatasKewaj = cekAneka{desKosong(func(o models.ObjekAneka) *apd.Decimal { return o.LimitOfLiability }), "limitOfLiability",
		`Batas kewajiban harus diisi \n`}
)

// kodeBisnisSama - `@String.equals(pyWorkPage.Quotation.BusinessCode, k1) || …` (perbandingan teks persis).
func kodeBisnisSama(kode ...string) func(string, kasusBisnis) bool {
	return func(_ string, kb kasusBisnis) bool {
		for _, k := range kode {
			if kb.businessCode == k {
				return true
			}
		}
		return false
	}
}

// predikatWhen - gerbang When lewat rules.Eval (galat = salah, dicatat; pola salahSatu).
func predikatWhen(nama string) func(string, kasusBisnis) bool {
	return func(id string, kb kasusBisnis) bool {
		v, err := rules.Eval(nama, kb)
		if err != nil {
			log.Printf("nbfacin: proteksi Aneka kasus %s: %s: %v", id, nama, err)
		}
		return err == nil && v
	}
}

// langkahProteksiAneka - 12 langkah dispatcher, urutan XML (L276 … L1640).
var langkahProteksiAneka = []struct {
	activity string
	gerbang  func(string, kasusBisnis) bool
	cek      []cekAneka
}{
	// 1 L372 `IsHe` (rule IsHE; resolusi nama tak peka huruf = dugaan). Langkah 3-6 (Quantity, UnitQuantity) ber-label `//`.
	{"ProtectionObjectHE_Act", predikatWhen("IsHE"), []cekAneka{cekTahun,
		{teksKosong(func(o models.ObjekAneka) string { return o.VehicleHE.ChassisNumber }), "vehicleHe.chassisNumber", `No Rangka harus diisi \n`},
		{teksKosong(func(o models.ObjekAneka) string { return o.VehicleHE.EngineNumber }), "vehicleHe.engineNumber", `No Mesin harus diisi \n`},
		{func(o models.ObjekAneka) bool { return o.VehicleHE.Type == "" && o.VehicleHE.TypeName == "" }, "vehicleHe.typeName", `Tipe harus diisi \n`},
		{func(o models.ObjekAneka) bool { return o.VehicleHE.ObjectIDHE == "" && o.VehicleHE.ObjectNameHE == "" }, "vehicleHe.objectNameHe",
			`Nama Obyek harus diisi \n`},
		{func(o models.ObjekAneka) bool { return o.VehicleHE.ObjectNameHE == "TBA" }, "vehicleHe.objectNameHe", `Nama Obyek tidak boleh TBA \n`},
		{func(o models.ObjekAneka) bool { return o.VehicleHE.Brand == "" && o.VehicleHE.BrandName == "" }, "vehicleHe.brandName",
			`Merk Obyek harus diisi \n`},
		{func(o models.ObjekAneka) bool { return o.VehicleHE.BrandName == "TBA" }, "vehicleHe.brandName", `Merk Obyek tidak boleh TBA \n`},
	}},
	{"ProtectionObjectLandRig_Act", kodeBisnisSama("B3"), cekMesinLima},               // 2 L495
	{"ProtectionObjectContractorPlant_Act", kodeBisnisSama("12", "S5"), cekMesinLima}, // 3 L618
	{"ProtectionObjectBoilerPresure_Act", kodeBisnisSama("11"), cekMesinLima},         // 4 L741
	{"ProtectionObjectMB_Act", kodeBisnisSama("10", "S2"), cekMesinLima},              // 5 L864
	{"ProtectionObjectAviation_Act", kodeBisnisSama("06"), []cekAneka{ // 6 L987
		{teksKosong(func(o models.ObjekAneka) string { return o.ObjectName }), "objectName", `Model dan Tipe harus di isi \n`},
		{teksKosong(func(o models.ObjekAneka) string { return o.Year }), "year", `Tahun Pembuatan harus diisi \n`},
		{teksKosong(func(o models.ObjekAneka) string { return o.GeographicalLimits }), "geographicalLimits", `Batas Geografis harus diisi \n`},
		{teksKosong(func(o models.ObjekAneka) string { return o.AviationHull.RegistrationMarks }), "aviationHull.registrationMarks",
			`Kode Registrasi harus diisi \n`},
		{teksKosong(func(o models.ObjekAneka) string { return o.AviationHull.MaxPassengers }), "aviationHull.maxPassengers",
			`Maksimal Penumpang harus diisi \n`},
		{teksKosong(func(o models.ObjekAneka) string { return o.AviationHull.Pilots }), "aviationHull.pilots", `Pilot harus diisi \n`},
		{teksKosong(func(o models.ObjekAneka) string { return o.AviationHull.StandardUses }), "aviationHull.standardUses",
			`Penggunaan Standar harus diisi \n`},
		{teksKosong(func(o models.ObjekAneka) string { return o.AviationHull.SpecialUses }), "aviationHull.specialUses",
			`Penggunaan Khusus harus diisi \n`},
		{teksKosong(func(o models.ObjekAneka) string { return o.AviationHull.SpecialRentalUses }), "aviationHull.specialRentalUses",
			`Penggunaan Khusus Rental harus diisi \n`},
	}},
	{"ProtectionObjectAllRisk_Act", kodeBisnisSama("16", "SR"), []cekAneka{cekKetObjek, // 7 L1110
		{teksKosong(func(o models.ObjekAneka) string { return o.Quantity }), "quantity", `Jumlah harus diisi \n`},
		// pyWorkPage.Aneka.TerritorialLimits - halaman kerja, tidak ada di model objek: selalu kosong.
		{func(models.ObjekAneka) bool { return true }, "aneka.territorialLimits", `Batas Teritorial harus diisi \n`},
	}},
	{"ProtectionObjectBurglary_Act", kodeBisnisSama("14", "SP"), []cekAneka{cekKetObjek, // 8 L1233
		{teksKosong(func(o models.ObjekAneka) string { return o.Quantity }), "quantity", `Jumlah harus diisi \n`},
		{teksKosong(func(o models.ObjekAneka) string { return o.InvoiceNo }), "invoiceNo", `Invoiceharus diisi \n`},
	}},
	{"ProtectionObjectGolfInsurance_Act", kodeBisnisSama("27", "SW"), []cekAneka{cekKetObjek, cekSeksiGolf, cekBatasKewaj}}, // 9 L1356
	{"ProtectionObjectGrowingTree_Act", kodeBisnisSama("28", "S6"), []cekAneka{cekKetObjek, // 10 L1479; TSI (9-10) ber-label `//`
		{teksKosong(func(o models.ObjekAneka) string { return o.Year }), "year", `Tahun tanam harus di input \n`},
		{desKosong(func(o models.ObjekAneka) *apd.Decimal { return o.NoOfTree }), "noOfTree", `Jumlah pohon harus di isi \n`},
		// Langkah 7 men-SET teks bila .NoOfTree kosong, langkah 8 menempel bila .AreaHectar kosong: teks hanya ada bila keduanya
		// kosong. AreaHectar kosong + NoOfTree terisi = pesan kosong di Pega (efek belum terverifikasi) - tidak dikirim.
		{func(o models.ObjekAneka) bool { return o.AreaHectar == nil && o.NoOfTree == nil }, "areaHectar", `area hectar harus di isi \n`},
	}},
	{"ProtectionObjectElectronicEq_Act", predikatWhen("isElectronicEquipment"), []cekAneka{cekKetObjek, // 11 L1602
		{teksKosong(func(o models.ObjekAneka) string { return o.Section }), "section", `Section harus di pilih \n`},
		{teksKosong(func(o models.ObjekAneka) string { return o.Quantity }), "noOfTree", `Jumlah harus di isi \n`},
		{desKosong(func(o models.ObjekAneka) *apd.Decimal { return o.TSI }), "tsi", `TSI Harus di input \n`},
		{teksKosong(func(o models.ObjekAneka) string { return o.Year }), "year", `Tahun Pembuatan di isi \n`},
		{desKosong(func(o models.ObjekAneka) *apd.Decimal { return o.LimitIdemnityDay }), "limitIdemnityDay",
			`Limit of idemnity Day harus di isi \n`},
		// Message `Local.cekIdemMonth` (L besar) vs set `local.cekIdemMonth` - dianggap variabel yang sama (dugaan).
		{desKosong(func(o models.ObjekAneka) *apd.Decimal { return o.LimitIdemnityMonth }), "limitIdemnityMonth",
			`Limit of idemnity Month harus di isi \n`},
	}},
	{"ProtectionObjectFidelity_Act", kodeBisnisSama("19", "SS"), []cekAneka{cekKetObjek, cekSeksiGolf, cekBatasKewaj}}, // 12 L1736
}

// proteksiObjekAneka - pelanggaran seluruh objek (okupasi lalu Golf per lokasi), urutan langkah dispatcher lalu urutan cek.
func proteksiObjekAneka(id string, kb kasusBisnis, lokasi []models.LokasiAneka) []PelanggaranAneka {
	var aktif []int
	for i, l := range langkahProteksiAneka {
		if l.gerbang(id, kb) {
			aktif = append(aktif, i)
		}
	}
	var hasil []PelanggaranAneka
	periksa := func(i int, okupasi *int, n int, o models.ObjekAneka) {
		for _, a := range aktif {
			for _, c := range langkahProteksiAneka[a].cek {
				if c.kosong(o) {
					hasil = append(hasil, PelanggaranAneka{Lokasi: i, Okupasi: okupasi, Objek: n, Medan: c.medan, Pesan: pesanProteksi(c.teks)})
				}
			}
		}
	}
	if len(aktif) == 0 {
		return nil
	}
	for i, l := range lokasi {
		for j, ok := range l.Okupasi {
			for n, o := range ok.Objek {
				periksa(i, &j, n, o)
			}
		}
		for n, o := range l.ObjekGolf {
			periksa(i, nil, n, o)
		}
	}
	return hasil
}

// String - "lokasi[i].okupasi[j].objek[n].medan: pesan" (log / galat).
func (p PelanggaranAneka) String() string {
	if p.Okupasi == nil {
		return fmt.Sprintf("lokasi[%d].objekGolf[%d].%s: %s", p.Lokasi, p.Objek, p.Medan, p.Pesan)
	}
	return fmt.Sprintf("lokasi[%d].okupasi[%d].objek[%d].%s: %s", p.Lokasi, *p.Okupasi, p.Objek, p.Medan, p.Pesan)
}
