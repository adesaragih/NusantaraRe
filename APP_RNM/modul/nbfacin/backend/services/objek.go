package services

// Tab Object FIRE - tiket 35. Baca daftar objek case; Save mengganti seluruh daftar.

import (
	"context"
	"errors"
	"fmt"
	"regexp"
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

// lebarObjek - lebar kolom (BYTE) tiap medan teks (migrasi 186 = rancangan).
var lebarObjek = []struct {
	nama  string
	nilai func(models.ObjekFire) string
	n     int
}{
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
}

// DenganObjek memasang penyimpan objek (tiket 35).
func (s *Service) DenganObjek(o repository.PenyimpanObjek) *Service {
	s.objek = o
	return s
}

// batasBarisObjek - T_LOCATIONLIST.SEQ_NO NUMBER(5): paling banyak 99999 baris (A112).
const batasBarisObjek = 99999

// periksaObjek - Object Type wajib per baris (sel 6 wajib), lantai, lebar kolom. Server TIDAK
// mencocokkan Object Type / Roof / Wall / Floor ke daftar pilihan (daftar layar G-9 sesi 0f;
// nilai disimpan teks apa adanya sesuai lebar kolom).
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
		for _, l := range lebarObjek {
			if len(l.nilai(o)) > l.n {
				masalah = append(masalah, fmt.Sprintf("baris[%d].%s paling banyak %d byte", i, l.nama, l.n))
			}
		}
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
// case (404); daftar diganti utuh di satu transaksi, lalu dibaca ulang.
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
	err := s.transaksi(ctx, func(tx *db.Tx) error { return s.objek.GantiObjek(ctx, tx, id, baris) })
	if errors.Is(err, repository.ErrKasusTidakAda) {
		return nil, ErrKasusTidakAda
	}
	if err != nil {
		return nil, err
	}
	return s.BacaObjek(ctx, id)
}
