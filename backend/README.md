# pseudogo — Konverter Pseudocode ALPRO ke Golang

CLI transpiler yang mengubah pseudocode sesuai *Spesifikasi Standar Pseudocode
ALPRO* (dokumen `rules_pseudocode.md`) menjadi kode Go yang valid dan bisa
langsung dijalankan.

## Struktur proyek

```
pseudogo/
  cmd/pseudogo/        entry point CLI
  internal/lexer/       tokenizer
  internal/ast/         definisi node AST
  internal/parser/      recursive-descent parser (token -> AST)
  internal/codegen/      AST -> kode Go (pakai go/format untuk auto-indent)
  testdata/              contoh pseudocode (termasuk 4 algoritma referensi
                          dari dokumen: sequential search, binary search,
                          selection sort, insertion sort)
  e2e/                    test otomatis: convert -> go build -> jalankan ->
                          cek output
```

## Pemakaian

```bash
go run ./cmd/pseudogo testdata/selection_sort.pseudo -o hasil.go
go run hasil.go
```

Tanpa `-o`, hasil konversi dicetak ke stdout.

## Menjalankan test

```bash
go test ./...
```

Setiap file di `testdata/*.pseudo` otomatis di-convert lalu di-`go build`
sebagai smoke test. Beberapa file (algoritma referensi + fitur umum) juga
dites hasil eksekusinya secara exact-match terhadap output yang diharapkan.

## Keputusan desain penting

Beberapa hal di dokumen aturan tidak sepenuhnya eksplisit, jadi diputuskan
sebagai berikut (didiskusikan dengan mr.ngx):

- **Indexing array**: pseudocode 1-based (`array[1..10]`, akses dari `[1]`)
  digeser otomatis jadi 0-based di Go. Setiap akses `arr[i]` menjadi
  `arr[i-lower]` (`lower` diambil dari batas bawah array itu di deklarasi;
  default 1 kalau tidak diketahui).
- **`output()`**: setiap panggilan otomatis diakhiri newline. Implementasinya
  pakai `fmt.Print(args..., "\n")` (bukan `fmt.Println`), supaya spasi antar
  argumen mengikuti aturan Go (`Print` tidak menambah spasi bila salah satu
  operand adalah string) — ini yang bikin `output("Nama: ", nama)` tercetak
  rapat tanpa spasi ganda, sesuai contoh di dokumen.
- **Parameter fungsi tanpa mode**: di dokumen, parameter `procedure` selalu
  diawali `in`/`out`/`inout`, tapi parameter `function` (mis.
  `CekGenap(angka : integer)`) tidak. Parser menganggap mode default `in`
  kalau tidak ditulis.
- **Parameter `out`/`inout` skalar** → jadi pointer di Go (`*int`, dll), dan
  setiap pemakaian variabel itu di dalam body otomatis di-dereference.
  Parameter **array** apapun mode-nya cukup `[]T` (slice Go sudah
  reference-like), tidak perlu pointer tambahan.
- **Array selalu jadi Go slice** (`[]T`), baik untuk deklarasi lokal
  (`array [1..10] of integer` → `make([]int, 10)`) maupun parameter dengan
  batas dinamis (`array [1..N] of integer` → `arr []int`).
- **`char`** dipetakan ke `rune`. Supaya `output(charVar)` mencetak
  karakternya (bukan kode Unicode-nya), codegen otomatis bungkus dengan
  `string(...)` saat tipe hasil inferensinya `char`. Untuk `input()` ke
  variabel `char`, codegen generate baca-string-lalu-ambil-rune-pertama
  (karena `fmt.Scan` ke `*rune` akan coba parse sebagai angka, bukan
  karakter).
- **`div` dan `/`** sama-sama dipetakan ke `/` di Go — karena Go otomatis
  melakukan pembagian integer kalau kedua operand bertipe `int`, jadi
  keduanya menghasilkan kode yang sama persis untuk operand integer.

## Batasan yang diketahui (untuk versi ini)

- Variabel yang dideklarasikan di `Kamus` tapi tidak pernah dipakai akan
  membuat kode Go hasil konversi gagal compile (`declared and not used`) —
  ini perilaku standar Go, sekaligus bagus untuk membiasakan penulisan
  pseudocode yang bersih.
- `input()` untuk `string` membaca satu token (dipisah spasi/newline), bukan
  satu baris penuh — cukup untuk nama tanpa spasi, tapi belum untuk kalimat.
- Belum ada validasi semantik (pengecekan tipe, variabel dipakai sebelum
  dideklarasikan, dst.) — kalau ada, errornya baru muncul dari Go compiler
  saat kode hasil konversi di-build, bukan dari pseudogo sendiri.
- Array multi-dimensi belum didukung (dokumen juga tidak mencontohkannya).

## Langkah berikutnya (belum dikerjakan)

Sesuai rencana awal: setelah CLI ini solid, bungkus jadi REST API (Go) +
frontend Angular, supaya mahasiswa bisa convert dari browser.
