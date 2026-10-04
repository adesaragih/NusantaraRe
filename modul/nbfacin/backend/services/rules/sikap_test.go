package rules

import (
	"strings"
	"testing"
)

// TestSikapKhusus - tiket NB-09: empat predikat yang sikapnya diputuskan sendiri.
// IsPKSASM (gagal keras dengan alasannya) diuji di TestEvalGagalKeras.
func TestSikapKhusus(t *testing.T) {
	periksaEval(t, []ujiEval{
		// K-002: IsOfferFacIn memakai ekspresi tersimpan .Quotation.BusinessFac = "F".
		{"IsOfferFacIn tersimpan benar", "IsOfferFacIn", kasusUji{".Quotation.BusinessFac": "F"}, true},
		// Kondisi teks tampilan (Kode Bisnis "02") TIDAK berlaku.
		{"IsOfferFacIn tampilan tidak berlaku", "IsOfferFacIn", kasusUji{".Quotation.BusinessCode": "02", "pyWorkPage.Quotation.BusinessCode": "02"}, false},
		// K-019: IsFacout menguji ProposalAcceptStatus = 4 (Banding), bukan fac out.
		{"IsFacout Banding", "IsFacout", kasusUji{"pyWorkPage.ProposalAcceptStatus": "4"}, true},
		{"IsFacout bukan Banding", "IsFacout", kasusUji{"pyWorkPage.ProposalAcceptStatus": "3"}, false},
		// K-003: IsSpreadingDepan = IsPA OR IsTravel OR IsMBU OR IsFire, salinan
		// folder Endorsment.
		{"IsSpreadingDepan lewat IsPA", "IsSpreadingDepan", kasusUji{"pyWorkPage.Quotation.BusinessType": "PA"}, true},
		{"IsSpreadingDepan lewat IsMBU", "IsSpreadingDepan", kasusUji{"pyWorkPage.OfferFacIn.QuotationData.BusinessType": "MBUCar"}, true},
		{"IsSpreadingDepan tidak cocok", "IsSpreadingDepan", kasusUji{"pyWorkPage.Quotation.BusinessType": "Aneka"}, false},
	})
}

// TestSikapKhususTercatat - setiap sikap khusus meninggalkan jejak di registry:
// asal salinan IsSpreadingDepan menyebut folder dan versinya (K-003), dan
// catatan IsOfferFacIn / IsFacout menyebut keputusannya.
func TestSikapKhususTercatat(t *testing.T) {
	asal := strings.Join(registry["ISSPREADINGDEPAN"].asal, " ")
	for _, mau := range []string{`Endorsment Fac In\When\IsSpreadingDepan.xml`, "pyRuleSetVersion 01-01-52", "K-003"} {
		if !strings.Contains(asal, mau) {
			t.Errorf("asal IsSpreadingDepan %q tidak memuat %q", asal, mau)
		}
	}
	for nama, mau := range map[string]string{"ISOFFERFACIN": "K-002", "ISFACOUT": "Banding"} {
		if c := registry[nama].catatan; !strings.Contains(c, mau) {
			t.Errorf("catatan %s %q tidak memuat %q", nama, c, mau)
		}
	}
}
