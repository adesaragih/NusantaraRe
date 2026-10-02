// Salinan `modul/nbfacin/backend/services/rules/logika/logika.go` (01-10-2026,
// sha256 6a9da6b4…) - lihat paket predikat. Isi tidak diubah.
//
// Package logika mengurai `<pyLogic>` rule `When`: label kondisi dirangkai AND/OR,
// kurung, dan negasi. Ia terpisah dari paket rules supaya pembangkit (`../bangkit`)
// dan penilai memakai pengurai yang SAMA tanpa pembangkit bergantung pada
// `registry_gen.go` yang ia tulis sendiri.
//
// [dugaan] `!` = NOT, `&&` = AND, `||` = OR - tafsir standar ekspresi Java/Pega,
// dikuatkan nama rule `IsNotPAandNotMBU` = `!A AND !B` atas IsPA/IsMBU (keputusan
// work owner 01-10-2026, `docs/KEPUTUSAN-30-09-2026.md` butir 23). Korpus NB
// memakai `!` 4 kali, `&&` sekali, `||` tidak pernah.
package logika

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

// Simpul - satu label, atau rangkaian anak dengan SATU operator.
//
// AND dan OR tidak pernah bercampur pada satu tingkat: urutan prioritasnya di
// Pega belum terverifikasi, dan korpus NB tidak pernah memakainya (sensus
// 01-10-2026). Logika seperti `A AND B OR C` ditolak `Urai`.
type Simpul struct {
	label string
	bukan bool // `!` di depan simpul ini
	dan   bool // true: anak dirangkai AND, false: OR
	anak  []Simpul
}

// Nilai menilai simpul atas hasil tiap label.
func (s Simpul) Nilai(hasil map[string]bool) bool {
	v := s.nilaiTanpaNegasi(hasil)
	if s.bukan {
		return !v
	}
	return v
}

func (s Simpul) nilaiTanpaNegasi(hasil map[string]bool) bool {
	if s.label != "" {
		return hasil[s.label]
	}
	for _, a := range s.anak {
		v := a.Nilai(hasil)
		// AND: satu anak salah → seluruhnya salah. OR: satu anak benar →
		// seluruhnya benar. Hasil anak sudah dihitung semua oleh pemanggil;
		// berhenti di sini hanya menghemat penggabungan.
		if s.dan && !v {
			return false
		}
		if !s.dan && v {
			return true
		}
	}
	return s.dan
}

// Label - label yang dipakai logika, urut seperti Pega menamainya: A, B, …, Z,
// AA, AB, … (pendek dulu, lalu abjad).
func (s Simpul) Label() []string {
	ada := map[string]bool{}
	s.kumpulkanLabel(ada)
	hasil := make([]string, 0, len(ada))
	for l := range ada {
		hasil = append(hasil, l)
	}
	sort.Slice(hasil, func(i, j int) bool {
		if len(hasil[i]) != len(hasil[j]) {
			return len(hasil[i]) < len(hasil[j])
		}
		return hasil[i] < hasil[j]
	})
	return hasil
}

func (s Simpul) kumpulkanLabel(ada map[string]bool) {
	if s.label != "" {
		ada[s.label] = true
	}
	for _, a := range s.anak {
		a.kumpulkanLabel(ada)
	}
}

// Urai mengurai teks `<pyLogic>`.
func Urai(teks string) (Simpul, error) {
	u := &pengurai{token: token(teks)}
	s, err := u.rangkaian()
	if err != nil {
		return Simpul{}, err
	}
	if u.i != len(u.token) {
		return Simpul{}, fmt.Errorf("token sisa %q", u.token[u.i:])
	}
	return s, nil
}

// token memecah logika: kurung, `!`, `&&`, `||`, kata kunci, dan label.
func token(s string) []string {
	return strings.Fields(strings.NewReplacer("(", " ( ", ")", " ) ", "!", " ! ", "&&", " && ", "||", " || ").Replace(s))
}

// operatorBiner - token operator biner → "AND" / "OR"; kosong bila bukan.
func operatorBiner(t string) string {
	switch strings.ToUpper(t) {
	case "AND", "&&":
		return "AND"
	case "OR", "||":
		return "OR"
	}
	return ""
}

type pengurai struct {
	token []string
	i     int
}

func (u *pengurai) rangkaian() (Simpul, error) {
	pertama, err := u.satuan()
	if err != nil {
		return Simpul{}, err
	}
	anak := []Simpul{pertama}
	operator := ""
	for u.i < len(u.token) && u.token[u.i] != ")" {
		op := operatorBiner(u.token[u.i])
		if op == "" {
			return Simpul{}, fmt.Errorf("operator %q", u.token[u.i])
		}
		if operator != "" && op != operator {
			return Simpul{}, errors.New("AND dan OR bercampur tanpa kurung")
		}
		operator = op
		u.i++
		s, err := u.satuan()
		if err != nil {
			return Simpul{}, err
		}
		anak = append(anak, s)
	}
	if len(anak) == 1 {
		return pertama, nil
	}
	return Simpul{dan: operator == "AND", anak: anak}, nil
}

func (u *pengurai) satuan() (Simpul, error) {
	if u.i >= len(u.token) {
		return Simpul{}, errors.New("logika terpotong")
	}
	t := u.token[u.i]
	u.i++
	switch {
	case t == "!":
		s, err := u.satuan()
		if err != nil {
			return Simpul{}, err
		}
		s.bukan = !s.bukan
		return s, nil
	case t == "(":
		s, err := u.rangkaian()
		if err != nil {
			return Simpul{}, err
		}
		if u.i >= len(u.token) || u.token[u.i] != ")" {
			return Simpul{}, errors.New("kurung tidak ditutup")
		}
		u.i++
		return s, nil
	case t == ")" || operatorBiner(t) != "" || strings.EqualFold(t, "NOT"):
		// NOT sebagai kata tidak ada di korpus NB; ditolak, bukan ditebak.
		return Simpul{}, fmt.Errorf("token %q di tempat label", t)
	}
	return Simpul{label: t}, nil
}
