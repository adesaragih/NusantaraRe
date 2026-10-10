package services

import (
	"reflect"
	"testing"

	"nusantarare/modul/treatyin/backend/models"
)

// `SetSpreadName` [3.2.2] - `TempSprd.TreatyGroupID` = Treaty Group Detail,
// yang `AddDelSpreadingTreatyin` [1] setel (ralat 9 Oktober 2026: [3.2.1]
// hanya membuang `pyReportContentPage`).
func TestSetSpreadNameMemakaiTreatyGroupDetail(t *testing.T) {
	src := &sumberPropUji{
		induk: []models.SusunanSpreading{{ReinsTypeID: "10263", ReinsTypeName: "2025 QS 181M TRT", TreatyYearID: "1000684"}},
		anak: map[string][]models.SusunanSpreading{"10263": {
			{ReinsTypeID: "10028", ReinsTypeName: "QS (OR)", Pct: "40"},
			{ReinsTypeID: "10004", ReinsTypeName: "QS (R/I)", Pct: "60"},
		}},
	}
	limits := pohonUji(t, `[{"TreatyType":"QUOTA SHARE","Detail":[{"TreatyGroup":"HOSPITAL","TreatyGroupID":"10015","RNMShare":"25",`+
		`"SpreadingList":[{"ReinsTypeID":"10263","Pct":"25"}],"RNMShareList":[{"Currency":"IDR","Value":"56250000"}]}]}]`)
	h := shareProp(t, MasukanShareProp{Aksi: AksiSharePropSebarNama, Limits: limits, Commencement: "20250701"}, src)
	if !reflect.DeepEqual(src.grupInduk, []string{"10015"}) {
		t.Errorf("grup RD induk %v, ingin [10015]", src.grupInduk)
	}
	r := larikSimpul(detailUji(h, 0, 0), "SpreadingList")[0]
	if n := len(larikSimpul(r, "BreakDownSprdList")); n != 2 {
		t.Errorf("rincian %d baris, ingin 2 (QS (OR), QS (R/I))", n)
	}
}
