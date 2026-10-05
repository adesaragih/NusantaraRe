// Package models memuat bentuk data modul Company Detail.
//
// Satu organisasi = satu baris `POOLDATA.CLIENT` (`FLAG` `Org`) + baris `CLIENT_PICLIST` (PIC) + baris
// `CLIENT_ADDRESS` (alamat; SATU BARIS PER NOMOR Phone and Fax, keputusan work owner 04-10-2026). Pilihan dropdown
// dari `M_ENUMERASI` (tabel modul ini) dan `NATION` (tabel datar, migrasi 802-804). Tanpa dokumen Pega
// (`M_CLIENT`) - keputusan work owner 03-10-2026: "aku tidak mau ada json lagi".
package models

// Nilai tetap baris organisasi - sama dengan baris buatan Pega SFAGIS (`RDBINSERTCLIENT`).
const (
	// FlagOrg - `CLIENT.FLAG` organisasi (DEV 04-10-2026: 18.345 baris).
	FlagOrg = "Org"
	// AwalanID - `CLIENT.ID` = AwalanID + `IDVIEW`.
	AwalanID = "ASM-SFAGIS-WORK-ORG "
	// AwalanIDView - `CLIENT.IDVIEW` = `ORG-` + nomor.
	AwalanIDView = "ORG-"
	// AwalanPIC - `CLIENT_PICLIST.USERIDENTIFIER` baris baru = `PIC-` + nomor per organisasi.
	AwalanPIC = "PIC-"
)

// Jenis `M_ENUMERASI.JENIS` yang dipakai layar (sama dengan tipe enumerasi Pega asalnya).
const (
	JenisTitle       = "title"
	JenisBidangUsaha = "bidangusaha"
	JenisPosisi      = "posisi"
	JenisAlamat      = "jenisalamat"
	JenisTelfax      = "telfax"
	JenisKodeHP      = "kodehp"
	JenisGender      = "gender"
)

// JenisPilihan - urutan jenis yang dibaca untuk pilihan form.
var JenisPilihan = []string{JenisTitle, JenisBidangUsaha, JenisPosisi, JenisAlamat, JenisTelfax, JenisKodeHP, JenisGender}

// Organisasi adalah satu baris `CLIENT` organisasi.
type Organisasi struct {
	// ID - `CLIENT.ID`, mis. `ASM-SFAGIS-WORK-ORG ORG-18351`; tidak pernah berubah.
	ID string `json:"id"`
	// IDView - `CLIENT.IDVIEW`, mis. `ORG-18351`.
	IDView string `json:"idView"`
	// Nama - Organization Name (`NAME`).
	Nama string `json:"nama"`
	// Title - LABEL title (`TITLE`, mis. `PT.`), seperti `RDBINSERTCLIENT`.
	Title string `json:"title"`
	NPWP  string `json:"npwp"`
	// Country - `COUNTRY` = `NATION.OLDID` (atau `NATION.ID` bila negara itu tanpa OLDID); CountryName = `NOTE`.
	Country     string `json:"country"`
	CountryName string `json:"countryName"`
	// BusinessField - kode `BU_ID` (`M_ENUMERASI` bidangusaha).
	BusinessField string `json:"businessField"`
	// ParentID - `PARENT_ID` (migrasi 805); ParentName - `GROUPNAME`, dibaca NB FacIn `GetGroupName_SQL`.
	ParentID   string `json:"parentId"`
	ParentName string `json:"parentName"`
	// Note - `NOTE` (migrasi 805).
	Note string `json:"note"`
	// Jejak (migrasi 805) - kosong untuk baris yang belum pernah disimpan lewat aplikasi Go.
	CreatedBy string `json:"createdBy"`
	CreatedAt string `json:"createdAt"`
	UpdatedBy string `json:"updatedBy"`
	UpdatedAt string `json:"updatedAt"`
}

// PIC adalah satu baris `CLIENT_PICLIST`.
type PIC struct {
	// UserIdentifier - kunci baris di organisasinya (`USERIDENTIFIER`); kosong = PIC baru.
	UserIdentifier string `json:"userIdentifier"`
	// Nama - Name (`NICKNAME`).
	Nama string `json:"nama"`
	// Position - teks bebas (`POSITION`); data lama Pega kebanyakan di luar daftar posisi.
	Position string `json:"position"`
	// Gender - kode `M_ENUMERASI` gender: `1` Male, `2` Female (migrasi 806).
	Gender string `json:"gender"`
	// Email - `EMAIL`.
	Email string `json:"email"`
	// DateOfBirth - `DD-MM-YYYY` di kabel; disimpan `YYYYMMDD` (bentuk tanggal Pega) di `DATEOFBIRTH`.
	DateOfBirth string `json:"dateOfBirth"`
	// Phone - Phone number (`PHONENUMBER`).
	Phone string `json:"phone"`
}

