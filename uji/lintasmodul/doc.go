// Package lintasmodul memuat uji INTEGRASI yang sengaja melintasi modul.
//
// Refactor bentuk B (30-09-2026). Modul di `modul/<nama>/` tidak pernah
// mengimpor modul lain - termasuk dari berkas ujinya. Sebagian kecil uji
// memang menguji PERTEMUAN dua modul di basis data yang sama (mis. Claim Life
// melahirkan kasus Komite, lalu Komite menimpa tangganya). Uji semacam itu
// tinggal di sini: paket ini bukan modul, tidak dipasang `cmd/api`, dan
// penjaga impor lintas modul mengecualikannya dengan nama.
package lintasmodul
