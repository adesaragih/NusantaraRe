package models

import (
	"reflect"
	"testing"
)

// Kotak saring portal (perintah work owner 07-10-2026): kata per spasi, huruf besar, paling banyak MaksKataCari.
func TestKataCari(t *testing.T) {
	if k := KataCari("  uji-pol\tMarsh  "); !reflect.DeepEqual(k, []string{"UJI-POL", "MARSH"}) {
		t.Fatalf("kata %v", k)
	}
	if k := KataCari(" "); len(k) != 0 {
		t.Fatalf("kosong %v", k)
	}
	if k := KataCari("a b c d e f g"); len(k) != MaksKataCari {
		t.Fatalf("batas %v", k)
	}
}

// CocokCari - padanan SQL portal di gudang tiruan: setiap kata memuat salah satu nilai, tanpa beda huruf.
func TestCocokCari(t *testing.T) {
	nilai := []string{"EDMT-7", "UJI-POL-1", "UJI Insured Satu", "UJI SOB"}
	for _, c := range []string{"", "edmt-7", "insured pol-1", "SOB satu"} {
		if !CocokCari(c, nilai...) {
			t.Errorf("%q harus cocok", c)
		}
	}
	for _, c := range []string{"edmt-8", "insured dua"} {
		if CocokCari(c, nilai...) {
			t.Errorf("%q tidak boleh cocok", c)
		}
	}
}