// Telfax adalah satu nomor Phone and Fax sebuah alamat - satu baris `CLIENT_ADDRESS` (migrasi 807).
type Telfax struct {
	// Jenis - kode `M_ENUMERASI` telfax (`TELFAX_TYPE`): 3 MOBILE PHONE, 5 OFFICE PHONE untuk nomor baru.
	Jenis string `json:"type"`
	// Code - kode area `M_ENUMERASI` kodehp (`TELFAX_CODE`).
	Code string `json:"code"`
	// No - nomornya (`TELFAX_NO`).
	No string `json:"no"`
}

// Alamat adalah satu alamat organisasi: satu atau lebih baris `CLIENT_ADDRESS` beralamat sama.
type Alamat struct {
	// Asal - `ASMADDRESS` alamat ini SEBELUM diubah (kosong = alamat baru). Dipakai untuk membawa kolom yang tidak
	// ada di layar (kota, kode pos, jejak Pega) ke baris yang ditulis ulang.
	Asal string `json:"asal"`
	// Jenis - kode `M_ENUMERASI` jenisalamat (`ASMADDRESSTYPE`).
	Jenis string `json:"type"`
	// Address - `ASMADDRESS`; unik di satu organisasi (kunci baris warisan).
	Address string   `json:"address"`
	Telfax  []Telfax `json:"telfax"`
}

// BarisAlamat adalah satu baris `CLIENT_ADDRESS` apa adanya.
type BarisAlamat struct {
	ClientID, Type, Address                             string
	City, CityName, ZipCode, DistrictName, ProvinceName string
	RWName, PxCreateOperator, PxCreateDateTime          string
	TelfaxType, TelfaxCode, TelfaxNo                    string
}

// Detail adalah satu organisasi lengkap.
type Detail struct {
	Organisasi
	PIC    []PIC    `json:"pic"`
	Alamat []Alamat `json:"alamat"`
}

// Isian adalah isian form Create dan ubah.
type Isian struct {
	Nama          string   `json:"nama"`
	Title         string   `json:"title"`
	NPWP          string   `json:"npwp"`
	Country       string   `json:"country"`
	BusinessField string   `json:"businessField"`
	ParentID      string   `json:"parentId"`
	Note          string   `json:"note"`
	PIC           []PIC    `json:"pic"`
	Alamat        []Alamat `json:"alamat"`
}

// Pilihan adalah satu baris `M_ENUMERASI`.
type Pilihan struct {
	Jenis string `json:"-"`
	Kode  string `json:"kode"`
	Label string `json:"label"`
	// Aktif - tampil di dropdown; yang tidak aktif hanya untuk menampilkan nilai lama.
	Aktif bool `json:"aktif"`
}

// Akun adalah satu akun login AKTIF `M_LOGIN_GO` (milik inti, dibaca saja) - pilihan PIC Name (perintah work owner
// 05-10-2026). `CLIENT_PICLIST.NICKNAME` menyimpan `Nama`-nya, seperti nama PIC lama Pega; `POSITION` menyimpan
// `Jabatan`-nya (`JOB_POSITION`).
type Akun struct {
	LoginID string `json:"loginId"`
	Nama    string `json:"nama"`
	Jabatan string `json:"jabatan"`
}

// Negara adalah satu baris `NATION`.
type Negara struct {
	ID            string `json:"id"`
	OldID         string `json:"oldId"`
	Nama          string `json:"nama"`
	NationInitial string `json:"nationInitial"`
}

// Kode - nilai `CLIENT.COUNTRY` negara ini: OLDID (bentuk seluruh data lama), atau ID bila tanpa OLDID.
func (n Negara) Kode() string {
	if n.OldID != "" {
		return n.OldID
	}
	return n.ID
}

// BarisDaftar adalah satu baris daftar organisasi.
type BarisDaftar struct {
	ID            string `json:"id"`
	IDView        string `json:"idView"`
	Nama          string `json:"nama"`
	Title         string `json:"title"`
	NPWP          string `json:"npwp"`
	CountryName   string `json:"countryName"`
	BusinessField string `json:"businessField"`
	ParentName    string `json:"parentName"`
}
