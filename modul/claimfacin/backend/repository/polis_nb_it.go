package repository

// Tab Loss Record dan Correspondence tampilan polis NB (hanya akun IT, `models.PolisNBIT`): baca-saja tabel nbfacin
// ber-induk kasus NB. Kueri = kueri baca nbfacin `repository/lossrecord.go` / `korespondensi.go` tanpa saringan LINI
// (kasus sudah ditentukan polis klaim); setiap tabel dibaca lewat indeks induknya (UQ_T_LOSSRECORD_PARENT,
// IX_LOSSRECORDLIST_PARENT, IX_CORRESPONDENCE_PARENT, IX_RECIPIENTLIST_PARENT - katalog DEV 11-10-2026).

import (
	"context"
	"fmt"

	"nusantarare/modul/claimfacin/backend/models"
)

// sqlJabatanPelaku - M_LOGIN_GO.JOB_POSITION akun (nbfacin `sqlJabatanAkun`).
func sqlJabatanPelaku(login string) string {
	return fmt.Sprintf(`SELECT JOB_POSITION FROM %s WHERE LOGIN_ID = :1`, login)
}

// kolomRasioLossRecord - kolom uang / rasio T_LOSSRECORD = urutan medan `models.LossRecordNB` (nbfacin
// `kolomKepalaLossRecord`).
var kolomRasioLossRecord = []string{"PREMI_PAR", "PREMI_MB", "PREMI_PL", "PREMI_OTHER", "LOSS_RATIO1_YEAR_PERCENT",
	"LOSS_RATIO3_YEAR_PERCENT", "LOSS_RATIO5_YEAR_PERCENT", "NET_LOSS_RATIO1_YEAR", "NET_LOSS_RATIO3_YEAR",
	"NET_LOSS_RATIO5_YEAR"}

func sqlKepalaLossRecordNB(kepala string) string {
	pilih := "TO_CHAR(ID)"
	for _, k := range kolomRasioLossRecord {
		pilih += ", " + dec(k)
	}
	return fmt.Sprintf(`SELECT %s FROM %s WHERE PARENT_ID = :1`, pilih, kepala)
}

func sqlBarisLossRecordNB(baris string) string {
	return fmt.Sprintf(`SELECT UW_YEAR, CLASS_OF_BUSINESS, DATE_OF_LOSS, CAUSE_OF_LOSS, CURRENCY, %s, %s FROM %s
		WHERE PARENT_ID = :1 ORDER BY SEQ_NO, ID`, dec("LOSS_AMOUNT"), dec("LOSS_AMOUNT_IDR"), baris)
}

func sqlKorespondensiNB(kor string) string {
	return fmt.Sprintf(`SELECT TO_CHAR(ID), TANGGAL, SENDER, SENDER_MAIL, SUBJECT, MESSAGE, IS_EDITABLE FROM %s
		WHERE PARENT_ID = :1 ORDER BY SEQ_NO, ID`, kor)
}

func sqlPenerimaNB(pen, kor string) string {
	return fmt.Sprintf(`SELECT TO_CHAR(PARENT_ID), RECIPIENT, RECIPIENT_MAIL FROM %s
		WHERE PARENT_ID IN (SELECT ID FROM %s WHERE PARENT_ID = :1) ORDER BY PARENT_ID, SEQ_NO, ID`, pen, kor)
}

// JabatanPelaku - lihat models.Acuan.
func (a *Acuan) JabatanPelaku(ctx context.Context, akun string) (string, error) {
	t, err := a.q("M_LOGIN_GO")
	if err != nil {
		return "", err
	}
	v, _, err := a.satu(ctx, sqlJabatanPelaku(t), akun)
	return v, err
}

// PolisNBIT - lihat models.Acuan.
func (a *Acuan) PolisNBIT(ctx context.Context, nopolis, prodke, kasus string) (models.PolisNBIT, error) {
	out := models.PolisNBIT{LossRecord: models.LossRecordNB{Baris: []models.BarisLossRecordNB{}},
		Korespondensi: []models.KorespondensiNB{}}
	if kasus == "" {
		var err error
		if kasus, err = a.kasusPolisDatar(ctx, nopolis, prodke); err != nil || kasus == "" {
			return out, err
		}
	}
	out.Kasus = kasus
	t := map[string]string{}
	for _, n := range []string{"T_LOSSRECORD", "T_LOSSRECORDLIST", "T_CORRESPONDENCELIST", "T_RECIPIENTLIST"} {
		q, err := a.q(n)
		if err != nil {
			return out, err
		}
		t[n] = q
	}
	kepala, err := a.banyak(ctx, sqlKepalaLossRecordNB(t["T_LOSSRECORD"]), 1+len(kolomRasioLossRecord), kasus)
	if err != nil {
		return out, err
	}
	if len(kepala) > 0 {
		k, l := kepala[0], &out.LossRecord
		for i, p := range []*string{&l.PremiPar, &l.PremiMb, &l.PremiPl, &l.PremiOther, &l.Gross1, &l.Gross3, &l.Gross5,
			&l.Net1, &l.Net3, &l.Net5} {
			*p = k[1+i]
		}
		baris, err := a.banyak(ctx, sqlBarisLossRecordNB(t["T_LOSSRECORDLIST"]), 7, k[0])
		if err != nil {
			return out, err
		}
		for _, b := range baris {
			l.Baris = append(l.Baris, models.BarisLossRecordNB{UwYear: b[0], ClassOfBusiness: b[1], DateOfLoss: b[2],
				CauseOfLoss: b[3], Currency: b[4], LossAmount: b[5], LossAmountIDR: b[6]})
		}
	}
	kor, err := a.banyak(ctx, sqlKorespondensiNB(t["T_CORRESPONDENCELIST"]), 7, kasus)
	if err != nil {
		return out, err
	}
	if len(kor) == 0 {
		return out, nil
	}
	pen, err := a.banyak(ctx, sqlPenerimaNB(t["T_RECIPIENTLIST"], t["T_CORRESPONDENCELIST"]), 3, kasus)
	if err != nil {
		return out, err
	}
	penerima := map[string][]models.PenerimaNB{}
	for _, p := range pen {
		penerima[p[0]] = append(penerima[p[0]], models.PenerimaNB{Recipient: p[1], RecipientMail: p[2]})
	}
	for _, k := range kor {
		p := penerima[k[0]]
		if p == nil {
			p = []models.PenerimaNB{}
		}
		out.Korespondensi = append(out.Korespondensi, models.KorespondensiNB{Tanggal: k[1], Sender: k[2], SenderMail: k[3],
			Subject: k[4], Message: k[5], Terkirim: k[6] == models.StatusKorespondensiTerkirim, Penerima: p})
	}
	return out, nil
}
