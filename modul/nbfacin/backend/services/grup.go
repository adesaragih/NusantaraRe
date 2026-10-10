package services

// Penanda case "group" (tiket 63, keputusan work owner 08-10-2026 D5: "saya ada kirim di folder, when IsGroup. itu penanda
// yang masuk group"). `[terverifikasi]` D:\migrasi\RNM\DDL\IsGroup.xml (ASM-FW-GISFW-WORK!ISGROUP, pyLogic "A OR B OR C OR D
// OR E"):
//   - A `.OfferFacIn.QuotationData.IsGroup = "Group"`. Properti itu TIDAK disimpan sistem baru; penyetelnya
//     Activity\InputOfferFacInEngineer_preACT.xml langkah 23-25: "Group" bila `@equalsIgnoreCase(pyWorkPage.pxCreateOperator,
//     <login>)` untuk salah satu dari 3 login literal, selain itu "NonGroup" ("SemiGroup" -> "NonGroup"). Diturunkan saat itu
//     juga dari T_WORK_POLIS.CREATE_OP, tidak peka huruf (tanpa migrasi).
//   - B / C / D `OperatorID.pyUserIdentifier = <login>` - akun yang MEMPROSES. Ketiga login literal = himpunan yang SAMA
//     dengan langkah 23 (dibandingkan tanpa mencetak nilainya: 3 = 3, irisan 3). Perbandingan Pega "=" `[dugaan]` peka huruf;
//     di sini tidak peka huruf `[keputusan agent]` (login akun sistem baru bisa berbeda kapitalisasi).
//   - E `pyWorkPage.OfferFacIn.QuotationData.MarketingCode = <kode>` (1 literal) - MarketingCode = MARKETINGOFFICER.CLIENTID
//     (RDBList GetMOID_Sql `clientid = MarketingCode`) lewat MOID case; dibandingkan persis (dipangkas).
// ⛔ Literal login (nama akun orang) dan kode marketing TIDAK ditulis di kode / tes / dokumen (CLAUDE.md §4.10): dibaca dari
// NBFACIN_GROUP_LOGIN / NBFACIN_GROUP_MARKETINGCODE (inti/backend/config), diisi work owner / admin dari IsGroup.xml.

import (
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/nbfacin/backend/models"
)

// KonfigGrup - NBFACIN_GROUP_LOGIN / NBFACIN_GROUP_MARKETINGCODE (dirakit modul.go). Kosong = tidak ada akun / kode group.
// BondingMarketingName - NBFACIN_BONDING_MARKETINGNAME (When IsTBonding, tiket 64).
type KonfigGrup struct {
	Login, MarketingCode, BondingMarketingName []string
}

// DenganGrup memasang penanda group dan gerbang produksi (IS_PEGA_PROD) - tiket 63.
func (s *Service) DenganGrup(g KonfigGrup, produksi bool) *Service {
	s.grup, s.produksi = g, produksi
	return s
}

// loginGroup - login itu termasuk daftar group (tidak peka huruf, dipangkas); login kosong tidak pernah.
func (g KonfigGrup) loginGroup(login string) bool {
	login = strings.TrimSpace(login)
	if login == "" {
		return false
	}
	for _, l := range g.Login {
		if strings.EqualFold(l, login) {
			return true
		}
	}
	return false
}

// pembuatGroup - A: QuotationData.IsGroup == "Group" (langkah 23 preACT atas pembuat case).
func (g KonfigGrup) pembuatGroup(k models.KonteksSubmit) bool { return g.loginGroup(k.CreateOp) }

// kasusGroup - When IsGroup: A (pembuat case) OR B-D (akun yang memproses) OR E (kode marketing case).
func (g KonfigGrup) kasusGroup(k models.KonteksSubmit, pelaku inti.Pelaku) bool {
	if g.pembuatGroup(k) || g.loginGroup(pelaku.AkunID) {
		return true
	}
	kode := strings.TrimSpace(k.MarketingCode)
	for _, m := range g.MarketingCode {
		if kode != "" && m == kode {
			return true
		}
	}
	return false
}
