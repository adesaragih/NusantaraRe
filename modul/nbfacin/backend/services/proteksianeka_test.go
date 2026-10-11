package services

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/nbfacin/backend/models"
)

// TestTabelProteksiAneka - instrumen: 12 langkah urut dispatcher; 60 pesan aktif (63 literal XML - 3 ber-label `//`, riset
// bagian A.3 dihitung dua cara); setiap teks verbatim berakhiran ` \n` literal (dua karakter), dipotong pesanProteksi.
func TestTabelProteksiAneka(t *testing.T) {
	mau := []string{"ProtectionObjectHE_Act", "ProtectionObjectLandRig_Act", "ProtectionObjectContractorPlant_Act",
		"ProtectionObjectBoilerPresure_Act", "ProtectionObjectMB_Act", "ProtectionObjectAviation_Act",
		"ProtectionObjectAllRisk_Act", "ProtectionObjectBurglary_Act", "ProtectionObjectGolfInsurance_Act",
		"ProtectionObjectGrowingTree_Act", "ProtectionObjectElectronicEq_Act", "ProtectionObjectFidelity_Act"}
	var nama []string
	n := 0
	for _, l := range langkahProteksiAneka {
		nama = append(nama, l.activity)
		for _, c := range l.cek {
			n++
			if !strings.HasSuffix(c.teks, ` \n`) || strings.Contains(c.teks, "\n") {
				t.Errorf("%s %s: akhiran bukan spasi + `\\n` literal: %q", l.activity, c.medan, c.teks)
			}
		}
	}
	if !reflect.DeepEqual(nama, mau) || n != 60 {
		t.Errorf("%d pesan aktif, langkah %v", n, nama)
	}
	if p := pesanProteksi(`Tahun harus diisi \n`); p != "Tahun harus diisi" {
		t.Errorf("potong: %q", p)
	}
}

// TestProteksiObjekAneka - gerbang dan isi: HE (IsHE) lewat BusinessType; LandRig hanya bila BusinessCode PERSIS "B3" (K73-4 -
// kode 5 digit FacIn tidak membuka gerbang 2 karakter); anomali Pega dibawa (Golf Section -> medan quantity, GrowingTree
// areaHectar hanya bila NoOfTree juga kosong); objek Golf ber-okupasi nil; urutan objek lalu cek.
func TestProteksiObjekAneka(t *testing.T) {
	he := kasusBisnis{businessType: "HE", businessCode: "10999"}
	lokasi := []models.LokasiAneka{{Okupasi: []models.OkupasiAneka{{Objek: []models.ObjekAneka{
		{Year: "2020", VehicleHE: models.VehicleHE{ChassisNumber: "C1", EngineNumber: "E1", TypeName: "T", ObjectNameHE: "TBA", Brand: "7"}},
	}}}, ObjekGolf: []models.ObjekAneka{{}}}}
	p := proteksiObjekAneka("K", he, lokasi)
	nol := 0
	mau := []PelanggaranAneka{
		{Lokasi: 0, Okupasi: &nol, Objek: 0, Medan: "vehicleHe.objectNameHe", Pesan: "Nama Obyek tidak boleh TBA"},
		{Lokasi: 0, Objek: 0, Medan: "year", Pesan: "Tahun harus diisi"},
		{Lokasi: 0, Objek: 0, Medan: "vehicleHe.chassisNumber", Pesan: "No Rangka harus diisi"},
		{Lokasi: 0, Objek: 0, Medan: "vehicleHe.engineNumber", Pesan: "No Mesin harus diisi"},
		{Lokasi: 0, Objek: 0, Medan: "vehicleHe.typeName", Pesan: "Tipe harus diisi"},
		{Lokasi: 0, Objek: 0, Medan: "vehicleHe.objectNameHe", Pesan: "Nama Obyek harus diisi"},
		{Lokasi: 0, Objek: 0, Medan: "vehicleHe.brandName", Pesan: "Merk Obyek harus diisi"},
	}
	if !reflect.DeepEqual(p, mau) {
		t.Errorf("HE:\n%v\nmau\n%v", p, mau)
	}
	if p := proteksiObjekAneka("K", kasusBisnis{businessType: "LandRig", businessCode: "10999"}, lokasi); len(p) != 0 {
		t.Errorf("K73-4: kode 5 digit membuka gerbang 2 karakter: %v", p)
	}
	lr := proteksiObjekAneka("K", kasusBisnis{businessCode: "B3"}, []models.LokasiAneka{{ObjekGolf: []models.ObjekAneka{{Year: "1"}}}})
	if len(lr) != 4 || lr[0].Pesan != "Jumlah harus diisi" || lr[3].Medan != "serialNo" || lr[0].Okupasi != nil {
		t.Errorf("LandRig B3: %v", lr)
	}
	golf := proteksiObjekAneka("K", kasusBisnis{businessCode: "SW"}, []models.LokasiAneka{{ObjekGolf: []models.ObjekAneka{{ObjectName: "G"}}}})
	if len(golf) != 2 || golf[0].Medan != "quantity" || golf[0].Pesan != "Seksi / Section harus dipilih" || golf[1].Medan != "limitOfLiability" {
		t.Errorf("Golf: %v", golf)
	}
	pohon := proteksiObjekAneka("K", kasusBisnis{businessCode: "28"}, []models.LokasiAneka{{ObjekGolf: []models.ObjekAneka{
		{ObjectName: "P", Year: "2001", NoOfTree: ds("10")}}}})
	if len(pohon) != 0 {
		t.Errorf("GrowingTree: areaHectar kosong + NoOfTree terisi = pesan kosong Pega, tidak dikirim: %v", pohon)
	}
	if p := proteksiObjekAneka("K", kasusBisnis{businessCode: "28"}, []models.LokasiAneka{{ObjekGolf: []models.ObjekAneka{
		{ObjectName: "P", Year: "2001"}}}}); len(p) != 2 || p[1].Medan != "areaHectar" {
		t.Errorf("GrowingTree keduanya kosong: %v", p)
	}
	// isElectronicEquipment membandingkan BusinessType dengan " ElectronicEquipment " (bertepi spasi) - BusinessType_DeT
	// memberi "ElectronicEquipment": gerbang TIDAK terbuka (perilaku registry apa adanya; pemangkasan Pega belum terverifikasi).
	if p := proteksiObjekAneka("K", kasusBisnis{businessType: "ElectronicEquipment"}, lokasi); len(p) != 0 {
		t.Errorf("ElectronicEquipment: %v", p)
	}
	if p := proteksiObjekAneka("K", kasusBisnis{businessType: " ElectronicEquipment "}, []models.LokasiAneka{{ObjekGolf: []models.ObjekAneka{
		{ObjectName: "E", Section: "1", Quantity: "1", TSI: ds("1"), Year: "2020", LimitIdemnityDay: ds("1")}}}}); len(p) != 1 ||
		p[0].Medan != "limitIdemnityMonth" {
		t.Errorf("ElectronicEquipment bertepi spasi: %v", p)
	}
}

