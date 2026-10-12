package services

// Bacaan paralel dan log kinerja GET / POST yang sering dibuka (laporan work owner 12-10-2026: tab Scoring "Memuat…" beberapa detik;
// "mau, jangan halu" untuk Payment / Deduction / Cedant / Submit). Prinsip: HASIL IDENTIK dengan alur berurutan lama - hanya cara
// baca; tanpa dependensi baru; tidak ada tx yang dipakai dua goroutine; galat yang dikembalikan = yang PERTAMA terjadi, dengan
// pembungkus yang sama dengan alur lama (langkah membungkusnya sendiri).

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/cockroachdb/apd/v3"
)

// ambangLogBaca - permintaan lebih lama dari ini dicatat satu baris dengan durasi tiap langkah.
const ambangLogBaca = 2 * time.Second

// langkahBaca - satu bacaan independen; f menerima ctx turunan yang dibatalkan begitu langkah lain gagal.
type langkahBaca struct {
	nama string
	f    func(context.Context) error
}

// bacaParalel - jalankan semua langkah bersamaan, tunggu semuanya; galat pertama dikembalikan, durasi tiap langkah dicatat di u.
func bacaParalel(ctx context.Context, u *ukurBaca, langkah ...langkahBaca) error {
	ctxBaca, batal := context.WithCancel(ctx)
	defer batal()
	var (
		wg      sync.WaitGroup
		mu      sync.Mutex
		pertama error
	)
	durasi := make([]time.Duration, len(langkah))
	for i, l := range langkah {
		wg.Add(1)
		go func(i int, l langkahBaca) {
			defer wg.Done()
			t0 := time.Now()
			err := l.f(ctxBaca)
			durasi[i] = time.Since(t0)
			if err != nil {
				mu.Lock()
				if pertama == nil {
					pertama = err
					batal()
				}
				mu.Unlock()
			}
		}(i, l)
	}
	wg.Wait()
	for i, l := range langkah {
		u.tambah(l.nama, durasi[i])
	}
	u.akhir = time.Now()
	return pertama
}

// ukurBaca - durasi langkah satu permintaan (log kinerja; tanpa data kasus selain ID).
type ukurBaca struct {
	mulai, akhir time.Time
	langkah      []string
}

func ukurBaru() *ukurBaca {
	n := time.Now()
	return &ukurBaca{mulai: n, akhir: n}
}

func (u *ukurBaca) tambah(nama string, d time.Duration) {
	u.langkah = append(u.langkah, nama+" "+d.String())
}

// tandai - langkah berurutan: durasi sejak tanda terakhir.
func (u *ukurBaca) tandai(nama string) {
	n := time.Now()
	u.tambah(nama, n.Sub(u.akhir))
	u.akhir = n
}

// lambat - satu baris log bila total melewati ambangLogBaca.
func (u *ukurBaca) lambat(apa, id string) {
	if total := time.Since(u.mulai); total > ambangLogBaca {
		log.Printf("nbfacin: %s kasus %s lambat %v: %s", apa, id, total, strings.Join(u.langkah, ", "))
	}
}

// kursMemo - sumberKurs dengan cache per permintaan (Scoring memanggil totalKeIDR empat kali atas mata uang yang sama). Hasil "tidak
// ada" (ID kosong / kurs nil) ikut di-cache; galat tidak.
type kursMemo struct {
	src  sumberKurs
	mu   sync.Mutex
	id   map[string]string
	kurs map[string]*apd.Decimal
}

func kursMemoBaru(src sumberKurs) *kursMemo {
	return &kursMemo{src: src, id: map[string]string{}, kurs: map[string]*apd.Decimal{}}
}

func (m *kursMemo) IDMataUang(ctx context.Context, nama string) (string, error) {
	m.mu.Lock()
	v, ada := m.id[nama]
	m.mu.Unlock()
	if ada {
		return v, nil
	}
	v, err := m.src.IDMataUang(ctx, nama)
	if err == nil {
		m.mu.Lock()
		m.id[nama] = v
		m.mu.Unlock()
	}
	return v, err
}

func (m *kursMemo) KursPada(ctx context.Context, idMataUang, waktu string) (*apd.Decimal, error) {
	kunci := idMataUang + "\x00" + waktu
	m.mu.Lock()
	v, ada := m.kurs[kunci]
	m.mu.Unlock()
	if ada {
		return v, nil
	}
	v, err := m.src.KursPada(ctx, idMataUang, waktu)
	if err == nil {
		m.mu.Lock()
		m.kurs[kunci] = v
		m.mu.Unlock()
	}
	return v, err
}
