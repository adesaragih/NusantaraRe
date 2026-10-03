package services

// Tab Object FIRE - tiket 35. Baca daftar objek case; Save mengganti seluruh daftar.

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
)

// ErrMasukanObjek - isian objek tidak sah (Object Type kosong, lantai bukan bilangan bulat
// >= 0, atau melebihi lebar kolom). 400; pesannya menyebut indeks baris.
var ErrMasukanObjek = errors.New("services: isian objek tidak sah")

// ErrObjekTanpaDatabase - tabel objek tidak terbaca. 503.
var ErrObjekTanpaDatabase = errors.New("services: basis data tidak dikonfigurasi, objek case NB tidak terbaca")

// polaLantai - SetErrorMessageFloorNumber_Act ("Floor number can't be minus"): kosong atau
// bilangan bulat >= 0 (A111: tanpa tanda, tanpa desimal).
var polaLantai = regexp.MustCompile(`^[0-9]+$`)

// polaJarak - Distance sisi Surrounding Risk (tiket 38, A129): bilangan >= 0, paling banyak dua
// desimal, titik sebagai pemisah; batas atas batasJarak. `[terverifikasi]` keempat kontrol Distance
// `Section\RiskAround.xml` (baris 2082-2140 dst.): pyMin 0, pyMax 100000, pyDecimalPlaces 2.
var polaJarak = regexp.MustCompile(`^[0-9]+(\.[0-9]{1,2})?$`)

// batasJarak - pyMax kontrol Distance.
const batasJarak = 100000

// jarakSah - kosong, atau cocok polaJarak dan <= batasJarak. Jarak bukan uang (CLAUDE.md §7):
// ParseFloat hanya untuk membandingkan batas.
func jarakSah(s string) bool {
	if s == "" {
		return true
	}
	if !polaJarak.MatchString(s) {
		return false
	}
	f, err := strconv.ParseFloat(s, 64)
	return err == nil && f <= batasJarak
}

type lebarMedan struct {
	nama  string
	nilai func(models.ObjekFire) string
	n     int
}

// sisiObjek - empat sisi Surrounding Risk, nama JSON-nya, dan cara mengambilnya.
var sisiObjek = []struct {
	nama  string
	ambil func(models.ObjekFire) models.SisiRisiko
}{
	{"front", func(o models.ObjekFire) models.SisiRisiko { return o.SurroundingRisk.Front }},
	{"left", func(o models.ObjekFire) models.SisiRisiko { return o.SurroundingRisk.Left }},
	{"back", func(o models.ObjekFire) models.SisiRisiko { return o.SurroundingRisk.Back }},
	{"right", func(o models.ObjekFire) models.SisiRisiko { return o.SurroundingRisk.Right }},
}

// lebarSekitar - lebar kolom T_SURROUNDINGRISK (migrasi 187) tiap medan surroundingRisk.
func lebarSekitar() []lebarMedan {
	var l []lebarMedan
	for _, s := range sisiObjek {
		ambil, awal := s.ambil, "surroundingRisk."+s.nama+"."
		l = append(l,
			lebarMedan{awal + "occupation", func(o models.ObjekFire) string { return ambil(o).Occupation }, 1000},
			lebarMedan{awal + "construction", func(o models.ObjekFire) string { return ambil(o).Construction }, 500},
			lebarMedan{awal + "distance", func(o models.ObjekFire) string { return ambil(o).Distance }, 50},
			lebarMedan{awal + "note", func(o models.ObjekFire) string { return ambil(o).Note }, 1000})
	}
	return append(l,
		lebarMedan{"surroundingRisk.housekeepingStatus", func(o models.ObjekFire) string { return o.SurroundingRisk.HousekeepingStatus }, 50},
		lebarMedan{"surroundingRisk.floodAreaStatus", func(o models.ObjekFire) string { return o.SurroundingRisk.FloodAreaStatus }, 50},
		lebarMedan{"surroundingRisk.floodArea", func(o models.ObjekFire) string { return o.SurroundingRisk.FloodArea }, 50},
		lebarMedan{"surroundingRisk.housekeepingRemark", func(o models.ObjekFire) string { return o.SurroundingRisk.HousekeepingRemark }, 500})
}