// TestGantiObjekAnekaProteksi - 422 sebelum transaksi (tidak ada yang ditulis); lolos -> tersimpan.
func TestGantiObjekAnekaProteksi(t *testing.T) {
	var tulis []models.LokasiAneka
	n := 0
	k := &models.Kasus{CaseID: "K", Opportunity: models.Opportunity{ClassOfBusiness: "UJI HE", GroupBusinessID: "3"}}
	svc := Baru(nil).DenganObjekAneka(anekaTiruan{ada: map[string][]models.LokasiAneka{"K": {}}, tulis: &tulis}).
		DenganTransaksi(tanpaTx).DenganKasus(kasusTiruan{ada: map[string]*models.Kasus{"K": k}}).
		DenganTableOfLimit(tolTiruan{bisnis: []models.BarisBisnisLini{{ID: "10999", OldID: "20", GroupPanel: "003"}}}).
		DenganPilihanItem(pilihanItemTiruan{panggil: &n})
	ctx, akun := context.Background(), inti.Pelaku{AkunID: "UJI"}
	_, err := svc.GantiObjekAneka(ctx, akun, "K", lokasiAnekaUji(models.ObjekAneka{ObjectName: "UJI"}), nil)
	var g GalatProteksiAneka
	if !errors.As(err, &g) || len(g.Pelanggaran) == 0 || tulis != nil {
		t.Fatalf("422: %v, tertulis %v", err, tulis)
	}
	lengkap := models.ObjekAneka{Year: "2020", VehicleHE: models.VehicleHE{ChassisNumber: "C", EngineNumber: "E", Type: "1",
		ObjectIDHE: "2", Brand: "3"}}
	l := lokasiAnekaUji(lengkap)
	l[0].ObjekGolf = []models.ObjekAneka{lengkap}
	if _, err := svc.GantiObjekAneka(ctx, akun, "K", l, nil); err != nil || len(tulis) != 1 {
		t.Errorf("lolos: %v", err)
	}
}
