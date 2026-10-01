package backend

// Pelaku permintaan HTTP - penunda sampai tiket 07 / IAM.
//
// ⛔ BERKAS INI BUKAN AUTENTIKASI, dan tidak boleh disalahartikan begitu.
// Header HTTP dapat ditulis siapa saja yang dapat mengirim permintaan, jadi
// `X-Pelaku` tidak membuktikan apa pun tentang siapa pengirimnya.
//
// `[keputusan work owner 26-09-2026, butir ab]`. Kenapa ia tetap ada: tiket 05
// menuntut peran `ReasLifeAdmin` sekarang, sedangkan sumber peran yang
// sebenarnya adalah satu tabel (ADR-U-0030) yang belum ada dan tiket 07 yang
// menegakkannya. Tanpa penunda ini, jalur wewenang tidak dapat dibangun
// maupun diuji sampai tiket 07 selesai.
//
// Tiga pagar, dan ketiganya perlu:
//  1. mati secara bawaan - AUTH_STUB harus disetel sadar
//  2. ditolak saat memuat konfigurasi bila IS_PEGA_PROD=true
//  3. seluruh uji wewenang berjalan di seam services dengan Pelaku langsung,
//     sehingga aturannya tidak bergantung pada jalur ini sama sekali
//
// ⭐ LOGIN SUNGGUHAN (keputusan work owner 01-10-2026): pelaku hasil login
// (`inti/backend/login`, cookie sesi + `M_LOGIN_GO`) ditaruh di context
// permintaan oleh middleware login, dan `PelakuDari` MENDAHULUKANNYA. Stub
// header di bawah tetap ada untuk pengembangan, hanya bila tidak ada sesi.
//
// Dibaca sesudah: handlers.go.

import (
	"context"
	"net/http"
	"strings"
)

type kunciPelakuSesi struct{}

// DenganPelakuSesi menaruh pelaku hasil login di context permintaan - hanya
// middleware login yang memanggilnya.
func DenganPelakuSesi(ctx context.Context, p Pelaku) context.Context {
	return context.WithValue(ctx, kunciPelakuSesi{}, p)
}

type kunciAksesMenu struct{}

// DenganAksesMenu menaruh KODE menu akun hasil login (`M_LOGIN_GO_MENU`,
// Kelola User 01-10-2026) di context permintaan - hanya middleware login yang
// memanggilnya, untuk sesi yang juga mendapat pelaku.
//
// ⛔ Terpisah dari `Pelaku`: menu menjawab "layar mana yang boleh dibuka"
// (sidebar dan gerbang 403 `cmd/api`), bukan "peran apa yang dipegang". Aturan
// dagang modul tetap membaca PERAN.
func DenganAksesMenu(ctx context.Context, kode []string) context.Context {
	return context.WithValue(ctx, kunciAksesMenu{}, append([]string{}, kode...))
}

// AksesMenuDari membaca menu akun hasil login. `ada` false = permintaan tanpa
// sesi login - BUKAN akun tanpa menu (yang `ada` true dengan daftar kosong).
func AksesMenuDari(ctx context.Context) (kode []string, ada bool) {
	kode, ada = ctx.Value(kunciAksesMenu{}).([]string)
	return kode, ada
}

// PunyaMenu menjawab apakah `kode` ada di daftar menu.
func PunyaMenu(daftar []string, kode string) bool {
	for _, k := range daftar {
		if k == kode {
			return true
		}
	}
	return false
}

// pelakuDari membaca pelaku permintaan.
//
// Tanpa stub, ia mengembalikan pelaku KOSONG - bukan pelaku istimewa. Jalur
// yang memerlukan identitas karena itu DITOLAK, dan itu memang yang benar
// sampai sumber peran yang sebenarnya ada.
//
// ⛔ Keadaan stub dibawa sebagai ARGUMEN, bukan variabel paket. Variabel paket
// membuat dua Router dalam satu proses - yang biasa terjadi di test - berbagi
// satu saklar: Router(svc, true) akan menyalakan stub bagi server yang dibuat
// Router(svc, false) dan masih hidup. Test menjadi bergantung urutan, dan
// `go test -race` menandainya sebagai tulis-baca serentak.
func PelakuDari(r *http.Request, stubAktif bool) Pelaku {
	if p, ok := r.Context().Value(kunciPelakuSesi{}).(Pelaku); ok {
		return p
	}
	if !stubAktif {
		return Pelaku{}
	}
	var peran []string
	for _, p := range strings.Split(r.Header.Get("X-Peran"), ",") {
		if p = strings.TrimSpace(p); p != "" {
			peran = append(peran, p)
		}
	}
	return Pelaku{
		AkunID: strings.TrimSpace(r.Header.Get("X-Pelaku")),
		Peran:  peran,
	}
}