// lebarObjek - lebar kolom (BYTE) tiap medan teks (migrasi 186 = rancangan; 187 tiket 38).
var lebarObjek = append([]lebarMedan{
	{"objectNo", func(o models.ObjekFire) string { return o.ObjectNo }, 50},
	{"objectType", func(o models.ObjekFire) string { return o.ObjectType }, 50},
	{"objectName", func(o models.ObjekFire) string { return o.ObjectName }, 500},
	{"roadType", func(o models.ObjekFire) string { return o.RoadType }, 50},
	{"roadName", func(o models.ObjekFire) string { return o.RoadName }, 500},
	{"buildingNo", func(o models.ObjekFire) string { return o.BuildingNo }, 50},
	{"zipCode", func(o models.ObjekFire) string { return o.ZipCode }, 50},
	{"country", func(o models.ObjekFire) string { return o.Country }, 50},
	{"riskLocation", func(o models.ObjekFire) string { return o.RiskLocation }, 50},
	{"territory", func(o models.ObjekFire) string { return o.Territory }, 50},
	{"city", func(o models.ObjekFire) string { return o.City }, 50},
	{"district", func(o models.ObjekFire) string { return o.District }, 50},
	{"province", func(o models.ObjekFire) string { return o.Province }, 50},
	{"riskAddressId", func(o models.ObjekFire) string { return o.RiskAddressID }, 50},
	{"numberOfFloor", func(o models.ObjekFire) string { return o.NumberOfFloor }, 50},
	{"roofType", func(o models.ObjekFire) string { return o.RoofType }, 50},
	{"wallType", func(o models.ObjekFire) string { return o.WallType }, 50},
	{"floorType", func(o models.ObjekFire) string { return o.FloorType }, 50},
	{"partitionType", func(o models.ObjekFire) string { return o.PartitionType }, 50},
	{"supportWallType", func(o models.ObjekFire) string { return o.SupportWallType }, 50},
	{"otherType", func(o models.ObjekFire) string { return o.OtherType }, 50},
	{"ownership", func(o models.ObjekFire) string { return o.Ownership }, 50},
}, lebarSekitar()...)

// DenganObjek memasang penyimpan objek (tiket 35).
func (s *Service) DenganObjek(o repository.PenyimpanObjek) *Service {
	s.objek = o
	return s
}

// batasBarisObjek - T_LOCATIONLIST.SEQ_NO NUMBER(5): paling banyak 99999 baris (A112).
const batasBarisObjek = 99999

// periksaObjek - Object Type wajib per baris (sel 6 wajib), lantai, jarak sisi, lebar kolom.
// Server TIDAK mencocokkan Object Type / Roof / Wall / Floor / Ownership / Construction /
// Housekeeping / Flood ke daftar pilihan (daftar layar G-9 sesi 0f; nilai disimpan teks apa
// adanya sesuai lebar kolom).
func periksaObjek(baris []models.ObjekFire) error {
	var masalah []string
	if len(baris) > batasBarisObjek {
		return fmt.Errorf("%w: paling banyak %d baris", ErrMasukanObjek, batasBarisObjek)
	}
	for i, o := range baris {
		if strings.TrimSpace(o.ObjectType) == "" {
			masalah = append(masalah, fmt.Sprintf("baris[%d].objectType wajib diisi", i))
		}
		if o.NumberOfFloor != "" && !polaLantai.MatchString(o.NumberOfFloor) {
			masalah = append(masalah, fmt.Sprintf("baris[%d].numberOfFloor harus bilangan bulat >= 0", i))
		}
		for _, s := range sisiObjek {
			if !jarakSah(s.ambil(o).Distance) {
				masalah = append(masalah, fmt.Sprintf("baris[%d].surroundingRisk.%s.distance harus angka 0..%d dengan paling banyak 2 desimal", i, s.nama, batasJarak))
			}
		}
		for _, l := range lebarObjek {
			if len(l.nilai(o)) > l.n {
				masalah = append(masalah, fmt.Sprintf("baris[%d].%s paling banyak %d byte", i, l.nama, l.n))
			}
		}
		masalah = append(masalah, periksaItem(i, o.Items)...)
		masalah = append(masalah, periksaOkupasi(i, o.Occupations)...)
		masalah = append(masalah, periksaFEA(i, o.FEA)...)
	}
	if len(masalah) > 0 {
		return fmt.Errorf("%w: %s", ErrMasukanObjek, strings.Join(masalah, "; "))
	}
	return nil
}

// BacaObjek - daftar objek case `id` urut SEQ_NO (selalu larik).
func (s *Service) BacaObjek(ctx context.Context, id string) ([]models.ObjekFire, error) {
	if s.objek == nil {
		return nil, ErrObjekTanpaDatabase
	}
	if !idKasusSah(id) {
		return nil, ErrKasusTidakAda
	}
	baris, err := s.objek.BacaObjek(ctx, id)
	if errors.Is(err, repository.ErrKasusTidakAda) {
		return nil, ErrKasusTidakAda
	}
	return baris, err
}

// GantiObjek - tombol Save tab Object: identitas (401) -> isian (400) -> basis data (503) ->
// mata uang item di CURRENCY (400, tiket 39) -> case (404); daftar diganti utuh di satu transaksi, lalu dibaca ulang.
func (s *Service) GantiObjek(ctx context.Context, pelaku inti.Pelaku, id string, baris []models.ObjekFire) ([]models.ObjekFire, error) {
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return nil, err
	}
	if err := periksaObjek(baris); err != nil {
		return nil, err
	}
	if s.objek == nil || s.transaksi == nil {
		return nil, ErrObjekTanpaDatabase
	}
	if !idKasusSah(id) {
		return nil, ErrKasusTidakAda
	}
	if err := s.periksaMataUang(ctx, baris); err != nil {
		return nil, err
	}
	err := s.transaksi(ctx, func(tx *db.Tx) error { return s.objek.GantiObjek(ctx, tx, id, baris) })
	if errors.Is(err, repository.ErrKasusTidakAda) {
		return nil, ErrKasusTidakAda
	}
	if err != nil {
		return nil, err
	}
	return s.BacaObjek(ctx, id)
}
